// Copyright 2026 The MathWorks, Inc.

package loader

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/loader/validator"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type toolsFile struct {
	Tools      []definition.Tool               `json:"tools"`
	Signatures map[string]definition.Signature `json:"signatures"`
}

type OSLayer interface {
	ReadFile(filePath string) ([]byte, error)
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type ToolValidator interface {
	Validate(toolDefinition definition.Tool, signatures map[string]definition.Signature) (definition.ValidatedTool, error)
}

type Loader struct {
	osLayer       OSLayer
	loggerFactory LoggerFactory
	toolValidator ToolValidator
}

func NewLoader(osLayer OSLayer, loggerFactory LoggerFactory, toolValidator ToolValidator) *Loader {
	return &Loader{
		osLayer:       osLayer,
		loggerFactory: loggerFactory,
		toolValidator: toolValidator,
	}
}

func (l *Loader) Load(filePath string) ([]definition.ValidatedTool, messages.Error) {
	logger, loggerErr := l.loggerFactory.GetGlobalLogger()
	if loggerErr != nil {
		return nil, loggerErr
	}

	data, err := l.osLayer.ReadFile(filePath)
	if err != nil {
		logger.WithError(err).Error("Failed to read custom tools extension file")
		return nil, messages.New_StartupErrors_FailedToReadExtensionFile_Error(filePath)
	}

	var parsed toolsFile
	if err := json.Unmarshal(data, &parsed); err != nil {
		logger.WithError(err).Error("Failed to parse custom tools extension file")
		return nil, messages.New_StartupErrors_FailedToParseExtensionFile_Error(filePath)
	}

	validatedTools := make([]definition.ValidatedTool, 0, len(parsed.Tools))
	for _, toolDefinition := range parsed.Tools {
		validatedTool, err := l.toolValidator.Validate(toolDefinition, parsed.Signatures)
		if err != nil {
			logger.WithError(err).Error("Invalid custom tool definition")
			return nil, validationErrorToMessage(err, toolDefinition.Name, filePath)
		}

		if isDuplicateToolName(validatedTool.Definition().Name, validatedTools) {
			logger.WithError(fmt.Errorf("duplicate tool name %q", validatedTool.Definition().Name)).Error("Invalid custom tool definition")
			return nil, messages.New_StartupErrors_DuplicateToolName_Error(validatedTool.Definition().Name, filePath)
		}

		validatedTools = append(validatedTools, validatedTool)
	}

	logger.With("count", len(validatedTools)).Info("Loaded custom tools from extension file")
	return validatedTools, nil
}

func validationErrorToMessage(err error, toolName string, filePath string) messages.Error {
	switch {
	case errors.Is(err, validator.ErrInvalidInputSchema):
		return messages.New_StartupErrors_InvalidToolInputSchema_Error(toolName, filePath)
	case errors.Is(err, validator.ErrSignatureNotFound):
		return messages.New_StartupErrors_MissingToolSignature_Error(toolName, filePath)
	case errors.Is(err, validator.ErrInvalidSignature):
		return messages.New_StartupErrors_InvalidToolSignature_Error(toolName, filePath)
	case errors.Is(err, validator.ErrInvalidToolDefinition):
		return messages.New_StartupErrors_InvalidToolDefinition_Error(filePath)
	default:
		return messages.New_StartupErrors_InvalidToolDefinition_Error(filePath)
	}
}

func isDuplicateToolName(name string, existing []definition.ValidatedTool) bool {
	for _, vt := range existing {
		if vt.Definition().Name == name {
			return true
		}
	}
	return false
}
