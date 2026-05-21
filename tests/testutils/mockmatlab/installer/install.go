// Copyright 2026 The MathWorks, Inc.

package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

const versionInfoXML = `<?xml version="1.0" encoding="UTF-8"?>
<MathWorks_version_info>
  <version>25.1.0.0</version>
  <release>R2025a</release>
  <description>Update 1</description>
  <date>Jan 01 2025</date>
</MathWorks_version_info>`

type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type ModuleRootFinder interface {
	FindModuleRoot() (string, error)
}

type BinaryBuilder interface {
	BuildPlatformSpecificBinaries(moduleDir, binDir string) error
}

type Installer struct {
	fileSystem       FileSystem
	moduleRootFinder ModuleRootFinder
	binaryBuilder    BinaryBuilder
}

func New(fileSystem FileSystem, moduleRootFinder ModuleRootFinder, binaryBuilder BinaryBuilder) *Installer {
	return &Installer{
		fileSystem:       fileSystem,
		moduleRootFinder: moduleRootFinder,
		binaryBuilder:    binaryBuilder,
	}
}

func (i *Installer) BuildAndInstall(matlabRoot string) error {
	binDir := filepath.Join(matlabRoot, "bin")
	if err := i.fileSystem.MkdirAll(binDir, 0o700); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	moduleDir, err := i.moduleRootFinder.FindModuleRoot()
	if err != nil {
		return fmt.Errorf("failed to find module root: %w", err)
	}

	if err := i.binaryBuilder.BuildPlatformSpecificBinaries(moduleDir, binDir); err != nil {
		return fmt.Errorf("failed to build mock MATLAB binaries: %w", err)
	}

	if err := i.fileSystem.WriteFile(filepath.Join(matlabRoot, "VersionInfo.xml"), []byte(versionInfoXML), 0o600); err != nil {
		return fmt.Errorf("failed to write VersionInfo.xml: %w", err)
	}

	return nil
}
