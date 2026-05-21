// Copyright 2026 The MathWorks, Inc.

package serverlogs_test

import (
	"testing"
	"testing/fstest"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/serverlogs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadErrorLogs_NoLogFiles(t *testing.T) {
	fsys := fstest.MapFS{}

	errorLogs, err := serverlogs.ReadErrorLogs(fsys)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no server log files found")
	assert.Nil(t, errorLogs)
}

func TestReadErrorLogs_NoErrors(t *testing.T) {
	fsys := fstest.MapFS{
		"server-abc123.log": &fstest.MapFile{
			Data: []byte(`{"level":"INFO","msg":"Server started"}
{"level":"DEBUG","msg":"Processing request"}`),
		},
	}

	errorLogs, err := serverlogs.ReadErrorLogs(fsys)

	require.NoError(t, err)
	assert.Empty(t, errorLogs)
}

func TestReadErrorLogs_WithErrors(t *testing.T) {
	fsys := fstest.MapFS{
		"server-abc123.log": &fstest.MapFile{
			Data: []byte(`{"level":"INFO","msg":"Server started"}
{"level":"ERROR","msg":"Something went wrong"}
{"level":"DEBUG","msg":"Processing request"}`),
		},
	}

	errorLogs, err := serverlogs.ReadErrorLogs(fsys)

	require.NoError(t, err)
	require.Len(t, errorLogs, 1)
	assert.Contains(t, errorLogs[0], "Something went wrong")
}

func TestReadErrorLogs_MultipleErrorsAcrossFiles(t *testing.T) {
	fsys := fstest.MapFS{
		"server-abc123.log": &fstest.MapFile{
			Data: []byte(`{"level":"ERROR","msg":"First error"}`),
		},
		"server-def456.log": &fstest.MapFile{
			Data: []byte(`{"level":"INFO","msg":"OK"}
{"level":"ERROR","msg":"Second error"}`),
		},
	}

	errorLogs, err := serverlogs.ReadErrorLogs(fsys)

	require.NoError(t, err)
	require.Len(t, errorLogs, 2)
}

func TestReadErrorLogs_IgnoresNonServerLogFiles(t *testing.T) {
	fsys := fstest.MapFS{
		"server-abc123.log": &fstest.MapFile{
			Data: []byte(`{"level":"INFO","msg":"OK"}`),
		},
		"watchdog-abc123.log": &fstest.MapFile{
			Data: []byte(`{"level":"ERROR","msg":"Watchdog error"}`),
		},
	}

	errorLogs, err := serverlogs.ReadErrorLogs(fsys)

	require.NoError(t, err)
	assert.Empty(t, errorLogs)
}

func TestReadErrorLogs_IgnoresShutdownEOFErrors(t *testing.T) {
	fsys := fstest.MapFS{
		"server-abc123.log": &fstest.MapFile{
			Data: []byte(`{"level":"ERROR","msg":"MCP server run method returned an unexpected error","error":"server is closing: EOF"}
{"level":"ERROR","msg":"Server failed with unexpected error","error":"server is closing: EOF"}
{"level":"ERROR","msg":"Genuine error"}`),
		},
	}

	errorLogs, err := serverlogs.ReadErrorLogs(fsys)

	require.NoError(t, err)
	require.Len(t, errorLogs, 1)
	assert.Contains(t, errorLogs[0], "Genuine error")
}

func TestReadErrorLogs_EmptyLogFile(t *testing.T) {
	fsys := fstest.MapFS{
		"server-abc123.log": &fstest.MapFile{
			Data: []byte(""),
		},
	}

	errorLogs, err := serverlogs.ReadErrorLogs(fsys)

	require.NoError(t, err)
	assert.Empty(t, errorLogs)
}
