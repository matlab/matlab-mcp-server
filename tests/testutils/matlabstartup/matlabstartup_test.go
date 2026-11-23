// Copyright 2026 The MathWorks, Inc.

package matlabstartup_test

import (
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-server/tests/testutils/matlabstartup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeout_DefaultsToTheServerStartupTimeout(t *testing.T) {
	t.Setenv(matlabstartup.EnvKey, "")

	timeout, err := matlabstartup.Timeout()

	require.NoError(t, err)
	assert.Equal(t, defaultparameters.EmbeddedConnectorDetailsTimeout().GetTypedDefaultValue(), timeout,
		"a test should give MATLAB the time that the server gives a MATLAB it launches")
}

func TestTimeout_UsesTheEnvironmentValue(t *testing.T) {
	t.Setenv(matlabstartup.EnvKey, "15m")

	timeout, err := matlabstartup.Timeout()

	require.NoError(t, err)
	assert.Equal(t, 15*time.Minute, timeout)
}

func TestTimeout_RejectsAValueThatIsNotADuration(t *testing.T) {
	t.Setenv(matlabstartup.EnvKey, "15")

	_, err := matlabstartup.Timeout()

	require.Error(t, err)
	assert.Contains(t, err.Error(), matlabstartup.EnvKey+" must be a duration")
}

func TestTimeout_RejectsAValueThatIsNotPositive(t *testing.T) {
	for _, value := range []string{"0s", "-1m"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv(matlabstartup.EnvKey, value)

			_, err := matlabstartup.Timeout()

			require.Error(t, err)
			assert.Contains(t, err.Error(), matlabstartup.EnvKey+" must be positive")
		})
	}
}

func TestClientTimeout_OutlastsTheServerStartupTimeoutAndCleanup(t *testing.T) {
	t.Setenv(matlabstartup.EnvKey, "15m")

	timeout, err := matlabstartup.ClientTimeout()

	require.NoError(t, err)
	assert.Equal(t, 18*time.Minute, timeout,
		"the server stops MATLAB for up to 2m before it reports a MATLAB that does not start, so the MCP client must wait longer")
}

func TestClientTimeout_RejectsAnInvalidStartupTimeout(t *testing.T) {
	t.Setenv(matlabstartup.EnvKey, "15")

	_, err := matlabstartup.ClientTimeout()

	require.Error(t, err)
	assert.Contains(t, err.Error(), matlabstartup.EnvKey)
}

func TestEnv_GivesTheServerTheSameStartupTimeout(t *testing.T) {
	serverEnvKey := defaultparameters.EmbeddedConnectorDetailsTimeout().GetEnvVarName()
	t.Setenv(matlabstartup.EnvKey, "15m")

	env, err := matlabstartup.Env([]string{"PATH=/usr/bin"})

	require.NoError(t, err)
	assert.Equal(t, []string{"PATH=/usr/bin", serverEnvKey + "=15m0s"}, env)
}

func TestEnv_ReplacesAServerStartupTimeoutThatIsAlreadySet(t *testing.T) {
	serverEnvKey := defaultparameters.EmbeddedConnectorDetailsTimeout().GetEnvVarName()
	t.Setenv(matlabstartup.EnvKey, "15m")

	env, err := matlabstartup.Env([]string{serverEnvKey + "=1m"})

	require.NoError(t, err)
	assert.Equal(t, []string{serverEnvKey + "=15m0s"}, env,
		"one setting must control both the server and the tests")
}

func TestEnv_RejectsAnInvalidStartupTimeout(t *testing.T) {
	t.Setenv(matlabstartup.EnvKey, "15")

	_, err := matlabstartup.Env(nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), matlabstartup.EnvKey)
}
