// Copyright 2026 The MathWorks, Inc.

package setupmatlab

import (
	"context"
	"fmt"
	"io"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type OSLayer interface {
	Stdout() io.Writer
}

type MessageCatalog interface {
	Get(message messages.MessageKey) string
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type WatchdogClient interface {
	Start() error
	Stop() error
}

type GlobalMATLAB interface {
	Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error)
}

type AddonManager interface {
	Install(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) messages.Error
}

type Mode struct {
	osLayer        OSLayer
	messageCatalog MessageCatalog
	loggerFactory  LoggerFactory
	watchdogClient WatchdogClient
	globalMATLAB   GlobalMATLAB
	addonManager   AddonManager
}

func New(
	osLayer OSLayer,
	messageCatalog MessageCatalog,
	loggerFactory LoggerFactory,
	watchdogClient WatchdogClient,
	globalMATLAB GlobalMATLAB,
	addonManager AddonManager,
) *Mode {
	return &Mode{
		osLayer:        osLayer,
		messageCatalog: messageCatalog,
		loggerFactory:  loggerFactory,
		watchdogClient: watchdogClient,
		globalMATLAB:   globalMATLAB,
		addonManager:   addonManager,
	}
}

func (m *Mode) StartAndWaitForCompletion(ctx context.Context) messages.Error {
	logger, messagesErr := m.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return messagesErr
	}

	logger.Debug("Starting watchdog")

	err := m.watchdogClient.Start()
	if err != nil {
		logger.
			WithError(err).
			Error("Failed to start watchdog")
		return messages.New_AddonManagerErrors_InstallFailed_Error()
	}
	defer func() {
		logger.Debug("Stopping watchdog")

		err := m.watchdogClient.Stop()
		if err != nil {
			logger.
				WithError(err).
				Warn("Watchdog shutdown failed")
		}
	}()

	logger.Info("Installing MATLAB Add-On")

	client, err := m.globalMATLAB.Client(ctx, logger)
	if err != nil {
		logger.
			WithError(err).
			Error("Failed to get MATLAB Client")
		return messages.New_AddonManagerErrors_InstallFailed_Error()
	}

	messagesErr = m.addonManager.Install(ctx, logger, client)
	if messagesErr != nil {
		return messagesErr
	}

	successMessage := m.messageCatalog.Get(messages.CLIMessages_SuccessfullySetupMATLAB)
	_, err = fmt.Fprintf(m.osLayer.Stdout(), "%s\n", successMessage)
	if err != nil {
		logger.
			WithError(err).
			Error("Failed to write success message")
		// Let's not fail the command for this
		return nil
	}

	return nil
}
