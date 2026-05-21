// Copyright 2026 The MathWorks, Inc.
//go:build windows

package resourcelimit

func (m *Manager) CapOpenFilesLimit(_ uint64) (func() error, error) {
	// no-op
	return func() error { return nil }, nil
}
