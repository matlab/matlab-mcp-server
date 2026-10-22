// Copyright 2026 The MathWorks, Inc.

package provider_test

import (
	"context"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel/meter/provider"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	otelmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry/otel"
	providermocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry/otel/meter/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const requiredShutdownTimeout = 3 * time.Second

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	// Act
	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockErrorHandler := &otelmocks.MockErrorHandler{}
	defer mockErrorHandler.AssertExpectations(t)

	expectedCollectionInterval := 1 * time.Minute

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectionInterval().
		Return(expectedCollectionInterval).
		Once()

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Once()

	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	// Act
	result, err := factory.New(mockExporter, mockErrorHandler, "matlab-mcp-server", "1.0.0")

	// Assert
	require.NoError(t, err, "Error should be nil")
	require.NotNil(t, result, "MeterProvider should not be nil")
}

func TestFactory_New_LoggerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	mockErrorHandler := &otelmocks.MockErrorHandler{}
	defer mockErrorHandler.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	// Act
	result, err := factory.New(mockExporter, mockErrorHandler, "matlab-mcp-server", "1.0.0")

	// Assert
	require.Nil(t, result)
	require.ErrorIs(t, err, expectedError)
}

func TestFactory_New_ShutdownPassesContextWithDeadline(t *testing.T) {
	// Arrange
	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockErrorHandler := &otelmocks.MockErrorHandler{}
	defer mockErrorHandler.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectionInterval().
		Return(time.Hour).
		Once()

	var shutdownCtx context.Context
	mockExporter.EXPECT().
		Export(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, _ *metricdata.ResourceMetrics) {
			shutdownCtx = ctx
		}).
		Return(nil).
		Once()

	mockExporter.EXPECT().
		Shutdown(mock.Anything).
		Return(nil).
		Once()

	var capturedShutdown func() error
	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Run(func(fn func() error) {
			capturedShutdown = fn
		}).
		Once()

	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	_, err := factory.New(mockExporter, mockErrorHandler, "matlab-mcp-server", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, capturedShutdown)

	// Act
	before := time.Now()
	_ = capturedShutdown()
	after := time.Now()

	// Assert
	require.NotNil(t, shutdownCtx, "exporter should have been called during shutdown")
	deadline, hasDeadline := shutdownCtx.Deadline()
	require.True(t, hasDeadline)
	assert.WithinRange(t, deadline, before.Add(requiredShutdownTimeout), after.Add(requiredShutdownTimeout))
}

func TestFactory_New_ShutdownErrorReachesTheErrorHandler(t *testing.T) {
	// Arrange
	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockErrorHandler := &otelmocks.MockErrorHandler{}
	defer mockErrorHandler.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectionInterval().
		Return(time.Hour).
		Once()

	mockExporter.EXPECT().
		Export(mock.Anything, mock.Anything).
		Return(assert.AnError).
		Once()

	var handledErr error
	mockErrorHandler.EXPECT().
		Handle(mock.Anything).
		Run(func(err error) {
			handledErr = err
		}).
		Once()

	mockExporter.EXPECT().
		Shutdown(mock.Anything).
		Return(nil).
		Once()

	var capturedShutdown func() error
	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Run(func(fn func() error) {
			capturedShutdown = fn
		}).
		Once()

	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	_, err := factory.New(mockExporter, mockErrorHandler, "matlab-mcp-server", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, capturedShutdown)

	// Act
	shutdownErr := capturedShutdown()

	// Assert
	require.NoError(t, shutdownErr, "a failed telemetry flush must not fail application shutdown")
	require.ErrorIs(t, handledErr, assert.AnError, "the failed flush should reach the error handler")
}

func TestFactory_New_ConfigError(t *testing.T) {
	// Arrange
	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	mockErrorHandler := &otelmocks.MockErrorHandler{}
	defer mockErrorHandler.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	// Act
	result, err := factory.New(mockExporter, mockErrorHandler, "matlab-mcp-server", "1.0.0")

	// Assert
	require.Nil(t, result, "MeterProvider should be nil when config fails")
	require.ErrorIs(t, err, expectedError)
}
