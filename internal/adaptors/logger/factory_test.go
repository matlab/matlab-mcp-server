// Copyright 2025-2026 The MathWorks, Inc.

package logger_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/logger"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	directorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	loggermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/logger"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	// Act
	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_NewMCPSessionLogger_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogFile := &osfacademocks.MockFile{}
	defer mockLogFile.AssertExpectations(t)

	expectedBaseDir := filepath.Join("some", "directory")
	expectedSuffix := "1337"
	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(entities.LogLevelInfo).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(mockLogFile, nil).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	logger, err := factory.NewMCPSessionLogger(&mcp.ServerSession{})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, logger, "Logger should not be nil")
}

func TestFactory_NewMCPSessionLogger_ReturnsErrorWhenConfigFails(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, messages.AnError).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	result, err := factory.NewMCPSessionLogger(&mcp.ServerSession{})

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	assert.Nil(t, result)
}

func TestFactory_NewMCPSessionLogger_ReturnsErrorWhenLogLevelIsUnknown(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	invalidLogLevel := entities.LogLevel("invalid")
	expectedError := messages.New_StartupErrors_InvalidLogLevel_Error(string(invalidLogLevel))

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(invalidLogLevel).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	result, err := factory.NewMCPSessionLogger(&mcp.ServerSession{})

	// Assert
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

func TestFactory_NewMCPSessionLogger_ReturnsErrorWhenDirectoryFails(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(entities.LogLevelDebug).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(nil, messages.AnError).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	result, err := factory.NewMCPSessionLogger(&mcp.ServerSession{})

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	assert.Nil(t, result)
}

func TestFactory_NewMCPSessionLogger_ReturnsErrorWhenLogFileCreationFails(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedBaseDir := filepath.Join("some", "directory")
	expectedSuffix := "1337"
	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")
	expectedError := messages.New_StartupErrors_FailedToCreateLogFile_Error(expectedLogFile)

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(entities.LogLevelDebug).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(nil, assert.AnError).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	result, err := factory.NewMCPSessionLogger(&mcp.ServerSession{})

	// Assert
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

func TestFactory_GetGlobalLogger_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogFile := &osfacademocks.MockFile{}
	defer mockLogFile.AssertExpectations(t)

	expectedBaseDir := filepath.Join("some", "directory")
	expectedSuffix := "1337"
	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(entities.LogLevelDebug).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(mockLogFile, nil).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	logger, err := factory.GetGlobalLogger()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, logger, "Global logger should not be nil")
}

func TestFactory_GetGlobalLogger_DuplicatesLogsToStderrWhenEnabled(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogFile := &osfacademocks.MockFile{}
	defer mockLogFile.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	expectedBaseDir := filepath.Join("some", "directory")
	expectedSuffix := "1337"
	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")
	expectedMessage := "stderr log message"
	var stderrWritten string

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(true).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(entities.LogLevelDebug).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(mockLogFile, nil).
		Once()

	mockOSLayer.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	mockStderr.EXPECT().
		Write(mock.AnythingOfType("[]uint8")).
		RunAndReturn(func(p []byte) (int, error) {
			stderrWritten += string(p)
			return len(p), nil
		}).
		Once()

	mockLogFile.EXPECT().
		Write(mock.AnythingOfType("[]uint8")).
		RunAndReturn(func(p []byte) (int, error) {
			return len(p), nil
		}).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	globalLogger, err := factory.GetGlobalLogger()
	require.NoError(t, err)
	globalLogger.Info(expectedMessage)

	// Assert
	assert.Contains(t, stderrWritten, expectedMessage, "Stderr should contain the logged message")
}

func TestFactory_GetGlobalLogger_UsesWatchdogLogFileInWatchdogMode(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogFile := &osfacademocks.MockFile{}
	defer mockLogFile.AssertExpectations(t)

	expectedBaseDir := filepath.Join("some", "directory")
	expectedSuffix := "1337"
	expectedLogFile := filepath.Join(expectedBaseDir, "watchdog.log")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(entities.LogLevelDebug).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(true).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.WatchdogLogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(mockLogFile, nil).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	logger, err := factory.GetGlobalLogger()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, logger, "Global logger should not be nil")
}

func TestFactory_GetGlobalLogger_IsSingleton(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogFile := &osfacademocks.MockFile{}
	defer mockLogFile.AssertExpectations(t)

	expectedBaseDir := filepath.Join("some", "directory")
	expectedSuffix := "1337"
	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(entities.LogLevelWarn).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(mockLogFile, nil).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	logger1, err1 := factory.GetGlobalLogger()
	logger2, err2 := factory.GetGlobalLogger()

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotNil(t, logger1, "First global logger should not be nil")
	assert.NotNil(t, logger2, "Second global logger should not be nil")
	assert.Same(t, logger1, logger2, "Global logger should be a singleton")
}

func TestFactory_GetGlobalLogger_ReturnsErrorWhenConfigFails(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, messages.AnError).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	result, err := factory.GetGlobalLogger()

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	assert.Nil(t, result)
}

func TestFactory_GetGlobalLogger_ReturnsErrorWhenLogLevelIsUnknown(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	unknownLogLevel := entities.LogLevel("unknown")
	expectedError := messages.New_StartupErrors_InvalidLogLevel_Error(string(unknownLogLevel))

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(unknownLogLevel).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	result, err := factory.GetGlobalLogger()

	// Assert
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

func TestFactory_GetGlobalLogger_ReturnsErrorWhenDirectoryFails(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(entities.LogLevelInfo).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(nil, messages.AnError).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	result, err := factory.GetGlobalLogger()

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	assert.Nil(t, result)
}

func TestFactory_GetGlobalLogger_ReturnsErrorWhenLogFileCreationFails(t *testing.T) {
	// Arrange
	mockConfigFactory := &loggermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectoryFactory := &loggermocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedBaseDir := filepath.Join("some", "directory")
	expectedSuffix := "1337"
	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")
	expectedError := messages.New_StartupErrors_FailedToCreateLogFile_Error(expectedLogFile)

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		DuplicateLogsToStderr().
		Return(false).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(entities.LogLevelInfo).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(nil, assert.AnError).
		Once()

	factory := logger.NewFactory(mockConfigFactory, mockDirectoryFactory, mockFilenameFactory, mockOSLayer)

	// Act
	result, err := factory.GetGlobalLogger()

	// Assert
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
}
