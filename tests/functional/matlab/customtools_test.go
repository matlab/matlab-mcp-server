// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/stretchr/testify/suite"
)

var (
	//go:embed testdata/customtools/single_tool.json
	singleToolJSON string

	//go:embed testdata/customtools/multiple_tools.json
	multipleToolsJSON string

	//go:embed testdata/customtools/name_conflict_tool.json
	nameConflictToolJSON string

	//go:embed testdata/customtools/malformed.json
	malformedJSON string
)

// CustomToolsTestSuite tests that using an extensions file with custom MATLAB tools works as expected.
type CustomToolsTestSuite struct {
	MockMATLABTestSuite
}

func TestCustomToolsTestSuite(t *testing.T) {
	suite.Run(t, new(CustomToolsTestSuite))
}

func (s *CustomToolsTestSuite) TestHappyPath_CustomToolsAreAvailable() {
	extensionFile := writeExtensionFile(s.T(), singleToolJSON)
	session, err := s.CreateSession(mockmatlab.HappyConfig(), "--extension-file="+extensionFile)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	ctx := s.T().Context()
	result, err := session.ListTools(ctx, nil)
	s.Require().NoError(err, "should list tools")

	var found bool
	for _, tool := range result.Tools {
		if tool.Name == "generate_magic_square" {
			found = true
			s.Equal("Generate Magic Square", tool.Title)
			s.Equal("Generates an n-by-n magic square matrix", tool.Description)
			break
		}
	}
	s.True(found, "custom tool should appear in tools list")
}

func (s *CustomToolsTestSuite) TestHappyPath_CallCustomTool_ReturnsOutput() {
	extensionFile := writeExtensionFile(s.T(), singleToolJSON)
	session, err := s.CreateSession(mockmatlab.HappyConfig(), "--extension-file="+extensionFile)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	ctx := s.T().Context()
	result, err := session.CallTool(ctx, "generate_magic_square", map[string]any{"n": float64(5)})
	s.Require().NoError(err, "should call custom tool")

	text, err := session.GetTextContent(result)
	s.Require().NoError(err, "should get text content")
	s.Contains(text, "magic(5)", "response should contain the assembled MATLAB function call")
}

func (s *CustomToolsTestSuite) TestHappyPath_MultipleCustomToolsLoaded() {
	extensionFile := writeExtensionFile(s.T(), multipleToolsJSON)
	session, err := s.CreateSession(mockmatlab.HappyConfig(), "--extension-file="+extensionFile)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	ctx := s.T().Context()
	result, err := session.ListTools(ctx, nil)
	s.Require().NoError(err, "should list tools")

	toolNames := make(map[string]bool)
	for _, tool := range result.Tools {
		toolNames[tool.Name] = true
	}
	s.True(toolNames["greet_user"], "greet_user should appear in tools list")
	s.True(toolNames["generate_magic_square"], "generate_magic_square should appear in tools list")

	callResult, err := session.CallTool(ctx, "greet_user", map[string]any{"name": "Alice"})
	s.Require().NoError(err, "should call greet_user tool")

	text, err := session.GetTextContent(callResult)
	s.Require().NoError(err, "should get text content")
	s.Contains(text, `greet_user("Alice")`, "response should contain the assembled MATLAB function call with string argument")
}

func (s *CustomToolsTestSuite) TestErrorPath_InvalidExtensionFile_ServerFails() {
	extensionFile := writeExtensionFile(s.T(), malformedJSON)

	_, err := s.CreateSession(mockmatlab.HappyConfig(), "--extension-file="+extensionFile)
	s.Require().Error(err, "session creation should fail with invalid extension file")
}

func (s *CustomToolsTestSuite) TestErrorPath_MissingExtensionFile_ServerFails() {
	_, err := s.CreateSession(mockmatlab.HappyConfig(), "--extension-file=/nonexistent/path/tools.json")
	s.Require().Error(err, "session creation should fail when extension file does not exist")
}

func (s *CustomToolsTestSuite) TestErrorPath_ToolNameConflict_ServerFails() {
	extensionFile := writeExtensionFile(s.T(), nameConflictToolJSON)

	_, err := s.CreateSession(mockmatlab.HappyConfig(), "--extension-file="+extensionFile)
	s.Require().Error(err, "session creation should fail when custom tool name conflicts with built-in tool")
}

func writeExtensionFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "tools.json")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write extension file: %v", err)
	}
	return path
}
