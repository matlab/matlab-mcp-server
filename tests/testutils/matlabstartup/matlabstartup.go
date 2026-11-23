// Copyright 2026 The MathWorks, Inc.

package matlabstartup

import (
	"fmt"
	"os"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-server/tests/testutils/pathcontrol"
)

const EnvKey = "MCP_MATLAB_STARTUP_TIMEOUT"

const (
	serverTimeToLaunchMATLAB = time.Minute
	serverTimeToStopMATLAB   = 2 * time.Minute
	serverReportMargin       = serverTimeToLaunchMATLAB + serverTimeToStopMATLAB
)

func Timeout() (time.Duration, error) {
	value := os.Getenv(EnvKey)
	if value == "" {
		return defaultparameters.EmbeddedConnectorDetailsTimeout().GetTypedDefaultValue(), nil
	}

	timeout, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%v must be a duration such as 15m, got %q", EnvKey, value)
	}

	if timeout <= 0 {
		return 0, fmt.Errorf("%v must be positive, got %q", EnvKey, value)
	}

	return timeout, nil
}

func ClientTimeout() (time.Duration, error) {
	timeout, err := Timeout()
	if err != nil {
		return 0, err
	}

	return timeout + serverReportMargin, nil
}

func Env(env []string) ([]string, error) {
	timeout, err := Timeout()
	if err != nil {
		return nil, err
	}

	serverEnvKey := defaultparameters.EmbeddedConnectorDetailsTimeout().GetEnvVarName()
	return pathcontrol.UpdateEnvEntry(env, serverEnvKey, timeout.String()), nil
}
