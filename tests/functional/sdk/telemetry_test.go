// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"io/fs"
	"maps"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/tests/functional/sdk/testbinaries"
	"github.com/matlab/matlab-mcp-server/tests/testutils/otel"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

const (
	telemetryCollectionInterval = 250 * time.Millisecond
	telemetrySessionDuration    = 2 * time.Second
)

type TelemetryTestSuite struct {
	SDKTestSuite

	serverDetails testbinaries.ServerDetails
}

func (s *TelemetryTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildEmptyServer(s.T())
}

func TestTelemetryTestSuite(t *testing.T) {
	suite.Run(t, new(TelemetryTestSuite))
}

func (s *TelemetryTestSuite) TestTelemetry_Enabled_InitializationMetricsRecordedOnce() {
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "Default",
		},
		{
			name: "FlagFalse",
			args: []string{"--disable-telemetry=false"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Arrange
			collector := otel.StartMockCollector(s.T())

			// Act
			session := s.CreateSession(s.serverDetails.BinaryLocation(), telemetryEnvironment(collector.Endpoint()), nil, tc.args...)
			time.Sleep(telemetrySessionDuration)
			s.CleanupSession(session, true)

			// Assert
			s.Equal([]int64{1}, recordedTotals(collector, "server.starts"), "server.starts should be recorded once, by the server only")
			s.Equal([]int64{1}, recordedTotals(collector, "server.client_connections"), "server.client_connections should be recorded once, by the server only")
			s.Positive(watchdogExportCount(collector.Metrics()), "control: the watchdog's exporter is live when telemetry is on, sending empty exports")
		})
	}
}

func (s *TelemetryTestSuite) TestTelemetry_Disabled_EndpointIsSilent() {
	testCases := []struct {
		name        string
		environment []string
		args        []string
	}{
		{
			name: "Flag",
			args: []string{"--disable-telemetry=true"},
		},
		{
			name:        "EnvVar",
			environment: []string{"MW_MCP_SERVER_DISABLE_TELEMETRY=true"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Arrange
			collector := otel.StartMockCollector(s.T())
			environment := append(telemetryEnvironment(collector.Endpoint()), tc.environment...)

			// Act
			session := s.CreateSession(s.serverDetails.BinaryLocation(), environment, nil, tc.args...)
			time.Sleep(telemetrySessionDuration)
			s.CleanupSession(session, true)

			// Assert
			s.Eventually(func() bool {
				return watchdogLogContains(session, "Telemetry disabled by configuration")
			}, 10*time.Second, 100*time.Millisecond, "the watchdog should start with telemetry disabled")

			s.Zero(collector.ConnectionCount(), "neither the server nor the watchdog should connect to the telemetry endpoint")
		})
	}
}

func watchdogLogContains(session *SDKSession, text string) bool {
	logFiles, err := fs.Glob(session.LogFS(), "watchdog-*.log")
	if err != nil {
		return false
	}

	for _, logFile := range logFiles {
		content, err := fs.ReadFile(session.LogFS(), logFile)
		if err != nil {
			return false
		}
		if strings.Contains(string(content), text) {
			return true
		}
	}

	return false
}

// The endpoint and interval go through env vars because the watchdog inherits the environment but, of the telemetry flags, receives only --disable-telemetry.
func telemetryEnvironment(endpoint string) []string {
	environment := []string{}
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if isTelemetryEnvVar(name) {
			continue
		}
		environment = append(environment, entry)
	}

	return append(environment,
		"MW_MCP_SERVER_TELEMETRY_COLLECTOR_ENDPOINT="+endpoint,
		"MW_MCP_SERVER_TELEMETRY_COLLECTION_INTERVAL="+telemetryCollectionInterval.String(),
	)
}

func isTelemetryEnvVar(name string) bool {
	upperName := strings.ToUpper(name)
	return strings.HasPrefix(upperName, "OTEL_") ||
		strings.HasPrefix(upperName, "MW_MCP_SERVER_TELEMETRY_") ||
		upperName == "MW_MCP_SERVER_DISABLE_TELEMETRY"
}

// Cumulative counters resend their total on every export, so recordings are counted per series start time, which is one per process.
func recordedTotals(collector *otel.MockCollector, metricName string) []int64 {
	totalsByStart := map[pcommon.Timestamp]int64{}
	for _, metric := range collector.AllMetricsByName(metricName) {
		if metric.Type() != pmetric.MetricTypeSum {
			continue
		}

		exportTotals := map[pcommon.Timestamp]int64{}
		for _, dataPoint := range metric.Sum().DataPoints().All() {
			exportTotals[dataPoint.StartTimestamp()] += dataPoint.IntValue()
		}

		for start, total := range exportTotals {
			totalsByStart[start] = max(totalsByStart[start], total)
		}
	}

	return slices.Collect(maps.Values(totalsByStart))
}

func watchdogExportCount(exports []pmetric.Metrics) int {
	count := 0
	for _, export := range exports {
		if export.DataPointCount() == 0 {
			count++
		}
	}
	return count
}
