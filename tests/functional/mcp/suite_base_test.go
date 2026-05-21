// Copyright 2026 The MathWorks, Inc.

package mcp_test

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/fakematlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpserver"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/pathcontrol"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/suite"
)

// MCPTestSuite provides common setup for functional MCP tests that use a fake
// MATLAB installation. It launches the real MCP server binary (from `make build`)
// with the fake MATLAB on PATH.
type MCPTestSuite struct {
	suite.Suite

	mcpServerPath string
	fakeMatlab    fakematlab.Installation
	baseEnv       []string

	// currentOutputFile is the per-test output file path for the fake MATLAB.
	currentOutputFile string
}

func (s *MCPTestSuite) SetupSuite() {
	mcpServerPath, err := mcpserver.NewLocator().GetPath()
	s.Require().NoError(err, "MCP server binary not found — run 'make build' first")
	s.mcpServerPath = mcpServerPath

	s.fakeMatlab = fakematlab.CreateExecutable(s.T())

	path := os.Getenv("PATH")
	path = pathcontrol.RemoveAllMATLABsFromPath(path)
	path = s.fakeMatlab.PathEntry + string(os.PathListSeparator) + path
	s.baseEnv = pathcontrol.UpdateEnvEntry(os.Environ(), "PATH", path)
}

// SetupTest creates a fresh output file for each test so results don't leak across tests.
func (s *MCPTestSuite) SetupTest() {
	s.currentOutputFile = filepath.Join(s.T().TempDir(), "startup-info.json")
}

// FakeMatlab returns an Installation with the current test's output file path.
func (s *MCPTestSuite) FakeMatlab() fakematlab.Installation {
	return fakematlab.Installation{
		Root:           s.fakeMatlab.Root,
		PathEntry:      s.fakeMatlab.PathEntry,
		OutputFilePath: s.currentOutputFile,
	}
}

// CreateSession creates an MCP client session connected to the real MCP server.
// The server is started with --initialize-matlab-on-startup so MATLAB launches
// eagerly on session creation, without needing a tool call.
func (s *MCPTestSuite) CreateSession(opts ...mcpclient.CreateSessionOption) *mcpclient.MCPClientSession {
	s.T().Helper()
	env := pathcontrol.UpdateEnvEntry(s.baseEnv, fakematlab.OutputFileEnvVar, s.currentOutputFile)
	client := mcpclient.NewClient(
		s.T().Context(),
		s.mcpServerPath,
		env,
		"--log-level=debug",
		"--initialize-matlab-on-startup",
	)
	session, err := client.CreateSession(s.T().Context(), opts...)
	s.Require().NoError(err, "should create MCP session")
	return session
}

// CreateSessionWithoutRootsCapability creates an MCP client session where the
// client does not advertise roots capability. The server will not request roots.
func (s *MCPTestSuite) CreateSessionWithoutRootsCapability() *mcpclient.MCPClientSession {
	s.T().Helper()
	env := pathcontrol.UpdateEnvEntry(s.baseEnv, fakematlab.OutputFileEnvVar, s.currentOutputFile)
	client := mcpclient.NewClientWithoutRootsCapability(
		s.T().Context(),
		s.mcpServerPath,
		env,
		"--log-level=debug",
		"--initialize-matlab-on-startup",
	)
	session, err := client.CreateSession(s.T().Context())
	s.Require().NoError(err, "should create MCP session")
	return session
}

// ToFileURI converts an absolute directory path to a file:// URI.
func ToFileURI(dir string) string {
	slashPath := filepath.ToSlash(dir)
	if !strings.HasPrefix(slashPath, "/") {
		slashPath = "/" + slashPath
	}
	return "file://" + slashPath
}

// NewRootFromDir creates an mcp.Root from a local directory path.
func NewRootFromDir(dir, name string) *mcp.Root {
	return &mcp.Root{URI: ToFileURI(dir), Name: name}
}

// WaitForStartupInfo polls the fake MATLAB output file until startup info is available.
func (s *MCPTestSuite) WaitForStartupInfo() fakematlab.StartupInfo {
	s.T().Helper()
	var startupInfo fakematlab.StartupInfo
	s.Require().Eventually(func() bool {
		info, err := s.FakeMatlab().ReadStartupInfo()
		startupInfo = info
		return err == nil
	}, 30*time.Second, 500*time.Millisecond, "fake MATLAB should have written its startup info")
	return startupInfo
}
