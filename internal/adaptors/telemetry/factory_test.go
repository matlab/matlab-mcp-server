// Copyright 2026 The MathWorks, Inc.

package telemetry_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	directorymocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/directory"
	telemetrymocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry"
	otelmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry/otel"
	instrumentsmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry/otel/instruments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	// Act
	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		nil,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Assert
	assert.NotNil(t, factory)
}

func TestFactory_Telemetry_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &telemetrymocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	mockInt64Counter := &instrumentsmocks.MockInt64Counter{}
	defer mockInt64Counter.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	noopMeterProvider := noop.NewMeterProvider()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockConfig.EXPECT().
		DisableTelemetry().
		Return(false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return("localhost:4317").
		Once()

	mockExporterFactory.EXPECT().
		New().
		Return(mockExporter, nil).
		Once()

	mockServerDefinition.EXPECT().
		Name().
		Return("matlab-mcp-server").
		Once()

	mockConfig.EXPECT().
		Version().
		Return("1.0.0").
		Once()

	mockMeterProviderFactory.EXPECT().
		New(mockExporter, mock.Anything, "matlab-mcp-server", "1.0.0").
		Return(noopMeterProvider, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockInt64Counter, nil).
		Times(3)

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		mockDirectoryFactory,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	result, err := factory.Telemetry()

	// Assert
	assert.NotNil(t, result)
	require.NoError(t, err)
}

func TestFactory_Telemetry_LoggerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		nil,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	result, err := factory.Telemetry()

	// Assert
	assert.Nil(t, result)
	require.ErrorIs(t, err, expectedError)
}

func TestFactory_Telemetry_ConfigError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		nil,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	result, err := factory.Telemetry()

	// Assert
	assert.Nil(t, result)
	require.ErrorIs(t, err, expectedError)
}

func TestFactory_Telemetry_DirectoryError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &telemetrymocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(nil, expectedError).
		Once()

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		mockDirectoryFactory,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	result, err := factory.Telemetry()

	// Assert
	assert.Nil(t, result)
	require.ErrorIs(t, err, expectedError)
}

func TestFactory_Telemetry_TelemetryDisabled(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &telemetrymocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockInt64Counter := &instrumentsmocks.MockInt64Counter{}
	defer mockInt64Counter.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return("localhost:4317").
		Once()

	mockConfig.EXPECT().
		DisableTelemetry().
		Return(true).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(
			mock.MatchedBy(func(m metric.Meter) bool {
				_, ok := m.(noop.Meter)
				return ok
			}),
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(mockInt64Counter, nil).
		Times(3)

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		mockDirectoryFactory,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	result, err := factory.Telemetry()

	// Assert
	assert.NotNil(t, result)
	require.NoError(t, err)
	assert.Contains(t, testLogger.InfoLogs(), "Telemetry disabled by configuration")
}

func TestFactory_Telemetry_EmptyCollectorEndpoint(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &telemetrymocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockInt64Counter := &instrumentsmocks.MockInt64Counter{}
	defer mockInt64Counter.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockConfig.EXPECT().
		DisableTelemetry().
		Return(false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return("").
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(
			mock.MatchedBy(func(m metric.Meter) bool {
				_, ok := m.(noop.Meter)
				return ok
			}),
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(mockInt64Counter, nil).
		Times(3)

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		mockDirectoryFactory,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	result, err := factory.Telemetry()

	// Assert
	assert.NotNil(t, result)
	require.NoError(t, err)
	assert.Contains(t, testLogger.InfoLogs(), "Telemetry off: no collector endpoint set")
}

func TestFactory_Telemetry_ExporterCreationError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &telemetrymocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockConfig.EXPECT().
		DisableTelemetry().
		Return(false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return("localhost:4317").
		Once()

	mockExporterFactory.EXPECT().
		New().
		Return(nil, expectedError).
		Once()

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		mockDirectoryFactory,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	result, err := factory.Telemetry()

	// Assert
	assert.Nil(t, result)
	require.ErrorIs(t, err, expectedError)
}

func TestFactory_Telemetry_MeterProviderCreationError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &telemetrymocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockConfig.EXPECT().
		DisableTelemetry().
		Return(false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return("localhost:4317").
		Once()

	mockExporterFactory.EXPECT().
		New().
		Return(mockExporter, nil).
		Once()

	mockServerDefinition.EXPECT().
		Name().
		Return("matlab-mcp-server").
		Once()

	mockConfig.EXPECT().
		Version().
		Return("1.0.0").
		Once()

	mockMeterProviderFactory.EXPECT().
		New(mockExporter, mock.Anything, "matlab-mcp-server", "1.0.0").
		Return(nil, expectedError).
		Once()

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		mockDirectoryFactory,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	result, err := factory.Telemetry()

	// Assert
	assert.Nil(t, result)
	require.ErrorIs(t, err, expectedError)
}

func TestFactory_Telemetry_InstrumentCreationError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &telemetrymocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	noopMeterProvider := noop.NewMeterProvider()
	instrumentError := assert.AnError
	expectedError := messages.New_StartupErrors_TelemetryInitializationFailed_Error()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockConfig.EXPECT().
		DisableTelemetry().
		Return(false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return("localhost:4317").
		Once()

	mockExporterFactory.EXPECT().
		New().
		Return(mockExporter, nil).
		Once()

	mockServerDefinition.EXPECT().
		Name().
		Return("matlab-mcp-server").
		Once()

	mockConfig.EXPECT().
		Version().
		Return("1.0.0").
		Once()

	mockMeterProviderFactory.EXPECT().
		New(mockExporter, mock.Anything, "matlab-mcp-server", "1.0.0").
		Return(noopMeterProvider, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, instrumentError).
		Once()

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		mockDirectoryFactory,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	result, err := factory.Telemetry()

	// Assert
	assert.Nil(t, result)
	require.Equal(t, expectedError, err)
}

func TestFactory_Telemetry_IsSingleton(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &telemetrymocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	mockInt64Counter := &instrumentsmocks.MockInt64Counter{}
	defer mockInt64Counter.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	noopMeterProvider := noop.NewMeterProvider()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockConfig.EXPECT().
		DisableTelemetry().
		Return(false).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectorEndpoint().
		Return("localhost:4317").
		Once()

	mockExporterFactory.EXPECT().
		New().
		Return(mockExporter, nil).
		Once()

	mockServerDefinition.EXPECT().
		Name().
		Return("matlab-mcp-server").
		Once()

	mockConfig.EXPECT().
		Version().
		Return("1.0.0").
		Once()

	mockMeterProviderFactory.EXPECT().
		New(mockExporter, mock.Anything, "matlab-mcp-server", "1.0.0").
		Return(noopMeterProvider, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockInt64Counter, nil).
		Times(3)

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		mockDirectoryFactory,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	firstResult, firstErr := factory.Telemetry()
	secondResult, secondErr := factory.Telemetry()

	// Assert
	assert.NotNil(t, firstResult)
	require.NoError(t, firstErr)
	assert.Equal(t, firstResult, secondResult)
	require.NoError(t, secondErr)
}

func TestFactory_Telemetry_SingletonPreservesError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &telemetrymocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &telemetrymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockExporterFactory := &telemetrymocks.MockExporterFactory{}
	defer mockExporterFactory.AssertExpectations(t)

	mockMeterProviderFactory := &telemetrymocks.MockMeterProviderFactory{}
	defer mockMeterProviderFactory.AssertExpectations(t)

	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockCorrelationIDProvider := &telemetrymocks.MockSessionCorrelationIDProvider{}
	defer mockCorrelationIDProvider.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := telemetry.NewFactory(
		mockLoggerFactory,
		mockConfigFactory,
		mockExporterFactory,
		mockMeterProviderFactory,
		mockInstrumentFactory,
		nil,
		mockOSLayer,
		mockOSVersionProvider,
		mockServerDefinition,
		mockCorrelationIDProvider,
	)

	// Act
	firstResult, firstErr := factory.Telemetry()
	secondResult, secondErr := factory.Telemetry()

	// Assert
	assert.Nil(t, firstResult)
	require.ErrorIs(t, firstErr, expectedError)
	assert.Nil(t, secondResult)
	require.ErrorIs(t, secondErr, expectedError)
}

func TestToolSource_String(t *testing.T) {
	cases := map[telemetry.ToolSource]string{
		telemetry.ToolSourceBuiltin:   "builtin",
		telemetry.ToolSourceExtension: "extension",
		telemetry.ToolSource(2):       "unknown",
	}

	for source, expected := range cases {
		assert.Equal(t, expected, source.String())
	}
}
