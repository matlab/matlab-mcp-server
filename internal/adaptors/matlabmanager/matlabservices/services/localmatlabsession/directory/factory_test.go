// Copyright 2025-2026 The MathWorks, Inc.

package directory_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	appconfigmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	applicationdirectorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	// Act
	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles, mockConfigFactory)

	// Assert
	assert.NotNil(t, factory)
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockApplicationDirectory := &applicationdirectorymocks.MockDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &appconfigmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSessionDir := filepath.Join("tmp", "matlab-session-12345")
	packageDir := filepath.Join(expectedSessionDir, "+matlab_mcp")
	expectedCertificateFile := filepath.Join(expectedSessionDir, "cert.pem")
	expectedCertificateKeyFile := filepath.Join(expectedSessionDir, "cert.key")
	expectedMATLABFiles := map[string][]byte{
		"initializeMCP.m": []byte("some content"),
		"eval.m":          []byte("some other content"),
	}

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(mockApplicationDirectory, nil).
		Once()

	mockApplicationDirectory.EXPECT().
		CreateSubDir("matlab-session-").
		Return(expectedSessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	mockMATLABFiles.EXPECT().
		GetAll().
		Return(expectedMATLABFiles).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(10 * time.Minute).
		Once()

	for fileName, fileContent := range expectedMATLABFiles {
		filePath := filepath.Join(packageDir, fileName)
		mockOSLayer.EXPECT().
			WriteFile(filePath, fileContent, os.FileMode(0o600)).
			Return(nil).
			Once()
	}

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles, mockConfigFactory)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, dir)
	assert.Equal(t, expectedSessionDir, dir.Path())
	assert.Equal(t, expectedCertificateFile, dir.CertificateFile())
	assert.Equal(t, expectedCertificateKeyFile, dir.CertificateKeyFile())
}

func TestFactory_New_ConfigFactoryError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles, mockConfigFactory)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)
}

func TestFactory_New_DirectoryFactoryError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &appconfigmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(nil, expectedError).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles, mockConfigFactory)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)
}

func TestFactory_New_CreateSubDirError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockApplicationDirectory := &applicationdirectorymocks.MockDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &appconfigmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(mockApplicationDirectory, nil).
		Once()

	mockApplicationDirectory.EXPECT().
		CreateSubDir("matlab-session-").
		Return("", expectedError).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles, mockConfigFactory)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)
}

func TestFactory_New_PackageDirectoryMkdirError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockApplicationDirectory := &applicationdirectorymocks.MockDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &appconfigmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	packageDir := filepath.Join(sessionDir, "+matlab_mcp")
	expectedError := assert.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(mockApplicationDirectory, nil).
		Once()

	mockApplicationDirectory.EXPECT().
		CreateSubDir("matlab-session-").
		Return(sessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(expectedError).
		Once()

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(nil).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles, mockConfigFactory)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)
}

func TestFactory_New_WriteFileError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockApplicationDirectory := &applicationdirectorymocks.MockDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &appconfigmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	packageDir := filepath.Join(sessionDir, "+matlab_mcp")
	expectedError := assert.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(mockApplicationDirectory, nil).
		Once()

	mockApplicationDirectory.EXPECT().
		CreateSubDir("matlab-session-").
		Return(sessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	expectedFailingFileName := "initializeMCP.m"

	expectedMATLABFiles := map[string][]byte{
		expectedFailingFileName: []byte("some other content"),
	}

	mockMATLABFiles.EXPECT().
		GetAll().
		Return(expectedMATLABFiles).
		Once()

	mockOSLayer.EXPECT().
		WriteFile(filepath.Join(packageDir, expectedFailingFileName), expectedMATLABFiles[expectedFailingFileName], os.FileMode(0o600)).
		Return(expectedError).
		Once()

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(nil).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles, mockConfigFactory)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)
}

func TestFactory_New_CleanupFailureOnError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockApplicationDirectory := &applicationdirectorymocks.MockDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &appconfigmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	packageDir := filepath.Join(sessionDir, "+matlab_mcp")
	expectedError := assert.AnError
	cleanupError := errors.New("cleanup failed")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(mockApplicationDirectory, nil).
		Once()

	mockApplicationDirectory.EXPECT().
		CreateSubDir("matlab-session-").
		Return(sessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(expectedError).
		Once()

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(cleanupError).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles, mockConfigFactory)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)

	warnLogs := mockLogger.WarnLogs()
	fields, found := warnLogs["Failed to cleanup session directory during error handling"]
	require.True(t, found)
	assert.Equal(t, cleanupError, fields["error"])
}
