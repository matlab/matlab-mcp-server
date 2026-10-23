// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector

import (
	"fmt"
	"log/slog"
	"strings"
)

type matlabError struct {
	message string
}

func newMATLABError(message string) matlabError {
	return matlabError{
		message: message,
	}
}

func (e matlabError) Error() string {
	return fmt.Sprintf("matlab error: %v", e.message)
}

type matlabException struct {
	Identifier string
	Message    string
}

type matlabExceptionsError struct {
	ConsoleOutput string
	Exceptions    []matlabException
}

func (e matlabExceptionsError) Error() string {
	blocks := make([]string, 0, len(e.Exceptions)+1)
	if consoleOutput := strings.TrimRight(e.ConsoleOutput, "\n"); consoleOutput != "" {
		blocks = append(blocks, consoleOutput)
	}

	for _, exception := range e.Exceptions {
		message := exception.Message
		if exception.Identifier != "" {
			message = fmt.Sprintf("(%s) %s", exception.Identifier, exception.Message)
		}
		blocks = append(blocks, newMATLABError(message).Error())
	}

	return strings.Join(blocks, "\n\n")
}

// Console output and error messages can carry user data, so they stay out of the logs.
func (e matlabExceptionsError) LogValue() slog.Value {
	identifiers := make([]string, 0, len(e.Exceptions))
	for _, exception := range e.Exceptions {
		if exception.Identifier != "" {
			identifiers = append(identifiers, exception.Identifier)
		}
	}

	return slog.GroupValue(
		slog.Int("exception-count", len(e.Exceptions)),
		slog.Any("identifiers", identifiers),
	)
}
