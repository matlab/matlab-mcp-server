// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type ConnectorPayload struct {
	Messages ConnectorMessage `json:"messages"`
}

type ConnectorMessage struct {
	Eval          []EvalMessage          `json:"Eval,omitempty"`
	FEval         []FevalMessage         `json:"FEval,omitempty"`
	EvalResponse  []EvalResponseMessage  `json:"EvalResponse,omitempty"`
	FevalResponse []FevalResponseMessage `json:"FEvalResponse,omitempty"`
	Ping          []PingMessage          `json:"Ping,omitempty"`
	PingResponse  []PingResponseMessage  `json:"PingResponse,omitempty"`
}

type EvalMessage struct {
	Code string `json:"mcode"`
}

type EvalResponseMessage struct {
	IsError     bool   `json:"isError"`
	ResponseStr string `json:"responseStr"`
	PromptType  int    `json:"promptType"`
}

type FevalMessage struct {
	Function  string   `json:"function"`
	Arguments []string `json:"arguments"`
	Nargout   int      `json:"nargout"`
}

type FevalResponseMessage struct {
	IsError       bool              `json:"isError"`
	MessageFaults []json.RawMessage `json:"messageFaults"`
	Results       []interface{}     `json:"results"`
}

type PingMessage struct {
}

type PingResponseMessage struct {
	MessageFaults []json.RawMessage `json:"messageFaults"`
}

type Fault struct {
	Message string `json:"message"`
}

type LiveEditorResponseEntry struct {
	Type     string            `json:"type"`
	MimeType []string          `json:"mimetype"`
	Value    []json.RawMessage `json:"value"`
	Content  struct {
		Text string `json:"text"`
		Name string `json:"name"`
	} `json:"content"`
	IsError    bool   `json:"isError"`
	Identifier string `json:"identifier"`
}

type responseProcessor struct {
	consoleOutput        []string
	images               [][]byte
	exceptions           []matlabException
	pendingStreamName    string
	pendingStreamContent string
}

func (p *responseProcessor) processEntry(entry LiveEditorResponseEntry) error {
	switch entry.Type {
	case "execute_result":
		p.flushPendingStream()
		err := p.processExecuteResult(entry)
		if err != nil {
			return err
		}
	case "stream":
		if entry.IsError {
			p.flushPendingStream()
			p.exceptions = append(p.exceptions, matlabException{
				Identifier: entry.Identifier,
				Message:    entry.Content.Text,
			})
			return nil
		}
		p.processStream(entry)
	}
	return nil
}

func (p *responseProcessor) processExecuteResult(entry LiveEditorResponseEntry) error {
	for i, mimeType := range entry.MimeType {
		if i >= len(entry.Value) {
			continue // Safety check
		}
		switch mimeType {
		case "text/plain":
			var value string
			err := json.Unmarshal(entry.Value[i], &value)
			if err != nil {
				return err
			}
			p.consoleOutput = append(p.consoleOutput, value)
		case "image/png":
			var value []byte
			err := json.Unmarshal(entry.Value[i], &value)
			if err != nil {
				return err
			}
			p.images = append(p.images, value)
		}
	}
	return nil
}

func (p *responseProcessor) processStream(entry LiveEditorResponseEntry) {
	// If we have a different stream name, flush the previous one
	if p.pendingStreamName != entry.Content.Name {
		p.flushPendingStream()
	}

	p.pendingStreamName = entry.Content.Name
	p.pendingStreamContent += entry.Content.Text
}

func (p *responseProcessor) flushPendingStream() {
	if p.pendingStreamName != "" {
		p.consoleOutput = append(p.consoleOutput, p.pendingStreamContent)
		p.pendingStreamName = ""
		p.pendingStreamContent = ""
	}
}

func parseEvalWithCaptureResponse(response entities.FEvalResponse) (entities.EvalResponse, error) {
	if len(response.Outputs) != 1 {
		return entities.EvalResponse{}, fmt.Errorf("unexpected number of outputs from MATLAB session")
	}

	consoleData, ok := response.Outputs[0].(string)
	if !ok {
		return entities.EvalResponse{}, fmt.Errorf("failed to cast output to string")
	}

	var parsedResponse []json.RawMessage
	err := json.Unmarshal([]byte(consoleData), &parsedResponse)
	if err != nil {
		return entities.EvalResponse{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	processor := &responseProcessor{}

	for _, rawEntry := range parsedResponse {
		var entry LiveEditorResponseEntry
		if err := json.Unmarshal(rawEntry, &entry); err != nil {
			continue // Ignore invalid entries
		}
		if err := processor.processEntry(entry); err != nil {
			return entities.EvalResponse{}, err
		}
	}

	processor.flushPendingStream()

	consoleOutput := strings.Join(processor.consoleOutput, "\n")
	if len(processor.exceptions) > 0 {
		return entities.EvalResponse{}, matlabExceptionsError{
			ConsoleOutput: consoleOutput,
			Exceptions:    processor.exceptions,
		}
	}

	return entities.EvalResponse{
		ConsoleOutput: consoleOutput,
		Images:        processor.images,
	}, nil
}
