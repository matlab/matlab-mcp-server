// Copyright 2026 The MathWorks, Inc.

package logs

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type LogFileSystem interface {
	Glob(logFS fs.FS, globPattern string) ([]string, error)
	ReadFile(logFS fs.FS, logFile string) ([]byte, error)
}

type Reader struct {
	fileSystem LogFileSystem
}

func NewReaderWithFileSystem(fileSystem LogFileSystem) (Reader, error) {
	if fileSystem == nil {
		return Reader{}, fmt.Errorf("fileSystem must not be nil")
	}

	return Reader{fileSystem: fileSystem}, nil
}

func (r Reader) ReadCombined(logFS fs.FS, globPattern string) (string, error) {
	return readCombined(r.fileSystem, logFS, globPattern)
}

func (r Reader) ReadEntries(logFS fs.FS, dumpPatterns []DumpPattern) ([]DumpEntry, error) {
	return readEntries(r.fileSystem, logFS, dumpPatterns)
}

func readCombined(fileSystem LogFileSystem, logFS fs.FS, globPattern string) (string, error) {
	logFiles, err := fileSystem.Glob(logFS, globPattern)
	if err != nil {
		return "", fmt.Errorf("failed to glob logs: %w", err)
	}
	if len(logFiles) == 0 {
		return "", fmt.Errorf("no logs found for pattern %s", globPattern)
	}

	var combined strings.Builder
	for _, logFile := range logFiles {
		content, err := fileSystem.ReadFile(logFS, logFile)
		if err != nil {
			return "", fmt.Errorf("failed to read log file %s: %w", logFile, err)
		}
		combined.Write(content)
	}

	return combined.String(), nil
}

func readEntries(fileSystem LogFileSystem, logFS fs.FS, dumpPatterns []DumpPattern) ([]DumpEntry, error) {
	entries := make([]DumpEntry, 0)
	for _, dumpPattern := range dumpPatterns {
		logFiles, err := fileSystem.Glob(logFS, dumpPattern.Glob)
		if err != nil {
			return nil, fmt.Errorf("failed to glob logs for pattern %s: %w", dumpPattern.Glob, err)
		}

		for _, logFile := range logFiles {
			content, err := fileSystem.ReadFile(logFS, logFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read log file %s: %w", logFile, err)
			}

			entries = append(entries, DumpEntry{
				Header:  dumpPattern.Header,
				File:    filepath.Base(logFile),
				Content: string(content),
			})
		}
	}

	return entries, nil
}
