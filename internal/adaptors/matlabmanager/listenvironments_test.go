// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMATLABManager_ListEnvironments_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABManager := &mocks.MockMATLABServices{}
	defer mockMATLABManager.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	expectedMatlabInfos := []datatypes.MatlabInfo{{
		Location: filepath.Join("path", "to", "matlab", "R2023a"),
		Version: datatypes.MatlabVersionInfo{
			ReleaseFamily: "R2023a",
			ReleasePhase:  "release",
			UpdateLevel:   0,
		},
	}, {
		Location: filepath.Join("path", "to", "matlab", "R2022b"),
		Version: datatypes.MatlabVersionInfo{
			ReleaseFamily: "R2022b",
			ReleasePhase:  "release",
			UpdateLevel:   1,
		},
	},
	}

	mockResponse := datatypes.ListMatlabInfo{
		MatlabInfo: expectedMatlabInfos,
	}
	mockMATLABManager.EXPECT().
		ListDiscoveredMatlabInfo(mockLogger.AsMockArg()).
		Return(mockResponse).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABManager, mockSessionStore, mockClientFactory, mockSessionSelector)
	ctx := t.Context()

	// Act
	result := manager.ListEnvironments(ctx, mockLogger)

	// Assert
	require.Len(t, result, 2)

	// Verify the outputs match the mock data
	for i := range expectedMatlabInfos {
		assert.Equal(t, expectedMatlabInfos[i].Location, result[i].MATLABRoot, "Output MATLAB root does not match input dummy data")
		assert.Equal(t, expectedMatlabInfos[i].Version.ReleaseFamily, result[i].Version, "Output MATLAB version does not match input dummy data")
	}
}

func TestMATLABManager_ListEnvironments_EmptyList(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABManager := &mocks.MockMATLABServices{}
	defer mockMATLABManager.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockResponse := datatypes.ListMatlabInfo{
		MatlabInfo: []datatypes.MatlabInfo{},
	}
	mockMATLABManager.EXPECT().
		ListDiscoveredMatlabInfo(mockLogger.AsMockArg()).
		Return(mockResponse).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABManager, mockSessionStore, mockClientFactory, mockSessionSelector)
	ctx := t.Context()

	// Act
	result := manager.ListEnvironments(ctx, mockLogger)

	// Assert
	assert.NotNil(t, result)
	assert.Empty(t, result)
}
