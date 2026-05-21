// Copyright 2025-2026 The MathWorks, Inc.

package server_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/server"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServerFactory_NewServerOverUDS_HappyPath(t *testing.T) {
	// Arrange
	factory := newServerFactory()

	testDataDir, err := os.MkdirTemp("", "mcp_test") //nolint:usetesting // We can't use t.TempDir() here, as it sometimes creates path that are too long for socket paths
	require.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(testDataDir); err != nil {
			t.Logf("Failed to remove test data dir (may be locked on Windows): %v", err)
		}
	}()

	socketPath := filepath.Join(testDataDir, "test.sock")

	expectedStatusCode := http.StatusOK
	expectedFirstBody := "first hello world"
	expectedSecondBody := "second hello world"

	handlers := map[string]http.HandlerFunc{
		"GET /first": func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(expectedFirstBody))
		},
		"POST /second": func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(expectedSecondBody))
		},
	}

	server, err := factory.NewServerOverUDS(handlers)
	require.NoError(t, err)

	serverStopped := make(chan error, 1)
	go func() {
		serverStopped <- server.Serve(socketPath)
	}()
	defer func() {
		require.NoError(t, server.Shutdown(t.Context()))
	}()

	socketFileExists := make(chan error, 1)
	go func() {
		socketFileExists <- waitForSocketFile(t, socketPath)
	}()

	select {
	case err := <-serverStopped:
		t.Fatalf("Server stopped unexpectedly: %v", err)
	case err := <-socketFileExists:
		require.NoError(t, err)
	}

	client := newUDSClient(socketPath)

	// Act & Assert
	req, err := http.NewRequest("GET", "http://unix/first", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, expectedFirstBody, string(body))

	req, err = http.NewRequest("POST", "http://unix/second", nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, expectedSecondBody, string(body))
}

func newServerFactory() *server.Factory {
	application := integration.NewEmptyApplication()
	return application.HTTPServerFactory
}

func newUDSClient(socketPath string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
}

func waitForSocketFile(t *testing.T, socketPath string) error {
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
	defer cancel()

	_, err := retry.Retry(ctx, func() (struct{}, bool, error) {
		var zeroValues struct{}

		_, err := os.Stat(socketPath)

		if err == nil {
			return zeroValues, true, nil
		}

		if !os.IsNotExist(err) {
			return zeroValues, false, err
		}

		return zeroValues, false, nil
	}, retry.NewLinearRetryStrategy(100*time.Millisecond))

	return err
}
