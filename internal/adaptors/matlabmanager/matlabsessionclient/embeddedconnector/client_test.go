// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/stretchr/testify/require"
)

func buildEvalWithCaptureResponse(t *testing.T, entries []embeddedconnector.LiveEditorResponseEntry) []byte {
	data, err := json.Marshal(entries)
	require.NoError(t, err, "Failed to marshal entries to JSON")

	return buildEvalWithCaptureResponseFromJSON(t, string(data))
}

func buildEvalWithCaptureResponseFromJSON(t *testing.T, entries string) []byte {
	responsePayload := embeddedconnector.ConnectorPayload{
		Messages: embeddedconnector.ConnectorMessage{
			FevalResponse: []embeddedconnector.FevalResponseMessage{
				{
					IsError: false,
					Results: []interface{}{entries},
				},
			},
		},
	}

	responseBody, err := json.Marshal(responsePayload)
	require.NoError(t, err, "Failed to marshal response payload to JSON")

	return responseBody
}

func parseConnectorRequest(req *http.Request) (embeddedconnector.ConnectorPayload, bool) {
	if req.Method != "POST" {
		return embeddedconnector.ConnectorPayload{}, false
	}

	if req.Header.Get("Content-Type") != "application/json" {
		return embeddedconnector.ConnectorPayload{}, false
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return embeddedconnector.ConnectorPayload{}, false
	}

	// Restore the body, so that the request can be read again
	req.Body = io.NopCloser(bytes.NewReader(body))

	var reqBody embeddedconnector.ConnectorPayload
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return embeddedconnector.ConnectorPayload{}, false
	}

	return reqBody, true
}

func validateConnectorRequest(req *http.Request) bool {
	_, ok := parseConnectorRequest(req)
	return ok
}
