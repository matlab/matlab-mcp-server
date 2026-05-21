// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/stretchr/testify/suite"
)

// ServerWithLoggingTestSuite tests SDK logging functionalities.
type ServerWithLoggingTestSuite struct {
	SDKTestSuite

	serverDetails testbinaries.ServerWithLoggingDetails
}

// SetupSuite runs once before all tests in a suite
func (s *ServerWithLoggingTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildServerWithLogging(s.T())
}

func TestServerWithLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWithLoggingTestSuite))
}

func (s *ServerWithLoggingTestSuite) TestSDK_Logging_DependenciesAndToolsProviderLogToFile() {
	// Arrange
	logFolder, err := os.MkdirTemp("", "server_session") // Can't use s.T().Tempdir() because too long for socket path
	s.Require().NoError(err)
	defer func() {
		if err := os.RemoveAll(logFolder); err != nil {
			s.T().Logf("Failed to remove log folder (may be locked on Windows): %v", err)
		}
	}()

	// This suite intentionally verifies logging behavior and may emit ERROR logs.
	session := s.CreateSession(s.serverDetails.BinaryLocation(), nil, nil, "--log-folder="+logFolder)
	defer s.CleanupSession(session, false)

	// Act
	_, err = session.CallTool(s.T().Context(), s.serverDetails.ToolThatLogsName(), map[string]any{"name": "World"})
	s.Require().NoError(err, "should call tool successfully")

	// Assert
	ctx, cancel := context.WithTimeout(s.T().Context(), 2*time.Second) // Timeout for the logs to write to disk
	defer cancel()

	_, err = retry.Retry(ctx, func() (struct{}, bool, error) {
		logContent, err := session.ReadServerLogs()
		if err != nil {
			return struct{}{}, false, err
		}

		foundDependenciesProviderLog := strings.Contains(logContent, "Creating Dependencies")
		foundToolsProviderLog := strings.Contains(logContent, "Creating Tools")

		return struct{}{}, foundDependenciesProviderLog && foundToolsProviderLog, nil
	}, retry.NewLinearRetryStrategy(200*time.Millisecond))

	s.Require().NoError(err)
}

func (s *ServerWithLoggingTestSuite) TestSDK_Logging_ToolHandlerLogsToFile() {
	// Arrange
	logFolder, err := os.MkdirTemp("", "server_session") // Can't use s.T().Tempdir() because too long for socket path
	s.Require().NoError(err)
	defer func() {
		if err := os.RemoveAll(logFolder); err != nil {
			s.T().Logf("Failed to remove log folder (may be locked on Windows): %v", err)
		}
	}()

	name := "World"

	// This suite intentionally verifies logging behavior and may emit ERROR logs.
	session := s.CreateSession(s.serverDetails.BinaryLocation(), nil, nil, "--log-folder="+logFolder)
	defer s.CleanupSession(session, false)

	// Act
	_, err = session.CallTool(s.T().Context(), s.serverDetails.ToolThatLogsName(), map[string]any{"name": name})
	s.Require().NoError(err, "should call unstructured tool successfully")

	_, err = session.CallTool(s.T().Context(), s.serverDetails.StructuredToolThatLogsName(), map[string]any{"name": name})
	s.Require().NoError(err, "should call structured tool successfully")

	// Assert
	ctx, cancel := context.WithTimeout(s.T().Context(), 2*time.Second) // Timeout for the logs to write to disk
	defer cancel()

	_, err = retry.Retry(ctx, func() (struct{}, bool, error) {
		logContent, err := session.ReadServerLogs()
		if err != nil {
			return struct{}{}, false, err
		}

		foundUnstructuredLogEntry := strings.Contains(logContent, "Logging from unstructured tool: "+name)
		foundStructuredLogEntry := strings.Contains(logContent, "Logging from structured tool: "+name)

		return struct{}{}, foundUnstructuredLogEntry && foundStructuredLogEntry, nil
	}, retry.NewLinearRetryStrategy(200*time.Millisecond))

	s.Require().NoError(err)
}
