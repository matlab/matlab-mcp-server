// Copyright 2025-2026 The MathWorks, Inc.

package evalmatlabcode

import (
	"context"
	"fmt"

	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type Args struct {
	Code          string
	ProjectPath   string
	CaptureOutput bool
}

type PathValidator interface {
	ValidateFolderPath(filePath string) (string, error)
}

type Usecase struct {
	pathValidator PathValidator
}

func New(
	pathValidator PathValidator,
) *Usecase {
	return &Usecase{
		pathValidator: pathValidator,
	}
}

func (u *Usecase) Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, request Args) (entities.EvalResponse, error) {
	sessionLogger.Debug("Entering EvalInlMATLAB Usecase")
	defer sessionLogger.Debug("Exiting EvalInMATLAB Usecase")

	if request.ProjectPath != "" {
		validatedPath, err := u.pathValidator.ValidateFolderPath(request.ProjectPath)
		if err != nil {
			sessionLogger.WithError(err).With("path", request.ProjectPath).Warn("Path validation failed")
			return entities.EvalResponse{}, fmt.Errorf("path validation failed: %w", err)
		}

		cdRequest := entities.FEvalRequest{
			Function:   "cd",
			Arguments:  []string{validatedPath},
			NumOutputs: 0,
		}
		_, err = client.FEval(ctx, sessionLogger, cdRequest)
		if err != nil {
			return entities.EvalResponse{}, err
		}
	}

	evalRequest := entities.EvalRequest{
		Code: request.Code,
	}

	if request.CaptureOutput {
		return client.EvalWithCapture(ctx, sessionLogger, evalRequest)
	}
	return client.Eval(ctx, sessionLogger, evalRequest)
}
