// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabbuild_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/runmatlabbuild"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	runmatlabbuilduc "github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabbuild"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/multisession/runmatlabbuild"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockMATLABManager := &entitiesmocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	// Act
	tool := runmatlabbuild.New(mockLoggerFactory, mockUsecase, mockMATLABManager)

	// Assert
	assert.NotNil(t, tool)
}

func TestHandler_HappyPath_Success(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockMATLABManager := &entitiesmocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	const sessionID = 42
	args := runmatlabbuild.Args{
		SessionID:        sessionID,
		WorkingDirectory: "/home/user/myproject",
		Tasks:            []string{"check", "test"},
	}

	ucResult := runmatlabbuilduc.Result{
		ConsoleOutput: "** Done check\n** Done test\n",
		Success:       true,
	}

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), entities.SessionID(sessionID)).
		Return(mockClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockClient, runmatlabbuilduc.Args{
			WorkingDirectory: args.WorkingDirectory,
			Tasks:            args.Tasks,
		}).
		Return(ucResult, nil).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockMATLABManager)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, ucResult.ConsoleOutput, result.Log)
	assert.True(t, result.Success)
}

func TestHandler_HappyPath_BuildFailed(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockMATLABManager := &entitiesmocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	const sessionID = 1
	args := runmatlabbuild.Args{SessionID: sessionID}

	ucResult := runmatlabbuilduc.Result{
		ConsoleOutput: "** Starting test\n** Failed test\n",
		Success:       false,
	}

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), entities.SessionID(sessionID)).
		Return(mockClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockClient, runmatlabbuilduc.Args{}).
		Return(ucResult, nil).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockMATLABManager)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, ucResult.ConsoleOutput, result.Log)
	assert.False(t, result.Success)
}

func TestHandler_ManagerError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockMATLABManager := &entitiesmocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	const sessionID = 7
	args := runmatlabbuild.Args{SessionID: sessionID}
	expectedError := assert.AnError

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), entities.SessionID(sessionID)).
		Return(nil, expectedError).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockMATLABManager)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestHandler_UsecaseError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockMATLABManager := &entitiesmocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	const sessionID = 3
	args := runmatlabbuild.Args{SessionID: sessionID, WorkingDirectory: "/bad/path"}
	expectedError := assert.AnError

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), entities.SessionID(sessionID)).
		Return(mockClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockClient, runmatlabbuilduc.Args{
			WorkingDirectory: args.WorkingDirectory,
		}).
		Return(runmatlabbuilduc.Result{}, expectedError).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockMATLABManager)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestHandler_AllInputsPassedThrough(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockMATLABManager := &entitiesmocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	const sessionID = 5
	args := runmatlabbuild.Args{
		SessionID:         sessionID,
		WorkingDirectory:  "/home/user/project",
		Tasks:             []string{"check", "test"},
		ContinueOnFailure: true,
		Parallel:          true,
		Verbosity:         "verbose",
		Skip:              []string{"clean"},
	}

	expectedUCArgs := runmatlabbuilduc.Args{
		WorkingDirectory:  args.WorkingDirectory,
		Tasks:             args.Tasks,
		ContinueOnFailure: args.ContinueOnFailure,
		Parallel:          args.Parallel,
		Verbosity:         args.Verbosity,
		Skip:              args.Skip,
	}

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), entities.SessionID(sessionID)).
		Return(mockClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockClient, expectedUCArgs).
		Return(runmatlabbuilduc.Result{ConsoleOutput: "** Done\n", Success: true}, nil).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockMATLABManager)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestHandler_SessionIDPassedToManager(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockMATLABManager := &entitiesmocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	const sessionID = 99
	args := runmatlabbuild.Args{SessionID: sessionID}

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), entities.SessionID(sessionID)).
		Return(mockClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockClient, runmatlabbuilduc.Args{}).
		Return(runmatlabbuilduc.Result{Success: true}, nil).
		Once()

	// Act
	_, err := runmatlabbuild.Handler(mockUsecase, mockMATLABManager)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err)
}

func TestRunMATLABBuild_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockMATLABManager := &entitiesmocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	expectedAnnotations := annotations.NewDestructiveAnnotations()

	// Act
	tool := runmatlabbuild.New(mockLoggerFactory, mockUsecase, mockMATLABManager)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations())
}
