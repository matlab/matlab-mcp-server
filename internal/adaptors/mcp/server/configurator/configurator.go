// Copyright 2025-2026 The MathWorks, Inc.

package configurator

import (
	"slices"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/codingguidelines"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/plaintextlivecodegeneration"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	evalmatlabcodemultisession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	evalmatlabcodesinglesession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type ApplicationDefinition interface {
	Features() definition.Features
}

type CustomToolFactory interface {
	LoadTools(filePath string) ([]tools.Tool, messages.Error)
}

type Configurator struct {
	configFactory    ConfigFactory
	featuresProvider ApplicationDefinition

	// Built-in tools
	multiSessionTools  []tools.Tool
	singleSessionTools []tools.Tool

	// Resources
	codingGuidelinesResource            resources.Resource
	plaintextlivecodegenerationResource resources.Resource

	// Custom tool dependencies
	customToolFactory CustomToolFactory
}

func New(
	configFactory ConfigFactory,

	featuresProvider ApplicationDefinition,

	listAvailableMATLABsTool *listavailablematlabs.Tool,
	startMATLABSessionTool *startmatlabsession.Tool,
	stopMATLABSessionTool *stopmatlabsession.Tool,
	evalInMATLABSessionTool *evalmatlabcodemultisession.Tool,

	evalInGlobalMATLABSessionTool *evalmatlabcodesinglesession.Tool,
	checkMATLABCodeInGlobalMATLABSession *checkmatlabcode.Tool,
	detectMATLABToolboxesInGlobalMATLABSessionTool *detectmatlabtoolboxes.Tool,
	runMATLABFileInGlobalMATLABSessionTool *runmatlabfile.Tool,
	runMATLABTestFileInGlobalMATLABSessionTool *runmatlabtestfile.Tool,

	codingGuidelinesResource *codingguidelines.Resource,
	plaintextlivecodegenerationResource *plaintextlivecodegeneration.Resource,

	customToolFactory CustomToolFactory,
) *Configurator {
	return &Configurator{
		configFactory: configFactory,

		featuresProvider: featuresProvider,

		multiSessionTools: []tools.Tool{
			listAvailableMATLABsTool,
			startMATLABSessionTool,
			stopMATLABSessionTool,
			evalInMATLABSessionTool,
		},

		singleSessionTools: []tools.Tool{
			evalInGlobalMATLABSessionTool,
			checkMATLABCodeInGlobalMATLABSession,
			detectMATLABToolboxesInGlobalMATLABSessionTool,
			runMATLABFileInGlobalMATLABSessionTool,
			runMATLABTestFileInGlobalMATLABSessionTool,
		},

		codingGuidelinesResource:            codingGuidelinesResource,
		plaintextlivecodegenerationResource: plaintextlivecodegenerationResource,

		customToolFactory: customToolFactory,
	}
}

func (c *Configurator) GetToolsToAdd() ([]tools.Tool, error) {
	if !c.featuresProvider.Features().MATLAB.Enabled {
		return []tools.Tool{}, nil
	}

	cfg, err := c.configFactory.Config()
	if err != nil {
		return nil, err
	}

	if cfg.UseSingleMATLABSession() {
		customTools, err := c.loadCustomTools(cfg)
		if err != nil {
			return nil, err
		}

		return slices.Concat(c.singleSessionTools, customTools), nil
	}

	return slices.Clone(c.multiSessionTools), nil
}

func (c *Configurator) loadCustomTools(cfg config.Config) ([]tools.Tool, error) {
	extensionFilePath := cfg.ExtensionFile()
	if extensionFilePath == "" {
		return nil, nil
	}

	customTools, err := c.customToolFactory.LoadTools(extensionFilePath)
	if err != nil {
		return nil, err
	}

	for _, t := range customTools {
		if c.isBuiltInSingleSessionToolName(t.Name()) {
			return nil, messages.New_StartupErrors_CustomToolNameConflict_Error(
				t.Name(),
				extensionFilePath,
			)
		}
	}

	return customTools, nil
}

func (c *Configurator) isBuiltInSingleSessionToolName(name string) bool {
	for _, t := range c.singleSessionTools {
		if t.Name() == name {
			return true
		}
	}
	return false
}

func (c *Configurator) GetResourcesToAdd() []resources.Resource {
	if !c.featuresProvider.Features().MATLAB.Enabled {
		return []resources.Resource{}
	}

	return []resources.Resource{
		c.codingGuidelinesResource,
		c.plaintextlivecodegenerationResource,
	}
}
