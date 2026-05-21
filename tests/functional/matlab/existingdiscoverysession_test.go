// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/sessiondetails"
	"github.com/stretchr/testify/suite"
)

type ExistingDiscoverySessionTestSuite struct {
	ExistingMATLABSessionTestSuite
}

func TestExistingDiscoverySessionTestSuite(t *testing.T) {
	suite.Run(t, new(ExistingDiscoverySessionTestSuite))
}

func (s *ExistingDiscoverySessionTestSuite) TestSingleSession_ConnectsSuccessfully() {
	ctx := s.T().Context()

	mockSession := s.startDiscoverableMockMATLAB()

	session := s.createSession("--matlab-session-discovery-timeout=10s")
	defer s.CleanupSession(session, true)

	output, err := session.EvaluateCode(ctx, "disp('discovered')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('discovered')")

	s.assertMockReceivedEval(mockSession, "disp('discovered')")
}

func (s *ExistingDiscoverySessionTestSuite) TestSessionDies_RediscoversNewSession() {
	ctx := s.T().Context()

	mockSessionA := s.startDiscoverableMockMATLAB()

	session := s.createSession("--matlab-session-discovery-timeout=15s")
	defer s.CleanupSession(session, false)

	output, err := session.EvaluateCode(ctx, "disp('session A')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('session A')")

	s.assertMockReceivedEval(mockSessionA, "disp('session A')")

	s.Require().NoError(mockSessionA.Stop())

	mockSessionB := s.startDiscoverableMockMATLAB()

	output, err = session.EvaluateCode(ctx, "disp('session B')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('session B')")

	s.assertMockReceivedEval(mockSessionB, "disp('session B')")
}

func (s *ExistingDiscoverySessionTestSuite) TestNoSessionBeforeTimeout_ReturnsError() {
	ctx := s.T().Context()

	session := s.createSession("--matlab-session-discovery-timeout=2s")
	defer s.CleanupSession(session, false)

	startTime := time.Now()
	_, err := session.EvaluateCode(ctx, "disp('should fail')", s.T().TempDir())
	s.assertTimeoutBoundedFailure(startTime, 2*time.Second, err)
}

func (s *ExistingDiscoverySessionTestSuite) TestLateSession_ConnectsSuccessfully() {
	ctx := s.T().Context()

	session := s.createSession("--matlab-session-discovery-timeout=10s")
	defer s.CleanupSession(session, true)

	type evalResult struct {
		output string
		err    error
	}
	resultCh := make(chan evalResult, 1)
	go func() {
		output, err := session.EvaluateCode(ctx, "disp('late discover')", s.T().TempDir())
		resultCh <- evalResult{output, err}
	}()

	s.waitForServerLog(session, "Discovering existing MATLAB sessions to attach to", 5*time.Second)

	mockSession := s.startDiscoverableMockMATLAB()

	select {
	case result := <-resultCh:
		s.Require().NoError(result.err)
		s.Contains(result.output, "disp('late discover')")
	case <-time.After(10 * time.Second):
		s.FailNow("evaluation did not complete after session details became discoverable")
	}

	s.assertMockReceivedEval(mockSession, "disp('late discover')")
}

func (s *ExistingDiscoverySessionTestSuite) TestSessionDiesWithoutReplacement_ReturnsError() {
	ctx := s.T().Context()

	mockSession := s.startDiscoverableMockMATLAB()

	session := s.createSession("--matlab-session-discovery-timeout=3s")
	defer s.CleanupSession(session, false)

	output, err := session.EvaluateCode(ctx, "disp('session A')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('session A')")

	s.assertMockReceivedEval(mockSession, "disp('session A')")

	s.Require().NoError(mockSession.Stop())
	s.Require().NoError(sessiondetails.Remove(s.testHome))

	startTime := time.Now()
	_, err = session.EvaluateCode(ctx, "disp('session should be gone')", s.T().TempDir())
	s.assertTimeoutBoundedFailure(startTime, 3*time.Second, err)
}

// --- Session Helpers ---

func (s *ExistingDiscoverySessionTestSuite) startDiscoverableMockMATLAB() *mockmatlab.Session {
	s.T().Helper()
	matlabSession := s.startMockMATLAB()
	_, err := matlabSession.ShareMATLABSession(s.testHome)
	s.Require().NoError(err, "should share MATLAB session for discovery")
	return matlabSession
}
