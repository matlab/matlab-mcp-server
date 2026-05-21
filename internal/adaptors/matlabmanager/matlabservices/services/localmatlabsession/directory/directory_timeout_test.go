// Copyright 2026 The MathWorks, Inc.

package directory_test

import (
	"os"
	"path/filepath"
	"testing"
	"testing/synctest"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectory_GetEmbeddedConnectorDetails_RespectsConfiguredTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockOSLayer := &mocks.MockOSLayer{}
		defer mockOSLayer.AssertExpectations(t)

		mockConfig := &mocks.MockConfig{}
		defer mockConfig.AssertExpectations(t)

		mockLogger := testutils.NewInspectableLogger()

		configuredTimeout := 10 * time.Millisecond
		configuredRetry := 5 * time.Millisecond

		mockConfig.EXPECT().
			EmbeddedConnectorDetailsTimeout().
			Return(10 * time.Millisecond).
			Once()

		dir := directory.NewDirectory(mockLogger, filepath.Join("tmp", "matlab-session-12345"), mockOSLayer, mockConfig)
		dir.SetEmbeddedConnectorDetailsRetry(configuredRetry)

		startupErrorFile := dir.StartupErrorFile()
		securePortFile := dir.SecurePortFile()

		mockOSLayer.EXPECT().
			ReadFile(startupErrorFile).
			Return(nil, os.ErrNotExist)

		mockOSLayer.EXPECT().
			Stat(securePortFile).
			Return(nil, os.ErrNotExist)

		// Act
		start := time.Now()
		_, _, err := dir.GetEmbeddedConnectorDetails()
		elapsed := time.Since(start)

		// Assert
		require.Error(t, err)
		assert.Equal(t, configuredTimeout, elapsed)
	})
}

func TestDirectory_NewDirectory_UsesConfiguredEmbeddedConnectorDetailsTimeout(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedTimeout := 10 * time.Minute

	mockConfig.EXPECT().
		EmbeddedConnectorDetailsTimeout().
		Return(10 * time.Minute).
		Once()

	dir := directory.NewDirectory(mockLogger, filepath.Join("tmp", "matlab-session-12345"), mockOSLayer, mockConfig)

	// Act
	timeout := dir.GetEmbeddedConnectorDetailsTimeout()

	// Assert
	assert.Equal(t, expectedTimeout, timeout)
}
