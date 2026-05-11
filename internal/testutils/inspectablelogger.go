// Copyright 2025-2026 The MathWorks, Inc.

package testutils

import (
	"fmt"
	"maps"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/stretchr/testify/mock"
)

type fields map[string]any

type InspectableLogs map[string]fields

// InspectableLogger is designed as a drop in replacement for a real logger that stores information
// internally and can be inspected during testing to assert on logs that have been made.
//
// When the logger is passed between functions the `AsMockArg` function should be added to the expectation
// so that internal details do not break the matching of arguments.
type InspectableLogger struct {
	// Store the log lists as pointers so they can be shared with child loggers created by With calls
	debugLogs InspectableLogs
	infoLogs  InspectableLogs
	warnLogs  InspectableLogs
	errorLogs InspectableLogs
	// Each logger needs to store its own list of fields, so Fields is not a pointer
	Fields fields

	lock *sync.Mutex
}

// NewInspectableLogger creates a new InspectableLogger with empty log and field lists
func NewInspectableLogger() *InspectableLogger {
	return &InspectableLogger{
		debugLogs: make(InspectableLogs),
		infoLogs:  make(InspectableLogs),
		warnLogs:  make(InspectableLogs),
		errorLogs: make(InspectableLogs),
		Fields:    make(fields),

		lock: new(sync.Mutex),
	}
}

// DebugLogs returns a copy of the debug logs stored by the InspectableLogger
func (il *InspectableLogger) DebugLogs() InspectableLogs {
	return il.cloneLogs(il.debugLogs)
}

// InfoLogs returns a copy of the info logs stored by the InspectableLogger
func (il *InspectableLogger) InfoLogs() InspectableLogs {
	return il.cloneLogs(il.infoLogs)
}

// WarnLogs returns a copy of the warn logs stored by the InspectableLogger
func (il *InspectableLogger) WarnLogs() InspectableLogs {
	return il.cloneLogs(il.warnLogs)
}

// ErrorLogs returns a copy of the error logs stored by the InspectableLogger
func (il *InspectableLogger) ErrorLogs() InspectableLogs {
	return il.cloneLogs(il.errorLogs)
}

func (il *InspectableLogger) cloneLogs(logs InspectableLogs) InspectableLogs {
	il.lock.Lock()
	defer il.lock.Unlock()

	return maps.Clone(logs)
}

// Debug stores a debug log in the InspectableLogger
func (il *InspectableLogger) Debug(msg string) {
	il.lock.Lock()
	defer il.lock.Unlock()

	addLog(il.debugLogs, msg, il.Fields)
}

// Info stores a info log in the InspectableLogger
func (il *InspectableLogger) Info(msg string) {
	il.lock.Lock()
	defer il.lock.Unlock()

	addLog(il.infoLogs, msg, il.Fields)
}

// Warn stores a warn log in the InspectableLogger
func (il *InspectableLogger) Warn(msg string) {
	il.lock.Lock()
	defer il.lock.Unlock()

	addLog(il.warnLogs, msg, il.Fields)
}

// Error stores a error log in the InspectableLogger
func (il *InspectableLogger) Error(msg string) {
	il.lock.Lock()
	defer il.lock.Unlock()

	addLog(il.errorLogs, msg, il.Fields)
}

// With returns a new InspectableLogger with references to the previous logger's log lists.
// The fields from the previous logger are copied over and the new field added.
// This way the new logger stores it's own fields, but the lists of logs are shared between
// all inspectable loggers that have been spawned from the original instance.
func (il *InspectableLogger) With(key string, value interface{}) entities.Logger {
	il.lock.Lock()
	defer il.lock.Unlock()

	newFields := copyFields(il.Fields)

	newFields[key] = value

	return &InspectableLogger{
		debugLogs: il.debugLogs,
		infoLogs:  il.infoLogs,
		warnLogs:  il.warnLogs,
		errorLogs: il.errorLogs,
		Fields:    newFields,

		lock: il.lock,
	}
}

// WithError calls With but with the key already defined
func (il *InspectableLogger) WithError(err error) entities.Logger {
	const key = "error"
	return il.With(key, err)
}

func (il *InspectableLogger) AsMockArg() any {
	return mock.AnythingOfType(fmt.Sprintf("%T", il))
}

func copyFields(m fields) fields {
	newFields := make(fields)
	maps.Copy(newFields, m)
	return newFields
}

func addLog(inspectableLogs InspectableLogs, msg string, fields fields) {
	newFields := copyFields(fields)
	inspectableLogs[msg] = newFields
}
