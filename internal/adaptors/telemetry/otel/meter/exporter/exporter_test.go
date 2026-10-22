// Copyright 2026 The MathWorks, Inc.

package exporter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel/meter/exporter"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	exportermocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry/otel/meter/exporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	// Act
	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	expectedEndpoint := "localhost:4317"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT").
		Return("", false).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT").
		Return("", false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return(expectedEndpoint).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpointInsecure().
		Return(false).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	result, err := factory.New()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestFactory_New_DoesNotLogInvalidEndpoint(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	const endpoint = "https://sensitive-collector.invalid/v1/%zz"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT").
		Return("", false).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT").
		Return("", false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return(endpoint).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	_, _ = factory.New()

	// Assert
	for _, logs := range []testutils.InspectableLogs{
		testLogger.DebugLogs(),
		testLogger.InfoLogs(),
		testLogger.WarnLogs(),
		testLogger.ErrorLogs(),
	} {
		for message, fields := range logs {
			assert.NotContains(t, message, endpoint)
			for _, field := range fields {
				assert.NotContains(t, fmt.Sprint(field), endpoint)
			}
		}
	}
}

func TestFactory_New_LoggerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	result, err := factory.New()

	// Assert
	require.Nil(t, result)
	require.ErrorIs(t, err, expectedError)
}

func TestFactory_New_ConfigError(t *testing.T) {
	// Arrange
	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	result, err := factory.New()

	// Assert
	require.Nil(t, result, "Exporter should be nil when config fails")
	require.ErrorIs(t, err, expectedError)
}

func TestFactory_New_InsecureEndpoint(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	expectedEndpoint := "localhost:4317"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT").
		Return("", false).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT").
		Return("", false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return(expectedEndpoint).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpointInsecure().
		Return(true).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	result, err := factory.New()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestFactory_New_EndpointWithoutPathExportsToDefaultMetricsPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedRequestPath := "/v1/metrics"
	requestPaths := make(chan string, 1)
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPaths <- r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer collector.Close()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT").
		Return("", false).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT").
		Return("", false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return(collector.URL).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpointInsecure().
		Return(false).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	result, messagesErr := factory.New()
	require.NoError(t, messagesErr)
	exportErr := result.Export(t.Context(), &metricdata.ResourceMetrics{})

	// Assert
	require.NoError(t, exportErr)
	require.Equal(t, expectedRequestPath, <-requestPaths)
}

func TestFactory_New_EndpointWithRootPathExportsToDefaultMetricsPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedRequestPath := "/v1/metrics"
	requestPaths := make(chan string, 1)
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPaths <- r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer collector.Close()

	expectedEndpoint := collector.URL + "/"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT").
		Return("", false).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT").
		Return("", false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return(expectedEndpoint).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpointInsecure().
		Return(false).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	result, messagesErr := factory.New()
	require.NoError(t, messagesErr)
	exportErr := result.Export(t.Context(), &metricdata.ResourceMetrics{})

	// Assert
	require.NoError(t, exportErr)
	require.Equal(t, expectedRequestPath, <-requestPaths)
}

func TestFactory_New_EndpointWithPathExportsToThatPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedRequestPath := "/otlp/v1/metrics"
	requestPaths := make(chan string, 1)
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPaths <- r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer collector.Close()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT").
		Return("", false).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT").
		Return("", false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return(collector.URL + expectedRequestPath).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpointInsecure().
		Return(false).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	result, messagesErr := factory.New()
	require.NoError(t, messagesErr)
	exportErr := result.Export(t.Context(), &metricdata.ResourceMetrics{})

	// Assert
	require.NoError(t, exportErr)
	require.Equal(t, expectedRequestPath, <-requestPaths)
}

func TestFactory_New_OTELEndpointEnvVarSet(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT").
		Return("http://env-endpoint:4318", true).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT").
		Return("", false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpointInsecure().
		Return(false).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	result, err := factory.New()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestFactory_New_OTELMetricsEndpointEnvVarSet(t *testing.T) {
	// Arrange
	mockLoggerFactory := &exportermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &exportermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockOSLayer := &exportermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT").
		Return("", false).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT").
		Return("http://metrics-endpoint:4318", true).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpointInsecure().
		Return(false).
		Once()

	factory := exporter.NewFactory(mockLoggerFactory, mockConfigFactory, mockOSLayer)

	// Act
	result, err := factory.New()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
}
