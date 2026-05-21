// Copyright 2026 The MathWorks, Inc.

package mockruntime_test

import (
	"crypto/tls"
	"errors"
	"os"
	"testing"

	mocks "github.com/matlab/matlab-mcp-core-server/tests/mocks/testutils/mockmatlab/mockruntime"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntime_GenerateAndWriteCert_WhenCertWriteFails_ReturnsError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("cert write failed")
	certPEM := []byte("cert")
	keyPEM := []byte("key")
	env := mocks.NewMockEnvironment(t)
	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().WriteFile("cert.pem", certPEM, os.FileMode(0o600)).Return(expectedErr)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)
	tlsProvider.EXPECT().GeneratePEM().Return(certPEM, keyPEM, nil)
	runtime := mockruntime.NewRuntime(env, fileSystem, tlsProvider)

	// Act
	_, err := runtime.GenerateAndWriteCert("cert.pem", "key.pem")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write cert file")
	assert.ErrorIs(t, err, expectedErr)
}

func TestRuntime_GenerateAndWriteCert_WhenKeyWriteFails_ReturnsError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("key write failed")
	certPEM := []byte("cert")
	keyPEM := []byte("key")
	env := mocks.NewMockEnvironment(t)
	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().WriteFile("cert.pem", certPEM, os.FileMode(0o600)).Return(nil)
	fileSystem.EXPECT().WriteFile("key.pem", keyPEM, os.FileMode(0o600)).Return(expectedErr)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)
	tlsProvider.EXPECT().GeneratePEM().Return(certPEM, keyPEM, nil)
	runtime := mockruntime.NewRuntime(env, fileSystem, tlsProvider)

	// Act
	_, err := runtime.GenerateAndWriteCert("cert.pem", "key.pem")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write key file")
	assert.ErrorIs(t, err, expectedErr)
}

func TestNewDefaultTLSMaterialProvider_TLSConfig_InvalidPEM_ReturnsError(t *testing.T) {
	// Arrange
	provider := mockruntime.NewDefaultTLSMaterialProvider()

	// Act
	_, err := provider.TLSConfig([]byte("not-a-cert"), []byte("not-a-key"))

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load X509 key pair")
}

func TestGenerateAndWriteCert_Success_UsesInjectedDependencies(t *testing.T) {
	// Arrange
	certPEM := []byte("cert")
	keyPEM := []byte("key")
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	env := mocks.NewMockEnvironment(t)
	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().WriteFile("cert.pem", certPEM, os.FileMode(0o600)).Return(nil)
	fileSystem.EXPECT().WriteFile("cert.key", keyPEM, os.FileMode(0o600)).Return(nil)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)
	tlsProvider.EXPECT().GeneratePEM().Return(certPEM, keyPEM, nil)
	tlsProvider.EXPECT().TLSConfig(certPEM, keyPEM).Return(tlsConfig, nil)
	runtime := mockruntime.NewRuntime(
		env,
		fileSystem,
		tlsProvider,
	)

	// Act
	actualTLSConfig, err := runtime.GenerateAndWriteCert("cert.pem", "cert.key")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, actualTLSConfig)
	assert.Equal(t, uint16(tls.VersionTLS12), actualTLSConfig.MinVersion)
}
