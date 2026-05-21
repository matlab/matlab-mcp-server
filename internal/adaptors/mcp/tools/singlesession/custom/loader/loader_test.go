// Copyright 2026 The MathWorks, Inc.

package loader_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/loader"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/loader/validator"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/custom/definition"
	loadermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/custom/loader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/single_tool.json
var singleToolJSON []byte

//go:embed testdata/multiple_tools.json
var multipleToolsJSON []byte

//go:embed testdata/empty_tools.json
var emptyToolsJSON []byte

//go:embed testdata/malformed.json
var malformedJSON []byte

//go:embed testdata/invalid_property.json
var invalidPropertyJSON []byte

func TestNewLoader_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	// Act
	result := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Assert
	require.NotNil(t, result)
}

func TestLoader_Load_SingleTool_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	mockValidatedTool := &definitionmocks.MockValidatedTool{}
	defer mockValidatedTool.AssertExpectations(t)

	logger := testutils.NewInspectableLogger()
	toolsFilePath := filepath.Join("config", "tools.json")

	var parsed struct {
		Tools      []definition.Tool               `json:"tools"`
		Signatures map[string]definition.Signature `json:"signatures"`
	}
	require.NoError(t, json.Unmarshal(singleToolJSON, &parsed))
	require.Len(t, parsed.Tools, 1)
	expectedDefinition := parsed.Tools[0]
	expectedSignatures := parsed.Signatures

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(toolsFilePath).
		Return(singleToolJSON, nil).
		Once()
	mockToolValidator.EXPECT().
		Validate(expectedDefinition, expectedSignatures).
		Return(mockValidatedTool, nil).
		Once()
	mockValidatedTool.EXPECT().
		Definition().
		Return(expectedDefinition)

	l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Act
	tools, err := l.Load(toolsFilePath)

	// Assert
	require.NoError(t, err)
	require.Len(t, tools, 1)
	assert.Equal(t, expectedDefinition.Name, tools[0].Definition().Name)
	assert.Equal(t, expectedDefinition.Title, tools[0].Definition().Title)
	assert.Equal(t, expectedDefinition.Description, tools[0].Definition().Description)
}

func TestLoader_Load_MultipleTools_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	mockValidatedToolA := &definitionmocks.MockValidatedTool{}
	defer mockValidatedToolA.AssertExpectations(t)

	mockValidatedToolB := &definitionmocks.MockValidatedTool{}
	defer mockValidatedToolB.AssertExpectations(t)

	logger := testutils.NewInspectableLogger()
	toolsFilePath := filepath.Join("config", "tools.json")

	var parsed struct {
		Tools      []definition.Tool               `json:"tools"`
		Signatures map[string]definition.Signature `json:"signatures"`
	}
	require.NoError(t, json.Unmarshal(multipleToolsJSON, &parsed))
	require.Len(t, parsed.Tools, 2)
	expectedDefinitionA := parsed.Tools[0]
	expectedDefinitionB := parsed.Tools[1]
	expectedSignatures := parsed.Signatures

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()
	mockOSLayer.EXPECT().
		ReadFile(toolsFilePath).
		Return(multipleToolsJSON, nil).
		Once()
	mockToolValidator.EXPECT().
		Validate(expectedDefinitionA, expectedSignatures).
		Return(mockValidatedToolA, nil).
		Once()
	mockToolValidator.EXPECT().
		Validate(expectedDefinitionB, expectedSignatures).
		Return(mockValidatedToolB, nil).
		Once()
	mockValidatedToolA.EXPECT().
		Definition().
		Return(expectedDefinitionA)
	mockValidatedToolB.EXPECT().
		Definition().
		Return(expectedDefinitionB)

	l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Act
	tools, err := l.Load(toolsFilePath)

	// Assert
	require.NoError(t, err)
	require.Len(t, tools, 2)

	actualDefinitions := make([]definition.Tool, len(tools))
	for i, tool := range tools {
		actualDefinitions[i] = tool.Definition()
	}
	assert.ElementsMatch(t, []definition.Tool{expectedDefinitionA, expectedDefinitionB}, actualDefinitions)
}

func TestLoader_Load_EmptyToolsArray_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	logger := testutils.NewInspectableLogger()
	toolsFilePath := filepath.Join("config", "tools.json")
	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()
	mockOSLayer.EXPECT().
		ReadFile(toolsFilePath).
		Return(emptyToolsJSON, nil).
		Once()

	l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Act
	tools, err := l.Load(toolsFilePath)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, tools)
}

func TestLoader_Load_FileNotFound_ReturnsError(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	logger := testutils.NewInspectableLogger()
	toolsFilePath := filepath.Join("config", "tools.json")

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()
	mockOSLayer.EXPECT().
		ReadFile(toolsFilePath).
		Return(nil, os.ErrNotExist).
		Once()

	l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Act
	tools, err := l.Load(toolsFilePath)

	// Assert
	expectedError := messages.New_StartupErrors_FailedToReadExtensionFile_Error(toolsFilePath)

	assert.Nil(t, tools)
	require.Equal(t, expectedError, err)
}

func TestLoader_Load_ReadFileError_ReturnsError(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	logger := testutils.NewInspectableLogger()
	toolsFilePath := filepath.Join("config", "tools.json")
	readFileError := fmt.Errorf("permission denied")

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()
	mockOSLayer.EXPECT().
		ReadFile(toolsFilePath).
		Return(nil, readFileError).
		Once()

	l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Act
	tools, err := l.Load(toolsFilePath)

	// Assert
	expectedError := messages.New_StartupErrors_FailedToReadExtensionFile_Error(toolsFilePath)

	assert.Nil(t, tools)
	require.Equal(t, expectedError, err)
}

func TestLoader_Load_MalformedJSON_ReturnsError(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	logger := testutils.NewInspectableLogger()
	toolsFilePath := filepath.Join("config", "tools.json")
	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()
	mockOSLayer.EXPECT().
		ReadFile(toolsFilePath).
		Return(malformedJSON, nil).
		Once()

	l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Act
	tools, err := l.Load(toolsFilePath)

	// Assert
	expectedError := messages.New_StartupErrors_FailedToParseExtensionFile_Error(toolsFilePath)

	assert.Nil(t, tools)
	require.Equal(t, expectedError, err)
}

func TestLoader_Load_InvalidPropertyDefinition_ReturnsError(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	logger := testutils.NewInspectableLogger()
	toolsFilePath := filepath.Join("config", "tools.json")
	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()
	mockOSLayer.EXPECT().
		ReadFile(toolsFilePath).
		Return(invalidPropertyJSON, nil).
		Once()

	l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Act
	tools, err := l.Load(toolsFilePath)

	// Assert
	expectedError := messages.New_StartupErrors_FailedToParseExtensionFile_Error(toolsFilePath)

	assert.Nil(t, tools)
	require.Equal(t, expectedError, err)
}

func TestLoader_Load_ValidationError_ReturnsError(t *testing.T) {
	toolsFilePath := filepath.Join("config", "tools.json")
	expectedToolName := "test_tool"

	tests := []struct {
		name            string
		validationError error
		expectedError   messages.Error
	}{
		{
			"invalid tool definition",
			fmt.Errorf("missing required field: name: %w", validator.ErrInvalidToolDefinition),
			messages.New_StartupErrors_InvalidToolDefinition_Error(toolsFilePath),
		},
		{
			"invalid input schema",
			fmt.Errorf("inputSchema is required: %w", validator.ErrInvalidInputSchema),
			messages.New_StartupErrors_InvalidToolInputSchema_Error(expectedToolName, toolsFilePath),
		},
		{
			"missing signature",
			validator.ErrSignatureNotFound,
			messages.New_StartupErrors_MissingToolSignature_Error(expectedToolName, toolsFilePath),
		},
		{
			"invalid signature",
			fmt.Errorf("signature must have a 'function' field: %w", validator.ErrInvalidSignature),
			messages.New_StartupErrors_InvalidToolSignature_Error(expectedToolName, toolsFilePath),
		},
		{
			"unknown validation error",
			fmt.Errorf("validation failed"),
			messages.New_StartupErrors_InvalidToolDefinition_Error(toolsFilePath),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &loadermocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockLoggerFactory := &loadermocks.MockLoggerFactory{}
			defer mockLoggerFactory.AssertExpectations(t)

			mockToolValidator := &loadermocks.MockToolValidator{}
			defer mockToolValidator.AssertExpectations(t)

			logger := testutils.NewInspectableLogger()

			mockLoggerFactory.EXPECT().
				GetGlobalLogger().
				Return(logger, nil).
				Once()
			mockOSLayer.EXPECT().
				ReadFile(toolsFilePath).
				Return(singleToolJSON, nil).
				Once()
			mockToolValidator.EXPECT().
				Validate(mock.Anything, mock.Anything).
				Return(nil, tt.validationError).
				Once()

			l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

			// Act
			tools, err := l.Load(toolsFilePath)

			// Assert
			assert.Nil(t, tools)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestLoader_Load_DuplicateToolName_ReturnsError(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	mockDuplicateValidatedTool := &definitionmocks.MockValidatedTool{}
	defer mockDuplicateValidatedTool.AssertExpectations(t)

	logger := testutils.NewInspectableLogger()
	toolsFilePath := filepath.Join("config", "tools.json")
	expectedDefinition := definition.Tool{Name: "same_name"}

	mockDuplicateValidatedTool.EXPECT().
		Definition().
		Return(expectedDefinition)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()
	mockOSLayer.EXPECT().
		ReadFile(toolsFilePath).
		Return(multipleToolsJSON, nil).
		Once()
	mockToolValidator.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(mockDuplicateValidatedTool, nil).
		Twice()

	l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Act
	tools, err := l.Load(toolsFilePath)

	// Assert
	expectedError := messages.New_StartupErrors_DuplicateToolName_Error("same_name", toolsFilePath)

	assert.Nil(t, tools)
	require.Equal(t, expectedError, err)
}

func TestLoader_Load_LoggerFactoryError_ReturnsError(t *testing.T) {
	// Arrange
	mockOSLayer := &loadermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &loadermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockToolValidator := &loadermocks.MockToolValidator{}
	defer mockToolValidator.AssertExpectations(t)

	toolsFilePath := filepath.Join("config", "tools.json")
	expectedErr := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedErr).
		Once()

	l := loader.NewLoader(mockOSLayer, mockLoggerFactory, mockToolValidator)

	// Act
	tools, err := l.Load(toolsFilePath)

	// Assert
	require.ErrorIs(t, err, expectedErr)
	assert.Nil(t, tools)
}
