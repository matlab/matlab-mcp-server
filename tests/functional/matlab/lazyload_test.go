// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/stretchr/testify/suite"
)

type LazyLoadTestSuite struct {
	MockMATLABTestSuite
}

func TestLazyLoadTestSuite(t *testing.T) {
	suite.Run(t, new(LazyLoadTestSuite))
}

func (s *LazyLoadTestSuite) TestLazyLoad_MATLABStartsOnFirstToolCall() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig())
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

	instanceEvents, err = session.ReadInstanceEvents()
	s.Require().NoError(err)
	s.Require().Len(instanceEvents, 1, "MATLAB should have started after first tool call")
	s.Equal("happy", instanceEvents[0].StartedMode())
	s.Equal(1, instanceEvents[0].CountEvent(mockruntime.EventStarted), "should have exactly one started event")
	s.True(instanceEvents[0].HasEvalMatching(isCdEval), "server should have sent a cd() eval to set the working directory")
	s.True(instanceEvents[0].HasEval("disp('hello')"), "should have recorded the user eval")
	s.Equal(
		[]string{mockruntime.EventStarted, mockruntime.EventEval, mockruntime.EventEval},
		instanceEvents[0].EventTypes(),
		"should have exactly: started, cd() eval, user eval",
	)
}

func (s *LazyLoadTestSuite) TestEagerLoad_MATLABStartsOnSessionCreation() {
	session, err := s.CreateSession(mockmatlab.HappyConfig(), "--initialize-matlab-on-startup")
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	s.Require().Eventually(func() bool {
		instanceEvents, err := session.ReadInstanceEvents()
		if err != nil {
			return false
		}
		if len(instanceEvents) != 1 {
			return false
		}
		return instanceEvents[0].HasEvent(mockruntime.EventStarted)
	}, 30*time.Second, 500*time.Millisecond,
		"MATLAB should have started before any tool call")

	instanceEvents, err := session.ReadInstanceEvents()
	s.Require().NoError(err)
	s.Equal("happy", instanceEvents[0].StartedMode())
	s.False(instanceEvents[0].HasEvent(mockruntime.EventEval), "no eval should have happened without a tool call")
}

func (s *LazyLoadTestSuite) TestLazyLoad_SecondToolCallReusesSession() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig())
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
	s.Equal(2, instanceEvents[0].CountEvent(mockruntime.EventEval)-countCdEvals(instanceEvents[0]),
		"should have exactly two user evals (excluding cd)")
}

func (s *LazyLoadTestSuite) TestReconnection_AfterExit_MATLABRestartsOnNextToolCall() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig())
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
	s.True(instanceEvents[0].HasEvalMatching(isCdEval), "first instance should have a cd() eval")
	s.True(instanceEvents[0].HasEvalsInOrder("disp('before exit')", "exit()"), "first instance should have eval then exit in order")
	s.True(instanceEvents[0].HasEvent(mockruntime.EventExitRequested), "first instance should have recorded exit_requested")

	s.Equal(2, instanceEvents[1].ID)
	s.Equal("happy", instanceEvents[1].StartedMode())
	s.True(instanceEvents[1].HasEvalMatching(isCdEval), "second instance should have a cd() eval")
	s.True(instanceEvents[1].HasEval("disp('after reconnect')"), "second instance should have recorded eval after reconnect")
	s.False(instanceEvents[1].HasEvent(mockruntime.EventExitRequested), "second instance should not have an exit event")
}

func isCdEval(code string) bool {
	return strings.HasPrefix(code, "cd('") && strings.HasSuffix(code, "')")
}

func countCdEvals(ie mockruntime.InstanceEvents) int {
	n := 0
	for _, e := range ie.Events {
		if e.Type == mockruntime.EventEval && isCdEval(e.Code) {
			n++
		}
	}
	return n
}
