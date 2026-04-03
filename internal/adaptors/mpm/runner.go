// Copyright 2026 The MathWorks, Inc.

package mpm

import (
	"context"
	"fmt"
	"os/exec"
)

// Runner executes shell commands and returns their combined output.
type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

// Run executes the named command with the given arguments and returns combined stdout/stderr.
func (r *Runner) Run(ctx context.Context, name string, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec // name and args are constructed internally
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}
	return string(output), nil
}
