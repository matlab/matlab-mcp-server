// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/tests/testutils/mockmatlab"
	"github.com/matlab/matlab-mcp-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/stretchr/testify/suite"
)

const expectedConnectTitleEvals = 3

type LazyLoadTestSuite struct {
	MockMATLABTestSuite
}

func TestLazyLoadTestSuite(t *testing.T) {
	suite.Run(t, new(LazyLoadTestSuite))
}

func (s *LazyLoadTestSuite) TestLazyLoad_MATLABStartsOnFirstToolCall() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	result, err := session.ListTools(ctx, nil)
	s.Require().NoError(err, "should list tools")
	s.NotEmpty(result.Tools, "should have tools registered")

	instanceEvents, err := session.ReadInstanceEvents()
	s.Require().NoError(err)
	s.Empty(instanceEvents, "MATLAB should not have started before a tool call")

	output, err := session.EvaluateCode(ctx, "disp('hello')", s.T().TempDir())
	s.Require().NoError(err, "first tool call should trigger MATLAB start")
	s.Contains(output, "disp('hello')")

	instanceEvents = s.WaitForInstanceEvents(session, 30*time.Second,
		func(events []mockruntime.InstanceEvents) bool {
			return len(events) == 1 && countTitleEvals(events[0]) == expectedConnectTitleEvals
		},
		"desktop connect should check the release to pick the title API, then read the title and write it back once each")

	s.Require().Len(instanceEvents, 1, "MATLAB should have started after first tool call")
	s.Equal("happy", instanceEvents[0].StartedMode())
	s.Equal(1, instanceEvents[0].CountEvent(mockruntime.EventStarted), "should have exactly one started event")
	s.True(instanceEvents[0].HasFeval("cd"), "server should have sent a cd feval to set the working directory")
	s.True(instanceEvents[0].HasEval("disp('hello')"), "should have recorded the user eval")
	s.Equal(
		[]string{mockruntime.EventStarted, mockruntime.EventFeval, mockruntime.EventFeval, mockruntime.EventEval},
		eventTypesExcludingInfra(instanceEvents[0]),
		"excluding infra evals (title + correlation ID), should have exactly: started, greeting feval, cd feval, user eval",
	)
	s.Equal(expectedConnectTitleEvals, countTitleEvals(instanceEvents[0]),
		"desktop connect should check the release to pick the title API, then read the title and write it back once each")
}

func (s *LazyLoadTestSuite) TestEagerLoad_MATLABStartsOnSessionCreation() {
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil, "--initialize-matlab-on-startup")
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	instanceEvents := s.WaitForInstanceEvents(session, 30*time.Second,
		func(events []mockruntime.InstanceEvents) bool {
			return len(events) == 1 && events[0].HasEvent(mockruntime.EventStarted)
		},
		"MATLAB should have started before any tool call")

	s.Equal("happy", instanceEvents[0].StartedMode())
	s.False(instanceEvents[0].HasFeval("cd"), "no cd feval should have happened without a tool call")
	s.Zero(instanceEvents[0].CountEvent(mockruntime.EventEval),
		"no evals should happen before a tool call; the greeting and desktop connection title ops are fevals")
}

func (s *LazyLoadTestSuite) TestLazyLoad_SecondToolCallReusesSession() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	output, err := session.EvaluateCode(ctx, "disp('first')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('first')")

	output, err = session.EvaluateCode(ctx, "disp('second')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('second')")

	instanceEvents, err := session.ReadInstanceEvents()
	s.Require().NoError(err)
	s.Require().Len(instanceEvents, 1, "only one mock MATLAB instance should have been created")
	s.True(instanceEvents[0].HasEvalsInOrder("disp('first')", "disp('second')"), "should have recorded both user evals in order")
	s.Equal(2, instanceEvents[0].CountEvent(mockruntime.EventEval)-countCorrelationIDEvals(instanceEvents[0]),
		"should have exactly two user evals (excluding correlation ID evals)")
}

func (s *LazyLoadTestSuite) TestReconnection_AfterExit_MATLABRestartsOnNextToolCall() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil)
	s.Require().NoError(err)
	defer s.CleanupSession(session, false)

	output, err := session.EvaluateCode(ctx, "disp('before exit')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('before exit')")

	_, _ = session.EvaluateCode(ctx, "exit()", s.T().TempDir())

	output, err = session.EvaluateCode(ctx, "disp('after reconnect')", s.T().TempDir())
	s.Require().NoError(err, "tool call after exit() should succeed via reconnection")
	s.Contains(output, "disp('after reconnect')")

	instanceEvents, err := session.ReadInstanceEvents()
	s.Require().NoError(err)
	s.Require().Len(instanceEvents, 2, "should have started two different mock MATLAB instances")

	s.Equal(1, instanceEvents[0].ID)
	s.Equal("happy", instanceEvents[0].StartedMode())
	s.True(instanceEvents[0].HasFeval("cd"), "first instance should have a cd feval")
	s.True(instanceEvents[0].HasEvalsInOrder("disp('before exit')", "exit()"), "first instance should have eval then exit in order")
	s.True(instanceEvents[0].HasEvent(mockruntime.EventExitRequested), "first instance should have recorded exit_requested")

	s.Equal(2, instanceEvents[1].ID)
	s.Equal("happy", instanceEvents[1].StartedMode())
	s.True(instanceEvents[1].HasFeval("cd"), "second instance should have a cd feval")
	s.True(instanceEvents[1].HasEval("disp('after reconnect')"), "second instance should have recorded eval after reconnect")
	s.False(instanceEvents[1].HasEvent(mockruntime.EventExitRequested), "second instance should not have an exit event")
}

func countCorrelationIDEvals(ie mockruntime.InstanceEvents) int {
	n := 0
	for _, e := range ie.Events {
		if isCorrelationIDEval(e) {
			n++
		}
	}
	return n
}

func isCorrelationIDEval(e mockruntime.Event) bool {
	return e.Type == mockruntime.EventEval && strings.Contains(e.Code, "serviceprocess.getClientId")
}

func isTitleEvent(e mockruntime.Event) bool {
	if e.Type != mockruntime.EventFeval {
		return false
	}
	return e.Function == "eval" || e.Function == "isMATLABReleaseOlderThan"
}

func isInfraEvent(e mockruntime.Event) bool {
	return isTitleEvent(e) || isCorrelationIDEval(e)
}

func countTitleEvals(ie mockruntime.InstanceEvents) int {
	n := 0
	for _, e := range ie.Events {
		if isTitleEvent(e) {
			n++
		}
	}
	return n
}

func eventTypesExcludingInfra(ie mockruntime.InstanceEvents) []string {
	var types []string
	for _, e := range ie.Events {
		if isInfraEvent(e) {
			continue
		}
		types = append(types, e.Type)
	}
	return types
}
