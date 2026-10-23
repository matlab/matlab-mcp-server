// Copyright 2025-2026 The MathWorks, Inc.

package process

import (
	"strconv"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
)

type Config interface {
	LogLevel() entities.LogLevel
	DisableTelemetry() bool
}

type Directory interface {
	BaseDir() string
	ID() string
}

func newProcess(
	osLayer OSLayer,
	logger entities.Logger,
	directory Directory,
	config Config,
) messages.Error {
	programPath, err := osLayer.Executable()
	if err != nil {
		logger.WithError(err).Error("Failed to get executable path")
		return messages.New_StartupErrors_FailedToGetExecutablePath_Error()
	}

	cmd := osLayer.Command(programPath,
		"--"+defaultparameters.WatchdogMode().GetFlagName(),
		"--"+defaultparameters.BaseDir().GetFlagName(), directory.BaseDir(),
		"--"+defaultparameters.ServerInstanceID().GetFlagName(), directory.ID(),
		"--"+defaultparameters.LogLevel().GetFlagName(), string(config.LogLevel()),
		"--"+defaultparameters.DisableTelemetry().GetFlagName()+"="+strconv.FormatBool(config.DisableTelemetry()),
	)

	cmd.SetSysProcAttr(getSysProcAttrForDetachingAProcess())

	if err := cmd.Start(); err != nil {
		logger.WithError(err).Error("Failed to start watchdog process")
		return messages.New_StartupErrors_FailedToStartWatchdogProcess_Error()
	}

	return nil
}
