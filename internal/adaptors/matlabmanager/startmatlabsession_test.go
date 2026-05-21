// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMATLABManager_StartMATLABSession_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedSessionID := entities.SessionID(123)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "1234",
	}

	sessionCleanupFunc := func() error { return nil }

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
	}

	expectedCtx := t.Context()

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.matlabSessionClientWithCleanup")).
		Return(expectedSessionID).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
	}

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionID, sessionID)
}

func TestMATLABManager_StartMATLABSession_MATLABServicesError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedError := assert.AnError

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
	}

	expectedCtx := t.Context()

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(embeddedconnector.ConnectionDetails{}, nil, expectedError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
	}

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, sessionID)
}

func TestMATLABManager_StartMATLABSession_ClientFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "12345",
	}
	cleanupCalled := false
	sessionCleanupFunc := func() error { cleanupCalled = true; return nil }
	expectedError := assert.AnError

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
	}

	expectedCtx := t.Context()

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(nil, expectedError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
	}

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, sessionID)
	assert.True(t, cleanupCalled, "session cleanup should be called when client factory fails")
}

func TestMATLABManager_StartMATLABSession_AttachToExistingSession_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockSessionClient.AssertExpectations(t)

	expectedSessionID := entities.SessionID(42)
	expectedConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "test-api-key",
		CertificatePEM: []byte("cert-content"),
	}
	expectedCtx := t.Context()

	mockSessionSelector.EXPECT().
		SelectSessionToAttachTo(mockLogger.AsMockArg()).
		Return(expectedConnectionDetails, nil).
		Once()

	mockClientFactory.EXPECT().
		New(expectedConnectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionClient.EXPECT().
		Ping(expectedCtx, mockLogger.AsMockArg()).
		Return(entities.PingResponse{IsAlive: true}).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.matlabSessionClientWithoutCleanup")).
		Return(expectedSessionID).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, entities.AttachToExistingSession{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionID, sessionID)
}

func TestMATLABManager_StartMATLABSession_AttachToExistingSession_SessionSelectorError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	expectedCtx := t.Context()

	mockSessionSelector.EXPECT().
		SelectSessionToAttachTo(mockLogger.AsMockArg()).
		Return(embeddedconnector.ConnectionDetails{}, assert.AnError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, entities.AttachToExistingSession{})

	// Assert
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, sessionID)
}

func TestMATLABManager_StartMATLABSession_AttachToExistingSession_ClientFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	expectedConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "key",
		CertificatePEM: []byte("cert"),
	}
	expectedCtx := t.Context()

	mockSessionSelector.EXPECT().
		SelectSessionToAttachTo(mockLogger.AsMockArg()).
		Return(expectedConnectionDetails, nil).
		Once()

	mockClientFactory.EXPECT().
		New(expectedConnectionDetails).
		Return(nil, assert.AnError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, entities.AttachToExistingSession{})

	// Assert
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, sessionID)
}

func TestMATLABManager_StartMATLABSession_AttachToExistingSession_PingFailure(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockSessionClient.AssertExpectations(t)

	expectedConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "test-api-key",
		CertificatePEM: []byte("cert-content"),
	}
	expectedCtx := t.Context()

	mockSessionSelector.EXPECT().
		SelectSessionToAttachTo(mockLogger.AsMockArg()).
		Return(expectedConnectionDetails, nil).
		Once()

	mockClientFactory.EXPECT().
		New(expectedConnectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionClient.EXPECT().
		Ping(expectedCtx, mockLogger.AsMockArg()).
		Return(entities.PingResponse{IsAlive: false}).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, entities.AttachToExistingSession{})

	// Assert
	require.ErrorIs(t, err, matlabmanager.ErrMATLABSessionNotAlive)
	assert.Empty(t, sessionID)
}
