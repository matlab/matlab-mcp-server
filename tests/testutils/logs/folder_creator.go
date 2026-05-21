// Copyright 2026 The MathWorks, Inc.

package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryFileSystem interface {
	MkdirTemp(dir string, pattern string) (string, error)
	MkdirAll(path string, perm os.FileMode) error
}

type FolderCreator struct {
	fileSystem DirectoryFileSystem
}

const (
	logLevelFlagPrefix  = "--log-level="
	logFolderFlagPrefix = "--log-folder="
)

type SessionCLIArgs struct {
	Args        []string
	LogDir      string
	TempBaseDir string
}

func NewFolderCreatorWithFileSystem(fileSystem DirectoryFileSystem) (*FolderCreator, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("fileSystem must not be nil")
	}

	return &FolderCreator{fileSystem: fileSystem}, nil
}

func (c *FolderCreator) CreateTempLogFolder(prefix string) (string, string, error) {
	baseDir, err := c.fileSystem.MkdirTemp("", prefix)
	if err != nil {
		return "", "", err
	}

	logDir := filepath.Join(baseDir, "logs")
	if err := c.fileSystem.MkdirAll(logDir, 0o750); err != nil {
		return "", "", err
	}

	return baseDir, logDir, nil
}

func (c *FolderCreator) PrepareSessionCLIArgs(args []string, defaultLogLevel string, tempFolderPrefix string) (SessionCLIArgs, error) {
	hasLogLevel := false
	hasLogFolder := false
	logDir := ""

	for _, arg := range args {
		if strings.HasPrefix(arg, logLevelFlagPrefix) {
			hasLogLevel = true
		}
		if strings.HasPrefix(arg, logFolderFlagPrefix) {
			hasLogFolder = true
			logDir = strings.TrimPrefix(arg, logFolderFlagPrefix)
		}
	}

	tempBaseDir := ""
	if !hasLogFolder {
		baseDir, createdLogDir, err := c.CreateTempLogFolder(tempFolderPrefix)
		if err != nil {
			return SessionCLIArgs{}, err
		}
		tempBaseDir = baseDir
		logDir = createdLogDir
	}

	defaults := make([]string, 0, 2)
	if !hasLogLevel {
		defaults = append(defaults, logLevelFlagPrefix+defaultLogLevel)
	}
	if !hasLogFolder {
		defaults = append(defaults, logFolderFlagPrefix+logDir)
	}

	return SessionCLIArgs{
		Args:        append(defaults, args...),
		LogDir:      logDir,
		TempBaseDir: tempBaseDir,
	}, nil
}
