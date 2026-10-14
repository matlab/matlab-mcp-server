// Copyright 2026 The MathWorks, Inc.

package matlabtempdir

import (
	"errors"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	removeTimeout  = 30 * time.Second
	removeInterval = 250 * time.Millisecond
)

func New(t *testing.T, pattern string) string {
	t.Helper()

	dir, err := os.MkdirTemp("", pattern) //nolint:usetesting // t.TempDir fails the test when an exiting MATLAB still holds the directory
	require.NoError(t, err, "should create temp dir matching %s", pattern)
	t.Cleanup(func() { removeWithRetry(t, dir) })

	return dir
}

func removeWithRetry(t *testing.T, dir string) {
	t.Helper()

	deadline := time.Now().Add(removeTimeout)
	for {
		err := os.RemoveAll(dir)
		if err == nil {
			return
		}
		if errors.Is(err, fs.ErrPermission) {
			makeOwnerWritable(dir)
		}
		if time.Now().After(deadline) {
			t.Logf("Could not remove %s within %s, leaving it on disk: %v", dir, removeTimeout, err)
			return
		}
		time.Sleep(removeInterval)
	}
}

func makeOwnerWritable(dir string) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return
	}
	defer func() { _ = root.Close() }()

	_ = fs.WalkDir(root.FS(), ".", func(path string, entry fs.DirEntry, err error) error {
		switch {
		case err != nil, !entry.IsDir() && !entry.Type().IsRegular():
			return nil
		case entry.IsDir():
			_ = root.Chmod(path, 0o700)
		default:
			_ = root.Chmod(path, 0o600)
		}
		return nil
	})
}
