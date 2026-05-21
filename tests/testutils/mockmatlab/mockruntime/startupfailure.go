// Copyright 2026 The MathWorks, Inc.

package mockruntime

import (
	"fmt"
	"path/filepath"
)

func (r *Runtime) WriteStartupFailureFile(sessionDir string) error {
	if sessionDir == "" {
		return fmt.Errorf("required environment variable not set: MW_MCP_SESSION_DIR")
	}

	errFile := filepath.Join(sessionDir, "mcp_startup_error.txt")
	content := "Simulated MATLAB startup failure: license checkout failed\n\nError using connector.ensureServiceOn\nNo licenses available for MATLAB."

	if err := r.FS.WriteFile(errFile, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to write startup error file: %w", err)
	}

	return nil
}
