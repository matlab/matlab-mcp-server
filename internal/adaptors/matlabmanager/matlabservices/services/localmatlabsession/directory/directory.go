// Copyright 2025-2026 The MathWorks, Inc.

package directory

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

const defaultEmbeddedConnectorDetailsRetry = 500 * time.Millisecond

const defaultCleanupTimeout = 2 * time.Minute
const defaultCleanupRetry = 500 * time.Millisecond

const securePortFile = "connector.securePort"
const certificateFile = "cert.pem"
const certificateKeyFile = "cert.key"
const startupErrorFile = "mcp_startup_error.txt"

var ErrMATLABStartup = errors.New("MATLAB startup failed")

type Config interface {
	EmbeddedConnectorDetailsTimeout() time.Duration
}

type directory struct {
	logger     entities.Logger
	sessionDir string
	osLayer    OSLayer

	embeddedConnectorDetailsTimeout time.Duration
	embeddedConnectorDetailsRetry   time.Duration
	cleanupTimeout                  time.Duration
	cleanupRetry                    time.Duration
}

func newDirectory(logger entities.Logger, sessionDir string, osLayer OSLayer, config Config) *directory {
	return &directory{
		logger:     logger,
		sessionDir: sessionDir,
		osLayer:    osLayer,

		embeddedConnectorDetailsTimeout: config.EmbeddedConnectorDetailsTimeout(),
		embeddedConnectorDetailsRetry:   defaultEmbeddedConnectorDetailsRetry,
		cleanupTimeout:                  defaultCleanupTimeout,
		cleanupRetry:                    defaultCleanupRetry,
	}
}

func (d *directory) Path() string {
	return d.sessionDir
}

func (d *directory) CertificateFile() string {
	return filepath.Join(d.sessionDir, certificateFile)
}

func (d *directory) CertificateKeyFile() string {
	return filepath.Join(d.sessionDir, certificateKeyFile)
}

func (d *directory) GetEmbeddedConnectorDetails() (string, []byte, error) {
	securePortFileFullPath := d.securePortFile()
	certificateFileFullPath := d.CertificateFile()
	type embeddedConnectorDetails struct {
		port        string
		certificate []byte
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.embeddedConnectorDetailsTimeout)
	defer cancel()

	d.logger.
		With("timeout", d.embeddedConnectorDetailsTimeout.String()).
		Debug("Watching for EC details")

	details, err := retry.Retry(ctx, func() (embeddedConnectorDetails, bool, error) {
		if err := d.checkStartupError(); err != nil {
			return embeddedConnectorDetails{}, false, err
		}
		if _, err := d.osLayer.Stat(securePortFileFullPath); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return embeddedConnectorDetails{}, false, fmt.Errorf("failed to stat secure port file: %w", err)
			}
			return embeddedConnectorDetails{}, false, nil
		}
		if _, err := d.osLayer.Stat(certificateFileFullPath); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return embeddedConnectorDetails{}, false, fmt.Errorf("failed to stat certificate file: %w", err)
			}
			return embeddedConnectorDetails{}, false, nil
		}
		securePort, err := d.osLayer.ReadFile(securePortFileFullPath)
		if err != nil {
			return embeddedConnectorDetails{}, false, fmt.Errorf("failed to read secure port file: %w", err)
		}
		if string(securePort) == "" {
			return embeddedConnectorDetails{}, false, nil
		}
		certificatePEM, err := d.osLayer.ReadFile(certificateFileFullPath)
		if err != nil {
			return embeddedConnectorDetails{}, false, fmt.Errorf("failed to read certificate path file: %w", err)
		}
		if string(certificatePEM) == "" {
			return embeddedConnectorDetails{}, false, nil
		}

		return embeddedConnectorDetails{
			port:        string(securePort),
			certificate: certificatePEM,
		}, true, nil
	}, retry.NewLinearRetryStrategy(d.embeddedConnectorDetailsRetry))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", nil, fmt.Errorf("timeout waiting for worker to start")
		}
		return "", nil, err
	}

	return details.port, details.certificate, nil
}

func (d *directory) Cleanup() error {
	if d.sessionDir == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.cleanupTimeout)
	defer cancel()

	_, err := retry.Retry(ctx, func() (struct{}, bool, error) {
		err := d.osLayer.RemoveAll(d.sessionDir)
		if err != nil {
			return struct{}{}, false, nil
		}

		return struct{}{}, true, nil
	}, retry.NewLinearRetryStrategy(d.cleanupRetry))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("timeout trying to delete session directory %s", d.sessionDir)
		}
		return err
	}

	return nil
}

func (d *directory) securePortFile() string {
	return filepath.Join(d.sessionDir, securePortFile)
}

func (d *directory) checkStartupError() error {
	path := filepath.Join(d.sessionDir, startupErrorFile)
	content, err := d.osLayer.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to read startup error file: %w", err)
	}
	if len(content) == 0 {
		return nil
	}
	return fmt.Errorf("%w: %s", ErrMATLABStartup, strings.TrimSpace(string(content)))
}
