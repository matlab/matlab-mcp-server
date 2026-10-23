// Copyright 2025-2026 The MathWorks, Inc.

package evalmatlabcode_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	"github.com/matlab/matlab-mcp-server/internal/usecases/evalmatlabcode"
	entitiesmocks "github.com/matlab/matlab-mcp-server/mocks/entities"
	mocks "github.com/matlab/matlab-mcp-server/mocks/usecases/evalmatlabcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	// Act
	usecase := evalmatlabcode.New(mockPathValidator)

	// Assert
	assert.NotNil(t, usecase, "Usecase should not be nil")
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	projectPath := filepath.Join("some", "path")
	validatedProjectPath := filepath.Join("some", "path")

	evalRequest := evalmatlabcode.Args{
		ProjectPath: projectPath,
		Code:        "disp('Hello, World!')",
	}

	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello, World!",
		Images:        nil,
	}

	ctx := t.Context()

	mockPathValidator.EXPECT().
		ValidateFolderPath(projectPath).
		Return(validatedProjectPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), entities.FEvalRequest{
			Function:   "cd",
			Arguments:  []string{validatedProjectPath},
			NumOutputs: 0,
		}).
		Return(entities.FEvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: evalRequest.Code}).
		Return(expectedResponse, nil).
		Once()

	usecase := evalmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, expectedResponse, response, "Response should match expected value")
}

func TestUsecase_Execute_SingleQuoteInPath_PassesPathUnchanged(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	projectPath := filepath.Join("Users", "O'Brien", "project")
	validatedProjectPath := projectPath

	evalRequest := evalmatlabcode.Args{
		ProjectPath: projectPath,
		Code:        "disp('Hello')",
	}

	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello",
	}

	mockPathValidator.EXPECT().
		ValidateFolderPath(projectPath).
		Return(validatedProjectPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), entities.FEvalRequest{
			Function:   "cd",
			Arguments:  []string{validatedProjectPath},
			NumOutputs: 0,
		}).
		Return(entities.FEvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: evalRequest.Code}).
		Return(expectedResponse, nil).
		Once()

	usecase := evalmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, expectedResponse, response, "Response should match expected value")
}

func TestUsecase_Execute_EmptyProjectPath_SkipsCd(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	evalRequest := evalmatlabcode.Args{
		ProjectPath: "",
		Code:        "disp('Hello, World!')",
	}

	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello, World!",
		Images:        nil,
	}

	ctx := t.Context()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: evalRequest.Code}).
		Return(expectedResponse, nil).
		Once()

	usecase := evalmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, expectedResponse, response, "Response should match expected value")
}

func TestUsecase_Execute_ValidatePathError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	projectPath := filepath.Join("some", "path")
	expectedError := assert.AnError

	evalRequest := evalmatlabcode.Args{
		ProjectPath: projectPath,
		Code:        "disp('Hello, World!')",
	}

	mockPathValidator.EXPECT().
		ValidateFolderPath(projectPath).
		Return("", expectedError).
		Once()

	usecase := evalmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalRequest)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the original error")
	assert.Empty(t, response, "Response should be empty when there's an error")
}

func TestUsecase_Execute_CDFEvalError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	projectPath := filepath.Join("some", "path")
	validatedProjectPath := filepath.Join("some", "path")

	evalRequest := evalmatlabcode.Args{
		Code:        "disp('Hello, World!')",
		ProjectPath: projectPath,
	}

	expectedError := assert.AnError

	ctx := t.Context()

	mockPathValidator.EXPECT().
		ValidateFolderPath(projectPath).
		Return(validatedProjectPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), entities.FEvalRequest{
			Function:   "cd",
			Arguments:  []string{validatedProjectPath},
			NumOutputs: 0,
		}).
		Return(entities.FEvalResponse{}, expectedError).
		Once()

	usecase := evalmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalRequest)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the original error")
	assert.Empty(t, response, "Response should be empty when there's an error")
}

func TestUsecase_Execute_EvalError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	projectPath := filepath.Join("some", "path")
	validatedProjectPath := filepath.Join("some", "path")

	evalRequest := evalmatlabcode.Args{
		Code:        "disp('Hello, World!')",
		ProjectPath: projectPath,
	}

	expectedError := assert.AnError

	ctx := t.Context()

	mockPathValidator.EXPECT().
		ValidateFolderPath(projectPath).
		Return(validatedProjectPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), entities.FEvalRequest{
			Function:   "cd",
			Arguments:  []string{validatedProjectPath},
			NumOutputs: 0,
		}).
		Return(entities.FEvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: evalRequest.Code}).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := evalmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalRequest)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the original error")
	assert.Empty(t, response, "Response should be empty when there's an error")
}

func TestUsecase_Execute_CaptureOutput_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	projectPath := filepath.Join("some", "path")
	validatedProjectPath := filepath.Join("some", "path")
	captureOutput := true

	evalRequest := evalmatlabcode.Args{
		ProjectPath:   projectPath,
		Code:          "disp('Hello, World!')",
		CaptureOutput: captureOutput,
	}

	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello, World!",
		Images:        nil,
	}

	mockPathValidator.EXPECT().
		ValidateFolderPath(projectPath).
		Return(validatedProjectPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), entities.FEvalRequest{
			Function:   "cd",
			Arguments:  []string{validatedProjectPath},
			NumOutputs: 0,
		}).
		Return(entities.FEvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: evalRequest.Code}).
		Return(expectedResponse, nil).
		Once()

	usecase := evalmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, expectedResponse, response, "Response should match expected value")
}

func TestUsecase_Execute_CaptureOutput_EvalWithCaptureError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	projectPath := filepath.Join("some", "path")
	validatedProjectPath := filepath.Join("some", "path")
	captureOutput := true
	expectedError := assert.AnError

	evalRequest := evalmatlabcode.Args{
		Code:          "disp('Hello, World!')",
		ProjectPath:   projectPath,
		CaptureOutput: captureOutput,
	}

	mockPathValidator.EXPECT().
		ValidateFolderPath(projectPath).
		Return(validatedProjectPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), entities.FEvalRequest{
			Function:   "cd",
			Arguments:  []string{validatedProjectPath},
			NumOutputs: 0,
		}).
		Return(entities.FEvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: evalRequest.Code}).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := evalmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalRequest)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the original error")
	assert.Empty(t, response, "Response should be empty when there's an error")
}
