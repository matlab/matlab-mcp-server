// Copyright 2026 The MathWorks, Inc.

//go:build !windows

package externalmatlab

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	standInScript     = "sleep 300 & echo $! > '%s'; wait"
	grandchildPidFile = "grandchild.pid"

	grandchildTimeout  = 10 * time.Second
	grandchildInterval = 50 * time.Millisecond
)

func TestTerminate_StopsTheProcessesThatMATLABStarted(t *testing.T) {
	pidFile := filepath.Join(t.TempDir(), grandchildPidFile)

	//nolint:gosec // The script is a constant, and the path comes from t.TempDir.
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf(standInScript, pidFile))
	configureProcessGroup(cmd)
	require.NoError(t, cmd.Start(), "should start the stand-in for MATLAB")

	grandchild := waitForGrandchild(t, pidFile)
	t.Cleanup(func() { _ = syscall.Kill(grandchild, syscall.SIGKILL) })

	terminate(cmd)
	_ = cmd.Wait()

	require.Eventually(t, func() bool {
		return syscall.Kill(grandchild, syscall.Signal(0)) != nil
	}, grandchildTimeout, grandchildInterval,
		"MATLAB starts Xvfb and fluxbox on a headless agent, so terminate must stop the whole group")
}

func waitForGrandchild(t *testing.T, pidFile string) int {
	t.Helper()

	var pid int

	require.Eventually(t, func() bool {
		content, err := os.ReadFile(pidFile) //nolint:gosec // Trusted test path
		if err != nil {
			return false
		}

		pid, err = strconv.Atoi(strings.TrimSpace(string(content)))

		return err == nil && pid > 0
	}, grandchildTimeout, grandchildInterval, "the stand-in should report the process it started")

	return pid
}
