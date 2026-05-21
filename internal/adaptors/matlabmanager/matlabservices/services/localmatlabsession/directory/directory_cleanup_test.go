// Copyright 2025-2026 The MathWorks, Inc.

package directory_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectory_Cleanup_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	cleanupTimeout := 100 * time.Millisecond
	cleanupRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(10 * time.Minute).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetCleanupTimeout(cleanupTimeout)
	dir.SetCleanupRetry(cleanupRetry)

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(nil).
		Once()

	// Act
	err := dir.Cleanup()

	// Assert
	require.NoError(t, err)
}

func TestDirectory_Cleanup_WaitsForRemoveAllToPass(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	cleanupTimeout := 100 * time.Millisecond
	cleanupRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(10 * time.Minute).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetCleanupTimeout(cleanupTimeout)
	dir.SetCleanupRetry(cleanupRetry)

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(assert.AnError).
		Once()

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(nil).
		Once()

	// Act
	err := dir.Cleanup()

	// Assert
	require.NoError(t, err)
}

func TestDirectory_Cleanup_Timesout(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	cleanupTimeout := 100 * time.Millisecond
	cleanupRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(10 * time.Minute).
		Once()

	dir := directory.NewDirectory(mockLogger, sessionDir, mockOSLayer, mockConfig)
	dir.SetCleanupTimeout(cleanupTimeout)
	dir.SetCleanupRetry(cleanupRetry)

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(assert.AnError) // Will be called many times with retry

	// Act
	err := dir.Cleanup()

	// Assert
	require.Error(t, err)
}

func TestDirectory_Cleanup_EmptySessionDir(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	cleanupTimeout := 100 * time.Millisecond
	cleanupRetry := 10 * time.Millisecond

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(10 * time.Minute).
		Once()

	dir := directory.NewDirectory(mockLogger, "", mockOSLayer, mockConfig)
	dir.SetCleanupTimeout(cleanupTimeout)
	dir.SetCleanupRetry(cleanupRetry)

	// Act
	err := dir.Cleanup()

	// Assert
	require.NoError(t, err)
}
