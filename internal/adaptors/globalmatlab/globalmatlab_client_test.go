// Copyright 2025-2026 The MathWorks, Inc.

package globalmatlab_test

import (
	"sync"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/sessionmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalMATLAB_Client_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer expectedSessionClient.AssertExpectations(t)

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client, err := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_StartSessionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	ctx := t.Context()
	expectedError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(entities.SessionID(0), expectedError).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client, err := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	require.Nil(t, client)
}

func TestGlobalMATLAB_Client_ReturnsMATLABStartupCachedErrorOnSubsequentClientCalls(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	ctx := t.Context()
	expectedError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(entities.SessionID(0), expectedError).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client1, err1 := globalMATLAB.Client(ctx, mockLogger)
	client2, err2 := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.Nil(t, client1)
	require.ErrorIs(t, err1, expectedError)

	require.Nil(t, client2)
	require.ErrorIs(t, err2, expectedError)
}

func TestGlobalMATLAB_Client_DiscoveryErrorNotCached_RetrySucceeds(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer expectedSessionClient.AssertExpectations(t)

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(entities.SessionID(0), sessionmanager.ErrFailedToAttachToMATLABSession).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client1, err1 := globalMATLAB.Client(ctx, mockLogger)
	client2, err2 := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.Nil(t, client1)
	require.ErrorIs(t, err1, sessionmanager.ErrFailedToAttachToMATLABSession)

	require.Equal(t, expectedSessionClient, client2)
	require.NoError(t, err2)
}

func TestGlobalMATLAB_Client_DiscoveryErrorNotCached_RetryAlsoFails(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	ctx := t.Context()

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(entities.SessionID(0), sessionmanager.ErrFailedToAttachToMATLABSession).
		Twice()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client1, err1 := globalMATLAB.Client(ctx, mockLogger)
	client2, err2 := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.Nil(t, client1)
	require.ErrorIs(t, err1, sessionmanager.ErrFailedToAttachToMATLABSession)

	require.Nil(t, client2)
	require.ErrorIs(t, err2, sessionmanager.ErrFailedToAttachToMATLABSession)
}

func TestGlobalMATLAB_Client_GetMATLABSessionClientError_RetrySucceeds(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer expectedSessionClient.AssertExpectations(t)

	ctx := t.Context()
	firstSessionID := entities.SessionID(123)
	secondSessionID := entities.SessionID(456)
	getClientError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(firstSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(nil, getClientError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		ShouldRestart().
		Return(true, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(secondSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), secondSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client, err := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_RestartOnGetClientFailure(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	firstSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer firstSessionClient.AssertExpectations(t)

	secondSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer secondSessionClient.AssertExpectations(t)

	ctx := t.Context()
	firstSessionID := entities.SessionID(123)
	secondSessionID := entities.SessionID(456)
	getClientError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(firstSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(firstSessionClient, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(nil, getClientError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		ShouldRestart().
		Return(true, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(secondSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), secondSessionID).
		Return(secondSessionClient, nil).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	firstClient, firstErr := globalMATLAB.Client(ctx, mockLogger)
	secondClient, secondErr := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, firstErr)
	require.Equal(t, firstSessionClient, firstClient)
	require.NoError(t, secondErr)
	require.Equal(t, secondSessionClient, secondClient)
}

func TestGlobalMATLAB_Client_DoesNotErrorIfStopSessionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer expectedSessionClient.AssertExpectations(t)

	ctx := t.Context()
	firstSessionID := entities.SessionID(123)
	secondSessionID := entities.SessionID(456)
	getClientError := assert.AnError
	stopError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(firstSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(nil, getClientError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(stopError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		ShouldRestart().
		Return(true, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(secondSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), secondSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client, err := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_RestartFailure_OnExistingSession(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	ctx := t.Context()
	firstSessionID := entities.SessionID(123)
	getClientError := assert.AnError
	expectedError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(firstSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(nil, getClientError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		ShouldRestart().
		Return(true, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(entities.SessionID(0), expectedError).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client, err := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	require.Nil(t, client)
}

func TestGlobalMATLAB_Client_RestartDiscoveryErrorNotCached_RetrySucceeds(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer expectedSessionClient.AssertExpectations(t)

	ctx := t.Context()
	firstSessionID := entities.SessionID(123)
	secondSessionID := entities.SessionID(456)
	getClientError := assert.AnError

	// First call: start session succeeds, get client fails, restart returns discovery error
	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(firstSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(nil, getClientError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), firstSessionID).
		Return(nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		ShouldRestart().
		Return(true, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(entities.SessionID(0), sessionmanager.ErrFailedToAttachToMATLABSession).
		Once()

	// Second call: discovery succeeds on retry
	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(secondSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), secondSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client1, err1 := globalMATLAB.Client(ctx, mockLogger)
	client2, err2 := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.Nil(t, client1)
	require.ErrorIs(t, err1, sessionmanager.ErrFailedToAttachToMATLABSession)

	require.Equal(t, expectedSessionClient, client2)
	require.NoError(t, err2)
}

func TestGlobalMATLAB_Client_ConcurrentCallsWaitForCompletion(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer expectedSessionClient.AssertExpectations(t)

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Times(3)

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	var wg sync.WaitGroup
	results := make([]entities.MATLABSessionClient, 3)
	errs := make([]error, 3)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			client, err := globalMATLAB.Client(ctx, mockLogger)
			results[index] = client
			errs[index] = err
		}(i)
	}

	wg.Wait()

	// Assert
	for i := 0; i < 3; i++ {
		require.NoError(t, errs[i])
		require.Equal(t, expectedSessionClient, results[i])
	}
}

func TestGlobalMATLAB_Client_LostConnectionToSpecifiedMATLAB(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	getClientError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(sessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil, getClientError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		ShouldRestart().
		Return(false, nil).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client, err := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, globalmatlab.ErrLostMATLABConnection)
	require.Nil(t, client)
}

func TestGlobalMATLAB_Client_LostConnectionToSpecifiedMATLAB_CachedOnSubsequentCall(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	getClientError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(sessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil, getClientError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		ShouldRestart().
		Return(false, nil).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client1, err1 := globalMATLAB.Client(ctx, mockLogger)
	client2, err2 := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err1, globalmatlab.ErrLostMATLABConnection)
	require.Nil(t, client1)

	require.ErrorIs(t, err2, globalmatlab.ErrLostMATLABConnection)
	require.Nil(t, client2)
}

func TestGlobalMATLAB_Client_ShouldRestartError_CachedOnSubsequentCall(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	getClientError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(sessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil, getClientError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		ShouldRestart().
		Return(false, messages.AnError).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client1, err1 := globalMATLAB.Client(ctx, mockLogger)
	client2, err2 := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err1, messages.AnError)
	require.Nil(t, client1)

	require.ErrorIs(t, err2, messages.AnError)
	require.Nil(t, client2)
}

func TestGlobalMATLAB_Client_ShouldRestartError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	getClientError := assert.AnError

	mockMATLABManagerAdaptor.EXPECT().
		StartSession(ctx, mockLogger.AsMockArg()).
		Return(sessionID, nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil, getClientError).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil).
		Once()

	mockMATLABManagerAdaptor.EXPECT().
		ShouldRestart().
		Return(false, messages.AnError).
		Once()

	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Act
	client, err := globalMATLAB.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	require.Nil(t, client)
}
