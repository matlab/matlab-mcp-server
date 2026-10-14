// Copyright 2026 The MathWorks, Inc.

//go:build !windows

package externalmatlab

import (
	"os/exec"
	"syscall"
)

func startupArgs(startupDir string) []string {
	return []string{"-nosplash", "-softwareopengl", "-nodesktop", "-sd", startupDir}
}

func configureProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func terminate(cmd *exec.Cmd) {
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	_ = cmd.Process.Kill()
}
