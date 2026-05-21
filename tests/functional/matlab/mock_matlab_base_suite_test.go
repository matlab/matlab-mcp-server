// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/facades/filefacade"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/logs"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpserver"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/serverlogs"
	"github.com/stretchr/testify/suite"
)

// MockMATLABBaseSuite provides the shared infrastructure for functional tests
// that use a mock MATLAB: the compiled mock binary, the MCP server binary, and
// the logged-session factory. It does NOT configure PATH — sub-suites own
// their environment setup.
type MockMATLABBaseSuite struct {
	suite.Suite
	mcpServerPath  string
	installation   *mockmatlab.Installation
	sessionFactory *mcpclient.LoggedSessionFactory
}

func (s *MockMATLABBaseSuite) SetupSuite() {
	s.installation = mockmatlab.BuildAndInstall(s.T())

	mcpServerPath, err := mcpserver.NewLocator().GetPath()
	s.Require().NoError(err, "MCP server binary not found — run 'make build' first")
	s.mcpServerPath = mcpServerPath

	sessionFactory, err := mcpclient.NewLoggedSessionFactory(logs.NewReader(), filefacade.RealFileSystem{})
	s.Require().NoError(err)
	s.sessionFactory = sessionFactory
}

// CreateSessionWithEnv creates an MCP session with a custom environment and CLI args.
func (s *MockMATLABBaseSuite) CreateSessionWithEnv(env []string, args ...string) (*mcpclient.LoggedSession, error) {
	return s.createLoggedSession(s.T().Context(), env, args...)
}

func (s *MockMATLABBaseSuite) createLoggedSession(ctx context.Context, env []string, args ...string) (*mcpclient.LoggedSession, error) {
	preparedArgs, err := logs.PrepareSessionCLIArgs(args, "debug", "mcp-functional-logs-")
	s.Require().NoError(err, "should prepare log args")
	s.T().Cleanup(func() {
		if err := os.RemoveAll(preparedArgs.TempBaseDir); err != nil {
			s.T().Logf("Failed to remove log temp dir (may be locked on Windows): %v", err)
		}
	})

	client := mcpclient.NewClient(ctx, s.mcpServerPath, env, preparedArgs.Args...)
	session, err := client.CreateSession(ctx)
	if err != nil {
		return nil, err
	}

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

	return loggedSession, nil
}

// AssertNoErrorLogs checks server log files for ERROR-level entries.
// Use assert (not require) so deferred cleanup continues if this fails.
func (s *MockMATLABBaseSuite) AssertNoErrorLogs(session *mcpclient.LoggedSession) {
	errorLogs, err := serverlogs.ReadErrorLogs(session.LogFS())
	s.NoError(err) //nolint:testifylint // assert in defer to avoid FailNow
	s.Empty(errorLogs, "unexpected ERROR logs in server logs")
}

func (s *MockMATLABBaseSuite) CleanupSession(session *mcpclient.LoggedSession, assertNoErrorLogs bool) {
	s.T().Helper()
	if err := session.Close(); err != nil {
		s.T().Logf("Ignoring session.Close() error (MCP go-sdk shutdown race): %v", err)
	}
	if assertNoErrorLogs {
		s.AssertNoErrorLogs(session)
	}
	session.DumpLogsOnFailure(s.T())
}
