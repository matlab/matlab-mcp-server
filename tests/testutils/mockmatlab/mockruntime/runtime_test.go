// Copyright 2026 The MathWorks, Inc.

package mockruntime_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	mocks "github.com/matlab/matlab-mcp-core-server/tests/mocks/testutils/mockmatlab/mockruntime"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRuntime_LoadConfigFromEnv_UsesInjectedEnvironment(t *testing.T) {
	// Arrange
	env := mocks.NewMockEnvironment(t)
	env.EXPECT().Getenv(mockruntime.EnvMockMATLABConfig).Return(`{"mode":"slow_startup","delayMs":250}`)

	fileSystem := mocks.NewMockFileSystem(t)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)

	runtime := mockruntime.NewRuntime(
		env,
		fileSystem,
		tlsProvider,
	)

	// Act
	cfg, err := runtime.LoadConfigFromEnv()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mockruntime.ModeSlowStartup, cfg.Mode)
	require.NotNil(t, cfg.DelayMs)
	assert.Equal(t, 250, *cfg.DelayMs)
}

func TestRuntime_WriteStartupFailureFile_UsesInjectedFileSystem(t *testing.T) {
	// Arrange
	env := mocks.NewMockEnvironment(t)
	fileSystem := mocks.NewMockFileSystem(t)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)

	sessionDir := filepath.Join("fake", "session", "dir")
	expectedPath := filepath.Join(sessionDir, "mcp_startup_error.txt")
	fileSystem.EXPECT().WriteFile(
		expectedPath,
		mock.MatchedBy(func(content []byte) bool {
			return len(content) > 0
		}),
		os.FileMode(0o600),
	).Return(nil)

	runtime := mockruntime.NewRuntime(env, fileSystem, tlsProvider)

	// Act
	err := runtime.WriteStartupFailureFile(sessionDir)

	// Assert
	require.NoError(t, err)
}

func TestRuntime_GenerateAndWriteCert_GeneratePEMError_ReturnsError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("generate failed")
	env := mocks.NewMockEnvironment(t)
	fileSystem := mocks.NewMockFileSystem(t)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)
	tlsProvider.EXPECT().GeneratePEM().Return(nil, nil, expectedErr)
	runtime := mockruntime.NewRuntime(env, fileSystem, tlsProvider)

	// Act
	_, err := runtime.GenerateAndWriteCert("cert.pem", "cert.key")

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}

func TestRuntime_GenerateAndWriteCert_TLSConfigError_ReturnsError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("tls config failed")
	certPEM := []byte("cert")
	keyPEM := []byte("key")
	env := mocks.NewMockEnvironment(t)
	fileSystem := mocks.NewMockFileSystem(t)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)
	tlsProvider.EXPECT().GeneratePEM().Return(certPEM, keyPEM, nil)
	fileSystem.EXPECT().WriteFile("cert.pem", certPEM, os.FileMode(0o600)).Return(nil)
	fileSystem.EXPECT().WriteFile("cert.key", keyPEM, os.FileMode(0o600)).Return(nil)
	tlsProvider.EXPECT().TLSConfig(certPEM, keyPEM).Return(nil, expectedErr)

	runtime := mockruntime.NewRuntime(
		env,
		fileSystem,
		tlsProvider,
	)

	// Act
	_, err := runtime.GenerateAndWriteCert("cert.pem", "cert.key")

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
