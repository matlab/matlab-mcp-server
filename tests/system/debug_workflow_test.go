// Copyright 2026 The MathWorks, Inc.

package system_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/tests/testutils/matlabtempdir"
	"github.com/matlab/matlab-mcp-server/tests/testutils/requirements"
	"github.com/stretchr/testify/suite"
)

type DebugWorkflowTestSuite struct {
	SystemTestSuite
}

// TestMATLABDebugWorkflow simulates a developer asking the AI to debug MATLAB code.
// The tool call returns while MATLAB is in debug mode, and the AI inspects
// and resumes MATLAB with its next tool calls.
//
// Scenario: Developer debugging MATLAB code with AI assistance
// - Sets a breakpoint and runs a function; the tool call returns where MATLAB stopped at the debug prompt
// - Reads a variable at the debug prompt, then continues; the function completes once
// - Runs the function again and quits the debugger; the call stopped at the debug prompt is abandoned
// - Stops at the debug prompt on an uncaught error; the tool call returns the error as a tool error
// - Stops at the debug prompt on a caught error and on a warning; the tool call returns them as output
// - Runs a file that calls keyboard; the tool call returns while MATLAB is at the keyboard prompt
// - Reads the file's workspace at the keyboard prompt with a project folder; the code runs in the project folder at the keyboard prompt
//
// MCP Tools tested:
// - evaluate_matlab_code (returns at the debug prompt; the next call runs in the workspace where MATLAB stopped)
// - run_matlab_file (returns at the keyboard prompt)
//
// nodesktop is skipped until, in that display mode, the tool call returns while MATLAB is in debug mode.
func (s *DebugWorkflowTestSuite) TestMATLABDebugWorkflow() {
	fixtureDir := filepath.Join(s.testDataDir, "debug")
	callLine := s.findLineNumber(filepath.Join(fixtureDir, "callLocalFunction.m"), "MCP_BREAK_CALL")

	for _, displayMode := range []string{"desktop", "nodesktop"} {
		s.Run(displayMode, func() {
			if displayMode == "nodesktop" {
				s.T().Skip("debugging is not supported in nodesktop mode: the tool call does not return while MATLAB is in debug mode")
			}

			ctx, session, cleanup := s.startInteractiveSession(displayMode, fixtureDir)
			defer cleanup()
			s.requireMATLABStarted(ctx, session)

			// Step 1: Set a breakpoint and run a function; the tool call returns where MATLAB stopped at the debug prompt
			s.requireEvaluate(ctx, session, fmt.Sprintf("dbstop in callLocalFunction at %d", callLine))

			status := s.requireEvaluate(ctx, session, `fprintf("MCP_BREAKPOINT_LINE=%d\n", dbstatus("callLocalFunction").line);`)
			s.Contains(strings.Fields(status), fmt.Sprintf("MCP_BREAKPOINT_LINE=%d", callLine), "dbstatus should report the breakpoint line")

			stopOutput := s.requireEvaluate(ctx, session, "callLocalFunction(5);")
			s.Contains(stopOutput, "MCP_BREAK_CALL", "the MCP client should see where MATLAB stopped")
			s.NotContains(strings.Fields(stopOutput), "MCP_DEBUG_RESULT=32", "the function should not complete before continuing")
			s.assertValue(ctx, session, "outerValue", "6")

			// Step 2: Continue; the function completes once
			result := s.requireEvaluate(ctx, session, "dbcont")
			s.Contains(strings.Fields(result), "MCP_DEBUG_RESULT=32", "dbcont should run the function to completion")
			s.Equal(1, strings.Count(stopOutput+result, "MCP_DEBUG_RESULT=32"), "the function should complete exactly once")
			s.assertNotInDebugMode(ctx, session)

			// Step 3: Run the function again and quit the debugger; the call stopped at the debug prompt is abandoned
			stopOutput = s.requireEvaluate(ctx, session, "callLocalFunction(5);")
			s.Contains(stopOutput, "MCP_BREAK_CALL", "the MCP client should see where MATLAB stopped")

			result = s.requireEvaluate(ctx, session, "dbquit")
			s.NotContains(strings.Fields(result), "MCP_DEBUG_RESULT=32", "dbquit should abandon the stopped call")
			s.assertNotInDebugMode(ctx, session)

			// Step 4: Stop at the debug prompt on an uncaught error; the tool call returns the error as a tool error
			s.requireEvaluate(ctx, session, "dbclear all")

			s.requireEvaluate(ctx, session, "dbstop if error MCP:ExpectedUncaught")

			stopOutput = s.requireEvaluateError(ctx, session, "throwUncaughtError();")
			s.Contains(stopOutput, "MCP_EXPECTED_UNCAUGHT state=7", "the MCP client should see the error that MATLAB stopped on")
			s.assertValue(ctx, session, "uncaughtState", "7")

			s.requireEvaluate(ctx, session, "dbcont")
			s.assertNotInDebugMode(ctx, session)

			// Step 5: Stop at the debug prompt on a caught error; the tool call returns the error as output
			s.requireEvaluate(ctx, session, "dbclear all")

			s.requireEvaluate(ctx, session, "dbstop if caught error MCP:ExpectedCaught")

			stopOutput = s.requireEvaluate(ctx, session, "throwAndCatchError();")
			s.Contains(stopOutput, "MCP_EXPECTED_CAUGHT", "the MCP client should see the error that MATLAB stopped on")
			s.NotContains(strings.Fields(stopOutput), "MCP_CAUGHT_HANDLED=MCP:ExpectedCaught", "the catch block should not run before continuing")
			s.assertValue(ctx, session, "caughtState", "8")

			result = s.requireEvaluate(ctx, session, "dbcont")
			s.Contains(strings.Fields(result), "MCP_CAUGHT_HANDLED=MCP:ExpectedCaught", "dbcont should let the catch block handle the error")
			s.assertNotInDebugMode(ctx, session)

			// Step 6: Stop at the debug prompt on a warning; the tool call returns the warning as output
			s.requireEvaluate(ctx, session, "dbclear all")

			s.requireEvaluate(ctx, session, "dbstop if warning MCP:ExpectedWarning")

			stopOutput = s.requireEvaluate(ctx, session, "issueWarning();")
			s.Contains(stopOutput, "MCP_EXPECTED_WARNING", "the MCP client should see the warning that MATLAB stopped on")
			s.NotContains(strings.Fields(stopOutput), "MCP_WARNING_COMPLETE", "the function should not complete before continuing")
			s.assertValue(ctx, session, "warningState", "9")

			result = s.requireEvaluate(ctx, session, "dbcont")
			s.Contains(strings.Fields(result), "MCP_WARNING_COMPLETE", "dbcont should run the function to completion")
			s.assertNotInDebugMode(ctx, session)

			// Step 7: run_matlab_file on a file that calls keyboard; the tool call returns while MATLAB is at the keyboard prompt
			s.requireEvaluate(ctx, session, "dbclear all")

			keyboardPath := filepath.Join(fixtureDir, "stopAtKeyboard.m")
			stopOutput, err := s.boundedStep(ctx, "run_matlab_file stopAtKeyboard.m", func(stepCtx context.Context) (string, error) {
				return session.RunFile(stepCtx, keyboardPath)
			})
			s.Require().NoError(err, "stopping at the keyboard prompt should not be reported as an error")
			s.NotContains(strings.Fields(stopOutput), "MCP_KEYBOARD_COMPLETE=11", "the file should not complete before continuing")
			s.assertValue(ctx, session, "keyboardState", "11")

			// Step 8: Read the file's workspace at the keyboard prompt with a project folder; the code runs in the project folder at the keyboard prompt
			projectDir := matlabtempdir.New(s.T(), "debug-project-")
			value, err := s.boundedStep(ctx, "evaluate_matlab_code with project_path", func(stepCtx context.Context) (string, error) {
				return session.EvaluateCode(stepCtx, `fprintf("MCP_VALUE=%g\nMCP_FOLDER=%s\n", keyboardState, pwd);`, projectDir)
			})
			s.Require().NoError(err, "reading the file's workspace at the keyboard prompt with a project folder should succeed")
			s.Contains(strings.Fields(value), "MCP_VALUE=11", "the code should run in the file's workspace at the keyboard prompt")
			s.Contains(value, "MCP_FOLDER="+projectDir, "the code should run in the project folder")

			// Step 9: Continue; the file completes once
			result = s.requireEvaluate(ctx, session, "dbcont")
			s.Contains(strings.Fields(result), "MCP_KEYBOARD_COMPLETE=11", "dbcont should resume the file after keyboard")
			s.Equal(1, strings.Count(stopOutput+result, "MCP_KEYBOARD_COMPLETE=11"), "the file should complete exactly once")
			s.assertNotInDebugMode(ctx, session)

			// Step 10: MATLAB accepts new code and no breakpoint remains
			output := s.requireEvaluate(ctx, session, `fprintf("MCP_DEBUG_WORKFLOW_HEALTHY\n");`)
			s.Contains(strings.Fields(output), "MCP_DEBUG_WORKFLOW_HEALTHY", "MATLAB should accept new code after all debug scenarios")

			breakpoints := s.requireEvaluate(ctx, session, `fprintf("MCP_BREAKPOINT_COUNT=%d\n", numel(dbstatus));`)
			s.Contains(strings.Fields(breakpoints), "MCP_BREAKPOINT_COUNT=0", "no breakpoint should remain after the workflow")
		})
	}
}

// TestQuitDebuggerAfterJavaError simulates a developer asking the AI to stop debugging
// after MATLAB stopped on a Java error.
//
// Scenario: Developer quits the debugger after a Java error
// - Stops at the debug prompt on a Java error; the tool call returns the error as a tool error
// - Quits all debugging; the same MATLAB keeps running
//
// MATLAB crashed on this dbquit all until initializeMCP initialized graphics at startup,
// so this test uses its own MATLAB session, where nothing else could initialize graphics first.
// The server starts a new MATLAB on the next call after a crash, so the test compares MATLAB process IDs.
//
// MCP Tools tested:
// - evaluate_matlab_code (returns at the debug prompt on a Java error)
//
// nodesktop is skipped until, in that display mode, the tool call returns while MATLAB is in debug mode.
// The test skips when the MATLAB under test has no JVM, unless MCP_SYSTEM_TEST_FORCE_JVM is set.
func (s *DebugWorkflowTestSuite) TestQuitDebuggerAfterJavaError() {
	fixtureDir := filepath.Join(s.testDataDir, "debug")

	for _, displayMode := range []string{"desktop", "nodesktop"} {
		s.Run(displayMode, func() {
			if displayMode == "nodesktop" {
				s.T().Skip("debugging is not supported in nodesktop mode: the tool call does not return while MATLAB is in debug mode")
			}

			ctx, session, cleanup := s.startInteractiveSession(displayMode, fixtureDir)
			defer cleanup()
			s.requireMATLABStarted(ctx, session)
			requirements.JVM(ctx, s.T(), session)

			// Step 1: Stop at the debug prompt on a Java error; the tool call returns the error as a tool error
			processID := s.requireMATLABProcessID(ctx, session)

			s.requireEvaluate(ctx, session, "dbstop if error")

			stopOutput := s.requireEvaluateError(ctx, session, "throwJavaException();")
			s.Contains(stopOutput, "java.net.ConnectException", "the MCP client should see the Java error that MATLAB stopped on")
			s.assertValue(ctx, session, "javaState", "12")

			// Step 2: Quit all debugging; MATLAB keeps running
			// MATLAB did not crash when dbquit all came right after MATLAB stopped, only when it came 10 seconds later.
			time.Sleep(10 * time.Second)
			s.requireEvaluate(ctx, session, "dbquit all")
			s.Equal(processID, s.requireMATLABProcessID(ctx, session), "the same MATLAB should keep running after dbquit all")
			s.assertNotInDebugMode(ctx, session)
		})
	}
}

func (s *DebugWorkflowTestSuite) findLineNumber(filePath string, marker string) int {
	s.T().Helper()
	content, err := os.ReadFile(filePath) //nolint:gosec // Test file path from the suite's test data folder
	s.Require().NoError(err)

	lineNumber := 0
	for i, line := range strings.Split(string(content), "\n") {
		if !strings.Contains(line, marker) {
			continue
		}
		s.Require().Zero(lineNumber, "%s should appear on only one line of %s", marker, filePath)
		lineNumber = i + 1
	}
	s.Require().NotZero(lineNumber, "%s should appear in %s", marker, filePath)
	return lineNumber
}

func (s *DebugWorkflowTestSuite) assertValue(ctx context.Context, session *SystemSession, variable string, expectedValue string) {
	s.T().Helper()
	output := s.requireEvaluate(ctx, session, fmt.Sprintf(`fprintf("MCP_VALUE=%%g\n", %s);`, variable))
	s.Contains(strings.Fields(output), "MCP_VALUE="+expectedValue, "%s should be %s in the workspace where MATLAB stopped", variable, expectedValue)
}

func (s *DebugWorkflowTestSuite) requireMATLABProcessID(ctx context.Context, session *SystemSession) string {
	s.T().Helper()
	output := s.requireEvaluate(ctx, session, `fprintf("MCP_MATLAB_PID=%d\n", feature("getpid"));`)
	for field := range strings.FieldsSeq(output) {
		if strings.HasPrefix(field, "MCP_MATLAB_PID=") {
			return field
		}
	}
	s.Require().Fail("the MCP client should see the MATLAB process ID", "output: %s", output)
	return ""
}

func (s *DebugWorkflowTestSuite) assertNotInDebugMode(ctx context.Context, session *SystemSession) {
	s.T().Helper()
	output := s.requireEvaluate(ctx, session, `fprintf("MCP_STACK_DEPTH=%d\n", numel(dbstack));`)
	s.Contains(strings.Fields(output), "MCP_STACK_DEPTH=0", "MATLAB should no longer be in debug mode")
}

func TestDebugWorkflowSuite(t *testing.T) {
	suite.Run(t, new(DebugWorkflowTestSuite))
}
