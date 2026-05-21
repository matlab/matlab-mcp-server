// Copyright 2026 The MathWorks, Inc.
//go:build !windows

package mockmatlab

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/config"
)

func buildPlatformSpecificBinaries(moduleDir, binDir string) error {
	binaryPath := filepath.Join(binDir, config.MATLABExeName)

	cmd := exec.Command("go", "build", "-o", binaryPath, "./tests/testutils/mockmatlab/executable/main") //nolint:gosec // Trusted test path
	cmd.Dir = moduleDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build mock MATLAB: %s: %w", string(output), err)
	}

	return nil
}

func mockMATLABBinaryPath(matlabRoot string) string {
	return filepath.Join(matlabRoot, "bin", config.MATLABExeName)
}
