// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	httpclientmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/http/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_Eval_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	expectedCode := "disp('Hello World')"
	expectedOutput := "Hello World"

	responsePayload := embeddedconnector.ConnectorPayload{
		Messages: embeddedconnector.ConnectorMessage{
			EvalResponse: []embeddedconnector.EvalResponseMessage{
				{
					IsError:     false,
					ResponseStr: expectedOutput,
				},
			},
		},
	}
	responseBody, _ := json.Marshal(responsePayload)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(func(req *http.Request) bool {
			payload, ok := parseConnectorRequest(req)
			if !ok {
				return false
			}
			if len(payload.Messages.Eval) != 1 {
				return false
			}
			return payload.Messages.Eval[0].Code == expectedCode
		})).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	// Act
	response, err := client.Eval(t.Context(), mockLogger, entities.EvalRequest{
		Code: expectedCode,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_Eval_HTTPError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(validateConnectorRequest)).
		Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     "500 Internal Server Error",
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	// Act
	response, err := client.Eval(t.Context(), mockLogger, entities.EvalRequest{Code: "ver"})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
	assert.Empty(t, response)
}

func TestClient_Eval_MATLABError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	expectedCode := "invalid_function()"
	expectedErrorMessage := "Undefined function 'invalid_function'"
	expectedError := embeddedconnector.MATLABExceptionsError{
		Exceptions: []embeddedconnector.MATLABException{
			{
				Message: expectedErrorMessage,
			},
		},
	}

	responsePayload := embeddedconnector.ConnectorPayload{
		Messages: embeddedconnector.ConnectorMessage{
			EvalResponse: []embeddedconnector.EvalResponseMessage{
				{
					IsError:     true,
					ResponseStr: expectedErrorMessage,
				},
			},
		},
	}
	responseBody, _ := json.Marshal(responsePayload)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(validateConnectorRequest)).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	evalRequest := entities.EvalRequest{
		Code: expectedCode,
	}

	// Act
	response, err := client.Eval(t.Context(), mockLogger, evalRequest)

	// Assert
	var matlabError embeddedconnector.MATLABExceptionsError
	require.ErrorAs(t, err, &matlabError)
	require.Equal(t, expectedError, matlabError)
	assert.Empty(t, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_Eval_InvalidJSONResponse(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(validateConnectorRequest)).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	// Act
	response, err := client.Eval(t.Context(), mockLogger, entities.EvalRequest{Code: "ver"})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal response")
	assert.Empty(t, response)
}

func TestClient_Eval_NoResponseMessages(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	responsePayload := embeddedconnector.ConnectorPayload{
		Messages: embeddedconnector.ConnectorMessage{
			EvalResponse: []embeddedconnector.EvalResponseMessage{},
		},
	}
	responseBody, _ := json.Marshal(responsePayload)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(validateConnectorRequest)).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	// Act
	response, err := client.Eval(t.Context(), mockLogger, entities.EvalRequest{Code: "ver"})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no response messages received")
	assert.Empty(t, response)
}

func TestClient_Eval_DoErrors(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(validateConnectorRequest)).
		Return(nil, assert.AnError).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	ctx := t.Context()
	evalRequest := entities.EvalRequest{
		Code: "ver",
	}

	// Act
	response, err := client.Eval(ctx, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Empty(t, response)
}

func TestClient_Eval_ReturnsPromptType(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	expectedPromptType := 1

	responsePayload := embeddedconnector.ConnectorPayload{
		Messages: embeddedconnector.ConnectorMessage{
			EvalResponse: []embeddedconnector.EvalResponseMessage{
				{
					IsError:     false,
					ResponseStr: "output",
					PromptType:  expectedPromptType,
				},
			},
		},
	}
	responseBody, _ := json.Marshal(responsePayload)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(validateConnectorRequest)).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	// Act
	response, err := client.Eval(t.Context(), mockLogger, entities.EvalRequest{Code: "dbstack"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPromptType, response.PromptType)
}

func TestClient_Eval_Error_StillReturnsPromptType(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	expectedPromptType := 2

	responsePayload := embeddedconnector.ConnectorPayload{
		Messages: embeddedconnector.ConnectorMessage{
			EvalResponse: []embeddedconnector.EvalResponseMessage{
				{
					IsError:     true,
					ResponseStr: "some error",
					PromptType:  expectedPromptType,
				},
			},
		},
	}
	responseBody, _ := json.Marshal(responsePayload)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(validateConnectorRequest)).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	// Act
	response, err := client.Eval(t.Context(), mockLogger, entities.EvalRequest{Code: "bad_code"})

	// Assert
	require.Error(t, err)
	assert.Equal(t, expectedPromptType, response.PromptType)
}

func TestClient_Eval_ContextPropagation(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	type contextKeyType string
	const contextKey contextKeyType = "uniqueKey"
	const contextKeyValue = "uniqueValue"

	expectedContext := context.WithValue(t.Context(), contextKey, contextKeyValue)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(func(request *http.Request) bool {
			return request.Context().Value(contextKey) == contextKeyValue
		})).
		Return(nil, assert.AnError).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	evalRequest := entities.EvalRequest{
		Code: "ver",
	}

	// Act
	response, err := client.Eval(expectedContext, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Empty(t, response)
}
