// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager"
	sessionstoremocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabsessionstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMATLABManager_StopMATLABSession_HappyPath(t *testing.T) {
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

	mockSessionClient := &sessionstoremocks.MockMATLABSessionClientWithCleanup{}
	defer mockSessionClient.AssertExpectations(t)

	expectedSessionID := entities.SessionID(123)
	ctx := t.Context()

	mockSessionStore.EXPECT().
		Get(expectedSessionID).
		Return(mockSessionClient, nil).
		Once()

	mockSessionClient.EXPECT().
		StopSession(ctx, mock.Anything).
		Return(nil).
		Once()

	mockSessionStore.EXPECT().
		Remove(expectedSessionID).
		Return().
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	err := manager.StopMATLABSession(ctx, mockLogger, expectedSessionID)

	// Assert
	require.NoError(t, err)
}

func TestMATLABManager_StopMATLABSession_SessionStoreGetError(t *testing.T) {
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

	expectedSessionID := entities.SessionID(123)
	ctx := t.Context()
	expectedError := assert.AnError

	mockSessionStore.EXPECT().
		Get(expectedSessionID).
		Return(nil, expectedError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	err := manager.StopMATLABSession(ctx, mockLogger, expectedSessionID)

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestMATLABManager_StopMATLABSession_StopSessionError(t *testing.T) {
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

	mockSessionClient := &sessionstoremocks.MockMATLABSessionClientWithCleanup{}
	defer mockSessionClient.AssertExpectations(t)

	expectedSessionID := entities.SessionID(123)
	ctx := t.Context()
	expectedError := assert.AnError

	mockSessionStore.EXPECT().
		Get(expectedSessionID).
		Return(mockSessionClient, nil).
		Once()

	mockSessionClient.EXPECT().
		StopSession(ctx, mock.Anything).
		Return(expectedError).
		Once()

	mockSessionStore.EXPECT().
		Remove(expectedSessionID).
		Return().
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	err := manager.StopMATLABSession(ctx, mockLogger, expectedSessionID)

	// Assert
	require.ErrorIs(t, err, expectedError)
}
