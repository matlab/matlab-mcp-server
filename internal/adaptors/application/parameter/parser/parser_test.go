// Copyright 2025-2026 The MathWorks, Inc.

package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/parser"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	parsermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/parameter/parser"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "test-param"
	paramEnvVar := "TEST_ENV_VAR"
	paramDefaultValue := "default-value"

	mockParam := newMockParam(
		t,
		paramID,
		"test-flag",
		paramEnvVar,
		paramDefaultValue,
		"Test description",
		false,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{}).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv(paramEnvVar).
		Return("", false).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, paramDefaultValue, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Empty(t, specifiedParameters)
}

func TestParser_Parse_FlagOverridesDefault(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "test-param"
	paramFlagName := "test-flag"
	paramEnvVar := "TEST_ENV_VAR"
	expectedFlagValue := "flag-value"

	mockParam := newMockParam(
		t,
		paramID,
		paramFlagName,
		paramEnvVar,
		"default-value",
		"Test description",
		false,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv(paramEnvVar).
		Return("", false).
		Once()

	args := []string{"--" + paramFlagName + "=" + expectedFlagValue}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedFlagValue, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_EnvVarOverridesDefault(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "test-param"
	paramEnvVar := "TEST_ENV_VAR"
	expectedEnvValue := "env-value"

	mockParam := newMockParam(
		t,
		paramID,
		"test-flag",
		paramEnvVar,
		"default-value",
		"Test description",
		false,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv(paramEnvVar).
		Return(expectedEnvValue, true).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedEnvValue, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_FlagOverridesEnvVar(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "test-param"
	paramFlagName := "test-flag"
	paramEnvVar := "TEST_ENV_VAR"
	envValue := "env-value"
	expectedFlagValue := "flag-value"

	mockParam := newMockParam(
		t,
		paramID,
		paramFlagName,
		paramEnvVar,
		"default-value",
		"Test description",
		false,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv(paramEnvVar).
		Return(envValue, true).
		Once()

	args := []string{"--" + paramFlagName + "=" + expectedFlagValue}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedFlagValue, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_StringArrayFlagAppendsToEnvVar(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "string-array-param"
	paramFlagName := "my-files"
	paramEnvVar := "FILES_ENV_VAR"

	mockParam := newMockParam(
		t,
		paramID,
		paramFlagName,
		paramEnvVar,
		[]string{},
		"Test string array description",
		false,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	envFile := filepath.Join("env", "file.json")
	cliFile := filepath.Join("cli", "file.json")

	mockOSLayer.EXPECT().
		LookupEnv(paramEnvVar).
		Return(envFile, true).
		Once()

	args := []string{"--" + paramFlagName + "=" + cliFile}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	_, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{envFile, cliFile}, result[paramID])
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_ParameterWithNoFlag(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "no-flag-param"
	paramEnvVar := "NO_FLAG_ENV_VAR"
	expectedEnvValue := "env-value"

	mockParam := newMockParam(
		t,
		paramID,
		"",
		paramEnvVar,
		"default",
		"",
		false,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv(paramEnvVar).
		Return(expectedEnvValue, true).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedEnvValue, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_EmptyParameterID(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	mockParam := &entitiesmocks.MockParameter{}
	defer mockParam.AssertExpectations(t)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	mockParam.EXPECT().
		GetID().
		Return("").
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_InvalidParameterKey_Error("")
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Parse_DuplicateParameterID(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	mockParam1 := &entitiesmocks.MockParameter{}
	defer mockParam1.AssertExpectations(t)

	mockParam2 := &entitiesmocks.MockParameter{}
	defer mockParam2.AssertExpectations(t)

	duplicateID := "duplicate-id"

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{mockParam1}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam2}).
		Once()

	mockParam1.EXPECT().
		GetID().
		Return(duplicateID).
		Once()

	mockParam1.EXPECT().
		GetFlagName().
		Return("flag1").
		Once()

	mockParam1.EXPECT().
		GetEnvVarName().
		Return("ENV1").
		Once()

	mockParam2.EXPECT().
		GetID().
		Return(duplicateID).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_DuplicateParameter_Error(duplicateID, "parameter ID", duplicateID)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Parse_DuplicateFlagName(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	mockParam1 := &entitiesmocks.MockParameter{}
	defer mockParam1.AssertExpectations(t)

	mockParam2 := &entitiesmocks.MockParameter{}
	defer mockParam2.AssertExpectations(t)

	duplicateFlagName := "duplicate-flag"
	param2ID := "param2"

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{mockParam1}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam2}).
		Once()

	mockParam1.EXPECT().
		GetID().
		Return("param1").
		Once()

	mockParam1.EXPECT().
		GetFlagName().
		Return(duplicateFlagName).
		Once()

	mockParam1.EXPECT().
		GetEnvVarName().
		Return("ENV1").
		Once()

	mockParam2.EXPECT().
		GetID().
		Return(param2ID).
		Once()

	mockParam2.EXPECT().
		GetFlagName().
		Return(duplicateFlagName).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_DuplicateParameter_Error(param2ID, "flag name", duplicateFlagName)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Parse_DuplicateEnvVarName(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	mockParam1 := &entitiesmocks.MockParameter{}
	defer mockParam1.AssertExpectations(t)

	mockParam2 := &entitiesmocks.MockParameter{}
	defer mockParam2.AssertExpectations(t)

	duplicateEnvVar := "DUPLICATE_ENV"
	param2ID := "param2"

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{mockParam1}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam2}).
		Once()

	mockParam1.EXPECT().
		GetID().
		Return("param1").
		Once()

	mockParam1.EXPECT().
		GetFlagName().
		Return("flag1").
		Once()

	mockParam1.EXPECT().
		GetEnvVarName().
		Return(duplicateEnvVar).
		Once()

	mockParam2.EXPECT().
		GetID().
		Return(param2ID).
		Once()

	mockParam2.EXPECT().
		GetFlagName().
		Return("flag2").
		Once()

	mockParam2.EXPECT().
		GetEnvVarName().
		Return(duplicateEnvVar).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_DuplicateParameter_Error(param2ID, "env var name", duplicateEnvVar)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Parse_DuplicateEnvVarNameCaseInsensitive(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	mockParam1 := &entitiesmocks.MockParameter{}
	defer mockParam1.AssertExpectations(t)

	mockParam2 := &entitiesmocks.MockParameter{}
	defer mockParam2.AssertExpectations(t)

	param2ID := "param2"

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{mockParam1}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam2}).
		Once()

	mockParam1.EXPECT().
		GetID().
		Return("param1").
		Once()

	mockParam1.EXPECT().
		GetFlagName().
		Return("flag1").
		Once()

	mockParam1.EXPECT().
		GetEnvVarName().
		Return("duplicate_env").
		Once()

	mockParam2.EXPECT().
		GetID().
		Return(param2ID).
		Once()

	mockParam2.EXPECT().
		GetFlagName().
		Return("flag2").
		Once()

	mockParam2.EXPECT().
		GetEnvVarName().
		Return("DUPLICATE_ENV").
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_DuplicateParameter_Error(param2ID, "env var name", "DUPLICATE_ENV")
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Usage_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramFlagName := "test-flag"
	paramDescription := "Test description"

	mockParam := newMockParam(
		t,
		"test-param",
		paramFlagName,
		"",
		"default",
		paramDescription,
		false,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	usage, err := p.Usage()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, usage, "--"+paramFlagName)
	assert.Contains(t, usage, paramDescription)
}

func TestParser_Usage_HiddenFlagNotShown(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	visibleFlagName := "visible-flag"
	hiddenFlagName := "hidden-flag"

	mockVisibleParam := newMockParam(
		t,
		"visible-param",
		visibleFlagName,
		"",
		"default",
		"Visible description",
		false,
		true,
	)
	mockHiddenParam := newMockParam(
		t,
		"hidden-param",
		hiddenFlagName,
		"",
		"default",
		"Hidden description",
		true,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockVisibleParam, mockHiddenParam}).
		Once()

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	usage, err := p.Usage()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, usage, "--"+visibleFlagName)
	assert.NotContains(t, usage, "--"+hiddenFlagName)
}

func TestParser_Usage_PropagatesInitError(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	mockParam1 := &entitiesmocks.MockParameter{}
	defer mockParam1.AssertExpectations(t)

	mockParam2 := &entitiesmocks.MockParameter{}
	defer mockParam2.AssertExpectations(t)

	duplicateID := "duplicate-id"

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{mockParam1}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam2}).
		Once()

	mockParam1.EXPECT().
		GetID().
		Return(duplicateID).
		Once()

	mockParam1.EXPECT().
		GetFlagName().
		Return("flag1").
		Once()

	mockParam1.EXPECT().
		GetEnvVarName().
		Return("ENV1").
		Once()

	mockParam2.EXPECT().
		GetID().
		Return(duplicateID).
		Once()

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	usage, err := p.Usage()

	// Assert
	expectedError := messages.New_StartupErrors_DuplicateParameter_Error(duplicateID, "parameter ID", duplicateID)
	require.Equal(t, expectedError, err)
	assert.Empty(t, usage)
}

// newMockParam creates a MockParameter with expectations configured.
// Note: Parameter methods (GetID, GetFlagName, etc.) do not use .Once() because
// the parser may call them multiple times during initialization (for duplicate
// detection, flag setup, usage generation) and during parsing. The exact call
// count depends on internal implementation details we don't want to couple to.
//
// When flagName is empty, GetDescription and GetHiddenFlag expectations are not
// set because the parser only calls these methods for parameters that have flags.
func newMockParam(
	t *testing.T,
	id string,
	flagName string,
	envVarName string,
	defaultValue any,
	description string,
	hidden bool,
	active bool,
) *entitiesmocks.MockParameter {
	mock := &entitiesmocks.MockParameter{}
	t.Cleanup(func() { mock.AssertExpectations(t) })

	mock.EXPECT().
		GetID().
		Return(id)

	mock.EXPECT().
		GetFlagName().
		Return(flagName)

	mock.EXPECT().
		GetEnvVarName().
		Return(envVarName)

	mock.EXPECT().
		GetActive().
		Return(active)

	mock.EXPECT().
		GetDefaultValue().
		Return(defaultValue)

	// Description and HiddenFlag are only used for flag setup of active parameters
	if flagName != "" && active {
		mock.EXPECT().
			GetDescription().
			Return(description)

		mock.EXPECT().
			GetHiddenFlag().
			Return(hidden)
	}

	return mock
}
