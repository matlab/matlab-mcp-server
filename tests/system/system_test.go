// Copyright 2025-2026 The MathWorks, Inc.

package system_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/pathcontrol"
	"github.com/stretchr/testify/suite"
)

// CLITestSuite tests command-line flags and basic server functionality
type CLITestSuite struct {
	SystemTestSuite
}

// TestVersionFlag validates the --version CLI flag
func (s *CLITestSuite) TestVersionFlag() {
	cmd := exec.Command(s.mcpServerPath, "--version") //nolint:gosec // Trusted test path
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "version flag should execute successfully")
	s.Contains(string(output), "github.com/matlab/matlab-mcp-core-server", "should display server package path")
}

func (s *CLITestSuite) TestMATLABRootFlag() {
	// Remove MATLABs from the path to directly test the --matlab-root flag
	newPath := pathcontrol.RemoveAllMATLABsFromPath(os.Getenv("PATH"))
	env := pathcontrol.UpdateEnvEntry(os.Environ(), "PATH", newPath)

	ctx := s.T().Context()
	session := s.CreateMCPSession(ctx, env, nil, "--matlab-root="+s.matlabRoot())
	defer s.CleanupSession(session, true)

	output, err := session.EvaluateCode(ctx, "disp('Server functional'); 2+2", s.testDataDir)
	s.Require().NoError(err, "evaluating MATLAB code should not error")
	s.Contains(output, "Server functional")
	s.Contains(output, "4")
}

// TestCLISuite runs the CLI test suite
func TestCLISuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}
