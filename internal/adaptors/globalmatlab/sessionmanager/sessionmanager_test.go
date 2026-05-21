// Copyright 2026 The MathWorks, Inc.

package sessionmanager_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/sessionmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab/sessionmanager"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	// Act
	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Assert
	assert.NotNil(t, starter)
}

func TestSessionManager_ShouldRestart_LocalInstallMode(t *testing.T) {
	// Arrange
	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionMode().
		Return(entities.MATLABSessionModeNew).
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	shouldRestart, err := starter.ShouldRestart()

	// Assert
	require.NoError(t, err)
	require.True(t, shouldRestart)
}

func TestSessionManager_ShouldRestart_AttachModeWithConnectionDetails(t *testing.T) {
	// Arrange
	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionMode().
		Return(entities.MATLABSessionModeExisting).
		Once()

	mockConfig.EXPECT().
		MATLABSessionConnectionDetails().
		Return("some-connection-details").
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	shouldRestart, err := starter.ShouldRestart()

	// Assert
	require.NoError(t, err)
	require.False(t, shouldRestart)
}

func TestSessionManager_ShouldRestart_AttachModeWithoutConnectionDetails(t *testing.T) {
	// Arrange
	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionMode().
		Return(entities.MATLABSessionModeExisting).
		Once()

	mockConfig.EXPECT().
		MATLABSessionConnectionDetails().
		Return("").
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	shouldRestart, err := starter.ShouldRestart()

	// Assert
	require.NoError(t, err)
	require.True(t, shouldRestart)
}

func TestSessionManager_ShouldRestart_ConfigError(t *testing.T) {
	// Arrange
	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, messages.AnError).
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	shouldRestart, err := starter.ShouldRestart()

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	require.False(t, shouldRestart)
}

func TestSessionManager_StopMATLABSession_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)

	mockMATLABManager.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(nil).
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	err := starter.StopMATLABSession(ctx, mockLogger, expectedSessionID)

	// Assert
	require.NoError(t, err)
}

func TestSessionManager_StopMATLABSession_Error(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedErr := assert.AnError

	mockMATLABManager.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedErr).
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	err := starter.StopMATLABSession(ctx, mockLogger, expectedSessionID)

	// Assert
	require.ErrorIs(t, err, expectedErr)
}

func TestSessionManager_GetMATLABSessionClient_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	client, err := starter.GetMATLABSessionClient(ctx, mockLogger, expectedSessionID)

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedSessionClient, client)
}

func TestSessionManager_GetMATLABSessionClient_Error(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedErr := assert.AnError

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(nil, expectedErr).
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	client, err := starter.GetMATLABSessionClient(ctx, mockLogger, expectedSessionID)

	// Assert
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, client)
}

func TestSessionManager_StartSession_ConfigError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	configErr := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, configErr).
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	sessionID, err := starter.StartSession(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, configErr)
	require.Equal(t, entities.SessionID(0), sessionID)
}
