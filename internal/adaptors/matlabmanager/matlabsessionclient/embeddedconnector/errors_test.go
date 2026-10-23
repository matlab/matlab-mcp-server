// Copyright 2026 The MathWorks, Inc.

package embeddedconnector_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/stretchr/testify/require"
)

func TestMATLABExceptionsError_LogValue_LogsOnlyIdentifiersAndCount(t *testing.T) {
	// Arrange
	target := embeddedconnector.MATLABExceptionsError{
		ConsoleOutput: "patientName = 'Jane Doe'",
		Exceptions: []embeddedconnector.MATLABException{
			{
				Identifier: "MATLAB:UndefinedFunction",
				Message:    "Unrecognized function or variable 'patientRecord'.",
			},
		},
	}
	expectedLoggedError := map[string]any{
		"exception-count": float64(1),
		"identifiers":     []any{"MATLAB:UndefinedFunction"},
	}

	// Act
	loggedError := logErrorAsJSON(t, target)

	// Assert
	require.Equal(t, expectedLoggedError, loggedError)
}

func TestMATLABExceptionsError_LogValue_EmptyIdentifier_OmitsIdentifier(t *testing.T) {
	// Arrange
	target := embeddedconnector.MATLABExceptionsError{
		Exceptions: []embeddedconnector.MATLABException{
			{
				Message: "Error: C:\\Users\\jdoe\\secret.m not found",
			},
		},
	}
	expectedLoggedError := map[string]any{
		"exception-count": float64(1),
		"identifiers":     []any{},
	}

	// Act
	loggedError := logErrorAsJSON(t, target)

	// Assert
	require.Equal(t, expectedLoggedError, loggedError)
}

func logErrorAsJSON(t *testing.T, err error) any {
	t.Helper()
	var logs bytes.Buffer
	slog.New(slog.NewJSONHandler(&logs, nil)).With("error", err).Warn("handler returned an error")

	var entry struct {
		Error any `json:"error"`
	}
	require.NoError(t, json.Unmarshal(logs.Bytes(), &entry))
	return entry.Error
}

func TestMATLABExceptionsError_Error_HappyPath(t *testing.T) {
	// Arrange
	target := embeddedconnector.MATLABExceptionsError{
		ConsoleOutput: "before\n",
		Exceptions: []embeddedconnector.MATLABException{
			{
				Identifier: "MATLAB:UndefinedFunction",
				Message:    "Unrecognized function or variable 'undefinedFnXyz'.",
			},
		},
	}
	expectedMessage := "before\n\nmatlab error: (MATLAB:UndefinedFunction) Unrecognized function or variable 'undefinedFnXyz'."

	// Act
	message := target.Error()

	// Assert
	require.Equal(t, expectedMessage, message)
}

func TestMATLABExceptionsError_Error_EmptyIdentifier_OmitsIdentifier(t *testing.T) {
	// Arrange
	target := embeddedconnector.MATLABExceptionsError{
		ConsoleOutput: "Warning: careful",
		Exceptions: []embeddedconnector.MATLABException{
			{
				Message: "boom",
			},
		},
	}
	expectedMessage := "Warning: careful\n\nmatlab error: boom"

	// Act
	message := target.Error()

	// Assert
	require.Equal(t, expectedMessage, message)
}

func TestMATLABExceptionsError_Error_MultipleExceptions_SeparatedByBlankLineInOrder(t *testing.T) {
	// Arrange
	target := embeddedconnector.MATLABExceptionsError{
		Exceptions: []embeddedconnector.MATLABException{
			{
				Identifier: "first:id",
				Message:    "first failure",
			},
			{
				Identifier: "second:id",
				Message:    "second failure",
			},
		},
	}
	expectedMessage := "matlab error: (first:id) first failure\n\nmatlab error: (second:id) second failure"

	// Act
	message := target.Error()

	// Assert
	require.Equal(t, expectedMessage, message)
}
