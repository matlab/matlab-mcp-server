// Copyright 2026 The MathWorks, Inc.

package release_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/tests/internal/requirements"
	"github.com/matlab/matlab-mcp-server/tests/testutils/mcpserver"
	"github.com/stretchr/testify/suite"
)

const (
	releaseArtifactExitTimeout = 5 * time.Second
	redactedEndpoint           = "[redacted]"
)

type ReleaseArtifactTestSuite struct {
	suite.Suite

	expectedEndpoint string
	expectedAddress  string
	endpointPattern  *regexp.Regexp
	expectedVersion  string
	serverPath       string
}

func TestReleaseArtifactTestSuite(t *testing.T) {
	suite.Run(t, new(ReleaseArtifactTestSuite))
}

func (s *ReleaseArtifactTestSuite) SetupSuite() {
	release := requirements.Release(s.T())
	s.expectedEndpoint = release.ExpectedTelemetryEndpoint
	s.setExpectedAddress(release.ExpectedTelemetryEndpoint)
	s.expectedVersion = release.ExpectedVersion

	serverPath, err := mcpserver.NewLocator().GetPath()
	s.Require().NoError(err, "release artifact could not be located")
	s.serverPath = serverPath
}

func (s *ReleaseArtifactTestSuite) TestArtifactVersionMatchesRelease() {
	t := s.T()
	var output strings.Builder

	// Step 1: request the version from the release artifact.
	err := s.runArtifact(t.Context(), "", &output, "--version", "--disable-telemetry")
	s.Require().NoError(err, "release artifact did not complete")

	// Step 2: confirm that the artifact reports the latest release version.
	version, err := mcpserver.ParseVersion(output.String())
	s.Require().NoError(err, "release artifact reported an invalid version")
	s.Equal(s.expectedVersion, version, "release artifact reported an unexpected version")
}

func (s *ReleaseArtifactTestSuite) TestArtifactContainsConfiguredTelemetryEndpoint() {
	// Step 1: read the release artifact.
	artifact, err := os.ReadFile(s.serverPath)
	s.Require().NoError(err, "release artifact could not be read")

	// Step 2: confirm that the artifact contains the full configured endpoint, including the path.
	s.True(bytes.Contains(artifact, []byte(s.expectedEndpoint)),
		"release artifact does not contain the configured telemetry endpoint")
}

func (s *ReleaseArtifactTestSuite) TestArtifactUsesConfiguredTelemetryEndpoint() {
	t := s.T()

	// Step 1: start a proxy that records the telemetry destination.
	requestAddresses := make(chan string, 1)
	proxy := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		select {
		case requestAddresses <- request.Host:
		default:
		}
		response.WriteHeader(http.StatusBadGateway)
	}))
	defer proxy.Close()

	// Step 2: run the release artifact with telemetry routed through the proxy.
	err := s.runArtifact(t.Context(), proxy.URL, io.Discard, "--version")
	s.Require().NoError(err, "release artifact did not complete")

	// Step 3: confirm that the artifact used the configured telemetry destination.
	select {
	case actualAddress := <-requestAddresses:
		s.True(strings.EqualFold(actualAddress, s.expectedAddress),
			"release artifact used an unexpected telemetry address")
	default:
		t.Fatal("release artifact did not contact the telemetry proxy")
	}
}

func (s *ReleaseArtifactTestSuite) TestArtifactDoesNotSendTelemetryWhenDisabled() {
	t := s.T()

	// Step 1: start a proxy that records any telemetry request.
	var contacted atomic.Bool
	proxy := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		contacted.Store(true)
	}))
	defer proxy.Close()

	// Step 2: run the release artifact with telemetry disabled.
	err := s.runArtifact(t.Context(), proxy.URL, io.Discard, "--version", "--disable-telemetry")
	s.Require().NoError(err, "release artifact did not complete")

	// Step 3: confirm that the compiled-in endpoint generated no network traffic.
	s.False(contacted.Load(), "release artifact contacted telemetry after telemetry was disabled")
}

func (s *ReleaseArtifactTestSuite) TestArtifactExitsWhenTelemetryProxyDoesNotRespond() {
	t := s.T()

	// Step 1: start a proxy that accepts the telemetry request without responding.
	var contacted atomic.Bool
	allowResponse := make(chan struct{})
	proxy := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		contacted.Store(true)
		<-allowResponse
		response.WriteHeader(http.StatusBadGateway)
	}))
	defer func() {
		close(allowResponse)
		proxy.Close()
	}()

	// Step 2: run the release artifact with a bounded context.
	ctx, cancel := context.WithTimeout(t.Context(), releaseArtifactExitTimeout)
	defer cancel()
	err := s.runArtifact(ctx, proxy.URL, io.Discard, "--version")

	// Step 3: confirm that an unresponsive telemetry destination did not delay exit.
	s.Require().NoError(err, "release artifact did not exit while the telemetry proxy was unresponsive")
	s.True(contacted.Load(), "release artifact did not contact the telemetry proxy")
}

func (s *ReleaseArtifactTestSuite) setExpectedAddress(endpoint string) {
	parsed, err := url.Parse(endpoint)
	s.Require().True(err == nil, "release telemetry endpoint configuration is invalid") //nolint:testifylint // NoError can disclose the endpoint.
	s.Require().NotNil(parsed, "release telemetry endpoint configuration is invalid")
	s.Require().NotEmpty(parsed.Hostname(), "release telemetry endpoint configuration is invalid")

	defaultPorts := map[string]string{
		"http":  "80",
		"https": "443",
	}
	defaultPort, supportedScheme := defaultPorts[strings.ToLower(parsed.Scheme)]
	s.Require().True(supportedScheme, "release telemetry endpoint configuration is invalid")
	port := parsed.Port()
	switch port {
	case "":
		port = defaultPort
	}

	s.expectedAddress = net.JoinHostPort(parsed.Hostname(), port)
	s.endpointPattern = regexp.MustCompile(
		"(?i)" + regexp.QuoteMeta(s.expectedEndpoint) + "|" + regexp.QuoteMeta(parsed.Hostname()),
	)
}

func (s *ReleaseArtifactTestSuite) runArtifact(
	ctx context.Context,
	proxyURL string,
	stdout io.Writer,
	arguments ...string,
) error {
	arguments = append(arguments, "--log-level=debug", "--duplicate-logs-to-stderr")
	command := exec.CommandContext(ctx, s.serverPath, arguments...) //nolint:gosec // Locator validates the artifact path.
	command.Env = s.releaseEndpointEnvironment(proxyURL)
	command.Stdout = stdout
	var stderr bytes.Buffer
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("%w\nstderr:\n%s", err, s.redact(stderr.String()))
	}
	return nil
}

func (s *ReleaseArtifactTestSuite) redact(text string) string {
	return s.endpointPattern.ReplaceAllLiteralString(text, redactedEndpoint)
}

func (s *ReleaseArtifactTestSuite) releaseEndpointEnvironment(proxyURL string) []string {
	s.T().Helper()

	environment := slices.DeleteFunc(os.Environ(), s.isExcludedReleaseEnvironmentEntry)

	return append(environment,
		"HTTP_PROXY="+proxyURL,
		"HTTPS_PROXY="+proxyURL,
		"NO_PROXY=",
	)
}

func (*ReleaseArtifactTestSuite) isExcludedReleaseEnvironmentEntry(entry string) bool {
	name, _, _ := strings.Cut(entry, "=")
	upperName := strings.ToUpper(name)

	return slices.Contains(
		[]string{
			"TELEMETRY_COLLECTOR_ENDPOINT",
			"TELEMETRY_ENDPOINT_VALUE",
			"HTTP_PROXY",
			"HTTPS_PROXY",
			"NO_PROXY",
			"ALL_PROXY",
		},
		upperName,
	) ||
		strings.HasPrefix(upperName, "MW_MCP_SERVER_") ||
		strings.HasPrefix(upperName, "OTEL_EXPORTER_OTLP_")
}
