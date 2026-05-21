// Copyright 2026 The MathWorks, Inc.

package evalcustomtool_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalcustomtool"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalcustomtool/functioncall"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	evalcustomtoolmocks "github.com/matlab/matlab-mcp-core-server/mocks/usecases/evalcustomtool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockFunctionCallAssembler := &evalcustomtoolmocks.MockFunctionCallAssembler{}
	defer mockFunctionCallAssembler.AssertExpectations(t)

	// Act
	usecase := evalcustomtool.New(mockFunctionCallAssembler)

	// Assert
	assert.NotNil(t, usecase)
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockFunctionCallAssembler := &evalcustomtoolmocks.MockFunctionCallAssembler{}
	defer mockFunctionCallAssembler.AssertExpectations(t)

	expectedFunctionCallArgs := functioncall.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{"n": float64(5)},
	}
	code := "magic(5)"
	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "result",
	}

	ctx := t.Context()

	mockFunctionCallAssembler.EXPECT().
		Assemble(expectedFunctionCallArgs).
		Return(code, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: code}).
		Return(expectedResponse, nil).
		Once()

	usecase := evalcustomtool.New(mockFunctionCallAssembler)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalcustomtool.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{"n": float64(5)},
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestUsecase_Execute_CaptureOutput_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockFunctionCallAssembler := &evalcustomtoolmocks.MockFunctionCallAssembler{}
	defer mockFunctionCallAssembler.AssertExpectations(t)

	expectedFunctionCallArgs := functioncall.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{"n": float64(5)},
	}
	code := "magic(5)"
	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "result",
	}

	ctx := t.Context()

	mockFunctionCallAssembler.EXPECT().
		Assemble(expectedFunctionCallArgs).
		Return(code, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: code}).
		Return(expectedResponse, nil).
		Once()

	usecase := evalcustomtool.New(mockFunctionCallAssembler)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalcustomtool.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{"n": float64(5)},
		CaptureOutput: true,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestUsecase_Execute_EvalError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockFunctionCallAssembler := &evalcustomtoolmocks.MockFunctionCallAssembler{}
	defer mockFunctionCallAssembler.AssertExpectations(t)

	expectedFunctionCallArgs := functioncall.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{"n": float64(5)},
	}
	code := "magic(5)"
	expectedError := assert.AnError

	ctx := t.Context()

	mockFunctionCallAssembler.EXPECT().
		Assemble(expectedFunctionCallArgs).
		Return(code, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: code}).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := evalcustomtool.New(mockFunctionCallAssembler)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalcustomtool.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{"n": float64(5)},
	})

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response)
}

func TestUsecase_Execute_CaptureOutput_EvalWithCaptureError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockFunctionCallAssembler := &evalcustomtoolmocks.MockFunctionCallAssembler{}
	defer mockFunctionCallAssembler.AssertExpectations(t)

	expectedFunctionCallArgs := functioncall.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{"n": float64(5)},
	}
	code := "magic(5)"
	expectedError := assert.AnError

	ctx := t.Context()

	mockFunctionCallAssembler.EXPECT().
		Assemble(expectedFunctionCallArgs).
		Return(code, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), entities.EvalRequest{Code: code}).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := evalcustomtool.New(mockFunctionCallAssembler)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalcustomtool.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{"n": float64(5)},
		CaptureOutput: true,
	})

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response)
}

func TestUsecase_Execute_AssembleError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockFunctionCallAssembler := &evalcustomtoolmocks.MockFunctionCallAssembler{}
	defer mockFunctionCallAssembler.AssertExpectations(t)

	expectedFunctionCallArgs := functioncall.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{},
	}
	expectedError := assert.AnError

	ctx := t.Context()

	mockFunctionCallAssembler.EXPECT().
		Assemble(expectedFunctionCallArgs).
		Return("", expectedError).
		Once()

	usecase := evalcustomtool.New(mockFunctionCallAssembler)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, evalcustomtool.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{},
	})

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response)
}
