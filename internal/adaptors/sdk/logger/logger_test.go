// Copyright 2026 The MathWorks, Inc.

package logger_test

import (
	"testing"

	loggeradaptor "github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/logger"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	// Act
	factory := loggeradaptor.NewFactory()

	// Assert
	require.NotNil(t, factory)
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	// Act
	adaptor := loggeradaptor.NewFactory().New(mockLogger)

	// Assert
	require.NotNil(t, adaptor)
}

func TestLogger_Debug(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	expectedMessage := "debug message"

	mockLogger.EXPECT().
		Debug(expectedMessage).
		Once()

	adaptor := loggeradaptor.NewFactory().New(mockLogger)

	// Act
	adaptor.Debug(expectedMessage)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLogger_Info(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	expectedMessage := "info message"

	mockLogger.EXPECT().
		Info(expectedMessage).
		Once()

	adaptor := loggeradaptor.NewFactory().New(mockLogger)

	// Act
	adaptor.Info(expectedMessage)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLogger_Warn(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	expectedMessage := "warn message"

	mockLogger.EXPECT().
		Warn(expectedMessage).
		Once()

	adaptor := loggeradaptor.NewFactory().New(mockLogger)

	// Act
	adaptor.Warn(expectedMessage)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLogger_Error(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	expectedMessage := "error message"

	mockLogger.EXPECT().
		Error(expectedMessage).
		Once()

	adaptor := loggeradaptor.NewFactory().New(mockLogger)

	// Act
	adaptor.Error(expectedMessage)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLogger_With(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockNewLogger := &entitiesmocks.MockLogger{}
	defer mockNewLogger.AssertExpectations(t)

	expectedKey := "request-id"
	expectedValue := "abc123"
	expectedMessage := "test message"

	mockLogger.EXPECT().
		With(expectedKey, expectedValue).
		Return(mockNewLogger).
		Once()

	mockNewLogger.EXPECT().
		Info(expectedMessage).
		Once()

	adaptor := loggeradaptor.NewFactory().New(mockLogger)

	// Act
	newAdaptor := adaptor.With(expectedKey, expectedValue)

	// Assert
	require.NotNil(t, newAdaptor)
	newAdaptor.Info(expectedMessage)
}

func TestLogger_WithError(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockNewLogger := &entitiesmocks.MockLogger{}
	defer mockNewLogger.AssertExpectations(t)

	expectedError := assert.AnError
	expectedMessage := "test message"

	mockLogger.EXPECT().
		WithError(expectedError).
		Return(mockNewLogger).
		Once()

	mockNewLogger.EXPECT().
		Info(expectedMessage).
		Once()

	adaptor := loggeradaptor.NewFactory().New(mockLogger)

	// Act
	newAdaptor := adaptor.WithError(expectedError)

	// Assert
	require.NotNil(t, newAdaptor)
	newAdaptor.Info(expectedMessage)
}
