// Copyright 2025-2026 The MathWorks, Inc.

package orchestrator

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MessageCatalog interface {
	GetFromError(err messages.Error) string
}

type LifecycleSignaler interface {
	RequestShutdown()
	WaitForShutdownToComplete() error
}

type ApplicationDefinition interface {
	Features() definition.Features
	Dependencies(resources definition.DependenciesProviderResources) (any, error)
	Tools(resources definition.ToolsProviderResources) []tools.Tool
}

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Server interface {
	Run(tools []tools.Tool) error
}

type WatchdogClient interface {
	Start() error
	RegisterProcessPIDWithWatchdog(pid int) error
	Stop() error
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
	NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error)
}

type OSSignaler interface {
	InterruptSignalChan() <-chan os.Signal
}

type DirectoryFactory interface {
	Directory() (directory.Directory, messages.Error)
}

type ResourceLimitManager interface {
	CapOpenFilesLimit(limit uint64) (func() error, error)
}

const UnixOpenFileDescriptorsSoftCap uint64 = 1 << 16

// Orchestrator
type Orchestrator struct {
	messageCatalog        MessageCatalog
	lifecycleSignaler     LifecycleSignaler
	applicationDefinition ApplicationDefinition
	configFactory         ConfigFactory
	server                Server
	watchdogClient        WatchdogClient
	loggerFactory         LoggerFactory
	osSignaler            OSSignaler
	directoryFactory      DirectoryFactory
	resourceLimitManager  ResourceLimitManager
}

func New(
	messageCatalog MessageCatalog,
	lifecycleSignaler LifecycleSignaler,
	applicationDefinition ApplicationDefinition,
	configFactory ConfigFactory,
	server Server,
	watchdogClient WatchdogClient,
	loggerFactory LoggerFactory,
	osSignaler OSSignaler,
	directoryFactory DirectoryFactory,
	resourceLimitManager ResourceLimitManager,
) *Orchestrator {
	orchestrator := &Orchestrator{
		messageCatalog:        messageCatalog,
		lifecycleSignaler:     lifecycleSignaler,
		applicationDefinition: applicationDefinition,
		configFactory:         configFactory,
		server:                server,
		watchdogClient:        watchdogClient,
		loggerFactory:         loggerFactory,
		osSignaler:            osSignaler,
		directoryFactory:      directoryFactory,
		resourceLimitManager:  resourceLimitManager,
	}
	return orchestrator
}

func (o *Orchestrator) StartAndWaitForCompletion(ctx context.Context) error {
	config, messagesErr := o.configFactory.Config()
	if messagesErr != nil {
		return messagesErr
	}

	logger, messagesErr := o.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return messagesErr
	}

	resetOpenFilesLimit, err := o.resourceLimitManager.CapOpenFilesLimit(UnixOpenFileDescriptorsSoftCap)
	if err != nil {
		return err
	}

	defer func() {
		logger.Info("Initiating application shutdown")
		o.lifecycleSignaler.RequestShutdown()

		err := o.lifecycleSignaler.WaitForShutdownToComplete()
		if err != nil {
			logger.WithError(err).Warn("Application shutdown failed")
		}

		logger.Debug("Shutdown functions have all completed, stopping the watchdog")
		err = o.watchdogClient.Stop()
		if err != nil {
			logger.WithError(err).Warn("Watchdog shutdown failed")
		}

		logger.Debug("Restoring open files limit")
		err = resetOpenFilesLimit()
		if err != nil {
			logger.WithError(err).Warn("Failed to restore open files limit")
		}

		logger.Info("Application shutdown complete")
	}()

	logger.
		With("version", config.Version()).
		Info("Initiating application startup")
	config.RecordToLogger(logger)

	directory, messagesErr := o.directoryFactory.Directory()
	if messagesErr != nil {
		return messagesErr
	}
	directory.RecordToLogger(logger)

	err = o.watchdogClient.Start()
	if err != nil {
		return err
	}

	logger.Debug("Building SDK dependencies")
	dependencies, err := o.applicationDefinition.Dependencies(definition.NewDependenciesProviderResources(
		logger,
		config,
		o.messageCatalog,
		o.watchdogClient,
	))
	if err != nil {
		return err
	}

	logger.Debug("Building SDK tools")
	tools := o.applicationDefinition.Tools(definition.NewToolsProviderResources(
		logger,
		config,
		o.messageCatalog,
		dependencies,
		o.loggerFactory,
	))

	serverErrC := make(chan error, 1)
	go func() {
		serverErrC <- o.server.Run(tools)
	}()

	logger.Info("Application startup complete")

	select {
	case <-o.osSignaler.InterruptSignalChan():
		logger.Info("Received termination signal")
		return nil
	case err := <-serverErrC:
		return err
	}
}
