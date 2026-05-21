// Copyright 2026 The MathWorks, Inc.

package fakematlab

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/tests/testconfig"
)

// Placeholder represents a non-functional MATLAB executable placeholder.
// Use this when tests only need a file to exist at the expected path.
type Placeholder struct {
	path string
}

// NewPlaceholder creates a fake MATLAB executable placeholder in the specified directory.
// The placeholder is a non-functional file that satisfies path detection checks.
//
// The executable is created at matlabRoot/bin/<matlab-exe-name>.
func NewPlaceholder(matlabRoot string) (*Placeholder, error) {
	matlabDir := filepath.Join(matlabRoot, "bin")
	if err := os.MkdirAll(matlabDir, 0o700); err != nil {
		return nil, fmt.Errorf("failed to create matlab directory: %w", err)
	}

	matlabPath := filepath.Join(matlabDir, testconfig.MATLABExeName)
	if err := os.WriteFile(matlabPath, []byte("fake matlab"), 0o700); err != nil { //nolint:gosec // Test file creation
		return nil, fmt.Errorf("failed to write fake matlab executable: %w", err)
	}

	return &Placeholder{path: matlabPath}, nil
}

// Path returns the full path to the fake MATLAB executable.
func (p *Placeholder) Path() string {
	return p.path
}

// Dir returns the directory containing the MATLAB executable (the bin directory).
func (p *Placeholder) Dir() string {
	return filepath.Dir(p.path)
}
