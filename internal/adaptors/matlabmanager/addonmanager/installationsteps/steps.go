// Copyright 2026 The MathWorks, Inc.

package installationsteps

import (
	"context"
	_ "embed"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

//go:embed assets/matlab/upload_mltbx.m
var uploadMLTBX string

//go:embed assets/matlab/cleanup_mltbx_installation_file.m
var cleanupMLTBXInstallationFile string

//go:embed assets/matlab/verify_mltbx_installation_file.m
var verifyMLTBXInstallationFile string

//go:embed assets/matlab/install_mltbx.m
var installMLTBX string

//go:embed assets/mltbx/MATLABMCPCoreServerToolbox.mltbx
var matlabAddOn []byte

type InstallationSteps struct{}

func New() *InstallationSteps {
	return &InstallationSteps{}
}

func (s *InstallationSteps) UploadMLTBX(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) (func(), error) {
	cleanupFunctionZeroValue := func() {}

	encoded := base64.StdEncoding.EncodeToString(matlabAddOn)

	writeCode := strings.ReplaceAll(uploadMLTBX,
		"MTBX_BINARY_CONTENT_AS_BASE_ENCODED_STRING",
		encoded,
	)

	logger.Debug("Writing MLTBX to temporary file")

	_, err := client.Eval(ctx, logger, entities.EvalRequest{Code: writeCode})
	if err != nil {
		logger.
			WithError(err).
			Error("Failed to write MLTBX to temporary file")
		return cleanupFunctionZeroValue, err
	}

	cleanupFunction := func() {
		cleanupCtx := context.WithoutCancel(ctx)
		if _, err := client.Eval(cleanupCtx, logger, entities.EvalRequest{
			Code: cleanupMLTBXInstallationFile,
		}); err != nil {
			logger.
				WithError(err).
				Warn("Failed to clean up temporary MLTBX file")
		}
	}

	return cleanupFunction, nil
}

func (s *InstallationSteps) VerifyMLTBXInstallationFile(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) error {
	expectedByteSize := len(matlabAddOn)

	verifyCode := strings.ReplaceAll(verifyMLTBXInstallationFile,
		"EXPECTED_BYTE_SIZE",
		strconv.Itoa(expectedByteSize),
	)

	logger.Debug("Verifying MLTBX temporary file")

	_, err := client.Eval(ctx, logger, entities.EvalRequest{
		Code: verifyCode,
	})
	if err != nil {
		logger.
			WithError(err).
			Error("Failed to verify MLTBX temporary file")
		return err
	}

	return nil
}

func (s *InstallationSteps) InstallMLTBX(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) error {
	logger.Debug("Installing MLTBX")

	_, err := client.Eval(ctx, logger, entities.EvalRequest{
		Code: installMLTBX,
	})
	if err != nil {
		logger.
			WithError(err).
			Debug("Failed to install MLTBX")
		return err
	}

	return nil
}
