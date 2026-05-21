// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/stretchr/testify/suite"
)

// ServerWithCustomDependenciesTestSuite tests SDK custom dependencies functionalities.
type ServerWithCustomDependenciesTestSuite struct {
	SDKTestSuite

	serverDetails testbinaries.ServerWithCustomDependenciesDetails
}

// SetupSuite runs once before all tests in a suite
func (s *ServerWithCustomDependenciesTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildServerWithCustomDependencies(s.T())
}

func TestServerWithCustomDependenciesTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWithCustomDependenciesTestSuite))
}

func (s *ServerWithCustomDependenciesTestSuite) TestSDK_CustomDependencies_ToolUsesDependency_HappyPath() {
	// Connect to a session
	session := s.CreateSession(s.serverDetails.BinaryLocation(), nil, nil)
	defer s.CleanupSession(session, true)

	name := "World"
	expectedTextOutput := "Service Hello " + name

	// Call the tool
	result, err := session.CallTool(s.T().Context(), s.serverDetails.GreetToolName(), map[string]any{"name": name})
	s.Require().NoError(err, "should call tool successfully")

	textContent, err := session.GetTextContent(result)
	s.Require().NoError(err, "should get text content")
	s.Require().Equal(expectedTextOutput, textContent, "should return greeting message using dependency")
}
