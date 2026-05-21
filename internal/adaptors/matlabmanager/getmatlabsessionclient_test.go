// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager_test

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager"
	sessionstoremocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabsessionstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMATLABManager_GetMATLABSessionClient_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockSessionClient := &sessionstoremocks.MockMATLABSessionClientWithCleanup{}
	defer mockSessionClient.AssertExpectations(t)

	expectedSessionID := entities.SessionID(123)
	ctx := t.Context()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionConnectionTimeout().
		Return(5 * time.Second).
		Once()

	mockSessionStore.EXPECT().
		Get(expectedSessionID).
		Return(mockSessionClient, nil).
		Once()

	mockSessionClient.EXPECT().
		Ping(mock.Anything, mockLogger.AsMockArg()).
		Return(entities.PingResponse{IsAlive: true}).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	client, err := manager.GetMATLABSessionClient(ctx, mockLogger, expectedSessionID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mockSessionClient, client)
}

func TestMATLABManager_GetMATLABSessionClient_Retries(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockConfigFactory := &mocks.MockConfigFactory{}
		defer mockConfigFactory.AssertExpectations(t)

		mockConfig := &configmocks.MockConfig{}
		defer mockConfig.AssertExpectations(t)

		mockMATLABServices := &mocks.MockMATLABServices{}
		defer mockMATLABServices.AssertExpectations(t)

		mockSessionStore := &mocks.MockMATLABSessionStore{}
		defer mockSessionStore.AssertExpectations(t)

		mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
		defer mockClientFactory.AssertExpectations(t)

		mockSessionSelector := &mocks.MockSessionSelector{}
		defer mockSessionSelector.AssertExpectations(t)

		mockSessionClient := &sessionstoremocks.MockMATLABSessionClientWithCleanup{}
		defer mockSessionClient.AssertExpectations(t)

		expectedSessionID := entities.SessionID(123)
		ctx := t.Context()
		retryInterval := 200 * time.Millisecond
		retryTimeout := 300 * time.Millisecond

		mockConfigFactory.EXPECT().
			Config().
			Return(mockConfig, nil).
			Once()

		mockConfig.EXPECT().
			MATLABSessionConnectionTimeout().
			Return(retryTimeout).
			Once()

		mockSessionStore.EXPECT().
			Get(expectedSessionID).
			Return(mockSessionClient, nil).
			Once()

		mockSessionClient.EXPECT().
			Ping(mock.Anything, mockLogger.AsMockArg()).
			Return(entities.PingResponse{IsAlive: false}).
			Once()

		mockSessionClient.EXPECT().
			Ping(mock.Anything, mockLogger.AsMockArg()).
			Return(entities.PingResponse{IsAlive: true}).
			Once()

		manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)
		manager.SetMATLABSessionConnectionRetryInterval(retryInterval)

		// Act
		client, err := manager.GetMATLABSessionClient(ctx, mockLogger, expectedSessionID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, mockSessionClient, client)
	})
}

func TestMATLABManager_GetMATLABSessionClient_RetryExhausted(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockConfigFactory := &mocks.MockConfigFactory{}
		defer mockConfigFactory.AssertExpectations(t)

		mockConfig := &configmocks.MockConfig{}
		defer mockConfig.AssertExpectations(t)

		mockMATLABServices := &mocks.MockMATLABServices{}
		defer mockMATLABServices.AssertExpectations(t)

		mockSessionStore := &mocks.MockMATLABSessionStore{}
		defer mockSessionStore.AssertExpectations(t)

		mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
		defer mockClientFactory.AssertExpectations(t)

		mockSessionSelector := &mocks.MockSessionSelector{}
		defer mockSessionSelector.AssertExpectations(t)

		mockSessionClient := &sessionstoremocks.MockMATLABSessionClientWithCleanup{}
		defer mockSessionClient.AssertExpectations(t)

		expectedSessionID := entities.SessionID(123)
		ctx := t.Context()
		retryInterval := 200 * time.Millisecond
		retryTimeout := 300 * time.Millisecond

		mockConfigFactory.EXPECT().
			Config().
			Return(mockConfig, nil).
			Once()

		mockConfig.EXPECT().
			MATLABSessionConnectionTimeout().
			Return(retryTimeout).
			Once()

		mockSessionStore.EXPECT().
			Get(expectedSessionID).
			Return(mockSessionClient, nil).
			Once()

		mockSessionClient.EXPECT().
			Ping(mock.Anything, mockLogger.AsMockArg()).
			Return(entities.PingResponse{IsAlive: false}).
			Twice()

		manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)
		manager.SetMATLABSessionConnectionRetryInterval(retryInterval)

		// Act
		client, err := manager.GetMATLABSessionClient(ctx, mockLogger, expectedSessionID)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not alive")
		assert.Nil(t, client)
	})
}

func TestMATLABManager_GetMATLABSessionClient_ConfigFactoryError(t *testing.T) {
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

	expectedSessionID := entities.SessionID(123)
	ctx := t.Context()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, messages.AnError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	client, err := manager.GetMATLABSessionClient(ctx, mockLogger, expectedSessionID)

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	assert.Nil(t, client)
}

func TestMATLABManager_GetMATLABSessionClient_SessionStoreError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	expectedSessionID := entities.SessionID(123)
	ctx := t.Context()
	expectedError := assert.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockSessionStore.EXPECT().
		Get(expectedSessionID).
		Return(nil, expectedError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Act
	client, err := manager.GetMATLABSessionClient(ctx, mockLogger, expectedSessionID)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}
