// Copyright 2025-2026 The MathWorks, Inc.

package matlablocator_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testconfig"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/fakematlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/matlablocator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetPath_HappyPath(t *testing.T) {
	// Arrange
	placeholder, err := fakematlab.NewPlaceholder(t.TempDir())
	require.NoError(t, err)

	t.Setenv("MCP_MATLAB_PATH", placeholder.Path())

	// Act
	matlabPath, err := matlablocator.GetPath()

	// Assert
	require.NoError(t, err)
	assert.True(t, filepath.IsAbs(matlabPath), "Expected absolute path, got: %s", matlabPath)
	assert.Equal(t, placeholder.Path(), matlabPath)
}

func Test_GetPath_NotSet(t *testing.T) {
	// Arrange

	t.Setenv("MCP_MATLAB_PATH", "")

	// Act
	_, err := matlablocator.GetPath()

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "MCP_MATLAB_PATH environment variable is empty")
}

func Test_GetPath_ErrorsWithRelativePath(t *testing.T) {
	// Arrange
	t.Setenv("MCP_MATLAB_PATH", "relative/path/to/matlab")

	// Act
	_, err := matlablocator.GetPath()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be an absolute path")
}

func Test_GetPath_ErrorsWithNonExistentAbsolutePath(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	nonExistentPath := filepath.Join(tempDir, "does_not_exist", testconfig.MATLABExeName)
	t.Setenv("MCP_MATLAB_PATH", nonExistentPath)

	// Act
	matlabPath, err := matlablocator.GetPath()

	// Assert
	require.Error(t, err)
	assert.Empty(t, matlabPath)
}
