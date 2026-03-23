// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabbuild

import (
	"context"
	"fmt"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/utils/matlabstring"
)

type Args struct {
	WorkingDirectory  string
	Tasks             []string
	ContinueOnFailure bool
	Parallel          bool
	Verbosity         string
	Skip              []string
}

type Result struct {
	ConsoleOutput string
	Success       bool
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

func (u *Usecase) Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, request Args) (Result, error) {
	sessionLogger.Debug("Entering RunMATLABBuild Usecase")
	defer sessionLogger.Debug("Exiting RunMATLABBuild Usecase")

	if request.WorkingDirectory != "" {
		validatedPath, err := u.pathValidator.ValidateFolderPath(request.WorkingDirectory)
		if err != nil {
			sessionLogger.WithError(err).With("path", request.WorkingDirectory).Warn("Path validation failed")
			return Result{}, fmt.Errorf("path validation failed: %w", err)
		}

		_, err = client.Eval(ctx, sessionLogger, entities.EvalRequest{
			Code: fmt.Sprintf("cd('%s')", matlabstring.EscapeSingleQuotes(validatedPath)),
		})
		if err != nil {
			return Result{}, err
		}
	}

	_, _ = client.Eval(ctx, sessionLogger, entities.EvalRequest{Code: "feature('HotLinks',0)"})
	defer func() {
		_, _ = client.Eval(ctx, sessionLogger, entities.EvalRequest{Code: "feature('HotLinks',1)"})
	}()

	response, evalErr := client.Eval(ctx, sessionLogger, entities.EvalRequest{
		Code: buildCommand(request),
	})
	if evalErr != nil {
		return Result{
			ConsoleOutput: evalErr.Error(),
			Success:       false,
		}, nil
	}

	return Result{
		ConsoleOutput: response.ConsoleOutput,
		Success:       true,
	}, nil
}

func buildCommand(request Args) string {
	parts := []string{"buildtool"}
	parts = append(parts, request.Tasks...)

	if request.ContinueOnFailure {
		parts = append(parts, "-continueOnFailure")
	}

	if request.Parallel {
		parts = append(parts, "-parallel")
	}

	if request.Verbosity != "" {
		parts = append(parts, "-verbosity", request.Verbosity)
	}

	for _, task := range request.Skip {
		parts = append(parts, "-skip", task)
	}

	return strings.Join(parts, " ")
}
