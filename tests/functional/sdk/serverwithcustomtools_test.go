// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/stretchr/testify/suite"
)

// ServerWithCustomToolsTestSuite tests SDK custom tools functionalities.
type ServerWithCustomToolsTestSuite struct {
	SDKTestSuite

	serverDetails testbinaries.ServerWithCustomToolsDetails
}

// SetupSuite runs once before all tests in a suite
func (s *ServerWithCustomToolsTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildServerWithCustomTools(s.T())
}

func TestServerWithCustomToolsTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWithCustomToolsTestSuite))
}

func (s *ServerWithCustomToolsTestSuite) TestSDK_CustomTools_HappyPath() {
	// Connect to a session
	session := s.CreateSession(s.serverDetails.BinaryLocation(), nil, nil)
	defer s.CleanupSession(session, true)

	name := "World"
	expectedTextOutput := "Hello " + name

	// Call the unstructured tool
	unstructuredResult, err := session.CallTool(s.T().Context(), s.serverDetails.GreetToolName(), map[string]any{"name": name})
	s.Require().NoError(err, "should call tool successfully")

	textContent, err := session.GetTextContent(unstructuredResult)
	s.Require().NoError(err, "should get text content")
	s.Require().Equal(expectedTextOutput, textContent, "should return greeting message")

	// Call the structured tool
	structuredResult, err := session.CallTool(s.T().Context(), s.serverDetails.GreetStructuredToolName(), map[string]any{"name": "World"})
	s.Require().NoError(err, "should call tool successfully")

	var output struct {
		Response string `json:"response"`
	}
	s.Require().NoError(session.UnmarshalStructuredContent(structuredResult, &output), "should unmarshal structured content")
	s.Require().Equal(expectedTextOutput, output.Response, "should return greeting message")
}
