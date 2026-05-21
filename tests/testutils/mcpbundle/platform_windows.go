// Copyright 2026 The MathWorks, Inc.

//go:build windows

package mcpbundle

import (
	"context"
	"os/exec"
)

const launcherFilename = "launch-matlab-mcp.cmd"
const pathWithSpaces = `C:\Program Files\MATLAB`

func execLauncherCommand(ctx context.Context, launcherPath string, args ...string) *exec.Cmd {
	cmdArgs := append([]string{"/c", launcherPath}, args...)
	return exec.CommandContext(ctx, "cmd", cmdArgs...) //nolint:gosec // Trusted test path
}
