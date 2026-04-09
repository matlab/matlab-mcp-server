// Copyright 2025-2026 The MathWorks, Inc.

//go:build windows

package client

import (
	"net"

	"github.com/Microsoft/go-winio"
)

// dialNamedPipe connects to a Windows named pipe
func dialNamedPipe(pipePath string) (net.Conn, error) {
	return winio.DialPipe(pipePath, nil)
}
