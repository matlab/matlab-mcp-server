// Copyright 2026 The MathWorks, Inc.

package sessiondetails

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

func MarshalJSON(port string, certificatePath, apiKey string, pid int) (string, error) {
	portValue, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("invalid port %q: %w", port, err)
	}

	details := map[string]any{
		"port":        portValue,
		"certificate": certificatePath,
		"apiKey":      apiKey,
		"pid":         pid,
	}
	data, err := json.Marshal(details)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session details: %w", err)
	}
	return string(data), nil
}

func ResolveAppDataDir(homeDir string) string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "AppData", "Roaming", "MathWorks", "MATLAB MCP Core Server")
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "MathWorks", "MATLAB MCP Core Server")
	default:
		return filepath.Join(homeDir, ".MathWorks", "MATLABMCPCoreServer")
	}
}

func ResolveSessionDetailsPath(homeDir string) string {
	return filepath.Join(ResolveAppDataDir(homeDir), "v1", "sessionDetails.json")
}
