// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type matlabSessionClientWithoutCleanup struct {
	entities.MATLABSessionClient
}

type matlabSessionClientWithCleanup struct {
	entities.MATLABSessionClient
	sessionCleanup func() error
}

func newMATLABSessionClientWithoutCleanup(matlabSessionClient entities.MATLABSessionClient) *matlabSessionClientWithoutCleanup {
	return &matlabSessionClientWithoutCleanup{
		MATLABSessionClient: matlabSessionClient,
	}
}

func (c *matlabSessionClientWithoutCleanup) StopSession(ctx context.Context, sessionLogger entities.Logger) error {
	sessionLogger.Debug("Skipping session stop for externally managed MATLAB session")
	return nil
}

func newMATLABSessionClientWithCleanup(matlabSessionClient entities.MATLABSessionClient, sessionCleanup func() error) *matlabSessionClientWithCleanup {
	return &matlabSessionClientWithCleanup{
		MATLABSessionClient: matlabSessionClient,
		sessionCleanup:      sessionCleanup,
	}
}

func (c *matlabSessionClientWithCleanup) StopSession(ctx context.Context, sessionLogger entities.Logger) error {
	_, err := c.Eval(ctx, sessionLogger, entities.EvalRequest{Code: "exit()"})
	if err != nil {
		return err
	}

	return c.sessionCleanup()
}
