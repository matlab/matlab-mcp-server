// Copyright 2025-2026 The MathWorks, Inc.

package responseconverter

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func ConvertEvalResponseToRichContent(response entities.EvalResponse) tools.RichContent {
	imageData := make([]tools.PNGImageData, len(response.Images))
	for i := range response.Images {
		imageData[i] = tools.PNGImageData(response.Images[i])
	}
	return tools.RichContent{
		TextContent:  []string{response.ConsoleOutput},
		ImageContent: imageData,
	}
}

func ConvertRichContentToCallToolResult(content tools.RichContent) *mcp.CallToolResult {
	result := &mcp.CallToolResult{
		Content: []mcp.Content{},
	}
	for _, text := range content.TextContent {
		result.Content = append(result.Content, &mcp.TextContent{Text: text})
	}
	for _, imageData := range content.ImageContent {
		result.Content = append(result.Content, &mcp.ImageContent{
			MIMEType: "image/png",
			Data:     imageData,
		})
	}
	return result
}
