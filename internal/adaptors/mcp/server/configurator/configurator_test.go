// Copyright 2025-2026 The MathWorks, Inc.

package configurator_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/codingguidelines"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/plaintextlivecodegeneration"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/configurator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	evalmatlabmultisession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	runmatlabbuildmultisession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/runmatlabbuild"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	evalmatlabsinglesession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	runmatlabbuildsinglesession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabbuild"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/server/configurator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	runMATLABBuildInGlobalMATLABSessionTool := &runmatlabbuildsinglesession.Tool{}
	runMATLABBuildInMATLABSessionTool := &runmatlabbuildmultisession.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	// Act
	result := configurator.New(
		mockConfigFactory,
		mockApplicationDefinition,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		runMATLABBuildInGlobalMATLABSessionTool,
		runMATLABBuildInMATLABSessionTool,
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
	)

	// Assert
	require.NotNil(t, result, "Configurator should not be nil")
}

func TestConfigurator_GetToolsToAdd_MultipleMATLABSession_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	runMATLABBuildInGlobalMATLABSessionTool := &runmatlabbuildsinglesession.Tool{}
	runMATLABBuildInMATLABSessionTool := &runmatlabbuildmultisession.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	mockApplicationDefinition.EXPECT().
		Features().
		Return(definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(false).
		Once()

	c := configurator.New(
		mockConfigFactory,
		mockApplicationDefinition,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		runMATLABBuildInGlobalMATLABSessionTool,
		runMATLABBuildInMATLABSessionTool,
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
	)

	// Act
	toolsToAdd, err := c.GetToolsToAdd()

	// Assert
	require.NoError(t, err, "GetToolsToAdd should not return an error")
	assert.ElementsMatch(t, toolsToAdd, []tools.Tool{
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		runMATLABBuildInMATLABSessionTool,
	}, "GetToolsToAdd should return all the injected tools for multi session")
}

func TestConfigurator_GetToolsToAdd_ConfigError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	runMATLABBuildInGlobalMATLABSessionTool := &runmatlabbuildsinglesession.Tool{}
	runMATLABBuildInMATLABSessionTool := &runmatlabbuildmultisession.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	expectedError := messages.AnError

	mockApplicationDefinition.EXPECT().
		Features().
		Return(definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	c := configurator.New(
		mockConfigFactory,
		mockApplicationDefinition,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		runMATLABBuildInGlobalMATLABSessionTool,
		runMATLABBuildInMATLABSessionTool,
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
	)

	// Act
	toolsToAdd, err := c.GetToolsToAdd()

	// Assert
	require.ErrorIs(t, err, expectedError, "GetToolsToAdd should return the error from Config")
	assert.Nil(t, toolsToAdd, "Tools should be nil when error occurs")
}

func TestConfigurator_GetToolsToAdd_SingleMATLABSession_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	runMATLABBuildInGlobalMATLABSessionTool := &runmatlabbuildsinglesession.Tool{}
	runMATLABBuildInMATLABSessionTool := &runmatlabbuildmultisession.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	mockApplicationDefinition.EXPECT().
		Features().
		Return(definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	c := configurator.New(
		mockConfigFactory,
		mockApplicationDefinition,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		runMATLABBuildInGlobalMATLABSessionTool,
		runMATLABBuildInMATLABSessionTool,
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
	)

	// Act
	toolsToAdd, err := c.GetToolsToAdd()

	// Assert
	require.NoError(t, err, "GetToolsToAdd should not return an error")
	assert.ElementsMatch(t, toolsToAdd, []tools.Tool{
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABBuildInGlobalMATLABSessionTool,
	}, "GetToolsToAdd should all injected tools for single session")
}

func TestConfigurator_GetResourcesToAdd_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	runMATLABBuildInGlobalMATLABSessionTool := &runmatlabbuildsinglesession.Tool{}
	runMATLABBuildInMATLABSessionTool := &runmatlabbuildmultisession.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	mockApplicationDefinition.EXPECT().
		Features().
		Return(definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}).
		Once()

	c := configurator.New(
		mockConfigFactory,
		mockApplicationDefinition,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		runMATLABBuildInGlobalMATLABSessionTool,
		runMATLABBuildInMATLABSessionTool,
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
	)

	// Act
	result := c.GetResourcesToAdd()

	// Assert
	assert.ElementsMatch(t, []resources.Resource{codingGuidelinesResource, plaintextlivecodegenerationResource}, result)
}

func TestConfigurator_GetToolsToAdd_MATLABFeatureDisabled(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	runMATLABBuildInGlobalMATLABSessionTool := &runmatlabbuildsinglesession.Tool{}
	runMATLABBuildInMATLABSessionTool := &runmatlabbuildmultisession.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	mockApplicationDefinition.EXPECT().
		Features().
		Return(definition.Features{MATLAB: definition.MATLABFeature{Enabled: false}}).
		Once()

	c := configurator.New(
		mockConfigFactory,
		mockApplicationDefinition,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		runMATLABBuildInGlobalMATLABSessionTool,
		runMATLABBuildInMATLABSessionTool,
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
	)

	// Act
	toolsToAdd, err := c.GetToolsToAdd()

	// Assert
	require.NoError(t, err)
	assert.Empty(t, toolsToAdd)
}

func TestConfigurator_GetResourcesToAdd_MATLABFeatureDisabled(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	runMATLABBuildInGlobalMATLABSessionTool := &runmatlabbuildsinglesession.Tool{}
	runMATLABBuildInMATLABSessionTool := &runmatlabbuildmultisession.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	mockApplicationDefinition.EXPECT().
		Features().
		Return(definition.Features{MATLAB: definition.MATLABFeature{Enabled: false}}).
		Once()

	c := configurator.New(
		mockConfigFactory,
		mockApplicationDefinition,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		runMATLABBuildInGlobalMATLABSessionTool,
		runMATLABBuildInMATLABSessionTool,
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
	)

	// Act
	result := c.GetResourcesToAdd()

	// Assert
	assert.Empty(t, result)
}
