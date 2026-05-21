// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type MATLABServices interface {
	ListDiscoveredMatlabInfo(logger entities.Logger) datatypes.ListMatlabInfo
	StartLocalMATLABSession(ctx context.Context, logger entities.Logger, request datatypes.LocalSessionDetails) (embeddedconnector.ConnectionDetails, func() error, error)
}

type MATLABSessionStore interface {
	Add(client matlabsessionstore.MATLABSessionClientWithCleanup) entities.SessionID
	Get(sessionID entities.SessionID) (matlabsessionstore.MATLABSessionClientWithCleanup, error)
	Remove(sessionID entities.SessionID)
}

type MATLABSessionClientFactory interface {
	New(endpoint embeddedconnector.ConnectionDetails) (entities.MATLABSessionClient, error)
}

type SessionSelector interface {
	SelectSessionToAttachTo(logger entities.Logger) (embeddedconnector.ConnectionDetails, error)
}

type MATLABManager struct {
	configFactory   ConfigFactory
	matlabServices  MATLABServices
	sessionStore    MATLABSessionStore
	clientFactory   MATLABSessionClientFactory
	sessionSelector SessionSelector

	matlabSessionConnectionRetryInterval time.Duration
}

var _ entities.MATLABManager = (*MATLABManager)(nil)

func New(
	configFactory ConfigFactory,
	matlabServices MATLABServices,
	sessionStore MATLABSessionStore,
	clientFactory MATLABSessionClientFactory,
	sessionSelector SessionSelector,
) *MATLABManager {
	return &MATLABManager{
		configFactory:   configFactory,
		matlabServices:  matlabServices,
		sessionStore:    sessionStore,
		clientFactory:   clientFactory,
		sessionSelector: sessionSelector,

		matlabSessionConnectionRetryInterval: defaultMATLABSessionConnectionRetryInterval,
	}
}
