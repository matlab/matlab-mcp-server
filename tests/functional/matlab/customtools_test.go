// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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

	//go:embed testdata/customtools/second_tool.json
	secondToolJSON string

	//go:embed testdata/customtools/cross_file_conflict_tool.json
	crossFileConflictToolJSON string

	//go:embed testdata/customtools/intra_file_duplicate_tool.json
	intraFileDuplicateToolJSON string
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
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil, "--extension-file="+extensionFile)
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
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil, "--extension-file="+extensionFile)
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
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil, "--extension-file="+extensionFile)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	ctx := s.T().Context()
	result, err := session.ListTools(ctx, nil)
	s.Require().NoError(err, "should list tools")

	toolNames := toolNamesFromResult(result)
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

	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil, "--extension-file="+extensionFile)
	s.Require().Error(err, "session creation should fail with invalid extension file")

	serverLogs, logErr := session.ReadServerLogs()
	s.Require().NoError(logErr)
	s.Contains(serverLogs, "Failed to parse custom tools extension file")
}

func (s *CustomToolsTestSuite) TestErrorPath_MissingExtensionFile_ServerFails() {
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil, "--extension-file=/nonexistent/path/tools.json")
	s.Require().Error(err, "session creation should fail when extension file does not exist")

	serverLogs, logErr := session.ReadServerLogs()
	s.Require().NoError(logErr)
	s.Contains(serverLogs, "Failed to read custom tools extension file")
}

func (s *CustomToolsTestSuite) TestErrorPath_ToolNameConflict_ServerFails() {
	extensionFile := writeExtensionFile(s.T(), nameConflictToolJSON)

	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil, "--extension-file="+extensionFile)
	s.Require().Error(err, "session creation should fail when custom tool name conflicts with built-in tool")

	serverLogs, logErr := session.ReadServerLogs()
	s.Require().NoError(logErr)
	s.Contains(serverLogs, "Failed to configure tools")
}

func (s *CustomToolsTestSuite) TestErrorPath_IntraFileDuplicateToolName_ServerFails() {
	extensionFile := writeExtensionFile(s.T(), intraFileDuplicateToolJSON)

	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil, "--extension-file="+extensionFile)
	s.Require().Error(err, "session creation should fail when same tool name appears twice in one file")

	serverLogs, logErr := session.ReadServerLogs()
	s.Require().NoError(logErr)
	s.Contains(serverLogs, "Invalid custom tool definition")
}

func (s *CustomToolsTestSuite) TestHappyPath_MultipleExtensionFiles_ToolsFromBothFilesAvailable() {
	fileA := writeExtensionFile(s.T(), singleToolJSON)
	fileB := writeExtensionFile(s.T(), secondToolJSON)

	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil,
		"--extension-file="+fileA,
		"--extension-file="+fileB,
	)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	ctx := s.T().Context()
	result, err := session.ListTools(ctx, nil)
	s.Require().NoError(err, "should list tools")

	toolNames := toolNamesFromResult(result)
	s.True(toolNames["generate_magic_square"], "tool from first file should be available")
	s.True(toolNames["compute_fibonacci"], "tool from second file should be available")

	magicResult, err := session.CallTool(ctx, "generate_magic_square", map[string]any{"n": float64(3)})
	s.Require().NoError(err, "should call tool from first file")
	magicText, err := session.GetTextContent(magicResult)
	s.Require().NoError(err)
	s.Contains(magicText, "magic(3)")

	fibResult, err := session.CallTool(ctx, "compute_fibonacci", map[string]any{"n": float64(10)})
	s.Require().NoError(err, "should call tool from second file")
	fibText, err := session.GetTextContent(fibResult)
	s.Require().NoError(err)
	s.Contains(fibText, "fibonacci(10)")
}

func (s *CustomToolsTestSuite) TestErrorPath_CrossFileToolNameCollision_ServerFails() {
	fileA := writeExtensionFile(s.T(), singleToolJSON)
	fileB := writeExtensionFile(s.T(), crossFileConflictToolJSON)

	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil,
		"--extension-file="+fileA,
		"--extension-file="+fileB,
	)
	s.Require().Error(err, "session creation should fail when same tool name appears in multiple extension files")

	serverLogs, logErr := session.ReadServerLogs()
	s.Require().NoError(logErr)
	s.Contains(serverLogs, "Failed to configure tools")
}

func (s *CustomToolsTestSuite) TestErrorPath_OneInvalidFileAmongMultiple_ServerFails() {
	validFile := writeExtensionFile(s.T(), singleToolJSON)
	malformedFile := writeExtensionFile(s.T(), malformedJSON)

	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil,
		"--extension-file="+validFile,
		"--extension-file="+malformedFile,
	)
	s.Require().Error(err, "session creation should fail when any extension file is invalid")

	serverLogs, logErr := session.ReadServerLogs()
	s.Require().NoError(logErr)
	s.Contains(serverLogs, "Failed to parse custom tools extension file")
}

func (s *CustomToolsTestSuite) TestErrorPath_DuplicateFilePath_ServerFails() {
	extensionFile := writeExtensionFile(s.T(), singleToolJSON)

	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil,
		"--extension-file="+extensionFile,
		"--extension-file="+extensionFile,
	)
	s.Require().Error(err, "session creation should fail when the same extension file is specified twice")

	serverLogs, logErr := session.ReadServerLogs()
	s.Require().NoError(logErr)
	s.Contains(serverLogs, "Failed to configure tools")
}

func (s *CustomToolsTestSuite) TestHappyPath_EnvVarWithPathSeparator_LoadsMultipleFiles() {
	fileA := writeExtensionFile(s.T(), singleToolJSON)
	fileB := writeExtensionFile(s.T(), secondToolJSON)
	envValue := fileA + string(filepath.ListSeparator) + fileB

	session, err := s.CreateSession(mockmatlab.HappyConfig(),
		[]string{"MW_MCP_SERVER_EXTENSION_FILE=" + envValue},
	)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	ctx := s.T().Context()
	result, err := session.ListTools(ctx, nil)
	s.Require().NoError(err, "should list tools")

	toolNames := toolNamesFromResult(result)
	s.True(toolNames["generate_magic_square"], "tool from first file should be available via env var")
	s.True(toolNames["compute_fibonacci"], "tool from second file should be available via env var")
}

func (s *CustomToolsTestSuite) TestHappyPath_EnvVarAndCLIFlag_MergesBothSources() {
	envFile := writeExtensionFile(s.T(), singleToolJSON)
	cliFile := writeExtensionFile(s.T(), secondToolJSON)

	session, err := s.CreateSession(mockmatlab.HappyConfig(),
		[]string{"MW_MCP_SERVER_EXTENSION_FILE=" + envFile},
		"--extension-file="+cliFile,
	)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	ctx := s.T().Context()
	result, err := session.ListTools(ctx, nil)
	s.Require().NoError(err, "should list tools")

	toolNames := toolNamesFromResult(result)
	s.True(toolNames["generate_magic_square"], "tool from env var file should be available")
	s.True(toolNames["compute_fibonacci"], "tool from CLI file should be available")
}

func toolNamesFromResult(result *mcp.ListToolsResult) map[string]bool {
	names := make(map[string]bool)
	for _, tool := range result.Tools {
		names[tool.Name] = true
	}
	return names
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
