// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabbuild_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabbuild"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/usecases/runmatlabbuild"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	disableHotLinks = "feature('HotLinks',0)"
	restoreHotLinks = "feature('HotLinks',1)"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	// Act
	usecase := runmatlabbuild.New(mockPathValidator)

	// Assert
	assert.NotNil(t, usecase, "Usecase should not be nil")
}

func TestExecute_NoWorkingDirectory_NoTasks(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	request := runmatlabbuild.Args{}

	expectedResponse := entities.EvalResponse{ConsoleOutput: "** Done test\n"}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool"}).
		Return(expectedResponse, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResponse.ConsoleOutput, result.ConsoleOutput)
	assert.True(t, result.Success)
}

func TestExecute_WithWorkingDirectory(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	workingDir := "/home/user/myproject"
	request := runmatlabbuild.Args{WorkingDirectory: workingDir}

	expectedResponse := entities.EvalResponse{ConsoleOutput: "** Done test\n"}

	mockPathValidator.EXPECT().
		ValidateFolderPath(workingDir).
		Return(workingDir, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: fmt.Sprintf("cd('%s')", workingDir)}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool"}).
		Return(expectedResponse, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResponse.ConsoleOutput, result.ConsoleOutput)
	assert.True(t, result.Success)
}

func TestExecute_WithWorkingDirectory_EscapesSingleQuotes(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	workingDir := "/home/O'Brien/myproject"
	escapedDir := strings.ReplaceAll(workingDir, "'", "''")
	request := runmatlabbuild.Args{WorkingDirectory: workingDir}

	mockPathValidator.EXPECT().
		ValidateFolderPath(workingDir).
		Return(workingDir, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: fmt.Sprintf("cd('%s')", escapedDir)}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool"}).
		Return(entities.EvalResponse{ConsoleOutput: "** Done\n"}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestExecute_WithTasks(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	request := runmatlabbuild.Args{Tasks: []string{"check", "test"}}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool check test"}).
		Return(entities.EvalResponse{ConsoleOutput: "** Done\n"}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestExecute_WithContinueOnFailure(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	request := runmatlabbuild.Args{ContinueOnFailure: true}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool -continueOnFailure"}).
		Return(entities.EvalResponse{ConsoleOutput: "** Done\n"}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestExecute_WithParallel(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	request := runmatlabbuild.Args{Parallel: true}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool -parallel"}).
		Return(entities.EvalResponse{ConsoleOutput: "** Done\n"}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestExecute_WithVerbosity(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	request := runmatlabbuild.Args{Verbosity: "verbose"}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool -verbosity verbose"}).
		Return(entities.EvalResponse{ConsoleOutput: "** Done\n"}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestExecute_WithSkipTasks(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	request := runmatlabbuild.Args{Skip: []string{"clean", "check"}}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool -skip clean -skip check"}).
		Return(entities.EvalResponse{ConsoleOutput: "** Done\n"}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestExecute_WithAllOptions(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	workingDir := "/home/user/myproject"
	request := runmatlabbuild.Args{
		WorkingDirectory:  workingDir,
		Tasks:             []string{"check", "test"},
		ContinueOnFailure: true,
		Parallel:          true,
		Verbosity:         "detailed",
		Skip:              []string{"clean"},
	}

	mockPathValidator.EXPECT().
		ValidateFolderPath(workingDir).
		Return(workingDir, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: fmt.Sprintf("cd('%s')", workingDir)}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{
			Code: "buildtool check test -continueOnFailure -parallel -verbosity detailed -skip clean",
		}).
		Return(entities.EvalResponse{ConsoleOutput: "** Done\n"}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestExecute_PathValidationError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	workingDir := "/nonexistent/path"
	request := runmatlabbuild.Args{WorkingDirectory: workingDir}
	expectedError := assert.AnError

	mockPathValidator.EXPECT().
		ValidateFolderPath(workingDir).
		Return("", expectedError).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.Error(t, err)
	assert.Empty(t, result)
}

func TestExecute_CdError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	workingDir := "/home/user/myproject"
	request := runmatlabbuild.Args{WorkingDirectory: workingDir}
	expectedError := assert.AnError

	mockPathValidator.EXPECT().
		ValidateFolderPath(workingDir).
		Return(workingDir, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: fmt.Sprintf("cd('%s')", workingDir)}).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestExecute_BuildtoolFailed_ReturnsSuccessFalse(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	request := runmatlabbuild.Args{}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool"}).
		Return(entities.EvalResponse{}, assert.AnError).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err, "Buildtool task failure should not propagate as an error")
	assert.Equal(t, assert.AnError.Error(), result.ConsoleOutput)
	assert.False(t, result.Success)
}

func TestExecute_HotLinksRestoredEvenOnBuildFailure(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	request := runmatlabbuild.Args{}

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: disableHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: "buildtool"}).
		Return(entities.EvalResponse{}, assert.AnError).
		Once()

	// Restore must be called even though buildtool failed
	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: restoreHotLinks}).
		Return(entities.EvalResponse{}, nil).
		Once()

	usecase := runmatlabbuild.New(mockPathValidator)

	// Act
	result, err := usecase.Execute(ctx, mockLogger, mockClient, request)

	// Assert
	require.NoError(t, err)
	assert.False(t, result.Success)
}
