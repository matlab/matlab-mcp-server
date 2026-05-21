// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/pathcontrol"
)

// ExistingMATLABSessionTestSuite tests connecting to an already-running MATLAB
// session. MATLAB is explicitly absent from PATH to prove the feature works
// without a local installation.
type ExistingMATLABSessionTestSuite struct {
	MockMATLABBaseSuite

	testHome   string
	defaultEnv []string
}

func (s *ExistingMATLABSessionTestSuite) SetupTest() {
	fakeHome := s.T().TempDir()
	s.testHome = fakeHome

	path := pathcontrol.RemoveAllMATLABsFromPath(os.Getenv("PATH"))
	env := pathcontrol.UpdateEnvEntry(os.Environ(), "PATH", path)
	env = pathcontrol.UpdateEnvEntry(env, "HOME", fakeHome)
	env = pathcontrol.UpdateEnvEntry(env, "APPDATA", filepath.Join(fakeHome, "AppData", "Roaming"))
	env = pathcontrol.UpdateEnvEntry(env, "MW_MCP_SERVER_MATLAB_SESSION_CONNECTION_TIMEOUT", "2s")
	env = pathcontrol.UpdateEnvEntry(env, "MW_MCP_SERVER_MATLAB_SESSION_DISCOVERY_TIMEOUT", "3s")
	s.defaultEnv = env
}

// --- Session Helpers ---

func (s *ExistingMATLABSessionTestSuite) createSession(extraArgs ...string) *mcpclient.LoggedSession {
	s.T().Helper()
	args := append([]string{"--matlab-session-mode=existing"}, extraArgs...)
	session, err := s.CreateSessionWithEnv(s.defaultEnv, args...)
	s.Require().NoError(err)
	return session
}

func (s *ExistingMATLABSessionTestSuite) startMockMATLAB() *mockmatlab.Session {
	s.T().Helper()
	ctx := s.T().Context()

	matlabSession, err := mockmatlab.StartSession(ctx, s.installation, mockmatlab.HappyConfig())
	s.Require().NoError(err, "should start mock MATLAB")
	s.T().Cleanup(func() { _ = matlabSession.Stop() })

	_, err = matlabSession.WaitForReady(ctx)
	s.Require().NoError(err, "mock MATLAB should become ready")

	return matlabSession
}

// --- Polling Helpers ---

func (s *ExistingMATLABSessionTestSuite) waitForServerLog(session *mcpclient.LoggedSession, substr string, timeout time.Duration) {
	s.T().Helper()
	s.Require().Eventually(func() bool {
		serverLogs, err := session.ReadServerLogs()
		return err == nil && strings.Contains(serverLogs, substr)
	}, timeout, 50*time.Millisecond, "server log did not contain %q within %s", substr, timeout)
}

// --- Assertions ---

func (s *ExistingMATLABSessionTestSuite) assertTimeoutBoundedFailure(startTime time.Time, configuredTimeout time.Duration, err error) {
	s.T().Helper()
	s.Require().Error(err)

	elapsed := time.Since(startTime)
	s.GreaterOrEqual(elapsed, configuredTimeout,
		"failure should not return before the configured timeout (immediate failure indicates an unrelated error)")
	s.LessOrEqual(elapsed, configuredTimeout+4*time.Second,
		"failure should happen near the configured timeout, not hang indefinitely")
	s.assertUserActionableError(err)
}

func (s *ExistingMATLABSessionTestSuite) assertMockReceivedEval(mockSession *mockmatlab.Session, codeSubstr string) {
	s.T().Helper()
	evals, err := mockSession.ReceivedEvals()
	s.Require().NoError(err, "should read mock MATLAB request log")
	s.Require().NotEmpty(evals, "mock MATLAB should have received at least one eval request")

	var found bool
	for _, eval := range evals {
		if strings.Contains(eval.Code, codeSubstr) {
			found = true
			break
		}
	}
	s.True(found, "mock MATLAB should have received an eval containing %q", codeSubstr)
}

func (s *ExistingMATLABSessionTestSuite) assertMockReceivedNoRequests(mockSession *mockmatlab.Session) {
	s.T().Helper()
	evals, err := mockSession.ReceivedEvals()
	s.Require().NoError(err, "should read mock MATLAB request log")
	s.Empty(evals, "mock MATLAB should not have received any requests")
}

func (s *ExistingMATLABSessionTestSuite) assertUserActionableError(err error) {
	s.T().Helper()
	s.Require().Error(err)

	message := strings.TrimSpace(err.Error())
	lowerMessage := strings.ToLower(message)

	s.NotEmpty(message)
	s.NotContains(lowerMessage, "stack trace")
	s.NotContains(lowerMessage, "panic:")
	s.NotContains(lowerMessage, "goroutine")
}
