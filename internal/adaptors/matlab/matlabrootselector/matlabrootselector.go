// Copyright 2025-2026 The MathWorks, Inc.

package matlabrootselector

import (
	"context"
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type MATLABManager interface {
	ListEnvironments(ctx context.Context, sessionLogger entities.Logger) []entities.EnvironmentInfo
}

type MATLABRootSelector struct {
	configFactory ConfigFactory
	matlabManager MATLABManager
}

func New(
	configFactory ConfigFactory,
	matlabManager MATLABManager,
) *MATLABRootSelector {
	return &MATLABRootSelector{
		configFactory: configFactory,
		matlabManager: matlabManager,
	}
}

func (m *MATLABRootSelector) SelectMATLABRoot(ctx context.Context, logger entities.Logger) (string, error) {
	config, err := m.configFactory.Config()
	if err != nil {
		return "", err
	}

	if preferredLocalMATLABRoot := config.PreferredLocalMATLABRoot(); preferredLocalMATLABRoot != "" {
		return preferredLocalMATLABRoot, nil
	}

	environments := m.matlabManager.ListEnvironments(ctx, logger)
	if len(environments) == 0 {
		return "", fmt.Errorf("no valid MATLAB environments found")
	}

	return environments[0].MATLABRoot, nil
}
