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

func TestParser_Parse_BoolFlagImplicitTrue(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "bool-param"
	paramFlagName := "bool-flag"

	mockParam := newMockParam(
		t,
		paramID,
		paramFlagName,
		"",
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

	args := []string{"--" + paramFlagName}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, true, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_HiddenFlag(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "hidden-param"
	paramFlagName := "hidden-flag"
	expectedValue := "hidden-value"

	mockParam := newMockParam(
		t,
		paramID,
		paramFlagName,
		"",
		"default",
		"Hidden flag description",
		true,
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

	args := []string{"--" + paramFlagName + "=" + expectedValue}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedValue, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_BadFlag(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{}).
		Once()

	badFlagName := "nonexistent"
	args := []string{"--" + badFlagName + "=value"}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	usage, _ := p.Usage()
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_BadFlag_Error(badFlagName, "\n", usage)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Parse_DurationFlag(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "duration-param"
	paramFlagName := "my-duration"

	mockParam := newMockParam(
		t,
		paramID,
		paramFlagName,
		"",
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

	args := []string{"--" + paramFlagName + "=5m30s"}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 5*time.Minute+30*time.Second, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_BadDurationFlagValue(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramFlagName := "my-duration"
	badValue := "notaduration"

	mockParam := newMockParam(
		t,
		"duration-param",
		paramFlagName,
		"",
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

	args := []string{"--" + paramFlagName + "=" + badValue}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_BadValue_Error(badValue, paramFlagName)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Parse_BadBoolFlagValue(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramFlagName := "bool-flag"
	badValue := "notabool"

	mockParam := newMockParam(
		t,
		"bool-param",
		paramFlagName,
		"",
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

	args := []string{"--" + paramFlagName + "=" + badValue}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_BadValue_Error(badValue, paramFlagName)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Parse_StringArrayFlag_SingleValue(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "string-array-param"
	paramFlagName := "my-files"

	mockParam := newMockParam(
		t,
		paramID,
		paramFlagName,
		"",
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

	expectedValue := filepath.Join("path", "to", "file.json")
	args := []string{"--" + paramFlagName + "=" + expectedValue}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	_, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{expectedValue}, result[paramID])
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_StringArrayFlag_MultipleValues(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "string-array-param"
	paramFlagName := "my-files"

	mockParam := newMockParam(
		t,
		paramID,
		paramFlagName,
		"",
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

	expectedValueA := filepath.Join("path", "to", "a.json")
	expectedValueB := filepath.Join("path", "to", "b.json")
	args := []string{
		"--" + paramFlagName + "=" + expectedValueA,
		"--" + paramFlagName + "=" + expectedValueB,
	}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	_, result, specifiedParameters, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{expectedValueA, expectedValueB}, result[paramID])
	assert.Equal(t, []string{paramID}, specifiedParameters)
}

func TestParser_Parse_InactiveParameterFlagSkipped(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	flagName := "inactive-flag"

	mockParam := newMockParam(
		t,
		"inactive-param",
		flagName,
		"",
		"default",
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

	args := []string{"--" + flagName + "=value"}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	usage, _ := p.Usage()
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_BadFlag_Error(flagName, "\n", usage)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}

func TestParser_Parse_BadFlagSyntax(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{}).
		Once()

	badArg := "---bad-syntax"
	args := []string{badArg}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	usage, _ := p.Usage()
	parameters, result, specifiedParameters, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_BadSyntax_Error(badArg, "\n", usage)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
	assert.Nil(t, specifiedParameters)
}
