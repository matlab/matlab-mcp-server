// Copyright 2026 The MathWorks, Inc.

package mcpb_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcpb/mcpbstagebuilder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuild_HappyPath(t *testing.T) {
	// Arrange
	stagingDir := filepath.Join(t.TempDir(), "staging")
	t.Setenv("MCPB_STAGING_DIR", stagingDir)

	// Act
	err := mcpbstagebuilder.Build("0.0.0-test")

	// Assert
	require.NoError(t, err)

	assertPackageJSONStaged(t, stagingDir)
	assertManifestStaged(t, stagingDir)
	assertStaticAssetsStaged(t, stagingDir)
}

func TestBuild_EnvVarNotSet(t *testing.T) {
	// Arrange
	t.Setenv("MCPB_STAGING_DIR", "")
	require.NoError(t, os.Unsetenv("MCPB_STAGING_DIR"))

	// Act
	err := mcpbstagebuilder.Build("0.0.0-test")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MCPB_STAGING_DIR")
}

func TestBuild_FailsIfStagingDirExists(t *testing.T) {
	// Arrange
	stagingDir := filepath.Join(t.TempDir(), "staging")
	t.Setenv("MCPB_STAGING_DIR", stagingDir)

	require.NoError(t, mcpbstagebuilder.Build("0.0.0-test"))

	// Act
	err := mcpbstagebuilder.Build("0.0.0-test")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestBuild_ManifestPreservesTemplateFields(t *testing.T) {
	// Arrange
	stagingDir := filepath.Join(t.TempDir(), "staging")
	t.Setenv("MCPB_STAGING_DIR", stagingDir)

	// Act
	require.NoError(t, mcpbstagebuilder.Build("0.0.0-test"))

	// Assert
	content, err := os.ReadFile(filepath.Join(stagingDir, "bundle", "manifest.json")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)

	var manifest map[string]any
	require.NoError(t, json.Unmarshal(content, &manifest))

	assert.Equal(t, "0.2", manifest["manifest_version"])
	assert.Equal(t, "matlab-mcp-core-server", manifest["name"])
	assert.Equal(t, "MATLAB", manifest["display_name"])
	assert.NotEmpty(t, manifest["description"])
	assert.NotEmpty(t, manifest["long_description"])

	server, ok := manifest["server"].(map[string]any)
	require.True(t, ok)

	mcpConfig, ok := server["mcp_config"].(map[string]any)
	require.True(t, ok)
	assert.NotEmpty(t, mcpConfig["command"])
	assert.NotEmpty(t, mcpConfig["env"])
	assert.NotEmpty(t, mcpConfig["platform_overrides"])
}

func TestBuild_LaunchersSetMCPBContextTags(t *testing.T) {
	// Arrange
	stagingDir := filepath.Join(t.TempDir(), "staging")
	t.Setenv("MCPB_STAGING_DIR", stagingDir)

	expectedShellTag := "export MW_CONTEXT_TAGS=\"${MW_CONTEXT_TAGS:+$MW_CONTEXT_TAGS,}MATLAB:MCPB:V1\""
	expectedCmdTagWithExisting := "set \"MW_CONTEXT_TAGS=%MW_CONTEXT_TAGS%,MATLAB:MCPB:V1\""
	expectedCmdTagNew := "set \"MW_CONTEXT_TAGS=MATLAB:MCPB:V1\""

	// Act
	require.NoError(t, mcpbstagebuilder.Build("0.0.0-test"))

	// Assert
	shContent, err := os.ReadFile(filepath.Join(stagingDir, "bundle", "bin", "launch-matlab-mcp.sh")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)
	assert.Contains(t, string(shContent), expectedShellTag)

	cmdContent, err := os.ReadFile(filepath.Join(stagingDir, "bundle", "bin", "launch-matlab-mcp.cmd")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)
	cmdText := string(cmdContent)
	assert.Contains(t, cmdText, expectedCmdTagWithExisting)
	assert.Contains(t, cmdText, expectedCmdTagNew)
}

func assertPackageJSONStaged(t *testing.T, stagingDir string) {
	t.Helper()

	content, err := os.ReadFile(filepath.Join(stagingDir, "package.json")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)

	var pkg map[string]any
	require.NoError(t, json.Unmarshal(content, &pkg))
	assert.Equal(t, "matlab-mcp-core-server-build", pkg["name"])
}

func assertManifestStaged(t *testing.T, stagingDir string) {
	t.Helper()

	content, err := os.ReadFile(filepath.Join(stagingDir, "bundle", "manifest.json")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)

	var manifest map[string]any
	require.NoError(t, json.Unmarshal(content, &manifest))

	assert.NotEmpty(t, manifest["version"])

	toolsRaw, ok := manifest["tools"].([]any)
	require.True(t, ok)
	assert.Len(t, toolsRaw, 5)

	for _, raw := range toolsRaw {
		tool, ok := raw.(map[string]any)
		require.True(t, ok)
		assert.NotEmpty(t, tool["name"])
		assert.NotEmpty(t, tool["description"])
	}

	userConfigRaw, ok := manifest["user_config"].(map[string]any)
	require.True(t, ok)
	assert.NotEmpty(t, userConfigRaw)

	assertUserConfigReferencedInMCPConfig(t, manifest, userConfigRaw)
}

func assertStaticAssetsStaged(t *testing.T, stagingDir string) {
	t.Helper()

	iconPath := filepath.Join(stagingDir, "bundle", "icon.png")
	info, err := os.Stat(iconPath)
	require.NoError(t, err)
	assert.Positive(t, info.Size())

	launcherSh := filepath.Join(stagingDir, "bundle", "bin", "launch-matlab-mcp.sh")
	launcherCmd := filepath.Join(stagingDir, "bundle", "bin", "launch-matlab-mcp.cmd")
	assertLauncherPermissions(t, launcherSh, launcherCmd)
}

func assertUserConfigReferencedInMCPConfig(t *testing.T, manifest map[string]any, userConfig map[string]any) {
	t.Helper()

	server, ok := manifest["server"].(map[string]any)
	require.True(t, ok, "server section not found in manifest")

	mcpConfig, ok := server["mcp_config"].(map[string]any)
	require.True(t, ok, "mcp_config section not found in manifest")

	env, ok := mcpConfig["env"].(map[string]any)
	require.True(t, ok, "env section not found in manifest")

	const prefix = "${user_config."
	const suffix = "}"

	referencedKeys := make(map[string]bool)

	// Check env vars for user_config references
	for envVar, value := range env {
		valueStr, ok := value.(string)
		require.True(t, ok, "env var %s value is not a string", envVar)

		require.True(t, strings.HasPrefix(valueStr, prefix) && strings.HasSuffix(valueStr, suffix),
			"env var %s should reference user_config, got %s", envVar, valueStr)

		key := valueStr[len(prefix) : len(valueStr)-len(suffix)]
		referencedKeys[key] = true

		assert.Contains(t, userConfig, key,
			"env var %s references ${user_config.%s} but no such user_config key exists", envVar, key)
	}

	// Check args for user_config references
	if args, ok := mcpConfig["args"].([]any); ok {
		for _, arg := range args {
			argStr, ok := arg.(string)
			if !ok {
				continue
			}
			if strings.HasPrefix(argStr, prefix) && strings.HasSuffix(argStr, suffix) {
				key := argStr[len(prefix) : len(argStr)-len(suffix)]
				referencedKeys[key] = true

				assert.Contains(t, userConfig, key,
					"args references ${user_config.%s} but no such user_config key exists", key)
			}
		}
	}

	for key := range userConfig {
		assert.True(t, referencedKeys[key],
			"user_config key %s is not referenced by any env var or arg", key)
	}
}

func TestBuild_EnvVarConsistencyAcrossManifestAndLaunchers(t *testing.T) {
	// Arrange
	stagingDir := filepath.Join(t.TempDir(), "staging")
	t.Setenv("MCPB_STAGING_DIR", stagingDir)

	// Act
	require.NoError(t, mcpbstagebuilder.Build("0.0.0-test"))

	// Assert — extract env var names from manifest.json
	manifestContent, err := os.ReadFile(filepath.Join(stagingDir, "bundle", "manifest.json")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)

	var manifest map[string]any
	require.NoError(t, json.Unmarshal(manifestContent, &manifest))

	server, ok := manifest["server"].(map[string]any)
	require.True(t, ok, "server section not found in manifest")

	mcpConfig, ok := server["mcp_config"].(map[string]any)
	require.True(t, ok, "mcp_config section not found in manifest")

	env, ok := mcpConfig["env"].(map[string]any)
	require.True(t, ok, "env section not found in manifest")

	var manifestEnvVars []string
	for key := range env {
		manifestEnvVars = append(manifestEnvVars, key)
	}
	sort.Strings(manifestEnvVars)

	// Assert — extract env var names from launch-matlab-mcp.sh
	shContent, err := os.ReadFile(filepath.Join(stagingDir, "bundle", "bin", "launch-matlab-mcp.sh")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)
	shEnvVars := extractEnvVarsFromContent(string(shContent))

	// Assert — extract env var names from launch-matlab-mcp.cmd
	cmdContent, err := os.ReadFile(filepath.Join(stagingDir, "bundle", "bin", "launch-matlab-mcp.cmd")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)
	cmdEnvVars := extractEnvVarsFromContent(string(cmdContent))

	// Assert — all three lists should be equal
	assert.Equal(t, manifestEnvVars, shEnvVars,
		"env var names in manifest.json and launch-matlab-mcp.sh should match")
	assert.Equal(t, manifestEnvVars, cmdEnvVars,
		"env var names in manifest.json and launch-matlab-mcp.cmd should match")
}

func TestBuild_LaunchersParseExtensionFilesArgs(t *testing.T) {
	// Arrange
	stagingDir := filepath.Join(t.TempDir(), "staging")
	t.Setenv("MCPB_STAGING_DIR", stagingDir)

	// Act
	require.NoError(t, mcpbstagebuilder.Build("0.0.0-test"))

	// Assert — bash script handles --extension-files and forwards as --extension-file
	shContent, err := os.ReadFile(filepath.Join(stagingDir, "bundle", "bin", "launch-matlab-mcp.sh")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)
	shText := string(shContent)
	assert.Contains(t, shText, "--extension-files)")
	assert.Contains(t, shText, `ARGS+=("--extension-file" "$1")`)

	// Assert — cmd script handles --extension-files and forwards as --extension-file
	cmdContent, err := os.ReadFile(filepath.Join(stagingDir, "bundle", "bin", "launch-matlab-mcp.cmd")) //nolint:gosec // Test file path from t.TempDir()
	require.NoError(t, err)
	cmdText := string(cmdContent)
	assert.Contains(t, cmdText, `"--extension-files"`)
	assert.Contains(t, cmdText, `--extension-file "%~1"`)
}

func extractEnvVarsFromContent(content string) []string {
	re := regexp.MustCompile(`__MATLAB_MCP_CORE_SERVER_MCPB_\w+`)
	matches := re.FindAllString(content, -1)

	seen := make(map[string]bool)
	var unique []string
	for _, m := range matches {
		if !seen[m] {
			seen[m] = true
			unique = append(unique, m)
		}
	}

	sort.Strings(unique)
	return unique
}
