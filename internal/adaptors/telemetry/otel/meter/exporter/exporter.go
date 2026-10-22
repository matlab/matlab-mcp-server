// Copyright 2026 The MathWorks, Inc.

package exporter

import (
	"context"
	"net/url"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

const defaultMetricsSignalPath = "/v1/metrics"

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type OSLayer interface {
	LookupEnv(key string) (string, bool)
}

type Factory struct {
	loggerFactory LoggerFactory
	configFactory ConfigFactory
	osLayer       OSLayer
}

func NewFactory(
	loggerFactory LoggerFactory,
	configFactory ConfigFactory,
	osLayer OSLayer,
) *Factory {
	return &Factory{
		loggerFactory: loggerFactory,
		configFactory: configFactory,
		osLayer:       osLayer,
	}
}

func (f *Factory) New() (otel.MetricExporter, messages.Error) {
	logger, messagesErr := f.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return nil, messagesErr
	}

	logger.Debug("Creating OTLP HTTP metric exporter")
	defer logger.Debug("Done creating OTLP HTTP metric exporter")

	cfg, messagesErr := f.configFactory.Config()
	if messagesErr != nil {
		return nil, messagesErr
	}

	options := []otlpmetrichttp.Option{}
	_, env1Exists := f.osLayer.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	_, env2Exists := f.osLayer.LookupEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT")

	if !env1Exists && !env2Exists {
		endpoint := cfg.TelemetryCollectorEndpoint()
		logger.Debug("Using configured OTLP HTTP endpoint")

		parsedEndpoint, parseErr := url.Parse(endpoint)
		if parseErr != nil {
			logger.Error("Failed to create OTLP HTTP metric exporter")
			return nil, messages.New_StartupErrors_TelemetryInitializationFailed_Error()
		}

		options = append(options, otlpmetrichttp.WithEndpointURL(endpoint))
		if parsedEndpoint.Path == "" || parsedEndpoint.Path == "/" {
			options = append(options, otlpmetrichttp.WithURLPath(defaultMetricsSignalPath))
		}
	}

	if cfg.TelemetryCollectorEndpointInsecure() {
		logger.Warn("Using insecure connection for OTLP HTTP endpoint")
		options = append(options, otlpmetrichttp.WithInsecure())
	}

	exporter, err := otlpmetrichttp.New(
		context.Background(),
		options...,
	)
	if err != nil {
		logger.Error("Failed to create OTLP HTTP metric exporter")
		return nil, messages.New_StartupErrors_TelemetryInitializationFailed_Error()
	}

	return exporter, nil
}
