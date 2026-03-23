// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabbuild

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabbuild"
)

type Usecase interface {
	Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, request runmatlabbuild.Args) (runmatlabbuild.Result, error)
}

type Tool struct {
	basetool.ToolWithStructuredContentOutput[Args, ReturnArgs]
}

func New(
	loggerFactory basetool.LoggerFactory,
	usecase Usecase,
	matlabManager entities.MATLABManager,
) *Tool {
	return &Tool{
		ToolWithStructuredContentOutput: basetool.NewToolWithStructuredAndTextContent(name, title, description, annotations.NewDestructiveAnnotations(), loggerFactory, Handler(usecase, matlabManager), func(r ReturnArgs) string { return r.Log }),
	}
}

func Handler(usecase Usecase, matlabManager entities.MATLABManager) basetool.HandlerWithStructuredContentOutput[Args, ReturnArgs] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (ReturnArgs, error) {
		sessionID := entities.SessionID(inputs.SessionID)

		sessionLogger = sessionLogger.With("session_id", sessionID)

		sessionLogger.Info("Executing Run MATLAB Build in Session tool")
		defer sessionLogger.Info("Done - Executing Run MATLAB Build in Session tool")

		client, err := matlabManager.GetMATLABSessionClient(ctx, sessionLogger, sessionID)
		if err != nil {
			return ReturnArgs{}, err
		}

		result, err := usecase.Execute(ctx, sessionLogger, client, runmatlabbuild.Args{
			WorkingDirectory:  inputs.WorkingDirectory,
			Tasks:             inputs.Tasks,
			ContinueOnFailure: inputs.ContinueOnFailure,
			Parallel:          inputs.Parallel,
			Verbosity:         inputs.Verbosity,
			Skip:              inputs.Skip,
		})
		if err != nil {
			return ReturnArgs{}, err
		}

		return ReturnArgs{
			Log:     result.ConsoleOutput,
			Success: result.Success,
		}, nil
	}
}
