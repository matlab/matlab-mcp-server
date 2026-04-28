// Copyright 2025-2026 The MathWorks, Inc.

package modeselector_test

import (
	"fmt"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/modeselector"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	modeselectormocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/modeselector"
	telemetrymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/telemetry"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	// Act
	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Assert
	assert.NotNil(t, modeSelectorInstance, "ModeSelector instance should not be nil")
}

func TestStartAndWaitForCompletion_ConfigError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Config")
}

func TestStartAndWaitForCompletion_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from GetGlobalLogger")
}

func TestStartAndWaitForCompletion_TelemetryError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(nil, expectedError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Telemetry")
}

func TestStartAndWaitForCompletion_VersionMode_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedCtx := t.Context()
	expectedVersion := "25.6.68"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(true).
		Once()

	mockOsLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockStdout.EXPECT().
		Write([]byte(fmt.Sprintf("%s\n", expectedVersion))).
		Return(len(expectedVersion)+1, nil).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.NoError(t, err, "StartAndWaitForCompletion should not return an error in version mode")
}

func TestStartAndWaitForCompletion_VersionMode_WriteError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedCtx := t.Context()
	expectedVersion := "25.6.68"
	writeError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(true).
		Once()

	mockOsLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockStdout.EXPECT().
		Write([]byte(fmt.Sprintf("%s\n", expectedVersion))).
		Return(0, writeError).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	var writeErr *messages.StartupErrors_WriteError_Error
	require.ErrorAs(t, err, &writeErr)
	require.Equal(t, "version", writeErr.Attr0)
	require.Equal(t, writeError.Error(), writeErr.Attr1)
}

func TestStartAndWaitForCompletion_VersionMode_ShutdownError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedCtx := t.Context()
	expectedVersion := "25.6.68"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(true).
		Once()

	mockOsLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockStdout.EXPECT().
		Write([]byte(fmt.Sprintf("%s\n", expectedVersion))).
		Return(len(expectedVersion)+1, nil).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(assert.AnError).
		Once()

	mockLogger.EXPECT().
		WithError(assert.AnError).
		Return(mockLogger).
		Once()

	mockLogger.EXPECT().
		Warn("Shutdown failed").
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.NoError(t, err, "StartAndWaitForCompletion should not return an error when only shutdown fails")
}

func TestStartAndWaitForCompletion_WatchdogMode_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(true).
		Once()

	mockWatchdogProcess.EXPECT().
		StartAndWaitForCompletion(expectedCtx).
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.NoError(t, err, "StartAndWaitForCompletion should not return an error in watchdog mode")
}

func TestStartAndWaitForCompletion_WatchdogMode_StartAndWaitError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	watchdogError := assert.AnError
	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(true).
		Once()

	mockWatchdogProcess.EXPECT().
		StartAndWaitForCompletion(expectedCtx).
		Return(watchdogError).
		Once()

	mockLogger.EXPECT().
		WithError(watchdogError).
		Return(mockLogger).
		Once()

	mockLogger.EXPECT().
		Error("Server failed with unexpected error").
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	var genericErr *messages.StartupErrors_GenericInitializeFailure_Error
	require.ErrorAs(t, err, &genericErr)
}

func TestStartAndWaitForCompletion_SetupMATLABMode_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		SetupMATLABMode().
		Return(true).
		Once()

	mockSetupMATLAB.EXPECT().
		StartAndWaitForCompletion(expectedCtx).
		Return(nil).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.NoError(t, err, "StartAndWaitForCompletion should not return an error in install MATLAB add-on mode")
}

func TestStartAndWaitForCompletion_SetupMATLABMode_Error(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedError := messages.AnError
	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		SetupMATLABMode().
		Return(true).
		Once()

	mockSetupMATLAB.EXPECT().
		StartAndWaitForCompletion(expectedCtx).
		Return(expectedError).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from SetupMATLAB")
}

func TestStartAndWaitForCompletion_DefaultMode_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		SetupMATLABMode().
		Return(false).
		Once()

	mockOrchestrator.EXPECT().
		StartAndWaitForCompletion(expectedCtx).
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.NoError(t, err, "StartAndWaitForCompletion should not return an error in default mode")
}

func TestStartAndWaitForCompletion_DefaultMode_StartAndWaitError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	orchestratorError := assert.AnError
	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		SetupMATLABMode().
		Return(false).
		Once()

	mockOrchestrator.EXPECT().
		StartAndWaitForCompletion(expectedCtx).
		Return(orchestratorError).
		Once()

	mockLogger.EXPECT().
		WithError(orchestratorError).
		Return(mockLogger).
		Once()

	mockLogger.EXPECT().
		Error("Server failed with unexpected error").
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	var genericErr *messages.StartupErrors_GenericInitializeFailure_Error
	require.ErrorAs(t, err, &genericErr)
}

func TestStartAndWaitForCompletion_DefaultMode_StartAndWaitMessagesError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedError := messages.AnError
	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		SetupMATLABMode().
		Return(false).
		Once()

	mockOrchestrator.EXPECT().
		StartAndWaitForCompletion(expectedCtx).
		Return(expectedError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should pass through messages.Error without wrapping")
}

func TestStartAndWaitForCompletion_HelpMode_StartAndWaitHappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	helpText := "Help me get my feet back on the ground."
	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(true).
		Once()

	mockParser.EXPECT().
		Usage().
		Return(helpText, nil).
		Once()

	mockOsLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockStdout.EXPECT().
		Write([]byte(fmt.Sprintf("%s\n", helpText))).
		Return(len(helpText)+1, nil).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.NoError(t, err)
}

func TestStartAndWaitForCompletion_HelpMode_UsageError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(true).
		Once()

	mockParser.EXPECT().
		Usage().
		Return("", messages.AnError).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	require.ErrorIs(t, err, messages.AnError)
}

func TestStartAndWaitForCompletion_HelpMode_StartAndWaitWriteError(t *testing.T) {
	// Arrange
	mockConfigFactory := &modeselectormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockTelemetryFactory := &modeselectormocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockWatchdogProcess := &modeselectormocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockOrchestrator := &modeselectormocks.MockOrchestrator{}
	defer mockOrchestrator.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockLoggerFactory := &modeselectormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockLifecycleSignaler := &modeselectormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockSetupMATLAB := &modeselectormocks.MockSetupMATLAB{}
	defer mockSetupMATLAB.AssertExpectations(t)

	helpText := "Help me get my feet back on the ground."
	writeError := assert.AnError
	expectedCtx := t.Context()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordServerStart(expectedCtx).
		Once()

	mockConfig.EXPECT().
		HelpMode().
		Return(true).
		Once()

	mockParser.EXPECT().
		Usage().
		Return(helpText, nil).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockOsLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockStdout.EXPECT().
		Write([]byte(fmt.Sprintf("%s\n", helpText))).
		Return(0, writeError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfigFactory,
		mockParser,
		mockTelemetryFactory,
		mockWatchdogProcess,
		mockOrchestrator,
		mockOsLayer,
		mockLifecycleSignaler,
		mockLoggerFactory,
		mockSetupMATLAB,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(expectedCtx)

	// Assert
	var writeErr *messages.StartupErrors_WriteError_Error
	require.ErrorAs(t, err, &writeErr)
	require.Equal(t, "help", writeErr.Attr0)
	require.Equal(t, writeError.Error(), writeErr.Attr1)
}
