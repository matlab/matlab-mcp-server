// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	publictypes "github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/publictypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/server"
	internalentities "github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/wire/adaptor"
	internaltoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
	publictypesmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/sdk/publictypes"
	servermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/sdk/server"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	adaptormocks "github.com/matlab/matlab-mcp-core-server/mocks/wire/adaptor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServer_StartAndWaitForCompletion_HappyPath(t *testing.T) {
	// Arrange
	mockDependenciesProviderFactory := &servermocks.MockDependenciesProviderFactory[struct{}]{}
	defer mockDependenciesProviderFactory.AssertExpectations(t)

	mockToolsProviderFactory := &servermocks.MockToolsProviderFactory[struct{}]{}
	defer mockToolsProviderFactory.AssertExpectations(t)

	mockParametersFactory := &servermocks.MockParametersFactory{}
	defer mockParametersFactory.AssertExpectations(t)

	mockFeaturesFactory := &servermocks.MockFeaturesFactory{}
	defer mockFeaturesFactory.AssertExpectations(t)

	mockApplicationFactory := &servermocks.MockApplicationFactory{}
	defer mockApplicationFactory.AssertExpectations(t)

	mockApplication := &adaptormocks.MockApplication{}
	defer mockApplication.AssertExpectations(t)

	mockModeSelector := &adaptormocks.MockModeSelector{}
	defer mockModeSelector.AssertExpectations(t)

	mockErrorWriter := &entitiesmocks.MockWriter{}
	defer mockErrorWriter.AssertExpectations(t)

	ctx := t.Context()
	expectedName := "test-server"
	expectedTitle := "Test Server"
	expectedInstructions := "Test instructions"
	expectedFeatures := publictypes.Features{
		MATLAB: publictypes.MATLABFeature{Enabled: true},
	}

	serverDefinition := server.Definition[struct{}]{
		Name:         expectedName,
		Title:        expectedTitle,
		Instructions: expectedInstructions,
		Features:     expectedFeatures,
	}

	expectedInternalFeatures := definition.Features{
		MATLAB: definition.MATLABFeature{Enabled: true},
	}
	expectedDefinition := definition.New(
		expectedName,
		expectedTitle,
		expectedInstructions,
		expectedInternalFeatures,
		nil,
		nil,
		nil,
	)

	mockFeaturesFactory.EXPECT().
		New(expectedFeatures).
		Return(expectedInternalFeatures).
		Once()

	mockParametersFactory.EXPECT().
		New([]publictypes.Parameter(nil)).
		Return(nil).
		Once()

	mockDependenciesProviderFactory.EXPECT().
		New(mock.Anything).
		Return(nil).
		Once()

	mockToolsProviderFactory.EXPECT().
		New(mock.Anything).
		Return(nil).
		Once()

	mockApplicationFactory.EXPECT().
		New(expectedDefinition).
		Return(mockApplication).
		Once()

	mockApplication.EXPECT().
		ModeSelector().
		Return(mockModeSelector).
		Once()

	mockModeSelector.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(nil).
		Once()

	s := server.New(serverDefinition, mockFeaturesFactory, mockParametersFactory, mockDependenciesProviderFactory, mockToolsProviderFactory, mockApplicationFactory, mockErrorWriter)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 0, exitCode)
}

func TestServer_StartAndWaitForCompletion_WithParameters(t *testing.T) {
	// Arrange
	mockDependenciesProviderFactory := &servermocks.MockDependenciesProviderFactory[struct{}]{}
	defer mockDependenciesProviderFactory.AssertExpectations(t)

	mockToolsProviderFactory := &servermocks.MockToolsProviderFactory[struct{}]{}
	defer mockToolsProviderFactory.AssertExpectations(t)

	mockParametersFactory := &servermocks.MockParametersFactory{}
	defer mockParametersFactory.AssertExpectations(t)

	mockFeaturesFactory := &servermocks.MockFeaturesFactory{}
	defer mockFeaturesFactory.AssertExpectations(t)

	mockApplicationFactory := &servermocks.MockApplicationFactory{}
	defer mockApplicationFactory.AssertExpectations(t)

	mockApplication := &adaptormocks.MockApplication{}
	defer mockApplication.AssertExpectations(t)

	mockModeSelector := &adaptormocks.MockModeSelector{}
	defer mockModeSelector.AssertExpectations(t)

	mockErrorWriter := &entitiesmocks.MockWriter{}
	defer mockErrorWriter.AssertExpectations(t)

	mockParameter := &publictypesmocks.MockParameter{}
	defer mockParameter.AssertExpectations(t)

	mockInternalParameter := &entitiesmocks.MockParameter{}
	defer mockInternalParameter.AssertExpectations(t)

	ctx := t.Context()
	expectedName := "test-server"
	expectedTitle := "Test Server"
	expectedInstructions := "Test instructions"

	serverDefinition := server.Definition[struct{}]{
		Name:         expectedName,
		Title:        expectedTitle,
		Instructions: expectedInstructions,
		Parameters:   []publictypes.Parameter{mockParameter},
	}

	expectedInternalParameters := []internalentities.Parameter{mockInternalParameter}
	expectedDefinition := definition.New(
		expectedName,
		expectedTitle,
		expectedInstructions,
		definition.Features{},
		expectedInternalParameters,
		nil,
		nil,
	)

	mockFeaturesFactory.EXPECT().
		New(publictypes.Features{}).
		Return(definition.Features{}).
		Once()

	mockParametersFactory.EXPECT().
		New([]publictypes.Parameter{mockParameter}).
		Return(expectedInternalParameters).
		Once()

	mockDependenciesProviderFactory.EXPECT().
		New(mock.Anything).
		Return(nil).
		Once()

	mockToolsProviderFactory.EXPECT().
		New(mock.Anything).
		Return(nil).
		Once()

	mockApplicationFactory.EXPECT().
		New(expectedDefinition).
		Return(mockApplication).
		Once()

	mockApplication.EXPECT().
		ModeSelector().
		Return(mockModeSelector).
		Once()

	mockModeSelector.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(nil).
		Once()

	s := server.New(serverDefinition, mockFeaturesFactory, mockParametersFactory, mockDependenciesProviderFactory, mockToolsProviderFactory, mockApplicationFactory, mockErrorWriter)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 0, exitCode)
}

func TestServer_StartAndWaitForCompletion_WithProviders(t *testing.T) {
	// Arrange
	mockDependenciesProviderFactory := &servermocks.MockDependenciesProviderFactory[struct{}]{}
	defer mockDependenciesProviderFactory.AssertExpectations(t)

	mockToolsProviderFactory := &servermocks.MockToolsProviderFactory[struct{}]{}
	defer mockToolsProviderFactory.AssertExpectations(t)

	mockParametersFactory := &servermocks.MockParametersFactory{}
	defer mockParametersFactory.AssertExpectations(t)

	mockFeaturesFactory := &servermocks.MockFeaturesFactory{}
	defer mockFeaturesFactory.AssertExpectations(t)

	mockApplicationFactory := &servermocks.MockApplicationFactory{}
	defer mockApplicationFactory.AssertExpectations(t)

	mockApplication := &adaptormocks.MockApplication{}
	defer mockApplication.AssertExpectations(t)

	mockModeSelector := &adaptormocks.MockModeSelector{}
	defer mockModeSelector.AssertExpectations(t)

	mockErrorWriter := &entitiesmocks.MockWriter{}
	defer mockErrorWriter.AssertExpectations(t)

	mockTool := &internaltoolsmocks.MockTool{}
	defer mockTool.AssertExpectations(t)

	ctx := t.Context()
	expectedName := "test-server"
	expectedTitle := "Test Server"
	expectedInstructions := "Test instructions"
	expectedDependencies := &struct{ Value string }{Value: "test"}
	expectedDepsErr := assert.AnError
	expectedTools := []tools.Tool{mockTool}

	expectedDepsProvider := definition.DependenciesProvider(func(resources definition.DependenciesProviderResources) (any, error) {
		return expectedDependencies, expectedDepsErr
	})
	expectedToolsProvider := definition.ToolsProvider(func(resources definition.ToolsProviderResources) []tools.Tool {
		return expectedTools
	})

	serverDefinition := server.Definition[struct{}]{
		Name:         expectedName,
		Title:        expectedTitle,
		Instructions: expectedInstructions,
	}

	mockFeaturesFactory.EXPECT().
		New(publictypes.Features{}).
		Return(definition.Features{}).
		Once()

	mockParametersFactory.EXPECT().
		New([]publictypes.Parameter(nil)).
		Return(nil).
		Once()

	mockDependenciesProviderFactory.EXPECT().
		New(mock.Anything).
		Return(expectedDepsProvider).
		Once()

	mockToolsProviderFactory.EXPECT().
		New(mock.Anything).
		Return(expectedToolsProvider).
		Once()

	mockApplicationFactory.EXPECT().
		New(mock.MatchedBy(func(def adaptor.ApplicationDefinition) bool {
			deps, depsErr := def.Dependencies(definition.DependenciesProviderResources{})
			if deps != expectedDependencies || depsErr != expectedDepsErr {
				return false
			}

			resultTools := def.Tools(definition.ToolsProviderResources{})
			if len(resultTools) != len(expectedTools) {
				return false
			}
			for i, tool := range resultTools {
				if tool != expectedTools[i] {
					return false
				}
			}

			return true
		})).
		Return(mockApplication).
		Once()

	mockApplication.EXPECT().
		ModeSelector().
		Return(mockModeSelector).
		Once()

	mockModeSelector.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(nil).
		Once()

	s := server.New(serverDefinition, mockFeaturesFactory, mockParametersFactory, mockDependenciesProviderFactory, mockToolsProviderFactory, mockApplicationFactory, mockErrorWriter)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 0, exitCode)
}

func TestServer_StartAndWaitForCompletion_KnownError(t *testing.T) {
	// Arrange
	mockDependenciesProviderFactory := &servermocks.MockDependenciesProviderFactory[struct{}]{}
	defer mockDependenciesProviderFactory.AssertExpectations(t)

	mockToolsProviderFactory := &servermocks.MockToolsProviderFactory[struct{}]{}
	defer mockToolsProviderFactory.AssertExpectations(t)

	mockParametersFactory := &servermocks.MockParametersFactory{}
	defer mockParametersFactory.AssertExpectations(t)

	mockFeaturesFactory := &servermocks.MockFeaturesFactory{}
	defer mockFeaturesFactory.AssertExpectations(t)

	mockApplicationFactory := &servermocks.MockApplicationFactory{}
	defer mockApplicationFactory.AssertExpectations(t)

	mockApplication := &adaptormocks.MockApplication{}
	defer mockApplication.AssertExpectations(t)

	mockModeSelector := &adaptormocks.MockModeSelector{}
	defer mockModeSelector.AssertExpectations(t)

	mockMessageCatalog := &adaptormocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockErrorWriter := &entitiesmocks.MockWriter{}
	defer mockErrorWriter.AssertExpectations(t)

	ctx := t.Context()
	expectedName := "test-server"
	expectedTitle := "Test Server"
	expectedInstructions := "Test instructions"
	expectedError := messages.AnError
	expectedErrorMessage := "A known error occurred"
	expectedFeatures := publictypes.Features{
		MATLAB: publictypes.MATLABFeature{Enabled: true},
	}

	serverDefinition := server.Definition[struct{}]{
		Name:         expectedName,
		Title:        expectedTitle,
		Instructions: expectedInstructions,
		Features:     expectedFeatures,
	}

	expectedInternalFeatures := definition.Features{
		MATLAB: definition.MATLABFeature{Enabled: true},
	}
	expectedDefinition := definition.New(
		expectedName,
		expectedTitle,
		expectedInstructions,
		expectedInternalFeatures,
		nil,
		nil,
		nil,
	)

	mockFeaturesFactory.EXPECT().
		New(expectedFeatures).
		Return(expectedInternalFeatures).
		Once()

	mockParametersFactory.EXPECT().
		New([]publictypes.Parameter(nil)).
		Return(nil).
		Once()

	mockDependenciesProviderFactory.EXPECT().
		New(mock.Anything).
		Return(nil).
		Once()

	mockToolsProviderFactory.EXPECT().
		New(mock.Anything).
		Return(nil).
		Once()

	mockApplicationFactory.EXPECT().
		New(expectedDefinition).
		Return(mockApplication).
		Once()

	mockApplication.EXPECT().
		ModeSelector().
		Return(mockModeSelector).
		Once()

	mockModeSelector.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(expectedError).
		Once()

	mockApplication.EXPECT().
		MessageCatalog().
		Return(mockMessageCatalog).
		Once()

	mockMessageCatalog.EXPECT().
		GetFromError(expectedError).
		Return(expectedErrorMessage).
		Once()

	mockErrorWriter.EXPECT().
		Write([]byte(expectedErrorMessage+"\n")).
		Return(len(expectedErrorMessage)+1, nil).
		Once()

	s := server.New(serverDefinition, mockFeaturesFactory, mockParametersFactory, mockDependenciesProviderFactory, mockToolsProviderFactory, mockApplicationFactory, mockErrorWriter)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 1, exitCode)
}
