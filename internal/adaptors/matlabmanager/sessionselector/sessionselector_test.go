// Copyright 2026 The MathWorks, Inc.

package sessionselector_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/sessionselector"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/sessionselector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockSessionDiscoverer := &mocks.MockSessionDiscoverer{}
	defer mockSessionDiscoverer.AssertExpectations(t)

	// Act
	attacher := sessionselector.New(mockConfigFactory, mockSessionDiscoverer)

	// Assert
	assert.NotNil(t, attacher)
}

func TestSessionSelector_SelectSessionToAttachTo_WithSessionDetails_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockSessionDiscoverer := &mocks.MockSessionDiscoverer{}
	defer mockSessionDiscoverer.AssertExpectations(t)

	sessionDetailsJSON := `{"port":31515,"certificate":"/path/to/cert.pem","apiKey":"test-api-key"}`
	expectedConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "test-api-key",
		CertificatePEM: []byte("cert-content"),
	}

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionConnectionDetails().
		Return(sessionDetailsJSON).
		Once()

	mockSessionDiscoverer.EXPECT().
		FromSessionDetails(mockLogger.AsMockArg(), []byte(sessionDetailsJSON)).
		Return(expectedConnectionDetails, nil).
		Once()

	attacher := sessionselector.New(mockConfigFactory, mockSessionDiscoverer)

	// Act
	connectionDetails, err := attacher.SelectSessionToAttachTo(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedConnectionDetails, connectionDetails)
}

func TestSessionSelector_SelectSessionToAttachTo_Discovery_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockSessionDiscoverer := &mocks.MockSessionDiscoverer{}
	defer mockSessionDiscoverer.AssertExpectations(t)

	expectedConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "test-api-key",
		CertificatePEM: []byte("cert-content"),
	}

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionConnectionDetails().
		Return("").
		Once()

	mockSessionDiscoverer.EXPECT().
		DiscoverSessions(mockLogger.AsMockArg()).
		Return([]embeddedconnector.ConnectionDetails{expectedConnectionDetails}).
		Once()

	attacher := sessionselector.New(mockConfigFactory, mockSessionDiscoverer)

	// Act
	connectionDetails, err := attacher.SelectSessionToAttachTo(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedConnectionDetails, connectionDetails)
}

func TestSessionSelector_SelectSessionToAttachTo_Discovery_MultipleSessionsUsesFirst(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockSessionDiscoverer := &mocks.MockSessionDiscoverer{}
	defer mockSessionDiscoverer.AssertExpectations(t)

	firstConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "first-api-key",
		CertificatePEM: []byte("first-cert"),
	}
	secondConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31516",
		APIKey:         "second-api-key",
		CertificatePEM: []byte("second-cert"),
	}
	expectedConnectionDetails := firstConnectionDetails

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionConnectionDetails().
		Return("").
		Once()

	mockSessionDiscoverer.EXPECT().
		DiscoverSessions(mockLogger.AsMockArg()).
		Return([]embeddedconnector.ConnectionDetails{firstConnectionDetails, secondConnectionDetails}).
		Once()

	attacher := sessionselector.New(mockConfigFactory, mockSessionDiscoverer)

	// Act
	connectionDetails, err := attacher.SelectSessionToAttachTo(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedConnectionDetails, connectionDetails)
}

func TestSessionSelector_SelectSessionToAttachTo_ConfigFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockSessionDiscoverer := &mocks.MockSessionDiscoverer{}
	defer mockSessionDiscoverer.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, messages.AnError).
		Once()

	attacher := sessionselector.New(mockConfigFactory, mockSessionDiscoverer)

	// Act
	connectionDetails, err := attacher.SelectSessionToAttachTo(mockLogger)

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	assert.Empty(t, connectionDetails)
}

func TestSessionSelector_SelectSessionToAttachTo_WithSessionDetails_FromSessionDetailsError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockSessionDiscoverer := &mocks.MockSessionDiscoverer{}
	defer mockSessionDiscoverer.AssertExpectations(t)

	sessionDetailsJSON := `invalid json`

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionConnectionDetails().
		Return(sessionDetailsJSON).
		Once()

	mockSessionDiscoverer.EXPECT().
		FromSessionDetails(mockLogger.AsMockArg(), []byte(sessionDetailsJSON)).
		Return(embeddedconnector.ConnectionDetails{}, assert.AnError).
		Once()

	attacher := sessionselector.New(mockConfigFactory, mockSessionDiscoverer)

	// Act
	connectionDetails, err := attacher.SelectSessionToAttachTo(mockLogger)

	// Assert
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, connectionDetails)
}

func TestSessionSelector_SelectSessionToAttachTo_Discovery_NoSessionsDiscovered(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockSessionDiscoverer := &mocks.MockSessionDiscoverer{}
	defer mockSessionDiscoverer.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionConnectionDetails().
		Return("").
		Once()

	mockSessionDiscoverer.EXPECT().
		DiscoverSessions(mockLogger.AsMockArg()).
		Return(nil).
		Once()

	attacher := sessionselector.New(mockConfigFactory, mockSessionDiscoverer)

	// Act
	connectionDetails, err := attacher.SelectSessionToAttachTo(mockLogger)

	// Assert
	require.ErrorIs(t, err, sessionselector.ErrNoMATLABSessionDiscovered)
	assert.Empty(t, connectionDetails)
}
