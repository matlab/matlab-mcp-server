// Copyright 2026 The MathWorks, Inc.

package userconfig

import (
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type userConfigEntry struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     any    `json:"default"`
}

type parameterForMCPB struct {
	parameter.ParameterWithDescriptionFromMessageCatalog
	Title        string
	TypeOverride string
}

func GetUserConfig() (map[string]userConfigEntry, error) {
	factory := newEntryFactory()

	parametersForMCPB := []parameterForMCPB{
		{
			ParameterWithDescriptionFromMessageCatalog: defaultparameters.PreferredLocalMATLABRoot(),
			Title:        "MATLAB Installation Path",
			TypeOverride: "directory",
		},
		{
			ParameterWithDescriptionFromMessageCatalog: defaultparameters.PreferredMATLABStartingDirectory(),
			Title:        "Initial Working Folder",
			TypeOverride: "directory",
		},
		{
			ParameterWithDescriptionFromMessageCatalog: defaultparameters.BaseDir(),
			Title:        "Log Folder",
			TypeOverride: "directory",
		},
		{
			ParameterWithDescriptionFromMessageCatalog: defaultparameters.InitializeMATLABOnStartup(),
			Title: "Initialize MATLAB on Startup",
		},
		{
			ParameterWithDescriptionFromMessageCatalog: defaultparameters.DisableTelemetry(),
			Title: "Disable Telemetry",
		},
		{
			ParameterWithDescriptionFromMessageCatalog: defaultparameters.MATLABDisplayMode(),
			Title: "MATLAB Display Mode",
		},
		{
			ParameterWithDescriptionFromMessageCatalog: defaultparameters.MATLABSessionMode(),
			Title: "MATLAB Session Mode",
		},
		{
			ParameterWithDescriptionFromMessageCatalog: defaultparameters.ExtensionFile(),
			Title:        "Extension File",
			TypeOverride: "file",
		},
	}

	config := make(map[string]userConfigEntry, len(parametersForMCPB))
	for _, p := range parametersForMCPB {
		var err error
		config[p.GetID()], err = factory.fromParameter(
			p.ParameterWithDescriptionFromMessageCatalog,
			p.Title,
			p.TypeOverride,
		)

		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

type entryFactory struct {
	messageCatalog *messages.Catalog
}

func newEntryFactory() *entryFactory {
	return &entryFactory{
		// Only support EN local bundles for now
		messageCatalog: messages.NewCatalog(messages.Locale_en_US),
	}
}

func (f *entryFactory) fromParameter(
	parameter parameter.ParameterWithDescriptionFromMessageCatalog,
	title string,
	typeOverride string,
) (userConfigEntry, error) {
	parameterType := typeOverride
	if parameterType == "" {
		switch parameterDefaultValue := parameter.GetDefaultValue().(type) {
		case string:
			parameterType = "string"
		case bool:
			parameterType = "boolean"
		default:
			return userConfigEntry{}, fmt.Errorf("unexpected type: %T", parameterDefaultValue)
		}
	}

	return userConfigEntry{
		Type:        parameterType,
		Title:       title,
		Description: f.messageCatalog.Get(parameter.GetDescriptionKey()),
		Required:    false, // All of our parameters are optional
		Default:     parameter.GetDefaultValue(),
	}, nil
}
