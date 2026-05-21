// Copyright 2026 The MathWorks, Inc.

package parser

import (
	"errors"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/spf13/pflag"
)

func (p *Parser) setupFlags() {
	for flagName, parameter := range p.flagToParameter {
		if !parameter.GetActive() {
			continue
		}

		switch defaultValue := parameter.GetDefaultValue().(type) {
		case bool:
			p.flagSet.Bool(flagName, defaultValue, parameter.GetDescription())
		case string:
			p.flagSet.String(flagName, defaultValue, parameter.GetDescription())
		case []string:
			p.flagSet.StringArray(flagName, defaultValue, parameter.GetDescription())
		case time.Duration:
			p.flagSet.Duration(flagName, defaultValue, parameter.GetDescription())
		}
		if parameter.GetHiddenFlag() {
			_ = p.flagSet.MarkHidden(flagName) // Logically impossible to hit NotExistError
		}
	}
}

func (p *Parser) parseFlags(args []string, specifiedArgs map[string]any, specifiedParameters map[string]struct{}) messages.Error {
	err := p.flagSet.Parse(args)
	if err != nil {
		return p.convertToUserFacingError(err)
	}

	var messagesErr messages.Error

	p.flagSet.Visit(func(f *pflag.Flag) {
		parameter := p.flagToParameter[f.Name]

		var val any
		var err error

		switch parameter.GetDefaultValue().(type) {
		case bool:
			val, err = p.flagSet.GetBool(f.Name)
		case string:
			val, err = p.flagSet.GetString(f.Name)
		case []string:
			var flagValues []string
			flagValues, err = p.flagSet.GetStringArray(f.Name)
			if err == nil {
				// Append CLI values after any existing values (e.g. from env var)
				if existing, ok := specifiedArgs[parameter.GetID()].([]string); ok {
					val = append(existing, flagValues...)
				} else {
					val = flagValues
				}
			}
		case time.Duration:
			val, err = p.flagSet.GetDuration(f.Name)
		default:
			// If you hit this error, it means this switch is not implementing a supported type in `pkg/config`
			messagesErr = messages.New_StartupErrors_ParseFailed_Error("\n", internalErrorText)
			return
		}

		if err != nil {
			messagesErr = p.convertToUserFacingError(err)
			return
		}

		specifiedArgs[parameter.GetID()] = val
		specifiedParameters[parameter.GetID()] = struct{}{}
	})

	return messagesErr
}

func (p *Parser) convertToUserFacingError(err error) messages.Error {
	var notExistError *pflag.NotExistError
	var invalidSyntaxError *pflag.InvalidSyntaxError
	var invalidValueError *pflag.InvalidValueError
	var valueRequiredError *pflag.ValueRequiredError

	switch {
	case errors.As(err, &notExistError):
		return messages.New_StartupErrors_BadFlag_Error(notExistError.GetSpecifiedName(), "\n", p.usageText)
	case errors.As(err, &invalidSyntaxError):
		return messages.New_StartupErrors_BadSyntax_Error(invalidSyntaxError.GetSpecifiedFlag(), "\n", p.usageText)
	case errors.As(err, &invalidValueError):
		return messages.New_StartupErrors_BadValue_Error(invalidValueError.GetValue(), invalidValueError.GetFlag().Name)
	case errors.As(err, &valueRequiredError):
		return messages.New_StartupErrors_MissingValue_Error(valueRequiredError.GetSpecifiedName())
	default:
		return messages.New_StartupErrors_ParseFailed_Error("\n", p.usageText)
	}
}
