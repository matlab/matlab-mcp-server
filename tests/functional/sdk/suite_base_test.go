// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/serverlogs"
	"github.com/stretchr/testify/suite"
)

type SDKSession struct {
	*mcpclient.MCPClientSession
	logDir string
	logFS  fs.FS
}

func (s *SDKSession) LogDir() string {
	return s.logDir
}

func (s *SDKSession) LogFS() fs.FS {
	return s.logFS
}

func (s *SDKSession) ReadServerLogs() (string, error) {
	serverLogFiles, err := fs.Glob(s.logFS, "server-*.log")
	if err != nil {
		return "", fmt.Errorf("failed to glob server logs: %w", err)
	}
	if len(serverLogFiles) == 0 {
		return "", fmt.Errorf("no server log files found")
	}

	var combined strings.Builder
	for _, logFile := range serverLogFiles {
		content, err := fs.ReadFile(s.logFS, logFile)
		if err != nil {
			return "", fmt.Errorf("failed to read server log file %s: %w", logFile, err)
		}
		combined.Write(content)
	}

	return combined.String(), nil
}

func (s *SDKSession) ReadAllServerLogs() (string, error) {
	return s.ReadServerLogs()
}

func (s *SDKSession) DumpLogsOnFailure(t *testing.T) {
	t.Helper()
	if !t.Failed() {
		return
	}
	stderr := s.Stderr()
	if stderr != "" {
		t.Logf("=== SDK Server Logs (stderr) ===\n%s\n=== End SDK Server Logs ===", stderr)
	}
	serverLogFiles, err := fs.Glob(s.LogFS(), "server-*.log")
	if err != nil || len(serverLogFiles) == 0 {
		return
	}
	for _, logFile := range serverLogFiles {
		content, err := fs.ReadFile(s.LogFS(), logFile)
		if err != nil {
			t.Logf("Failed to read server log file: %s", err.Error())
		} else {
			t.Logf("=== SDK Server Log File (%s) ===\n%s\n=== End SDK Server Log File ===", logFile, content)
		}
	}
}

type SDKTestSuite struct {
	suite.Suite
}

// CreateSession creates an MCP client session with debug logging enabled.
// The returned session provides helper methods for reading and dumping logs.
//
// The caller is responsible for closing the session and making any assertions.
//
// Usage:
//
//	session := s.CreateSession(serverPath, nil, nil)
//	defer s.CleanupSession(session, true)
func (s *SDKTestSuite) CreateSession(serverPath string, env []string, sessionOpts []mcpclient.CreateSessionOption, args ...string) *SDKSession {
	s.T().Helper()
	hasLogLevel := false
	hasLogFolder := false
	logFolder := ""
	for _, a := range args {
		if strings.HasPrefix(a, "--log-level=") {
			hasLogLevel = true
		}
		if strings.HasPrefix(a, "--log-folder=") {
			hasLogFolder = true
			logFolder = strings.TrimPrefix(a, "--log-folder=")
		}
	}

	if !hasLogFolder {
		base, err := os.MkdirTemp("", "sdk-logs-")
		s.Require().NoError(err, "should create SDK log temp dir")
		logFolder = filepath.Join(base, "logs")
		s.Require().NoError(os.MkdirAll(logFolder, 0750), "should create SDK log folder")
		s.T().Cleanup(func() {
			if err := os.RemoveAll(base); err != nil {
				s.T().Logf("Failed to remove SDK log temp dir (may be locked on Windows): %v", err)
			}
		})
	}

	defaults := make([]string, 0, 2)
	if !hasLogLevel {
		defaults = append(defaults, "--log-level=debug")
	}
	if !hasLogFolder {
		defaults = append(defaults, "--log-folder="+logFolder)
	}
	args = append(defaults, args...)
	client := mcpclient.NewClient(s.T().Context(), serverPath, env, args...)
	mcpSession, err := client.CreateSession(s.T().Context(), sessionOpts...)
	s.Require().NoError(err, "should create MCP session")

	session := &SDKSession{
		MCPClientSession: mcpSession,
		logDir:           logFolder,
		logFS:            os.DirFS(logFolder),
	}

	return session
}

// AssertNoErrorLogs checks server log files for ERROR-level entries.
// Use assert (not require) so deferred cleanup continues if this fails.
func (s *SDKTestSuite) AssertNoErrorLogs(session *SDKSession) {
	errorLogs, err := serverlogs.ReadErrorLogs(session.LogFS())
	s.NoError(err) //nolint:testifylint // assert in defer to avoid FailNow
	s.Empty(errorLogs, "unexpected ERROR logs in server logs")
}

func (s *SDKTestSuite) CleanupSession(session *SDKSession, assertNoErrorLogs bool) {
	s.T().Helper()
	if err := session.Close(); err != nil {
		s.T().Logf("Ignoring session.Close() error (MCP go-sdk shutdown race): %v", err)
	}
	if assertNoErrorLogs {
		s.AssertNoErrorLogs(session)
	}
	session.DumpLogsOnFailure(s.T())
}
