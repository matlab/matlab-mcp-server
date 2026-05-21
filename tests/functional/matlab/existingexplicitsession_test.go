// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/pathcontrol"
	"github.com/stretchr/testify/suite"
)

type ExistingExplicitSessionTestSuite struct {
	ExistingMATLABSessionTestSuite
}

func TestExistingExplicitSessionTestSuite(t *testing.T) {
	suite.Run(t, new(ExistingExplicitSessionTestSuite))
}

func (s *ExistingExplicitSessionTestSuite) TestHappyPath_EvaluateCode_ReturnsOutput() {
	ctx := s.T().Context()

	mockSession := s.startMockMATLAB()
	session := s.startExplicitSession(mockSession)
	defer s.CleanupSession(session, true)

	output, err := session.EvaluateCode(ctx, "disp('hello')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('hello')")

	s.assertMockReceivedEval(mockSession, "disp('hello')")
	s.assertSecretNotInServerLogs(session, mockSession.APIKey)
}

func (s *ExistingExplicitSessionTestSuite) TestTimeout_ReturnsErrorWithoutDiscoveryFallback() {
	ctx := s.T().Context()

	mockSession := s.startMockMATLAB()

	details := map[string]any{
		"port":        1,
		"certificate": mockSession.CertificatePath(),
		"apiKey":      mockSession.APIKey,
	}
	detailsBytes, err := json.Marshal(details)
	s.Require().NoError(err)

	session := s.createSession(
		"--matlab-session-connection-details="+string(detailsBytes),
		"--matlab-session-connection-timeout=1s",
	)
	defer s.CleanupSession(session, false)

	startTime := time.Now()
	_, err = session.EvaluateCode(ctx, "disp('should fail')", s.T().TempDir())
	s.assertTimeoutBoundedFailure(startTime, 1*time.Second, err)

	s.assertMockReceivedNoRequests(mockSession)
}

func (s *ExistingExplicitSessionTestSuite) TestMCPCapabilities_MatchLocalMode() {
	ctx := s.T().Context()

	mockSession := s.startMockMATLAB()
	session := s.startExplicitSession(mockSession)
	defer s.CleanupSession(session, true)

	localSession := s.createLocalSession()
	defer s.CleanupSession(localSession, true)

	existingToolsResult, err := session.ListTools(ctx, nil)
	s.Require().NoError(err)
	s.Require().NotNil(existingToolsResult)
	localToolsResult, err := localSession.ListTools(ctx, nil)
	s.Require().NoError(err)
	s.Require().NotNil(localToolsResult)

	existingToolNames := make([]string, 0, len(existingToolsResult.Tools))
	for _, tool := range existingToolsResult.Tools {
		existingToolNames = append(existingToolNames, tool.Name)
	}
	localToolNames := make([]string, 0, len(localToolsResult.Tools))
	for _, tool := range localToolsResult.Tools {
		localToolNames = append(localToolNames, tool.Name)
	}
	s.ElementsMatch(localToolNames, existingToolNames, "existing mode should expose the same tools as local mode")

	existingResourcesResult, err := session.ListResources(ctx, nil)
	s.Require().NoError(err)
	s.Require().NotNil(existingResourcesResult)
	localResourcesResult, err := localSession.ListResources(ctx, nil)
	s.Require().NoError(err)
	s.Require().NotNil(localResourcesResult)

	existingResourceURIs := make([]string, 0, len(existingResourcesResult.Resources))
	for _, resource := range existingResourcesResult.Resources {
		existingResourceURIs = append(existingResourceURIs, resource.URI)
	}
	localResourceURIs := make([]string, 0, len(localResourcesResult.Resources))
	for _, resource := range localResourcesResult.Resources {
		localResourceURIs = append(localResourceURIs, resource.URI)
	}
	s.ElementsMatch(localResourceURIs, existingResourceURIs, "existing mode should expose the same resources as local mode")

	existingInitResult := session.InitializeResult()
	s.Require().NotNil(existingInitResult, "existing-mode session should have an InitializeResult")
	localInitResult := localSession.InitializeResult()
	s.Require().NotNil(localInitResult, "local-mode session should have an InitializeResult")

	s.Equal(localInitResult.ServerInfo.Name, existingInitResult.ServerInfo.Name,
		"existing mode should report the same server name as local mode")
	s.Equal(localInitResult.ServerInfo.Version, existingInitResult.ServerInfo.Version,
		"existing mode should report the same server version as local mode")
}

func (s *ExistingExplicitSessionTestSuite) TestSessionDies_ReturnsErrorWithoutDiscoveryFallback() {
	ctx := s.T().Context()

	mockSession := s.startMockMATLAB()
	session := s.startExplicitSession(mockSession)
	defer s.CleanupSession(session, false)

	output, err := session.EvaluateCode(ctx, "disp('alive')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('alive')")

	s.assertMockReceivedEval(mockSession, "disp('alive')")

	s.Require().NoError(mockSession.Stop())

	_, err = session.EvaluateCode(ctx, "disp('dead')", s.T().TempDir())
	s.Require().Error(err, "evaluation should fail after session dies with explicit details")
}

func (s *ExistingExplicitSessionTestSuite) TestCloseSession_DoesNotStopExternalMATLAB() {
	ctx := s.T().Context()

	mockSession := s.startMockMATLAB()
	session := s.startExplicitSession(mockSession)

	output, err := session.EvaluateCode(ctx, "disp('before close')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('before close')")

	s.Require().NoError(session.Close())
	s.AssertNoErrorLogs(session)
	session.DumpLogsOnFailure(s.T())

	detailsJSON, err := mockSession.ToSessionDetailsJSON()
	s.Require().NoError(err)
	session2 := s.createSession("--matlab-session-connection-details=" + detailsJSON)
	defer s.CleanupSession(session2, true)

	output, err = session2.EvaluateCode(ctx, "disp('after close')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('after close')")
}

// --- Session Helpers ---

func (s *ExistingExplicitSessionTestSuite) startExplicitSession(mockSession *mockmatlab.Session) *mcpclient.LoggedSession {
	s.T().Helper()
	detailsJSON, err := mockSession.ToSessionDetailsJSON()
	s.Require().NoError(err, "should serialize session details")
	return s.createSession("--matlab-session-connection-details=" + detailsJSON)
}

func (s *ExistingExplicitSessionTestSuite) createLocalSession() *mcpclient.LoggedSession {
	s.T().Helper()
	mockMATLABBinDir := filepath.Join(s.installation.MATLABRoot, "bin")
	path := pathcontrol.AddToPath(os.Getenv("PATH"), []string{mockMATLABBinDir})
	env := pathcontrol.UpdateEnvEntry(os.Environ(), "PATH", path)
	env = pathcontrol.UpdateEnvEntry(env, "MW_MCP_SERVER_EMBEDDED_CONNECTOR_DETAILS_TIMEOUT", "10s")

	cfg := mockmatlab.HappyConfig()
	value, err := cfg.ToEnvValue()
	s.Require().NoError(err, "failed to serialize mock config")
	env = pathcontrol.UpdateEnvEntry(env, mockmatlab.EnvMockMATLABConfig, value)

	session, err := s.CreateSessionWithEnv(env)
	s.Require().NoError(err)
	return session
}

// --- Assertions ---

func (s *ExistingExplicitSessionTestSuite) assertSecretNotInServerLogs(session *mcpclient.LoggedSession, secret string) {
	s.T().Helper()
	s.Require().NotEmpty(secret, "secret must not be empty — test is misconfigured")

	serverLogs, err := session.ReadServerLogs()
	s.Require().NoError(err)
	s.NotContains(serverLogs, secret, "server logs must not contain secrets in plaintext")
}
