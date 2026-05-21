// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

const defaultMATLABSessionConnectionRetryInterval = 100 * time.Millisecond

func (m *MATLABManager) GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error) {
	config, messagesErr := m.configFactory.Config()
	if messagesErr != nil {
		return nil, messagesErr
	}

	client, err := m.sessionStore.Get(sessionID)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, config.MATLABSessionConnectionTimeout())
	defer cancel()

	_, err = retry.Retry(pingCtx, func() (struct{}, bool, error) {
		pingResponse := client.Ping(pingCtx, sessionLogger)
		return struct{}{}, pingResponse.IsAlive, nil
	}, retry.NewLinearRetryStrategy(m.matlabSessionConnectionRetryInterval))

	if err != nil {
		return nil, fmt.Errorf("MATLAB session %v is not alive", sessionID)
	}

	return client, nil
}
