// Copyright 2026 The MathWorks, Inc.

package appdatadir

import (
	"fmt"
	"path/filepath"
)

const (
	appDataDirNameLinux   = "MATLABMCPCoreServer"
	appDataDirNameDarwin  = "MATLAB MCP Core Server"
	appDataDirNameWindows = "MATLAB MCP Core Server"
)

type OSLayer interface {
	UserHomeDir() (string, error)
	Getenv(key string) string
	GOOS() string
}

type Getter struct {
	osLayer OSLayer
}

func New(osLayer OSLayer) *Getter {
	return &Getter{
		osLayer: osLayer,
	}
}

func (g *Getter) AppDataDir() (string, error) {
	goos := g.osLayer.GOOS()

	switch goos {
	case "linux":
		home, err := g.osLayer.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		return filepath.Join(home, ".MathWorks", appDataDirNameLinux), nil
	case "darwin":
		home, err := g.osLayer.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		return filepath.Join(home, "Library", "Application Support", "MathWorks", appDataDirNameDarwin), nil
	case "windows":
		appData := g.osLayer.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable is not set")
		}
		return filepath.Join(appData, "MathWorks", appDataDirNameWindows), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", goos)
	}
}
