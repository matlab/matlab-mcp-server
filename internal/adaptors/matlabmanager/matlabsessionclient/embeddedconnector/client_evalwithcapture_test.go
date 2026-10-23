// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
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

//go:embed testdata/matlab_error.json
var matlabErrorJSON string

//go:embed testdata/matlab_error_after_warning.json
var matlabErrorAfterWarningJSON string

//go:embed testdata/multiple_matlab_errors.json
var multipleMATLABErrorsJSON string

func TestClient_EvalWithCapture_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	expectedCode := "disp('Hello World')"
	expectedOutput := "Hello World"

	entries := []embeddedconnector.LiveEditorResponseEntry{
		{
			Type:     "execute_result",
			MimeType: []string{"text/plain"},
			Value:    []json.RawMessage{json.RawMessage(`"` + expectedOutput + `"`)},
		},
	}
	responseBody := buildEvalWithCaptureResponse(t, entries)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(func(req *http.Request) bool {
			payload, ok := parseConnectorRequest(req)
			if !ok {
				return false
			}
			if len(payload.Messages.FEval) != 1 {
				return false
			}
			feval := payload.Messages.FEval[0]
			return feval.Function == "matlab_mcp.mcpEval" &&
				len(feval.Arguments) == 1 &&
				feval.Arguments[0] == expectedCode
		})).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	// Act
	response, err := client.EvalWithCapture(t.Context(), mockLogger, entities.EvalRequest{
		Code: expectedCode,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_EvalWithCapture_ReturnImages(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	expectedImageData := []byte("image data")
	expectedImageBase64 := base64.StdEncoding.EncodeToString(expectedImageData)

	entries := []embeddedconnector.LiveEditorResponseEntry{
		{
			Type:     "execute_result",
			MimeType: []string{"image/png"},
			Value:    []json.RawMessage{json.RawMessage(`"` + expectedImageBase64 + `"`)},
		},
	}
	responseBody := buildEvalWithCaptureResponse(t, entries)

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
	response, err := client.EvalWithCapture(t.Context(), mockLogger, entities.EvalRequest{Code: "plot(1:10)"})

	// Assert
	require.NoError(t, err)
	assert.Empty(t, response.ConsoleOutput)
	assert.Equal(t, [][]byte{expectedImageData}, response.Images)
}

func TestClient_EvalWithCapture_ReturnStreams(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	ctx := t.Context()
	evalRequest := entities.EvalRequest{
		Code: "warning('some warning')",
	}
	expectedOutput := "Warning: some warning"

	entries := []embeddedconnector.LiveEditorResponseEntry{
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: expectedOutput,
				Name: "stderr",
			},
		},
	}
	responseBody := buildEvalWithCaptureResponse(t, entries)

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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_EvalWithCapture_MultipleStreams_SameName(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	entries := []embeddedconnector.LiveEditorResponseEntry{
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: "line1\n",
				Name: "stdout",
			},
		},
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: "line2\n",
				Name: "stdout",
			},
		},
	}
	responseBody := buildEvalWithCaptureResponse(t, entries)

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
	response, err := client.EvalWithCapture(t.Context(), mockLogger, entities.EvalRequest{Code: "disp('line1'); disp('line2')"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "line1\nline2\n", response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_EvalWithCapture_MultipleStreams_DifferentNames(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	entries := []embeddedconnector.LiveEditorResponseEntry{
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: "output",
				Name: "stdout",
			},
		},
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: "Warning: warning message",
				Name: "stderr",
			},
		},
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: " continued",
				Name: "stderr",
			},
		},
	}
	responseBody := buildEvalWithCaptureResponse(t, entries)

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
	response, err := client.EvalWithCapture(t.Context(), mockLogger, entities.EvalRequest{Code: "fprintf('output'); warning('warning message')"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "output\nWarning: warning message continued", response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_EvalWithCapture_MixedStreamsAndResults(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	expectedImageData := []byte("plot_image_data")
	expectedImageBase64 := base64.StdEncoding.EncodeToString(expectedImageData)

	entries := []embeddedconnector.LiveEditorResponseEntry{
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: "calculating\n",
				Name: "stdout",
			},
		},
		{
			Type:     "execute_result",
			MimeType: []string{"text/plain"},
			Value:    []json.RawMessage{json.RawMessage(`"y = 10"`)},
		},
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: "more output",
				Name: "stdout",
			},
		},
		{
			Type:     "execute_result",
			MimeType: []string{"image/png"},
			Value:    []json.RawMessage{json.RawMessage(`"` + expectedImageBase64 + `"`)},
		},
	}
	responseBody := buildEvalWithCaptureResponse(t, entries)

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
	response, err := client.EvalWithCapture(t.Context(), mockLogger, entities.EvalRequest{Code: "x = 5; disp('calculating'); y = x * 2; plot(1:y)"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "calculating\n\ny = 10\nmore output", response.ConsoleOutput)
	assert.Equal(t, [][]byte{expectedImageData}, response.Images)
}

func TestClient_EvalWithCapture_StreamsWithInterruptionByExecuteResult(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	entries := []embeddedconnector.LiveEditorResponseEntry{
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: "Warning: first",
				Name: "stderr",
			},
		},
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: " warning message",
				Name: "stderr",
			},
		},
		{
			Type:     "execute_result",
			MimeType: []string{"text/plain"},
			Value:    []json.RawMessage{json.RawMessage(`"x = 1"`)},
		},
		{
			Type: "stream",
			Content: struct {
				Text string `json:"text"`
				Name string `json:"name"`
			}{
				Text: "Warning: second warning",
				Name: "stderr",
			},
		},
	}
	responseBody := buildEvalWithCaptureResponse(t, entries)

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
	response, err := client.EvalWithCapture(t.Context(), mockLogger, entities.EvalRequest{Code: "warning('first'); x = 1; warning('second')"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Warning: first warning message\nx = 1\nWarning: second warning", response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_EvalWithCapture_MATLABError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	ctx := t.Context()
	evalRequest := entities.EvalRequest{
		Code: "disp('before'); undefinedFnXyz()",
	}
	expectedError := embeddedconnector.MATLABExceptionsError{
		ConsoleOutput: "before\n",
		Exceptions: []embeddedconnector.MATLABException{
			{
				Identifier: "MATLAB:UndefinedFunction",
				Message:    "Unrecognized function or variable 'undefinedFnXyz'.",
			},
		},
	}
	responseBody := buildEvalWithCaptureResponseFromJSON(t, matlabErrorJSON)

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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	var matlabError embeddedconnector.MATLABExceptionsError
	require.ErrorAs(t, err, &matlabError)
	require.Equal(t, expectedError, matlabError)
	assert.Empty(t, response)
}

func TestClient_EvalWithCapture_MATLABError_AfterWarning(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	ctx := t.Context()
	evalRequest := entities.EvalRequest{
		Code: "warning('careful'); error('boom')",
	}
	expectedError := embeddedconnector.MATLABExceptionsError{
		ConsoleOutput: "Warning: careful",
		Exceptions: []embeddedconnector.MATLABException{
			{
				Message: "boom",
			},
		},
	}
	responseBody := buildEvalWithCaptureResponseFromJSON(t, matlabErrorAfterWarningJSON)

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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	var matlabError embeddedconnector.MATLABExceptionsError
	require.ErrorAs(t, err, &matlabError)
	require.Equal(t, expectedError, matlabError)
	assert.Empty(t, response)
}

func TestClient_EvalWithCapture_MultipleMATLABErrors(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	ctx := t.Context()
	evalRequest := entities.EvalRequest{
		Code: "runAll()",
	}
	expectedError := embeddedconnector.MATLABExceptionsError{
		Exceptions: []embeddedconnector.MATLABException{
			{
				Identifier: "first:id",
				Message:    "first failure",
			},
			{
				Identifier: "second:id",
				Message:    "second failure",
			},
		},
	}
	responseBody := buildEvalWithCaptureResponseFromJSON(t, multipleMATLABErrorsJSON)

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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	var matlabError embeddedconnector.MATLABExceptionsError
	require.ErrorAs(t, err, &matlabError)
	require.Equal(t, expectedError, matlabError)
	assert.Empty(t, response)
}

func TestClient_EvalWithCapture_DoErrors(t *testing.T) {
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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Empty(t, response)
}

func TestClient_EvalWithCapture_ContextPropagation(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	type contextKeyType string
	const contextKey contextKeyType = "uniqueKey"
	const contextKeyValue = "uniqueValue"

	expectedContext := context.WithValue(t.Context(), contextKey, contextKeyValue)
	expectedError := assert.AnError

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
	response, err := client.EvalWithCapture(expectedContext, mockLogger, evalRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response)
}
