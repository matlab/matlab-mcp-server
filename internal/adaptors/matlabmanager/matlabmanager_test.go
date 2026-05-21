// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager"
	"github.com/stretchr/testify/assert"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	// Act
	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector)

	// Assert
	assert.NotNil(t, manager, "MATLABManager should not be nil")
}
