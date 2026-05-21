// Copyright 2025-2026 The MathWorks, Inc.

package basetool

import (
	"context"
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/utils/responseconverter"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/mcpfacade"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ToolWithUnstructuredContentOutput[ToolInput any] struct {
	tool[ToolInput, any]
	unstructuredContentHandler HandlerWithUnstructuredContentOutput[ToolInput]
}

type HandlerWithUnstructuredContentOutput[ToolInput any] func(context.Context, entities.Logger, ToolInput) (tools.RichContent, error)

func NewToolWithUnstructuredContent[ToolInput any](
	name string,
	title string,
	description string,
	annotations AnnotationProvider,
	loggerFactory LoggerFactory,
	handler func(context.Context, entities.Logger, ToolInput) (tools.RichContent, error),
) ToolWithUnstructuredContentOutput[ToolInput] {
	return ToolWithUnstructuredContentOutput[ToolInput]{
		tool: tool[ToolInput, any]{
			name:          name,
			title:         title,
			description:   description,
			annotations:   annotations,
			loggerFactory: loggerFactory,
			// Manually inject adder as only have type information at compile time
			toolAdder: mcpfacade.NewToolAdder[ToolInput, any](),
		},
		unstructuredContentHandler: handler,
	}
}

func (t ToolWithUnstructuredContentOutput[_]) AddToServer(server *mcp.Server) error {
	if t.annotations == nil {
		return fmt.Errorf(UnexpectedErrorPrefixForLLM + "annotations must not be nil")
	}

	inputSchema, err := t.GetInputSchema()
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
			OutputSchema: nil,
		},
		t.Handler(),
	)

	return nil
}

func (t ToolWithUnstructuredContentOutput[ToolInput]) Handler() mcp.ToolHandlerFor[ToolInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ToolInput) (*mcp.CallToolResult, any, error) {
		logger, messagesErr := t.loggerFactory.NewMCPSessionLogger(req.Session)
		if messagesErr != nil {
			return nil, nil, messagesErr
		}

		logger = logger.With("tool-name", t.name)
		logger.Debug("Handling tool call request")
		defer logger.Debug("Handled tool call request")

		if t.unstructuredContentHandler == nil {
			err := fmt.Errorf(UnexpectedErrorPrefixForLLM + "no unstructured handler available")
			logger.WithError(err).Warn("Unstructured content handler is nil")
			return nil, nil, err
		}

		richContent, err := t.unstructuredContentHandler(ctx, logger, input)
		if err != nil {
			logger.WithError(err).Warn("Unstructured handler returned an error")
			return nil, nil, err
		}
		return responseconverter.ConvertRichContentToCallToolResult(richContent), nil, nil
	}
}
