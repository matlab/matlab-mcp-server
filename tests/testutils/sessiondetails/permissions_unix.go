// Copyright 2026 The MathWorks, Inc.

//go:build !windows

package sessiondetails

import (
	"fmt"
	"os"
)

func ensureFileSecure(path string) error {
	return os.Chmod(path, 0o600)
}

func AssertFileSecure(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat session details file: %w", err)
	}

	if fileInfo.Mode().Perm() != 0o600 {
		return fmt.Errorf("expected permissions 0600, got %04o", fileInfo.Mode().Perm())
	}

	return nil
}
