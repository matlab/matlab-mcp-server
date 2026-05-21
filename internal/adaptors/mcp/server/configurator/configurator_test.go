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
		ExtensionFiles().
		Return([]string{}).
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
		ExtensionFiles().
		Return([]string{expectedExtensionFilePath}).
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
		ExtensionFiles().
		Return([]string{expectedExtensionFilePath}).
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

func TestConfigurator_GetToolsToAdd_SingleMATLABSession_CrossFileNameCollision(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	mockCustomToolA := &toolsmocks.MockTool{}
	defer mockCustomToolA.AssertExpectations(t)

	mockCustomToolB := &toolsmocks.MockTool{}
	defer mockCustomToolB.AssertExpectations(t)

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

	expectedFilePathA := filepath.Join("config", "toolbox_a.json")
	expectedFilePathB := filepath.Join("config", "toolbox_b.json")
	expectedCollidingToolName := "run_model"

	mockCustomToolA.EXPECT().
		Name().
		Return(expectedCollidingToolName).
		Once()

	mockCustomToolB.EXPECT().
		Name().
		Return(expectedCollidingToolName).
		Once()

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
		ExtensionFiles().
		Return([]string{expectedFilePathA, expectedFilePathB}).
		Once()

	mockCustomToolFactory.EXPECT().
		LoadTools(expectedFilePathA).
		Return([]tools.Tool{mockCustomToolA}, nil).
		Once()

	mockCustomToolFactory.EXPECT().
		LoadTools(expectedFilePathB).
		Return([]tools.Tool{mockCustomToolB}, nil).
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
	require.Error(t, err, "GetToolsToAdd should return an error for cross-file name collision")
	assert.Nil(t, toolsToAdd, "Tools should be nil when cross-file name collision occurs")
	var collisionError *messages.StartupErrors_CustomToolNameCollisionAcrossFiles_Error
	require.ErrorAs(t, err, &collisionError)
	assert.Equal(t, expectedCollidingToolName, collisionError.Attr0)
	assert.Equal(t, expectedFilePathA, collisionError.Attr1)
	assert.Equal(t, expectedFilePathB, collisionError.Attr2)
}

func TestConfigurator_GetToolsToAdd_SingleMATLABSession_WithMultipleExtensionFiles_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	mockCustomToolA := &toolsmocks.MockTool{}
	defer mockCustomToolA.AssertExpectations(t)

	mockCustomToolB := &toolsmocks.MockTool{}
	defer mockCustomToolB.AssertExpectations(t)

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

	expectedFilePathA := filepath.Join("config", "toolbox_a.json")
	expectedFilePathB := filepath.Join("config", "toolbox_b.json")

	mockCustomToolA.EXPECT().
		Name().
		Return("tool_from_file_a").
		Once()

	mockCustomToolB.EXPECT().
		Name().
		Return("tool_from_file_b").
		Once()

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
		ExtensionFiles().
		Return([]string{expectedFilePathA, expectedFilePathB}).
		Once()

	mockCustomToolFactory.EXPECT().
		LoadTools(expectedFilePathA).
		Return([]tools.Tool{mockCustomToolA}, nil).
		Once()

	mockCustomToolFactory.EXPECT().
		LoadTools(expectedFilePathB).
		Return([]tools.Tool{mockCustomToolB}, nil).
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
	assert.Contains(t, toolsToAdd, mockCustomToolA, "GetToolsToAdd should include tool from file A")
	assert.Contains(t, toolsToAdd, mockCustomToolB, "GetToolsToAdd should include tool from file B")
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
		ExtensionFiles().
		Return([]string{expectedExtensionFilePath}).
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

func TestConfigurator_GetToolsToAdd_SingleMATLABSession_LoaderErrorOnSecondFile(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockApplicationDefinition := &mocks.MockApplicationDefinition{}
	defer mockApplicationDefinition.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCustomToolFactory := &mocks.MockCustomToolFactory{}
	defer mockCustomToolFactory.AssertExpectations(t)

	mockCustomToolA := &toolsmocks.MockTool{}
	defer mockCustomToolA.AssertExpectations(t)

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

	expectedFilePathA := filepath.Join("config", "toolbox_a.json")
	expectedFilePathB := filepath.Join("config", "toolbox_b.json")
	expectedError := messages.New_StartupErrors_FailedToParseExtensionFile_Error(expectedFilePathB)

	mockCustomToolA.EXPECT().
		Name().
		Return("tool_from_file_a").
		Once()

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
		ExtensionFiles().
		Return([]string{expectedFilePathA, expectedFilePathB}).
		Once()

	mockCustomToolFactory.EXPECT().
		LoadTools(expectedFilePathA).
		Return([]tools.Tool{mockCustomToolA}, nil).
		Once()

	mockCustomToolFactory.EXPECT().
		LoadTools(expectedFilePathB).
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
	require.ErrorAs(t, err, &expectedError, "GetToolsToAdd should return the error from the second file's loader")
	assert.Nil(t, toolsToAdd, "Tools should be nil when loader error occurs on second file")
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
