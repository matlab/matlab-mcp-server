// Copyright 2026 The MathWorks, Inc.

package selector

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ApplicationDefinition interface {
	Features() definition.Features
}

type MessageCatalog interface {
	Get(message messages.MessageKey) string
}

type Selector struct {
	applicationDefinition ApplicationDefinition
	messageCatalog        MessageCatalog
}

func New(
	applicationDefinition ApplicationDefinition,
	messageCatalog MessageCatalog,
) *Selector {
	return &Selector{
		applicationDefinition: applicationDefinition,
		messageCatalog:        messageCatalog,
	}
}

func (s *Selector) DefaultParameters() []entities.Parameter {
	parameterDefs := []parameter.ParameterWithDescriptionFromMessageCatalog{
		defaultparameters.HelpMode(),
		defaultparameters.VersionMode(),
		defaultparameters.SetupMATLABMode(),
		defaultparameters.BaseDir(),
		defaultparameters.LogLevel(),
		defaultparameters.DuplicateLogsToStderr(),
		defaultparameters.WatchdogMode(),
		defaultparameters.ServerInstanceID(),
		defaultparameters.DisableTelemetry(),
		defaultparameters.TelemetryCollectorEndpoint(),
		defaultparameters.TelemetryCollectionInterval(),
		defaultparameters.TelemetryCollectorEndpointInsecure(),
	}

	matlabParameters := []parameter.ParameterWithDescriptionFromMessageCatalog{
		defaultparameters.PreferredLocalMATLABRoot(),
		defaultparameters.PreferredMATLABStartingDirectory(),
		defaultparameters.UseSingleMATLABSession(),
		defaultparameters.InitializeMATLABOnStartup(),
		defaultparameters.MATLABDisplayMode(),
		defaultparameters.MATLABSessionMode(),
		defaultparameters.MATLABSessionConnectionDetails(),
		defaultparameters.MATLABSessionConnectionTimeout(),
		defaultparameters.MATLABSessionDiscoveryTimeout(),
		defaultparameters.EmbeddedConnectorDetailsTimeout(),
		defaultparameters.ExtensionFiles(),
	}

	matlabFeature := s.applicationDefinition.Features().MATLAB
	if !matlabFeature.Enabled {
		for _, matlabParameter := range matlabParameters {
			matlabParameter.SetActive(false)
		}
	}

	parameterDefs = append(parameterDefs, matlabParameters...)

	parameters := make([]entities.Parameter, len(parameterDefs))
	for i, parameterDef := range parameterDefs {
		// Only resolve the description of not hidden (visible) flags
		if !parameterDef.GetHiddenFlag() {
			s.resolveDescription(parameterDef)
		}
		parameters[i] = parameterDef
	}
	return parameters
}

func (s *Selector) resolveDescription(p parameter.ParameterWithDescriptionFromMessageCatalog) {
	p.SetDescription(s.messageCatalog.Get(p.GetDescriptionKey()))
}
