// Copyright 2025-2026 The MathWorks, Inc.

package process_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/watchdog/process"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	processmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/watchdog/process"
	osfacademocks "github.com/matlab/matlab-mcp-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewProcess_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfig := &processmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	expectedProgramPath := filepath.Join("path", "to", "program")
	expectedBaseDir := filepath.Join("tmp", "base", "dir")
	expectedServerID := "server-id"
	expectedLogLevel := entities.LogLevelInfo
	expectedDisableTelemetry := true

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedServerID).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(expectedLogLevel).
		Once()

	mockConfig.EXPECT().
		DisableTelemetry().
		Return(expectedDisableTelemetry).
		Once()

	mockOSLayer.EXPECT().
		Executable().
		Return(expectedProgramPath, nil).
		Once()

	mockOSLayer.EXPECT().
		Command(expectedProgramPath, []string{
			"--" + defaultparameters.WatchdogMode().GetFlagName(),
			"--" + defaultparameters.BaseDir().GetFlagName(), expectedBaseDir,
			"--" + defaultparameters.ServerInstanceID().GetFlagName(), expectedServerID,
			"--" + defaultparameters.LogLevel().GetFlagName(), string(expectedLogLevel),
			"--" + defaultparameters.DisableTelemetry().GetFlagName() + "=true",
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		SetSysProcAttr(mock.Anything).
		Once()

	mockCmd.EXPECT().
		Start().
		Return(nil).
		Once()

	// Act
	err := process.NewProcess(mockOSLayer, mockLogger, mockDirectory, mockConfig)

	// Assert
	require.NoError(t, err)
}

func TestNewProcess_ExecutableError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfig := &processmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	expectedError := messages.New_StartupErrors_FailedToGetExecutablePath_Error()

	mockOSLayer.EXPECT().
		Executable().
		Return("", assert.AnError).
		Once()

	// Act
	err := process.NewProcess(mockOSLayer, mockLogger, mockDirectory, mockConfig)

	// Assert
	require.Equal(t, expectedError, err)
}

func TestNewProcess_CommandStartError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfig := &processmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	expectedProgramPath := filepath.Join("path", "to", "program")
	expectedBaseDir := filepath.Join("tmp", "base", "dir")
	expectedServerID := "server-id"
	expectedLogLevel := entities.LogLevelInfo
	expectedDisableTelemetry := true
	expectedError := messages.New_StartupErrors_FailedToStartWatchdogProcess_Error()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedServerID).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(expectedLogLevel).
		Once()

	mockConfig.EXPECT().
		DisableTelemetry().
		Return(expectedDisableTelemetry).
		Once()

	mockOSLayer.EXPECT().
		Executable().
		Return(expectedProgramPath, nil).
		Once()

	mockOSLayer.EXPECT().
		Command(expectedProgramPath, []string{
			"--" + defaultparameters.WatchdogMode().GetFlagName(),
			"--" + defaultparameters.BaseDir().GetFlagName(), expectedBaseDir,
			"--" + defaultparameters.ServerInstanceID().GetFlagName(), expectedServerID,
			"--" + defaultparameters.LogLevel().GetFlagName(), string(expectedLogLevel),
			"--" + defaultparameters.DisableTelemetry().GetFlagName() + "=true",
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		SetSysProcAttr(mock.Anything).
		Once()

	mockCmd.EXPECT().
		Start().
		Return(assert.AnError).
		Once()

	// Act
	err := process.NewProcess(mockOSLayer, mockLogger, mockDirectory, mockConfig)

	// Assert
	require.Equal(t, expectedError, err)
}
