// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/facades/filefacade"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
)

const (
	envSessionDir = "MW_MCP_SESSION_DIR"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	runtime := mockruntime.NewRuntime(
		osfacade.RealEnvironment{},
		filefacade.RealFileSystem{},
		mockruntime.NewDefaultTLSMaterialProvider(),
	)

	cfg, err := runtime.LoadConfigFromEnv()
	if err != nil {
		log.Printf("mock MATLAB configuration error: %v", err)
		os.Exit(1)
	}

	recorder, err := initEventRecorder()
	if err != nil {
		log.Printf("mock MATLAB event recorder error: %v", err)
		os.Exit(1)
	}
	defer recorder.Close()
	recorder.RecordStarted(cfg.Mode)

	switch cfg.Mode {
	case mockruntime.ModeHangBeforeFiles:
		<-ctx.Done()
		return
	case mockruntime.ModeExitImmediately:
		if cfg.ExitCode == nil {
			os.Exit(1)
		}
		os.Exit(*cfg.ExitCode)
	case mockruntime.ModeSlowStartup:
		delayMs := 0
		if cfg.DelayMs != nil {
			delayMs = *cfg.DelayMs
		}

		select {
		case <-time.After(time.Duration(delayMs) * time.Millisecond):
		case <-ctx.Done():
			return
		}

		if err := runHappyPath(ctx, runtime, recorder); err != nil {
			log.Printf("mock MATLAB startup failed: %v", err)
			os.Exit(1)
		}
	case mockruntime.ModeStartupFailure:
		if err := runtime.WriteStartupFailureFile(os.Getenv(envSessionDir)); err != nil {
			log.Printf("mock MATLAB startup failure flow failed: %v", err)
		}
		os.Exit(1)
	default:
		if err := runHappyPath(ctx, runtime, recorder); err != nil {
			log.Printf("mock MATLAB startup failed: %v", err)
			os.Exit(1)
		}
	}
}

func initEventRecorder() (*mockruntime.EventRecorder, error) {
	logDir := os.Getenv(mockruntime.EnvMockMATLABLogDir)
	if logDir == "" {
		return nil, fmt.Errorf("%s environment variable is not set", mockruntime.EnvMockMATLABLogDir)
	}
	return mockruntime.NewEventRecorder(logDir)
}

func runHappyPath(ctx context.Context, runtime *mockruntime.Runtime, recorder *mockruntime.EventRecorder) error {
	sessionDir := os.Getenv(envSessionDir)
	apiKey := os.Getenv("MWAPIKEY")
	certFile := os.Getenv("MW_CERTFILE")
	keyFile := os.Getenv("MW_PKEYFILE")

	if sessionDir == "" || apiKey == "" || certFile == "" || keyFile == "" {
		return fmt.Errorf("required environment variables not set: %s, MWAPIKEY, MW_CERTFILE, MW_PKEYFILE", envSessionDir)
	}

	tlsCfg, err := runtime.GenerateAndWriteCert(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to generate TLS cert: %w", err)
	}

	if err := startConnectorServer(ctx, sessionDir, apiKey, tlsCfg, recorder); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
