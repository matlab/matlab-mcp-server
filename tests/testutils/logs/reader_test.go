// Copyright 2026 The MathWorks, Inc.

package logs_test

import (
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"

	mocks "github.com/matlab/matlab-mcp-core-server/tests/mocks/testutils/logs"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReaderWithFileSystem_NilDependencyReturnsError(t *testing.T) {
	// Act
	_, err := logs.NewReaderWithFileSystem(nil)

	// Assert
	require.EqualError(t, err, "fileSystem must not be nil")
}

func TestReader_ReadCombined_HappyPath(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockLogFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	var logFS fs.FS = fstest.MapFS{}

	mockFileSystem.EXPECT().Glob(logFS, "server-*.log").Return([]string{"server-a.log", "server-b.log"}, nil).Once()
	mockFileSystem.EXPECT().ReadFile(logFS, "server-a.log").Return([]byte("a\n"), nil).Once()
	mockFileSystem.EXPECT().ReadFile(logFS, "server-b.log").Return([]byte("b\n"), nil).Once()

	reader, err := logs.NewReaderWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	content, err := reader.ReadCombined(logFS, "server-*.log")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "a\nb\n", content)
}

func TestReader_ReadCombined_GlobError(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockLogFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	var logFS fs.FS = fstest.MapFS{}

	mockFileSystem.EXPECT().Glob(logFS, "server-*.log").Return(nil, assert.AnError).Once()

	reader, err := logs.NewReaderWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	_, err = reader.ReadCombined(logFS, "server-*.log")

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	assert.Contains(t, err.Error(), "failed to glob logs")
}

func TestReader_ReadCombined_NoMatchingFiles(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockLogFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	var logFS fs.FS = fstest.MapFS{}

	mockFileSystem.EXPECT().Glob(logFS, "server-*.log").Return([]string{}, nil).Once()

	reader, err := logs.NewReaderWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	_, err = reader.ReadCombined(logFS, "server-*.log")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no logs found")
}

func TestReader_ReadCombined_ReadFileError(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockLogFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	var logFS fs.FS = fstest.MapFS{}

	mockFileSystem.EXPECT().Glob(logFS, "server-*.log").Return([]string{"server-a.log"}, nil).Once()
	mockFileSystem.EXPECT().ReadFile(logFS, "server-a.log").Return(nil, assert.AnError).Once()

	reader, err := logs.NewReaderWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	_, err = reader.ReadCombined(logFS, "server-*.log")

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	assert.Contains(t, err.Error(), "failed to read log file")
}

func TestReader_ReadEntries_HappyPath(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockLogFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	var logFS fs.FS = fstest.MapFS{}
	dumpPatterns := []logs.DumpPattern{
		{Glob: "server-*.log", Header: "Server"},
		{Glob: "watchdog-*.log", Header: "Watchdog"},
	}

	mockFileSystem.EXPECT().Glob(logFS, "server-*.log").Return([]string{"nested/server-1.log"}, nil).Once()
	mockFileSystem.EXPECT().ReadFile(logFS, "nested/server-1.log").Return([]byte("server"), nil).Once()
	mockFileSystem.EXPECT().Glob(logFS, "watchdog-*.log").Return([]string{"watchdog-1.log"}, nil).Once()
	mockFileSystem.EXPECT().ReadFile(logFS, "watchdog-1.log").Return([]byte("watchdog"), nil).Once()

	reader, err := logs.NewReaderWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	entries, err := reader.ReadEntries(logFS, dumpPatterns)

	// Assert
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "Server", entries[0].Header)
	assert.Equal(t, filepath.Base("nested/server-1.log"), entries[0].File)
	assert.Equal(t, "server", entries[0].Content)
	assert.Equal(t, "Watchdog", entries[1].Header)
	assert.Equal(t, "watchdog-1.log", entries[1].File)
	assert.Equal(t, "watchdog", entries[1].Content)
}

func TestReader_ReadEntries_GlobError(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockLogFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	var logFS fs.FS = fstest.MapFS{}
	dumpPatterns := []logs.DumpPattern{{Glob: "server-*.log", Header: "Server"}}

	mockFileSystem.EXPECT().Glob(logFS, "server-*.log").Return(nil, assert.AnError).Once()

	reader, err := logs.NewReaderWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	_, err = reader.ReadEntries(logFS, dumpPatterns)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	assert.Contains(t, err.Error(), "failed to glob logs for pattern")
}

func TestReader_ReadEntries_ReadFileError(t *testing.T) {
	// Arrange
	mockFileSystem := mocks.NewMockLogFileSystem(t)
	defer mockFileSystem.AssertExpectations(t)
	var logFS fs.FS = fstest.MapFS{}
	dumpPatterns := []logs.DumpPattern{{Glob: "server-*.log", Header: "Server"}}

	mockFileSystem.EXPECT().Glob(logFS, "server-*.log").Return([]string{"server-1.log"}, nil).Once()
	mockFileSystem.EXPECT().ReadFile(logFS, "server-1.log").Return(nil, assert.AnError).Once()

	reader, err := logs.NewReaderWithFileSystem(mockFileSystem)
	require.NoError(t, err)

	// Act
	_, err = reader.ReadEntries(logFS, dumpPatterns)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	assert.Contains(t, err.Error(), "failed to read log file")
}
