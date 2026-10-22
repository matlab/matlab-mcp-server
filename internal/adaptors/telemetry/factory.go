// Copyright 2026 The MathWorks, Inc.

package telemetry

import (
	"context"
	"sync"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel/instruments"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

const name = "github.com/matlab/matlab-mcp-server"

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type ExporterFactory interface {
	New() (otel.MetricExporter, messages.Error)
}

type MeterProviderFactory interface {
	New(
		exporter otel.MetricExporter,
		errorHandler otel.ErrorHandler,
		serviceName, serviceVersion string,
	) (otel.MeterProvider, messages.Error)
}

type InstrumentFactory interface {
	NewInt64Counter(meter metric.Meter, name, description, unit string) (instruments.Int64Counter, error)
}

type DirectoryFactory interface {
	Directory() (directory.Directory, messages.Error)
}

type OSLayer interface {
	GOOS() string
}

type OSVersionProvider interface {
	Version() (string, error)
}

type Definition interface {
	Name() string
}

type SessionCorrelationIDProvider interface {
	CurrentCorrelationID(ctx context.Context) string
}

type ClientConnectionInfo struct {
	Name             string
	Title            string
	WebsiteURL       string
	Version          string
	Capabilities     []string
	CapabilitiesJSON string
}

type ToolSource int

const (
	ToolSourceBuiltin ToolSource = iota
	ToolSourceExtension
)

func (s ToolSource) String() string {
	switch s {
	case ToolSourceBuiltin:
		return "builtin"
	case ToolSourceExtension:
		return "extension"
	default:
		return "unknown"
	}
}

type Telemetry interface {
	RecordServerStart(ctx context.Context)
	RecordClientConnection(ctx context.Context, info ClientConnectionInfo)
	RecordToolCallRequest(ctx context.Context, toolName string, source ToolSource)
}

type Factory struct {
	loggerFactory                LoggerFactory
	configFactory                ConfigFactory
	exporterFactory              ExporterFactory
	meterProviderFactory         MeterProviderFactory
	instrumentFactory            InstrumentFactory
	directoryFactory             DirectoryFactory
	osLayer                      OSLayer
	osVersionProvider            OSVersionProvider
	serverDefinition             Definition
	sessionCorrelationIDProvider SessionCorrelationIDProvider

	telemetryOnce  sync.Once
	telemetryError messages.Error
	telemetry      Telemetry
}

func NewFactory(
	loggerFactory LoggerFactory,
	configFactory ConfigFactory,
	exporterFactory ExporterFactory,
	meterProviderFactory MeterProviderFactory,
	instrumentFactory InstrumentFactory,
	directoryFactory DirectoryFactory,
	osLayer OSLayer,
	osVersionProvider OSVersionProvider,
	definition Definition,
	sessionCorrelationIDProvider SessionCorrelationIDProvider,
) *Factory {
	return &Factory{
		loggerFactory:                loggerFactory,
		configFactory:                configFactory,
		exporterFactory:              exporterFactory,
		meterProviderFactory:         meterProviderFactory,
		instrumentFactory:            instrumentFactory,
		directoryFactory:             directoryFactory,
		osLayer:                      osLayer,
		osVersionProvider:            osVersionProvider,
		serverDefinition:             definition,
		sessionCorrelationIDProvider: sessionCorrelationIDProvider,
	}
}

func (f *Factory) Telemetry() (Telemetry, messages.Error) {
	f.telemetryOnce.Do(func() {
		telemetry, err := f.newOTELTelemetry()
		if err != nil {
			f.telemetryError = err
			return
		}

		f.telemetry = telemetry
	})

	if f.telemetryError != nil {
		return nil, f.telemetryError
	}

	return f.telemetry, nil
}

func (f *Factory) newOTELTelemetry() (Telemetry, messages.Error) {
	logger, err := f.loggerFactory.GetGlobalLogger()
	if err != nil {
		return nil, err
	}

	logger.Debug("Initializing telemetry")
	defer logger.Debug("Done initializing telemetry")

	cfg, err := f.configFactory.Config()
	if err != nil {
		return nil, err
	}

	directory, err := f.directoryFactory.Directory()
	if err != nil {
		return nil, err
	}

	endpoint := cfg.TelemetryCollectorEndpoint()
	errorHandler := otel.NewErrorHandler(logger, endpoint)
	otel.SetErrorHandler(errorHandler)

	var meterProvider metric.MeterProvider

	switch {
	case cfg.DisableTelemetry():
		logger.Info("Telemetry disabled by configuration")
		meterProvider = noop.NewMeterProvider()
	case endpoint == "":
		logger.Info("Telemetry off: no collector endpoint set")
		meterProvider = noop.NewMeterProvider()
	default:
		logger.Debug("Creating OTEL metric exporter")
		exporter, err := f.exporterFactory.New()
		if err != nil {
			return nil, err
		}

		logger.Debug("Creating OTEL meter provider")
		concreteMeterProvider, err := f.meterProviderFactory.New(exporter, errorHandler, f.serverDefinition.Name(), cfg.Version())
		if err != nil {
			return nil, err
		}

		meterProvider = concreteMeterProvider
	}

	logger.Debug("Creating OTEL telemetry instance")
	return newOTELTelemetry(
		logger,
		meterProvider.Meter(name),
		f.instrumentFactory,
		cfg,
		directory,
		f.osLayer,
		f.osVersionProvider,
		f.serverDefinition,
		f.sessionCorrelationIDProvider,
	)
}
