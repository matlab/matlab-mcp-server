// Copyright 2026 The MathWorks, Inc.

package mockmatlab

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/sessiondetails"
)

const (
	securePortFile     = "connector.securePort"
	certificateFile    = "cert.pem"
	certificateKeyFile = "cert.key"
	eventLogSubDir     = "events"

	EnvMockMATLABConfig = mockruntime.EnvMockMATLABConfig
	EnvMockMATLABLogDir = mockruntime.EnvMockMATLABLogDir

	defaultReadyTimeout = 10 * time.Second
	defaultReadyPoll    = 100 * time.Millisecond
)

type Config = mockruntime.Config

func HappyConfig() Config {
	return mockruntime.HappyConfig()
}

func HangBeforeFilesConfig() Config {
	return mockruntime.HangBeforeFilesConfig()
}

func ExitImmediatelyConfig(exitCode int) Config {
	return mockruntime.ExitImmediatelyConfig(exitCode)
}

func SlowStartupConfig(delayMs int) Config {
	return mockruntime.SlowStartupConfig(delayMs)
}

func StartupFailureConfig() Config {
	return mockruntime.StartupFailureConfig()
}

type Session struct {
	cmd               *exec.Cmd
	SessionDir        string
	APIKey            string
	eventLogDir       string
	connectionDetails embeddedconnector.ConnectionDetails
}

func StartSession(ctx context.Context, installation *Installation, cfg Config) (*Session, error) {
	sessionDir, err := os.MkdirTemp("", "mock-matlab-session-")
	if err != nil {
		return nil, fmt.Errorf("failed to create session dir: %w", err)
	}

	apiKey := "mock-api-key" //nolint:gosec // Not a real credential
	certFile := filepath.Join(sessionDir, certificateFile)
	keyFile := filepath.Join(sessionDir, certificateKeyFile)
	eventLogDir := filepath.Join(sessionDir, eventLogSubDir)

	configJSON, err := cfg.ToEnvValue()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize config: %w", err)
	}

	binaryPath := mockMATLABBinaryPath(installation.MATLABRoot)
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec // Trusted test path
	cmd.Env = append(os.Environ(),
		"MW_MCP_SESSION_DIR="+sessionDir,
		"MWAPIKEY="+apiKey,
		"MW_CERTFILE="+certFile,
		"MW_PKEYFILE="+keyFile,
		mockruntime.EnvMockMATLABConfig+"="+configJSON,
		mockruntime.EnvMockMATLABLogDir+"="+eventLogDir,
	)

	if err := cmd.Start(); err != nil {
		_ = os.RemoveAll(sessionDir)
		return nil, fmt.Errorf("failed to start mock MATLAB: %w", err)
	}

	return &Session{
		cmd:         cmd,
		SessionDir:  sessionDir,
		APIKey:      apiKey,
		eventLogDir: eventLogDir,
	}, nil
}

func (s *Session) WaitForReady(ctx context.Context) (embeddedconnector.ConnectionDetails, error) {
	portPath := filepath.Join(s.SessionDir, securePortFile)
	certPath := filepath.Join(s.SessionDir, certificateFile)

	ctx, cancel := context.WithTimeout(ctx, defaultReadyTimeout)
	defer cancel()

	details, err := retry.Retry(ctx, func() (embeddedconnector.ConnectionDetails, bool, error) {
		port, readErr := readNonEmptyFile(portPath)
		if readErr != nil {
			return embeddedconnector.ConnectionDetails{}, false, nil
		}
		certPEM, readErr := readNonEmptyFile(certPath)
		if readErr != nil {
			return embeddedconnector.ConnectionDetails{}, false, nil
		}
		return embeddedconnector.ConnectionDetails{
			Host:           "localhost",
			Port:           string(port),
			APIKey:         s.APIKey,
			CertificatePEM: certPEM,
		}, true, nil
	}, retry.NewLinearRetryStrategy(defaultReadyPoll))

	if err != nil {
		return embeddedconnector.ConnectionDetails{}, fmt.Errorf("timeout waiting for mock MATLAB to become ready")
	}

	s.connectionDetails = details
	return details, nil
}

// ToSessionDetailsJSON returns a JSON string in the format expected by
// --matlab-session-connection-details. The certificate field is the file path
// to the PEM cert in the session directory.
func (s *Session) ToSessionDetailsJSON() (string, error) {
	return sessiondetails.MarshalJSON(
		s.connectionDetails.Port,
		s.CertificatePath(),
		s.APIKey,
		s.cmd.Process.Pid,
	)
}

// CertificatePath returns the path to the PEM certificate file in the session directory.
func (s *Session) CertificatePath() string {
	return filepath.Join(s.SessionDir, certificateFile)
}

// ShareMATLABSession publishes session details to the standard discovery
// location under homeDir, mimicking what a real MATLAB session would do.
func (s *Session) ShareMATLABSession(homeDir string) (string, error) {
	detailsJSON, err := s.ToSessionDetailsJSON()
	if err != nil {
		return "", err
	}
	return sessiondetails.Publish(homeDir, detailsJSON)
}

// ReceivedEvents reads all events recorded by this mock MATLAB session.
func (s *Session) ReceivedEvents() ([]mockruntime.InstanceEvents, error) {
	return mockruntime.ReadEventsForAllInstances(s.eventLogDir)
}

// ReceivedEvals returns the eval/feval events received by this mock session.
func (s *Session) ReceivedEvals() ([]mockruntime.Event, error) {
	instances, err := s.ReceivedEvents()
	if err != nil {
		return nil, err
	}
	var evals []mockruntime.Event
	for _, inst := range instances {
		for _, e := range inst.Events {
			if e.Type == mockruntime.EventEval || e.Type == mockruntime.EventFeval {
				evals = append(evals, e)
			}
		}
	}
	return evals, nil
}

func (s *Session) Stop() error {
	var firstErr error
	recordErr := func(err error) {
		if firstErr == nil {
			firstErr = err
		}
	}

	if s.cmd != nil && s.cmd.Process != nil {
		if err := s.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			recordErr(fmt.Errorf("failed to kill mock MATLAB process: %w", err))
		}

		if err := s.cmd.Wait(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) {
				recordErr(fmt.Errorf("failed to wait for mock MATLAB process: %w", err))
			}
		}
	}

	if err := os.RemoveAll(s.SessionDir); err != nil {
		recordErr(fmt.Errorf("failed to remove mock MATLAB session directory: %w", err))
	}

	return firstErr
}

func (s *Session) Wait() error {
	return s.cmd.Wait()
}

func (s *Session) ProcessExited() bool {
	return s.cmd.ProcessState != nil
}

func readNonEmptyFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path) //nolint:gosec // Trusted test path
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("file is empty")
	}
	return data, nil
}
