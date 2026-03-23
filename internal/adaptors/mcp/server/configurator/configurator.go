// Copyright 2025-2026 The MathWorks, Inc.

package configurator

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/codingguidelines"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/plaintextlivecodegeneration"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	evalmatlabcodemultisession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	runmatlabbuildmultisession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/runmatlabbuild"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	evalmatlabcodesinglesession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	runmatlabbuildsinglesession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabbuild"
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

type Configurator struct {
	configFactory    ConfigFactory
	featuresProvider ApplicationDefinition

	// Multi Session tools
	listAvailableMATLABsTool tools.Tool
	startMATLABSessionTool   tools.Tool
	stopMATLABSessionTool    tools.Tool
	evalInMATLABSessionTool  tools.Tool

	// Single Session tools
	evalInGlobalMATLABSessionTool                  tools.Tool
	checkMATLABCodeInGlobalMATLABSessionTool       tools.Tool
	detectMATLABToolboxesInGlobalMATLABSessionTool tools.Tool
	runMATLABFileInGlobalMATLABSessionTool         tools.Tool
	runMATLABTestFileInGlobalMATLABSessionTool     tools.Tool
	runMATLABBuildInGlobalMATLABSessionTool        tools.Tool

	// Multi Session build tool
	runMATLABBuildInMATLABSessionTool tools.Tool

	// Resources
	codingGuidelinesResource            resources.Resource
	plaintextlivecodegenerationResource resources.Resource
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
	runMATLABBuildInGlobalMATLABSessionTool *runmatlabbuildsinglesession.Tool,
	runMATLABBuildInMATLABSessionTool *runmatlabbuildmultisession.Tool,

	codingGuidelinesResource *codingguidelines.Resource,
	plaintextlivecodegenerationResource *plaintextlivecodegeneration.Resource,
) *Configurator {
	return &Configurator{
		configFactory: configFactory,

		featuresProvider: featuresProvider,

		listAvailableMATLABsTool: listAvailableMATLABsTool,
		startMATLABSessionTool:   startMATLABSessionTool,
		stopMATLABSessionTool:    stopMATLABSessionTool,
		evalInMATLABSessionTool:  evalInMATLABSessionTool,

		evalInGlobalMATLABSessionTool:                  evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSessionTool:       checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInGlobalMATLABSessionTool: detectMATLABToolboxesInGlobalMATLABSessionTool,
		runMATLABFileInGlobalMATLABSessionTool:         runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool:     runMATLABTestFileInGlobalMATLABSessionTool,
		runMATLABBuildInGlobalMATLABSessionTool:        runMATLABBuildInGlobalMATLABSessionTool,
		runMATLABBuildInMATLABSessionTool:              runMATLABBuildInMATLABSessionTool,

		codingGuidelinesResource:            codingGuidelinesResource,
		plaintextlivecodegenerationResource: plaintextlivecodegenerationResource,
	}
}

func (c *Configurator) GetToolsToAdd() ([]tools.Tool, error) {
	if !c.featuresProvider.Features().MATLAB.Enabled {
		return []tools.Tool{}, nil
	}

	config, err := c.configFactory.Config()
	if err != nil {
		return nil, err
	}

	if config.UseSingleMATLABSession() {
		return []tools.Tool{
			c.evalInGlobalMATLABSessionTool,
			c.checkMATLABCodeInGlobalMATLABSessionTool,
			c.detectMATLABToolboxesInGlobalMATLABSessionTool,
			c.runMATLABFileInGlobalMATLABSessionTool,
			c.runMATLABTestFileInGlobalMATLABSessionTool,
			c.runMATLABBuildInGlobalMATLABSessionTool,
		}, nil
	}

	return []tools.Tool{
		c.listAvailableMATLABsTool,
		c.startMATLABSessionTool,
		c.stopMATLABSessionTool,
		c.evalInMATLABSessionTool,
		c.runMATLABBuildInMATLABSessionTool,
	}, nil
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
