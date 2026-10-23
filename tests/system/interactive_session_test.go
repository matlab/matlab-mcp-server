// Copyright 2026 The MathWorks, Inc.

package system_test

import (
	"context"
	"time"
)

const (
	matlabStartupTimeout       = 5 * time.Minute
	interactionStepTimeout     = 30 * time.Second
	interactiveWorkflowTimeout = 15 * time.Minute
)

func (s *SystemTestSuite) startInteractiveSession(displayMode string, workingFolder string) (context.Context, *SystemSession, func()) {
	t := s.T()
	t.Helper()
	ctx, abort := context.WithTimeout(t.Context(), interactiveWorkflowTimeout)
	session := s.CreateMCPSession(ctx, nil, nil,
		"--matlab-session-mode=new",
		"--matlab-display-mode="+displayMode,
		"--initial-working-folder="+workingFolder,
	)

	cleanup := func() {
		// A MATLAB stuck at a prompt can use up mcpclient's 3-minute graceful Close, so a
		// failed test kills the server instead and leaves MATLAB to the watchdog.
		if t.Failed() {
			abort()
		}
		s.CleanupSession(session, assertNoErrorLogs)
		abort()
	}
	return ctx, session, cleanup
}

func (s *SystemTestSuite) requireMATLABStarted(ctx context.Context, session *SystemSession) {
	s.T().Helper()
	startupCtx, cancel := context.WithTimeout(ctx, matlabStartupTimeout)
	defer cancel()

	output, err := session.EvaluateCode(startupCtx, `fprintf("MCP_MATLAB_STARTED\n");`)
	s.Require().NoError(err, "MATLAB should start within %s", matlabStartupTimeout)
	s.Require().Contains(output, "MCP_MATLAB_STARTED")
}

func (s *SystemTestSuite) boundedStep(ctx context.Context, description string, call func(context.Context) (string, error)) (string, error) {
	s.T().Helper()
	stepCtx, cancel := context.WithTimeout(ctx, interactionStepTimeout)
	defer cancel()

	output, err := call(stepCtx)
	s.Require().NoError(stepCtx.Err(), "%s should return within %s (call error: %v)", description, interactionStepTimeout, err)
	return output, err
}

func (s *SystemTestSuite) evaluateStep(ctx context.Context, session *SystemSession, code string) (string, error) {
	s.T().Helper()
	return s.boundedStep(ctx, "evaluate_matlab_code "+code, func(stepCtx context.Context) (string, error) {
		return session.EvaluateCode(stepCtx, code)
	})
}

func (s *SystemTestSuite) requireEvaluate(ctx context.Context, session *SystemSession, code string) string {
	s.T().Helper()
	output, err := s.evaluateStep(ctx, session, code)
	s.Require().NoError(err, "%s should succeed", code)
	return output
}

func (s *SystemTestSuite) requireEvaluateError(ctx context.Context, session *SystemSession, code string) string {
	s.T().Helper()
	_, err := s.evaluateStep(ctx, session, code)
	s.Require().Error(err, "%s should return a tool error", code)
	return err.Error()
}
