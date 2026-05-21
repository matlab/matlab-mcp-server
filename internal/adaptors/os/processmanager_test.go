// Copyright 2025-2026 The MathWorks, Inc.

package os_test

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"testing/synctest"
	"time"

	osadaptor "github.com/matlab/matlab-mcp-core-server/internal/adaptors/os"
	osmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/os"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
)

func TestNewProcessManager_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &osmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		GOOS().
		Return("linux").
		Once()

	// Act
	pm := osadaptor.NewProcessManager(mockOSLayer)

	// Assert
	assert.NotNil(t, pm, "ProcessManager instance should not be nil")
}

func TestProcessManager_FindProcess_Unix_HappyPath(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		t.Run(goos, func(t *testing.T) {
			// Arrange
			mockOSLayer := &osmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockProcess := &osfacademocks.MockProcess{}
			defer mockProcess.AssertExpectations(t)

			processPid := 1234

			mockOSLayer.EXPECT().
				GOOS().
				Return(goos).
				Once()

			mockOSLayer.EXPECT().
				FindProcess(processPid).
				Return(mockProcess, nil).
				Once()

			mockProcess.EXPECT().
				Signal(syscall.Signal(0)).
				Return(nil).
				Once()

			pm := osadaptor.NewProcessManager(mockOSLayer)

			// Act
			result := pm.FindProcess(processPid)

			// Assert
			assert.Equal(t, mockProcess, result, "Should return the found process")
		})
	}
}

func TestProcessManager_FindProcess_Windows_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &osmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcess := &osfacademocks.MockProcess{}
	defer mockProcess.AssertExpectations(t)

	processPid := 1234

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		FindProcess(processPid).
		Return(mockProcess, nil).
		Once()

	pm := osadaptor.NewProcessManager(mockOSLayer)

	// Act
	result := pm.FindProcess(processPid)

	// Assert
	assert.Equal(t, mockProcess, result, "Should return the found process")
}

func TestProcessManager_FindProcess_OSLayerFindProcessError(t *testing.T) {
	// Arrange
	mockOSLayer := &osmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	processPid := 1234
	expectedError := assert.AnError

	mockOSLayer.EXPECT().
		GOOS().
		Return("linux").
		Once()

	mockOSLayer.EXPECT().
		FindProcess(processPid).
		Return(nil, expectedError).
		Once()

	pm := osadaptor.NewProcessManager(mockOSLayer)

	// Act
	result := pm.FindProcess(processPid)

	// Assert
	assert.Nil(t, result, "Should return nil when OSLayer.FindProcess returns error")
}

func TestProcessManager_FindProcess_Unix_ProcessSignalError(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		t.Run(goos, func(t *testing.T) {
			// Arrange
			mockOSLayer := &osmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockProcess := &osfacademocks.MockProcess{}
			defer mockProcess.AssertExpectations(t)

			processPid := 1234
			signalError := assert.AnError

			mockOSLayer.EXPECT().
				GOOS().
				Return(goos).
				Once()

			mockOSLayer.EXPECT().
				FindProcess(processPid).
				Return(mockProcess, nil).
				Once()

			mockProcess.EXPECT().
				Signal(syscall.Signal(0)).
				Return(signalError).
				Once()

			pm := osadaptor.NewProcessManager(mockOSLayer)

			// Act
			result := pm.FindProcess(processPid)

			// Assert
			assert.Nil(t, result, "Should return nil when process signal check fails on Unix")
		})
	}
}

func TestProcessManager_WaitForProcessToComplete_Unix_HappyPath(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		t.Run(goos, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				// Arrange
				mockOSLayer := &osmocks.MockOSLayer{}
				defer mockOSLayer.AssertExpectations(t)

				mockProcess := &osfacademocks.MockProcess{}
				defer mockProcess.AssertExpectations(t)

				processPid := 1234

				mockOSLayer.EXPECT().
					GOOS().
					Return(goos).
					Once()

				mockOSLayer.EXPECT().
					FindProcess(processPid).
					Return(mockProcess, nil).
					Times(3)

				mockProcess.EXPECT().
					Signal(syscall.Signal(0)).
					Return(nil).
					Twice()

				mockProcess.EXPECT().
					Signal(syscall.Signal(0)).
					Return(fmt.Errorf("process not found")).
					Once()

				pm := osadaptor.NewProcessManager(mockOSLayer)

				// Act & Assert: mock expectations verify it was called 3 times with 2 retries
				pm.WaitForProcessToComplete(processPid)
			})
		})
	}
}

func TestProcessManager_WaitForProcessToComplete_Windows_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &osmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcess := &osfacademocks.MockProcess{}
	defer mockProcess.AssertExpectations(t)

	processPid := 1234

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		FindProcess(processPid).
		Return(mockProcess, nil).
		Once()

	processRunTime := 10 * time.Millisecond
	mockProcess.EXPECT().
		Wait().
		RunAndReturn(func() (*os.ProcessState, error) {
			<-time.After(processRunTime)
			return nil, nil
		}).
		Once()

	pm := osadaptor.NewProcessManager(mockOSLayer)

	// Act
	startTime := time.Now()
	pm.WaitForProcessToComplete(processPid)
	duration := time.Since(startTime)

	// Assert
	assert.GreaterOrEqual(t, duration, processRunTime, "Should wait for process to complete on Windows")
}

func TestProcessManager_InterruptSignalChan_ReturnsChannel(t *testing.T) {
	// Arrange
	mockOSLayer := &osmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		GOOS().
		Return("linux").
		Once()

	pm := osadaptor.NewProcessManager(mockOSLayer)

	// Act
	signalChan := pm.InterruptSignalChan()

	// Assert
	assert.NotNil(t, signalChan, "InterruptSignalChan should return a non-nil channel")
}
