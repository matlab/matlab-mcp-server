// Copyright 2026 The MathWorks, Inc.

package sessiondetails_test

import (
	"os"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/sessiondetails"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublish_WritesFileToExpectedPath(t *testing.T) {
	homeDir := t.TempDir()
	detailsJSON := `{"port":8080,"certificate":"/cert.pem","apiKey":"key","pid":1}`

	_, err := sessiondetails.Publish(homeDir, detailsJSON)
	require.NoError(t, err)

	expectedPath := sessiondetails.ResolveSessionDetailsPath(homeDir)
	content, err := os.ReadFile(expectedPath) //nolint:gosec // test reads from t.TempDir()
	require.NoError(t, err)
	assert.JSONEq(t, detailsJSON, string(content))
}

func TestPublish_FileHasSecurePermissions(t *testing.T) {
	homeDir := t.TempDir()
	detailsJSON := `{"port":8080}`

	_, err := sessiondetails.Publish(homeDir, detailsJSON)
	require.NoError(t, err)

	expectedPath := sessiondetails.ResolveSessionDetailsPath(homeDir)
	require.NoError(t, sessiondetails.AssertFileSecure(expectedPath))
}

func TestRemove_DeletesPublishedFile(t *testing.T) {
	homeDir := t.TempDir()
	detailsJSON := `{"port":8080}`

	_, err := sessiondetails.Publish(homeDir, detailsJSON)
	require.NoError(t, err)

	require.NoError(t, sessiondetails.Remove(homeDir))

	expectedPath := sessiondetails.ResolveSessionDetailsPath(homeDir)
	_, err = os.Stat(expectedPath)
	assert.True(t, os.IsNotExist(err))
}
