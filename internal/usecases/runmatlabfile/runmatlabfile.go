// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabfile

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/usecases/utils/pathextractor"
)

// matlabIdentifierPattern matches a valid MATLAB file name: it must start with a
// letter and contain only letters, digits, and underscores.
var matlabIdentifierPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

type Args struct {
	ScriptPath    string
	CaptureOutput bool
}

type PathValidator interface {
	ValidateMATLABCodeFile(filePath string) (string, error)
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
	sessionLogger.Debug("Entering RunMATLABFile Usecase")
	defer sessionLogger.Debug("Exiting RunMATLABFile Usecase")

	validatedPath, err := u.pathValidator.ValidateMATLABCodeFile(request.ScriptPath)
	if err != nil {
		return entities.EvalResponse{}, err
	}

	scriptDir, scriptName := pathextractor.ExtractPathComponents(validatedPath)

	if !matlabIdentifierPattern.MatchString(scriptName) {
		return entities.EvalResponse{}, fmt.Errorf("invalid MATLAB file name %q: a MATLAB script file name must start with a letter and contain only letters, digits, and underscores", filepath.Base(validatedPath))
	}

	_, err = client.FEval(ctx, sessionLogger, entities.FEvalRequest{
		Function:   "cd",
		Arguments:  []string{scriptDir},
		NumOutputs: 0,
	})
	if err != nil {
		return entities.EvalResponse{}, err
	}

	runCodeRequest := entities.EvalRequest{
		Code: scriptName,
	}

	if request.CaptureOutput {
		return client.EvalWithCapture(ctx, sessionLogger, runCodeRequest)
	}
	return client.Eval(ctx, sessionLogger, runCodeRequest)
}
