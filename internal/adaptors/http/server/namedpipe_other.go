// Copyright 2025-2026 The MathWorks, Inc.

//go:build !windows

package server

import (
	"errors"
	"net"
)

// listenNamedPipe is not supported on non-Windows platforms
func listenNamedPipe(pipePath string) (net.Listener, error) {
	return nil, errors.New("named pipes are only supported on Windows")
}
