// Copyright 2026 The MathWorks, Inc.

package mcpbundle

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testconfig"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpbundle/launcherflags"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpbundle/launchersyntax"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpbundle/unpack"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmcpbinary"
	"github.com/stretchr/testify/require"
)

type CommandRunner interface {
	Run(ctx context.Context, launcherPath string, env []string, args []string) (LaunchResult, error)
}

type Bundle struct {
	dir           string
	commandRunner CommandRunner
}

func Open(t *testing.T) *Bundle {
	t.Helper()

	mcpbPath := os.Getenv("MCPB_ARTIFACT_PATH")
	require.NotEmpty(t, mcpbPath, "MCPB_ARTIFACT_PATH must be set")

	f, err := os.Open(mcpbPath) //nolint:gosec // Path from trusted env var
	require.NoError(t, err)
	defer func() { _ = f.Close() }()

	info, err := f.Stat()
	require.NoError(t, err)

	bundleDir := filepath.Join(t.TempDir(), "bundle")
	u := unpack.New()
	require.NoError(t, u.Unpack(f, info.Size(), bundleDir))

	return &Bundle{dir: bundleDir, commandRunner: &shellCommandRunner{baseDir: t.TempDir()}}
}

func newBundle(bundleDir string, commandRunner CommandRunner) *Bundle {
	return &Bundle{dir: bundleDir, commandRunner: commandRunner}
}

func (b *Bundle) Dir() string {
	return b.dir
}

func (b *Bundle) BinaryPath() string {
	return filepath.Join(b.dir, "bin", "matlab-mcp-core-server-"+testconfig.OSDescriptor+testconfig.ExecutableExtension)
}

func (b *Bundle) ManifestVersion() (string, error) {
	content, err := os.ReadFile(filepath.Join(b.dir, "manifest.json"))
	if err != nil {
		return "", fmt.Errorf("reading manifest.json: %w", err)
	}

	var manifest map[string]any
	if err := json.Unmarshal(content, &manifest); err != nil {
		return "", fmt.Errorf("parsing manifest.json: %w", err)
	}

	version, ok := manifest["version"].(string)
	if !ok {
		return "", fmt.Errorf("manifest version is not a string")
	}
	return version, nil
}

func (b *Bundle) LauncherFlags() ([]string, error) {
	content, err := os.ReadFile(filepath.Join(b.dir, "bin", launcherFilename))
	if err != nil {
		return nil, fmt.Errorf("reading launcher script: %w", err)
	}
	return launcherflags.Parse(string(content)), nil
}

type LaunchResult struct {
	Args []string
	Env  []string
}

func (b *Bundle) Launch(t *testing.T, envVars map[string]string, args ...string) LaunchResult {
	t.Helper()

	launcher := filepath.Join(b.dir, "bin", launcherFilename)
	env := filterMCPBEnvVars(os.Environ())
	for k, v := range envVars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	result, err := b.commandRunner.Run(t.Context(), launcher, env, args)
	require.NoError(t, err, "launcher failed")

	return result
}

func (b *Bundle) CheckLauncherSyntax() error {
	return launchersyntax.Check(filepath.Join(b.dir, "bin", launcherFilename))
}

func PathWithSpaces() string {
	return pathWithSpaces
}

func ExcludedFlags() map[string]string {
	return map[string]string{
		"help":           "meta flag, not a user configuration",
		"version":        "meta flag, not a user configuration",
		"setup-matlab":   "internal operational mode, not user-facing in MCPB",
		"extension-file": "passed via manifest args array, not launcher env var mapping",
	}
}

type shellCommandRunner struct {
	baseDir string
	seq     int
}

func (r *shellCommandRunner) Run(ctx context.Context, launcherPath string, env []string, args []string) (LaunchResult, error) {
	r.seq++
	outputFile := filepath.Join(r.baseDir, fmt.Sprintf("invocation-%d.jsonl", r.seq))

	env = append(env, fmt.Sprintf("%s=%s", mockmcpbinary.OutputFileEnvVar, outputFile))

	cmd := execLauncherCommand(ctx, launcherPath, args...)
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		return LaunchResult{}, fmt.Errorf("%s: %w", output, err)
	}

	record := &mockmcpbinary.Installation{OutputFile: outputFile}
	invocations, err := record.Invocations()
	if err != nil {
		return LaunchResult{}, fmt.Errorf("reading invocations: %w", err)
	}
	if len(invocations) != 1 {
		return LaunchResult{}, fmt.Errorf("expected 1 invocation, got %d", len(invocations))
	}

	return LaunchResult{
		Args: invocations[0].Args,
		Env:  invocations[0].Env,
	}, nil
}

func filterMCPBEnvVars(env []string) []string {
	filtered := make([]string, 0, len(env))
	for _, e := range env {
		if !strings.HasPrefix(e, "__MATLAB_MCP_CORE_SERVER_MCPB_") {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
