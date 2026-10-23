// Copyright 2026 The MathWorks, Inc.

package embeddedconnector_test

import (
	"net/http"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	"github.com/matlab/matlab-mcp-server/tests/integration"
	"github.com/matlab/matlab-mcp-server/tests/testutils/mockembeddedconnector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Eval_HappyPath(t *testing.T) {
	// Arrange
	logger := testutils.NewInspectableLogger()

	const expectedCode = "disp('Hello World')"
	const expectedOutput = "mock echo: disp('Hello World')"

	server := mockembeddedconnector.New(t,
		func(response http.ResponseWriter, request *http.Request) {
			req := mockembeddedconnector.ReadConnectorRequest(t, request)
			assert.Len(t, req.Messages.Eval, 1)
			assert.Equal(t, expectedCode, req.Messages.Eval[0].Code)

			mockembeddedconnector.RespondWithJSON(t, response, embeddedconnector.ConnectorPayload{
				Messages: embeddedconnector.ConnectorMessage{
					EvalResponse: []embeddedconnector.EvalResponseMessage{
						{
							IsError:     false,
							ResponseStr: expectedOutput,
						},
					},
				},
			})
		},
		nil,
	)
	defer server.Stop()

	client := newClient(t, server.ConnectionDetails())

	// Act
	response, err := client.Eval(t.Context(), logger, entities.EvalRequest{
		Code: expectedCode,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_FEval_HappyPath(t *testing.T) {
	// Arrange
	logger := testutils.NewInspectableLogger()

	expectedFunction := "sum"
	expectedArguments := []string{"1", "2"}
	expectedNumOutputs := 1
	expectedResults := []any{expectedFunction, []any{"1", "2"}, float64(expectedNumOutputs)}

	server := mockembeddedconnector.New(t,
		func(response http.ResponseWriter, request *http.Request) {
			req := mockembeddedconnector.ReadConnectorRequest(t, request)
			assert.Len(t, req.Messages.FEval, 1)
			assert.Equal(t, expectedFunction, req.Messages.FEval[0].Function)
			assert.Equal(t, expectedArguments, req.Messages.FEval[0].Arguments)
			assert.Equal(t, expectedNumOutputs, req.Messages.FEval[0].Nargout)

			mockembeddedconnector.RespondWithJSON(t, response, embeddedconnector.ConnectorPayload{
				Messages: embeddedconnector.ConnectorMessage{
					FevalResponse: []embeddedconnector.FevalResponseMessage{
						{
							IsError: false,
							Results: expectedResults,
						},
					},
				},
			})
		},
		nil,
	)
	defer server.Stop()

	client := newClient(t, server.ConnectionDetails())

	// Act
	response, err := client.FEval(t.Context(), logger, entities.FEvalRequest{
		Function:   expectedFunction,
		Arguments:  expectedArguments,
		NumOutputs: expectedNumOutputs,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResults, response.Outputs)
}

func TestClient_Ping_HappyPath(t *testing.T) {
	// Arrange
	logger := testutils.NewInspectableLogger()

	server := mockembeddedconnector.New(t,
		nil,
		func(response http.ResponseWriter, request *http.Request) {
			req := mockembeddedconnector.ReadConnectorRequest(t, request)
			assert.Len(t, req.Messages.Ping, 1)

			mockembeddedconnector.RespondWithJSON(t, response, embeddedconnector.ConnectorPayload{
				Messages: embeddedconnector.ConnectorMessage{
					PingResponse: []embeddedconnector.PingResponseMessage{
						{},
					},
				},
			})
		},
	)
	defer server.Stop()

	client := newClient(t, server.ConnectionDetails())

	// Act
	response := client.Ping(t.Context(), logger)

	// Assert
	assert.True(t, response.IsAlive)
}

func newClient(t *testing.T, connectionDetails embeddedconnector.ConnectionDetails) entities.MATLABSessionClient {
	application := integration.NewEmptyApplication()

	client, err := application.MATLABClientFactory.New(connectionDetails)
	require.NoError(t, err)

	return client
}
