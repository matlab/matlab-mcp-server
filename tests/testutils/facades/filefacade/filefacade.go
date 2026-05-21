// Copyright 2025-2026 The MathWorks, Inc.

package filefacade

import (
	"io/fs"
	"os"
	"path/filepath"
)

// RealFileSystem implements FileSystem using the os and filepath packages
type RealFileSystem struct{}

func (RealFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (RealFileSystem) MkdirTemp(dir string, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

func (RealFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (RealFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (RealFileSystem) ReadFile(fileSystem fs.FS, name string) ([]byte, error) {
	return fs.ReadFile(fileSystem, name)
}

func (RealFileSystem) Glob(fileSystem fs.FS, pattern string) ([]string, error) {
	return fs.Glob(fileSystem, pattern)
}

func (RealFileSystem) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

func (RealFileSystem) DirFS(path string) fs.FS {
	return os.DirFS(path)
}
