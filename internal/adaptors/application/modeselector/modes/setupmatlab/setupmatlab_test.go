// Copyright 2026 The MathWorks, Inc.

package setupmatlab_test

import (
	"bytes"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/modeselector/modes/setupmatlab"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	setupmatlabmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/modeselector/modes/setupmatlab"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &setupmatlabmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMessageCatalog := &setupmatlabmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &setupmatlabmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockWatchdogClient := &setupmatlabmocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockGlobalMATLAB := &setupmatlabmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockAddonManager := &setupmatlabmocks.MockAddonManager{}
	defer mockAddonManager.AssertExpectations(t)

	// Act
	mode := setupmatlab.New(mockOSLayer, mockMessageCatalog, mockLoggerFactory, mockWatchdogClient, mockGlobalMATLAB, mockAddonManager)

	// Assert
	assert.NotNil(t, mode)
}

func TestMode_StartAndWaitForCompletion_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &setupmatlabmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMessageCatalog := &setupmatlabmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &setupmatlabmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockWatchdogClient := &setupmatlabmocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockGlobalMATLAB := &setupmatlabmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockAddonManager := &setupmatlabmocks.MockAddonManager{}
	defer mockAddonManager.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()
	successMessage := "Successfully installed MATLAB Add-On."

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(expectedCtx, mockLogger.AsMockArg()).
		Return(mockClient, nil).
		Once()

	mockAddonManager.EXPECT().
		Install(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(nil).
		Once()

	mockMessageCatalog.EXPECT().
		Get(messages.CLIMessages_SuccessfullySetupMATLAB).
		Return(successMessage).
		Once()

	mockOSLayer.EXPECT().
		Stdout().
		Return(&bytes.Buffer{}).
		Once()

	mode := setupmatlab.New(mockOSLayer, mockMessageCatalog, mockLoggerFactory, mockWatchdogClient, mockGlobalMATLAB, mockAddonManager)

	// Act
	err := mode.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.NoError(t, err)
}

func TestMode_StartAndWaitForCompletion_LoggerFactoryError(t *testing.T) {
	// Arrange
	mockOSLayer := &setupmatlabmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMessageCatalog := &setupmatlabmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &setupmatlabmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockWatchdogClient := &setupmatlabmocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockGlobalMATLAB := &setupmatlabmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockAddonManager := &setupmatlabmocks.MockAddonManager{}
	defer mockAddonManager.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, messages.AnError).
		Once()

	mode := setupmatlab.New(mockOSLayer, mockMessageCatalog, mockLoggerFactory, mockWatchdogClient, mockGlobalMATLAB, mockAddonManager)

	// Act
	err := mode.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, messages.AnError)
}

func TestMode_StartAndWaitForCompletion_WatchdogStartError(t *testing.T) {
	// Arrange
	mockOSLayer := &setupmatlabmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMessageCatalog := &setupmatlabmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &setupmatlabmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockWatchdogClient := &setupmatlabmocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockGlobalMATLAB := &setupmatlabmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockAddonManager := &setupmatlabmocks.MockAddonManager{}
	defer mockAddonManager.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(assert.AnError).
		Once()

	mode := setupmatlab.New(mockOSLayer, mockMessageCatalog, mockLoggerFactory, mockWatchdogClient, mockGlobalMATLAB, mockAddonManager)

	// Act
	err := mode.StartAndWaitForCompletion(t.Context())

	// Assert
	expectedError := messages.New_AddonManagerErrors_InstallFailed_Error()
	require.Equal(t, expectedError, err)
}

func TestMode_StartAndWaitForCompletion_MATLABClientError(t *testing.T) {
	// Arrange
	mockOSLayer := &setupmatlabmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMessageCatalog := &setupmatlabmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &setupmatlabmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockWatchdogClient := &setupmatlabmocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockGlobalMATLAB := &setupmatlabmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockAddonManager := &setupmatlabmocks.MockAddonManager{}
	defer mockAddonManager.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(expectedCtx, mockLogger.AsMockArg()).
		Return(nil, assert.AnError).
		Once()

	mode := setupmatlab.New(mockOSLayer, mockMessageCatalog, mockLoggerFactory, mockWatchdogClient, mockGlobalMATLAB, mockAddonManager)

	// Act
	err := mode.StartAndWaitForCompletion(expectedCtx)

	// Assert
	expectedError := messages.New_AddonManagerErrors_InstallFailed_Error()
	require.Equal(t, expectedError, err)
}

func TestMode_StartAndWaitForCompletion_AddonManagerInstallError(t *testing.T) {
	// Arrange
	mockOSLayer := &setupmatlabmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMessageCatalog := &setupmatlabmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &setupmatlabmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockWatchdogClient := &setupmatlabmocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockGlobalMATLAB := &setupmatlabmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockAddonManager := &setupmatlabmocks.MockAddonManager{}
	defer mockAddonManager.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(expectedCtx, mockLogger.AsMockArg()).
		Return(mockClient, nil).
		Once()

	mockAddonManager.EXPECT().
		Install(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(messages.AnError).
		Once()

	mode := setupmatlab.New(mockOSLayer, mockMessageCatalog, mockLoggerFactory, mockWatchdogClient, mockGlobalMATLAB, mockAddonManager)

	// Act
	err := mode.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.ErrorIs(t, err, messages.AnError)
}

func TestMode_StartAndWaitForCompletion_WatchdogStopError(t *testing.T) {
	// Arrange
	mockOSLayer := &setupmatlabmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMessageCatalog := &setupmatlabmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &setupmatlabmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockWatchdogClient := &setupmatlabmocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockGlobalMATLAB := &setupmatlabmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockAddonManager := &setupmatlabmocks.MockAddonManager{}
	defer mockAddonManager.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()
	successMessage := "Successfully installed MATLAB Add-On."

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(assert.AnError).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(expectedCtx, mockLogger.AsMockArg()).
		Return(mockClient, nil).
		Once()

	mockAddonManager.EXPECT().
		Install(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(nil).
		Once()

	mockMessageCatalog.EXPECT().
		Get(messages.CLIMessages_SuccessfullySetupMATLAB).
		Return(successMessage).
		Once()

	mockOSLayer.EXPECT().
		Stdout().
		Return(&bytes.Buffer{}).
		Once()

	mode := setupmatlab.New(mockOSLayer, mockMessageCatalog, mockLoggerFactory, mockWatchdogClient, mockGlobalMATLAB, mockAddonManager)

	// Act
	err := mode.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.NoError(t, err)
}
