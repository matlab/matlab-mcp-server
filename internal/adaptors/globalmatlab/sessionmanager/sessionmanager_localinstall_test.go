// Copyright 2026 The MathWorks, Inc.

package sessionmanager_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/sessionmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab/sessionmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionManager_StartSession_LocalInstall_HappyPath(t *testing.T) {
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

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMATLABStartingDir(mockLogger.AsMockArg()).
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionMode().
		Return(entities.MATLABSessionModeNew).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(ctx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedSessionID, nil).
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
	require.NoError(t, err)
	require.Equal(t, expectedSessionID, sessionID)
}

func TestSessionManager_StartSession_LocalInstall_NoStartingDirectory(t *testing.T) {
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

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	shouldShowMATLABDesktop := false

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		StartingDirectory:      "",
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMATLABStartingDir(mockLogger.AsMockArg()).
		Return("", assert.AnError).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionMode().
		Return(entities.MATLABSessionModeNew).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(ctx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedSessionID, nil).
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
	require.NoError(t, err)
	require.Equal(t, expectedSessionID, sessionID)
}

func TestSessionManager_StartSession_LocalInstall_SelectMATLABRootError(t *testing.T) {
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

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	ctx := t.Context()
	expectedError := assert.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionMode().
		Return(entities.MATLABSessionModeNew).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(false).
		Once()

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return("", expectedError).
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
	require.ErrorIs(t, err, expectedError)
	require.Equal(t, entities.SessionID(0), sessionID)
}

func TestSessionManager_StartSession_LocalInstall_StartMATLABSessionError(t *testing.T) {
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

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	ctx := t.Context()
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	shouldShowMATLABDesktop := true
	expectedError := assert.AnError

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMATLABStartingDir(mockLogger.AsMockArg()).
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionMode().
		Return(entities.MATLABSessionModeNew).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(ctx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(entities.SessionID(0), expectedError).
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
	require.ErrorIs(t, err, expectedError)
	require.Equal(t, entities.SessionID(0), sessionID)
}
