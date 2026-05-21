// Copyright 2025-2026 The MathWorks, Inc.

package pathvalidator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/usecases/utils/pathvalidator"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/usecases/utils/pathvalidator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOsLayer := &mocks.MockOSLayer{}

	// Act
	validator := pathvalidator.New(mockOsLayer)

	// Assert
	assert.NotNil(t, validator, "New() should return a non-nil Validator")
}

func TestValidator_ValidateMATLABScript_HappyPath(t *testing.T) {
	// Arrange
	mockOsLayer := &mocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	validator := pathvalidator.New(mockOsLayer)

	testPath, absErr := filepath.Abs("test.m")
	require.NoError(t, absErr)

	mockOsLayer.EXPECT().
		Stat(testPath).
		Return(mockFileInfo, nil).
		Once()

	mockFileInfo.EXPECT().
		IsDir().
		Return(false).
		Once()

	// Act
	result, err := validator.ValidateMATLABScript(testPath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, testPath, result)
}

func TestValidator_ValidateMATLABScript_InvalidPath(t *testing.T) {
	folderPath, absErr := filepath.Abs("./")
	require.NoError(t, absErr)

	tests := []struct {
		name     string
		filePath string
	}{
		{
			name:     "MATLAB file with relative path",
			filePath: filepath.Join(".", "relative", "folder", "test.m"),
		},
		{
			name:     "Folder path",
			filePath: folderPath,
		},
		{
			name:     "Empty path",
			filePath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockOsLayer := &mocks.MockOSLayer{}
			defer mockOsLayer.AssertExpectations(t)

			validator := pathvalidator.New(mockOsLayer)

			// Act
			_, err := validator.ValidateMATLABScript(tt.filePath)

			// Assert
			require.Error(t, err)
		})
	}
}

func TestValidator_ValidateMATLABScript_NotMATLABScript(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "Non-.m file",
			fileName: "file.txt",
		},
		{
			name:     "File without extension",
			fileName: "noextension",
		},
		{
			name:     "Capital M",
			fileName: "script.M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockOsLayer := &mocks.MockOSLayer{}
			defer mockOsLayer.AssertExpectations(t)

			validator := pathvalidator.New(mockOsLayer)

			filePath, absErr := filepath.Abs(tt.fileName)
			require.NoError(t, absErr)

			// Act
			_, err := validator.ValidateMATLABScript(filePath)

			// Assert
			require.Error(t, err)
		})
	}
}

func TestValidator_ValidateMATLABScript_MPathIsAFolder(t *testing.T) {
	// Arrange
	mockOsLayer := &mocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	validator := pathvalidator.New(mockOsLayer)

	// path has .m extension to pass suffix check but is registered as a folder
	testPath, absErr := filepath.Abs("folder.m")
	require.NoError(t, absErr)

	mockOsLayer.EXPECT().
		Stat(testPath).
		Return(mockFileInfo, nil).
		Once()

	mockFileInfo.EXPECT().
		IsDir().
		Return(true).
		Once()

	// Act
	_, err := validator.ValidateMATLABScript(testPath)

	// Assert
	require.Error(t, err)
}

func TestValidator_ValidateMATLABScript_StatFails(t *testing.T) {
	// Arrange
	mockOsLayer := &mocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	testPath, absErr := filepath.Abs("test.m")
	require.NoError(t, absErr)

	mockOsLayer.EXPECT().
		Stat(testPath).
		Return(nil, os.ErrNotExist).
		Once()

	validator := pathvalidator.New(mockOsLayer)

	// Act
	_, err := validator.ValidateMATLABScript(testPath)

	// Assert
	require.Error(t, err)
}

func TestValidator_ValidateMATLABScript_PathCleaning(t *testing.T) {
	testPath, absErr := filepath.Abs("test.m")
	require.NoError(t, absErr)

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "Path with double slashes",
			filePath: strings.ReplaceAll(testPath, string(filepath.Separator), string(filepath.Separator)+string(filepath.Separator)),
			expected: testPath,
		},
		{
			name:     "Path with dot segments",
			filePath: filepath.Join(filepath.Dir(testPath), ".", "test.m"),
			expected: testPath,
		},
		{
			name:     "Path with dot-dot segments",
			filePath: filepath.Join(filepath.Dir(testPath), "subdir", "..", "test.m"),
			expected: testPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockOsLayer := &mocks.MockOSLayer{}
			defer mockOsLayer.AssertExpectations(t)

			mockFileInfo := &osfacademocks.MockFileInfo{}
			defer mockFileInfo.AssertExpectations(t)

			validator := pathvalidator.New(mockOsLayer)

			mockOsLayer.EXPECT().
				Stat(tt.expected).
				Return(mockFileInfo, nil).
				Once()

			mockFileInfo.EXPECT().
				IsDir().
				Return(false).
				Once()

			// Act
			result, err := validator.ValidateMATLABScript(tt.filePath)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidator_ValidateFolderPath_HappyPath(t *testing.T) {
	// Arrange
	mockOsLayer := &mocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	validator := pathvalidator.New(mockOsLayer)

	testPath, absErr := filepath.Abs("./")
	require.NoError(t, absErr)

	mockOsLayer.EXPECT().
		Stat(testPath).
		Return(mockFileInfo, nil).
		Once()

	mockFileInfo.EXPECT().
		IsDir().
		Return(true).
		Once()

	// Act
	result, err := validator.ValidateFolderPath(testPath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, testPath, result)
}

func TestValidator_ValidateFolderPath_FailsForRelativePath(t *testing.T) {
	// Arrange
	mockOsLayer := &mocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	validator := pathvalidator.New(mockOsLayer)

	testPath := filepath.Join(".", "relative", "folder")

	// Act
	_, err := validator.ValidateFolderPath(testPath)

	// Assert
	require.Error(t, err)
}

func TestValidator_ValidateFolderPath_FailsForFilePath(t *testing.T) {
	// Arrange
	mockOsLayer := &mocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	validator := pathvalidator.New(mockOsLayer)

	testPath, absErr := filepath.Abs("test.m")
	require.NoError(t, absErr)

	mockOsLayer.EXPECT().
		Stat(testPath).
		Return(mockFileInfo, nil).
		Once()

	mockFileInfo.EXPECT().
		IsDir().
		Return(false).
		Once()

	// Act
	_, err := validator.ValidateFolderPath(testPath)

	// Assert
	require.Error(t, err)
}

func TestValidator_ValidateFolderPath_StatFails(t *testing.T) {
	// Arrange
	mockOsLayer := &mocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	testPath, absErr := filepath.Abs("./")
	require.NoError(t, absErr)

	mockOsLayer.EXPECT().
		Stat(testPath).
		Return(nil, os.ErrNotExist).
		Once()

	validator := pathvalidator.New(mockOsLayer)

	// Act
	_, err := validator.ValidateFolderPath(testPath)

	// Assert
	require.Error(t, err)
}
