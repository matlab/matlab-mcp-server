// Copyright 2025-2026 The MathWorks, Inc.

package responseconverter_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/utils/responseconverter"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertEvalResponseToRichContent_HappyPath(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		response entities.EvalResponse
		expected tools.RichContent
	}{
		{
			name: "EmptyResponse",
			response: entities.EvalResponse{
				ConsoleOutput: "",
				Images:        [][]byte{},
			},
			expected: tools.RichContent{
				TextContent:  []string{""},
				ImageContent: []tools.PNGImageData{},
			},
		},
		{
			name: "ConsoleOutputOnly",
			response: entities.EvalResponse{
				ConsoleOutput: "Hello World",
				Images:        [][]byte{},
			},
			expected: tools.RichContent{
				TextContent:  []string{"Hello World"},
				ImageContent: []tools.PNGImageData{},
			},
		},
		{
			name: "ImagesOnly",
			response: entities.EvalResponse{
				ConsoleOutput: "",
				Images:        [][]byte{[]byte("image1"), []byte("image2")},
			},
			expected: tools.RichContent{
				TextContent:  []string{""},
				ImageContent: []tools.PNGImageData{tools.PNGImageData("image1"), tools.PNGImageData("image2")},
			},
		},
		{
			name: "BothConsoleOutputAndImages",
			response: entities.EvalResponse{
				ConsoleOutput: "Processing complete",
				Images:        [][]byte{[]byte("chart")},
			},
			expected: tools.RichContent{
				TextContent:  []string{"Processing complete"},
				ImageContent: []tools.PNGImageData{tools.PNGImageData("chart")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := responseconverter.ConvertEvalResponseToRichContent(tt.response)

			// Assert
			assert.Equal(t, tt.expected.TextContent, result.TextContent, "TextContent should match expected value")
			assert.Equal(t, tt.expected.ImageContent, result.ImageContent, "ImageContent should match expected value")
		})
	}
}

func TestConvertRichContentToCallToolResult_HappyPath(t *testing.T) {
	// Arrange
	tests := []struct {
		name            string
		content         tools.RichContent
		expectedContent []mcp.Content
	}{
		{
			name: "EmptyContent",
			content: tools.RichContent{
				TextContent:  []string{},
				ImageContent: []tools.PNGImageData{},
			},
			expectedContent: []mcp.Content{},
		},
		{
			name: "TextContentOnly",
			content: tools.RichContent{
				TextContent:  []string{"Hello World"},
				ImageContent: []tools.PNGImageData{},
			},
			expectedContent: []mcp.Content{
				&mcp.TextContent{Text: "Hello World"},
			},
		},
		{
			name: "ImageContentOnly",
			content: tools.RichContent{
				TextContent:  []string{},
				ImageContent: []tools.PNGImageData{tools.PNGImageData("image1"), tools.PNGImageData("image2")},
			},
			expectedContent: []mcp.Content{
				&mcp.ImageContent{MIMEType: "image/png", Data: []byte("image1")},
				&mcp.ImageContent{MIMEType: "image/png", Data: []byte("image2")},
			},
		},
		{
			name: "BothTextAndImageContent",
			content: tools.RichContent{
				TextContent:  []string{"Processing complete"},
				ImageContent: []tools.PNGImageData{tools.PNGImageData("chart")},
			},
			expectedContent: []mcp.Content{
				&mcp.TextContent{Text: "Processing complete"},
				&mcp.ImageContent{MIMEType: "image/png", Data: []byte("chart")},
			},
		},
		{
			name: "MultipleTextEntries",
			content: tools.RichContent{
				TextContent:  []string{"line1", "line2"},
				ImageContent: []tools.PNGImageData{},
			},
			expectedContent: []mcp.Content{
				&mcp.TextContent{Text: "line1"},
				&mcp.TextContent{Text: "line2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := responseconverter.ConvertRichContentToCallToolResult(tt.content)

			// Assert
			require.NotNil(t, result)
			assert.Equal(t, tt.expectedContent, result.Content)
		})
	}
}
