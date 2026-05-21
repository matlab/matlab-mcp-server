// Copyright 2026 The MathWorks, Inc.

package addonmanager_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/addonmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	addonmanagermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/addonmanager"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockInstallationSteps := &addonmanagermocks.MockInstallationSteps{}

	// Act
	manager := addonmanager.New(mockInstallationSteps)

	// Assert
	assert.NotNil(t, manager)
}

func TestAddonManager_Install_HappyPath(t *testing.T) {
	// Arrange
	mockInstallationSteps := &addonmanagermocks.MockInstallationSteps{}
	defer mockInstallationSteps.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()
	cleanup := func() {}

	mockInstallationSteps.EXPECT().
		UploadMLTBX(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(cleanup, nil).
		Once()

	mockInstallationSteps.EXPECT().
		VerifyMLTBXInstallationFile(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(nil).
		Once()

	mockInstallationSteps.EXPECT().
		InstallMLTBX(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(nil).
		Once()

	manager := addonmanager.New(mockInstallationSteps)

	// Act
	err := manager.Install(expectedCtx, mockLogger, mockClient)

	// Assert
	require.NoError(t, err)
}

func TestAddonManager_Install_UploadMLTBXError(t *testing.T) {
	// Arrange
	mockInstallationSteps := &addonmanagermocks.MockInstallationSteps{}
	defer mockInstallationSteps.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()
	cleanup := func() {}

	mockInstallationSteps.EXPECT().
		UploadMLTBX(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(cleanup, assert.AnError).
		Once()

	manager := addonmanager.New(mockInstallationSteps)

	// Act
	err := manager.Install(expectedCtx, mockLogger, mockClient)

	// Assert
	require.ErrorIs(t, err, assert.AnError)
}

func TestAddonManager_Install_VerifyMLTBXInstallationFileError(t *testing.T) {
	// Arrange
	mockInstallationSteps := &addonmanagermocks.MockInstallationSteps{}
	defer mockInstallationSteps.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()
	cleanup := func() {}

	mockInstallationSteps.EXPECT().
		UploadMLTBX(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(cleanup, nil).
		Once()

	mockInstallationSteps.EXPECT().
		VerifyMLTBXInstallationFile(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(assert.AnError).
		Once()

	manager := addonmanager.New(mockInstallationSteps)

	// Act
	err := manager.Install(expectedCtx, mockLogger, mockClient)

	// Assert
	require.ErrorIs(t, err, assert.AnError)
}

func TestAddonManager_Install_InstallRetryExhausted(t *testing.T) {
	// Arrange
	mockInstallationSteps := &addonmanagermocks.MockInstallationSteps{}
	defer mockInstallationSteps.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()
	cleanup := func() {}

	mockInstallationSteps.EXPECT().
		UploadMLTBX(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(cleanup, nil).
		Once()

	mockInstallationSteps.EXPECT().
		VerifyMLTBXInstallationFile(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(nil).
		Once()

	mockInstallationSteps.EXPECT().
		InstallMLTBX(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(assert.AnError).
		Times(2)

	manager := addonmanager.New(mockInstallationSteps)

	// Act
	err := manager.Install(expectedCtx, mockLogger, mockClient)

	// Assert
	require.ErrorIs(t, err, assert.AnError)
}

func TestAddonManager_Install_InstallRetrySucceeds(t *testing.T) {
	// Arrange
	mockInstallationSteps := &addonmanagermocks.MockInstallationSteps{}
	defer mockInstallationSteps.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()
	cleanup := func() {}

	mockInstallationSteps.EXPECT().
		UploadMLTBX(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(cleanup, nil).
		Once()

	mockInstallationSteps.EXPECT().
		VerifyMLTBXInstallationFile(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(nil).
		Once()

	mockInstallationSteps.EXPECT().
		InstallMLTBX(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(assert.AnError).
		Once()

	mockInstallationSteps.EXPECT().
		InstallMLTBX(expectedCtx, mockLogger.AsMockArg(), mockClient).
		Return(nil).
		Once()

	manager := addonmanager.New(mockInstallationSteps)

	// Act
	err := manager.Install(expectedCtx, mockLogger, mockClient)

	// Assert
	require.NoError(t, err)
}
