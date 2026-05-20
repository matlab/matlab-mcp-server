// Copyright 2025-2026 The MathWorks, Inc.

package system_test

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/system/testdata"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/suite"
)

// WorkflowTestSuite tests realistic end-to-end workflows based on how users actually
// interact with MATLAB through the MCP server in different scenarios.
type WorkflowTestSuite struct {
	SystemTestSuite
}

// TestInteractiveDevelopmentWorkflow simulates a user doing interactive MATLAB development
// with an AI application in a single conversation context. This is the typical use case
// for single-session mode where the user wants continuity: variables persist between
// requests, the workspace builds up over time, and there's a consistent working directory.
//
// Scenario: Data scientist developing MATLAB code interactively
// - Discovers available toolboxes to understand capabilities
// - Writes and tests code snippets iteratively (variables persist)
// - Checks code quality before committing
// - Runs existing scripts and test suites
// - All work happens in one continuous MATLAB session
//
// MCP Resources tested:
// - matlab_coding_guidelines (coding standards reference)
//
// MCP Tools tested:
// - detect_matlab_toolboxes (feature discovery)
// - evaluate_matlab_code (iterative development with persistent workspace)
// - check_matlab_code (code quality analysis)
// - run_matlab_file (script execution)
// - run_matlab_test_file (test execution)
func (s *WorkflowTestSuite) TestInteractiveDevelopmentWorkflow() {
	for _, tc := range []struct {
		displayMode string
	}{
		{displayMode: "desktop"},
		{displayMode: "nodesktop"},
	} {
		s.Run(tc.displayMode, func() {
			ctx := s.T().Context()

			// Provide the test data directory as an MCP root so the server uses it as
			// MATLAB's starting directory (fallback when no --initial-working-folder is set).
			// File URIs require the path to start with "/" (e.g., file:///C:/path on Windows).
			slashPath := filepath.ToSlash(s.testDataDir)
			if !strings.HasPrefix(slashPath, "/") {
				slashPath = "/" + slashPath
			}
			testDataURI := "file://" + slashPath
			sessionOpts := []mcpclient.CreateSessionOption{
				mcpclient.WithRoots(&mcp.Root{URI: testDataURI, Name: "test-data"}),
			}
			session := s.CreateMCPSession(ctx, nil, sessionOpts, "--matlab-display-mode="+tc.displayMode)
			defer s.CleanupSession(session, true)

			// Step 1: Read coding guidelines (AI references standards before writing code)
			guidelines, err := session.ReadResource(ctx, "guidelines://coding")
			s.Require().NoError(err, "should read coding guidelines resource")
			s.Contains(guidelines, "MATLAB Coding Standards", "should contain coding standards title")

			// Step 2: Feature discovery - check what toolboxes are available
			info, err := session.DetectToolboxes(ctx)
			s.Require().NoError(err, "should detect toolboxes")
			s.Contains(info, "MATLAB Version:", "should discover MATLAB version")

			// Step 3a: Iterative development with explicit integer math
			output, err := session.EvaluateCode(ctx, `a = int32(2); b = int32(3);`)
			s.Require().NoError(err)
			s.Empty(output, "semicolon-terminated statements should produce no output")

			// Step 3b: Verify variables persist and computation works
			output, err = session.EvaluateCode(ctx, `
				c = a + b;
				fprintf('a=%d\n', a);
				fprintf('b=%d\n', b);
				fprintf('c=%d\n', c);
			`)
			s.Require().NoError(err)
			s.Require().NotEmpty(output, "fprintf should produce output")
			lines := strings.Split(output, "\n")
			s.Contains(lines, "a=2", "variable 'a' should persist on its own line")
			s.Contains(lines, "b=3", "variable 'b' should persist on its own line")
			s.Contains(lines, "c=5", "should compute 2 + 3 = 5 on its own line")

			// Step 4: Verify MATLAB started in the MCP root directory
			// The server should have used the client's root as MATLAB's initial working folder.
			output, err = session.EvaluateCode(ctx, "pwd")
			s.Require().NoError(err)
			s.Contains(output, s.testDataDir, "MATLAB should start in the MCP root directory (test data dir)")

			// Verify the test script is visible from the root directory
			scriptName := strings.TrimSuffix(filepath.Base(s.testScriptPath()), filepath.Ext(s.testScriptPath()))
			output, err = session.EvaluateCode(ctx, "which "+scriptName)
			s.Require().NoError(err)
			s.Contains(output, s.testScriptPath(), "script should be found from the MCP root directory")

			// Step 4b: Verify projectPath parameter changes MATLAB's working directory
			tmpDir := s.T().TempDir()
			output, err = session.EvaluateCode(ctx, "pwd", tmpDir)
			s.Require().NoError(err)
			s.Contains(output, tmpDir, "projectPath should change MATLAB's working directory")

			output, err = session.EvaluateCode(ctx, "which "+scriptName)
			s.Require().NoError(err)
			s.Contains(output, "'"+scriptName+"' not found", "script should not be found after cd to temp dir")

			// Step 5: Code quality checking - analyze existing code for issues
			// First check code with problems to see what issues are detected
			problematicMessages, err := session.CheckCode(ctx, s.problematicCodePath())
			s.Require().NoError(err, "should check code without error")
			testdata.AssertProblematicCodeIssues(s.T(), problematicMessages)

			// Then check well-written code
			cleanMessages, err := session.CheckCode(ctx, s.testMathFunctionsPath())
			s.Require().NoError(err, "should check code without error")
			testdata.AssertCleanCode(s.T(), cleanMessages)

			// Step 6: Script execution - run a MATLAB script file
			scriptOutput, err := session.RunFile(ctx, s.testScriptPath())
			s.Require().NoError(err, "should execute script file without error")
			testdata.TestScript.Assert(s.T(), scriptOutput)

			// Step 7: Test execution - run test suite (TDD workflow)
			testOutput, err := session.RunTestFile(ctx, s.testMathFunctionsPath())
			s.Require().NoError(err, "should execute test suite without error")
			testdata.TestMathFunctions.Assert(s.T(), testOutput)
		})
	}
}

// TestParallelExperimentationWorkflow simulates a developer running isolated
// experiments in parallel MATLAB sessions. This is the typical use case for
// multi-session mode where the user needs independent workspaces: each session
// has its own variables, and results from one experiment don't affect another.
//
// Scenario: Data scientist comparing algorithm variants
// - Discovers available MATLAB installations
// - Starts multiple isolated sessions for parallel experimentation
// - Runs different computations in each session
// - Verifies results are isolated (no variable leakage between sessions)
// - Cleans up sessions when done
//
// MCP Tools tested:
// - list_available_matlabs (discover MATLAB installations)
// - start_matlab_session (create isolated session)
// - eval_in_matlab_session (run code in specific session)
// - stop_matlab_session (clean up session)
func (s *WorkflowTestSuite) TestParallelExperimentationWorkflow() {
	ctx := s.T().Context()
	mcpSession := s.CreateMCPSession(ctx, nil, nil, "--use-single-matlab-session=false")
	defer s.CleanupSession(mcpSession, true)

	sm := mcpSession.NewSessionManager()

	// Step 1: Discover available MATLAB installations
	matlabs, err := sm.ListAvailableMATLABs(ctx)
	s.Require().NoError(err, "should list available MATLABs")
	s.Require().NotEmpty(matlabs, "should have at least one MATLAB available")
	matlabRoot := matlabs[0].Path

	// Step 2: Start two independent sessions
	session1, err := sm.StartSession(ctx, matlabRoot)
	s.Require().NoError(err, "should start session 1")
	defer sm.CleanupSession(context.Background(), session1)

	session2, err := sm.StartSession(ctx, matlabRoot)
	s.Require().NoError(err, "should start session 2")
	defer sm.CleanupSession(context.Background(), session2)

	s.NotEqual(session1, session2, "sessions should have different IDs")

	// Step 3: Set different values in each session
	session1Value := 100
	session2Value := 200

	_, err = sm.EvaluateInSession(ctx, session1, fmt.Sprintf(`x = int32(%d);`, session1Value), s.testDataDir)
	s.Require().NoError(err)

	_, err = sm.EvaluateInSession(ctx, session2, fmt.Sprintf(`x = int32(%d);`, session2Value), s.testDataDir)
	s.Require().NoError(err)

	// Step 4: Verify isolation - each session has its own value
	output1, err := sm.EvaluateInSession(ctx, session1, `fprintf('x=%d\n', x);`, s.testDataDir)
	s.Require().NoError(err)
	s.Require().NotEmpty(output1)
	lines1 := strings.Split(output1, "\n")
	s.Contains(lines1, fmt.Sprintf("x=%d", session1Value), "session 1 should have its own value")

	output2, err := sm.EvaluateInSession(ctx, session2, `fprintf('x=%d\n', x);`, s.testDataDir)
	s.Require().NoError(err)
	s.Require().NotEmpty(output2)
	lines2 := strings.Split(output2, "\n")
	s.Contains(lines2, fmt.Sprintf("x=%d", session2Value), "session 2 should have its own value")

	// Step 5: Explicitly stop sessions (cleanup defers handle errors)
	err = sm.StopSession(ctx, session1)
	s.Require().NoError(err, "should stop session 1")

	err = sm.StopSession(ctx, session2)
	s.Require().NoError(err, "should stop session 2")
}

// TestInstallMATLABAddOnWorkflow verifies the --setup-matlab flag
// correctly installs the MATLAB MCP Core Server Toolbox add-on.
//
// Scenario: User installing the MATLAB add-on
// - Ensures the add-on is not installed (uninstalls if present)
// - Runs the MCP server with --setup-matlab
// - Verifies the add-on is installed
// - Always cleans up by uninstalling the add-on
//
// CLI flags tested:
// - --setup-matlab (one-shot add-on installation mode)
func (s *WorkflowTestSuite) TestSetupMATLABWorkflow() {
	ctx := s.T().Context()
	s.T().Setenv("PATH", s.pathEnvWithMATLAB)

	checkInstalledCode := strings.Join([]string{
		`tbxs = matlab.addons.toolbox.installedToolboxes();`,
		`if isempty(tbxs), installed = false;`,
		`else, installed = any(strcmp({tbxs.Name}, 'MATLAB MCP Core Server Toolbox')); end;`,
		`fprintf('installed=%d\n', installed);`,
	}, " ")

	uninstallCode := strings.Join([]string{
		`tbxs = matlab.addons.toolbox.installedToolboxes();`,
		`if ~isempty(tbxs),`,
		`idx = find(strcmp({tbxs.Name}, 'MATLAB MCP Core Server Toolbox'));`,
		`for i = 1:numel(idx), matlab.addons.toolbox.uninstallToolbox(tbxs(idx(i))); end;`,
		`end`,
	}, " ")

	runMATLAB := func(code string) (string, error) {
		s.T().Helper()
		cmd := exec.CommandContext(ctx, "matlab", "-batch", code)
		output, err := cmd.CombinedOutput()
		return string(output), err
	}

	// Step 1: Try to uninstall the add-on to start from a clean state
	output, err := runMATLAB(uninstallCode)
	s.Require().NoError(err, "should uninstall add-on if present:\n%s", output)

	output, err = runMATLAB(checkInstalledCode)
	s.Require().NoError(err, "should check add-on installation status:\n%s", output)
	wasPreInstalled := strings.Contains(output, "installed=1")

	if !wasPreInstalled {
		// Only uninstall on cleanup if we're the ones installing it
		defer func() {
			_, _ = runMATLAB(uninstallCode)
		}()
	}

	// Step 2: Run MCP server with --setup-matlab
	installCmd := exec.CommandContext(ctx, s.mcpServerPath, //nolint:gosec // Trusted path in tests
		"--setup-matlab",
	)
	installCmd.Env = s.defaultEnv
	installOutput, err := installCmd.CombinedOutput()
	s.Require().NoError(err, "setup-matlab should succeed:\n%s", string(installOutput))
	s.Contains(string(installOutput), "Successfully setup MATLAB.", "--setup-matlab should print success message")

	// Step 3: Verify the add-on is installed
	output, err = runMATLAB(checkInstalledCode)
	s.Require().NoError(err, "should check add-on installation status:\n%s", output)
	s.Contains(output, "installed=1", "add-on should be installed after running --setup-matlab")
}

// TestWorkflowSuite runs the workflow test suite
func TestWorkflowSuite(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}
