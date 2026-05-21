// Copyright 2026 The MathWorks, Inc.

package resourcelimit

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	unixfacade "github.com/matlab/matlab-mcp-core-server/internal/facades/unix"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type SyscallLayer interface {
	Getrlimit(resource int, rlim *unixfacade.Rlimit) error
	Setrlimit(resource int, rlim *unixfacade.Rlimit) error
}

type Manager struct {
	loggerFactory LoggerFactory
	syscallLayer  SyscallLayer
}

func New(
	loggerFactory LoggerFactory,
	syscallLayer SyscallLayer,
) *Manager {
	return &Manager{
		loggerFactory: loggerFactory,
		syscallLayer:  syscallLayer,
	}
}
