// Copyright 2026 The MathWorks, Inc.

package installer_test

import (
	"errors"
	"os"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/installer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoListModuleRootFinder_FindModuleRoot_ReturnsExistingDirectory(t *testing.T) {
	// Arrange
	finder := installer.GoListModuleRootFinder{}

	// Act
	moduleDir, err := finder.FindModuleRoot()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, moduleDir)
	info, statErr := os.Stat(moduleDir)
	require.NoError(t, statErr)
	assert.True(t, info.IsDir())
}

func TestBinaryBuilderFunc_BuildPlatformSpecificBinaries_DelegatesToWrappedFunction(t *testing.T) {
	// Arrange
	expectedErr := errors.New("build failed")
	called := false
	builder := installer.BinaryBuilderFunc(func(moduleDir, binDir string) error {
		called = true
		assert.Equal(t, "module-dir", moduleDir)
		assert.Equal(t, "bin-dir", binDir)
		return expectedErr
	})

	// Act
	err := builder.BuildPlatformSpecificBinaries("module-dir", "bin-dir")

	// Assert
	assert.True(t, called)
	assert.ErrorIs(t, err, expectedErr)
}
