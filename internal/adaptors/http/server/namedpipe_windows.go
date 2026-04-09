// Copyright 2025-2026 The MathWorks, Inc.

//go:build windows

package server

import (
	"net"

	"github.com/Microsoft/go-winio"
)

// listenNamedPipe creates a named pipe listener on Windows
func listenNamedPipe(pipePath string) (net.Listener, error) {
	return winio.ListenPipe(pipePath, nil)
}
