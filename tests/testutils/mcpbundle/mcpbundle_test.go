// Copyright 2026 The MathWorks, Inc.

package mcpbundle_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	mocks "github.com/matlab/matlab-mcp-core-server/tests/mocks/testutils/mcpbundle"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpbundle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestManifestVersion_HappyPath(t *testing.T) {
	bundleDir := createBundleDir(t, map[string]string{
		"manifest.json": `{"version": "1.2.3"}`,
	})
	bundle := mcpbundle.NewBundleForTest(bundleDir, nil)

	version, err := bundle.ManifestVersion()

	require.NoError(t, err)
	assert.Equal(t, "1.2.3", version)
}

func TestManifestVersion_MissingFile(t *testing.T) {
	bundle := mcpbundle.NewBundleForTest(t.TempDir(), nil)

	_, err := bundle.ManifestVersion()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "reading manifest.json")
}

func TestManifestVersion_InvalidJSON(t *testing.T) {
	bundleDir := createBundleDir(t, map[string]string{
		"manifest.json": `not json`,
	})
	bundle := mcpbundle.NewBundleForTest(bundleDir, nil)

	_, err := bundle.ManifestVersion()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing manifest.json")
}

func TestManifestVersion_VersionNotString(t *testing.T) {
	bundleDir := createBundleDir(t, map[string]string{
		"manifest.json": `{"version": 123}`,
	})
	bundle := mcpbundle.NewBundleForTest(bundleDir, nil)

	_, err := bundle.ManifestVersion()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a string")
}

func TestLauncherFlags_MissingFile(t *testing.T) {
	bundle := mcpbundle.NewBundleForTest(t.TempDir(), nil)

	_, err := bundle.LauncherFlags()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "reading launcher script")
}

func TestLauncherFlags_EmptyScript(t *testing.T) {
	bundleDir := createBundleDir(t, map[string]string{
		"bin/" + mcpbundle.LauncherFilename(): "#!/bin/bash\necho hello",
	})
	bundle := mcpbundle.NewBundleForTest(bundleDir, nil)

	flags, err := bundle.LauncherFlags()

	require.NoError(t, err)
	assert.Empty(t, flags)
}

func TestLaunch_PassesEnvVarsToRunner(t *testing.T) {
	bundleDir := createBundleDir(t, map[string]string{
		"bin/" + mcpbundle.LauncherFilename(): "#!/bin/bash",
	})

	runner := mocks.NewMockCommandRunner(t)

	runner.EXPECT().
		Run(t.Context(), mock.AnythingOfType("string"), mock.MatchedBy(func(env []string) bool {
			return slices.Contains(env, "__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_ROOT=/opt/matlab")
		}), []string(nil)).
		Return(mcpbundle.LaunchResult{Args: []string{"--matlab-root", "/opt/matlab"}}, nil).
		Once()

	bundle := mcpbundle.NewBundleForTest(bundleDir, runner)

	result := bundle.Launch(t, map[string]string{
		"__MATLAB_MCP_CORE_SERVER_MCPB_MATLAB_ROOT": "/opt/matlab",
	})

	assert.Equal(t, []string{"--matlab-root", "/opt/matlab"}, result.Args)
}

func TestLaunch_FiltersExistingMCPBEnvVars(t *testing.T) {
	bundleDir := createBundleDir(t, map[string]string{
		"bin/" + mcpbundle.LauncherFilename(): "#!/bin/bash",
	})

	runner := mocks.NewMockCommandRunner(t)

	runner.EXPECT().
		Run(t.Context(), mock.AnythingOfType("string"), mock.MatchedBy(func(env []string) bool {
			for _, e := range env {
				if strings.HasPrefix(e, "__MATLAB_MCP_CORE_SERVER_MCPB_") {
					return false
				}
			}
			return true
		}), []string(nil)).
		Return(mcpbundle.LaunchResult{}, nil).
		Once()

	bundle := mcpbundle.NewBundleForTest(bundleDir, runner)

	bundle.Launch(t, nil)
}

func createBundleDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0750))
		require.NoError(t, os.WriteFile(path, []byte(content), 0600))
	}
	return dir
}
