// Copyright 2026 The MathWorks, Inc.

package mcpclient

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/logs"
)

type LogReader interface {
	ReadCombined(logFS fs.FS, globPattern string) (string, error)
	ReadEntries(logFS fs.FS, dumpPatterns []logs.DumpPattern) ([]logs.DumpEntry, error)
}

type FileSystemProvider interface {
	DirFS(path string) fs.FS
}

type StderrProvider interface {
	Stderr() string
}

type LoggedSessionFactory struct {
	logReader          LogReader
	fileSystemProvider FileSystemProvider
}

type LoggedSession struct {
	*MCPClientSession

	logDir         string
	logFS          fs.FS
	stderrHeader   string
	dumpPatterns   []logs.DumpPattern
	logReader      LogReader
	stderrProvider StderrProvider
}

func NewLoggedSessionFactory(logReader LogReader, fileSystemProvider FileSystemProvider) (*LoggedSessionFactory, error) {
	if logReader == nil {
		return nil, fmt.Errorf("logReader must not be nil")
	}
	if fileSystemProvider == nil {
		return nil, fmt.Errorf("fileSystemProvider must not be nil")
	}

	return &LoggedSessionFactory{
		logReader:          logReader,
		fileSystemProvider: fileSystemProvider,
	}, nil
}

func (f *LoggedSessionFactory) New(
	session *MCPClientSession,
	logDir string,
	stderrHeader string,
	dumpPatterns []logs.DumpPattern,
) (*LoggedSession, error) {
	var stderrProvider StderrProvider
	if session != nil {
		stderrProvider = session
	}

	logFS := f.fileSystemProvider.DirFS(logDir)

	return NewLoggedSession(session, logDir, logFS, stderrHeader, dumpPatterns, f.logReader, stderrProvider)
}

func NewLoggedSession(
	session *MCPClientSession,
	logDir string,
	logFS fs.FS,
	stderrHeader string,
	dumpPatterns []logs.DumpPattern,
	logReader LogReader,
	stderrProvider StderrProvider,
) (*LoggedSession, error) {
	if logReader == nil {
		return nil, fmt.Errorf("logReader must not be nil")
	}
	if logFS == nil {
		return nil, fmt.Errorf("logFS must not be nil")
	}

	return &LoggedSession{
		MCPClientSession: session,
		logDir:           logDir,
		logFS:            logFS,
		stderrHeader:     stderrHeader,
		dumpPatterns:     append([]logs.DumpPattern(nil), dumpPatterns...),
		logReader:        logReader,
		stderrProvider:   stderrProvider,
	}, nil
}

func (s *LoggedSession) LogDir() string {
	return s.logDir
}

func (s *LoggedSession) LogFS() fs.FS {
	return s.logFS
}

func (s *LoggedSession) ReadLogs(globPattern string) (string, error) {
	return s.logReader.ReadCombined(s.logFS, globPattern)
}

func (s *LoggedSession) ReadServerLogs() (string, error) {
	return s.ReadLogs("server-*.log")
}

func (s *LoggedSession) ReadAllServerLogs() (string, error) {
	return s.ReadServerLogs()
}

func (s *LoggedSession) ReadWatchdogLogs() (string, error) {
	return s.ReadLogs("watchdog-*.log")
}

func (s *LoggedSession) CollectDumpData() (string, []logs.DumpEntry, error) {
	entries, err := s.logReader.ReadEntries(s.logFS, s.dumpPatterns)
	if err != nil {
		return "", nil, err
	}
	if s.stderrProvider == nil {
		return "", entries, nil
	}

	return s.stderrProvider.Stderr(), entries, nil
}

func (s *LoggedSession) DumpLogsOnFailure(t *testing.T) {
	t.Helper()
	if !t.Failed() {
		return
	}

	stderr, entries, err := s.CollectDumpData()
	if err != nil {
		t.Logf("Failed to read log files: %s", err.Error())
		return
	}

	if stderr != "" {
		t.Logf("=== %s ===\n%s\n=== End %s ===", s.stderrHeader, stderr, s.stderrHeader)
	}

	for _, entry := range entries {
		t.Logf("=== %s (%s) ===\n%s\n=== End %s ===", entry.Header, entry.File, entry.Content, entry.Header)
	}
}
