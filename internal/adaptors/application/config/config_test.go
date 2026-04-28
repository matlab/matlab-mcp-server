// Copyright 2025-2026 The MathWorks, Inc.

package config_test

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultParameters() []entities.Parameter {
	return []entities.Parameter{
		defaultparameters.HelpMode(),
		defaultparameters.VersionMode(),
		defaultparameters.WatchdogMode(),
		defaultparameters.SetupMATLABMode(),

		defaultparameters.BaseDir(),
		defaultparameters.ServerInstanceID(),

		defaultparameters.LogLevel(),
		defaultparameters.DuplicateLogsToStderr(),

		defaultparameters.UseSingleMATLABSession(),
		defaultparameters.PreferredLocalMATLABRoot(),
		defaultparameters.PreferredMATLABStartingDirectory(),
		defaultparameters.InitializeMATLABOnStartup(),
		defaultparameters.MATLABDisplayMode(),
		defaultparameters.MATLABSessionMode(),
		defaultparameters.MATLABSessionConnectionDetails(),
		defaultparameters.MATLABSessionConnectionTimeout(),
		defaultparameters.MATLABSessionDiscoveryTimeout(),
		defaultparameters.EmbeddedConnectorDetailsTimeout(),

		defaultparameters.DisableTelemetry(),
		defaultparameters.ExtensionFile(),
		defaultparameters.TelemetryCollectorEndpoint(),
		defaultparameters.TelemetryCollectionInterval(),
		defaultparameters.TelemetryCollectorEndpointInsecure(),
	}
}

func configDefaultParsedArgs() map[string]any {
	result := make(map[string]any)
	for _, p := range defaultParameters() {
		result[p.GetID()] = p.GetDefaultValue()
	}
	return result
}

func TestNewConfig_InvalidParameterType(t *testing.T) {
	testCases := []struct {
		key          string
		invalidValue any
		expectedType string
	}{
		{key: defaultparameters.VersionMode().GetID(), invalidValue: "false", expectedType: "bool"},
		{key: defaultparameters.HelpMode().GetID(), invalidValue: "false", expectedType: "bool"},
		{key: defaultparameters.WatchdogMode().GetID(), invalidValue: "false", expectedType: "bool"},
		{key: defaultparameters.SetupMATLABMode().GetID(), invalidValue: "false", expectedType: "bool"},

		{key: defaultparameters.BaseDir().GetID(), invalidValue: 123, expectedType: "string"},
		{key: defaultparameters.ServerInstanceID().GetID(), invalidValue: 123, expectedType: "string"},

		{key: defaultparameters.LogLevel().GetID(), invalidValue: 123, expectedType: "string"},
		{key: defaultparameters.DuplicateLogsToStderr().GetID(), invalidValue: "false", expectedType: "bool"},

		{key: defaultparameters.UseSingleMATLABSession().GetID(), invalidValue: "true", expectedType: "bool"},
		{key: defaultparameters.InitializeMATLABOnStartup().GetID(), invalidValue: "false", expectedType: "bool"},
		{key: defaultparameters.PreferredLocalMATLABRoot().GetID(), invalidValue: 123, expectedType: "string"},
		{key: defaultparameters.PreferredMATLABStartingDirectory().GetID(), invalidValue: 123, expectedType: "string"},
		{key: defaultparameters.MATLABDisplayMode().GetID(), invalidValue: 123, expectedType: "string"},
		{key: defaultparameters.MATLABSessionMode().GetID(), invalidValue: 123, expectedType: "string"},
		{key: defaultparameters.MATLABSessionConnectionDetails().GetID(), invalidValue: 123, expectedType: "string"},
		{key: defaultparameters.MATLABSessionConnectionTimeout().GetID(), invalidValue: "5s", expectedType: "time.Duration"},
		{key: defaultparameters.MATLABSessionDiscoveryTimeout().GetID(), invalidValue: "30s", expectedType: "time.Duration"},
		{key: defaultparameters.EmbeddedConnectorDetailsTimeout().GetID(), invalidValue: "1m", expectedType: "time.Duration"},
		{key: defaultparameters.ExtensionFile().GetID(), invalidValue: 123, expectedType: "string"},

		{key: defaultparameters.DisableTelemetry().GetID(), invalidValue: "false", expectedType: "bool"},
		{key: defaultparameters.TelemetryCollectorEndpoint().GetID(), invalidValue: 123, expectedType: "string"},
		{key: defaultparameters.TelemetryCollectionInterval().GetID(), invalidValue: "1m", expectedType: "time.Duration"},
		{key: defaultparameters.TelemetryCollectorEndpointInsecure().GetID(), invalidValue: "false", expectedType: "bool"},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			mockBuildInfo := &configmocks.MockBuildInfo{}
			defer mockBuildInfo.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			parsedArgs := configDefaultParsedArgs()
			parsedArgs[tc.key] = tc.invalidValue

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
				Once()

			expectedError := messages.New_StartupErrors_InvalidParameterType_Error(tc.key, tc.expectedType)

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

			// Assert
			require.Equal(t, expectedError, err)
			assert.Nil(t, cfg)
		})
	}
}

func TestNewConfig_MissingParameter(t *testing.T) {
	parameters := []entities.Parameter{
		defaultparameters.VersionMode(),
		defaultparameters.HelpMode(),
		defaultparameters.WatchdogMode(),
		defaultparameters.SetupMATLABMode(),
		defaultparameters.BaseDir(),
		defaultparameters.ServerInstanceID(),
		defaultparameters.LogLevel(),
		defaultparameters.DuplicateLogsToStderr(),
		defaultparameters.UseSingleMATLABSession(),
		defaultparameters.InitializeMATLABOnStartup(),
		defaultparameters.PreferredLocalMATLABRoot(),
		defaultparameters.PreferredMATLABStartingDirectory(),
		defaultparameters.MATLABDisplayMode(),
		defaultparameters.MATLABSessionMode(),
		defaultparameters.MATLABSessionConnectionDetails(),
		defaultparameters.MATLABSessionConnectionTimeout(),
		defaultparameters.MATLABSessionDiscoveryTimeout(),
		defaultparameters.EmbeddedConnectorDetailsTimeout(),
		defaultparameters.ExtensionFile(),
		defaultparameters.DisableTelemetry(),
		defaultparameters.TelemetryCollectorEndpoint(),
		defaultparameters.TelemetryCollectionInterval(),
		defaultparameters.TelemetryCollectorEndpointInsecure(),
	}

	for _, parameter := range parameters {
		t.Run(parameter.GetID(), func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			mockBuildInfo := &configmocks.MockBuildInfo{}
			defer mockBuildInfo.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			parsedArgs := configDefaultParsedArgs()
			delete(parsedArgs, parameter.GetID())

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
				Once()

			expectedError := messages.New_StartupErrors_InvalidParameterKey_Error(parameter.GetID())

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

			// Assert
			require.Equal(t, expectedError, err)
			assert.Nil(t, cfg)
		})
	}
}

func TestNewConfig_ParseError(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(nil, nil, nil, messages.AnError).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	assert.Nil(t, cfg, "Config should be nil")
}

func TestNewConfig_InvalidLogLevel(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}
	invalidLevel := "invalid-level"

	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.LogLevel().GetID()] = invalidLevel

	expectedError := messages.New_StartupErrors_InvalidLogLevel_Error(invalidLevel)

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

	// Assert
	require.Equal(t, expectedError, err)
	assert.Nil(t, cfg, "Config should be nil")
}

func TestNewConfig_InvalidDisplayMode(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}
	invalidMode := "invalid-mode"

	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.MATLABDisplayMode().GetID()] = invalidMode

	expectedError := messages.New_StartupErrors_InvalidDisplayMode_Error(invalidMode)

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

	// Assert
	require.Equal(t, expectedError, err)
	assert.Nil(t, cfg, "Config should be nil")
}

func TestNewConfig_InvalidMATLABSessionMode(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}
	invalidMode := "invalid-mode"

	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.MATLABSessionMode().GetID()] = invalidMode

	expectedError := messages.New_StartupErrors_InvalidMATLABSessionMode_Error(invalidMode)

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

	// Assert
	require.Equal(t, expectedError, err)
	assert.Nil(t, cfg, "Config should be nil")
}

func TestConfig_Version_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	expectedVersion := "github.com/matlab/matlab-mcp-core-server v1.2.3"

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return([]entities.Parameter{}, configDefaultParsedArgs(), []string{}, nil).
		Once()

	mockBuildInfo.EXPECT().
		FullVersion().
		Return(expectedVersion).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)
	require.NoError(t, err)

	version := cfg.Version()

	// Assert
	require.Equal(t, expectedVersion, version)
}

func TestConfig_SpecifiedParameters_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}
	expectedSpecifiedParameters := []string{"DisableTelemetry", "LogLevel"}

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(defaultParameters(), configDefaultParsedArgs(), expectedSpecifiedParameters, nil).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)
	require.NoError(t, err)

	result := cfg.SpecifiedParameters()

	// Assert
	require.Equal(t, expectedSpecifiedParameters, result)
}

func TestConfig_InitializeMATLABOnStartup_DisabledWhenNotSingleSession(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.UseSingleMATLABSession().GetID()] = false
	parsedArgs[defaultparameters.InitializeMATLABOnStartup().GetID()] = true

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

	// Assert
	require.NoError(t, err)
	assert.False(t, cfg.InitializeMATLABOnStartup(), "InitializeMATLABOnStartup should be false when UseSingleMATLABSession is false")
}

func TestConfig_ShouldShowMATLABDesktop_DefaultsToNoDesktopInInstallAddOnMode(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.SetupMATLABMode().GetID()] = true

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

	// Assert
	require.NoError(t, err)
	assert.False(t, cfg.ShouldShowMATLABDesktop(), "ShouldShowMATLABDesktop should default to false in install add-on mode")
}

func TestNewConfig_MATLABSessionConnectionTimeout_FallsBackToDefaultWhenNotPositive(t *testing.T) {
	testCases := []struct {
		name    string
		timeout time.Duration
	}{
		{name: "zero timeout", timeout: 0},
		{name: "negative timeout", timeout: -time.Minute},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			mockBuildInfo := &configmocks.MockBuildInfo{}
			defer mockBuildInfo.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			parsedArgs := configDefaultParsedArgs()
			parsedArgs[defaultparameters.MATLABSessionConnectionTimeout().GetID()] = tc.timeout
			expectedTimeout := defaultparameters.MATLABSessionConnectionTimeout().GetTypedDefaultValue()

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
				Once()

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTimeout, cfg.MATLABSessionConnectionTimeout())
		})
	}
}

func TestNewConfig_MATLABSessionDiscoveryTimeout_FallsBackToDefaultWhenNotPositive(t *testing.T) {
	testCases := []struct {
		name    string
		timeout time.Duration
	}{
		{name: "zero timeout", timeout: 0},
		{name: "negative timeout", timeout: -time.Minute},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			mockBuildInfo := &configmocks.MockBuildInfo{}
			defer mockBuildInfo.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			parsedArgs := configDefaultParsedArgs()
			parsedArgs[defaultparameters.MATLABSessionDiscoveryTimeout().GetID()] = tc.timeout
			expectedTimeout := defaultparameters.MATLABSessionDiscoveryTimeout().GetTypedDefaultValue()

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
				Once()

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTimeout, cfg.MATLABSessionDiscoveryTimeout())
		})
	}
}

func TestNewConfig_EmbeddedConnectorDetailsTimeout_FallsBackToDefaultWhenNotPositive(t *testing.T) {
	testCases := []struct {
		name    string
		timeout time.Duration
	}{
		{name: "zero timeout", timeout: 0},
		{name: "negative timeout", timeout: -time.Minute},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			mockBuildInfo := &configmocks.MockBuildInfo{}
			defer mockBuildInfo.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			parsedArgs := configDefaultParsedArgs()
			parsedArgs[defaultparameters.EmbeddedConnectorDetailsTimeout().GetID()] = tc.timeout
			expectedTimeout := defaultparameters.EmbeddedConnectorDetailsTimeout().GetTypedDefaultValue()

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
				Once()

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTimeout, cfg.EmbeddedConnectorDetailsTimeout())
		})
	}
}

func TestNewConfig_TelemetryCollectionInterval_FallsBackToDefaultWhenNotPositive(t *testing.T) {
	testCases := []struct {
		name     string
		interval time.Duration
	}{
		{name: "zero interval", interval: 0},
		{name: "negative interval", interval: -time.Minute},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			mockBuildInfo := &configmocks.MockBuildInfo{}
			defer mockBuildInfo.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			parsedArgs := configDefaultParsedArgs()
			parsedArgs[defaultparameters.TelemetryCollectionInterval().GetID()] = tc.interval
			expectedInterval := defaultparameters.TelemetryCollectionInterval().GetTypedDefaultValue()

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
				Once()

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedInterval, cfg.TelemetryCollectionInterval())
		})
	}
}

func TestConfig_RecordToLogger_HappyPath(t *testing.T) {
	// Arrange
	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.DisableTelemetry().GetID()] = true
	parsedArgs[defaultparameters.PreferredMATLABStartingDirectory().GetID()] = filepath.Join("home", "user")
	parsedArgs[defaultparameters.LogLevel().GetID()] = string(entities.LogLevelDebug)
	parsedArgs[defaultparameters.PreferredLocalMATLABRoot().GetID()] = filepath.Join("home", "matlab")
	parsedArgs[defaultparameters.UseSingleMATLABSession().GetID()] = false

	expectedLogMessage := "Configuration state"
	expectedConfigField := map[string]any{
		defaultparameters.DisableTelemetry().GetID():                 true,
		defaultparameters.PreferredMATLABStartingDirectory().GetID(): filepath.Join("home", "user"),
		defaultparameters.LogLevel().GetID():                         string(entities.LogLevelDebug),
		defaultparameters.PreferredLocalMATLABRoot().GetID():         filepath.Join("home", "matlab"),
		defaultparameters.UseSingleMATLABSession().GetID():           false,
		defaultparameters.InitializeMATLABOnStartup().GetID():        false,
		defaultparameters.MATLABSessionMode().GetID():                string(entities.MATLABSessionModeNew),
		defaultparameters.MATLABSessionConnectionTimeout().GetID():   5 * time.Second,
		defaultparameters.MATLABSessionDiscoveryTimeout().GetID():    30 * time.Second,
		defaultparameters.DuplicateLogsToStderr().GetID():            false,
	}

	parameters := []entities.Parameter{
		defaultparameters.DisableTelemetry(),
		defaultparameters.UseSingleMATLABSession(),
		defaultparameters.LogLevel(),
		defaultparameters.DuplicateLogsToStderr(),
		defaultparameters.PreferredLocalMATLABRoot(),
		defaultparameters.PreferredMATLABStartingDirectory(),
		defaultparameters.InitializeMATLABOnStartup(),
		defaultparameters.MATLABSessionMode(),
		defaultparameters.MATLABSessionConnectionTimeout(),
		defaultparameters.MATLABSessionDiscoveryTimeout(),
		defaultparameters.HelpMode(),
		defaultparameters.VersionMode(),
		defaultparameters.BaseDir(),
		defaultparameters.WatchdogMode(),
		defaultparameters.SetupMATLABMode(),
		defaultparameters.ServerInstanceID(),
	}

	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(parameters, parsedArgs, []string{}, nil)

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)
	require.NoError(t, err)

	testLogger := testutils.NewInspectableLogger()

	// Act
	cfg.RecordToLogger(testLogger)

	// Assert
	infoLogs := testLogger.InfoLogs()
	require.Len(t, infoLogs, 1)

	fields, found := infoLogs[expectedLogMessage]
	require.True(t, found, "Expected log message not found")

	for expectedField, expectedValue := range expectedConfigField {
		actualValue, exists := fields[expectedField]
		require.True(t, exists, "%s field not found in log", expectedField)
		assert.Equal(t, expectedValue, actualValue, "%s field has incorrect value", expectedField)
	}
}

func TestNewConfig_ExistingSessionMode_DisallowedParameter(t *testing.T) {
	disallowedParameters := []entities.Parameter{
		defaultparameters.PreferredLocalMATLABRoot(),
		defaultparameters.PreferredMATLABStartingDirectory(),
		defaultparameters.MATLABDisplayMode(),
	}

	for _, param := range disallowedParameters {
		t.Run(param.GetID(), func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			mockBuildInfo := &configmocks.MockBuildInfo{}
			defer mockBuildInfo.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			parsedArgs := configDefaultParsedArgs()
			parsedArgs[defaultparameters.MATLABSessionMode().GetID()] = string(entities.MATLABSessionModeExisting)

			specifiedParameters := []string{param.GetID()}

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, parsedArgs, specifiedParameters, nil).
				Once()

			expectedError := messages.New_StartupErrors_ArgumentNotAllowedInSessionMode_Error(
				param.GetFlagName(),
				string(entities.MATLABSessionModeExisting),
			)

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

			// Assert
			require.Equal(t, expectedError, err)
			assert.Nil(t, cfg)
		})
	}
}

func TestNewConfig_ExistingSessionMode_AllowedWithoutSpecifiedParameters(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.MATLABSessionMode().GetID()] = string(entities.MATLABSessionModeExisting)

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return([]entities.Parameter{}, parsedArgs, []string{}, nil).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestConfig_AsPIISafeJSONString_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockBuildInfo := &configmocks.MockBuildInfo{}
	defer mockBuildInfo.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	parameters := defaultParameters()
	parsedArgs := configDefaultParsedArgs()

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(parameters, parsedArgs, []string{}, nil).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser, mockBuildInfo)
	require.NoError(t, err)

	result := cfg.AsPIISafeJSONString()

	// Assert
	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(result), &parsed))

	piiSafeParams := []entities.Parameter{
		defaultparameters.HelpMode(),
		defaultparameters.VersionMode(),
		defaultparameters.SetupMATLABMode(),
		defaultparameters.DisableTelemetry(),
		defaultparameters.UseSingleMATLABSession(),
		defaultparameters.LogLevel(),
		defaultparameters.InitializeMATLABOnStartup(),
		defaultparameters.MATLABDisplayMode(),
		defaultparameters.MATLABSessionMode(),
		defaultparameters.MATLABSessionConnectionTimeout(),
		defaultparameters.MATLABSessionDiscoveryTimeout(),
		defaultparameters.WatchdogMode(),
		defaultparameters.TelemetryCollectionInterval(),
		defaultparameters.TelemetryCollectorEndpointInsecure(),
		defaultparameters.DuplicateLogsToStderr(),
	}
	for _, param := range piiSafeParams {
		var expected any
		raw, _ := json.Marshal(parsedArgs[param.GetID()])
		_ = json.Unmarshal(raw, &expected)
		assert.Equal(t, expected, parsed[param.GetID()], "%s should show actual value", param.GetID())
	}

	redactedParams := []entities.Parameter{
		defaultparameters.PreferredLocalMATLABRoot(),
		defaultparameters.PreferredMATLABStartingDirectory(),
		defaultparameters.BaseDir(),
		defaultparameters.ServerInstanceID(),
		defaultparameters.MATLABSessionConnectionDetails(),
		defaultparameters.TelemetryCollectorEndpoint(),
	}
	for _, param := range redactedParams {
		assert.Equal(t, config.RedactedValue, parsed[param.GetID()], "%s should be redacted", param.GetID())
	}
}
