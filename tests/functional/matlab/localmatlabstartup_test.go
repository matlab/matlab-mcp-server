// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/stretchr/testify/suite"
)

// LocalMATLABStartupTestSuite tests local MATLAB startup behavior through the MCP server
// using a mock MATLAB installation.
type LocalMATLABStartupTestSuite struct {
	MockMATLABTestSuite
}

func TestLocalMATLABStartupTestSuite(t *testing.T) {
	suite.Run(t, new(LocalMATLABStartupTestSuite))
}

func (s *LocalMATLABStartupTestSuite) TestHappyPath_EvaluateCode_ReturnsOutput() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig())
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	output, err := session.EvaluateCode(ctx, "disp('hello world')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('hello world')")
}

func (s *LocalMATLABStartupTestSuite) TestErrorPath_MATLABExitsImmediately_EvaluationFails() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.ExitImmediatelyConfig(1))
	s.Require().NoError(err)
	defer s.CleanupSession(session, false)

	_, err = session.EvaluateCode(ctx, "disp('should fail')", s.T().TempDir())
	s.Error(err, "evaluation should fail when MATLAB exits immediately")
}

func (s *LocalMATLABStartupTestSuite) TestErrorPath_MATLABStartupFailure_ReturnsStartupError() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.StartupFailureConfig())
	s.Require().NoError(err)
	defer s.CleanupSession(session, false)

	_, err = session.EvaluateCode(ctx, "disp('should fail')", s.T().TempDir())
	s.Require().Error(err, "evaluation should fail when MATLAB startup fails")
	errMsg := err.Error()
	s.Contains(errMsg, "MATLAB startup failed")
	s.Contains(errMsg, "Simulated MATLAB startup failure: license checkout failed")
}
