// Copyright 2026 The MathWorks, Inc.

package definition

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Tool struct {
	Name        string               `json:"name"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	InputSchema *jsonschema.Schema   `json:"inputSchema"`
	Annotations *mcp.ToolAnnotations `json:"annotations,omitempty"`
}

type ValidatedTool interface {
	Definition() Tool
	Signature() Signature
}

type Signature struct {
	Function string         `json:"function"`
	Input    SignatureInput `json:"input"`
}

type SignatureInput struct {
	Order []string `json:"order"`
}
