// Copyright 2025-2026 The MathWorks, Inc.

package basetool_test

import (
	"context"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestInput struct {
	Message string `json:"message"`
}

type TestOutput struct {
	Result string `json:"result"`
}

const (
	testToolName        = "test-tool"
	testToolTitle       = "Test Tool"
	testToolDescription = "A test tool for unit testing"
)

func TestNewToolWithStructuredContent_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "success"}, nil
	}

	// Act
	tool := basetool.NewToolWithStructuredContent(
		testToolName,
		testToolTitle,
		testToolDescription,
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	// Assert
	assert.Equal(t, testToolName, tool.Name(), "Tool name should match")
	assert.Equal(t, testToolTitle, tool.Title(), "Tool title should match")
	assert.Equal(t, testToolDescription, tool.Description(), "Tool description should match")

	expectedInputSchema, err := jsonschema.For[TestInput](&jsonschema.ForOptions{})
	require.NoError(t, err, "Input schema generation should succeed")
	inputSchema, err := tool.GetInputSchema()
	require.NoError(t, err, "Input schema generation should succeed")
	require.Equal(t, expectedInputSchema, inputSchema, "Input schema should not be nil")

	expectedOutputSchema, err := jsonschema.For[TestOutput](&jsonschema.ForOptions{})
	require.NoError(t, err, "Output schema generation should succeed")
	outputSchema, err := tool.GetOutputSchema()
	require.NoError(t, err, "Output schema generation should succeed")
	require.Equal(t, expectedOutputSchema, outputSchema, "Output schema should not be nil")
}

func TestToolWithStructuredContentOutput_AddToServer_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockAdder := &mocks.MockToolAdder[TestInput, TestOutput]{}
	defer mockAdder.AssertExpectations(t)

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "success"}, nil
	}

	expectedAnnotations := annotations.NewReadOnlyAnnotations()

	tool := basetool.NewToolWithStructuredContent(
		testToolName,
		testToolTitle,
		testToolDescription,
		expectedAnnotations,
		mockLoggerFactory,
		handler,
	)

	expectedToolInputSchema, err := tool.GetInputSchema()
	require.NoError(t, err, "GetInputSchema should not return an error")

	expectedToolOutputSchema, err := tool.GetOutputSchema()
	require.NoError(t, err, "GetOutputSchema should not return an error")

	expectedServer := mcp.NewServer(&mcp.Implementation{}, &mcp.ServerOptions{})

	mockAdder.EXPECT().AddTool(
		expectedServer,
		&mcp.Tool{
			Name:         testToolName,
			Title:        testToolTitle,
			Description:  testToolDescription,
			Annotations:  expectedAnnotations.ToToolAnnotations(),
			InputSchema:  expectedToolInputSchema,
			OutputSchema: expectedToolOutputSchema,
		},
		mock.Anything,
	)

	tool.SetToolAdder(mockAdder)

	// Act
	err = tool.AddToServer(expectedServer)

	// Assert
	require.NoError(t, err, "AddToServer should not return an error")
}

func TestToolWithStructuredContentOutput_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestInput{Message: "test message"}
	expectedOutput := TestOutput{Result: "processed: test message"}
	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "processed: " + input.Message}, nil
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithStructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Nil(t, result, "Result should be nil for structured content output")
	assert.Equal(t, expectedOutput, output, "Output should match expected output")
}

func TestToolWithStructuredContentOutput_Handler_StructuredHandlerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestInput{Message: "test message"}
	expectedError := assert.AnError
	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{}, expectedError
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithStructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return an error")
	assert.Nil(t, result, "Result should be nil when error occurs")
	assert.Empty(t, output, "Output should be zero value when error occurs")
}

func TestToolWithStructuredContentOutput_Handler_NewMCPSessionLoggerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestInput{Message: "test message"}
	expectedError := messages.AnError

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "should not be called"}, nil
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(nil, expectedError).
		Once()

	tool := basetool.NewToolWithStructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return the NewMCPSessionLogger error")
	assert.Nil(t, result, "Result should be nil when error occurs")
	assert.Empty(t, output, "Output should be zero value when error occurs")
}

func TestToolWithStructuredContentOutput_Handler_ContextPropagation(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestInput{Message: "test message"}
	expectedOutput := TestOutput{Result: "success"}
	mockSessionLogger := testutils.NewInspectableLogger()
	var capturedContext context.Context

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		capturedContext = ctx
		return expectedOutput, nil
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithStructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, _, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Equal(t, t.Context(), capturedContext, "Context should be propagated to handler")
}

func TestToolWithStructuredContent_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "success"}, nil
	}

	expectedAnnotations := annotations.NewDestructiveAnnotations()

	// Act
	tool := basetool.NewToolWithStructuredContent(
		"",
		"",
		"",
		expectedAnnotations,
		mockLoggerFactory,
		handler,
	)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have destructive annotations")
}

func TestToolWithStructuredAndTextContent_Handler_IncludesLogAsTextContent(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestInput{Message: "test message"}
	expectedOutput := TestOutput{Result: "processed: test message"}
	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "processed: " + input.Message}, nil
	}
	extractor := func(output TestOutput) string { return output.Result }

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithStructuredAndTextContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
		extractor,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, output, "Structured output should match")
	require.NotNil(t, result, "Result should not be nil when text extractor is set")
	require.Len(t, result.Content, 1, "Result should have one text content item")
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok, "Content should be TextContent")
	assert.Equal(t, expectedOutput.Result, textContent.Text, "Text content should be the extracted value")
}

func TestToolWithStructuredAndTextContent_Handler_Error_ReturnsNilResult(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestInput{Message: "test message"}
	expectedError := assert.AnError
	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{}, expectedError
	}
	extractor := func(output TestOutput) string { return output.Result }

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithStructuredAndTextContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
		extractor,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, result, "Result should be nil when error occurs")
	assert.Empty(t, output)
}

func TestToolWithStructuredContentOutput_AddToServer_NilAnnotationInterface(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "success"}, nil
	}

	tool := basetool.NewToolWithStructuredContent(
		testToolName,
		testToolTitle,
		testToolDescription,
		nil,
		mockLoggerFactory,
		handler,
	)

	// Act
	err := tool.AddToServer(nil)

	// Assert
	require.Error(t, err, "AddToServer should return an error for nil annotations")
	assert.Contains(t, err.Error(), "annotations must not be nil", "Error message should indicate nil annotations")
}
