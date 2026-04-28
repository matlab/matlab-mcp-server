// Copyright 2025-2026 The MathWorks, Inc.

package modeselector

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Parser interface {
	Usage() (string, messages.Error)
}

type TelemetryFactory interface {
	Telemetry() (telemetry.Telemetry, messages.Error)
}

type WatchdogProcess interface { //nolint:iface // Intentional interface for deps injection
	StartAndWaitForCompletion(ctx context.Context) error
}

type Orchestrator interface { //nolint:iface // Intentional interface for deps injection
	StartAndWaitForCompletion(ctx context.Context) error
}

type OSLayer interface {
	Stdout() io.Writer
}

type LifecycleSignaler interface {
	RequestShutdown()
	WaitForShutdownToComplete() error
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type SetupMATLAB interface {
	StartAndWaitForCompletion(ctx context.Context) messages.Error
}

type ModeSelector struct {
	configFactory     ConfigFactory
	telemetryFactory  TelemetryFactory
	watchdogProcess   WatchdogProcess
	orchestrator      Orchestrator
	osLayer           OSLayer
	parser            Parser
	lifecycleSignaler LifecycleSignaler
	loggerFactory     LoggerFactory
	setupMATLAB       SetupMATLAB
}

func New(
	configFactory ConfigFactory,
	parser Parser,
	telemetryFactory TelemetryFactory,
	watchdogProcess WatchdogProcess,
	orchestrator Orchestrator,
	osLayer OSLayer,
	lifecycleSignaler LifecycleSignaler,
	loggerFactory LoggerFactory,
	setupMATLAB SetupMATLAB,
) *ModeSelector {
	return &ModeSelector{
		configFactory:     configFactory,
		parser:            parser,
		telemetryFactory:  telemetryFactory,
		watchdogProcess:   watchdogProcess,
		orchestrator:      orchestrator,
		osLayer:           osLayer,
		lifecycleSignaler: lifecycleSignaler,
		loggerFactory:     loggerFactory,
		setupMATLAB:       setupMATLAB,
	}
}

func (m *ModeSelector) StartAndWaitForCompletion(ctx context.Context) messages.Error {
	config, err := m.configFactory.Config()
	if err != nil {
		return err
	}

	logger, err := m.loggerFactory.GetGlobalLogger()
	if err != nil {
		return err
	}

	telemetryInstance, err := m.telemetryFactory.Telemetry()
	if err != nil {
		return err
	}

	telemetryInstance.RecordServerStart(ctx)

	switch {
	case config.HelpMode():
		usage, messagesErr := m.parser.Usage()
		if messagesErr != nil {
			return m.shutdownAndReturn(logger, messagesErr)
		}
		_, err := fmt.Fprintf(m.osLayer.Stdout(), "%s\n", usage)
		if err != nil {
			messagesErr := messages.New_StartupErrors_WriteError_Error("help", err.Error())
			return m.shutdownAndReturn(logger, messagesErr)
		}
		return m.shutdownAndReturn(logger, nil)
	case config.VersionMode():
		_, err := fmt.Fprintf(m.osLayer.Stdout(), "%s\n", config.Version())
		if err != nil {
			messagesErr := messages.New_StartupErrors_WriteError_Error("version", err.Error())
			return m.shutdownAndReturn(logger, messagesErr)
		}
		return m.shutdownAndReturn(logger, nil)
	case config.WatchdogMode():
		return m.toMessagesError(logger, m.watchdogProcess.StartAndWaitForCompletion(ctx))
	case config.SetupMATLABMode():
		err := m.setupMATLAB.StartAndWaitForCompletion(ctx)
		return m.shutdownAndReturn(logger, err)
	default:
		return m.toMessagesError(logger, m.orchestrator.StartAndWaitForCompletion(ctx))
	}
}

func (m *ModeSelector) toMessagesError(logger entities.Logger, err error) messages.Error {
	if err == nil {
		return nil
	}
	var messagesErr messages.Error
	if errors.As(err, &messagesErr) {
		return messagesErr
	}
	logger.WithError(err).Error("Server failed with unexpected error")
	return messages.New_StartupErrors_GenericInitializeFailure_Error()
}

func (m *ModeSelector) shutdownAndReturn(logger entities.Logger, err messages.Error) messages.Error {
	m.lifecycleSignaler.RequestShutdown()
	if shutdownErr := m.lifecycleSignaler.WaitForShutdownToComplete(); shutdownErr != nil {
		logger.WithError(shutdownErr).Warn("Shutdown failed")
	}
	return err
}
