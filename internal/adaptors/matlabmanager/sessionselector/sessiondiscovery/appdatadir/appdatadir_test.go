// Copyright 2026 The MathWorks, Inc.

package appdatadir_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/sessionselector/sessiondiscovery/appdatadir"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/sessionselector/sessiondiscovery/appdatadir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	// Act
	result := appdatadir.New(mockOSLayer)

	// Assert
	require.NotNil(t, result)
}

func TestGetter_AppDataDir_Linux(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedHome := filepath.Join("home", "user")

	mockOSLayer.EXPECT().
		GOOS().
		Return("linux").
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(expectedHome, nil).
		Once()

	getter := appdatadir.New(mockOSLayer)

	// Act
	result, err := getter.AppDataDir()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(expectedHome, ".MathWorks", "MATLABMCPCoreServer"), result)
}

func TestGetter_AppDataDir_Darwin(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedHome := filepath.Join("Users", "user")

	mockOSLayer.EXPECT().
		GOOS().
		Return("darwin").
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(expectedHome, nil).
		Once()

	getter := appdatadir.New(mockOSLayer)

	// Act
	result, err := getter.AppDataDir()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(expectedHome, "Library", "Application Support", "MathWorks", "MATLAB MCP Core Server"), result)
}

func TestGetter_AppDataDir_Windows(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedAppData := filepath.Join("C:", "Users", "user", "AppData", "Roaming")

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		Getenv("APPDATA").
		Return(expectedAppData).
		Once()

	getter := appdatadir.New(mockOSLayer)

	// Act
	result, err := getter.AppDataDir()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(expectedAppData, "MathWorks", "MATLAB MCP Core Server"), result)
}

func TestGetter_AppDataDir_DarwinHomeDirError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		GOOS().
		Return("darwin").
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return("", assert.AnError).
		Once()

	getter := appdatadir.New(mockOSLayer)

	// Act
	result, err := getter.AppDataDir()

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, result)
}

func TestGetter_AppDataDir_HomeDirError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		GOOS().
		Return("linux").
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return("", assert.AnError).
		Once()

	getter := appdatadir.New(mockOSLayer)

	// Act
	result, err := getter.AppDataDir()

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, result)
}

func TestGetter_AppDataDir_WindowsEmptyAPPDATA(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		Getenv("APPDATA").
		Return("").
		Once()

	getter := appdatadir.New(mockOSLayer)

	// Act
	result, err := getter.AppDataDir()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "APPDATA")
	assert.Empty(t, result)
}

func TestGetter_AppDataDir_UnsupportedOS(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		GOOS().
		Return("freebsd").
		Once()

	getter := appdatadir.New(mockOSLayer)

	// Act
	result, err := getter.AppDataDir()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "freebsd")
	assert.Empty(t, result)
}
