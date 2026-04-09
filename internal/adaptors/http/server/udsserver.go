// Copyright 2025-2026 The MathWorks, Inc.

package server

import (
	"context"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
)

const defaultReadHeaderTimeout = 10 * time.Second

type udsServer struct {
	httpServer *http.Server
	osLayer    OSLayer
	socketPath string

	lock *sync.Mutex
}

func newUDSServer(
	osLayer OSLayer,
	handlers map[string]http.HandlerFunc,
) *udsServer {
	mux := http.NewServeMux()
	for pattern, handler := range handlers {
		mux.HandleFunc(pattern, handler)
	}

	return &udsServer{
		httpServer: &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
		},
		osLayer: osLayer,
		lock:    new(sync.Mutex),
	}
}

func (s *udsServer) Serve(socketPath string) error {
	s.setSocketPath(socketPath)

	// Check if this is a Windows Named Pipe path
	if isNamedPipePath(socketPath) {
		return s.serveNamedPipe(socketPath)
	}

	// Try Unix Domain Socket first
	if err := s.serveUnixSocket(socketPath); err != nil {
		// On Windows, if Unix socket fails, fall back to Named Pipe
		if runtime.GOOS == "windows" {
			// Convert socket path to named pipe path
			pipePath := socketPathToNamedPipe(socketPath)
			return s.serveNamedPipe(pipePath)
		}
		return err
	}
	return nil
}

// isNamedPipePath checks if the path is a Windows named pipe path
func isNamedPipePath(path string) bool {
	return strings.HasPrefix(path, `\\.\pipe\`)
}

// socketPathToNamedPipe converts a Unix socket path to a named pipe path
func socketPathToNamedPipe(socketPath string) string {
	// Extract the ID from the socket path (e.g., "watchdog-123.sock" -> "123")
	base := socketPath
	if idx := strings.LastIndex(socketPath, "\\"); idx != -1 {
		base = socketPath[idx+1:]
	} else if idx := strings.LastIndex(socketPath, "/"); idx != -1 {
		base = socketPath[idx+1:]
	}
	
	// Remove .sock extension if present
	base = strings.TrimSuffix(base, ".sock")
	
	return `\\.\pipe\matlab-mcp-` + base
}

// serveUnixSocket tries to serve over Unix Domain Socket
func (s *udsServer) serveUnixSocket(socketPath string) error {
	// Socket path max length is 108 characters, but for safety using 105
	if len(socketPath) > 105 {
		return socket.ErrSocketPathTooLong
	}

	if err := s.osLayer.RemoveAll(socketPath); err != nil {
		// Log but don't fail - file might not exist
		// This is expected on Windows where Unix sockets are not supported
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// serveNamedPipe serves over Windows Named Pipe
func (s *udsServer) serveNamedPipe(pipePath string) error {
	listener, err := listenNamedPipe(pipePath)
	if err != nil {
		return err
	}

	if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *udsServer) Shutdown(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	if s.socketPath == "" {
		return nil
	}

	if err := s.osLayer.RemoveAll(s.socketPath); err != nil {
		return err
	}

	return nil
}

func (s *udsServer) setSocketPath(socketPath string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.socketPath = socketPath
}
