// Copyright 2026 The MathWorks, Inc.

package requirements_test

import (
	"testing"

	mocks "github.com/matlab/matlab-mcp-server/tests/mocks/testutils/requirements"
	"github.com/matlab/matlab-mcp-server/tests/testutils/requirements"
)

const (
	testForceEnvVar = "MCP_SYSTEM_TEST_FORCE_SOMETHING"
	testReason      = "something is missing"
)

func TestUnavailable_ForceEnvVarUnset_Skips(t *testing.T) {
	// Arrange
	mockT := mocks.NewMockTestingT(t)
	defer mockT.AssertExpectations(t)

	t.Setenv(testForceEnvVar, "")

	mockT.EXPECT().Helper().Once()
	mockT.EXPECT().Skipf("%s. Set %s=1 to run the test anyway", []any{testReason, testForceEnvVar}).Once()

	// Act
	requirements.Unavailable(mockT, testForceEnvVar, testReason)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestUnavailable_ForceEnvVarZero_Skips(t *testing.T) {
	// Arrange
	mockT := mocks.NewMockTestingT(t)
	defer mockT.AssertExpectations(t)

	t.Setenv(testForceEnvVar, "0")

	mockT.EXPECT().Helper().Once()
	mockT.EXPECT().Skipf("%s. Set %s=1 to run the test anyway", []any{testReason, testForceEnvVar}).Once()

	// Act
	requirements.Unavailable(mockT, testForceEnvVar, testReason)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestUnavailable_ForceEnvVarOne_RunsAnyway(t *testing.T) {
	// Arrange
	mockT := mocks.NewMockTestingT(t)
	defer mockT.AssertExpectations(t)

	t.Setenv(testForceEnvVar, "1")

	mockT.EXPECT().Helper().Once()
	mockT.EXPECT().Logf("%s. Running anyway because %s is set", []any{testReason, testForceEnvVar}).Once()

	// Act
	requirements.Unavailable(mockT, testForceEnvVar, testReason)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestUnavailable_ForceEnvVarInvalid_Fails(t *testing.T) {
	// Arrange
	mockT := mocks.NewMockTestingT(t)
	defer mockT.AssertExpectations(t)

	t.Setenv(testForceEnvVar, "yes")

	mockT.EXPECT().Helper().Once()
	mockT.EXPECT().Fatalf("%s should be 1 or 0, got %q", []any{testForceEnvVar, "yes"}).Once()

	// Act
	requirements.Unavailable(mockT, testForceEnvVar, testReason)

	// Assert
	// Assertions are verified via deferred mock expectations.
}
