// Copyright 2026 The MathWorks, Inc.

package sessiondetails_test

import (
	"encoding/json"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/sessiondetails"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON_HappyPath(t *testing.T) {
	result, err := sessiondetails.MarshalJSON("8080", "/path/to/cert.pem", "test-api-key", 12345)
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(result), &parsed))

	assert.InDelta(t, float64(8080), parsed["port"], 0)
	assert.Equal(t, "/path/to/cert.pem", parsed["certificate"])
	assert.Equal(t, "test-api-key", parsed["apiKey"])
	assert.InDelta(t, float64(12345), parsed["pid"], 0)
}

func TestMarshalJSON_InvalidPort(t *testing.T) {
	_, err := sessiondetails.MarshalJSON("not-a-number", "/cert.pem", "key", 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid port")
}

func TestResolveSessionDetailsPath(t *testing.T) {
	path := sessiondetails.ResolveSessionDetailsPath("/home/testuser")
	assert.Contains(t, path, "v1")
	assert.Contains(t, path, "sessionDetails.json")
}
