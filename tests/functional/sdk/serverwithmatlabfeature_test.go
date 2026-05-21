// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/stretchr/testify/suite"
)

// ServerWithMATLABFeatureTestSuite tests SDK MATLAB feature functionalities.
type ServerWithMATLABFeatureTestSuite struct {
	SDKTestSuite

	serverDetails testbinaries.ServerDetails
}

// SetupSuite runs once before all tests in a suite
func (s *ServerWithMATLABFeatureTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildServerWithMATLABFeature(s.T())
}

func TestServerWithMATLABFeatureTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWithMATLABFeatureTestSuite))
}

func (s *ServerWithMATLABFeatureTestSuite) TestSDK_MATLABFeature_HappyPath() {
	// Arrange
	session := s.CreateSession(s.serverDetails.BinaryLocation(), nil, nil)
	defer s.CleanupSession(session, true)

	// Act
	listToolsResponse, err := session.ListTools(s.T().Context(), nil)
	s.Require().NoError(err)

	listResourcesResponse, err := session.ListResources(s.T().Context(), nil)
	s.Require().NoError(err)

	// Assert
	s.Require().NotNil(listToolsResponse)
	s.Len(listToolsResponse.Tools, 5)

	s.Require().NotNil(listResourcesResponse)
	s.Len(listResourcesResponse.Resources, 2)
}
