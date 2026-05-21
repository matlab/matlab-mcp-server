// Copyright 2026 The MathWorks, Inc.

package mcp_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/stretchr/testify/suite"
)

// RootsTestSuite tests that MCP roots are correctly forwarded to MATLAB's working directory.
//
// The following scenarios require a mock MATLAB that stays alive (see mockmatlab package)
// and are not yet covered:
//   - ListChanged: adding a root after session initialization updates the root store.
//   - ListChanged: replacing roots after session initialization updates the root store.
//   - updateRoots error path: ListRoots fails and the error is logged as a warning.
//   - updateRoots happy path: ListRoots succeeds and roots are stored + logged.
type RootsTestSuite struct {
	MCPTestSuite
}

func TestRootsTestSuite(t *testing.T) {
	suite.Run(t, new(RootsTestSuite))
}

func (s *RootsTestSuite) TestMATLABStartsInFirstRootDirectory() {
	// Arrange
	workspaceDir := s.T().TempDir()

	// Act — session creation triggers eager MATLAB launch with the root as working dir.
	// The fake MATLAB writes its startup info and exits. The server will detect MATLAB
	// died but that is expected — we only care about the startup info.
	_ = s.CreateSession(
		mcpclient.WithRoots(NewRootFromDir(workspaceDir, "workspace")),
	)

	// Assert
	startupInfo := s.WaitForStartupInfo()
	s.Equal(
		s.normalizedPath(workspaceDir),
		s.normalizedPath(startupInfo.WorkingDir),
		"MATLAB should start in the MCP root directory",
	)
	s.NotEmpty(startupInfo.Args, "fake MATLAB should have received command-line arguments")
}

func (s *RootsTestSuite) TestMultipleRoots_MATLABStartsInFirstRoot() {
	// Arrange
	projectDir1 := s.T().TempDir()
	projectDir2 := s.T().TempDir()

	// Act
	_ = s.CreateSession(
		mcpclient.WithRoots(
			NewRootFromDir(projectDir1, "project-1"),
			NewRootFromDir(projectDir2, "project-2"),
		),
	)

	// Assert — server should use the first root as MATLAB's working directory.
	startupInfo := s.WaitForStartupInfo()
	s.Equal(
		s.normalizedPath(projectDir1),
		s.normalizedPath(startupInfo.WorkingDir),
		"MATLAB should start in the first MCP root directory",
	)
}

func (s *RootsTestSuite) TestNoRoots_MATLABStillLaunches() {
	// Arrange & Act — no roots provided, but MATLAB should still launch.
	_ = s.CreateSession()

	// Assert — the fake MATLAB should still have been launched (with some default working directory).
	startupInfo := s.WaitForStartupInfo()
	s.NotEmpty(startupInfo.WorkingDir, "MATLAB should have a working directory even without roots")
	s.NotEmpty(startupInfo.Args, "fake MATLAB should have received command-line arguments")
}

func (s *RootsTestSuite) TestNoRootsCapability_MATLABStillLaunches() {
	// Arrange & Act — client does not advertise roots capability.
	// The server should not request roots, but MATLAB should still launch.
	_ = s.CreateSessionWithoutRootsCapability()

	// Assert — MATLAB launches with a default working directory (not root-derived).
	startupInfo := s.WaitForStartupInfo()
	s.NotEmpty(startupInfo.WorkingDir, "MATLAB should have a working directory even without roots capability")
	s.NotEmpty(startupInfo.Args, "fake MATLAB should have received command-line arguments")
}

// normalizedPath returns a cleaned, lowercased absolute path for comparison.
func (s *RootsTestSuite) normalizedPath(dir string) string {
	s.T().Helper()
	abs, err := filepath.Abs(dir)
	s.Require().NoError(err)
	return strings.ToLower(filepath.Clean(abs))
}
