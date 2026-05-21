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

func TestRuntime_WriteStartupFailureFile_WhenWriteFails_ReturnsError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("write failed")
	sessionDir := filepath.Join("fake", "session", "dir")
	expectedPath := filepath.Join(sessionDir, "mcp_startup_error.txt")
	env := mocks.NewMockEnvironment(t)
	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().WriteFile(
		expectedPath,
		mock.MatchedBy(func(content []byte) bool {
			return len(content) > 0
		}),
		os.FileMode(0o600),
	).Return(expectedErr)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)
	runtime := mockruntime.NewRuntime(env, fileSystem, tlsProvider)

	// Act
	err := runtime.WriteStartupFailureFile(sessionDir)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write startup error file")
	assert.ErrorIs(t, err, expectedErr)
}

func TestWriteStartupFailureFile_MissingSessionDir_ReturnsError(t *testing.T) {
	// Arrange
	fileSystem := mocks.NewMockFileSystem(t)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)
	runtime := mockruntime.NewRuntime(mocks.NewMockEnvironment(t), fileSystem, tlsProvider)
	sessionDir := ""

	// Act
	err := runtime.WriteStartupFailureFile(sessionDir)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MW_MCP_SESSION_DIR")
	fileSystem.AssertNotCalled(t, "WriteFile", mock.Anything, mock.Anything, mock.Anything)
}

func TestWriteStartupFailureFile_Success_UsesInjectedFileSystem(t *testing.T) {
	// Arrange
	sessionDir := filepath.Join("fake", "session", "dir")
	expectedPath := filepath.Join(sessionDir, "mcp_startup_error.txt")
	fileSystem := mocks.NewMockFileSystem(t)
	fileSystem.EXPECT().WriteFile(
		expectedPath,
		mock.MatchedBy(func(content []byte) bool {
			return len(content) > 0 && string(content) != ""
		}),
		mock.Anything,
	).Return(nil)
	runtime := mockruntime.NewRuntime(
		mocks.NewMockEnvironment(t),
		fileSystem,
		mocks.NewMockTLSMaterialProvider(t),
	)

	// Act
	err := runtime.WriteStartupFailureFile(sessionDir)

	// Assert
	require.NoError(t, err)
}
