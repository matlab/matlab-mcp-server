// Copyright 2026 The MathWorks, Inc.

package parser

import (
	"strconv"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

const internalErrorText = "Unimplemented parameter type"

func (p *Parser) parseEnvVars(specifiedArgs map[string]any, specifiedParameters map[string]struct{}) messages.Error {
	for _, parameter := range p.parameters {
		if !parameter.GetActive() {
			continue
		}

		envVarName := parameter.GetEnvVarName()
		if envVarName == "" {
			continue
		}

		val, ok := p.osLayer.LookupEnv(envVarName)
		if !ok {
			continue
		}

		var parsedVal any

		switch parameter.GetDefaultValue().(type) {
		case bool:
			boolVal, err := strconv.ParseBool(val)
			if err != nil {
				return messages.New_StartupErrors_BadValueForEnvVar_Error(val, envVarName)
			}
			parsedVal = boolVal
		case string:
			parsedVal = val
		case []string:
			parsedVal = []string{val}
		case time.Duration:
			durationVal, err := time.ParseDuration(val)
			if err != nil {
				return messages.New_StartupErrors_BadValueForEnvVar_Error(val, envVarName)
			}
			parsedVal = durationVal
		default:
			// If you hit this error, it means this switch is not implementing a supported type in `pkg/config`
			return messages.New_StartupErrors_ParseFailed_Error("\n", internalErrorText)
		}

		specifiedArgs[parameter.GetID()] = parsedVal
		specifiedParameters[parameter.GetID()] = struct{}{}
	}
	return nil
}
