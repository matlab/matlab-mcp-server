// Copyright 2026 The MathWorks, Inc.

package mockruntime_test

import (
	"encoding/json"
	"testing"

	mocks "github.com/matlab/matlab-mcp-core-server/tests/mocks/testutils/mockmatlab/mockruntime"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigHelpers_CreateExpectedModes(t *testing.T) {
	// Arrange
	happy := mockruntime.HappyConfig()
	hangBeforeFiles := mockruntime.HangBeforeFilesConfig()
	exitImmediately := mockruntime.ExitImmediatelyConfig(3)
	slowStartup := mockruntime.SlowStartupConfig(250)
	startupFailure := mockruntime.StartupFailureConfig()

	// Assert
	assert.Equal(t, mockruntime.ModeHappy, happy.Mode)
	assert.Equal(t, mockruntime.ModeHangBeforeFiles, hangBeforeFiles.Mode)
	assert.Equal(t, mockruntime.ModeExitImmediately, exitImmediately.Mode)
	require.NotNil(t, exitImmediately.ExitCode)
	assert.Equal(t, 3, *exitImmediately.ExitCode)
	assert.Equal(t, mockruntime.ModeSlowStartup, slowStartup.Mode)
	require.NotNil(t, slowStartup.DelayMs)
	assert.Equal(t, 250, *slowStartup.DelayMs)
	assert.Equal(t, mockruntime.ModeStartupFailure, startupFailure.Mode)
}

func TestConfig_ToEnvValue_EncodesConfigAsJSON(t *testing.T) {
	// Arrange
	cfg := mockruntime.SlowStartupConfig(500)

	// Act
	value, err := cfg.ToEnvValue()

	// Assert
	require.NoError(t, err)
	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(value), &parsed))
	assert.Equal(t, mockruntime.ModeSlowStartup, parsed["mode"])
	delayMs, ok := parsed["delayMs"].(float64)
	require.True(t, ok)
	assert.InDelta(t, 500, delayMs, 0)
}

func TestLoadConfigFromEnv_EmptyEnv_ReturnsDefaultHappyMode(t *testing.T) {
	// Arrange
	env := mocks.NewMockEnvironment(t)
	env.EXPECT().Getenv(mockruntime.EnvMockMATLABConfig).Return("")
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
	assert.Equal(t, mockruntime.ModeHappy, cfg.Mode)
	assert.Nil(t, cfg.ExitCode)
	assert.Nil(t, cfg.DelayMs)
}

func TestLoadConfigFromEnv_InvalidJSON_ReturnsError(t *testing.T) {
	// Arrange
	env := mocks.NewMockEnvironment(t)
	env.EXPECT().Getenv(mockruntime.EnvMockMATLABConfig).Return("{invalid-json")
	fileSystem := mocks.NewMockFileSystem(t)
	tlsProvider := mocks.NewMockTLSMaterialProvider(t)

	runtime := mockruntime.NewRuntime(
		env,
		fileSystem,
		tlsProvider,
	)

	// Act
	_, err := runtime.LoadConfigFromEnv()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), mockruntime.EnvMockMATLABConfig)
}

func TestLoadConfigFromEnv_EmptyMode_DefaultsToHappyMode(t *testing.T) {
	// Arrange
	env := mocks.NewMockEnvironment(t)
	env.EXPECT().Getenv(mockruntime.EnvMockMATLABConfig).Return(`{"delayMs":250}`)
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
	assert.Equal(t, mockruntime.ModeHappy, cfg.Mode)
	require.NotNil(t, cfg.DelayMs)
	assert.Equal(t, 250, *cfg.DelayMs)
}
