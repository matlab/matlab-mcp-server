// Copyright 2026 The MathWorks, Inc.

package serverlogs

import (
	"fmt"
	"io/fs"
	"strings"
)

// ReadErrorLogs scans server-*.log files in fsys for lines containing "level":"ERROR"
// and returns them. It returns an error if no server log files are found.
func ReadErrorLogs(fsys fs.FS) ([]string, error) {
	logFiles, err := fs.Glob(fsys, "server-*.log")
	if err != nil {
		return nil, fmt.Errorf("failed to glob server logs: %w", err)
	}
	if len(logFiles) == 0 {
		return nil, fmt.Errorf("no server log files found")
	}

	errorLogs := make([]string, 0)
	for _, logFile := range logFiles {
		content, err := fs.ReadFile(fsys, logFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read server log file %s: %w", logFile, err)
		}
		for _, line := range strings.Split(string(content), "\n") {
			if strings.Contains(line, "\"level\":\"ERROR\"") && !isShutdownEOFError(line) {
				errorLogs = append(errorLogs, line)
			}
		}
	}

	return errorLogs, nil
}

// isShutdownEOFError filters benign EOF errors from the MCP go-sdk shutdown race.
func isShutdownEOFError(line string) bool {
	return strings.Contains(line, `"error":"server is closing: EOF"`)
}
