// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

var ErrMATLABSessionNotAlive = errors.New("session is not alive")

func (m *MATLABManager) StartMATLABSession(ctx context.Context, sessionLogger entities.Logger, startRequest entities.SessionDetails) (entities.SessionID, error) {
	var zeroValue entities.SessionID
	var client matlabsessionstore.MATLABSessionClientWithCleanup

	switch request := startRequest.(type) {
	case entities.LocalSessionDetails:
		localSessionLogger := sessionLogger.With("matlab-root", request.MATLABRoot)
		// For now, we return embedded connector details, to decouple the session start logic from the client creation.
		embeddedConnectorEndpoint, sessionCleanup, err := m.matlabServices.StartLocalMATLABSession(
			ctx,
			localSessionLogger,
			datatypes.LocalSessionDetails{
				MATLABRoot:             request.MATLABRoot,
				IsStartingDirectorySet: request.IsStartingDirectorySet,
				StartingDirectory:      request.StartingDirectory,
				ShowMATLABDesktop:      request.ShowMATLABDesktop,
			},
		)
		if err != nil {
			return zeroValue, err
		}
		embeddedConnectorClient, err := m.clientFactory.New(embeddedConnectorEndpoint)
		if err != nil {
			if cleanupErr := sessionCleanup(); cleanupErr != nil {
				sessionLogger.WithError(cleanupErr).Error("Failed to clean up session after client factory error")
			}
			return zeroValue, err
		}
		client = newMATLABSessionClientWithCleanup(embeddedConnectorClient, sessionCleanup)
	case entities.AttachToExistingSession:
		sessionLogger.Info("Attaching to existing session")

		connectionDetails, err := m.sessionSelector.SelectSessionToAttachTo(sessionLogger)
		if err != nil {
			return zeroValue, err
		}

		embeddedConnectorClient, err := m.clientFactory.New(connectionDetails)
		if err != nil {
			return zeroValue, err
		}

		response := embeddedConnectorClient.Ping(ctx, sessionLogger)
		if !response.IsAlive {
			return zeroValue, ErrMATLABSessionNotAlive
		}

		client = newMATLABSessionClientWithoutCleanup(embeddedConnectorClient)
	default:
		return zeroValue, fmt.Errorf("unknown request type: %T", request)
	}

	return m.sessionStore.Add(client), nil
}
