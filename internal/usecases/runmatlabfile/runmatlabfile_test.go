// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabfile_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	"github.com/matlab/matlab-mcp-server/internal/usecases/runmatlabfile"
	entitiesmocks "github.com/matlab/matlab-mcp-server/mocks/entities"
	mocks "github.com/matlab/matlab-mcp-server/mocks/usecases/runmatlabfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	// Act
	usecase := runmatlabfile.New(mockPathValidator)

	// Assert
	assert.NotNil(t, usecase, "Usecase should not be nil")
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	tests := []struct {
		name      string
		extension string
	}{
		{
			name:      "MATLAB source file",
			extension: ".m",
		},
		{
			name:      "P-code file",
			extension: ".p",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockLogger := testutils.NewInspectableLogger()

			mockPathValidator := &mocks.MockPathValidator{}
			defer mockPathValidator.AssertExpectations(t)

			mockClient := &entitiesmocks.MockMATLABSessionClient{}
			defer mockClient.AssertExpectations(t)

			ctx := t.Context()
			fileName := "file"
			scriptDir := filepath.Join("some", "path", "to")
			scriptPath := filepath.Join(scriptDir, fileName+tt.extension)

			usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

			expectedCdRequest := entities.FEvalRequest{
				Function:   "cd",
				Arguments:  []string{scriptDir},
				NumOutputs: 0,
			}

			expectedEvalRequest := entities.EvalRequest{
				Code: fileName,
			}

			expectedResponse := entities.EvalResponse{
				ConsoleOutput: "Hello, World!",
				Images:        nil,
			}

			mockPathValidator.EXPECT().
				ValidateMATLABCodeFile(scriptPath).
				Return(scriptPath, nil).
				Once()

			mockClient.EXPECT().
				FEval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
				Return(entities.FEvalResponse{}, nil).
				Once()

			mockClient.EXPECT().
				Eval(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
				Return(expectedResponse, nil).
				Once()

			usecase := runmatlabfile.New(mockPathValidator)

			// Act
			response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

			// Assert
			require.NoError(t, err, "Execute should not return an error")
			assert.Equal(t, expectedResponse, response, "Response should match expected value")
		})
	}
}

func TestUsecase_Execute_SingleQuoteInPath_PassesPathUnchanged(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "test"
	scriptDir := filepath.Join("Users", "O'Brien", "scripts")
	scriptPath := filepath.Join(scriptDir, fileName+".m")

	usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

	expectedCdRequest := entities.FEvalRequest{
		Function:   "cd",
		Arguments:  []string{scriptDir},
		NumOutputs: 0,
	}

	expectedEvalRequest := entities.EvalRequest{
		Code: fileName,
	}

	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello, World!",
		Images:        nil,
	}

	mockPathValidator.EXPECT().
		ValidateMATLABCodeFile(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(expectedResponse, nil).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, expectedResponse, response, "Response should match expected value")
}

func TestUsecase_Execute_InvalidMATLABFileName(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "Name with spaces",
			fileName: "my script",
		},
		{
			name:     "Name starting with a digit",
			fileName: "1script",
		},
		{
			name:     "Name with a hyphen",
			fileName: "my-script",
		},
		{
			name:     "Name with a dot",
			fileName: "my.script",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockLogger := testutils.NewInspectableLogger()

			mockPathValidator := &mocks.MockPathValidator{}
			defer mockPathValidator.AssertExpectations(t)

			mockClient := &entitiesmocks.MockMATLABSessionClient{}
			defer mockClient.AssertExpectations(t)

			ctx := t.Context()
			scriptDir := filepath.Join("some", "path", "to")
			scriptPath := filepath.Join(scriptDir, tt.fileName+".m")

			usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

			mockPathValidator.EXPECT().
				ValidateMATLABCodeFile(scriptPath).
				Return(scriptPath, nil).
				Once()

			usecase := runmatlabfile.New(mockPathValidator)

			// Act
			response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid MATLAB file name")
			assert.Empty(t, response, "Response should be empty")
		})
	}
}

func TestUsecase_Execute_ValidateMATLABCodeFileError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	expectedError := assert.AnError

	usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

	mockPathValidator.EXPECT().
		ValidateMATLABCodeFile(scriptPath).
		Return("", expectedError).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response, "Response should be empty")
}

func TestUsecase_Execute_CdFEvalError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	expectedError := assert.AnError

	usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

	expectedCdRequest := entities.FEvalRequest{
		Function:   "cd",
		Arguments:  []string{scriptDir},
		NumOutputs: 0,
	}

	mockPathValidator.EXPECT().
		ValidateMATLABCodeFile(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.FEvalResponse{}, expectedError).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response, "Response should be empty")
}

func TestUsecase_Execute_RunMATLABFileEvalError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	expectedError := assert.AnError

	usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

	expectedCdRequest := entities.FEvalRequest{
		Function:   "cd",
		Arguments:  []string{scriptDir},
		NumOutputs: 0,
	}

	expectedEvalRequest := entities.EvalRequest{
		Code: fileName,
	}

	mockPathValidator.EXPECT().
		ValidateMATLABCodeFile(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response, "Response should be empty")
}

func TestUsecase_Execute_CaptureOutput_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	captureOutput := true

	usecaseRequest := runmatlabfile.Args{
		ScriptPath:    scriptPath,
		CaptureOutput: captureOutput,
	}

	expectedCdRequest := entities.FEvalRequest{
		Function:   "cd",
		Arguments:  []string{scriptDir},
		NumOutputs: 0,
	}

	expectedEvalRequest := entities.EvalRequest{
		Code: fileName,
	}

	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello, World!",
		Images:        nil,
	}

	mockPathValidator.EXPECT().
		ValidateMATLABCodeFile(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(expectedResponse, nil).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

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
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	captureOutput := true
	expectedError := assert.AnError

	usecaseRequest := runmatlabfile.Args{
		ScriptPath:    scriptPath,
		CaptureOutput: captureOutput,
	}

	expectedCdRequest := entities.FEvalRequest{
		Function:   "cd",
		Arguments:  []string{scriptDir},
		NumOutputs: 0,
	}

	expectedEvalRequest := entities.EvalRequest{
		Code: fileName,
	}

	mockPathValidator.EXPECT().
		ValidateMATLABCodeFile(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response, "Response should be empty")
}
