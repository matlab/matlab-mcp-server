// Copyright 2025-2026 The MathWorks, Inc.

package localmatlabsession

import (
	"context"
	"runtime"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

const startupCode = "sessionPath = getenv('MW_MCP_SESSION_DIR');addpath(sessionPath);matlab_mcp.initializeMCP(); clear sessionPath;"

type SessionDirectoryFactory interface {
	New(logger entities.Logger) (directory.Directory, error)
}

type ProcessDetails interface {
	NewAPIKey() string
	EnvironmentVariables(sessionDirPath string, apiKey string, certificateFile string, certificateKey string) []string
	StartupFlag(os string, showMATLAB bool, startupCode string) []string
}

type MATLABProcessLauncher interface {
	Launch(ctx context.Context, logger entities.Logger, sessionRoot string, matlabRoot string, workingDir string, args []string, env []string) (int, func(), <-chan struct{}, error)
}

type Watchdog interface {
	RegisterProcessPIDWithWatchdog(processPID int) error
}

type Starter struct {
	directoryFactory      SessionDirectoryFactory
	processDetails        ProcessDetails
	matlabProcessLauncher MATLABProcessLauncher
	watchdog              Watchdog
}

func NewStarter(
	directoryFactory SessionDirectoryFactory,
	processDetails ProcessDetails,
	matlabProcessLauncher MATLABProcessLauncher,
	watchdog Watchdog,
) *Starter {
	return &Starter{
		directoryFactory:      directoryFactory,
		processDetails:        processDetails,
		matlabProcessLauncher: matlabProcessLauncher,
		watchdog:              watchdog,
	}
}

func (m *Starter) StartLocalMATLABSession(ctx context.Context, logger entities.Logger, request datatypes.LocalSessionDetails) (embeddedconnector.ConnectionDetails, func() error, error) {
	logger.Debug("Starting a local MATLAB session")

	sessionDir, err := m.directoryFactory.New(logger)
	if err != nil {
		return embeddedconnector.ConnectionDetails{}, nil, err
	}

	sessionDirPath := sessionDir.Path()

	logger = logger.With("session_dir", sessionDirPath)
	logger.Debug("Created session directory")

	if !request.IsStartingDirectorySet {
		request.StartingDirectory = sessionDirPath
	}

	uniqueAPIKey := m.processDetails.NewAPIKey()

	env := m.processDetails.EnvironmentVariables(
		sessionDirPath,
		uniqueAPIKey,
		sessionDir.CertificateFile(),
		sessionDir.CertificateKeyFile(),
	)

	startupFlags := m.processDetails.StartupFlag(runtime.GOOS, request.ShowMATLABDesktop, startupCode)

	processID, processCleanup, _, err := m.matlabProcessLauncher.Launch(ctx, logger, sessionDirPath, request.MATLABRoot, request.StartingDirectory, startupFlags, env)
	if err != nil {
		if cleanupErr := sessionDir.Cleanup(); cleanupErr != nil {
			logger.WithError(cleanupErr).Warn("Failed to cleanup session directory after launch error")
		}
		return embeddedconnector.ConnectionDetails{}, nil, err
	}

	logger = logger.With("pid", processID)
	logger.Debug("Started MATLAB process")

	cleanup := func() error {
		if processCleanup != nil {
			processCleanup()
		}
		return sessionDir.Cleanup()
	}

	logger.Debug("Registering process with watchdog")

	if err = m.watchdog.RegisterProcessPIDWithWatchdog(processID); err != nil {
		logger.WithError(err).Warn("Failed to register process with watchdog")
	}

	logger.Debug("Retrieving EC details")

	securePort, certificatePEM, err := sessionDir.GetEmbeddedConnectorDetails()
	if err != nil {
		if cleanupErr := cleanup(); cleanupErr != nil {
			logger.WithError(cleanupErr).Warn("Failed to cleanup after startup error")
		}
		return embeddedconnector.ConnectionDetails{}, nil, err
	}

	logger.Debug("Retrieved EC details")

	return embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           securePort,
		APIKey:         uniqueAPIKey,
		CertificatePEM: certificatePEM,
	}, cleanup, nil
}
