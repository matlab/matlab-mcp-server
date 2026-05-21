// Copyright 2026 The MathWorks, Inc.

package addonmanager

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type InstallationSteps interface {
	UploadMLTBX(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) (func(), error)
	VerifyMLTBXInstallationFile(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) error
	InstallMLTBX(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) error
}

type AddonManager struct {
	installationSteps InstallationSteps
}

func New(
	installationSteps InstallationSteps,
) *AddonManager {
	return &AddonManager{
		installationSteps: installationSteps,
	}
}

func (a *AddonManager) Install(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) error {
	logger.Debug("Installing MATLAB Add-On")

	cleanup, err := a.installationSteps.UploadMLTBX(ctx, logger, client)
	if err != nil {
		return err
	}
	defer cleanup()

	err = a.installationSteps.VerifyMLTBXInstallationFile(ctx, logger, client)
	if err != nil {
		return err
	}

	// installToolbox sometimes throws strange errors, which go away on retry
	var lastErr error
	_, retryErr := retry.Retry(ctx, func() (struct{}, bool, error) {
		lastErr = a.installationSteps.InstallMLTBX(ctx, logger, client)
		if lastErr != nil {
			return struct{}{}, false, nil
		}

		return struct{}{}, true, nil
	}, retry.NewFixedCountRetryStrategy(2))

	if retryErr != nil {
		if lastErr != nil {
			return lastErr
		}
		logger.
			Debug("Failed to install MLTBX in MATLAB")
		return retryErr
	}

	logger.Info("MATLAB Add-On installation complete")

	return nil
}
