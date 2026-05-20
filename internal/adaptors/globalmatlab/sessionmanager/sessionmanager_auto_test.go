// Copyright 2026 The MathWorks, Inc.

package sessionmanager_test

import (
	"context"
	"path/filepath"
	"testing"
	"testing/synctest"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/sessionmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab/sessionmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSessionManager_StartSession_AutoMode_DiscoverySucceeds(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
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

		expectedSessionID := entities.SessionID(789)
		discoveryTimeout := time.Duration(0)
		expectedAttachSessionDetails := entities.AttachToExistingSession{}

		mockConfigFactory.EXPECT().
			Config().
			Return(mockConfig, nil).
			Once()

		mockConfig.EXPECT().
			MATLABSessionMode().
			Return(entities.MATLABSessionModeAuto).
			Once()

		mockConfig.EXPECT().
			MATLABSessionDiscoveryTimeout().
			Return(discoveryTimeout).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedAttachSessionDetails).
			Return(expectedSessionID, nil).
			Once()

		starter := sessionmanager.New(
			mockMATLABManager,
			mockConfigFactory,
			mockMATLABRootSelector,
			mockMATLABStartingDirSelector,
		)

		// Act
		sessionID, err := starter.StartSession(t.Context(), mockLogger)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSessionID, sessionID)
	})
}

func TestSessionManager_StartSession_AutoMode_DiscoveryFails_FallsBackToLocalInstall(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
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
		expectedSessionID := entities.SessionID(456)
		discoveryTimeout := time.Duration(0)
		shouldShowMATLABDesktop := true
		expectedMATLABRoot := filepath.Join("some", "matlab", "root")
		expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")

		expectedAttachSessionDetails := entities.AttachToExistingSession{}
		expectedLocalSessionDetails := entities.LocalSessionDetails{
			MATLABRoot:             expectedMATLABRoot,
			IsStartingDirectorySet: true,
			StartingDirectory:      expectedMATLABStartingDir,
			ShowMATLABDesktop:      shouldShowMATLABDesktop,
		}

		mockConfigFactory.EXPECT().
			Config().
			Return(mockConfig, nil).
			Once()

		mockConfig.EXPECT().
			MATLABSessionMode().
			Return(entities.MATLABSessionModeAuto).
			Once()

		mockConfig.EXPECT().
			MATLABSessionDiscoveryTimeout().
			Return(discoveryTimeout).
			Once()

		mockConfig.EXPECT().
			ShouldShowMATLABDesktop().
			Return(shouldShowMATLABDesktop).
			Once()

		mockMATLABRootSelector.EXPECT().
			SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
			Return(expectedMATLABRoot, nil).
			Once()

		mockMATLABStartingDirSelector.EXPECT().
			SelectMATLABStartingDir(mockLogger.AsMockArg()).
			Return(expectedMATLABStartingDir, nil).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedAttachSessionDetails).
			Return(entities.SessionID(0), assert.AnError).
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
	})
}

func TestSessionManager_StartSession_AutoMode_BothFail(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
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
		discoveryTimeout := time.Duration(0)
		shouldShowMATLABDesktop := true
		expectedMATLABRoot := filepath.Join("some", "matlab", "root")
		expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")

		expectedAttachSessionDetails := entities.AttachToExistingSession{}
		expectedLocalSessionDetails := entities.LocalSessionDetails{
			MATLABRoot:             expectedMATLABRoot,
			IsStartingDirectorySet: true,
			StartingDirectory:      expectedMATLABStartingDir,
			ShowMATLABDesktop:      shouldShowMATLABDesktop,
		}

		mockConfigFactory.EXPECT().
			Config().
			Return(mockConfig, nil).
			Once()

		mockConfig.EXPECT().
			MATLABSessionMode().
			Return(entities.MATLABSessionModeAuto).
			Once()

		mockConfig.EXPECT().
			MATLABSessionDiscoveryTimeout().
			Return(discoveryTimeout).
			Once()

		mockConfig.EXPECT().
			ShouldShowMATLABDesktop().
			Return(shouldShowMATLABDesktop).
			Once()

		mockMATLABRootSelector.EXPECT().
			SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
			Return(expectedMATLABRoot, nil).
			Once()

		mockMATLABStartingDirSelector.EXPECT().
			SelectMATLABStartingDir(mockLogger.AsMockArg()).
			Return(expectedMATLABStartingDir, nil).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedAttachSessionDetails).
			Return(entities.SessionID(0), assert.AnError).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(ctx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
			Return(entities.SessionID(0), assert.AnError).
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
		require.ErrorIs(t, err, assert.AnError)
		require.Equal(t, entities.SessionID(0), sessionID)
	})
}

func TestSessionManager_StartSession_AutoMode_WithExplicitTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
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

		type contextKeyType string
		const contextKey contextKeyType = "testKey"
		const contextKeyValue = "testValue"

		ctx := context.WithValue(t.Context(), contextKey, contextKeyValue)
		expectedSessionID := entities.SessionID(321)
		retryInterval := 200 * time.Millisecond
		discoveryTimeout := 300 * time.Millisecond

		expectedAttachSessionDetails := entities.AttachToExistingSession{}

		mockConfigFactory.EXPECT().
			Config().
			Return(mockConfig, nil).
			Once()

		mockConfig.EXPECT().
			MATLABSessionMode().
			Return(entities.MATLABSessionModeAuto).
			Once()

		mockConfig.EXPECT().
			MATLABSessionDiscoveryTimeout().
			Return(discoveryTimeout).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(mock.MatchedBy(func(ctx context.Context) bool {
				return ctx.Value(contextKey) == contextKeyValue
			}), mockLogger.AsMockArg(), expectedAttachSessionDetails).
			Return(entities.SessionID(0), assert.AnError).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(mock.MatchedBy(func(ctx context.Context) bool {
				return ctx.Value(contextKey) == contextKeyValue
			}), mockLogger.AsMockArg(), expectedAttachSessionDetails).
			Return(expectedSessionID, nil).
			Once()

		starter := sessionmanager.New(
			mockMATLABManager,
			mockConfigFactory,
			mockMATLABRootSelector,
			mockMATLABStartingDirSelector,
		)
		starter.SetDiscoveryRetryInterval(retryInterval)

		// Act
		sessionID, err := starter.StartSession(ctx, mockLogger)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSessionID, sessionID)
	})
}

func TestSessionManager_ShouldRestart_AutoMode(t *testing.T) {
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
		Return(entities.MATLABSessionModeAuto).
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
	assert.True(t, shouldRestart)
}
