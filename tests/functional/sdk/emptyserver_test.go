// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/stretchr/testify/suite"
)

// EmptyServerTestSuite tests SDK definition functionnalities, using a mock MATLAB.
type EmptyServerTestSuite struct {
	SDKTestSuite

	serverDetails testbinaries.ServerDetails
}

// SetupSuite runs once before all tests in a suite
func (s *EmptyServerTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildEmptyServer(s.T())
}

func TestEmptyServerTestSuite(t *testing.T) {
	suite.Run(t, new(EmptyServerTestSuite))
}

func (s *EmptyServerTestSuite) TestSDK_EmptyServer_Version() {
	// Arrange

	// Act
	cmd := exec.Command(s.serverDetails.BinaryLocation(), "--version") //nolint:gosec // Trusted test path
	output, err := cmd.CombinedOutput()

	// Assert
	s.Require().NoError(err, "version flag should execute successfully")
	s.Require().Contains(string(output), s.serverDetails.ModuleName(), "should display server package path")
}

func (s *EmptyServerTestSuite) TestSDK_EmptyServer_HasNoMATLABFlags() {
	for _, flag := range []string{
		"--matlab-root=/somewhere",
		"--initial-working-folder=/somewhere",
		"--use-single-matlab-session",
		"--initialize-matlab-on-startup",
		"--matlab-display-mode=desktop",
	} {
		s.Run(flag, func() {
			// Arrange

			// Act
			cmd := exec.Command(s.serverDetails.BinaryLocation(), flag) //nolint:gosec // Trusted test path
			output, err := cmd.CombinedOutput()

			// Assert
			s.Require().Error(err, "The flag shouldn't exist, and cause an error")
			s.Contains(string(output), "non-existent option")
		})
	}
}

func (s *EmptyServerTestSuite) TestSDK_EmptyServer_NameTitleAndInstructionNoToolsAndNoResources() {
	// Arrange
	session := s.CreateSession(s.serverDetails.BinaryLocation(), nil, nil)
	defer s.CleanupSession(session, true)

	// Act
	result := session.InitializeResult()

	listToolsResponse, err := session.ListTools(s.T().Context(), nil)
	s.Require().NoError(err)

	listResourcesResponse, err := session.ListResources(s.T().Context(), nil)
	s.Require().NoError(err)

	// Assert
	s.Require().NotNil(result)
	s.Equal(s.serverDetails.Name(), result.ServerInfo.Name)
	s.Equal(s.serverDetails.Title(), result.ServerInfo.Title)
	s.Equal(s.serverDetails.Instructions(), result.Instructions)

	s.Require().NotNil(listToolsResponse)
	s.Empty(listToolsResponse.Tools)

	s.Require().NotNil(listResourcesResponse)
	s.Empty(listResourcesResponse.Resources)

	// Verify client info is logged
	s.Eventually(func() bool {
		logContent, err := session.ReadServerLogs()
		if err != nil {
			return false
		}

		foundClientSessionLog := strings.Contains(logContent, "New client session")
		foundClientName := strings.Contains(logContent, session.ClientName())
		foundClientVersion := strings.Contains(logContent, session.ClientVersion())
		foundClientTitle := strings.Contains(logContent, session.ClientTitle())
		foundClientURL := strings.Contains(logContent, session.ClientWebsiteURL())

		return foundClientSessionLog &&
			foundClientName &&
			foundClientVersion &&
			foundClientTitle &&
			foundClientURL
	}, 2*time.Second, 200*time.Millisecond)
}
