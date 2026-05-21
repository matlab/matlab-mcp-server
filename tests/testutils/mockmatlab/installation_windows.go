// Copyright 2026 The MathWorks, Inc.
//go:build windows

package mockmatlab

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/config"
)

func buildPlatformSpecificBinaries(moduleDir, binDir string) error {
	archDir := filepath.Join(binDir, config.ArchFolder)
	if err := os.MkdirAll(archDir, 0o700); err != nil {
		return fmt.Errorf("failed to create arch directory: %w", err)
	}

	binaryPath := filepath.Join(archDir, config.ArchSpecificExeName)
	cmd := exec.Command("go", "build", "-o", binaryPath, "./tests/testutils/mockmatlab/executable/main") //nolint:gosec // Trusted test path
	cmd.Dir = moduleDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build mock MATLAB: %s: %w", string(output), err)
	}

	// The MATLAB locator discovers installations by finding "matlab.exe" on PATH
	// in the bin/ directory. Create a placeholder so discovery succeeds.
	placeholder := filepath.Join(binDir, config.MATLABExeName)
	if err := os.WriteFile(placeholder, nil, 0o600); err != nil {
		return fmt.Errorf("failed to create discovery placeholder: %w", err)
	}

	return nil
}

func mockMATLABBinaryPath(matlabRoot string) string {
	return filepath.Join(matlabRoot, "bin", config.ArchFolder, config.ArchSpecificExeName)
}
