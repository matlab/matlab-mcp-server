// Copyright 2026 The MathWorks, Inc.

package custom_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	basetoolmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	custommocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/custom"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/custom/definition"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactory_LoadTools_HappyPath(t *testing.T) {
	// Arrange
	mockLoader := &custommocks.MockLoader{}
	defer mockLoader.AssertExpectations(t)

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &custommocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockConfigFactory := &custommocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedFilePath := "tools.json"

	mockValidatedTool1 := &definitionmocks.MockValidatedTool{}
	defer mockValidatedTool1.AssertExpectations(t)
	mockValidatedTool2 := &definitionmocks.MockValidatedTool{}
	defer mockValidatedTool2.AssertExpectations(t)

	validatedTools := []definition.ValidatedTool{mockValidatedTool1, mockValidatedTool2}

	mockValidatedTool1.EXPECT().
		Definition().
		Return(definition.Tool{Name: "tool1"}).
		Twice()
	mockValidatedTool1.EXPECT().
		Signature().
		Return(definition.Signature{}).
		Once()
	mockValidatedTool2.EXPECT().
		Definition().
		Return(definition.Tool{Name: "tool2"}).
		Twice()
	mockValidatedTool2.EXPECT().
		Signature().
		Return(definition.Signature{}).
		Once()

	mockLoader.EXPECT().
		Load(expectedFilePath).
		Return(validatedTools, nil).
		Once()

	factory := custom.NewFactory(mockLoader, mockLoggerFactory, mockUsecase, mockGlobalMATLAB, mockConfigFactory)

	// Act
	tools, err := factory.LoadTools(expectedFilePath)

	// Assert
	require.Nil(t, err)
	assert.Len(t, tools, 2)
	assert.Equal(t, "tool1", tools[0].Name())
	assert.Equal(t, "tool2", tools[1].Name())
}

func TestFactory_LoadTools_EmptyList(t *testing.T) {
	// Arrange
	mockLoader := &custommocks.MockLoader{}
	defer mockLoader.AssertExpectations(t)

	expectedFilePath := "tools.json"

	mockLoader.EXPECT().
		Load(expectedFilePath).
		Return([]definition.ValidatedTool{}, nil).
		Once()

	factory := custom.NewFactory(mockLoader, nil, nil, nil, nil)

	// Act
	tools, err := factory.LoadTools(expectedFilePath)

	// Assert
	require.Nil(t, err)
	assert.Empty(t, tools)
}

func TestFactory_LoadTools_LoaderError_ReturnsError(t *testing.T) {
	// Arrange
	mockLoader := &custommocks.MockLoader{}
	defer mockLoader.AssertExpectations(t)

	expectedFilePath := "tools.json"
	expectedError := messages.New_StartupErrors_FailedToParseExtensionFile_Error(expectedFilePath)

	mockLoader.EXPECT().
		Load(expectedFilePath).
		Return(nil, expectedError).
		Once()

	factory := custom.NewFactory(mockLoader, nil, nil, nil, nil)

	// Act
	tools, err := factory.LoadTools(expectedFilePath)

	// Assert
	assert.Nil(t, tools)
	require.Equal(t, expectedError, err)
}
