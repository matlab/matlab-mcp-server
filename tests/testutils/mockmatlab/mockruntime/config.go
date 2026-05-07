// Copyright 2026 The MathWorks, Inc.

package mockruntime

import (
	"encoding/json"
	"fmt"
)

const (
	ModeHappy           = "happy"
	ModeHangBeforeFiles = "hang_before_files"
	ModeExitImmediately = "exit_immediately"
	ModeSlowStartup     = "slow_startup"
	ModeStartupFailure  = "startup_failure"

	EnvMockMATLABConfig = "MW_MCP_MOCK_MATLAB_CONFIG"
	EnvMockMATLABLogDir = "MW_MCP_MOCK_MATLAB_LOG_DIR"
)

type Config struct {
	Mode     string `json:"mode"`
	ExitCode *int   `json:"exitCode,omitempty"`
	DelayMs  *int   `json:"delayMs,omitempty"`
}

func (c Config) ToEnvValue() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func HappyConfig() Config {
	return Config{Mode: ModeHappy}
}

func HangBeforeFilesConfig() Config {
	return Config{Mode: ModeHangBeforeFiles}
}

func ExitImmediatelyConfig(exitCode int) Config {
	return Config{Mode: ModeExitImmediately, ExitCode: &exitCode}
}

func SlowStartupConfig(delayMs int) Config {
	return Config{Mode: ModeSlowStartup, DelayMs: &delayMs}
}

func StartupFailureConfig() Config {
	return Config{Mode: ModeStartupFailure}
}

func (r *Runtime) LoadConfigFromEnv() (Config, error) {
	raw := r.Env.Getenv(EnvMockMATLABConfig)

	if raw == "" {
		return Config{Mode: ModeHappy}, nil
	}

	var cfg Config
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse %s: %w", EnvMockMATLABConfig, err)
	}

	if cfg.Mode == "" {
		cfg.Mode = ModeHappy
	}

	return cfg, nil
}
