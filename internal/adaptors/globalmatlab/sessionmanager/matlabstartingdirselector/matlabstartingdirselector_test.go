// Copyright 2025-2026 The MathWorks, Inc.

package matlabstartingdirselector_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/sessionmanager/matlabstartingdirselector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab/sessionmanager/matlabstartingdirselector"
	osFacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	// Act
	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Assert
	assert.NotNil(t, selector)
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_ConfigError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError, "SelectMATLABStartingDir should return the error from Config")
	assert.Empty(t, result)
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_HappyPath(t *testing.T) {
	testCases := []struct {
		name        string
		os          string
		homeDir     string
		expectedDir string
	}{
		{
			name:        "Windows",
			os:          "windows",
			homeDir:     filepath.Join("Users", "testuser"),
			expectedDir: filepath.Join("Users", "testuser", "Documents"),
		},
		{
			name:        "Darwin",
			os:          "darwin",
			homeDir:     filepath.Join("Users", "testuser"),
			expectedDir: filepath.Join("Users", "testuser", "Documents"),
		},
		{
			name:        "Linux",
			os:          "linux",
			homeDir:     filepath.Join("home", "testuser"),
			expectedDir: filepath.Join("home", "testuser"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &mocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockConfigFactory := &mocks.MockConfigFactory{}
			defer mockConfigFactory.AssertExpectations(t)

			mockConfig := &configmocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			mockRootStore := &mocks.MockRootStore{}
			defer mockRootStore.AssertExpectations(t)

			mockRootPathResolver := &mocks.MockRootPathResolver{}
			defer mockRootPathResolver.AssertExpectations(t)

			mockFileInfo := &osFacademocks.MockFileInfo{}
			defer mockFileInfo.AssertExpectations(t)

			mockLogger := testutils.NewInspectableLogger()

			mockConfigFactory.EXPECT().
				Config().
				Return(mockConfig, nil).
				Once()

			mockConfig.EXPECT().
				PreferredMATLABStartingDirectory().
				Return("").
				Once()

			mockRootStore.EXPECT().
				GetRoots().
				Return([]entities.MCPRoot{}).
				Once()

			mockOSLayer.EXPECT().
				UserHomeDir().
				Return(tc.homeDir, nil).
				Once()

			mockOSLayer.EXPECT().
				GOOS().
				Return(tc.os).
				Once()

			mockOSLayer.EXPECT().
				Stat(tc.expectedDir).
				Return(mockFileInfo, nil).
				Once()

			selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

			// Act
			result, err := selector.SelectMATLABStartingDir(mockLogger)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.expectedDir, result)
		})
	}
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_PreferredStartingDirectorySetHappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedPreferredMATLABStartingDir := filepath.Join("custom", "preferred", "directory")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return(expectedPreferredMATLABStartingDir).
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedPreferredMATLABStartingDir).
		Return(mockFileInfo, nil).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPreferredMATLABStartingDir, result)
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_UnknownOSHappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedHomeDir := filepath.Join("home", "testuser")
	unknownOS := "freebsd"

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockRootStore.EXPECT().
		GetRoots().
		Return([]entities.MCPRoot{}).
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(expectedHomeDir, nil).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return(unknownOS).
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedHomeDir).
		Return(mockFileInfo, nil).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedHomeDir, result)
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_UserHomeDirError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := assert.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockRootStore.EXPECT().
		GetRoots().
		Return([]entities.MCPRoot{}).
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return("", expectedError).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_StatErrorOnHomeDir(t *testing.T) {
	testCases := []struct {
		name    string
		os      string
		homeDir string
	}{
		{
			name:    "Windows - Stat Error",
			os:      "windows",
			homeDir: filepath.Join("Users", "testuser"),
		},
		{
			name:    "Darwin - Stat Error",
			os:      "darwin",
			homeDir: filepath.Join("Users", "testuser"),
		},
		{
			name:    "Linux - Stat Error",
			os:      "linux",
			homeDir: filepath.Join("home", "testuser"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &mocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockConfigFactory := &mocks.MockConfigFactory{}
			defer mockConfigFactory.AssertExpectations(t)

			mockConfig := &configmocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			mockRootStore := &mocks.MockRootStore{}
			defer mockRootStore.AssertExpectations(t)

			mockRootPathResolver := &mocks.MockRootPathResolver{}
			defer mockRootPathResolver.AssertExpectations(t)

			mockLogger := testutils.NewInspectableLogger()

			expectedDir := tc.homeDir
			expectedError := assert.AnError
			if tc.os == "windows" || tc.os == "darwin" {
				expectedDir = filepath.Join(tc.homeDir, "Documents")
			}

			mockConfigFactory.EXPECT().
				Config().
				Return(mockConfig, nil).
				Once()

			mockConfig.EXPECT().
				PreferredMATLABStartingDirectory().
				Return("").
				Once()

			mockRootStore.EXPECT().
				GetRoots().
				Return([]entities.MCPRoot{}).
				Once()

			mockOSLayer.EXPECT().
				UserHomeDir().
				Return(tc.homeDir, nil).
				Once()

			mockOSLayer.EXPECT().
				GOOS().
				Return(tc.os).
				Once()

			mockOSLayer.EXPECT().
				Stat(expectedDir).
				Return(nil, expectedError).
				Once()

			selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

			// Act
			result, err := selector.SelectMATLABStartingDir(mockLogger)

			// Assert
			require.ErrorIs(t, err, expectedError)
			assert.Empty(t, result)
		})
	}
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_StatErrorOnPreferredMATLABStartingDir(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedPreferredMATLABStartingDir := filepath.Join("some", "path", "that", "doesnt", "exist")
	expectedError := assert.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return(expectedPreferredMATLABStartingDir).
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedPreferredMATLABStartingDir).
		Return(nil, expectedError).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_UsesFirstRootDirHappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedDir := filepath.Join("C:", "Users", "project")
	rootURI := "file:///C:/Users/project"

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockRootStore.EXPECT().
		GetRoots().
		Return([]entities.MCPRoot{entities.NewMCPRoot(rootURI, "")}).
		Once()

	mockRootPathResolver.EXPECT().
		Resolve(entities.NewMCPRoot(rootURI, "")).
		Return(expectedDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedDir).
		Return(mockFileInfo, nil).
		Once()

	mockFileInfo.EXPECT().
		IsDir().
		Return(true).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedDir, result)
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_RootStatErrorFallsBackToDefault(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	rootURI := "file:///C:/Users/project"
	rootDir := filepath.Join("C:", "Users", "project")
	homeDir := filepath.Join("Users", "testuser")
	expectedDir := filepath.Join("Users", "testuser", "Documents")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockRootStore.EXPECT().
		GetRoots().
		Return([]entities.MCPRoot{entities.NewMCPRoot(rootURI, "")}).
		Once()

	mockRootPathResolver.EXPECT().
		Resolve(entities.NewMCPRoot(rootURI, "")).
		Return(rootDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(rootDir).
		Return(nil, assert.AnError).
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(homeDir, nil).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedDir).
		Return(mockFileInfo, nil).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedDir, result)
	_, hasWarnLog := mockLogger.WarnLogs()["failed to use MCP root as starting directory, falling back to default"]
	assert.True(t, hasWarnLog, "expected warning log about falling back to default")
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_RootIsNotDirectoryFallsBackToDefault(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockRootFileInfo := &osFacademocks.MockFileInfo{}
	defer mockRootFileInfo.AssertExpectations(t)

	mockDocFileInfo := &osFacademocks.MockFileInfo{}
	defer mockDocFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	rootURI := "file:///C:/Users/project"
	rootDir := filepath.Join("C:", "Users", "project")
	homeDir := filepath.Join("Users", "testuser")
	expectedDir := filepath.Join("Users", "testuser", "Documents")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockRootStore.EXPECT().
		GetRoots().
		Return([]entities.MCPRoot{entities.NewMCPRoot(rootURI, "")}).
		Once()

	mockRootPathResolver.EXPECT().
		Resolve(entities.NewMCPRoot(rootURI, "")).
		Return(rootDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(rootDir).
		Return(mockRootFileInfo, nil).
		Once()

	mockRootFileInfo.EXPECT().
		IsDir().
		Return(false).
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(homeDir, nil).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedDir).
		Return(mockDocFileInfo, nil).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedDir, result)
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_RootUNCPathFallsBackToDefault(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	rootURI := "file://server/share"
	homeDir := filepath.Join("Users", "testuser")
	expectedDir := filepath.Join("Users", "testuser", "Documents")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockRootStore.EXPECT().
		GetRoots().
		Return([]entities.MCPRoot{entities.NewMCPRoot(rootURI, "")}).
		Once()

	mockRootPathResolver.EXPECT().
		Resolve(entities.NewMCPRoot(rootURI, "")).
		Return("", assert.AnError).
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(homeDir, nil).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedDir).
		Return(mockFileInfo, nil).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedDir, result)
	_, hasWarnLog := mockLogger.WarnLogs()["failed to use MCP root as starting directory, falling back to default"]
	assert.True(t, hasWarnLog, "expected warning log about falling back to default")
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_RootInvalidURIFallsBackToDefault(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	homeDir := filepath.Join("Users", "testuser")
	expectedDir := filepath.Join("Users", "testuser", "Documents")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockRootStore.EXPECT().
		GetRoots().
		Return([]entities.MCPRoot{entities.NewMCPRoot("://invalid", "")}).
		Once()

	mockRootPathResolver.EXPECT().
		Resolve(entities.NewMCPRoot("://invalid", "")).
		Return("", assert.AnError).
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(homeDir, nil).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedDir).
		Return(mockFileInfo, nil).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedDir, result)
	_, hasWarnLog := mockLogger.WarnLogs()["failed to use MCP root as starting directory, falling back to default"]
	assert.True(t, hasWarnLog, "expected warning log about falling back to default")
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_RootFileURIWithoutDriveLetterOnWindowsFallsBackToDefault(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	homeDir := filepath.Join("Users", "testuser")
	expectedDir := filepath.Join("Users", "testuser", "Documents")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockRootStore.EXPECT().
		GetRoots().
		Return([]entities.MCPRoot{entities.NewMCPRoot("file:///home/user/project", "")}).
		Once()

	mockRootPathResolver.EXPECT().
		Resolve(entities.NewMCPRoot("file:///home/user/project", "")).
		Return("", nil).
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(homeDir, nil).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedDir).
		Return(mockFileInfo, nil).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedDir, result)
}

func TestMATLABStartingDirSelector_SelectMATLABStartingDir_RootNonFileSchemeFallsBackToDefault(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockRootPathResolver := &mocks.MockRootPathResolver{}
	defer mockRootPathResolver.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	homeDir := filepath.Join("Users", "testuser")
	expectedDir := filepath.Join("Users", "testuser", "Documents")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockRootStore.EXPECT().
		GetRoots().
		Return([]entities.MCPRoot{entities.NewMCPRoot("https://example.com/repo", "")}).
		Once()

	mockRootPathResolver.EXPECT().
		Resolve(entities.NewMCPRoot("https://example.com/repo", "")).
		Return("", nil).
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(homeDir, nil).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedDir).
		Return(mockFileInfo, nil).
		Once()

	selector := matlabstartingdirselector.New(mockConfigFactory, mockOSLayer, mockRootStore, mockRootPathResolver)

	// Act
	result, err := selector.SelectMATLABStartingDir(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedDir, result)
}
