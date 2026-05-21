// Copyright 2026 The MathWorks, Inc.

package fakematlab_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testconfig"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/fakematlab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlaceholder_CreatesExecutableFile(t *testing.T) {
	// Arrange
	baseDir := t.TempDir()

	// Act
	placeholder, err := fakematlab.NewPlaceholder(baseDir)

	// Assert
	require.NoError(t, err)
	require.FileExists(t, placeholder.Path())
}

func TestNewPlaceholder_CreatesInBinSubdirectory(t *testing.T) {
	// Arrange
	baseDir := t.TempDir()

	// Act
	placeholder, err := fakematlab.NewPlaceholder(baseDir)

	// Assert
	require.NoError(t, err)
	expectedDir := filepath.Join(baseDir, "bin")
	assert.Equal(t, expectedDir, filepath.Dir(placeholder.Path()))
}

func TestNewPlaceholder_UsesCorrectExecutableName(t *testing.T) {
	// Arrange
	baseDir := t.TempDir()

	// Act
	placeholder, err := fakematlab.NewPlaceholder(baseDir)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, testconfig.MATLABExeName, filepath.Base(placeholder.Path()))
}

func TestPlaceholder_Dir_ReturnsBinDirectory(t *testing.T) {
	// Arrange
	baseDir := t.TempDir()
	placeholder, err := fakematlab.NewPlaceholder(baseDir)
	require.NoError(t, err)

	// Act
	dir := placeholder.Dir()

	// Assert
	assert.Equal(t, filepath.Join(baseDir, "bin"), dir)
}

func TestNewPlaceholder_FailsWithInvalidPath(t *testing.T) {
	// Arrange
	invalidPath := fmt.Sprintf("%s%c%s", filepath.Join("nonexistent", "path", "that", "cannot", "be", "created"), 0, "invalid")

	// Act
	_, err := fakematlab.NewPlaceholder(invalidPath)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create matlab directory")
}
