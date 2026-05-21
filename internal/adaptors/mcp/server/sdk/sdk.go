// Copyright 2025-2026 The MathWorks, Inc.

package sdk

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Definition interface {
	Name() string
	Title() string
	Instructions() string
	Features() definition.Features
}

type RootStore interface {
	UpdateRoots(roots []*mcp.Root)
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type GlobalMATLAB interface {
	Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error)
}

type MCPSession interface {
	InitializeParams() *mcp.InitializeParams
	ListRoots(ctx context.Context, params *mcp.ListRootsParams) (*mcp.ListRootsResult, error)
}

type Factory struct {
	configFactory ConfigFactory
	definition    Definition
	rootStore     RootStore
	loggerFactory LoggerFactory
	globalMATLAB  GlobalMATLAB
}

type serverCallbackHandler struct {
	config       config.Config
	logger       entities.Logger
	features     definition.Features
	rootStore    RootStore
	globalMATLAB GlobalMATLAB
}

func NewFactory(
	configFactory ConfigFactory,
	definition Definition,
	rootStore RootStore,
	loggerFactory LoggerFactory,
	globalMATLAB GlobalMATLAB,
) *Factory {
	return &Factory{
		configFactory: configFactory,
		definition:    definition,
		rootStore:     rootStore,
		loggerFactory: loggerFactory,
		globalMATLAB:  globalMATLAB,
	}
}

func (f *Factory) NewServer() (*mcp.Server, messages.Error) {
	cfg, err := f.configFactory.Config()
	if err != nil {
		return nil, err
	}

	logger, err := f.loggerFactory.GetGlobalLogger()
	if err != nil {
		return nil, err
	}

	s := &serverCallbackHandler{
		config:       cfg,
		logger:       logger,
		features:     f.definition.Features(),
		rootStore:    f.rootStore,
		globalMATLAB: f.globalMATLAB,
	}

	impl := &mcp.Implementation{
		Name:    f.definition.Name(),
		Title:   f.definition.Title(),
		Version: cfg.Version(),
	}
	options := &mcp.ServerOptions{
		Instructions:            f.definition.Instructions(),
		InitializedHandler:      s.handleInitialized,
		RootsListChangedHandler: s.handleRootsListChanged,
	}

	return mcp.NewServer(impl, options), nil
}

func (s *serverCallbackHandler) handleInitialized(ctx context.Context, req *mcp.InitializedRequest) {
	if req == nil ||
		req.Session == nil {
		return
	}

	s.logClientDetails(req.Session)

	if err := s.updateRoots(ctx, req.Session); err != nil {
		s.logger.WithError(err).Warn("failed to update MCP roots, using fallback starting folder")
	}

	matlabEnabled := s.features.MATLAB.Enabled
	if matlabEnabled && s.config.UseSingleMATLABSession() && s.config.InitializeMATLABOnStartup() {
		if _, err := s.globalMATLAB.Client(ctx, s.logger); err != nil {
			s.logger.WithError(err).Warn("MATLAB eager initialization failed")
		}
	}
}

func (s *serverCallbackHandler) handleRootsListChanged(ctx context.Context, req *mcp.RootsListChangedRequest) {
	if err := s.updateRoots(ctx, req.Session); err != nil {
		s.logger.WithError(err).Warn("failed to update MCP roots, using fallback starting folder")
	}
}

func (s *serverCallbackHandler) logClientDetails(session MCPSession) {
	initializeParams := session.InitializeParams()
	if initializeParams != nil &&
		initializeParams.ClientInfo != nil {
		clientInfo := initializeParams.ClientInfo
		s.logger.
			With("client-name", clientInfo.Name).
			With("client-title", clientInfo.Title).
			With("client-url", clientInfo.WebsiteURL).
			With("client-version", clientInfo.Version).
			Info("New client session")
	}
}

func (s *serverCallbackHandler) updateRoots(ctx context.Context, session MCPSession) error {
	// RootsV2 is the correct pointer field for this check.
	// The legacy Roots field is a value type (not a pointer) due to a go-sdk bug (issue #607),
	// making it impossible to distinguish "no roots support" from "empty roots support".
	params := session.InitializeParams()
	if params == nil || params.Capabilities == nil || params.Capabilities.RootsV2 == nil {
		return nil
	}

	result, err := session.ListRoots(ctx, nil)
	if err != nil {
		return err
	}

	s.rootStore.UpdateRoots(result.Roots)
	s.logger.With("roots", result.Roots).Debug("Updated MCP roots from client")

	return nil
}
