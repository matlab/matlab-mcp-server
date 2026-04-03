// Copyright 2026 The MathWorks, Inc.

package installmatlab

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/installmatlab"
)

// Usecase defines the business logic interface for installing MATLAB products.
type Usecase interface {
	Execute(ctx context.Context, logger entities.Logger, args installmatlab.Args) (installmatlab.ReturnArgs, error)
}

// Tool is the MCP tool for installing MATLAB and toolboxes via mpm.
type Tool struct {
	basetool.ToolWithStructuredContentOutput[Args, ReturnArgs]
}

// New creates a new install_matlab tool.
func New(
	loggerFactory basetool.LoggerFactory,
	usecase Usecase,
) *Tool {
	return &Tool{
		ToolWithStructuredContentOutput: basetool.NewToolWithStructuredContent(
			name, title, description,
			annotations.NewDestructiveAnnotations(),
			loggerFactory,
			Handler(usecase),
		),
	}
}

// Handler returns the tool handler function.
func Handler(usecase Usecase) basetool.HandlerWithStructuredContentOutput[Args, ReturnArgs] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (ReturnArgs, error) {
		sessionLogger.Info("Executing install MATLAB tool")
		defer sessionLogger.Info("Done - Executing install MATLAB tool")

		result, err := usecase.Execute(ctx, sessionLogger, installmatlab.Args{
			Release:     inputs.Release,
			Destination: inputs.Destination,
			Products:    inputs.Products,
		})
		if err != nil {
			return ReturnArgs{Output: result.Output}, err
		}

		return ReturnArgs{
			Output: result.Output,
		}, nil
	}
}
