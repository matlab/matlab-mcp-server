// Copyright 2026 The MathWorks, Inc.

package mcpb_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpbundle"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmcpbinary"
	"github.com/stretchr/testify/suite"
)

// MCPBLauncherBehaviorSuite verifies the launcher script correctly translates env vars into CLI flags for the binary.
type MCPBLauncherBehaviorSuite struct {
	suite.Suite

	bundle *mcpbundle.Bundle
}

func (s *MCPBLauncherBehaviorSuite) SetupSuite() {
	s.bundle = mcpbundle.Open(s.T())
	mockmcpbinary.BuildAndInstall(s.T(), filepath.Join(s.bundle.Dir(), "bin"))
}

func TestMCPBLauncherBehaviorSuite(t *testing.T) {
	suite.Run(t, new(MCPBLauncherBehaviorSuite))
}

func (s *MCPBLauncherBehaviorSuite) TestString_PassesFlagAndValue() {
	result := s.bundle.Launch(s.T(), map[string]string{
		"__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_ROOT": mcpbundle.PathWithSpaces(),
	})

	s.Contains(result.Args, "--matlab-root")
	s.Contains(result.Args, mcpbundle.PathWithSpaces())
}

func (s *MCPBLauncherBehaviorSuite) TestString_OmittedWhenEmpty() {
	result := s.bundle.Launch(s.T(), map[string]string{
		"__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_ROOT": "",
	})

	s.Empty(result.Args, "empty string env var should produce no flags")
}

func (s *MCPBLauncherBehaviorSuite) TestString_OmittedWhenUnset() {
	result := s.bundle.Launch(s.T(), nil)

	s.Empty(result.Args, "no flags should be passed when no env vars are set")
}

func (s *MCPBLauncherBehaviorSuite) TestBool_PassesFlagWhenTrue() {
	result := s.bundle.Launch(s.T(), map[string]string{
		"__MATLAB_MCP_CORE_SERVER_MCPB_INIT_ON_START": "true",
	})

	s.Contains(result.Args, "--initialize-matlab-on-startup")
}

func (s *MCPBLauncherBehaviorSuite) TestBool_OmittedWhenFalse() {
	result := s.bundle.Launch(s.T(), map[string]string{
		"__MATLAB_MCP_CORE_SERVER_MCPB_INIT_ON_START": "false",
	})

	s.Empty(result.Args, "false bool env var should produce no flags")
}

func (s *MCPBLauncherBehaviorSuite) TestEnvVars_UnsetBeforeExec() {
	result := s.bundle.Launch(s.T(), map[string]string{
		"__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_ROOT":   "/opt/matlab",
		"__MATLAB_MCP_CORE_SERVER_MCPB_INIT_ON_START": "true",
	})

	s.Empty(result.Env, "MCPB env vars should be unset before exec-ing the binary")
}

func (s *MCPBLauncherBehaviorSuite) TestArgs_ExtensionFilesExpandedToRepeatedFlags() {
	result := s.bundle.Launch(s.T(), nil,
		"--extension-files", "/path/to/tools-a.json", "/path/to/tools-b.json",
	)

	s.Equal([]string{
		"--extension-file", "/path/to/tools-a.json",
		"--extension-file", "/path/to/tools-b.json",
	}, result.Args)
}

func (s *MCPBLauncherBehaviorSuite) TestArgs_ExtensionFilePathsWithSpacesPreserved() {
	pathWithSpaces := mcpbundle.PathWithSpaces() + "/my-tools.json"

	result := s.bundle.Launch(s.T(), nil, "--extension-files", pathWithSpaces)

	s.Equal([]string{"--extension-file", pathWithSpaces}, result.Args)
}

func (s *MCPBLauncherBehaviorSuite) TestArgs_ExtensionFilesCombineWithEnvVarFlags() {
	result := s.bundle.Launch(s.T(), map[string]string{
		"__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_ROOT": "/opt/matlab",
	}, "--extension-files", "/path/to/tools.json")

	s.Contains(result.Args, "--matlab-root")
	s.Contains(result.Args, "/opt/matlab")
	s.Contains(result.Args, "--extension-file")
	s.Contains(result.Args, "/path/to/tools.json")
}

func (s *MCPBLauncherBehaviorSuite) TestLauncher_ValidSyntax() {
	s.Require().NoError(s.bundle.CheckLauncherSyntax())
}
