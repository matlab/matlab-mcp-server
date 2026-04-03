// Copyright 2026 The MathWorks, Inc.

package mpm

import (
	"os"
	"runtime"
)

// OSLayer provides OS-level operations for the install MATLAB usecase.
type OSLayer struct{}

func NewOSLayer() *OSLayer {
	return &OSLayer{}
}

func (o *OSLayer) GOOS() string {
	return runtime.GOOS
}

func (o *OSLayer) GOARCH() string {
	return runtime.GOARCH
}

func (o *OSLayer) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (o *OSLayer) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o *OSLayer) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (o *OSLayer) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}
