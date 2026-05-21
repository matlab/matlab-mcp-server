// Copyright 2026 The MathWorks, Inc.
//go:build !windows

package resourcelimit

import (
	"fmt"

	unixfacade "github.com/matlab/matlab-mcp-core-server/internal/facades/unix"
	"golang.org/x/sys/unix"
)

func (m *Manager) CapOpenFilesLimit(limit uint64) (func() error, error) {
	// CapOpenFilesLimit lowers the RLIMIT_NOFILE soft limit to limit if it is
	// currently above that value (or unlimited). It returns a function that
	// restores the original limit. If the current limit is already at or below
	// the requested value, the limit is left untouched and the returned function
	// is a no-op.
	logger, loggerErr := m.loggerFactory.GetGlobalLogger()
	if loggerErr != nil {
		return nil, loggerErr
	}

	var original unixfacade.Rlimit
	if err := m.syscallLayer.Getrlimit(unix.RLIMIT_NOFILE, &original); err != nil {
		logger.WithError(err).Error("Failed to get RLIMIT_NOFILE")
		return nil, fmt.Errorf("failed to get RLIMIT_NOFILE: %w", err)
	}

	if original.Cur <= limit {
		// no-op
		return func() error { return nil }, nil
	}

	limited := original
	limited.Cur = limit
	if err := m.syscallLayer.Setrlimit(unix.RLIMIT_NOFILE, &limited); err != nil {
		logger.WithError(err).Error(fmt.Sprintf("Failed to set RLIMIT_NOFILE to %d", limit))
		return nil, fmt.Errorf("failed to set RLIMIT_NOFILE to %d: %w", limit, err)
	}
	logger.Debug(fmt.Sprintf("Lowered RLIMIT_NOFILE soft limit from %d to %d", original.Cur, limit))

	return func() error {
		if err := m.syscallLayer.Setrlimit(unix.RLIMIT_NOFILE, &original); err != nil {
			return fmt.Errorf("failed to restore RLIMIT_NOFILE: %w", err)
		}
		return nil
	}, nil
}
