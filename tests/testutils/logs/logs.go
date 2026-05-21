// Copyright 2026 The MathWorks, Inc.

package logs

import (
	"io/fs"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/facades/filefacade"
)

func NewReader() Reader {
	return Reader{fileSystem: filefacade.RealFileSystem{}}
}

func CreateTempLogFolder(prefix string) (string, string, error) {
	creator, err := NewFolderCreatorWithFileSystem(filefacade.RealFileSystem{})
	if err != nil {
		return "", "", err
	}

	return creator.CreateTempLogFolder(prefix)
}

func PrepareSessionCLIArgs(args []string, defaultLogLevel string, tempFolderPrefix string) (SessionCLIArgs, error) {
	creator, err := NewFolderCreatorWithFileSystem(filefacade.RealFileSystem{})
	if err != nil {
		return SessionCLIArgs{}, err
	}

	return creator.PrepareSessionCLIArgs(args, defaultLogLevel, tempFolderPrefix)
}

func ReadCombined(logFS fs.FS, globPattern string) (string, error) {
	return readCombined(filefacade.RealFileSystem{}, logFS, globPattern)
}

func ReadEntries(logFS fs.FS, dumpPatterns []DumpPattern) ([]DumpEntry, error) {
	return readEntries(filefacade.RealFileSystem{}, logFS, dumpPatterns)
}
