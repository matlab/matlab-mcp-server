// Copyright 2026 The MathWorks, Inc.

package mcp_test

import (
	"runtime"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/fakematlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/pathcontrol"
	"github.com/stretchr/testify/suite"
)

// DisplayModeTestSuite tests that the --matlab-display-mode flag is correctly
// forwarded to MATLAB as the appropriate startup flags.
type DisplayModeTestSuite struct {
	MCPTestSuite
}

func TestDisplayModeTestSuite(t *testing.T) {
	suite.Run(t, new(DisplayModeTestSuite))
}

func (s *DisplayModeTestSuite) TestDisplayMode_Desktop_PassesDesktopFlag() {
	// Arrange & Act
	_, err := s.createSessionWithDisplayMode("desktop")
	s.Require().NoError(err, "should create MCP session with desktop mode")

	// Assert
	startupInfo := s.WaitForStartupInfo()
	s.Contains(startupInfo.Args, "-desktop", "desktop mode should pass -desktop to MATLAB")
	s.NotContains(startupInfo.Args, "-nodesktop", "desktop mode should NOT pass -nodesktop to MATLAB")
	s.NotContains(startupInfo.Args, "-nosplash", "desktop mode should NOT pass -nosplash to MATLAB")
	s.NotContains(startupInfo.Args, "-softwareopengl", "desktop mode should NOT pass -softwareopengl to MATLAB")
}

func (s *DisplayModeTestSuite) TestDisplayMode_Nodesktop_PassesCorrectFlags() {
	// Arrange & Act
	_, err := s.createSessionWithDisplayMode("nodesktop")
	s.Require().NoError(err, "should create MCP session with nodesktop mode")

	// Assert
	startupInfo := s.WaitForStartupInfo()
	s.Contains(startupInfo.Args, "-nodesktop", "nodesktop mode should pass -nodesktop to MATLAB")
	s.Contains(startupInfo.Args, "-nosplash", "nodesktop mode should pass -nosplash to MATLAB")
	s.Contains(startupInfo.Args, "-softwareopengl", "nodesktop mode should pass -softwareopengl to MATLAB")
	s.NotContains(startupInfo.Args, "-desktop", "nodesktop mode should NOT pass -desktop to MATLAB")

	if runtime.GOOS == "windows" {
		s.Contains(startupInfo.Args, "-noDisplayDesktop", "nodesktop mode on Windows should pass -noDisplayDesktop to MATLAB")
	}
}

// createSessionWithDisplayMode creates an MCP client session with the given
// display mode flag appended to the server args.
func (s *DisplayModeTestSuite) createSessionWithDisplayMode(displayMode string, opts ...mcpclient.CreateSessionOption) (*mcpclient.MCPClientSession, error) {
	s.T().Helper()
	env := pathcontrol.UpdateEnvEntry(s.baseEnv, fakematlab.OutputFileEnvVar, s.currentOutputFile)
	client := mcpclient.NewClient(
		s.T().Context(),
		s.mcpServerPath,
		env,
		"--log-level=debug",
		"--initialize-matlab-on-startup",
		"--matlab-display-mode="+displayMode,
	)
	return client.CreateSession(s.T().Context(), opts...)
}
