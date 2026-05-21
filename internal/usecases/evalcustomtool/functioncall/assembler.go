// Copyright 2026 The MathWorks, Inc.

package functioncall

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

type Args struct {
	Function      string
	Order         []string
	ArgumentTypes map[string]string
	Arguments     map[string]any
}

type Assembler struct{}

func NewAssembler() *Assembler {
	return &Assembler{}
}

func (a *Assembler) Assemble(args Args) (string, error) {
	if err := validateOrderedArgsExist(args); err != nil {
		return "", err
	}
	if err := validateNoExtraArgs(args); err != nil {
		return "", err
	}

	if len(args.Order) == 0 {
		return args.Function + "()", nil
	}

	formattedArgs := make([]string, 0, len(args.Order))
	for _, paramName := range args.Order {
		formatted, err := formatArgument(args.Arguments[paramName], args.ArgumentTypes[paramName])
		if err != nil {
			return "", fmt.Errorf("failed to format argument %q: %w", paramName, err)
		}
		formattedArgs = append(formattedArgs, formatted)
	}

	return args.Function + "(" + strings.Join(formattedArgs, ", ") + ")", nil
}

func validateOrderedArgsExist(args Args) error {
	for _, paramName := range args.Order {
		if _, exists := args.ArgumentTypes[paramName]; !exists {
			return fmt.Errorf("argument type for %q not defined", paramName)
		}
		if _, exists := args.Arguments[paramName]; !exists {
			return fmt.Errorf("argument %q not provided", paramName)
		}
	}
	return nil
}

func validateNoExtraArgs(args Args) error {
	orderSet := make(map[string]struct{}, len(args.Order))
	for _, paramName := range args.Order {
		orderSet[paramName] = struct{}{}
	}
	for typeName := range args.ArgumentTypes {
		if _, exists := orderSet[typeName]; !exists {
			return fmt.Errorf("unexpected argument type %q not in order", typeName)
		}
	}
	for argName := range args.Arguments {
		if _, exists := orderSet[argName]; !exists {
			return fmt.Errorf("unexpected argument %q not in order", argName)
		}
	}
	return nil
}

func formatArgument(value any, argType string) (string, error) {
	switch argType {
	case "string":
		s, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("expected string, got %T", value)
		}
		s = stripControlCharacters(s)
		escaped := strings.ReplaceAll(s, `"`, `""`)
		return fmt.Sprintf(`"%s"`, escaped), nil
	case "number":
		v, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("expected number, got %T", value)
		}
		return fmt.Sprintf("%g", v), nil
	case "integer":
		v, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("expected number, got %T", value)
		}
		if v != math.Trunc(v) || v > math.MaxInt64 || v < math.MinInt64 {
			return "", fmt.Errorf("integer value %g is out of range or not a whole number", v)
		}
		return fmt.Sprintf("%d", int64(v)), nil
	case "boolean":
		v, ok := value.(bool)
		if !ok {
			return "", fmt.Errorf("expected boolean, got %T", value)
		}
		if v {
			return "true", nil
		}
		return "false", nil
	default:
		return "", fmt.Errorf("unsupported type %q", argType)
	}
}

func stripControlCharacters(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}
