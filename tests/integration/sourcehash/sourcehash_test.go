// Copyright 2026 The MathWorks, Inc.

package sourcehash_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/sourcehash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeHash_HappyPath(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "hello.txt", "hello world")

	// Act
	hash, err := sourcehash.ComputeHash(repoDir)

	// Assert
	require.NoError(t, err)
	assert.Len(t, hash, 64) // SHA-256 hex string
}

func TestComputeHash_Deterministic(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "hello.txt", "hello world")

	// Act
	hash1, err1 := sourcehash.ComputeHash(repoDir)
	hash2, err2 := sourcehash.ComputeHash(repoDir)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, hash1, hash2)
}

func TestComputeHash_IncludesSubdirectories(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	require.NoError(t, os.MkdirAll(filepath.Join(repoDir, "sub", "dir"), 0o750))
	writeFileAndCommit(t, repoDir, filepath.Join("sub", "dir", "nested.txt"), "nested content")

	// Act
	hash, err := sourcehash.ComputeHash(repoDir)

	// Assert
	require.NoError(t, err)
	assert.Len(t, hash, 64)
}

func TestComputeHash_ChangesWhenContentChanges(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "hello.txt", "hello world")

	hashBefore, err := sourcehash.ComputeHash(repoDir)
	require.NoError(t, err)

	writeFileAndCommit(t, repoDir, "hello.txt", "goodbye world")

	// Act
	hashAfter, err := sourcehash.ComputeHash(repoDir)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, hashBefore, hashAfter)
}

func TestComputeHash_ChangesWhenFileRenamed(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "original.txt", "content")

	hashBefore, err := sourcehash.ComputeHash(repoDir)
	require.NoError(t, err)

	renameAndCommit(t, repoDir, "original.txt", "renamed.txt")

	// Act
	hashAfter, err := sourcehash.ComputeHash(repoDir)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, hashBefore, hashAfter)
}

func TestComputeHash_ChangesWhenFileAdded(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "first.txt", "first")

	hashBefore, err := sourcehash.ComputeHash(repoDir)
	require.NoError(t, err)

	writeFileAndCommit(t, repoDir, "second.txt", "second")

	// Act
	hashAfter, err := sourcehash.ComputeHash(repoDir)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, hashBefore, hashAfter)
}

func TestComputeHash_ChangesWhenFileRemoved(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "first.txt", "first")
	writeFileAndCommit(t, repoDir, "second.txt", "second")

	hashBefore, err := sourcehash.ComputeHash(repoDir)
	require.NoError(t, err)

	removeAndCommit(t, repoDir, "second.txt")

	// Act
	hashAfter, err := sourcehash.ComputeHash(repoDir)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, hashBefore, hashAfter)
}

func TestWrite_HappyPath(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "hello.txt", "hello world")

	hashFile := filepath.Join(t.TempDir(), "hash.txt")

	// Act
	err := sourcehash.Write(hashFile, repoDir)

	// Assert
	require.NoError(t, err)
	require.NoError(t, sourcehash.Check(hashFile, repoDir))
}

func TestCheck_HappyPath(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "hello.txt", "hello world")

	hashFile := filepath.Join(t.TempDir(), "hash.txt")
	require.NoError(t, sourcehash.Write(hashFile, repoDir))

	// Act
	err := sourcehash.Check(hashFile, repoDir)

	// Assert
	require.NoError(t, err)
}

func TestCheck_MismatchError(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "hello.txt", "hello world")

	hashFile := filepath.Join(t.TempDir(), "hash.txt")
	require.NoError(t, sourcehash.Write(hashFile, repoDir))

	writeFileAndCommit(t, repoDir, "hello.txt", "changed content")

	// Act
	err := sourcehash.Check(hashFile, repoDir)

	// Assert
	require.Error(t, err)
	var mismatch *sourcehash.MismatchError
	require.ErrorAs(t, err, &mismatch)
	assert.NotEmpty(t, mismatch.Stored)
	assert.NotEmpty(t, mismatch.Current)
	assert.NotEqual(t, mismatch.Stored, mismatch.Current)
}

func TestCheck_HashFileNotFound(t *testing.T) {
	// Arrange
	repoDir := initGitRepo(t)
	writeFileAndCommit(t, repoDir, "hello.txt", "hello world")

	hashFile := filepath.Join(t.TempDir(), "nonexistent.txt")

	// Act
	err := sourcehash.Check(hashFile, repoDir)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read hash file")
}

// --- test helpers ---

func initGitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@test.com")
	git(t, dir, "config", "user.name", "Test")
	return dir
}

func writeFileAndCommit(t *testing.T, repoDir, name, content string) {
	t.Helper()

	fullPath := filepath.Join(repoDir, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o750))
	require.NoError(t, os.WriteFile(fullPath, []byte(content), 0o600))
	git(t, repoDir, "add", name)
	git(t, repoDir, "commit", "-m", "add "+name)
}

func renameAndCommit(t *testing.T, repoDir, oldName, newName string) {
	t.Helper()

	git(t, repoDir, "mv", oldName, newName)
	git(t, repoDir, "commit", "-m", "rename "+oldName+" to "+newName)
}

func removeAndCommit(t *testing.T, repoDir, name string) {
	t.Helper()

	git(t, repoDir, "rm", name)
	git(t, repoDir, "commit", "-m", "remove "+name)
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, out)
}
