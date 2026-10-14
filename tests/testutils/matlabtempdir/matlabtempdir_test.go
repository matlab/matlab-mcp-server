// Copyright 2026 The MathWorks, Inc.

package matlabtempdir_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-server/tests/testutils/matlabtempdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_CleanupRemovesReadOnlyContent(t *testing.T) {
	// Arrange
	var dir string
	t.Run("create read-only content", func(t *testing.T) {
		dir = matlabtempdir.New(t, "matlabtempdir-test-")
		readOnlyDir := filepath.Join(dir, "toolbox")
		require.NoError(t, os.Mkdir(readOnlyDir, 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(readOnlyDir, "function.m"), nil, 0o400))
		require.NoError(t, os.Chmod(readOnlyDir, 0o500)) //nolint:gosec // A directory needs the execute bit
	})

	// Assert
	assert.NoDirExists(t, dir)
}

func TestNew_CleanupKeepsSymlinkTargetPermissions(t *testing.T) {
	// Arrange
	target := t.TempDir()
	targetFile := filepath.Join(target, "prefs.txt")
	require.NoError(t, os.WriteFile(targetFile, nil, 0o400))
	before, err := os.Stat(targetFile)
	require.NoError(t, err)
	if err := os.Symlink(targetFile, filepath.Join(target, "probe")); err != nil {
		t.Skipf("cannot create a symlink: %v", err)
	}
	var dir string
	t.Run("link outside content from a read-only folder", func(t *testing.T) {
		dir = matlabtempdir.New(t, "matlabtempdir-test-")
		readOnlyDir := filepath.Join(dir, "toolbox")
		require.NoError(t, os.Mkdir(readOnlyDir, 0o700))
		require.NoError(t, os.Symlink(targetFile, filepath.Join(readOnlyDir, "link")))
		require.NoError(t, os.Chmod(readOnlyDir, 0o500)) //nolint:gosec // A directory needs the execute bit
	})

	// Assert
	assert.NoDirExists(t, dir)
	after, err := os.Stat(targetFile)
	require.NoError(t, err)
	assert.Equal(t, before.Mode().Perm(), after.Mode().Perm())
}
