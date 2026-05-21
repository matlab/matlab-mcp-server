// Copyright 2025-2026 The MathWorks, Inc.

package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type PNGImageData []byte

// RichContent is used as a tool output, when unstructured content should be used.
// That is, the tool will have no output schema and `structuredContent` will be `nil`.
// This should only be used when the tool needs to return content like images, sound, or resources.
type RichContent struct {
	TextContent  []string
	ImageContent []PNGImageData
}

type Tool interface {
	Name() string
	AddToServer(server *mcp.Server) error
}

type ToolWithUnstructuredContentOutput[ToolInput any] interface {
	Tool
	Handler() mcp.ToolHandlerFor[ToolInput, any]
}

type ToolWithStructuredContentOutput[ToolInput any, ToolOutput any] interface {
	Tool
	Handler() mcp.ToolHandlerFor[ToolInput, ToolOutput]
}
