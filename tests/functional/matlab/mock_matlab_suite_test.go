// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/logs"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/pathcontrol"
)

type MockMATLABSession struct {
	*mcpclient.LoggedSession
	mockMATLABLogDir string
}

func (s *MockMATLABSession) ReadInstanceEvents() ([]mockruntime.InstanceEvents, error) {
	return mockruntime.ReadEventsForAllInstances(s.mockMATLABLogDir)
}

// MockMATLABTestSuite extends the base with an environment that has mock MATLAB
// on PATH, suitable for testing local-install mode.
type MockMATLABTestSuite struct {
	MockMATLABBaseSuite
	defaultEnv []string
}

func (s *MockMATLABTestSuite) SetupSuite() {
	s.MockMATLABBaseSuite.SetupSuite()

	mockMATLABBinDir := filepath.Join(s.installation.MATLABRoot, "bin")
	path := pathcontrol.RemoveAllMATLABsFromPath(os.Getenv("PATH"))
	path = pathcontrol.AddToPath(path, []string{mockMATLABBinDir})

	env := pathcontrol.UpdateEnvEntry(os.Environ(), "PATH", path)
	env = pathcontrol.UpdateEnvEntry(env, "MW_MCP_SERVER_EMBEDDED_CONNECTOR_DETAILS_TIMEOUT", "10s")

	s.defaultEnv = env
}

// CreateSession starts the MCP server with mock MATLAB on PATH and returns an active MCP session.
// Pass extra env entries as KEY=VALUE strings to augment the default environment, or nil for defaults.
// Additional CLI args (e.g. "--extension-file=...") are forwarded to the MCP server binary.
func (s *MockMATLABTestSuite) CreateSession(cfg mockmatlab.Config, extraEnv []string, args ...string) (*MockMATLABSession, error) {
	env := s.defaultEnv
	if len(extraEnv) > 0 {
		env = append(slices.Clone(s.defaultEnv), extraEnv...)
	}

	value, err := cfg.ToEnvValue()
	s.Require().NoError(err, "failed to serialize mock config")
	env = pathcontrol.UpdateEnvEntry(env, mockmatlab.EnvMockMATLABConfig, value)

	preparedArgs, err := logs.PrepareSessionCLIArgs(args, "debug", "mcp-functional-logs-")
	s.Require().NoError(err, "should prepare log args")
	s.T().Cleanup(func() {
		if err := os.RemoveAll(preparedArgs.TempBaseDir); err != nil {
			s.T().Logf("Failed to remove log temp dir %s: %v", preparedArgs.TempBaseDir, err)
		}
	})

	mockMATLABLogDir := filepath.Join(preparedArgs.TempBaseDir, "mock-matlab-logs")
	env = pathcontrol.UpdateEnvEntry(env, mockmatlab.EnvMockMATLABLogDir, mockMATLABLogDir)

	ctx := s.T().Context()
	client := mcpclient.NewClient(ctx, s.mcpServerPath, env, preparedArgs.Args...)
	session, sessionErr := client.CreateSession(ctx)

	loggedSession, err := s.sessionFactory.New(
		session,
		preparedArgs.LogDir,
		"MCP Server Logs (stderr)",
		[]logs.DumpPattern{
			{Glob: "server-*.log", Header: "MCP Server Log File"},
			{Glob: "watchdog-*.log", Header: "MCP Watchdog Log File"},
		},
	)
	if err != nil {
		return nil, err
	}

	return &MockMATLABSession{
		LoggedSession:    loggedSession,
		mockMATLABLogDir: mockMATLABLogDir,
	}, sessionErr
}

func (s *MockMATLABTestSuite) CleanupSession(session *MockMATLABSession, assertNoErrorLogs bool) {
	s.T().Helper()
	s.MockMATLABBaseSuite.CleanupSession(session.LoggedSession, assertNoErrorLogs)
}
