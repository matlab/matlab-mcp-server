// Copyright 2026 The MathWorks, Inc.

package validator_test

import (
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/loader/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validToolDefinition() definition.Tool {
	return definition.Tool{
		Name:        "test_tool",
		Title:       "Test Tool",
		Description: "A test tool",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"n": {Type: "number", Description: "A number"},
			},
			Required: []string{"n"},
		},
	}
}

func validSignatures() map[string]definition.Signature {
	return map[string]definition.Signature{
		"test_tool": {Function: "testFunc", Input: definition.SignatureInput{Order: []string{"n"}}},
	}
}

func TestValidator_Validate_HappyPath(t *testing.T) {
	// Arrange
	v := validator.NewValidator()
	td := validToolDefinition()
	signatures := validSignatures()

	// Act
	result, err := v.Validate(td, signatures)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, td.Name, result.Definition().Name)
	assert.Equal(t, "testFunc", result.Signature().Function)
}

func TestValidator_Validate_NoArgs_HappyPath(t *testing.T) {
	// Arrange
	v := validator.NewValidator()
	td := definition.Tool{
		Name:        "no_arg_tool",
		Title:       "No Arg Tool",
		Description: "A tool with no arguments",
		InputSchema: &jsonschema.Schema{Type: "object"},
	}
	signatures := map[string]definition.Signature{
		"no_arg_tool": {Function: "noArgFunc", Input: definition.SignatureInput{Order: []string{}}},
	}

	// Act
	result, err := v.Validate(td, signatures)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "no_arg_tool", result.Definition().Name)
}

func TestValidator_Validate_MissingRequiredField_ReturnsError(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*definition.Tool)
	}{
		{"missing name", func(td *definition.Tool) { td.Name = "" }},
		{"missing title", func(td *definition.Tool) { td.Title = "" }},
		{"missing description", func(td *definition.Tool) { td.Description = "" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			v := validator.NewValidator()
			td := validToolDefinition()
			tt.modify(&td)

			// Act
			_, err := v.Validate(td, validSignatures())

			// Assert
			require.Error(t, err)
			assert.ErrorIs(t, err, validator.ErrInvalidToolDefinition)
		})
	}
}

func TestValidator_Validate_InvalidSchema_ReturnsError(t *testing.T) {
	tests := []struct {
		name        string
		inputSchema *jsonschema.Schema
	}{
		{"nil schema", nil},
		{"missing type", &jsonschema.Schema{Properties: map[string]*jsonschema.Schema{}}},
		{"type not object", &jsonschema.Schema{Type: "array"}},
		{"unsupported property type", &jsonschema.Schema{
			Type:       "object",
			Properties: map[string]*jsonschema.Schema{"arr": {Type: "array"}},
		}},
		{"nil property", &jsonschema.Schema{
			Type:       "object",
			Properties: map[string]*jsonschema.Schema{"x": nil},
		}},
		{"property missing type", &jsonschema.Schema{
			Type:       "object",
			Properties: map[string]*jsonschema.Schema{"x": {}},
		}},
		{"required not in properties", &jsonschema.Schema{
			Type:       "object",
			Properties: map[string]*jsonschema.Schema{"x": {Type: "string"}},
			Required:   []string{"missing_prop"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			v := validator.NewValidator()
			td := validToolDefinition()
			td.InputSchema = tt.inputSchema

			// Act
			_, err := v.Validate(td, validSignatures())

			// Assert
			require.Error(t, err)
			assert.ErrorIs(t, err, validator.ErrInvalidInputSchema)
		})
	}
}

func TestValidator_Validate_SignatureNotFound(t *testing.T) {
	// Arrange
	v := validator.NewValidator()
	td := validToolDefinition()
	signatures := map[string]definition.Signature{}

	// Act
	_, err := v.Validate(td, signatures)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, validator.ErrSignatureNotFound)
}

func TestValidator_Validate_SignatureMissingFunction(t *testing.T) {
	// Arrange
	v := validator.NewValidator()
	td := validToolDefinition()
	signatures := map[string]definition.Signature{
		"test_tool": {Function: "", Input: definition.SignatureInput{Order: []string{}}},
	}

	// Act
	_, err := v.Validate(td, signatures)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, validator.ErrInvalidSignature)
}

func TestValidator_Validate_InvalidFunctionName_ReturnsError(t *testing.T) {
	tests := []struct {
		name     string
		function string
	}{
		{"injection attempt", "magic; system('rm -rf /')"},
		{"contains spaces", "my func"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			v := validator.NewValidator()
			td := validToolDefinition()
			signatures := map[string]definition.Signature{
				"test_tool": {Function: tt.function, Input: definition.SignatureInput{Order: []string{"n"}}},
			}

			// Act
			_, err := v.Validate(td, signatures)

			// Assert
			require.Error(t, err)
			assert.ErrorIs(t, err, validator.ErrInvalidSignature)
		})
	}
}

func TestValidator_Validate_SignatureDottedFunctionName_HappyPath(t *testing.T) {
	// Arrange
	v := validator.NewValidator()
	td := validToolDefinition()
	signatures := map[string]definition.Signature{
		"test_tool": {Function: "pkg.myFunc", Input: definition.SignatureInput{Order: []string{"n"}}},
	}

	// Act
	result, err := v.Validate(td, signatures)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "pkg.myFunc", result.Signature().Function)
}

func TestValidator_Validate_SignatureDuplicateOrderEntry(t *testing.T) {
	// Arrange
	v := validator.NewValidator()
	td := validToolDefinition()
	signatures := map[string]definition.Signature{
		"test_tool": {Function: "testFunc", Input: definition.SignatureInput{Order: []string{"n", "n"}}},
	}

	// Act
	_, err := v.Validate(td, signatures)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, validator.ErrInvalidSignature)
}

func TestValidator_Validate_NilPropertiesWithNonEmptyOrder_ReturnsError(t *testing.T) {
	// Arrange
	v := validator.NewValidator()
	td := definition.Tool{
		Name:        "test_tool",
		Title:       "Test Tool",
		Description: "A test tool",
		InputSchema: &jsonschema.Schema{Type: "object"},
	}
	signatures := map[string]definition.Signature{
		"test_tool": {Function: "testFunc", Input: definition.SignatureInput{Order: []string{"n"}}},
	}

	// Act
	_, err := v.Validate(td, signatures)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, validator.ErrInvalidSignature)
}

func TestValidator_Validate_SignatureOrderEntryNotInProperties(t *testing.T) {
	// Arrange
	v := validator.NewValidator()
	td := validToolDefinition()
	signatures := map[string]definition.Signature{
		"test_tool": {Function: "testFunc", Input: definition.SignatureInput{Order: []string{"nonexistent"}}},
	}

	// Act
	_, err := v.Validate(td, signatures)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, validator.ErrInvalidSignature)
}

func TestValidator_Validate_SchemaPropertyNotInOrder(t *testing.T) {
	// Arrange
	v := validator.NewValidator()
	td := validToolDefinition()
	signatures := map[string]definition.Signature{
		"test_tool": {Function: "testFunc", Input: definition.SignatureInput{Order: []string{}}},
	}

	// Act
	_, err := v.Validate(td, signatures)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, validator.ErrInvalidSignature)
}
