// Copyright 2025-2026 The MathWorks, Inc.

package config

import (
	"sync"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type Parser interface {
	Parse(args []string) ([]entities.Parameter, map[string]any, []string, messages.Error)
}

type OSLayer interface {
	Args() []string
}

type BuildInfo interface {
	FullVersion() string
}

type GenericConfig interface {
	Get(key string) (any, messages.Error)
}

type Config interface {
	GenericConfig

	Version() string
	HelpMode() bool
	VersionMode() bool
	WatchdogMode() bool
	SetupMATLABMode() bool

	BaseDir() string
	ServerInstanceID() string

	// Logger
	LogLevel() entities.LogLevel
	DuplicateLogsToStderr() bool
	RecordToLogger(logger entities.Logger)

	// MATLAB
	UseSingleMATLABSession() bool
	InitializeMATLABOnStartup() bool
	PreferredLocalMATLABRoot() string
	PreferredMATLABStartingDirectory() string
	ShouldShowMATLABDesktop() bool
	MATLABSessionMode() entities.MATLABSessionMode
	MATLABSessionConnectionDetails() string
	MATLABSessionConnectionTimeout() time.Duration
	MATLABSessionDiscoveryTimeout() time.Duration
	EmbeddedConnectorDetailsTimeout() time.Duration
	ExtensionFile() string

	// Telemetry
	DisableTelemetry() bool
	TelemetryCollectorEndpoint() string
	TelemetryCollectionInterval() time.Duration
	TelemetryCollectorEndpointInsecure() bool
	SpecifiedParameters() []string
	AsPIISafeJSONString() string
}

type Factory struct {
	parser    Parser
	osLayer   OSLayer
	buildInfo BuildInfo

	initOnce       sync.Once
	initError      messages.Error
	configInstance *config
}

func NewFactory(parser Parser, osLayer OSLayer, buildInfo BuildInfo) *Factory {
	return &Factory{
		parser:    parser,
		osLayer:   osLayer,
		buildInfo: buildInfo,
	}
}

func (f *Factory) Config() (Config, messages.Error) {
	f.initOnce.Do(func() {
		configInstance, err := newConfig(f.osLayer, f.parser, f.buildInfo)
		if err != nil {
			f.initError = err
			return
		}

		f.configInstance = configInstance
	})

	if f.initError != nil {
		return nil, f.initError
	}

	return f.configInstance, nil
}
