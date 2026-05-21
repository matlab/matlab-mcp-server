// Copyright 2025-2026 The MathWorks, Inc.

package directory_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectory_GetEmbeddedConnectorDetails_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond
	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_WaitsForSecurePortFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond
	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_WaitsForCertificateFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond
	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_WaitsForNotEmptyPortFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond
	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(""), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_WaitsForNotEmptyCertificateFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond
	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return([]byte(""), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_TimesoutWaitingForFilesToExists(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist)

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(nil, os.ErrNotExist)

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_TimesoutWaitingForFileContent(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(""), nil) // Will be called multiple times in wait loop

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_ReadSecurePortFileError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return(nil, assert.AnError).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read secure port file")
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_StatSecurePortFileError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(nil, assert.AnError).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stat secure port file")
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_ReadCertificateFileError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond
	expectedPort := "9999"

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(nil, assert.AnError).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read certificate path file")
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_StatCertificateFileError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()
	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(nil, assert.AnError).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stat certificate file")
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_StartupErrorFileDetected(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond
	startupErrorContent := "License checkout failed: No license available"

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return([]byte(startupErrorContent), nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, directory.ErrMATLABStartup)
	assert.Contains(t, err.Error(), "MATLAB startup failed")
	assert.Contains(t, err.Error(), startupErrorContent)
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_StartupErrorFileReadError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	embeddedConnectorDetailsRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(100 * time.Millisecond).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetEmbeddedConnectorDetailsRetry(embeddedConnectorDetailsRetry)

	startupErrorFile := dir.StartupErrorFile()

	mockOSLayer.EXPECT().
		ReadFile(startupErrorFile).
		Return(nil, assert.AnError).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read startup error file")
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}
