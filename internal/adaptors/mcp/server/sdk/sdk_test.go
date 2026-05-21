// Copyright 2025-2026 The MathWorks, Inc.

package sdk_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/sdk"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/server/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	// Act
	factory := sdk.NewFactory(mockConfigFactory, mockDefinition, mockRootStore, mockLoggerFactory, mockGlobalMATLAB)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_NewServer_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedVersion := "1.0.0"
	expectedName := "test-server"
	expectedTitle := "Test Server"
	expectedInstructions := "test instructions"

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockDefinition.EXPECT().
		Features().
		Return(definition.Features{}).
		Once()

	mockDefinition.EXPECT().
		Name().
		Return(expectedName).
		Once()

	mockDefinition.EXPECT().
		Title().
		Return(expectedTitle).
		Once()

	mockDefinition.EXPECT().
		Instructions().
		Return(expectedInstructions).
		Once()

	factory := sdk.NewFactory(mockConfigFactory, mockDefinition, mockRootStore, mockLoggerFactory, mockGlobalMATLAB)

	// Act
	server, err := factory.NewServer()

	// Assert
	require.NoError(t, err, "NewServer should not return an error")
	assert.NotNil(t, server, "Server should not be nil")
}

func TestFactory_NewServer_ConfigError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := sdk.NewFactory(mockConfigFactory, mockDefinition, mockRootStore, mockLoggerFactory, mockGlobalMATLAB)

	// Act
	server, err := factory.NewServer()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, server, "Server should be nil when error occurs")
}

func TestFactory_NewServer_LoggerError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	factory := sdk.NewFactory(mockConfigFactory, mockDefinition, mockRootStore, mockLoggerFactory, mockGlobalMATLAB)

	// Act
	server, err := factory.NewServer()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, server, "Server should be nil when error occurs")
}

func TestHandleInitialized_EagerMATLABInit_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(true).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, nil).
		Once()

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	assert.Empty(t, mockLogger.WarnLogs(), "no warnings should be logged on successful eager initialization")
}

func TestHandleInitialized_NilRequest(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	handler(ctx, nil)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestHandleInitialized_NilSession(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	handler(ctx, &mcp.InitializedRequest{})

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLogClientDetails_HappyPath(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedClientName := "test-client"
	expectedClientTitle := "Test Client"
	expectedClientURL := "https://example.com"
	expectedClientVersion := "1.2.3"

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			ClientInfo: &mcp.Implementation{
				Name:       expectedClientName,
				Title:      expectedClientTitle,
				WebsiteURL: expectedClientURL,
				Version:    expectedClientVersion,
			},
		}).
		Once()

	logClientDetails := sdk.LogClientDetails(mockLogger)

	// Act
	logClientDetails(mockSession)

	// Assert
	logs := mockLogger.InfoLogs()
	fields, found := logs["New client session"]
	require.True(t, found, "Expected info log for new client session")
	assert.Equal(t, expectedClientName, fields["client-name"])
	assert.Equal(t, expectedClientTitle, fields["client-title"])
	assert.Equal(t, expectedClientURL, fields["client-url"])
	assert.Equal(t, expectedClientVersion, fields["client-version"])
}

func TestLogClientDetails_NilInitializeParams(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockSession.EXPECT().
		InitializeParams().
		Return(nil).
		Once()

	logClientDetails := sdk.LogClientDetails(mockLogger)

	// Act
	logClientDetails(mockSession)

	// Assert
	assert.Empty(t, mockLogger.InfoLogs(), "No info logs should be emitted when InitializeParams is nil")
}

func TestLogClientDetails_NilClientInfo(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{}).
		Once()

	logClientDetails := sdk.LogClientDetails(mockLogger)

	// Act
	logClientDetails(mockSession)

	// Assert
	assert.Empty(t, mockLogger.InfoLogs(), "No info logs should be emitted when ClientInfo is nil")
}

func TestHandleInitialized_EagerMATLABInit_MATLABFeatureDisabled(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: false}}

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	mockGlobalMATLAB.AssertNotCalled(t, "Client")
}

func TestHandleInitialized_EagerMATLABInit_MultipleSession(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(false).
		Once()

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	mockGlobalMATLAB.AssertNotCalled(t, "Client")
}

func TestHandleInitialized_EagerMATLABInit_InitializeMATLABOnStartupFalse(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(false).
		Once()

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	mockGlobalMATLAB.AssertNotCalled(t, "Client")
}

func TestHandleInitialized_EagerMATLABInit_ErrorLogsWarning(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(true).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	logs := mockLogger.WarnLogs()
	fields, found := logs["MATLAB eager initialization failed"]
	require.True(t, found, "Expected warning log for MATLAB initialization failure")
	assert.Equal(t, expectedError, fields["error"])
}

func TestUpdateRoots_NilInitializeParams_ReturnsNil(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	mockSession.EXPECT().
		InitializeParams().
		Return(nil).
		Once()

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.NoError(t, err)
	mockRootStore.AssertNotCalled(t, "UpdateRoots")
}

func TestUpdateRoots_NilCapabilities_ReturnsNil(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{}).
		Once()

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.NoError(t, err)
	mockRootStore.AssertNotCalled(t, "UpdateRoots")
}

func TestUpdateRoots_NilRootsV2_ReturnsNil(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			Capabilities: &mcp.ClientCapabilities{},
		}).
		Once()

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.NoError(t, err)
	mockRootStore.AssertNotCalled(t, "UpdateRoots")
}

func TestUpdateRoots_ListRootsError_ReturnsError(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError
	features := definition.Features{}

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			Capabilities: &mcp.ClientCapabilities{
				RootsV2: &mcp.RootCapabilities{},
			},
		}).
		Once()

	mockSession.EXPECT().
		ListRoots(ctx, (*mcp.ListRootsParams)(nil)).
		Return(nil, expectedError).
		Once()

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.ErrorIs(t, err, expectedError)
	mockRootStore.AssertNotCalled(t, "UpdateRoots")
}

func TestUpdateRoots_HappyPath_UpdatesRootStore(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}
	expectedRoots := []*mcp.Root{
		{URI: "file:///home/user/project", Name: "project"},
	}

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			Capabilities: &mcp.ClientCapabilities{
				RootsV2: &mcp.RootCapabilities{},
			},
		}).
		Once()

	mockSession.EXPECT().
		ListRoots(ctx, (*mcp.ListRootsParams)(nil)).
		Return(&mcp.ListRootsResult{Roots: expectedRoots}, nil).
		Once()

	mockRootStore.EXPECT().
		UpdateRoots(expectedRoots).
		Once()

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.NoError(t, err)
}
