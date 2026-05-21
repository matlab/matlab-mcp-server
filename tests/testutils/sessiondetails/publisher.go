// Copyright 2026 The MathWorks, Inc.

package sessiondetails

import (
	"fmt"
	"os"
	"path/filepath"
)

func Publish(homeDir string, detailsJSON string) (string, error) {
	appDataDir := ResolveAppDataDir(homeDir)

	sessionDir := filepath.Join(appDataDir, "v1")
	if err := os.MkdirAll(sessionDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create session dir: %w", err)
	}

	sessionFilePath := ResolveSessionDetailsPath(homeDir)
	tmp, err := os.CreateTemp(sessionDir, "sessionDetails-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp session details file: %w", err)
	}
	tmpName := tmp.Name()
	cleanupTmp := true
	defer func() {
		if cleanupTmp {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmp.Write([]byte(detailsJSON)); err != nil {
		_ = tmp.Close()
		return "", fmt.Errorf("failed to write temp session details: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp session details file: %w", err)
	}

	if err := ensureFileSecure(tmpName); err != nil {
		return "", fmt.Errorf("failed to secure temp session details file: %w", err)
	}

	if err := os.Rename(tmpName, sessionFilePath); err != nil {
		if removeErr := os.Remove(sessionFilePath); removeErr == nil {
			err = os.Rename(tmpName, sessionFilePath)
		}
		if err != nil {
			return "", fmt.Errorf("failed to replace session details file: %w", err)
		}
	}
	cleanupTmp = false

	return appDataDir, nil
}

func Remove(homeDir string) error {
	return os.Remove(ResolveSessionDetailsPath(homeDir))
}
