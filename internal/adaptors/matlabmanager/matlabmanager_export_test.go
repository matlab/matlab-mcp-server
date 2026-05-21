// Copyright 2026 The MathWorks, Inc.

package matlabmanager

import (
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

func (m *MATLABManager) SetMATLABSessionConnectionRetryInterval(matlabSessionConnectionRetryInterval time.Duration) {
	m.matlabSessionConnectionRetryInterval = matlabSessionConnectionRetryInterval
}

func NewMATLABSessionClientWithoutCleanup(matlabSessionClient entities.MATLABSessionClient) *matlabSessionClientWithoutCleanup {
	return newMATLABSessionClientWithoutCleanup(matlabSessionClient)
}

func NewMATLABSessionClientWithCleanup(matlabSessionClient entities.MATLABSessionClient, sessionCleanup func() error) *matlabSessionClientWithCleanup {
	return newMATLABSessionClientWithCleanup(matlabSessionClient, sessionCleanup)
}
