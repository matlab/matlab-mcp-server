// Copyright 2025-2026 The MathWorks, Inc.

package basetool

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/mcpfacade"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ToolWithStructuredContentOutput[ToolInput, ToolOutput any] struct {
	tool[ToolInput, ToolOutput]
	structuredContentHandler HandlerWithStructuredContentOutput[ToolInput, ToolOutput]
	textContentExtractor     func(ToolOutput) string
}
type HandlerWithStructuredContentOutput[ToolInput, ToolOutput any] func(context.Context, entities.Logger, ToolInput) (ToolOutput, error)

func NewToolWithStructuredAndTextContent[ToolInput, ToolOutput any](
	name string,
	title string,
	description string,
	annotations AnnotationProvider,
	loggerFactory LoggerFactory,
	handler func(context.Context, entities.Logger, ToolInput) (ToolOutput, error),
	textContentExtractor func(ToolOutput) string,
) ToolWithStructuredContentOutput[ToolInput, ToolOutput] {
	t := NewToolWithStructuredContent(name, title, description, annotations, loggerFactory, handler)
	t.textContentExtractor = textContentExtractor
	return t
}

func NewToolWithStructuredContent[ToolInput, ToolOutput any](
	name string,
	title string,
	description string,
	annotations AnnotationProvider,
	loggerFactory LoggerFactory,
	handler func(context.Context, entities.Logger, ToolInput) (ToolOutput, error),
) ToolWithStructuredContentOutput[ToolInput, ToolOutput] {
	return ToolWithStructuredContentOutput[ToolInput, ToolOutput]{
		tool: tool[ToolInput, ToolOutput]{
			name:          name,
			title:         title,
			description:   description,
			annotations:   annotations,
			loggerFactory: loggerFactory,
			// Manually inject adder as only have type information at compile time
			toolAdder: mcpfacade.NewToolAdder[ToolInput, ToolOutput](),
		},
		structuredContentHandler: handler,
	}
}

func (t ToolWithStructuredContentOutput[_, _]) AddToServer(server *mcp.Server) error {
	if t.annotations == nil {
		return fmt.Errorf(UnexpectedErrorPrefixForLLM + "annotations must not be nil")
	}

	inputSchema, err := t.GetInputSchema()
	if err != nil {
		return err
	}

	outputSchema, err := t.GetOutputSchema()
	if err != nil {
		return err
	}

	t.toolAdder.AddTool(
		server,
		&mcp.Tool{
			Name:         t.name,
			Title:        t.title,
			Description:  t.description,
			Annotations:  t.annotations.ToToolAnnotations(),
			InputSchema:  inputSchema,
			OutputSchema: outputSchema,
		},
		t.Handler(),
	)

	return nil
}

func (t ToolWithStructuredContentOutput[ToolInput, ToolOutput]) Handler() mcp.ToolHandlerFor[ToolInput, ToolOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ToolInput) (*mcp.CallToolResult, ToolOutput, error) {
		var toolOutputZeroValue ToolOutput

		logger, messagesErr := t.loggerFactory.NewMCPSessionLogger(req.Session)
		if messagesErr != nil {
			return nil, toolOutputZeroValue, messagesErr
		}

		logger = logger.With("tool-name", t.name)
		logger.Debug("Handling tool call request")
		defer logger.Debug("Handled tool call request")

		if t.structuredContentHandler == nil {
			err := fmt.Errorf(UnexpectedErrorPrefixForLLM + "no structured handler available")
			logger.WithError(err).Warn("Structured content handler is nil")
			return nil, toolOutputZeroValue, err
		}

		toolOutput, err := t.structuredContentHandler(ctx, logger, input)
		if err != nil {
			logger.WithError(err).Warn("Structured handler returned an error")
			return nil, toolOutputZeroValue, err
		}
		if t.textContentExtractor != nil {
			text := t.textContentExtractor(toolOutput)
			result := &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}
			return result, toolOutput, nil
		}
		return nil, toolOutput, nil
	}
}

func (_ ToolWithStructuredContentOutput[_, ToolOutput]) GetOutputSchema() (any, error) {
	return jsonschema.For[ToolOutput](&jsonschema.ForOptions{})
}
