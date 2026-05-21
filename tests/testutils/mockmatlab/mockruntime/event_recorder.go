// Copyright 2026 The MathWorks, Inc.

package mockruntime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	EventStarted       = "started"
	EventEval          = "eval"
	EventFeval         = "feval"
	EventExitRequested = "exit_requested"
)

type Event struct {
	Type     string `json:"type"`
	Mode     string `json:"mode,omitempty"`
	Code     string `json:"code,omitempty"`
	Function string `json:"function,omitempty"`
}

type EventRecorder struct {
	ID      int
	logFile *os.File
	encoder *json.Encoder
}

// NewEventRecorder claims the next auto-incrementing ID in logDir and opens
// a JSONL file at mock-matlab-<id>.jsonl. The caller must call Close when done.
func NewEventRecorder(logDir string) (*EventRecorder, error) {
	if logDir == "" {
		return nil, fmt.Errorf("log directory must not be empty")
	}

	if err := os.MkdirAll(logDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	id := nextRecorderID(logDir)

	logPath := filepath.Join(logDir, fmt.Sprintf("mock-matlab-%d.jsonl", id))
	f, err := os.Create(logPath) //nolint:gosec // Trusted test path
	if err != nil {
		return nil, fmt.Errorf("failed to create event recorder file: %w", err)
	}

	return &EventRecorder{ID: id, logFile: f, encoder: json.NewEncoder(f)}, nil
}

func (r *EventRecorder) RecordStarted(mode string) {
	r.record(Event{Type: EventStarted, Mode: mode})
}

func (r *EventRecorder) RecordEval(code string) {
	r.record(Event{Type: EventEval, Code: code})
}

func (r *EventRecorder) RecordFeval(function string) {
	r.record(Event{Type: EventFeval, Function: function})
}

func (r *EventRecorder) RecordExitRequested() {
	r.record(Event{Type: EventExitRequested})
}

func (r *EventRecorder) Close() {
	_ = r.logFile.Close()
}

func (r *EventRecorder) record(event Event) {
	_ = r.encoder.Encode(event)
}

// nextRecorderID scans logDir for existing mock-matlab-*.jsonl files and returns max+1.
func nextRecorderID(logDir string) int {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return 1
	}

	maxID := 0
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "mock-matlab-") || !strings.HasSuffix(name, ".jsonl") {
			continue
		}
		numStr := strings.TrimPrefix(name, "mock-matlab-")
		numStr = strings.TrimSuffix(numStr, ".jsonl")
		n, err := strconv.Atoi(numStr)
		if err != nil {
			continue // Skip files that don't match the mock-matlab-<ID>.jsonl naming convention.
		}
		if n > maxID {
			maxID = n
		}
	}
	return maxID + 1
}
