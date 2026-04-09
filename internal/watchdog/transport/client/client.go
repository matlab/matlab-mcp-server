// Copyright 2025-2026 The MathWorks, Inc.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	httpclient "github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/client"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
)

const (
	defaultSocketWaitTimeout   = 10 * time.Second
	defaultSocketRetryInterval = 500 * time.Millisecond
)

var (
	ErrSocketFileInaccessible      = errors.New("access denied for socket file")
	ErrTimeoutWaitingForSocketFile = errors.New("socket file access timed out")
	ErrClientNotConnected          = errors.New("client not connected")
	ErrHTTP                        = errors.New("HTTP request failed")
)

type Client struct {
	httpClient          httpclient.HttpClient
	osLayer             OSLayer
	httpClientFactory   HTTPClientFactory
	loggerFactory       LoggerFactory
	logger              entities.Logger
	socketWaitTimeout   time.Duration
	socketRetryInterval time.Duration
}

func newClient(
	osLayer OSLayer,
	httpClientFactory HTTPClientFactory,
	loggerFactory LoggerFactory,
) *Client {
	return &Client{
		osLayer:             osLayer,
		httpClientFactory:   httpClientFactory,
		loggerFactory:       loggerFactory,
		socketWaitTimeout:   defaultSocketWaitTimeout,
		socketRetryInterval: defaultSocketRetryInterval,
	}
}

func (c *Client) Connect(socketPath string) error {
	logger, messagesErr := c.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return messagesErr
	}
	c.logger = logger

	c.logger.
		With("socketPath", socketPath).
		Debug("Connecting to socket")

	// Check if this is a Windows named pipe path
	if isNamedPipePath(socketPath) {
		c.logger.Debug("Using Windows named pipe, skipping file check")
		c.httpClient = c.httpClientFactory.NewClientOverUDS(socketPath)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.socketWaitTimeout)
	defer cancel()

	httpClient, err := retry.Retry(ctx, func() (httpclient.HttpClient, bool, error) {
		_, err := c.osLayer.Stat(socketPath)

		if err == nil {
			c.logger.Debug("Socket file found")
			httpClient := c.httpClientFactory.NewClientOverUDS(socketPath)
			return httpClient, true, nil
		}

		if !os.IsNotExist(err) {
			c.logger.Debug("Socket file inaccessible")
			return nil, false, ErrSocketFileInaccessible
		}

		return nil, false, nil
	}, retry.NewLinearRetryStrategy(c.socketRetryInterval))

	if err != nil {
		// On Windows, if socket file is not found, try named pipe as fallback
		if errors.Is(err, context.DeadlineExceeded) && runtime.GOOS == "windows" {
			pipePath := socketPathToNamedPipe(socketPath)
			c.logger.With("pipePath", pipePath).Debug("Socket file not found, trying named pipe fallback")
			c.httpClient = c.httpClientFactory.NewClientOverUDS(pipePath)
			return nil
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimeoutWaitingForSocketFile
		}

		return err
	}

	c.httpClient = httpClient
	return nil
}

// isNamedPipePath checks if the path is a Windows named pipe path
func isNamedPipePath(path string) bool {
	return len(path) >= 9 && path[:9] == `\\.\pipe\`
}

func (c *Client) SendProcessPID(pid int) (messages.ProcessToKillResponse, error) {
	return post[messages.ProcessToKillRequest, messages.ProcessToKillResponse](
		c.httpClient,
		c.logger,
		messages.ProcessToKillPath,
		messages.ProcessToKillRequest{
			PID: pid,
		},
	)
}

func (c *Client) SendStop() (messages.ShutdownResponse, error) {
	return post[messages.ShutdownRequest, messages.ShutdownResponse](
		c.httpClient,
		c.logger,
		messages.ShutdownPath,
		messages.ShutdownRequest{},
	)
}

func post[RequestType any, ResponseType any](httpClient httpclient.HttpClient, logger entities.Logger, path string, reqBody RequestType) (ResponseType, error) {
	var zeroValueResp ResponseType

	if httpClient == nil {
		return zeroValueResp, ErrClientNotConnected
	}

	logger.
		With("path", path).
		Debug("Sending request")

	body, err := json.Marshal(reqBody)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal request")
		return zeroValueResp, ErrHTTP
	}

	// Host is ignored for UDS, but required for valid HTTP
	req, err := http.NewRequest("POST", "http://watchdog"+path, bytes.NewReader(body))
	if err != nil {
		logger.WithError(err).Error("Failed to create request")
		return zeroValueResp, ErrHTTP
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.WithError(err).Error("Failed to send request")
		return zeroValueResp, ErrHTTP
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithError(err).Error("Failed to read response body")
		return zeroValueResp, ErrHTTP
	}

	if err := resp.Body.Close(); err != nil {
		logger.WithError(err).Error("Failed to close response body")
		return zeroValueResp, ErrHTTP
	}

	if resp.StatusCode != http.StatusOK {
		logger.With("status", resp.Status).Error("Unexpected status")
		return zeroValueResp, ErrHTTP
	}

	var response ResponseType
	if err := json.Unmarshal(respBody, &response); err != nil {
		logger.WithError(err).Error("Failed to decode response")
		return zeroValueResp, ErrHTTP
	}

	logger.With("path", path).Debug("Request completed")

	return response, nil
}

func (c *Client) Close() error {
	c.logger.Info("Client closing")
	c.httpClient.CloseIdleConnections()
	return nil
}

// socketPathToNamedPipe converts a Unix socket path to a named pipe path
func socketPathToNamedPipe(socketPath string) string {
	// Extract the ID from the socket path (e.g., "watchdog-123.sock" -> "123")
	base := filepath.Base(socketPath)
	
	// Remove .sock extension if present
	base = strings.TrimSuffix(base, ".sock")
	
	return `\\.\pipe\matlab-mcp-` + base
}
