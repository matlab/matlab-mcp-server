// Copyright 2026 The MathWorks, Inc.

//go:build darwin

package os_test

import (
	"testing"

	osadaptor "github.com/matlab/matlab-mcp-server/internal/adaptors/os"
	osmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/os"
	osfacademocks "github.com/matlab/matlab-mcp-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOS_Version_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &osmocks.MockVersionOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockRegistryLayer := &osmocks.MockRegistryLayer{}
	defer mockRegistryLayer.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	expectedVersion := "macOS 15.3.1"

	mockOSLayer.EXPECT().
		Command("sw_vers", []string{"-productVersion"}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		Output().
		Return([]byte("15.3.1\n"), nil).
		Once()

	osInstance := osadaptor.New(mockOSLayer, mockRegistryLayer)

	// Act
	version, err := osInstance.Version()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedVersion, version)
}

func TestOS_Version_CommandOutputError(t *testing.T) {
	// Arrange
	mockOSLayer := &osmocks.MockVersionOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockRegistryLayer := &osmocks.MockRegistryLayer{}
	defer mockRegistryLayer.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	mockOSLayer.EXPECT().
		Command("sw_vers", []string{"-productVersion"}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		Output().
		Return(nil, assert.AnError).
		Once()

	osInstance := osadaptor.New(mockOSLayer, mockRegistryLayer)

	// Act
	version, err := osInstance.Version()

	// Assert
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, version)
}
