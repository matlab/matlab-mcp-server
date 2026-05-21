// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	httpclient "github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/client"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

const defaultPingRetry = 100 * time.Millisecond
const defaultPingTimeout = 1 * time.Second

type HttpClientFactory interface {
	NewClientForSelfSignedTLSServer(certificatePEM []byte) (httpclient.HttpClient, error)
}

type ConnectionDetails struct {
	Host           string
	Port           string
	APIKey         string
	CertificatePEM []byte
}

type Client struct {
	host       string
	port       string
	apiKey     string
	httpClient httpclient.HttpClient

	pingRetry   time.Duration
	pingTimeout time.Duration
}

func NewClient(
	endpoint ConnectionDetails,
	httpClientFactory HttpClientFactory,
) (*Client, error) {
	httpClient, err := httpClientFactory.NewClientForSelfSignedTLSServer(endpoint.CertificatePEM)
	if err != nil {
		return nil, err
	}

	return &Client{
		host:       endpoint.Host,
		port:       endpoint.Port,
		apiKey:     endpoint.APIKey,
		httpClient: httpClient,

		pingRetry:   defaultPingRetry,
		pingTimeout: defaultPingTimeout,
	}, nil
}

func (c *Client) SetPingTimeout(timeout time.Duration) {
	c.pingTimeout = timeout
}

func (c *Client) SetPingRetry(retry time.Duration) {
	c.pingRetry = retry
}

func (c *Client) Eval(ctx context.Context, logger entities.Logger, input entities.EvalRequest) (entities.EvalResponse, error) {
	payload := ConnectorPayload{
		Messages: ConnectorMessage{
			Eval: []EvalMessage{
				{
					Code: input.Code,
				},
			},
		},
	}

	response, err := c.sendRequestToEvaluationEndpoint(ctx, logger, payload)
	if err != nil {
		return entities.EvalResponse{}, err
	}

	if len(response.Messages.EvalResponse) == 0 {
		logger.
			Debug("No EvalResponse messages received")
		return entities.EvalResponse{}, fmt.Errorf("no response messages received")
	}

	if response.Messages.EvalResponse[0].IsError {
		return entities.EvalResponse{}, newMATLABError(response.Messages.EvalResponse[0].ResponseStr)
	}

	return entities.EvalResponse{
		ConsoleOutput: response.Messages.EvalResponse[0].ResponseStr,
		Images:        nil,
	}, nil
}

func (c *Client) EvalWithCapture(ctx context.Context, logger entities.Logger, input entities.EvalRequest) (entities.EvalResponse, error) {
	fevalRequest := entities.FEvalRequest{
		Function:   "matlab_mcp.mcpEval",
		Arguments:  []string{input.Code},
		NumOutputs: 1,
	}

	response, err := c.FEval(ctx, logger, fevalRequest)
	if err != nil {
		return entities.EvalResponse{}, err
	}

	outputs, err := parseEvalWithCaptureResponse(response)
	if err != nil {
		return entities.EvalResponse{}, err
	}

	return outputs, nil
}

func (c *Client) FEval(ctx context.Context, logger entities.Logger, input entities.FEvalRequest) (entities.FEvalResponse, error) {
	payload := ConnectorPayload{
		Messages: ConnectorMessage{
			FEval: []FevalMessage{
				{
					Function:  input.Function,
					Arguments: input.Arguments,
					Nargout:   input.NumOutputs,
				},
			},
		},
	}

	response, err := c.sendRequestToEvaluationEndpoint(ctx, logger, payload)
	if err != nil {
		return entities.FEvalResponse{}, err
	}

	if len(response.Messages.FevalResponse) == 0 {
		logger.
			Debug("No FEvalResponse messages received")
		return entities.FEvalResponse{}, fmt.Errorf("no response messages received")
	}

	if response.Messages.FevalResponse[0].IsError {
		if len(response.Messages.FevalResponse[0].MessageFaults) == 0 {
			logger.
				Debug("Response was in error state but no fault messages received")
			return entities.FEvalResponse{}, fmt.Errorf("response was in error state but no fault messages received")
		}

		var errorMessage string
		for _, rawFault := range response.Messages.FevalResponse[0].MessageFaults {
			var f Fault
			if err := json.Unmarshal(rawFault, &f); err != nil {
				logger.
					WithError(err).
					Debug("Failed to deserialize fault message into a fault")
			}
			errorMessage += f.Message + "\n\n"
		}
		return entities.FEvalResponse{}, newMATLABError(errorMessage)
	}

	return entities.FEvalResponse{
		Outputs: response.Messages.FevalResponse[0].Results,
	}, nil
}

func (m *Client) Ping(ctx context.Context, sessionLogger entities.Logger) entities.PingResponse {
	pingCtx, cancel := context.WithTimeout(ctx, m.pingTimeout)
	defer cancel()

	_, err := retry.Retry(pingCtx, func() (struct{}, bool, error) {
		status, err := m.pingMATLAB(pingCtx, sessionLogger)
		if err != nil {
			sessionLogger.WithError(err).Debug("Ping to MATLAB session failed")
		}
		return struct{}{}, status, nil
	}, retry.NewLinearRetryStrategy(m.pingRetry))

	if err != nil {
		sessionLogger.Warn("timeout waiting for matlab to be ready")
		return entities.PingResponse{
			IsAlive: false,
		}
	}

	return entities.PingResponse{
		IsAlive: true,
	}
}

func (c *Client) pingMATLAB(ctx context.Context, logger entities.Logger) (bool, error) {
	payload := ConnectorPayload{
		Messages: ConnectorMessage{
			Ping: []PingMessage{{}},
		},
	}

	response, err := c.sendRequestToStateEndpoint(ctx, logger, payload)
	if err != nil {
		return false, err
	}

	if len(response.Messages.PingResponse) == 0 {
		return false, fmt.Errorf("no response messages received")
	}

	messageFaults := response.Messages.PingResponse[0].MessageFaults
	if len(messageFaults) > 0 {
		var errorMessage strings.Builder
		for _, rawFault := range response.Messages.PingResponse[0].MessageFaults {
			var f Fault
			if err := json.Unmarshal(rawFault, &f); err != nil {
				logger.
					WithError(err).
					Debug("Failed to deserialize fault message into a fault")
			}
			errorMessage.WriteString(f.Message)
			errorMessage.WriteString("\n\n")
		}
		return false, newMATLABError(errorMessage.String())
	}

	return true, nil
}

func (c *Client) sendRequestToEvaluationEndpoint(ctx context.Context, logger entities.Logger, payload ConnectorPayload) (ConnectorPayload, error) {
	endpoint := fmt.Sprintf("https://%s:%s/messageservice/json/secure", c.host, c.port)
	return c.sendRequest(ctx, logger, endpoint, payload)
}

func (c *Client) sendRequestToStateEndpoint(ctx context.Context, logger entities.Logger, payload ConnectorPayload) (ConnectorPayload, error) {
	endpoint := fmt.Sprintf("https://%s:%s/messageservice/json/state", c.host, c.port)
	return c.sendRequest(ctx, logger, endpoint, payload)
}

func (c *Client) sendRequest(ctx context.Context, logger entities.Logger, endpoint string, payload ConnectorPayload) (ConnectorPayload, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return ConnectorPayload{}, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.
			Debug("Failed to create HTTP request")
		return ConnectorPayload{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("mwapikey", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.
			Debug("Failed to send HTTP request")
		return ConnectorPayload{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.
				WithError(err).
				Debug("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logger.
			With("status", resp.Status).
			With("status-code", resp.StatusCode).
			Debug("Request failed")
		return ConnectorPayload{}, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.
			Debug("Failed to read response body")
		return ConnectorPayload{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var response ConnectorPayload
	if err := json.Unmarshal(body, &response); err != nil {
		logger.
			Debug("Failed to unmarshal response")
		return ConnectorPayload{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}
