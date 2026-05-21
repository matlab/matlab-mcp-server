// Copyright 2026 The MathWorks, Inc.

package selector_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/defaultparameters/selector"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	selectormocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/parameter/defaultparameters/selector"
)

func TestSelector_DefaultParameters_DescriptionsResolvedForVisibleParameters(t *testing.T) {
	// Arrange
	mockAppDef := &selectormocks.MockApplicationDefinition{}
	defer mockAppDef.AssertExpectations(t)

	mockMessageCatalog := &selectormocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedDescriptions := map[messages.MessageKey]struct {
		description string
	}{
		messages.CLIMessages_HelpDescription: {
			description: "Help description",
		},
		messages.CLIMessages_VersionDescription: {
			description: "Version description",
		},
		messages.CLIMessages_SetupMATLABDescription: {
			description: "Install MATLAB Add-On description",
		},
		messages.CLIMessages_DisableTelemetryDescription: {
			description: "Disable telemetry description",
		},
		messages.CLIMessages_BaseDirDescription: {
			description: "Base dir description",
		},
		messages.CLIMessages_LogLevelDescription: {
			description: "Log level description",
		},
		messages.CLIMessages_PreferredLocalMATLABRootDescription: {
			description: "MATLAB root description",
		},
		messages.CLIMessages_PreferredMATLABStartingDirectoryDescription: {
			description: "MATLAB starting directory description",
		},
		messages.CLIMessages_InitializeMATLABOnStartupDescription: {
			description: "Initialize MATLAB on startup description",
		},
		messages.CLIMessages_DisplayModeDescription: {
			description: "Display mode description",
		},
		messages.CLIMessages_MATLABSessionModeDescription: {
			description: "MATLAB session mode description",
		},
		messages.CLIMessages_ExtensionFileDescription: {
			description: "Extension file description",
		},
	}

	mockAppDef.EXPECT().
		Features().
		Return(definition.Features{
			MATLAB: definition.MATLABFeature{
				Enabled: true,
			},
		}).
		Once()

	for key, expected := range expectedDescriptions {
		mockMessageCatalog.EXPECT().
			Get(key).
			Return(expected.description).
			Once()
	}

	sut := selector.New(mockAppDef, mockMessageCatalog)

	// Act
	parameters := sut.DefaultParameters()

	// Assert
	expectedNumParametersWithNonEmptyDescription := len(expectedDescriptions)
	numParametersWithNonEmptyDescription := 0
	for _, p := range parameters {
		if !p.GetHiddenFlag() {
			numParametersWithNonEmptyDescription += 1
		}
	}

	assert.Equal(t,
		expectedNumParametersWithNonEmptyDescription,
		numParametersWithNonEmptyDescription,
	)
}

func TestSelector_DefaultParameters_MATLABEnabled(t *testing.T) {
	// Arrange
	mockAppDef := &selectormocks.MockApplicationDefinition{}
	defer mockAppDef.AssertExpectations(t)

	mockMessageCatalog := &selectormocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockAppDef.EXPECT().
		Features().
		Return(definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}).
		Once()

	mockMessageCatalog.EXPECT().
		Get(mock.Anything).
		Return("description")

	sut := selector.New(mockAppDef, mockMessageCatalog)

	// Act
	parameters := sut.DefaultParameters()

	// Assert
	assert.Len(t, parameters, 23)

	for _, p := range parameters {
		assert.True(t, p.GetActive(), "parameter %s should be active", p.GetID())
	}
}

func TestSelector_DefaultParameters_MATLABDisabled(t *testing.T) {
	// Arrange
	mockAppDef := &selectormocks.MockApplicationDefinition{}
	defer mockAppDef.AssertExpectations(t)

	mockMessageCatalog := &selectormocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedActiveStateByParameterID := map[string]bool{
		"HelpMode":                           true,
		"VersionMode":                        true,
		"SetupMATLABMode":                    true,
		"DisableTelemetry":                   true,
		"BaseDir":                            true,
		"LogLevel":                           true,
		"DuplicateLogsToStderr":              true,
		"WatchdogMode":                       true,
		"ServerInstanceID":                   true,
		"TelemetryCollectorEndpoint":         true,
		"TelemetryCollectionInterval":        true,
		"TelemetryCollectorEndpointInsecure": true,
		"PreferredLocalMATLABRoot":           false,
		"PreferredMATLABStartingDirectory":   false,
		"UseSingleMATLABSession":             false,
		"InitializeMATLABOnStartup":          false,
		"MATLABDisplayMode":                  false,
		"MATLABSessionMode":                  false,
		"MATLABSessionConnectionDetails":     false,
		"MATLABSessionConnectionTimeout":     false,
		"MATLABSessionDiscoveryTimeout":      false,
		"EmbeddedConnectorDetailsTimeout":    false,
		"ExtensionFile":                      false,
	}

	mockAppDef.EXPECT().
		Features().
		Return(definition.Features{MATLAB: definition.MATLABFeature{Enabled: false}}).
		Once()

	mockMessageCatalog.EXPECT().
		Get(mock.Anything).
		Return("description")

	sut := selector.New(mockAppDef, mockMessageCatalog)

	// Act
	parameters := sut.DefaultParameters()

	// Assert
	assert.Len(t, parameters, 23)

	for _, p := range parameters {
		expectedState, exists := expectedActiveStateByParameterID[p.GetID()]
		assert.True(t, exists, "unexpected parameter %s", p.GetID())
		assert.Equal(t, expectedState, p.GetActive(), "unexpected active state for parameter %s", p.GetID())
	}
}
