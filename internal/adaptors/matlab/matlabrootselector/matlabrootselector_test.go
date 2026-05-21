// Copyright 2025-2026 The MathWorks, Inc.

package matlabrootselector_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlab/matlabrootselector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlab/matlabrootselector"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	// Act
	selector := matlabrootselector.New(mockConfigFactory, mockMATLABManager)

	// Assert
	assert.NotNil(t, selector)
}

func TestMATLABRootSelector_SelectMATLABRoot_ConfigError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	ctx := t.Context()
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	selector := matlabrootselector.New(mockConfigFactory, mockMATLABManager)

	// Act
	result, err := selector.SelectMATLABRoot(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError, "SelectMATLABRoot should return the error from Config")
	assert.Empty(t, result)
}

func TestMATLABRootSelector_SelectMATLABRoot_HappyPath(t *testing.T) {
	testCases := []struct {
		name         string
		environments []entities.EnvironmentInfo
		expected     string
	}{
		{
			name: "single MATLAB installation",
			environments: []entities.EnvironmentInfo{
				{
					MATLABRoot: filepath.Join("usr", "local", "MATLAB", "R2024b"),
					Version:    "R2024b",
				},
			},
			expected: filepath.Join("usr", "local", "MATLAB", "R2024b"),
		},
		{
			name: "multiple MATLAB installations - returns first one",
			environments: []entities.EnvironmentInfo{
				{
					MATLABRoot: filepath.Join("usr", "local", "MATLAB", "R2023b"),
					Version:    "R2023b",
				},
				{
					MATLABRoot: filepath.Join("usr", "local", "MATLAB", "R2024b"),
					Version:    "R2024b",
				},
				{
					MATLABRoot: filepath.Join("usr", "local", "MATLAB", "R2024a"),
					Version:    "R2024a",
				},
			},
			expected: filepath.Join("usr", "local", "MATLAB", "R2023b"),
		},
		{
			name: "Windows paths",
			environments: []entities.EnvironmentInfo{
				{
					MATLABRoot: filepath.Join("C:", "Program Files", "MATLAB", "R2024b"),
					Version:    "R2024b",
				},
				{
					MATLABRoot: filepath.Join("C:", "Program Files", "MATLAB", "R2024a"),
					Version:    "R2024a",
				},
			},
			expected: filepath.Join("C:", "Program Files", "MATLAB", "R2024b"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockLogger := testutils.NewInspectableLogger()

			mockConfigFactory := &mocks.MockConfigFactory{}
			defer mockConfigFactory.AssertExpectations(t)

			mockConfig := &configmocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			mockMATLABManager := &mocks.MockMATLABManager{}
			defer mockMATLABManager.AssertExpectations(t)

			ctx := t.Context()

			mockConfigFactory.EXPECT().
				Config().
				Return(mockConfig, nil).
				Once()

			mockConfig.EXPECT().
				PreferredLocalMATLABRoot().
				Return("").
				Once()

			mockMATLABManager.EXPECT().
				ListEnvironments(ctx, mockLogger.AsMockArg()).
				Return(tc.environments).
				Once()

			selector := matlabrootselector.New(mockConfigFactory, mockMATLABManager)

			// Act
			result, err := selector.SelectMATLABRoot(ctx, mockLogger)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMATLABRootSelector_SelectMATLABRoot_PreferredMATLABRootSet_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	ctx := t.Context()

	expectedPreferredMATLABRoot := filepath.Join("usr", "local", "MATLAB", "R2024b")

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		PreferredLocalMATLABRoot().
		Return(expectedPreferredMATLABRoot).
		Once()

	selector := matlabrootselector.New(mockConfigFactory, mockMATLABManager)

	// Act
	result, err := selector.SelectMATLABRoot(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPreferredMATLABRoot, result)
}

func TestMATLABRootSelector_SelectMATLABRoot_ListEnvironmentsEmpty(t *testing.T) {
	testCases := []struct {
		name         string
		environments []entities.EnvironmentInfo
	}{
		{
			name:         "empty environments list",
			environments: []entities.EnvironmentInfo{},
		},
		{
			name:         "nil environments list",
			environments: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockLogger := testutils.NewInspectableLogger()

			mockConfigFactory := &mocks.MockConfigFactory{}
			defer mockConfigFactory.AssertExpectations(t)

			mockConfig := &configmocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			mockMATLABManager := &mocks.MockMATLABManager{}
			defer mockMATLABManager.AssertExpectations(t)

			ctx := t.Context()

			mockConfigFactory.EXPECT().
				Config().
				Return(mockConfig, nil).
				Once()

			mockConfig.EXPECT().
				PreferredLocalMATLABRoot().
				Return("").
				Once()

			mockMATLABManager.EXPECT().
				ListEnvironments(ctx, mockLogger.AsMockArg()).
				Return(tc.environments).
				Once()

			selector := matlabrootselector.New(mockConfigFactory, mockMATLABManager)

			// Act
			result, err := selector.SelectMATLABRoot(ctx, mockLogger)

			// Assert
			require.Error(t, err)
			assert.Empty(t, result)
		})
	}
}
