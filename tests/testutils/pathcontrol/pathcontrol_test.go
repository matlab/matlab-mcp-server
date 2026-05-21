// Copyright 2025-2026 The MathWorks, Inc.

package pathcontrol_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/fakematlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/pathcontrol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveFromPath_HappyPath(t *testing.T) {
	// Arrange
	path1 := filepath.Join("test", "path1")
	path2 := filepath.Join("test", "path2")
	expectedFinalPath := filepath.Join("original", "path")
	initialPath := strings.Join([]string{path1, path2, expectedFinalPath}, string(os.PathListSeparator))

	// Act
	result := pathcontrol.RemoveFromPath(initialPath, []string{path1, path2})

	// Assert
	actualElements := strings.Split(result, string(os.PathListSeparator))
	expectedElements := []string{expectedFinalPath}
	assert.ElementsMatch(t, expectedElements, actualElements)
}

func TestAddToPath_HappyPath(t *testing.T) {
	// Arrange
	path1 := filepath.Join("original", "path")
	path2 := filepath.Join("another", "original", "path")
	initialPath := strings.Join([]string{path1, path2}, string(os.PathListSeparator))

	newPath1 := filepath.Join("test", "path1")
	newPath2 := filepath.Join("test", "path2")

	// Act
	result := pathcontrol.AddToPath(initialPath, []string{newPath1, newPath2})

	// Assert
	actualElements := strings.Split(result, string(os.PathListSeparator))
	expectedElements := []string{path1, path2, newPath1, newPath2}
	assert.ElementsMatch(t, expectedElements, actualElements)
}

func TestRemoveAllMATLABsFromPath_HappyPath(t *testing.T) {
	// Arrange - Create a temporary directory with a fake MATLAB executable
	placeholder, err := fakematlab.NewPlaceholder(t.TempDir())
	require.NoError(t, err)
	fakeMatlabDir := placeholder.Dir()

	// Add some other paths to preserve
	otherPath1 := filepath.Join("usr", "local", "bin")
	otherPath2 := filepath.Join("opt", "tools", "bin")

	// Set PATH with fake MATLAB and other paths
	initialPath := strings.Join([]string{otherPath1, fakeMatlabDir, otherPath2}, string(os.PathListSeparator))

	// Act
	result := pathcontrol.RemoveAllMATLABsFromPath(initialPath)

	// Assert
	actualElements := strings.Split(result, string(os.PathListSeparator))
	expectedElements := []string{otherPath1, otherPath2}
	assert.ElementsMatch(t, expectedElements, actualElements, "PATH should only contain non-MATLAB paths")
}

func TestRemoveAllMATLABsFromPath_HappyPath_RemovesMultipleFakeMATLABs(t *testing.T) {
	// Arrange - Create multiple temporary directories with fake MATLAB executables
	placeholder1, err := fakematlab.NewPlaceholder(t.TempDir())
	require.NoError(t, err)
	fakeMatlabDir1 := placeholder1.Dir()

	// Create second fake MATLAB
	placeholder2, err := fakematlab.NewPlaceholder(t.TempDir())
	require.NoError(t, err)
	fakeMatlabDir2 := placeholder2.Dir()

	// Add some other paths to preserve
	otherPath := filepath.Join("usr", "local", "bin")

	// Set PATH with multiple fake MATLABs - first in PATH will be found first
	initialPath := strings.Join([]string{fakeMatlabDir1, fakeMatlabDir2, otherPath}, string(os.PathListSeparator))

	// Act
	result := pathcontrol.RemoveAllMATLABsFromPath(initialPath)

	// Assert
	actualElements := strings.Split(result, string(os.PathListSeparator))
	expectedElements := []string{otherPath}
	assert.ElementsMatch(t, expectedElements, actualElements, "PATH should only contain the non-MATLAB path")
}

func TestRemoveAllMATLABsFromPath_HappyPath_EmptyPath(t *testing.T) {
	// Arrange
	initialPath := ""

	// Act
	result := pathcontrol.RemoveAllMATLABsFromPath(initialPath)

	// Assert
	assert.Empty(t, result)
}

func TestRemoveFromPath_HandlesLeadingSeparator(t *testing.T) {
	// Arrange
	path1 := filepath.Join("test", "path1")
	path2 := filepath.Join("test", "path2")
	initialPath := strings.Join([]string{path1, path2}, string(os.PathListSeparator))

	// Act
	// Remove first path
	result := pathcontrol.RemoveFromPath(initialPath, []string{path1})

	// Assert
	actualElements := strings.Split(result, string(os.PathListSeparator))
	expectedElements := []string{path2}
	assert.ElementsMatch(t, expectedElements, actualElements)
}

func TestRemoveFromPath_HandlesDoubleSeparators(t *testing.T) {
	// Arrange
	path1 := filepath.Join("test", "path1")
	path2 := filepath.Join("test", "path2")
	path3 := filepath.Join("test", "path3")

	initialPath := strings.Join([]string{path1, path2, path3}, string(os.PathListSeparator))

	// Act
	result := pathcontrol.RemoveFromPath(initialPath, []string{path2})

	// Assert
	assert.NotContains(t, result, string(os.PathListSeparator)+string(os.PathListSeparator))
	actualElements := strings.Split(result, string(os.PathListSeparator))
	expectedElements := []string{path1, path3}
	assert.ElementsMatch(t, expectedElements, actualElements)
}

func TestRemoveAllMATLABsFromPath_NoMATLABs(t *testing.T) {
	// Arrange
	path1 := filepath.Join("usr", "local", "bin")
	path2 := filepath.Join("opt", "tools", "bin")
	initialPath := strings.Join([]string{path1, path2}, string(os.PathListSeparator))

	// Act
	result := pathcontrol.RemoveAllMATLABsFromPath(initialPath)

	// Assert
	// PATH should remain unchanged if no MATLAB is found
	actualElements := strings.Split(result, string(os.PathListSeparator))
	expectedElements := []string{path1, path2}
	assert.ElementsMatch(t, expectedElements, actualElements)
}

func TestUpdateEnvEntry_HappyPath(t *testing.T) {
	// Arrange
	oldPath := filepath.Join("old", "path")
	newPath := filepath.Join("new", "path")
	env := []string{"VAR1=val1", "PATH=" + oldPath, "VAR2=val2"}

	// Act
	newEnv := pathcontrol.UpdateEnvEntry(env, "PATH", newPath)

	// Assert
	assert.Contains(t, newEnv, "PATH="+newPath)
	assert.NotContains(t, newEnv, "PATH="+oldPath)
	assert.Contains(t, newEnv, "VAR1=val1")
	assert.Contains(t, newEnv, "VAR2=val2")
	assert.Len(t, newEnv, 3)
}

func TestUpdateEnvEntry_NoExistingKey(t *testing.T) {
	// Arrange
	env := []string{"VAR1=val1", "VAR2=val2"}
	newPath := filepath.Join("new", "path")

	// Act
	newEnv := pathcontrol.UpdateEnvEntry(env, "PATH", newPath)

	// Assert
	assert.Contains(t, newEnv, "PATH="+newPath)
	assert.Contains(t, newEnv, "VAR1=val1")
	assert.Contains(t, newEnv, "VAR2=val2")
	assert.Len(t, newEnv, 3)
}

func TestUpdateEnvEntry_MultipleSimilarVars(t *testing.T) {
	// Arrange
	oldPath := filepath.Join("old", "path")
	newPath := filepath.Join("new", "path")
	env := []string{"PATH_LIKE=something", "PATH=" + oldPath, "MY_PATH=elsewhere"}

	// Act
	newEnv := pathcontrol.UpdateEnvEntry(env, "PATH", newPath)

	// Assert
	assert.Contains(t, newEnv, "PATH="+newPath)
	assert.NotContains(t, newEnv, "PATH="+oldPath)
	assert.Contains(t, newEnv, "PATH_LIKE=something")
	assert.Contains(t, newEnv, "MY_PATH=elsewhere")
	assert.Len(t, newEnv, 3)
}

func TestUpdateEnvEntry_OtherVariable(t *testing.T) {
	// Arrange
	env := []string{"VAR1=val1", "VAR2=val2"}
	newValue := "newValue"

	// Act
	newEnv := pathcontrol.UpdateEnvEntry(env, "MY_VAR", newValue)

	// Assert
	assert.Contains(t, newEnv, "MY_VAR=newValue")
	assert.Contains(t, newEnv, "VAR1=val1")
	assert.Contains(t, newEnv, "VAR2=val2")
	assert.Len(t, newEnv, 3)
}
