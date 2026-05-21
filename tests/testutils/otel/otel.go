// Copyright 2026 The MathWorks, Inc.

package otel

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

const (
	containerImage     = "otel/opentelemetry-collector:latest"
	containerHTTPPort  = "4318"
	healthCheckPort    = "13133"
	telemetryFileName  = "telemetry.json"
	startupTimeout     = 30 * time.Second
	shutdownTimeout    = 10 * time.Second
	defaultPermissions = 0o777
)

type Collector struct {
	containerID    string
	hostPort       string
	healthHostPort string
	telemetryDir   string
	telemetryFile  string
	configFile     string
}

func StartCollector(t *testing.T, cfg collectorConfig) *Collector {
	t.Helper()

	hostPort, err := findFreePort()
	require.NoError(t, err, "failed to find free port for OTLP")

	healthHostPort, err := findFreePort()
	require.NoError(t, err, "failed to find free port for health check")

	testDir := t.TempDir()
	telemetryDir := filepath.Join(testDir, "telemetry")
	require.NoError(t, os.MkdirAll(telemetryDir, defaultPermissions), "failed to create telemetry dir")
	require.NoError(t, os.Chmod(telemetryDir, defaultPermissions), "failed to set telemetry dir permissions")

	configContent := cfg.generateConfig()
	configFile := filepath.Join(telemetryDir, "config.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(configContent), defaultPermissions), "failed to write config file")

	telemetryFile := filepath.Join(telemetryDir, telemetryFileName)

	containerID, err := startContainer(t, hostPort, healthHostPort, configFile, telemetryDir)
	require.NoError(t, err, "failed to start container")

	c := &Collector{
		containerID:    containerID,
		hostPort:       hostPort,
		healthHostPort: healthHostPort,
		telemetryDir:   telemetryDir,
		telemetryFile:  telemetryFile,
		configFile:     configFile,
	}

	err = c.waitForReady(t.Context())
	if err != nil {
		c.Stop(t)
	}
	require.NoError(t, err, "collector failed to become ready")

	return c
}

func (c *Collector) Endpoint() string {
	return fmt.Sprintf("http://localhost:%s", c.hostPort)
}

func (c *Collector) Stop(t *testing.T) {
	t.Helper()

	if c.containerID == "" {
		return
	}

	ctx, cancel := context.WithTimeout(t.Context(), shutdownTimeout)
	defer cancel()

	require.NoError(t, exec.CommandContext(ctx, "docker", "stop", c.containerID).Run())     //nolint:gosec // Trusted container ID from test setup
	require.NoError(t, exec.CommandContext(ctx, "docker", "rm", "-f", c.containerID).Run()) //nolint:gosec // Trusted container ID from test setup

	c.containerID = ""
}

func (c *Collector) WaitForTelemetry(t *testing.T, timeout time.Duration) {
	t.Helper()

	pollInterval := 100 * time.Millisecond
	require.Eventually(t, func() bool {
		info, err := os.Stat(c.telemetryFile)
		return err == nil && info.Size() > 0
	}, timeout, pollInterval, "telemetry file was not written: %s", c.telemetryFile)
}

func (c *Collector) ReadTelemetry(t *testing.T) (pmetric.Metrics, error) {
	data, err := os.ReadFile(c.telemetryFile)
	if err != nil {
		return pmetric.Metrics{}, err
	}

	unmarshaler := &pmetric.JSONUnmarshaler{}
	return unmarshaler.UnmarshalMetrics(data)
}

func (c *Collector) waitForReady(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, startupTimeout)
	defer cancel()

	healthURL := fmt.Sprintf("http://localhost:%s/", c.healthHostPort)

	_, err := retry.Retry(ctx, func() (struct{}, bool, error) {
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
		if reqErr != nil {
			return struct{}{}, false, reqErr
		}
		resp, httpErr := http.DefaultClient.Do(req)
		if httpErr != nil {
			return struct{}{}, false, nil
		}
		if err := resp.Body.Close(); err != nil {
			return struct{}{}, false, err
		}
		return struct{}{}, resp.StatusCode == http.StatusOK, nil
	}, retry.NewLinearRetryStrategy(100*time.Millisecond))

	if err != nil {
		logs := c.getLogs()
		return fmt.Errorf("collector failed to become ready: %w\nlogs: %s", err, logs)
	}

	return nil
}

func (c *Collector) getLogs() string {
	output, _ := exec.Command("docker", "logs", c.containerID).CombinedOutput() //nolint:gosec // Trusted container ID from test setup
	return string(output)
}

func findFreePort() (string, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	addr := listener.Addr().(*net.TCPAddr)

	return fmt.Sprintf("%d", addr.Port), listener.Close()
}

func startContainer(t *testing.T, hostPort, healthHostPort, configFile, telemetryDir string) (string, error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()

	containerPathToConfig := "/etc/otelcol/config.yaml"
	containerPathForFileExporter := "/tmp/telemetry"

	args := []string{
		"run", "-d",
		"-p", fmt.Sprintf("%s:%s", hostPort, containerHTTPPort),
		"-p", fmt.Sprintf("%s:%s", healthHostPort, healthCheckPort),
		"-v", fmt.Sprintf("%s:%s:ro", configFile, containerPathToConfig),
		"-v", fmt.Sprintf("%s:%s", telemetryDir, containerPathForFileExporter),
		containerImage,
	}

	cmd := exec.CommandContext(ctx, "docker", args...) //nolint:gosec // Trusted test arguments
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("docker run failed: %w, stderr: %s", err, string(exitErr.Stderr))
		}
		return "", fmt.Errorf("docker run failed: %w", err)
	}

	containerID := string(output)
	if len(containerID) > 12 {
		containerID = containerID[:12]
	}

	return containerID, nil
}
