// Copyright 2026 The MathWorks, Inc.

package installationsteps_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/addonmanager/installationsteps"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Act
	steps := installationsteps.New()

	// Assert
	assert.NotNil(t, steps)
}

func TestInstallationSteps_UploadMLTBX_HappyPath(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()

	mockClient.EXPECT().
		Eval(expectedCtx, mockLogger.AsMockArg(), mock.Anything).
		Return(entities.EvalResponse{}, nil).
		Once()

	steps := installationsteps.New()

	// Act
	cleanup, err := steps.UploadMLTBX(expectedCtx, mockLogger, mockClient)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cleanup)
}

func TestInstallationSteps_UploadMLTBX_EvalError(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()

	mockClient.EXPECT().
		Eval(expectedCtx, mockLogger.AsMockArg(), mock.Anything).
		Return(entities.EvalResponse{}, assert.AnError).
		Once()

	steps := installationsteps.New()

	// Act
	cleanup, err := steps.UploadMLTBX(expectedCtx, mockLogger, mockClient)

	// Assert
	require.ErrorIs(t, err, assert.AnError)
	assert.NotNil(t, cleanup)
}

func TestInstallationSteps_VerifyMLTBXInstallationFile_HappyPath(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()

	mockClient.EXPECT().
		Eval(expectedCtx, mockLogger.AsMockArg(), mock.Anything).
		Return(entities.EvalResponse{}, nil).
		Once()

	steps := installationsteps.New()

	// Act
	err := steps.VerifyMLTBXInstallationFile(expectedCtx, mockLogger, mockClient)

	// Assert
	require.NoError(t, err)
}

func TestInstallationSteps_VerifyMLTBXInstallationFile_EvalError(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()

	mockClient.EXPECT().
		Eval(expectedCtx, mockLogger.AsMockArg(), mock.Anything).
		Return(entities.EvalResponse{}, assert.AnError).
		Once()

	steps := installationsteps.New()

	// Act
	err := steps.VerifyMLTBXInstallationFile(expectedCtx, mockLogger, mockClient)

	// Assert
	require.ErrorIs(t, err, assert.AnError)
}

func TestInstallationSteps_InstallMLTBX_HappyPath(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()

	mockClient.EXPECT().
		Eval(expectedCtx, mockLogger.AsMockArg(), mock.Anything).
		Return(entities.EvalResponse{}, nil).
		Once()

	steps := installationsteps.New()

	// Act
	err := steps.InstallMLTBX(expectedCtx, mockLogger, mockClient)

	// Assert
	require.NoError(t, err)
}

func TestInstallationSteps_InstallMLTBX_EvalError(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCtx := t.Context()

	mockClient.EXPECT().
		Eval(expectedCtx, mockLogger.AsMockArg(), mock.Anything).
		Return(entities.EvalResponse{}, assert.AnError).
		Once()

	steps := installationsteps.New()

	// Act
	err := steps.InstallMLTBX(expectedCtx, mockLogger, mockClient)

	// Assert
	require.ErrorIs(t, err, assert.AnError)
}
