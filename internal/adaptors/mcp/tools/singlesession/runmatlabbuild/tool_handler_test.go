// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabbuild_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabbuild"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	runmatlabbuilduc "github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabbuild"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/runmatlabbuild"
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

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	// Act
	tool := runmatlabbuild.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

	// Assert
	assert.NotNil(t, tool)
}

func TestHandler_HappyPath_Success(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	args := runmatlabbuild.Args{
		WorkingDirectory: "/home/user/myproject",
		Tasks:            []string{"check", "test"},
	}

	ucResult := runmatlabbuilduc.Result{
		ConsoleOutput: "** Done check\n** Done test\n",
		Success:       true,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
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
	result, err := runmatlabbuild.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, ucResult.ConsoleOutput, result.Log)
	assert.True(t, result.Success)
}

func TestHandler_HappyPath_BuildFailed(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	args := runmatlabbuild.Args{}

	ucResult := runmatlabbuilduc.Result{
		ConsoleOutput: "** Starting test\n** Failed test\n",
		Success:       false,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockClient, runmatlabbuilduc.Args{}).
		Return(ucResult, nil).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, ucResult.ConsoleOutput, result.Log)
	assert.False(t, result.Success)
}

func TestHandler_ClientError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	args := runmatlabbuild.Args{}
	expectedError := assert.AnError

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestHandler_UsecaseError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	args := runmatlabbuild.Args{WorkingDirectory: "/bad/path"}
	expectedError := assert.AnError

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockClient, runmatlabbuilduc.Args{
			WorkingDirectory: args.WorkingDirectory,
		}).
		Return(runmatlabbuilduc.Result{}, expectedError).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestHandler_AllInputsPassedThrough(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	args := runmatlabbuild.Args{
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

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockClient, expectedUCArgs).
		Return(runmatlabbuilduc.Result{ConsoleOutput: "** Done\n", Success: true}, nil).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestRunMATLABBuild_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	expectedAnnotations := annotations.NewDestructiveAnnotations()

	// Act
	tool := runmatlabbuild.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations())
}

func TestHandler_ReturnsEmptyLogWhenNoOutput(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	args := runmatlabbuild.Args{}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockClient, runmatlabbuilduc.Args{}).
		Return(runmatlabbuilduc.Result{ConsoleOutput: "", Success: true}, nil).
		Once()

	// Act
	result, err := runmatlabbuild.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, result.Log)
	assert.True(t, result.Success)
}
