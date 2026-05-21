// Copyright 2026 The MathWorks, Inc.

package custom

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/utils/responseconverter"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/mcpfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalcustomtool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Usecase interface {
	Execute(
		ctx context.Context,
		sessionLogger entities.Logger,
		client entities.MATLABSessionClient,
		request evalcustomtool.Args,
	) (entities.EvalResponse, error)
}

type Tool struct {
	validatedTool definition.ValidatedTool
	handler       mcp.ToolHandlerFor[map[string]any, any]
	toolAdder     basetool.ToolAdder[map[string]any, any]
}

func NewTool(
	validatedTool definition.ValidatedTool,
	loggerFactory basetool.LoggerFactory,
	configFactory ConfigFactory,
	usecase Usecase,
	globalMATLAB entities.GlobalMATLAB,
) *Tool {
	return &Tool{
		validatedTool: validatedTool,
		handler:       Handler(validatedTool, loggerFactory, configFactory, usecase, globalMATLAB),
		toolAdder:     mcpfacade.NewToolAdder[map[string]any, any](),
	}
}

func (t *Tool) Name() string {
	return t.validatedTool.Definition().Name
}

func (t *Tool) AddToServer(server *mcp.Server) error {
	toolDef := t.validatedTool.Definition()
	t.toolAdder.AddTool(
		server,
		&mcp.Tool{
			Name:        toolDef.Name,
			Title:       toolDef.Title,
			Description: toolDef.Description,
			Annotations: toolDef.Annotations,
			InputSchema: toolDef.InputSchema,
		},
		t.handler,
	)
	return nil
}

func Handler(
	validatedTool definition.ValidatedTool,
	loggerFactory basetool.LoggerFactory,
	configFactory ConfigFactory,
	usecase Usecase,
	globalMATLAB entities.GlobalMATLAB,
) mcp.ToolHandlerFor[map[string]any, any] {
	toolDef := validatedTool.Definition()
	toolSig := validatedTool.Signature()

	return func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		logger, messagesErr := loggerFactory.NewMCPSessionLogger(req.Session)
		if messagesErr != nil {
			return nil, nil, messagesErr
		}
		logger = logger.With("tool-name", toolDef.Name)
		logger.Debug("Handling custom tool call request")
		defer logger.Debug("Handled custom tool call request")

		var argumentTypes map[string]string
		if toolDef.InputSchema != nil {
			argumentTypes = make(map[string]string, len(toolDef.InputSchema.Properties))
			for name, prop := range toolDef.InputSchema.Properties {
				argumentTypes[name] = prop.Type
			}
		}

		cfg, cfgErr := configFactory.Config()
		if cfgErr != nil {
			return nil, nil, cfgErr
		}

		client, err := globalMATLAB.Client(ctx, logger)
		if err != nil {
			return nil, nil, err
		}

		response, err := usecase.Execute(ctx, logger, client, evalcustomtool.Args{
			Function:      toolSig.Function,
			Order:         toolSig.Input.Order,
			ArgumentTypes: argumentTypes,
			Arguments:     args,
			CaptureOutput: !cfg.ShouldShowMATLABDesktop(),
		})
		if err != nil {
			return nil, nil, err
		}

		return responseconverter.ConvertRichContentToCallToolResult(
			responseconverter.ConvertEvalResponseToRichContent(response),
		), nil, nil
	}
}
