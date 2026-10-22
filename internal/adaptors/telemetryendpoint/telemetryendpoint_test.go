// Copyright 2026 The MathWorks, Inc.

package telemetryendpoint_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetryendpoint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_HappyPath(t *testing.T) {
	// Arrange
	outputPath := filepath.Join(t.TempDir(), "endpoint_generated.go")
	const endpoint = "https://telemetry.example.test/v1/metrics"

	// Act
	err := telemetryendpoint.Generate(outputPath, endpoint)

	// Assert
	require.NoError(t, err)

	source, err := os.ReadFile(outputPath) //nolint:gosec // test controls the path
	require.NoError(t, err)

	content := string(source)
	assert.Contains(t, content, "//go:build release")
	assert.Contains(t, content, "package defaultparameters")
	assert.Contains(t, content, "const defaultTelemetryCollectorEndpoint = "+strconv.Quote(endpoint))
}

func TestGenerate_EmptyEndpointReturnsError(t *testing.T) {
	// Arrange
	outputPath := filepath.Join(t.TempDir(), "endpoint_generated.go")

	// Act
	err := telemetryendpoint.Generate(outputPath, "")

	// Assert
	require.ErrorIs(t, err, telemetryendpoint.ErrEmptyEndpoint)
	assert.NoFileExists(t, outputPath)
}

func TestGenerate_InvalidEndpointIsNotDisclosed(t *testing.T) {
	// Arrange
	testCases := []struct {
		name          string
		endpoint      string
		expectedError error
	}{
		{
			name:          "missing scheme",
			endpoint:      "sensitive-collector.invalid/v1/metrics",
			expectedError: telemetryendpoint.ErrUnsupportedScheme,
		},
		{
			name:          "unsupported scheme",
			endpoint:      "ftp://sensitive-collector.invalid/v1/metrics",
			expectedError: telemetryendpoint.ErrUnsupportedScheme,
		},
		{
			name:          "missing host",
			endpoint:      "https:///v1/metrics",
			expectedError: telemetryendpoint.ErrMissingHost,
		},
		{
			name:          "malformed",
			endpoint:      "https://sensitive-collector.invalid/v1/%zz",
			expectedError: telemetryendpoint.ErrInvalidEndpoint,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			outputPath := filepath.Join(t.TempDir(), "endpoint_generated.go")

			// Act
			err := telemetryendpoint.Generate(outputPath, testCase.endpoint)

			// Assert
			require.ErrorIs(t, err, testCase.expectedError)
			assert.NotContains(t, err.Error(), testCase.endpoint)
		})
	}
}
