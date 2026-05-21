// Copyright 2026 The MathWorks, Inc.

package mockruntime

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type InstanceEvents struct {
	ID     int
	Events []Event
}

// HasEvent returns true if the instance recorded an event with the given type and matching fields.
func (ie InstanceEvents) HasEvent(eventType string) bool {
	for _, e := range ie.Events {
		if e.Type == eventType {
			return true
		}
	}
	return false
}

// HasEval returns true if the instance recorded an eval event with the given code.
func (ie InstanceEvents) HasEval(code string) bool {
	for _, e := range ie.Events {
		if e.Type == EventEval && e.Code == code {
			return true
		}
	}
	return false
}

// HasFeval returns true if the instance recorded a feval event with the given function name.
func (ie InstanceEvents) HasFeval(function string) bool {
	for _, e := range ie.Events {
		if e.Type == EventFeval && e.Function == function {
			return true
		}
	}
	return false
}

// StartedMode returns the mode from the started event, or empty string if none.
func (ie InstanceEvents) StartedMode() string {
	for _, e := range ie.Events {
		if e.Type == EventStarted {
			return e.Mode
		}
	}
	return ""
}

// EventTypes returns the type of each event, in order.
func (ie InstanceEvents) EventTypes() []string {
	types := make([]string, len(ie.Events))
	for i, e := range ie.Events {
		types[i] = e.Type
	}
	return types
}

// HasEvalsInOrder returns true if the given codes appear as eval events in the
// given relative order (not necessarily contiguous — other evals may appear between them).
func (ie InstanceEvents) HasEvalsInOrder(codes ...string) bool {
	idx := 0
	for _, e := range ie.Events {
		if idx < len(codes) && e.Type == EventEval && e.Code == codes[idx] {
			idx++
		}
	}
	return idx == len(codes)
}

// HasEvalMatching returns true if the instance recorded an eval event whose code
// satisfies the predicate.
func (ie InstanceEvents) HasEvalMatching(predicate func(code string) bool) bool {
	for _, e := range ie.Events {
		if e.Type == EventEval && predicate(e.Code) {
			return true
		}
	}
	return false
}

// CountEvent returns how many events of the given type were recorded.
func (ie InstanceEvents) CountEvent(eventType string) int {
	n := 0
	for _, e := range ie.Events {
		if e.Type == eventType {
			n++
		}
	}
	return n
}

// ReadEventsForAllInstances reads all mock-matlab-*.jsonl files from logDir, sorted by instance ID.
// Each entry in the returned slice represents one mock MATLAB process lifetime.
func ReadEventsForAllInstances(logDir string) ([]InstanceEvents, error) {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read log directory: %w", err)
	}

	var instances []InstanceEvents
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "mock-matlab-") || !strings.HasSuffix(name, ".jsonl") {
			continue
		}
		numStr := strings.TrimPrefix(name, "mock-matlab-")
		numStr = strings.TrimSuffix(numStr, ".jsonl")
		id, err := strconv.Atoi(numStr)
		if err != nil {
			continue // Skip files that don't match the mock-matlab-<ID>.jsonl naming convention.
		}

		events, err := readEventsFromFile(filepath.Join(logDir, name))
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", name, err)
		}

		instances = append(instances, InstanceEvents{ID: id, Events: events})
	}

	sort.Slice(instances, func(i, j int) bool { return instances[i].ID < instances[j].ID })
	return instances, nil
}

func readEventsFromFile(path string) ([]Event, error) {
	f, err := os.Open(path) //nolint:gosec // Trusted test path
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var events []Event
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var event Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, fmt.Errorf("failed to parse event: %w", err)
		}
		events = append(events, event)
	}
	return events, scanner.Err()
}
