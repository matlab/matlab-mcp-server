// Copyright 2026 The MathWorks, Inc.

package system_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-server/tests/testutils/matlabtempdir"
	"github.com/stretchr/testify/suite"
)

type InputWorkflowTestSuite struct {
	SystemTestSuite
}

// TestInteractiveInputWorkflow simulates a developer asking the AI to run code that
// asks the user for a value. The tool call returns while MATLAB waits at the prompt,
// and the AI answers the prompt with its next tool call.
//
// Scenario: Developer running interactive MATLAB code with AI assistance
// - Runs code that calls input; the tool call returns and shows the prompt
// - Answers with an expression, which runs in the waiting workspace
// - Answers a text prompt; the text reaches MATLAB exactly as sent
// - Runs a script that prompts; the tool call returns while MATLAB waits
// - Answers the script prompt; the script completes once
// - Answers a prompt with a project folder; the answer reaches the prompt and runs in the project folder
//
// MCP Tools tested:
// - evaluate_matlab_code (returns at a prompt; the next call answers it)
// - run_matlab_file (returns at a prompt)
//
// nodesktop is skipped until the tool call returns while MATLAB waits at a prompt in that mode.
func (s *InputWorkflowTestSuite) TestInteractiveInputWorkflow() {
	fixtureDir := filepath.Join(s.testDataDir, "input")

	for _, displayMode := range []string{"desktop", "nodesktop"} {
		s.Run(displayMode, func() {
			if displayMode == "nodesktop" {
				s.T().Skip("answering input prompts is not supported in nodesktop mode: the tool call does not return while MATLAB waits at an input prompt")
			}

			ctx, session, cleanup := s.startInteractiveSession(displayMode, fixtureDir)
			defer cleanup()
			s.requireMATLABStarted(ctx, session)

			// Step 1: Run code that asks for a value; the tool call returns at the prompt
			output, err := s.evaluateStep(ctx, session, `
				baseValue = 40;
				inputValue = input("MCP_INPUT_EXPRESSION> ");
				fprintf("MCP_INPUT_EXPRESSION_RESULT=%g\n", inputValue);
			`)
			s.requireWaitingForInput(output, err, "MCP_INPUT_EXPRESSION>")

			// Step 2: Answer with an expression; it runs in the waiting workspace, once
			result := s.requireEvaluate(ctx, session, "baseValue + 2")
			s.Contains(strings.Fields(result), "MCP_INPUT_EXPRESSION_RESULT=42", "the answer should be evaluated where input is waiting")
			s.Equal(1, strings.Count(output+result, "MCP_INPUT_EXPRESSION_RESULT="), "the code after input should run exactly once")

			// Step 3: Answer a text prompt; MATLAB stores the answer exactly as sent
			output, err = s.evaluateStep(ctx, session, `
				inputText = input("MCP_INPUT_STRING> ", "s");
				fprintf("MCP_INPUT_STRING_RESULT=%s\n", inputText);
			`)
			s.requireWaitingForInput(output, err, "MCP_INPUT_STRING>")

			result = s.requireEvaluate(ctx, session, "1+1")
			s.Contains(strings.Fields(result), "MCP_INPUT_STRING_RESULT=1+1", "string input should store the literal text")
			s.NotContains(strings.Fields(result), "MCP_INPUT_STRING_RESULT=2", "string input should not evaluate the answer")
			s.Equal(1, strings.Count(output+result, "MCP_INPUT_STRING_RESULT="), "the code after input should run exactly once")

			// Step 4: run_matlab_file on a script that prompts; the tool call returns at the prompt
			scriptPath := filepath.Join(fixtureDir, "promptInScript.m")
			output, err = s.boundedStep(ctx, "run_matlab_file promptInScript.m", func(stepCtx context.Context) (string, error) {
				return session.RunFile(stepCtx, scriptPath)
			})
			s.requireWaitingForInput(output, err, "MCP_INPUT_SCRIPT>")

			// Step 5: Answer the script prompt; the script completes once
			result = s.requireEvaluate(ctx, session, "scriptBaseValue + 2")
			s.Contains(strings.Fields(result), "MCP_INPUT_SCRIPT_RESULT=42", "the answer should be evaluated in the script workspace")
			s.Equal(1, strings.Count(output+result, "MCP_INPUT_SCRIPT_RESULT="), "the code after input should run exactly once")

			// Step 6: Answer a prompt with a project folder; the answer reaches the prompt and runs in the project folder
			output, err = s.evaluateStep(ctx, session, `
				projectInput = input("MCP_INPUT_PROJECT> ");
				fprintf("MCP_INPUT_PROJECT_RESULT=%s\n", projectInput);
			`)
			s.requireWaitingForInput(output, err, "MCP_INPUT_PROJECT>")

			projectDir := matlabtempdir.New(s.T(), "input-project-")
			result, err = s.boundedStep(ctx, "evaluate_matlab_code with project_path", func(stepCtx context.Context) (string, error) {
				return session.EvaluateCode(stepCtx, "pwd", projectDir)
			})
			s.Require().NoError(err, "answering with a project folder should succeed")
			s.Contains(result, "MCP_INPUT_PROJECT_RESULT="+projectDir, "the answer should reach the prompt and run in the project folder")
			s.Equal(1, strings.Count(output+result, "MCP_INPUT_PROJECT_RESULT="), "the code after input should run exactly once")
		})
	}
}

func (s *InputWorkflowTestSuite) requireWaitingForInput(output string, err error, expectedPrompt string) {
	s.T().Helper()
	s.Require().NoError(err, "waiting for input should not be reported as an error")
	s.Require().Contains(output, expectedPrompt, "the input prompt should be visible to the MCP client")
}

func TestInputWorkflowSuite(t *testing.T) {
	suite.Run(t, new(InputWorkflowTestSuite))
}
