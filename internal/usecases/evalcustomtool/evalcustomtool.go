// Copyright 2026 The MathWorks, Inc.

package evalcustomtool

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalcustomtool/functioncall"
)

type FunctionCallAssembler interface {
	Assemble(args functioncall.Args) (string, error)
}

type Args struct {
	Function      string
	Order         []string
	ArgumentTypes map[string]string
	Arguments     map[string]any
	CaptureOutput bool
}

type Usecase struct {
	functionCallAssembler FunctionCallAssembler
}

func New(functionCallAssembler FunctionCallAssembler) *Usecase {
	return &Usecase{
		functionCallAssembler: functionCallAssembler,
	}
}

func (u *Usecase) Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, request Args) (entities.EvalResponse, error) {
	sessionLogger.Debug("Entering EvalCustomTool Usecase")
	defer sessionLogger.Debug("Exiting EvalCustomTool Usecase")

	code, err := u.functionCallAssembler.Assemble(functioncall.Args{
		Function:      request.Function,
		Order:         request.Order,
		ArgumentTypes: request.ArgumentTypes,
		Arguments:     request.Arguments,
	})
	if err != nil {
		return entities.EvalResponse{}, err
	}

	evalRequest := entities.EvalRequest{
		Code: code,
	}

	if request.CaptureOutput {
		return client.EvalWithCapture(ctx, sessionLogger, evalRequest)
	}
	return client.Eval(ctx, sessionLogger, evalRequest)
}
