// Copyright 2025-2026 The MathWorks, Inc.

package parser_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/parser"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	parsermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/parameter/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse_BoolEnvVar(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "bool-param"
	paramEnvVar := "BOOL_ENV_VAR"

	mockParam := newMockParam(
		t,
		paramID,
		"bool-flag",
		paramEnvVar,
		false,
		"Test bool description",
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
		Return("true", true).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, true, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_InactiveParameterEnvVarSkipped(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "inactive-param"
	paramEnvVar := "INACTIVE_ENV_VAR"
	paramDefaultValue := "default-value"

	mockParam := newMockParam(
		t,
		paramID,
		"",
		paramEnvVar,
		paramDefaultValue,
		"",
		false,
		false,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
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

func TestParser_Parse_BadEnvVarBoolValue(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramEnvVar := "BOOL_ENV_VAR"
	badEnvValue := "notabool"

	mockParam := newMockParam(
		t,
		"bool-param",
		"bool-flag",
		paramEnvVar,
		false,
		"Test bool description",
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
		Return(badEnvValue, true).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_BadValueForEnvVar_Error(badEnvValue, paramEnvVar)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Parse_DurationEnvVar(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "duration-param"
	paramEnvVar := "DURATION_ENV_VAR"

	mockParam := newMockParam(
		t,
		paramID,
		"duration-flag",
		paramEnvVar,
		time.Minute,
		"Test duration description",
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
		Return("5m30s", true).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 5*time.Minute+30*time.Second, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_StringArrayEnvVar(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "string-array-param"
	paramEnvVar := "FILES_ENV_VAR"

	mockParam := newMockParam(
		t,
		paramID,
		"my-files",
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

	envValue := filepath.Join("path", "to", "a.json")

	mockOSLayer.EXPECT().
		LookupEnv(paramEnvVar).
		Return(envValue, true).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	_, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{envValue}, result[paramID])
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_BadEnvVarDurationValue(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramEnvVar := "DURATION_ENV_VAR"
	badEnvValue := "notaduration"

	mockParam := newMockParam(
		t,
		"duration-param",
		"duration-flag",
		paramEnvVar,
		time.Minute,
		"Test duration description",
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
		Return(badEnvValue, true).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_BadValueForEnvVar_Error(badEnvValue, paramEnvVar)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}
