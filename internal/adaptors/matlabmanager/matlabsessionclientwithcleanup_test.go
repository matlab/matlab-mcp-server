// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMATLABSessionClientWithoutCleanup_HappyPath(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}

	// Act
	result := matlabmanager.NewMATLABSessionClientWithoutCleanup(mockClient)

	// Assert
	require.NotNil(t, result)
}

func TestMATLABSessionClientWithoutCleanup_StopSession_NoOp(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	client := matlabmanager.NewMATLABSessionClientWithoutCleanup(mockClient)

	// Act
	err := client.StopSession(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	_, hasDebugLog := mockLogger.DebugLogs()["Skipping session stop for externally managed MATLAB session"]
	assert.True(t, hasDebugLog, "should log that session stop was skipped")
}

func TestNewMATLABSessionClientWithCleanup_HappyPath(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}

	cleanupCalled := false
	cleanup := func() error {
		cleanupCalled = true
		return nil
	}

	// Act
	result := matlabmanager.NewMATLABSessionClientWithCleanup(mockClient, cleanup)

	// Assert
	require.NotNil(t, result)
	require.False(t, cleanupCalled, "Cleanup should not be called during construction")
}

func TestMATLABSessionClientWithCleanup_StopSession_HappyPath(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	cleanupCalled := false
	cleanup := func() error {
		cleanupCalled = true
		return nil
	}

	expectedEvalRequest := entities.EvalRequest{Code: "exit()"}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{}, nil).
		Once()

	client := matlabmanager.NewMATLABSessionClientWithCleanup(mockClient, cleanup)

	// Act
	err := client.StopSession(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	require.True(t, cleanupCalled, "Cleanup should be called")
}

func TestMATLABSessionClientWithCleanup_StopSession_EvalError(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	cleanupCalled := false
	cleanup := func() error {
		cleanupCalled = true
		return nil
	}

	expectedEvalRequest := entities.EvalRequest{Code: "exit()"}
	expectedError := assert.AnError

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	client := matlabmanager.NewMATLABSessionClientWithCleanup(mockClient, cleanup)

	// Act
	err := client.StopSession(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	require.False(t, cleanupCalled, "Cleanup should not be called when eval fails")
}

func TestMATLABSessionClientWithCleanup_StopSession_CleanupError(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	expectedCleanupError := assert.AnError
	cleanup := func() error {
		return expectedCleanupError
	}

	expectedEvalRequest := entities.EvalRequest{Code: "exit()"}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{}, nil).
		Once()

	client := matlabmanager.NewMATLABSessionClientWithCleanup(mockClient, cleanup)

	// Act
	err := client.StopSession(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedCleanupError)
}
