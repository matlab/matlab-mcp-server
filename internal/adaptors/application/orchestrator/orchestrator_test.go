// Copyright 2025-2026 The MathWorks, Inc.

package orchestrator_test

import (
	"os"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/orchestrator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	directorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	orchestratormocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/orchestrator"
	toolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	//Act
	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Assert
	assert.NotNil(t, orchestratorInstance, "Orchestrator instance should not be nil")
}

func TestOrchestrator_StartAndWaitForCompletion_ConfigError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	ctx := t.Context()
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Config")
}

func TestOrchestrator_StartAndWaitForCompletion_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	ctx := t.Context()
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from GetGlobalLogger")
}

func TestOrchestrator_StartAndWaitForCompletion_CapOpenFilesLimitError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockResourceLimitManager.EXPECT().
		CapOpenFilesLimit(orchestrator.UnixOpenFileDescriptorsSoftCap).
		Return(nil, expectedError).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from CapOpenFilesLimit")
}

func TestOrchestrator_StartAndWaitForCompletion_DirectoryError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockResourceLimitManager.EXPECT().
		CapOpenFilesLimit(orchestrator.UnixOpenFileDescriptorsSoftCap).
		Return(func() error { return nil }, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return("test-version").
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(nil, expectedError).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Directory")
}

func TestOrchestrator_StartAndWaitForCompletion_WatchdogStartError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockResourceLimitManager.EXPECT().
		CapOpenFilesLimit(orchestrator.UnixOpenFileDescriptorsSoftCap).
		Return(func() error { return nil }, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return("test-version").
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(expectedError).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from watchdogClient.Start")
}

func TestOrchestrator_StartAndWaitForCompletion_DependenciesError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError
	expectedDependenciesProviderResources := definition.NewDependenciesProviderResources(mockLogger, mockConfig, mockMessageCatalog, mockWatchdogClient)

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockResourceLimitManager.EXPECT().
		CapOpenFilesLimit(orchestrator.UnixOpenFileDescriptorsSoftCap).
		Return(func() error { return nil }, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return("test-version").
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockApplicationDefinition.EXPECT().
		Dependencies(expectedDependenciesProviderResources).
		Return(nil, expectedError).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Dependencies")
}

func TestOrchestrator_StartAndWaitForCompletion_HappyPath(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockTool := &toolsmocks.MockTool{}
	defer mockTool.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()
	serverStarted := make(chan struct{})
	stopServer := make(chan struct{})
	defer close(stopServer)

	expectedDependencies := &struct{}{}
	expectedDependenciesProviderResources := definition.NewDependenciesProviderResources(mockLogger, mockConfig, mockMessageCatalog, mockWatchdogClient)
	expectedToolProviderResources := definition.NewToolsProviderResources(mockLogger, mockConfig, mockMessageCatalog, expectedDependencies, mockLoggerFactory)
	expectedTools := []tools.Tool{mockTool}
	expectedVersion := "test-version"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	wasResetOpenFilesCapCalled := false
	mockResourceLimitManager.EXPECT().
		CapOpenFilesLimit(orchestrator.UnixOpenFileDescriptorsSoftCap).
		Return(func() error {
			wasResetOpenFilesCapCalled = true
			return nil
		}, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockApplicationDefinition.EXPECT().
		Dependencies(expectedDependenciesProviderResources).
		Return(expectedDependencies, nil).
		Once()

	mockApplicationDefinition.EXPECT().
		Tools(expectedToolProviderResources).
		Return(expectedTools).
		Once()

	// Server should run indefinitely (simulate with a blocking channel)
	mockServer.EXPECT().
		Run(expectedTools).
		RunAndReturn(func(_ []tools.Tool) error {
			close(serverStarted)
			<-stopServer
			return nil
		}).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	<-serverStarted

	sendInterruptSignal(interruptC)

	// Assert
	require.NoError(t, <-errC, "StartAndWaitForCompletion should not return an error on signal interrupt")

	logs := mockLogger.InfoLogs()
	fields, found := logs["Initiating application startup"]
	require.True(t, found, "Expected info log for application startup")
	require.True(t, wasResetOpenFilesCapCalled, "Expected to see CapOpenFilesLimit callback to be called")
	assert.Equal(t, expectedVersion, fields["version"])
}

func TestOrchestrator_StartAndWaitForCompletion_ServerError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()
	expectedError := assert.AnError
	var expectedDependencies any
	expectedDependenciesProviderResources := definition.NewDependenciesProviderResources(mockLogger, mockConfig, mockMessageCatalog, mockWatchdogClient)
	expectedToolProviderResources := definition.NewToolsProviderResources(mockLogger, mockConfig, mockMessageCatalog, expectedDependencies, mockLoggerFactory)
	var expectedTools []tools.Tool

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockResourceLimitManager.EXPECT().
		CapOpenFilesLimit(orchestrator.UnixOpenFileDescriptorsSoftCap).
		Return(func() error { return nil }, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return("test-version").
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockApplicationDefinition.EXPECT().
		Dependencies(expectedDependenciesProviderResources).
		Return(expectedDependencies, nil).
		Once()

	mockApplicationDefinition.EXPECT().
		Tools(expectedToolProviderResources).
		Return(expectedTools).
		Once()

	mockServer.EXPECT().
		Run(expectedTools).
		Return(expectedError).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the server error")
}

func TestOrchestrator_StartAndWaitForCompletion_WaitForShutdownToCompleteError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()
	expectedError := assert.AnError
	var expectedDependencies any
	expectedDependenciesProviderResources := definition.NewDependenciesProviderResources(mockLogger, mockConfig, mockMessageCatalog, mockWatchdogClient)
	expectedToolProviderResources := definition.NewToolsProviderResources(mockLogger, mockConfig, mockMessageCatalog, expectedDependencies, mockLoggerFactory)
	var expectedTools []tools.Tool

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockResourceLimitManager.EXPECT().
		CapOpenFilesLimit(orchestrator.UnixOpenFileDescriptorsSoftCap).
		Return(func() error { return nil }, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return("test-version").
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockApplicationDefinition.EXPECT().
		Dependencies(expectedDependenciesProviderResources).
		Return(expectedDependencies, nil).
		Once()

	mockApplicationDefinition.EXPECT().
		Tools(expectedToolProviderResources).
		Return(expectedTools).
		Once()

	mockServer.EXPECT().
		Run(expectedTools).
		Return(nil).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(expectedError).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	// Assert
	require.NoError(t, <-errC, "StartAndWaitForCompletion should not return an error on signal interrupt")

	// This is mostly optional
	logs := mockLogger.WarnLogs()

	fields, found := logs["Application shutdown failed"]
	require.True(t, found, "Expected a warning log about shutdown failure")

	errField, found := fields["error"]
	require.True(t, found, "Expected an error field in the warning log")

	err, ok := errField.(error)
	require.True(t, ok, "Error field should be of type error")
	require.ErrorIs(t, err, expectedError, "Logged error should match the shutdown error")
}

func TestOrchestrator_StartAndWaitForCompletion_WatchdogStopError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockApplicationDefinition := &orchestratormocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockResourceLimitManager := &orchestratormocks.MockResourceLimitManager{}
	defer mockResourceLimitManager.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	ctx := t.Context()
	interruptC := getInterruptChannel()
	expectedError := assert.AnError
	var expectedDependencies any
	expectedDependenciesProviderResources := definition.NewDependenciesProviderResources(mockLogger, mockConfig, mockMessageCatalog, mockWatchdogClient)
	expectedToolProviderResources := definition.NewToolsProviderResources(mockLogger, mockConfig, mockMessageCatalog, expectedDependencies, mockLoggerFactory)
	var expectedTools []tools.Tool

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockResourceLimitManager.EXPECT().
		CapOpenFilesLimit(orchestrator.UnixOpenFileDescriptorsSoftCap).
		Return(func() error { return nil }, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return("test-version").
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	serverStarted := make(chan struct{})
	stopServer := make(chan struct{})
	defer close(stopServer)

	mockApplicationDefinition.EXPECT().
		Dependencies(expectedDependenciesProviderResources).
		Return(expectedDependencies, nil).
		Once()

	mockApplicationDefinition.EXPECT().
		Tools(expectedToolProviderResources).
		Return(expectedTools).
		Once()

	mockServer.EXPECT().
		Run(expectedTools).
		RunAndReturn(func(_ []tools.Tool) error {
			close(serverStarted)
			<-stopServer
			return nil
		}).
		Once()

	// Signal
	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	// Shutdown sequence
	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	// Watchdog Stop Fails
	mockWatchdogClient.EXPECT().
		Stop().
		Return(expectedError).
		Once()

	orchestratorInstance := orchestrator.New(
		mockMessageCatalog,
		mockLifecycleSignaler,
		mockApplicationDefinition,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockDirectoryFactory,
		mockResourceLimitManager,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	<-serverStarted
	sendInterruptSignal(interruptC)
	err := <-errC

	// Assert
	require.NoError(t, err)

	// Verify Log
	logs := mockLogger.WarnLogs()
	fields, found := logs["Watchdog shutdown failed"]
	require.True(t, found, "Expected warning log for watchdog failure")
	assert.Equal(t, expectedError, fields["error"])
}

func getInterruptChannel() chan os.Signal {
	return make(chan os.Signal, 1)
}

func sendInterruptSignal(interruptC chan os.Signal) {
	interruptC <- os.Interrupt
}
