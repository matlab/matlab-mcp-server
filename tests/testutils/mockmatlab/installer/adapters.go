// Copyright 2026 The MathWorks, Inc.

package installer

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type GoListModuleRootFinder struct{}

func (GoListModuleRootFinder) FindModuleRoot() (string, error) {
	output, err := commandOutput("go", "list", "-m", "-json")
	if err != nil {
		return "", fmt.Errorf("failed to run 'go list -m -json': %w", err)
	}

	var mod struct{ Dir string }
	if err := json.Unmarshal(output, &mod); err != nil {
		return "", fmt.Errorf("failed to parse module info: %w", err)
	}

	return mod.Dir, nil
}

type BinaryBuilderFunc func(moduleDir, binDir string) error

func (f BinaryBuilderFunc) BuildPlatformSpecificBinaries(moduleDir, binDir string) error {
	return f(moduleDir, binDir)
}

func commandOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}
