// Copyright 2026 The MathWorks, Inc.

package mcpclient_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	mocks "github.com/matlab/matlab-mcp-core-server/tests/mocks/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/logs"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggedSession_ReadLogMethods(t *testing.T) {
	// Arrange
	mockLogReader := mocks.NewMockLogReader(t)
	defer mockLogReader.AssertExpectations(t)
	mockFileSystemProvider := mocks.NewMockFileSystemProvider(t)
	defer mockFileSystemProvider.AssertExpectations(t)

	logDir := "/tmp/logs"
	var logFS fs.FS = fstest.MapFS{}

	mockFileSystemProvider.EXPECT().DirFS(logDir).Return(logFS).Once()
	mockLogReader.EXPECT().ReadCombined(logFS, "server-*.log").Return("server-data", nil).Twice()
	mockLogReader.EXPECT().ReadCombined(logFS, "watchdog-*.log").Return("watchdog-data", nil).Once()
	factory, err := mcpclient.NewLoggedSessionFactory(mockLogReader, mockFileSystemProvider)
	require.NoError(t, err)

	session, err := factory.New(
		nil,
		logDir,
		"Server stderr",
		nil,
	)
	require.NoError(t, err)

	// Act
	serverLogs, err := session.ReadServerLogs()
	require.NoError(t, err)

	allServerLogs, err := session.ReadAllServerLogs()
	require.NoError(t, err)

	watchdogLogs, err := session.ReadWatchdogLogs()
	require.NoError(t, err)

	// Assert
	assert.Equal(t, logDir, session.LogDir())
	assert.Equal(t, logFS, session.LogFS())
	assert.Equal(t, "server-data", serverLogs)
	assert.Equal(t, "server-data", allServerLogs)
	assert.Equal(t, "watchdog-data", watchdogLogs)
}

func TestLoggedSession_ReadLogs_NoMatchingFiles(t *testing.T) {
	// Arrange
	mockLogReader := mocks.NewMockLogReader(t)
	defer mockLogReader.AssertExpectations(t)

	logDir := "/tmp/logs"
	var logFS fs.FS = fstest.MapFS{}

	mockLogReader.EXPECT().ReadCombined(logFS, "server-*.log").Return("", assert.AnError).Once()

	session, err := mcpclient.NewLoggedSession(
		nil,
		logDir,
		logFS,
		"Server stderr",
		nil,
		mockLogReader,
		nil,
	)
	require.NoError(t, err)

	// Act
	_, err = session.ReadServerLogs()

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestLoggedSession_CollectDumpData_HappyPath(t *testing.T) {
	// Arrange
	mockLogReader := mocks.NewMockLogReader(t)
	defer mockLogReader.AssertExpectations(t)
	mockStderrProvider := mocks.NewMockStderrProvider(t)
	defer mockStderrProvider.AssertExpectations(t)

	logDir := "/tmp/logs"
	var logFS fs.FS = fstest.MapFS{}
	dumpPatterns := []logs.DumpPattern{
		{Glob: "server-*.log", Header: "Server Log"},
		{Glob: "watchdog-*.log", Header: "Watchdog Log"},
	}
	expectedEntries := []logs.DumpEntry{
		{Header: "Server Log", File: "server-1.log", Content: "server-data"},
		{Header: "Watchdog Log", File: "watchdog-1.log", Content: "watchdog-data"},
	}

	mockLogReader.EXPECT().ReadEntries(logFS, dumpPatterns).Return(expectedEntries, nil).Once()
	mockStderrProvider.EXPECT().Stderr().Return("stderr-data").Once()

	session, err := mcpclient.NewLoggedSession(
		nil,
		logDir,
		logFS,
		"Server stderr",
		dumpPatterns,
		mockLogReader,
		mockStderrProvider,
	)
	require.NoError(t, err)

	// Act
	stderr, entries, err := session.CollectDumpData()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "stderr-data", stderr)
	assert.Equal(t, expectedEntries, entries)
}

func TestLoggedSession_CollectDumpData_InvalidPattern(t *testing.T) {
	// Arrange
	mockLogReader := mocks.NewMockLogReader(t)
	defer mockLogReader.AssertExpectations(t)
	mockStderrProvider := mocks.NewMockStderrProvider(t)
	defer mockStderrProvider.AssertExpectations(t)

	logDir := "/tmp/logs"
	var logFS fs.FS = fstest.MapFS{}
	dumpPatterns := []logs.DumpPattern{{Glob: "[", Header: "Server Log"}}

	mockLogReader.EXPECT().ReadEntries(logFS, dumpPatterns).Return(nil, assert.AnError).Once()

	session, err := mcpclient.NewLoggedSession(
		nil,
		logDir,
		logFS,
		"Server stderr",
		dumpPatterns,
		mockLogReader,
		mockStderrProvider,
	)
	require.NoError(t, err)

	// Act
	_, _, err = session.CollectDumpData()

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestLoggedSession_CollectDumpData_DefaultsToEmptyStderrWithoutSession(t *testing.T) {
	// Arrange
	mockLogReader := mocks.NewMockLogReader(t)
	defer mockLogReader.AssertExpectations(t)

	logDir := "/tmp/logs"
	var logFS fs.FS = fstest.MapFS{}
	dumpPatterns := []logs.DumpPattern{{Glob: "server-*.log", Header: "Server Log"}}
	expectedEntries := []logs.DumpEntry{{Header: "Server Log", File: "server-1.log", Content: "server-data"}}

	mockLogReader.EXPECT().ReadEntries(logFS, dumpPatterns).Return(expectedEntries, nil).Once()

	session, err := mcpclient.NewLoggedSession(
		nil,
		logDir,
		logFS,
		"Server stderr",
		dumpPatterns,
		mockLogReader,
		nil,
	)
	require.NoError(t, err)

	// Act
	stderr, entries, err := session.CollectDumpData()

	// Assert
	require.NoError(t, err)
	assert.Empty(t, stderr)
	assert.Equal(t, expectedEntries, entries)
}

func TestNewLoggedSessionFactory_NilDependencyReturnsError(t *testing.T) {
	// Arrange
	mockFileSystemProvider := mocks.NewMockFileSystemProvider(t)
	defer mockFileSystemProvider.AssertExpectations(t)

	// Act
	_, err := mcpclient.NewLoggedSessionFactory(nil, mockFileSystemProvider)

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "logReader must not be nil")

	// Act
	mockLogReader := mocks.NewMockLogReader(t)
	defer mockLogReader.AssertExpectations(t)
	_, err = mcpclient.NewLoggedSessionFactory(mockLogReader, nil)

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "fileSystemProvider must not be nil")
}

func TestNewLoggedSession_NilInputReturnsError(t *testing.T) {
	// Act
	_, err := mcpclient.NewLoggedSession(nil, "/tmp/logs", fstest.MapFS{}, "Server stderr", nil, nil, nil)

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "logReader must not be nil")

	// Arrange
	logReader := mocks.NewMockLogReader(t)
	defer logReader.AssertExpectations(t)

	// Act
	_, err = mcpclient.NewLoggedSession(nil, "/tmp/logs", nil, "Server stderr", nil, logReader, nil)

	// Assert
	require.Error(t, err)
	require.EqualError(t, err, "logFS must not be nil")
}
