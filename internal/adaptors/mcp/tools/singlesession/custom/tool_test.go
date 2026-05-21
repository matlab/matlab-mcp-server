// Copyright 2026 The MathWorks, Inc.

package custom_test

import (
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/definition"
	basetoolmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/custom/definition"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testToolName        = "test_tool"
	testToolTitle       = "Test Tool"
	testToolDescription = "A test tool"
)

func TestTool_AddToServer_HappyPath(t *testing.T) {
	// Arrange
	mockAdder := &basetoolmocks.MockToolAdder[map[string]any, any]{}
	defer mockAdder.AssertExpectations(t)

	expectedInputSchema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"n": {Type: "number", Description: "Size"},
		},
		Required: []string{"n"},
	}
	expectedAnnotations := &mcp.ToolAnnotations{
		ReadOnlyHint: true,
	}
	mockValidatedTool := &definitionmocks.MockValidatedTool{}
	defer mockValidatedTool.AssertExpectations(t)

	expectedDefinition := definition.Tool{
		Name:        testToolName,
		Title:       testToolTitle,
		Description: testToolDescription,
		InputSchema: expectedInputSchema,
		Annotations: expectedAnnotations,
	}
	expectedServer := mcp.NewServer(&mcp.Implementation{}, &mcp.ServerOptions{})

	mockValidatedTool.EXPECT().
		Definition().
		Return(expectedDefinition).
		Twice()
	mockValidatedTool.EXPECT().
		Signature().
		Return(definition.Signature{}).
		Once()

	mockAdder.EXPECT().
		AddTool(
			expectedServer,
			&mcp.Tool{
				Name:        testToolName,
				Title:       testToolTitle,
				Description: testToolDescription,
				Annotations: expectedAnnotations,
				InputSchema: expectedInputSchema,
			},
			mock.Anything,
		).
		Once()

	tool := custom.NewTool(mockValidatedTool, nil, nil, nil, nil)
	tool.SetToolAdder(mockAdder)

	// Act
	err := tool.AddToServer(expectedServer)

	// Assert
	require.NoError(t, err)
}
