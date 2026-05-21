// Copyright 2026 The MathWorks, Inc.

package installer_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	mocks "github.com/matlab/matlab-mcp-core-server/tests/mocks/testutils/mockmatlab/installer"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/installer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInstaller_BuildAndInstall_Success(t *testing.T) {
	// Arrange
	matlabRoot := filepath.Join("fake", "matlab", "root")
	binDir := filepath.Join(matlabRoot, "bin")
	moduleDir := filepath.Join("mock", "module")

	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().MkdirAll(binDir, os.FileMode(0o700)).Return(nil)
	fileSystem.EXPECT().WriteFile(
		filepath.Join(matlabRoot, "VersionInfo.xml"),
		mock.MatchedBy(func(content []byte) bool { return len(content) > 0 }),
		os.FileMode(0o600),
	).Return(nil)

	finder := mocks.NewMockModuleRootFinder(t)
	finder.EXPECT().FindModuleRoot().Return(moduleDir, nil)

	builder := mocks.NewMockBinaryBuilder(t)
	builder.EXPECT().BuildPlatformSpecificBinaries(moduleDir, binDir).Return(nil)

	inst := installer.New(fileSystem, finder, builder)

	// Act
	err := inst.BuildAndInstall(matlabRoot)

	// Assert
	require.NoError(t, err)
}

func TestInstaller_BuildAndInstall_MkdirAllFailure(t *testing.T) {
	// Arrange
	matlabRoot := filepath.Join("fake", "matlab", "root")
	expectedErr := errors.New("mkdir failed")

	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().MkdirAll(filepath.Join(matlabRoot, "bin"), os.FileMode(0o700)).Return(expectedErr)

	finder := mocks.NewMockModuleRootFinder(t)
	builder := mocks.NewMockBinaryBuilder(t)
	inst := installer.New(fileSystem, finder, builder)

	// Act
	err := inst.BuildAndInstall(matlabRoot)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create bin directory")
	assert.ErrorIs(t, err, expectedErr)
}

func TestInstaller_BuildAndInstall_FindModuleRootFailure(t *testing.T) {
	// Arrange
	matlabRoot := filepath.Join("fake", "matlab", "root")
	expectedErr := errors.New("find module root failed")

	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().MkdirAll(filepath.Join(matlabRoot, "bin"), os.FileMode(0o700)).Return(nil)

	finder := mocks.NewMockModuleRootFinder(t)
	finder.EXPECT().FindModuleRoot().Return("", expectedErr)

	builder := mocks.NewMockBinaryBuilder(t)
	inst := installer.New(fileSystem, finder, builder)

	// Act
	err := inst.BuildAndInstall(matlabRoot)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find module root")
	assert.ErrorIs(t, err, expectedErr)
}

func TestInstaller_BuildAndInstall_BuildBinariesFailure(t *testing.T) {
	// Arrange
	matlabRoot := filepath.Join("fake", "matlab", "root")
	binDir := filepath.Join(matlabRoot, "bin")
	moduleDir := filepath.Join("mock", "module")
	expectedErr := errors.New("build binaries failed")

	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().MkdirAll(binDir, os.FileMode(0o700)).Return(nil)

	finder := mocks.NewMockModuleRootFinder(t)
	finder.EXPECT().FindModuleRoot().Return(moduleDir, nil)

	builder := mocks.NewMockBinaryBuilder(t)
	builder.EXPECT().BuildPlatformSpecificBinaries(moduleDir, binDir).Return(expectedErr)

	inst := installer.New(fileSystem, finder, builder)

	// Act
	err := inst.BuildAndInstall(matlabRoot)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to build mock MATLAB binaries")
	assert.ErrorIs(t, err, expectedErr)
}

func TestInstaller_BuildAndInstall_WriteVersionInfoFailure(t *testing.T) {
	// Arrange
	matlabRoot := filepath.Join("fake", "matlab", "root")
	binDir := filepath.Join(matlabRoot, "bin")
	moduleDir := filepath.Join("mock", "module")
	expectedErr := errors.New("write version info failed")

	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().MkdirAll(binDir, os.FileMode(0o700)).Return(nil)
	fileSystem.EXPECT().WriteFile(
		filepath.Join(matlabRoot, "VersionInfo.xml"),
		mock.MatchedBy(func(content []byte) bool { return len(content) > 0 }),
		os.FileMode(0o600),
	).Return(expectedErr)

	finder := mocks.NewMockModuleRootFinder(t)
	finder.EXPECT().FindModuleRoot().Return(moduleDir, nil)

	builder := mocks.NewMockBinaryBuilder(t)
	builder.EXPECT().BuildPlatformSpecificBinaries(moduleDir, binDir).Return(nil)

	inst := installer.New(fileSystem, finder, builder)

	// Act
	err := inst.BuildAndInstall(matlabRoot)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write VersionInfo.xml")
	assert.ErrorIs(t, err, expectedErr)
}
