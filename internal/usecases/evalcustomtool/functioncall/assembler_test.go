// Copyright 2026 The MathWorks, Inc.

package functioncall_test

import (
	"math"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalcustomtool/functioncall"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAssembler_HappyPath(t *testing.T) {
	// Act
	assembler := functioncall.NewAssembler()

	// Assert
	assert.NotNil(t, assembler)
}

func TestAssembler_Assemble_SingleArg_HappyPath(t *testing.T) {
	tests := []struct {
		name           string
		propType       string
		value          any
		expectedResult string
	}{
		{"positive number", "number", float64(5), "myFunc(5)"},
		{"negative number", "number", float64(-3.14), "myFunc(-3.14)"},
		{"floating point number", "number", float64(3.14), "myFunc(3.14)"},
		{"positive integer", "integer", float64(5), "myFunc(5)"},
		{"negative integer", "integer", float64(-7), "myFunc(-7)"},
		{"boolean true", "boolean", true, "myFunc(true)"},
		{"boolean false", "boolean", false, "myFunc(false)"},
		{"string", "string", "hello", `myFunc("hello")`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			assembler := functioncall.NewAssembler()

			args := functioncall.Args{
				Function:      "myFunc",
				Order:         []string{"arg"},
				ArgumentTypes: map[string]string{"arg": tt.propType},
				Arguments:     map[string]any{"arg": tt.value},
			}

			// Act
			result, err := assembler.Assemble(args)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestAssembler_Assemble_StringWithDoubleQuotes_HappyPath(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function:      "disp",
		Order:         []string{"text"},
		ArgumentTypes: map[string]string{"text": "string"},
		Arguments:     map[string]any{"text": `say "hi"`},
	}

	// Act
	result, err := assembler.Assemble(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, `disp("say ""hi""")`, result)
}

func TestAssembler_Assemble_StringWithControlCharacters_HappyPath(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function:      "disp",
		Order:         []string{"text"},
		ArgumentTypes: map[string]string{"text": "string"},
		Arguments:     map[string]any{"text": "line1\nline2\ttab\r"},
	}

	// Act
	result, err := assembler.Assemble(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, `disp("line1line2tab")`, result)
}

func TestAssembler_Assemble_MultipleArgsInOrder_HappyPath(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function: "myFunc",
		Order:    []string{"name", "count", "flag"},
		ArgumentTypes: map[string]string{
			"name":  "string",
			"count": "number",
			"flag":  "boolean",
		},
		Arguments: map[string]any{
			"flag":  true,
			"name":  "test",
			"count": float64(42),
		},
	}

	// Act
	result, err := assembler.Assemble(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, `myFunc("test", 42, true)`, result)
}

func TestAssembler_Assemble_NoArguments_HappyPath(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function:      "version",
		Order:         []string{},
		ArgumentTypes: map[string]string{},
		Arguments:     map[string]any{},
	}

	// Act
	result, err := assembler.Assemble(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "version()", result)
}

func TestAssembler_Assemble_InvalidInteger_ReturnsError(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"non-whole integer", 3.5},
		{"too large integer", float64(math.MaxInt64) * 2},
		{"too small integer", float64(math.MinInt64) * 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			assembler := functioncall.NewAssembler()

			args := functioncall.Args{
				Function:      "myFunc",
				Order:         []string{"n"},
				ArgumentTypes: map[string]string{"n": "integer"},
				Arguments:     map[string]any{"n": tt.value},
			}

			// Act
			_, err := assembler.Assemble(args)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), "out of range or not a whole number")
		})
	}
}

func TestAssembler_Assemble_UnsupportedType_ReturnsError(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function:      "myFunc",
		Order:         []string{"data"},
		ArgumentTypes: map[string]string{"data": "array"},
		Arguments:     map[string]any{"data": []any{"a", "b"}},
	}

	// Act
	_, err := assembler.Assemble(args)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestAssembler_Assemble_TypeMismatch_ReturnsError(t *testing.T) {
	tests := []struct {
		name     string
		propType string
		value    any
	}{
		{"bool passed as string", "string", true},
		{"string passed as number", "number", "hello"},
		{"string passed as integer", "integer", "hello"},
		{"string passed as boolean", "boolean", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			assembler := functioncall.NewAssembler()

			args := functioncall.Args{
				Function:      "myFunc",
				Order:         []string{"arg"},
				ArgumentTypes: map[string]string{"arg": tt.propType},
				Arguments:     map[string]any{"arg": tt.value},
			}

			// Act
			_, err := assembler.Assemble(args)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), "expected")
		})
	}
}

func TestAssembler_Assemble_NilArgumentTypes_ReturnsError(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function:  "myFunc",
		Order:     []string{"arg"},
		Arguments: map[string]any{"arg": "value"},
	}

	// Act
	_, err := assembler.Assemble(args)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "argument type for \"arg\" not defined")
}

func TestAssembler_Assemble_MissingArgumentType_ReturnsError(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{},
		Arguments:     map[string]any{"n": float64(5)},
	}

	// Act
	_, err := assembler.Assemble(args)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "argument type for \"n\" not defined")
}

func TestAssembler_Assemble_UnexpectedArgument_ReturnsError(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{"n": float64(5), "extra": "hello"},
	}

	// Act
	_, err := assembler.Assemble(args)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected argument \"extra\" not in order")
}

func TestAssembler_Assemble_UnexpectedArgumentType_ReturnsError(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number", "extra": "string"},
		Arguments:     map[string]any{"n": float64(5)},
	}

	// Act
	_, err := assembler.Assemble(args)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected argument type \"extra\" not in order")
}

func TestAssembler_Assemble_MissingArgument_ReturnsError(t *testing.T) {
	// Arrange
	assembler := functioncall.NewAssembler()

	args := functioncall.Args{
		Function:      "magic",
		Order:         []string{"n"},
		ArgumentTypes: map[string]string{"n": "number"},
		Arguments:     map[string]any{},
	}

	// Act
	_, err := assembler.Assemble(args)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "n")
}
