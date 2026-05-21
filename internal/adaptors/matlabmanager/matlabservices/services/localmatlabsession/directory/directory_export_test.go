// Copyright 2025-2026 The MathWorks, Inc.

package directory

import (
	"path/filepath"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

func NewDirectory(logger entities.Logger, sessionDir string, osLayer OSLayer, config Config) *directory {
	return newDirectory(logger, sessionDir, osLayer, config)
}

func (d *directory) SecurePortFile() string {
	return d.securePortFile()
}

func (d *directory) StartupErrorFile() string {
	return filepath.Join(d.sessionDir, startupErrorFile)
}

func (d *directory) SetEmbeddedConnectorDetailsRetry(retry time.Duration) {
	d.embeddedConnectorDetailsRetry = retry
}

func (d *directory) GetEmbeddedConnectorDetailsTimeout() time.Duration {
	return d.embeddedConnectorDetailsTimeout
}

func (d *directory) SetCleanupTimeout(timeout time.Duration) {
	d.cleanupTimeout = timeout
}

func (d *directory) SetCleanupRetry(retry time.Duration) {
	d.cleanupRetry = retry
}
