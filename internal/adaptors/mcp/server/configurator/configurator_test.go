// Copyright 2025-2026 The MathWorks, Inc.

package configurator_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/codingguidelines"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/plaintextlivecodegeneration"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/configurator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	evalmatlabmultisession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	evalmatlabsinglesession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/server/configurator"
	toolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
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

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
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
	}, "GetToolsToAdd should return all the injected tools for multi session")
}

func TestConfigurator_GetToolsToAdd_ConfigError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
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

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
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

	mockConfig.EXPECT().
		ExtensionFile().
		Return("").
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
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
	}, "GetToolsToAdd should return all injected tools for single session")
}

func TestConfigurator_GetToolsToAdd_SingleMATLABSession_WithCustomTools_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	mockCustomTool := &toolsmocks.MockTool{}
	defer mockCustomTool.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	expectedExtensionFilePath := filepath.Join("config", "tools.json")

	mockCustomTool.EXPECT().
		Name().
		Return("generate_magic_square")

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

	mockConfig.EXPECT().
		ExtensionFile().
		Return(expectedExtensionFilePath).
		Once()

	mockCustomToolFactory.EXPECT().
		LoadTools(expectedExtensionFilePath).
		Return([]tools.Tool{mockCustomTool}, nil).
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
	)

	// Act
	toolsToAdd, err := c.GetToolsToAdd()

	// Assert
	require.NoError(t, err, "GetToolsToAdd should not return an error")
	assert.Contains(t, toolsToAdd, mockCustomTool, "GetToolsToAdd should include the custom tool")
}

func TestConfigurator_GetToolsToAdd_SingleMATLABSession_CustomToolNameConflict(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	mockCustomTool := &toolsmocks.MockTool{}
	defer mockCustomTool.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := evalmatlabsinglesession.New(nil, nil, nil, nil)
	checkMATLABCodeInGlobalMATLABSession := checkmatlabcode.New(nil, nil, nil)
	detectMATLABToolboxesInSingleSessionTool := detectmatlabtoolboxes.New(nil, nil, nil)
	runMATLABFileInGlobalMATLABSessionTool := runmatlabfile.New(nil, nil, nil, nil)
	runMATLABTestFileInGlobalMATLABSessionTool := runmatlabtestfile.New(nil, nil, nil)
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	expectedExtensionFilePath := filepath.Join("config", "tools.json")
	expectedConflictingToolName := "evaluate_matlab_code"

	mockCustomTool.EXPECT().
		Name().
		Return(expectedConflictingToolName)

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

	mockConfig.EXPECT().
		ExtensionFile().
		Return(expectedExtensionFilePath).
		Once()

	mockCustomToolFactory.EXPECT().
		LoadTools(expectedExtensionFilePath).
		Return([]tools.Tool{mockCustomTool}, nil).
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
	)

	// Act
	toolsToAdd, err := c.GetToolsToAdd()

	// Assert
	require.Error(t, err, "GetToolsToAdd should return an error for conflicting tool name")
	assert.Nil(t, toolsToAdd, "Tools should be nil when name conflict occurs")
	var nameConflictError *messages.StartupErrors_CustomToolNameConflict_Error
	require.ErrorAs(t, err, &nameConflictError)
	assert.Equal(t, expectedConflictingToolName, nameConflictError.Attr0)
	assert.Equal(t, expectedExtensionFilePath, nameConflictError.Attr1)
}

func TestConfigurator_GetToolsToAdd_SingleMATLABSession_LoaderError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}
	plaintextlivecodegenerationResource := &plaintextlivecodegeneration.Resource{}

	expectedExtensionFilePath := filepath.Join("config", "tools.json")
	expectedError := messages.AnError

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

	mockConfig.EXPECT().
		ExtensionFile().
		Return(expectedExtensionFilePath).
		Once()

	mockCustomToolFactory.EXPECT().
		LoadTools(expectedExtensionFilePath).
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
	)

	// Act
	toolsToAdd, err := c.GetToolsToAdd()

	// Assert
	require.ErrorIs(t, err, expectedError, "GetToolsToAdd should return the error from the loader")
	assert.Nil(t, toolsToAdd, "Tools should be nil when loader error occurs")
}

func TestConfigurator_GetResourcesToAdd_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
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

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
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

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
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
		codingGuidelinesResource,
		plaintextlivecodegenerationResource,
		mockCustomToolFactory,
	)

	// Act
	result := c.GetResourcesToAdd()

	// Assert
	assert.Empty(t, result)
}
