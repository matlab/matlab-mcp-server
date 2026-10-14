// Copyright 2026 The MathWorks, Inc.

//go:build windows

package externalmatlab

import (
	"os/exec"
	"strconv"
)

// matlab.exe is a launcher, not MATLAB.
func startupArgs(startupDir string) []string {
	return []string{
		"-nosplash",
		"-softwareopengl",
		"-nodesktop",
		"-noDisplayDesktop",
		"-wait",
		"-log",
		"-sd",
		startupDir,
	}
}

// taskkill /T already stops the whole tree, so MATLAB does not need its own group.
func configureProcessGroup(*exec.Cmd) {}

// Kill alone stops only the launcher.
func terminate(cmd *exec.Cmd) {
	kill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid)) //nolint:gosec // Trusted test arguments
	_ = kill.Run()
	_ = cmd.Process.Kill()
}
