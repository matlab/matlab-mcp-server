// Copyright 2026 The MathWorks, Inc.

package requirements

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const releaseEndpointEnvironmentVariable = "TELEMETRY_ENDPOINT_VALUE"

var errNoRemoteReleaseTag = errors.New("remote does not contain a semantic version tag")

type ReleaseConfiguration struct {
	ExpectedVersion           string
	ExpectedTelemetryEndpoint string
}

func Release(t TestingT) ReleaseConfiguration {
	t.Helper()

	endpoint, found := os.LookupEnv(releaseEndpointEnvironmentVariable)
	if !found {
		t.Skip("release test requirements are not available")
		return ReleaseConfiguration{}
	}

	version, err := latestReleaseVersion(t.Context())
	if err != nil {
		t.Fatal("release version requirement is not available")
		return ReleaseConfiguration{}
	}

	return ReleaseConfiguration{
		ExpectedVersion:           version,
		ExpectedTelemetryEndpoint: endpoint,
	}
}

func latestReleaseVersion(ctx context.Context) (string, error) {
	command := exec.CommandContext(
		ctx,
		"git",
		"ls-remote",
		"--tags",
		"--refs",
		"--sort=-version:refname",
		"origin",
		"refs/tags/v[0-9]*.[0-9]*.[0-9]*",
	)
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("list remote release tags: %w", err)
	}

	fields := strings.Fields(string(output))
	if len(fields) < 2 {
		return "", errNoRemoteReleaseTag
	}

	return strings.TrimPrefix(fields[1], "refs/tags/"), nil
}
