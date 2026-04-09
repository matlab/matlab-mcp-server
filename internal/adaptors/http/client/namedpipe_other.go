// Copyright 2025-2026 The MathWorks, Inc.

//go:build !windows

package client

import (
	"errors"
	"net"
)

// dialNamedPipe is not supported on non-Windows platforms
func dialNamedPipe(pipePath string) (net.Conn, error) {
	return nil, errors.New("named pipes are only supported on Windows")
}
