// Copyright 2025-2026 The MathWorks, Inc.

package matlabstartingdirselector

import (
	"fmt"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type OSLayer interface {
	UserHomeDir() (string, error)
	Stat(path string) (osfacade.FileInfo, error)
	GOOS() string
}

type RootStore interface {
	GetRoots() []entities.MCPRoot
}

type RootPathResolver interface {
	Resolve(root entities.MCPRoot) (string, error)
}

type MATLABStartingDirSelector struct {
	configFactory    ConfigFactory
	osLayer          OSLayer
	rootStore        RootStore
	rootPathResolver RootPathResolver
}

func New(
	configFactory ConfigFactory,
	osLayer OSLayer,
	rootStore RootStore,
	rootPathResolver RootPathResolver,
) *MATLABStartingDirSelector {
	return &MATLABStartingDirSelector{
		configFactory:    configFactory,
		osLayer:          osLayer,
		rootStore:        rootStore,
		rootPathResolver: rootPathResolver,
	}
}

func (s *MATLABStartingDirSelector) SelectMATLABStartingDir(logger entities.Logger) (string, error) {
	config, configErr := s.configFactory.Config()
	if configErr != nil {
		return "", configErr
	}

	// 1. Try preferred directory (--initial-working-folder flag)
	if preferredDir := config.PreferredMATLABStartingDirectory(); preferredDir != "" {
		if _, err := s.osLayer.Stat(preferredDir); err != nil {
			return "", err
		}
		return preferredDir, nil
	}

	// 2. Try first MCP root from the client
	if dir, err := s.getFirstRootDir(); err != nil {
		logger.WithError(err).Warn("failed to use MCP root as starting directory, falling back to default")
	} else if dir != "" {
		return dir, nil
	}

	// 3. Fall back to documents directory
	dir, err := s.getDocumentsDir()
	if err != nil {
		return "", err
	}

	if _, err := s.osLayer.Stat(dir); err != nil {
		return "", err
	}

	return dir, nil
}

func (s *MATLABStartingDirSelector) getFirstRootDir() (string, error) {
	roots := s.rootStore.GetRoots()
	if len(roots) == 0 {
		return "", nil
	}

	dir, err := s.rootPathResolver.Resolve(roots[0])
	if err != nil {
		return "", err
	}

	if dir == "" {
		return "", nil
	}

	info, err := s.osLayer.Stat(dir)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf("MCP root path is not a directory: %s", dir)
	}

	return dir, nil
}

func (s *MATLABStartingDirSelector) getDocumentsDir() (string, error) {
	home, err := s.osLayer.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch s.osLayer.GOOS() {
	case "windows", "darwin":
		return filepath.Join(home, "Documents"), nil
	default: // Linux - Documents less commonly used
		return home, nil // Just return home for Linux
	}
}
