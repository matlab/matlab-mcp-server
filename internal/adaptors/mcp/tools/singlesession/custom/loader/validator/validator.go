// Copyright 2026 The MathWorks, Inc.

package validator

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/definition"
)

var (
	ErrInvalidToolDefinition = errors.New("invalid tool definition")
	ErrInvalidInputSchema    = errors.New("invalid input schema")
	ErrSignatureNotFound     = errors.New("signature not found")
	ErrInvalidSignature      = errors.New("invalid signature")
)

// validMATLABFunctionName matches a valid MATLAB function name:
//   - Starts with an alphabetic character
//   - Followed by zero or more word characters
//   - Followed by zero or more dot-separated segments of word characters (each starting with an alphabetic character)
//   - Example: pkg.myFunc
var validMATLABFunctionName = regexp.MustCompile(`^[A-Za-z]\w*(\.[A-Za-z]\w*)*$`)

type validatedTool struct {
	definition definition.Tool
	signature  definition.Signature
}

func (v *validatedTool) Definition() definition.Tool     { return v.definition }
func (v *validatedTool) Signature() definition.Signature { return v.signature }

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(toolDefinition definition.Tool, signatures map[string]definition.Signature) (definition.ValidatedTool, error) {
	if err := validateToolDefinition(toolDefinition); err != nil {
		return nil, err
	}

	if err := validateInputSchema(toolDefinition.InputSchema); err != nil {
		return nil, err
	}

	sig, err := findSignature(toolDefinition.Name, signatures)
	if err != nil {
		return nil, err
	}

	if err := validateSignature(sig, toolDefinition.InputSchema); err != nil {
		return nil, err
	}

	return &validatedTool{
		definition: toolDefinition,
		signature:  sig,
	}, nil
}

func validateToolDefinition(toolDefinition definition.Tool) error {
	if toolDefinition.Name == "" {
		return fmt.Errorf("missing required field: name: %w", ErrInvalidToolDefinition)
	}
	if toolDefinition.Title == "" {
		return fmt.Errorf("missing required field: title: %w", ErrInvalidToolDefinition)
	}
	if toolDefinition.Description == "" {
		return fmt.Errorf("missing required field: description: %w", ErrInvalidToolDefinition)
	}
	return nil
}

func validateInputSchema(schema *jsonschema.Schema) error {
	if schema == nil {
		return fmt.Errorf("inputSchema is required: %w", ErrInvalidInputSchema)
	}

	if schema.Type != "object" {
		return fmt.Errorf("inputSchema type must be 'object', got '%s': %w", schema.Type, ErrInvalidInputSchema)
	}

	for propName, prop := range schema.Properties {
		if prop == nil {
			return fmt.Errorf("property %q is nil: %w", propName, ErrInvalidInputSchema)
		}
		if prop.Type == "" {
			return fmt.Errorf("property %q must have a 'type' field: %w", propName, ErrInvalidInputSchema)
		}
		if !isSupportedType(prop.Type) {
			return fmt.Errorf("property %q has unsupported type %q (supported: string, number, integer, boolean): %w", propName, prop.Type, ErrInvalidInputSchema)
		}
	}

	for _, reqName := range schema.Required {
		if _, exists := schema.Properties[reqName]; !exists {
			return fmt.Errorf("required property %q is not defined in properties: %w", reqName, ErrInvalidInputSchema)
		}
	}

	return nil
}

func isSupportedType(t string) bool {
	switch t {
	case "string", "number", "integer", "boolean":
		return true
	default:
		return false
	}
}

func findSignature(toolName string, signatures map[string]definition.Signature) (definition.Signature, error) {
	sig, ok := signatures[toolName]
	if !ok {
		return definition.Signature{}, ErrSignatureNotFound
	}
	return sig, nil
}

func validateSignature(sig definition.Signature, schema *jsonschema.Schema) error {
	if sig.Function == "" {
		return fmt.Errorf("signature must have a 'function' field: %w", ErrInvalidSignature)
	}

	if !validMATLABFunctionName.MatchString(sig.Function) {
		return fmt.Errorf("signature function %q is not a valid MATLAB function name: %w", sig.Function, ErrInvalidSignature)
	}

	if len(sig.Input.Order) > 0 && schema.Properties == nil {
		return fmt.Errorf("inputSchema properties not defined: %w", ErrInvalidSignature)
	}

	orderSet := make(map[string]struct{}, len(sig.Input.Order))
	for _, orderEntry := range sig.Input.Order {
		if _, duplicate := orderSet[orderEntry]; duplicate {
			return fmt.Errorf("duplicate entry %q in input.order: %w", orderEntry, ErrInvalidSignature)
		}
		orderSet[orderEntry] = struct{}{}

		if _, exists := schema.Properties[orderEntry]; !exists {
			return fmt.Errorf("input.order entry %q not found in inputSchema properties: %w", orderEntry, ErrInvalidSignature)
		}
	}

	for propName := range schema.Properties {
		if _, exists := orderSet[propName]; !exists {
			return fmt.Errorf("inputSchema property %q is not included in input.order: %w", propName, ErrInvalidSignature)
		}
	}

	return nil
}
