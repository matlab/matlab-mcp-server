// Copyright 2026 The MathWorks, Inc.

package externalmatlab

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/tests/testutils/matlabtempdir"
	"github.com/stretchr/testify/require"
)

// Must stay equal to MATLAB_SYSTEM_TEST_PROCESS_MARKER in the Makefile.
const startupDirName = "mcp-system-test-matlab"

const (
	outputFileName = "matlab_output.log"
	readyFileName  = "ready"

	commandSucceeded = "done"

	startupTimeout   = 3 * time.Minute
	commandTimeout   = 1 * time.Minute
	filePollInterval = 100 * time.Millisecond
)

const matlabAssetsDir = "assets/matlab"

//go:embed assets/matlab/*.m
var matlabAssets embed.FS

type Instance struct {
	dir          string
	commandDir   string
	exited       chan struct{}
	exitErr      error
	commandCount int
}

func Start(t *testing.T, matlabExe string, env []string, prepCode string) *Instance {
	t.Helper()

	dir := matlabtempdir.New(t, "external-matlab-")
	startupDir := filepath.Join(dir, startupDirName)
	require.NoError(t, os.MkdirAll(startupDir, 0o700), "should create MATLAB startup directory")
	writeMATLABAssets(t, startupDir)

	output := createOutputFile(t, dir)

	// Not CommandContext: t.Context() ends before the cleanup, so terminate would have no tree to stop.
	cmd := exec.Command(matlabExe, startupArgs(startupDir)...) //nolint:gosec // Trusted test path
	cmd.Dir = startupDir
	cmd.Env = env
	cmd.Stdout = output
	cmd.Stderr = output
	configureProcessGroup(cmd)
	// The open stdin pipe keeps MATLAB running until the cleanup.
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "should open MATLAB stdin pipe")
	require.NoError(t, cmd.Start(), "should start external MATLAB")

	matlab := &Instance{dir: dir, commandDir: startupDir, exited: make(chan struct{})}
	go func() {
		matlab.exitErr = cmd.Wait()
		close(matlab.exited)
	}()
	t.Cleanup(func() {
		_ = stdin.Close()
		terminate(cmd)
		<-matlab.exited
	})

	matlab.waitForFile(t, filepath.Join(startupDir, readyFileName), "MATLAB startup", startupTimeout)
	matlab.run(t, prepCode, "the preparation code")

	return matlab
}

func (m *Instance) Share(t *testing.T) {
	t.Helper()
	m.run(t, "shareMATLABSession();", "shareMATLABSession()")
}

// Commands travel in files, not stdin, so that the same channel works on Windows.
func (m *Instance) run(t *testing.T, code, description string) {
	t.Helper()

	m.commandCount++
	name := fmt.Sprintf("command_%d", m.commandCount)
	requestPath := filepath.Join(m.commandDir, name+".request")
	responsePath := filepath.Join(m.commandDir, name+".response")

	partialPath := requestPath + ".partial"
	require.NoError(t, os.WriteFile(partialPath, []byte(code), 0o600),
		"should stage the %s request for MATLAB", description)
	require.NoError(t, os.Rename(partialPath, requestPath),
		"should publish the %s request for MATLAB", description)

	m.waitForFile(t, responsePath, description, commandTimeout)

	response, err := os.ReadFile(responsePath) //nolint:gosec // Trusted test path
	require.NoError(t, err, "should read the %s response from MATLAB", description)
	require.Equal(t, commandSucceeded, strings.TrimSpace(string(response)),
		"%s should succeed in the external MATLAB\nOutput: %s", description, m.output())
}

func (m *Instance) AssertAlive(t *testing.T) {
	t.Helper()
	m.run(t, "pause(0);", "liveness check")
}

func (m *Instance) output() string {
	data, err := os.ReadFile(filepath.Join(m.dir, outputFileName))
	if err != nil {
		return ""
	}
	return string(data)
}

func (m *Instance) waitForFile(t *testing.T, path, description string, timeout time.Duration) {
	t.Helper()

	deadline := time.After(timeout)
	ticker := time.NewTicker(filePollInterval)
	defer ticker.Stop()

	for {
		if _, err := os.Stat(path); err == nil {
			return
		}

		select {
		case <-m.exited:
			require.Fail(t, "MATLAB exited before "+description,
				"exit error: %v\nOutput: %s", m.exitErr, m.output())
		case <-deadline:
			require.Fail(t, "timed out waiting for "+description,
				"waited %s\nOutput: %s", timeout, m.output())
		case <-ticker.C:
		}
	}
}

func createOutputFile(t *testing.T, dir string) *os.File {
	t.Helper()
	output, err := os.OpenFile( //nolint:gosec // Trusted test path
		filepath.Join(dir, outputFileName),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o600,
	)
	require.NoError(t, err, "should create MATLAB output file")
	t.Cleanup(func() { _ = output.Close() })
	return output
}

func writeMATLABAssets(t *testing.T, dir string) {
	t.Helper()

	entries, err := fs.ReadDir(matlabAssets, matlabAssetsDir)
	require.NoError(t, err, "should list the embedded MATLAB files")
	for _, entry := range entries {
		content, err := matlabAssets.ReadFile(matlabAssetsDir + "/" + entry.Name())
		require.NoError(t, err, "should read the embedded %s", entry.Name())
		require.NoError(t, os.WriteFile(filepath.Join(dir, entry.Name()), content, 0o600),
			"should write %s", entry.Name())
	}
}
