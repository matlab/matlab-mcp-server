// Copyright 2026 The MathWorks, Inc.

package otel_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	entitiesmocks "github.com/matlab/matlab-mcp-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestErrorHandler_Handle_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	handler := otel.NewErrorHandler(mockLogger, "https://sensitive-collector.invalid/v1/metrics")

	// Act
	handler.Handle(assert.AnError)

	// Assert
	warnLogs := mockLogger.WarnLogs()
	require.Len(t, warnLogs, 1)
	for _, fields := range warnLogs {
		_, hasError := fields["error"]
		assert.False(t, hasError, "warning must not carry the error, which can contain the endpoint URL")
	}

	debugLogs := mockLogger.DebugLogs()
	require.Len(t, debugLogs, 1)
	for _, fields := range debugLogs {
		assert.EqualError(t, fields["error"].(error), assert.AnError.Error())
	}
}

func TestErrorHandler_Handle_RedactsEndpoint(t *testing.T) {
	const endpoint = "https://sensitive-collector.invalid/v1/metrics"
	const host = "sensitive-collector.invalid"

	testCases := []struct {
		name      string
		exportErr error
	}{
		{
			name:      "endpoint URL",
			exportErr: errors.New(`failed to upload metrics: Post "` + endpoint + `": request failed`),
		},
		{
			name:      "endpoint host",
			exportErr: errors.New(`dial tcp: lookup ` + host + `: no such host`),
		},
		{
			name:      "endpoint host in other case",
			exportErr: errors.New(`dial tcp: lookup SENSITIVE-COLLECTOR.INVALID: no such host`),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockLogger := testutils.NewInspectableLogger()

			handler := otel.NewErrorHandler(mockLogger, endpoint)

			// Act
			handler.Handle(testCase.exportErr)

			// Assert
			for _, logs := range []testutils.InspectableLogs{
				mockLogger.DebugLogs(),
				mockLogger.InfoLogs(),
				mockLogger.WarnLogs(),
				mockLogger.ErrorLogs(),
			} {
				requireLogsExcludeValue(t, logs, host)
			}

			debugLogs := mockLogger.DebugLogs()
			require.Len(t, debugLogs, 1)
			for _, fields := range debugLogs {
				assert.Contains(t, fmt.Sprint(fields["error"]), otel.RedactedEndpoint)
			}
		})
	}
}

func TestErrorHandler_Handle_NoEndpointLogsErrorUnchanged(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	handler := otel.NewErrorHandler(mockLogger, "")

	// Act
	handler.Handle(assert.AnError)

	// Assert
	debugLogs := mockLogger.DebugLogs()
	require.Len(t, debugLogs, 1)
	for _, fields := range debugLogs {
		assert.Equal(t, assert.AnError, fields["error"])
	}
}

func TestErrorHandler_Handle_SurfacesWarningOnce(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}

	anyWarningMessage := mock.Anything
	mockLogger.EXPECT().Warn(anyWarningMessage).Once()
	mockLogger.EXPECT().WithError(mock.Anything).Return(mockLogger).Times(2)
	mockLogger.EXPECT().Debug(mock.Anything).Times(2)

	handler := otel.NewErrorHandler(mockLogger, "https://sensitive-collector.invalid/v1/metrics")

	// Act
	handler.Handle(assert.AnError)
	handler.Handle(assert.AnError)

	// Assert
	mockLogger.AssertExpectations(t)
}

func requireLogsExcludeValue(t *testing.T, logs testutils.InspectableLogs, value string) {
	t.Helper()

	for message, fields := range logs {
		assert.NotContains(t, strings.ToLower(message), strings.ToLower(value))
		for _, field := range fields {
			assert.NotContains(t, strings.ToLower(fmt.Sprint(field)), strings.ToLower(value))
		}
	}
}
