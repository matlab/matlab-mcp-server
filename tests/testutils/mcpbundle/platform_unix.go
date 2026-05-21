// Copyright 2026 The MathWorks, Inc.

//go:build !windows

package mcpbundle

import (
	"context"
	"os/exec"
)

const launcherFilename = "launch-matlab-mcp.sh"
const pathWithSpaces = "/opt/my matlab/R2025b"

func execLauncherCommand(ctx context.Context, launcherPath string, args ...string) *exec.Cmd {
	cmdArgs := append([]string{launcherPath}, args...)
	return exec.CommandContext(ctx, "bash", cmdArgs...) //nolint:gosec // Trusted test path
}
