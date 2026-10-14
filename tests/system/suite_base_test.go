// Copyright 2025-2026 The MathWorks, Inc.

package system_test

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/usecases/utils/matlabstring"
	"github.com/matlab/matlab-mcp-server/tests/system/testdata"
	"github.com/matlab/matlab-mcp-server/tests/testconfig"
	"github.com/matlab/matlab-mcp-server/tests/testutils/matlablocator"
	"github.com/matlab/matlab-mcp-server/tests/testutils/matlabtempdir"
	"github.com/matlab/matlab-mcp-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-server/tests/testutils/mcpserver"
	"github.com/matlab/matlab-mcp-server/tests/testutils/pathcontrol"
	"github.com/matlab/matlab-mcp-server/tests/testutils/serverlogs"
	"github.com/stretchr/testify/suite"
)

type errorLogPolicy bool

const (
	assertNoErrorLogs errorLogPolicy = true
	allowErrorLogs    errorLogPolicy = false
)

// CreateMCPSession defaults to newSessionMode, so a test does not attach to a MATLAB that the user shared.
const (
	newSessionMode      = "--matlab-session-mode=new"
	existingSessionMode = "--matlab-session-mode=existing"
	autoSessionMode     = "--matlab-session-mode=auto"
)

// SystemTestSuite provides common setup for system tests
type SystemTestSuite struct {
	suite.Suite
	mcpServerPath     string
	matlabPath        string
	testDataDir       string
	pathEnvWithMATLAB string
	defaultEnv        []string
}

type SystemSession struct {
	*mcpclient.MCPClientSession
	logDir string
	logFS  fs.FS
}

func (s *SystemSession) LogDir() string {
	return s.logDir
}

func (s *SystemSession) LogFS() fs.FS {
	return s.logFS
}

func (s *SystemSession) ReadServerLogs() (string, error) {
	return s.readLogs("server-*.log")
}

func (s *SystemSession) ReadWatchdogLogs() (string, error) {
	return s.readLogs("watchdog-*.log")
}

func (s *SystemSession) readLogs(globPattern string) (string, error) {
	logFiles, err := fs.Glob(s.logFS, globPattern)
	if err != nil {
		return "", fmt.Errorf("failed to glob logs: %w", err)
	}
	if len(logFiles) == 0 {
		return "", fmt.Errorf("no logs found for pattern %s", globPattern)
	}

	var combined strings.Builder
	for _, logFile := range logFiles {
		content, err := fs.ReadFile(s.logFS, logFile)
		if err != nil {
			return "", fmt.Errorf("failed to read log file %s: %w", logFile, err)
		}
		combined.Write(content)
	}

	return combined.String(), nil
}

func (s *SystemSession) DumpLogsOnFailure(t *testing.T) {
	t.Helper()
	if !t.Failed() {
		return
	}
	stderr := s.Stderr()
	if stderr != "" {
		t.Logf("=== MCP Server Logs (stderr) ===\n%s\n=== End MCP Server Logs ===", stderr)
	}
	serverLogFiles, err := fs.Glob(s.LogFS(), "server-*.log")
	if err == nil && len(serverLogFiles) > 0 {
		for _, logFile := range serverLogFiles {
			serverLog, err := fs.ReadFile(s.LogFS(), logFile)
			if err != nil {
				t.Logf("Failed to read server log file: %s", err.Error())
			} else {
				t.Logf("=== MCP Server Log File (%s) ===\n%s\n=== End MCP Server Log File ===", filepath.Base(logFile), serverLog)
			}
		}
	}
	watchdogLogFiles, err := fs.Glob(s.LogFS(), "watchdog-*.log")
	if err == nil && len(watchdogLogFiles) > 0 {
		for _, logFile := range watchdogLogFiles {
			watchdogLog, err := fs.ReadFile(s.LogFS(), logFile)
			if err != nil {
				t.Logf("Failed to read watchdog log file: %s", err.Error())
			} else {
				t.Logf("=== MCP Watchdog Log File (%s) ===\n%s\n=== End MCP Watchdog Log File ===", filepath.Base(logFile), watchdogLog)
			}
		}
	}
}

func (s *SystemTestSuite) AssertNoErrorLogs(session *SystemSession) {
	errorLogs, err := serverlogs.ReadErrorLogs(session.LogFS())
	s.NoError(err) //nolint:testifylint // assert in defer to avoid FailNow
	s.Empty(errorLogs, "unexpected ERROR logs in server logs")
}

// SetupSuite runs once before all tests in a suite
func (s *SystemTestSuite) SetupSuite() {
	// Get MCP Server binary path
	mcpServerPath, err := mcpserver.NewLocator().GetPath()
	s.Require().NoError(err, "Failed to get MCP Server binary path")
	s.Require().NotEmpty(mcpServerPath, "MCP Server binary path cannot be empty")
	s.mcpServerPath = mcpServerPath

	// Get MATLAB path
	matlabPath, err := matlablocator.GetPath()
	s.Require().NoError(err, "Failed to get MATLAB path")
	s.Require().NotEmpty(matlabPath, "A non empty MATLAB path is required to run the system tests")
	s.matlabPath = matlabPath

	// Extract test assets to a temporary directory
	// This allows the system tests to be self-contained and compiled with the binary
	testDataDir := matlabtempdir.New(s.T(), "system-test-data-")

	// Extract files from embedded FS
	err = testdata.CopyToDir(testDataDir)
	s.Require().NoError(err, "Failed to extract test assets")

	s.testDataDir = testDataDir
}

// SetupTest runs before each test
func (s *SystemTestSuite) SetupTest() {
	s.pathEnvWithMATLAB, s.defaultEnv = s.matlabEnv()
}

func (s *SystemTestSuite) matlabEnv() (pathWithMATLAB string, envWithMATLAB []string) {
	path := os.Getenv("PATH")
	path = pathcontrol.RemoveAllMATLABsFromPath(path)
	path = pathcontrol.AddToPath(path, []string{s.matlabPath})
	s.Require().Contains(path, s.matlabPath, "MATLAB directory should be in the PATH environment variable")

	return path, pathcontrol.UpdateEnvEntry(os.Environ(), "PATH", path)
}

// CreateMCPSession creates an MCP client session with debug logging enabled.
// The returned session provides helper methods for reading and dumping logs.
//
// The caller is responsible for closing the session.
//
// Usage:
//
//	session := s.CreateMCPSession(ctx, nil, nil)
//	defer s.CleanupSession(session, assertNoErrorLogs)
//
// If env is nil, the suite's defaultEnv is used.
func (s *SystemTestSuite) CreateMCPSession(ctx context.Context, env []string, sessionOpts []mcpclient.CreateSessionOption, args ...string) *SystemSession {
	if env == nil {
		env = s.defaultEnv
	}

	hasLogLevel := false
	hasLogFolder := false
	hasSessionMode := false
	logFolderLocation := ""
	for _, arg := range args {
		if strings.HasPrefix(arg, "--log-level=") {
			hasLogLevel = true
		}
		if strings.HasPrefix(arg, "--log-folder=") {
			hasLogFolder = true
			logFolderLocation = strings.TrimPrefix(arg, "--log-folder=")
		}
		if strings.HasPrefix(arg, "--matlab-session-mode=") {
			hasSessionMode = true
		}
	}

	if !hasLogFolder {
		logFolderLocation = filepath.Join(matlabtempdir.New(s.T(), "mcp-logs-"), "logs")
		s.Require().NoError(os.MkdirAll(logFolderLocation, 0750), "should create log folder")
	}

	defaults := make([]string, 0, 3)
	if !hasLogLevel {
		defaults = append(defaults, "--log-level=debug")
	}
	if !hasLogFolder {
		defaults = append(defaults, "--log-folder="+logFolderLocation)
	}
	if !hasSessionMode {
		defaults = append(defaults, newSessionMode)
	}
	args = append(defaults, args...)

	client := mcpclient.NewClient(ctx, s.mcpServerPath, env, args...)
	mcpSession, err := client.CreateSession(ctx, sessionOpts...)
	s.Require().NoError(err, "should create MCP session")

	return &SystemSession{
		MCPClientSession: mcpSession,
		logDir:           logFolderLocation,
		logFS:            os.DirFS(logFolderLocation),
	}
}

func (s *SystemTestSuite) CleanupSession(session *SystemSession, errorLogs errorLogPolicy) {
	s.T().Helper()
	if err := session.Close(); err != nil {
		s.T().Logf("Ignoring session.Close() error (MCP go-sdk shutdown race): %v", err)
	}
	if errorLogs == assertNoErrorLogs {
		s.AssertNoErrorLogs(session)
	}
	session.DumpLogsOnFailure(s.T())
}

// Test file paths
func (s *SystemTestSuite) problematicCodePath() string {
	return filepath.Join(s.testDataDir, "problematic_code.m")
}

func (s *SystemTestSuite) testScriptPath() string {
	return filepath.Join(s.testDataDir, "test_script.m")
}

func (s *SystemTestSuite) testMathFunctionsPath() string {
	return filepath.Join(s.testDataDir, "test_math_functions.m")
}

func (s *SystemTestSuite) testPCodablePath() string {
	return filepath.Join(s.testDataDir, "test_pcodable.m")
}

// preparePCodeFiles compiles test_script.m and test_pcodable.m to P-code in a
// fresh temporary folder and returns the paths of the resulting .p files.
func (s *SystemTestSuite) preparePCodeFiles(ctx context.Context, session *SystemSession) (scriptPCodePath, testPCodePath string) {
	s.T().Helper()

	pcodeDir := s.T().TempDir()
	scriptPCodePath = filepath.Join(pcodeDir, "test_script.p")
	testPCodePath = filepath.Join(pcodeDir, "test_pcodable.p")

	_, err := session.EvaluateCode(ctx,
		fmt.Sprintf("pcode('%s', '%s')",
			matlabstring.EscapeSingleQuotes(s.testScriptPath()), matlabstring.EscapeSingleQuotes(s.testPCodablePath())), pcodeDir)
	s.Require().NoError(err, "should compile p-code files")
	s.Require().FileExists(scriptPCodePath, "pcode should produce a .p script")
	s.Require().FileExists(testPCodePath, "pcode should produce a .p test file")

	return scriptPCodePath, testPCodePath
}

// matlabRoot returns the MATLAB root directory (without /bin)
// This is needed for the --matlab-root flag
func (s *SystemTestSuite) matlabRoot() string {
	// matlabPath is expected to be the bin directory, so get its parent
	return filepath.Dir(s.matlabPath)
}

func (s *SystemTestSuite) matlabExe() string {
	return filepath.Join(s.matlabPath, testconfig.MATLABExeName)
}
