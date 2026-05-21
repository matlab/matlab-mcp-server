// Copyright 2025-2026 The MathWorks, Inc.

package os

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
)

const defaultCheckParentAliveInterval = 1 * time.Second

type OSLayer interface {
	FindProcess(pid int) (osfacade.Process, error)
	GOOS() string
}

type ProcessManager struct {
	osLayer OSLayer

	goos                     string
	checkParentAliveInterval time.Duration
}

func NewProcessManager(osLayer OSLayer) *ProcessManager {
	return &ProcessManager{
		osLayer: osLayer,

		goos:                     osLayer.GOOS(),
		checkParentAliveInterval: defaultCheckParentAliveInterval,
	}
}

func (pm *ProcessManager) FindProcess(processPid int) osfacade.Process {
	proc, err := pm.osLayer.FindProcess(processPid)
	if err != nil {
		return nil
	}

	if pm.goos != "windows" {
		err = proc.Signal(syscall.Signal(0))
		if err != nil {
			return nil
		}
	}

	return proc
}

func (pm *ProcessManager) WaitForProcessToComplete(processPid int) {
	if pm.goos == "windows" {
		if process := pm.FindProcess(processPid); process != nil {
			process.Wait() //nolint:gosec,errcheck // It doesn't matter why we stopped waiting
		}
	} else {
		retry.Retry(context.Background(), func() (struct{}, bool, error) { //nolint:errcheck,gosec // Runs indefinitely until process exits
			return struct{}{}, pm.FindProcess(processPid) == nil, nil
		}, retry.NewLinearRetryStrategy(pm.checkParentAliveInterval))
	}
}

func (pm *ProcessManager) InterruptSignalChan() <-chan os.Signal {
	interruptC := make(chan os.Signal, 1)
	signal.Notify(interruptC, os.Interrupt, syscall.SIGTERM)
	return interruptC
}
