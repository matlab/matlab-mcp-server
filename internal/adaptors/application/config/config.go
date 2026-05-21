// Copyright 2025-2026 The MathWorks, Inc.

package config

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

const redactedValue = "[REDACTED]"

type validatedArguments struct {
	versionMode     bool
	helpMode        bool
	watchdogMode    bool
	setupMATLABMode bool

	baseDirectory    string
	serverInstanceID string

	// Logger
	logLevel              entities.LogLevel
	duplicateLogsToStderr bool

	// MATLAB
	useSingleMATLABSession           bool
	initializeMATLABOnStartup        bool
	preferredLocalMATLABRoot         string
	preferredMATLABStartingDirectory string
	displayMode                      entities.DisplayMode
	matlabSessionMode                entities.MATLABSessionMode
	matlabSessionConnectionDetails   string
	matlabSessionConnectionTimeout   time.Duration
	matlabSessionDiscoveryTimeout    time.Duration
	embeddedConnectorDetailsTimeout  time.Duration
	extensionFile                    string

	// Telemetry
	disableTelemetry                   bool
	telemetryCollectorEndpoint         string
	telemetryCollectionInterval        time.Duration
	telemetryCollectorEndpointInsecure bool
}

type rawConfig struct {
	parameters          []entities.Parameter
	parsedArgs          map[string]any
	specifiedParameters []string
}

func (c *rawConfig) Get(key string) (any, messages.Error) {
	return getForKey(c.parsedArgs, key)
}

type config struct {
	buildInfo BuildInfo

	*rawConfig
	validatedArguments
}

func newConfig(osLayer OSLayer, parser Parser, buildInfo BuildInfo) (*config, messages.Error) {
	parameters, parsedArgs, specifiedParameters, err := parser.Parse(osLayer.Args()[1:])
	if err != nil {
		return nil, err
	}

	rawCfg := &rawConfig{
		parameters:          parameters,
		parsedArgs:          parsedArgs,
		specifiedParameters: specifiedParameters,
	}

	validated, err := validateArguments(rawCfg)
	if err != nil {
		return nil, err
	}

	return &config{
		buildInfo:          buildInfo,
		rawConfig:          rawCfg,
		validatedArguments: validated,
	}, nil
}

func (c *config) Get(key string) (any, messages.Error) {
	return getForKey(c.parsedArgs, key)
}

func (c *config) Version() string {
	return c.buildInfo.FullVersion()
}

func (c *config) LogLevel() entities.LogLevel {
	return c.logLevel
}

func (c *config) DuplicateLogsToStderr() bool {
	return c.duplicateLogsToStderr
}

func (c *config) VersionMode() bool {
	return c.versionMode
}

func (c *config) HelpMode() bool {
	return c.helpMode
}

func (c *config) WatchdogMode() bool {
	return c.watchdogMode
}

func (c *config) SetupMATLABMode() bool {
	return c.setupMATLABMode
}

func (c *config) UseSingleMATLABSession() bool {
	return c.useSingleMATLABSession
}

func (c *config) PreferredLocalMATLABRoot() string {
	return c.preferredLocalMATLABRoot
}

func (c *config) PreferredMATLABStartingDirectory() string {
	return c.preferredMATLABStartingDirectory
}

func (c *config) InitializeMATLABOnStartup() bool {
	return c.initializeMATLABOnStartup
}

func (c *config) ShouldShowMATLABDesktop() bool {
	switch c.displayMode {
	case entities.DisplayModeDesktop:
		return true
	case entities.DisplayModeNoDesktop:
		return false
	default:
		return true
	}
}

func (c *config) MATLABSessionMode() entities.MATLABSessionMode {
	return c.matlabSessionMode
}

func (c *config) MATLABSessionConnectionDetails() string {
	return c.matlabSessionConnectionDetails
}

func (c *config) MATLABSessionConnectionTimeout() time.Duration {
	return c.matlabSessionConnectionTimeout
}

func (c *config) MATLABSessionDiscoveryTimeout() time.Duration {
	return c.matlabSessionDiscoveryTimeout
}

func (c *config) ExtensionFile() string {
	return c.extensionFile
}

func (c *config) BaseDir() string {
	return c.baseDirectory
}

func (c *config) ServerInstanceID() string {
	return c.serverInstanceID
}

func (c *config) EmbeddedConnectorDetailsTimeout() time.Duration {
	return c.embeddedConnectorDetailsTimeout
}

func (c *config) DisableTelemetry() bool {
	return c.disableTelemetry
}

func (c *config) TelemetryCollectorEndpoint() string {
	return c.telemetryCollectorEndpoint
}

func (c *config) TelemetryCollectionInterval() time.Duration {
	return c.telemetryCollectionInterval
}

func (c *config) TelemetryCollectorEndpointInsecure() bool {
	return c.telemetryCollectorEndpointInsecure
}

func (c *config) SpecifiedParameters() []string {
	return slices.Clone(c.specifiedParameters)
}

func (c *config) AsPIISafeJSONString() string {
	details := map[string]any{}

	for _, param := range c.parameters {
		var value any = redactedValue
		if param.GetPIISafe() {
			if actualValue, err := c.Get(param.GetID()); err == nil {
				value = actualValue
			}
		}
		details[param.GetID()] = value
	}

	jsonBytes, err := json.Marshal(details)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

func (c *config) RecordToLogger(logger entities.Logger) {
	for _, param := range c.parameters {
		if param.GetRecordToLog() {
			value, err := c.Get(param.GetID())
			if err == nil {
				logger = logger.With(param.GetID(), value)
			}
		}
	}
	logger.Info("Configuration state")
}

func validateArguments(rawCfg *rawConfig) (validatedArguments, messages.Error) {
	versionMode, err := get(rawCfg, defaultparameters.VersionMode())
	if err != nil {
		return validatedArguments{}, err
	}

	helpMode, err := get(rawCfg, defaultparameters.HelpMode())
	if err != nil {
		return validatedArguments{}, err
	}

	watchdogMode, err := get(rawCfg, defaultparameters.WatchdogMode())
	if err != nil {
		return validatedArguments{}, err
	}

	setupMATLABMode, err := get(rawCfg, defaultparameters.SetupMATLABMode())
	if err != nil {
		return validatedArguments{}, err
	}

	baseDirectory, err := get(rawCfg, defaultparameters.BaseDir())
	if err != nil {
		return validatedArguments{}, err
	}

	serverInstanceID, err := get(rawCfg, defaultparameters.ServerInstanceID())
	if err != nil {
		return validatedArguments{}, err
	}

	logLevel, err := get(rawCfg, defaultparameters.LogLevel())
	if err != nil {
		return validatedArguments{}, err
	}

	switch logLevel {
	case string(entities.LogLevelDebug), string(entities.LogLevelInfo), string(entities.LogLevelWarn), string(entities.LogLevelError):
	default:
		return validatedArguments{}, messages.New_StartupErrors_InvalidLogLevel_Error(logLevel)
	}

	duplicateLogsToStderr, err := get(rawCfg, defaultparameters.DuplicateLogsToStderr())
	if err != nil {
		return validatedArguments{}, err
	}

	useSingleMATLABSession, err := get(rawCfg, defaultparameters.UseSingleMATLABSession())
	if err != nil {
		return validatedArguments{}, err
	}

	initializeMATLABOnStartup, err := get(rawCfg, defaultparameters.InitializeMATLABOnStartup())
	if err != nil {
		return validatedArguments{}, err
	}

	if !useSingleMATLABSession {
		initializeMATLABOnStartup = false
	}

	preferredLocalMATLABRoot, err := get(rawCfg, defaultparameters.PreferredLocalMATLABRoot())
	if err != nil {
		return validatedArguments{}, err
	}

	preferredMATLABStartingDirectory, err := get(rawCfg, defaultparameters.PreferredMATLABStartingDirectory())
	if err != nil {
		return validatedArguments{}, err
	}

	displayMode, err := get(rawCfg, defaultparameters.MATLABDisplayMode())
	if err != nil {
		return validatedArguments{}, err
	}

	switch displayMode {
	case string(entities.DisplayModeDesktop), string(entities.DisplayModeNoDesktop):
		break
	default:
		return validatedArguments{}, messages.New_StartupErrors_InvalidDisplayMode_Error(displayMode)
	}

	extensionFile, err := get(rawCfg, defaultparameters.ExtensionFile())
	if err != nil {
		return validatedArguments{}, err
	}

	matlabSessionMode, err := get(rawCfg, defaultparameters.MATLABSessionMode())
	if err != nil {
		return validatedArguments{}, err
	}

	switch matlabSessionMode {
	case string(entities.MATLABSessionModeNew), string(entities.MATLABSessionModeExisting):
	default:
		return validatedArguments{}, messages.New_StartupErrors_InvalidMATLABSessionMode_Error(matlabSessionMode)
	}

	matlabSessionConnectionDetails, err := get(rawCfg, defaultparameters.MATLABSessionConnectionDetails())
	if err != nil {
		return validatedArguments{}, err
	}

	matlabSessionConnectionTimeout, err := get(rawCfg, defaultparameters.MATLABSessionConnectionTimeout())
	if err != nil {
		return validatedArguments{}, err
	}

	if matlabSessionConnectionTimeout <= 0 {
		matlabSessionConnectionTimeout = defaultparameters.MATLABSessionConnectionTimeout().GetTypedDefaultValue()
	}

	matlabSessionDiscoveryTimeout, err := get(rawCfg, defaultparameters.MATLABSessionDiscoveryTimeout())
	if err != nil {
		return validatedArguments{}, err
	}

	if matlabSessionDiscoveryTimeout <= 0 {
		matlabSessionDiscoveryTimeout = defaultparameters.MATLABSessionDiscoveryTimeout().GetTypedDefaultValue()
	}

	disableTelemetry, err := get(rawCfg, defaultparameters.DisableTelemetry())
	if err != nil {
		return validatedArguments{}, err
	}

	embeddedConnectorDetailsTimeout, err := get(rawCfg, defaultparameters.EmbeddedConnectorDetailsTimeout())
	if err != nil {
		return validatedArguments{}, err
	}

	if embeddedConnectorDetailsTimeout <= 0 {
		embeddedConnectorDetailsTimeout = defaultparameters.EmbeddedConnectorDetailsTimeout().GetTypedDefaultValue()
	}

	telemetryCollectorEndpoint, err := get(rawCfg, defaultparameters.TelemetryCollectorEndpoint())
	if err != nil {
		return validatedArguments{}, err
	}

	telemetryCollectionInterval, err := get(rawCfg, defaultparameters.TelemetryCollectionInterval())
	if err != nil {
		return validatedArguments{}, err
	}

	if telemetryCollectionInterval <= 0 {
		telemetryCollectionInterval = defaultparameters.TelemetryCollectionInterval().GetTypedDefaultValue()
	}

	telemetryCollectorEndpointInsecure, err := get(rawCfg, defaultparameters.TelemetryCollectorEndpointInsecure())
	if err != nil {
		return validatedArguments{}, err
	}

	args := validatedArguments{
		versionMode:     versionMode,
		helpMode:        helpMode,
		watchdogMode:    watchdogMode,
		setupMATLABMode: setupMATLABMode,

		baseDirectory:    baseDirectory,
		serverInstanceID: serverInstanceID,

		// Logger
		logLevel:              entities.LogLevel(logLevel),
		duplicateLogsToStderr: duplicateLogsToStderr,

		// MATLAB
		useSingleMATLABSession:           useSingleMATLABSession,
		initializeMATLABOnStartup:        initializeMATLABOnStartup,
		preferredLocalMATLABRoot:         preferredLocalMATLABRoot,
		preferredMATLABStartingDirectory: preferredMATLABStartingDirectory,
		displayMode:                      entities.DisplayMode(displayMode),
		matlabSessionMode:                entities.MATLABSessionMode(matlabSessionMode),
		matlabSessionConnectionDetails:   matlabSessionConnectionDetails,
		matlabSessionConnectionTimeout:   matlabSessionConnectionTimeout,
		matlabSessionDiscoveryTimeout:    matlabSessionDiscoveryTimeout,
		embeddedConnectorDetailsTimeout:  embeddedConnectorDetailsTimeout,
		extensionFile:                    extensionFile,

		// Telemetry
		disableTelemetry:                   disableTelemetry,
		telemetryCollectorEndpoint:         telemetryCollectorEndpoint,
		telemetryCollectionInterval:        telemetryCollectionInterval,
		telemetryCollectorEndpointInsecure: telemetryCollectorEndpointInsecure,
	}

	args, err = checkArgumentCompatibilityAndAdjustDefaults(args, rawCfg.specifiedParameters)
	if err != nil {
		return validatedArguments{}, err
	}

	return args, nil
}

func checkArgumentCompatibilityAndAdjustDefaults(args validatedArguments, specifiedParameters []string) (validatedArguments, messages.Error) {
	// If installing the MATLAB Add-On, and displayMode isn't specified
	// it's a better user experience to not flash the desktop
	if args.setupMATLABMode && !slices.Contains(specifiedParameters, defaultparameters.MATLABDisplayMode().GetID()) {
		args.displayMode = entities.DisplayModeNoDesktop
	}

	// If using MATLAB Session Mode `existing`, most of the MATLAB flags are unsupported
	if args.matlabSessionMode == entities.MATLABSessionModeExisting {
		disallowedParametersInExistingSessionMode := []entities.Parameter{
			defaultparameters.PreferredLocalMATLABRoot(),
			defaultparameters.PreferredMATLABStartingDirectory(),
			defaultparameters.MATLABDisplayMode(),
		}
		for _, parameter := range disallowedParametersInExistingSessionMode {
			if slices.Contains(specifiedParameters, parameter.GetID()) {
				return validatedArguments{}, messages.New_StartupErrors_ArgumentNotAllowedInSessionMode_Error(parameter.GetFlagName(), string(entities.MATLABSessionModeExisting))
			}
		}
	}

	return args, nil
}

func getForKey(args map[string]any, key string) (any, messages.Error) {
	if value, ok := args[key]; ok {
		return value, nil
	}
	return nil, messages.New_StartupErrors_InvalidParameterKey_Error(key)
}
