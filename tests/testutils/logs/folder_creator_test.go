// Copyright 2026 The MathWorks, Inc.

package logs_test

import (
	"io/fs"
	"path/filepath"
	"testing"

	mocks "github.com/matlab/matlab-mcp-core-server/tests/mocks/testutils/logs"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFolderCreatorWithFileSystem_NilDependencyReturnsError(t *testing.T) {
	// Act
	_, err := logs.NewFolderCreatorWithFileSystem(nil)

	// Assert
	require.EqualError(t, err, "fileSystem must not be nil")
}

func TestFolderCreator_CreateTempLogFolder_HappyPath(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockDirectoryFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	prefix := "test-logs-"
	baseDir := "/tmp/base"
	expectedLogDir := filepath.Join(baseDir, "logs")

	mockFileSystem.EXPECT().MkdirTemp("", prefix).Return(baseDir, nil).Once()
	mockFileSystem.EXPECT().MkdirAll(expectedLogDir, fs.FileMode(0o750)).Return(nil).Once()

	creator, err := logs.NewFolderCreatorWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	returnedBaseDir, returnedLogDir, err := creator.CreateTempLogFolder(prefix)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, baseDir, returnedBaseDir)
	assert.Equal(t, expectedLogDir, returnedLogDir)
}

func TestFolderCreator_CreateTempLogFolder_MkdirTempError(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockDirectoryFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	prefix := "test-logs-"

	mockFileSystem.EXPECT().MkdirTemp("", prefix).Return("", assert.AnError).Once()

	creator, err := logs.NewFolderCreatorWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	_, _, err = creator.CreateTempLogFolder(prefix)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestFolderCreator_CreateTempLogFolder_MkdirAllError(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockDirectoryFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	prefix := "test-logs-"
	baseDir := "/tmp/base"
	expectedLogDir := filepath.Join(baseDir, "logs")

	mockFileSystem.EXPECT().MkdirTemp("", prefix).Return(baseDir, nil).Once()
	mockFileSystem.EXPECT().MkdirAll(expectedLogDir, fs.FileMode(0o750)).Return(assert.AnError).Once()

	creator, err := logs.NewFolderCreatorWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	_, _, err = creator.CreateTempLogFolder(prefix)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestFolderCreator_PrepareSessionCLIArgs_UsesProvidedLogFlags(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockDirectoryFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)

	creator, err := logs.NewFolderCreatorWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	inputArgs := []string{"--log-level=info", "--log-folder=/custom/logs", "--other=1"}

	// Act
	result, err := creator.PrepareSessionCLIArgs(inputArgs, "debug", "unused-")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, inputArgs, result.Args)
	assert.Equal(t, "/custom/logs", result.LogDir)
	assert.Empty(t, result.TempBaseDir)
}

func TestFolderCreator_PrepareSessionCLIArgs_AddsDefaultsAndCreatesFolder(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockDirectoryFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	prefix := "test-logs-"
	baseDir := "/tmp/base"
	expectedLogDir := filepath.Join(baseDir, "logs")

	mockFileSystem.EXPECT().MkdirTemp("", prefix).Return(baseDir, nil).Once()
	mockFileSystem.EXPECT().MkdirAll(expectedLogDir, fs.FileMode(0o750)).Return(nil).Once()

	creator, err := logs.NewFolderCreatorWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	inputArgs := []string{"--other=1"}

	// Act
	result, err := creator.PrepareSessionCLIArgs(inputArgs, "debug", prefix)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"--log-level=debug", "--log-folder=" + expectedLogDir, "--other=1"}, result.Args)
	assert.Equal(t, expectedLogDir, result.LogDir)
	assert.Equal(t, baseDir, result.TempBaseDir)
}

func TestFolderCreator_PrepareSessionCLIArgs_AddsMissingLogLevelOnly(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockDirectoryFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)

	creator, err := logs.NewFolderCreatorWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	inputArgs := []string{"--log-folder=/custom/logs", "--other=1"}

	// Act
	result, err := creator.PrepareSessionCLIArgs(inputArgs, "debug", "unused-")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"--log-level=debug", "--log-folder=/custom/logs", "--other=1"}, result.Args)
	assert.Equal(t, "/custom/logs", result.LogDir)
	assert.Empty(t, result.TempBaseDir)
}
