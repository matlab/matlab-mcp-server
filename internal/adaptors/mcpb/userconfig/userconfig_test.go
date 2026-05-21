// Copyright 2026 The MathWorks, Inc.

package userconfig_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcpb/userconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserConfig_HappyPath(t *testing.T) {
	// Arrange

	// Act
	config, err := userconfig.GetUserConfig()

	// Assert
	require.NoError(t, err)
	require.Len(t, config, 8)

	expectedKeys := []string{
		"PreferredLocalMATLABRoot",
		"PreferredMATLABStartingDirectory",
		"BaseDir",
		"InitializeMATLABOnStartup",
		"DisableTelemetry",
		"MATLABDisplayMode",
		"MATLABSessionMode",
		"ExtensionFile",
	}

	for _, key := range expectedKeys {
		entry, exists := config[key]
		if assert.True(t, exists, "expected key %s not found", key) {
			assert.NotEmpty(t, entry.Type, "entry %s should have a type", key)
			assert.NotEmpty(t, entry.Title, "entry %s should have a title", key)
			assert.NotEmpty(t, entry.Description, "entry %s should have a description", key)
		}
	}
}

func TestGetUserConfig_TypesAreValid(t *testing.T) {
	// Arrange
	validTypes := map[string]bool{"string": true, "boolean": true, "directory": true, "file": true}

	// Act
	config, err := userconfig.GetUserConfig()

	// Assert
	require.NoError(t, err)

	for key, entry := range config {
		assert.True(t, validTypes[entry.Type], "entry %s has invalid type %s", key, entry.Type)
	}
}

func TestGetUserConfig_DirectoryTypeOverrides(t *testing.T) {
	// Arrange
	expectedDirectoryKeys := []string{
		"PreferredLocalMATLABRoot",
		"PreferredMATLABStartingDirectory",
		"BaseDir",
	}

	// Act
	config, err := userconfig.GetUserConfig()

	// Assert
	require.NoError(t, err)

	for _, key := range expectedDirectoryKeys {
		entry, exists := config[key]
		if assert.True(t, exists, "expected key %s not found", key) {
			assert.Equal(t, "directory", entry.Type, "entry %s should have type 'directory'", key)
		}
	}
}

func TestGetUserConfig_FileTypeOverrides(t *testing.T) {
	// Arrange
	expectedFileKeys := []string{
		"ExtensionFile",
	}

	// Act
	config, err := userconfig.GetUserConfig()

	// Assert
	require.NoError(t, err)

	for _, key := range expectedFileKeys {
		entry, exists := config[key]
		if assert.True(t, exists, "expected key %s not found", key) {
			assert.Equal(t, "file", entry.Type, "entry %s should have type 'file'", key)
		}
	}
}

func TestGetUserConfig_InferredTypes(t *testing.T) {
	// Arrange

	// Act
	config, err := userconfig.GetUserConfig()

	// Assert
	require.NoError(t, err)

	assert.Equal(t, "string", config["MATLABDisplayMode"].Type)
	assert.Equal(t, "boolean", config["InitializeMATLABOnStartup"].Type)
	assert.Equal(t, "boolean", config["DisableTelemetry"].Type)
}

func TestGetUserConfig_DefaultValues(t *testing.T) {
	// Arrange

	// Act
	config, err := userconfig.GetUserConfig()

	// Assert
	require.NoError(t, err)

	assert.Empty(t, config["PreferredLocalMATLABRoot"].Default)
	assert.Empty(t, config["PreferredMATLABStartingDirectory"].Default)
	assert.Empty(t, config["BaseDir"].Default)
	assert.Equal(t, false, config["InitializeMATLABOnStartup"].Default)
	assert.Equal(t, false, config["DisableTelemetry"].Default)
	assert.Equal(t, "desktop", config["MATLABDisplayMode"].Default)
}

func TestGetUserConfig_AllEntriesNotRequired(t *testing.T) {
	// Arrange

	// Act
	config, err := userconfig.GetUserConfig()

	// Assert
	require.NoError(t, err)

	for key, entry := range config {
		assert.False(t, entry.Required, "entry %s should not be required", key)
	}
}
