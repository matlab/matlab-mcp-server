// Copyright 2026 The MathWorks, Inc.

package sdk

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func HandleInitialized(cfg config.Config, logger entities.Logger, features definition.Features, rs RootStore, gm GlobalMATLAB) func(context.Context, *mcp.InitializedRequest) {
	s := &serverCallbackHandler{config: cfg, logger: logger, features: features, rootStore: rs, globalMATLAB: gm}
	return s.handleInitialized
}

func HandleRootsListChanged(cfg config.Config, logger entities.Logger, features definition.Features, rs RootStore, gm GlobalMATLAB) func(context.Context, *mcp.RootsListChangedRequest) {
	s := &serverCallbackHandler{config: cfg, logger: logger, features: features, rootStore: rs, globalMATLAB: gm}
	return s.handleRootsListChanged
}

func UpdateRoots(cfg config.Config, logger entities.Logger, features definition.Features, rs RootStore, gm GlobalMATLAB) func(context.Context, MCPSession) error {
	s := &serverCallbackHandler{config: cfg, logger: logger, features: features, rootStore: rs, globalMATLAB: gm}
	return s.updateRoots
}

func LogClientDetails(logger entities.Logger) func(MCPSession) {
	s := &serverCallbackHandler{logger: logger}
	return s.logClientDetails
}
