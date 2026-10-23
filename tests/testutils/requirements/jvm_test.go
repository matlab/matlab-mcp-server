// Copyright 2026 The MathWorks, Inc.

package requirements_test

import (
	"testing"

	mocks "github.com/matlab/matlab-mcp-server/tests/mocks/testutils/requirements"
	"github.com/matlab/matlab-mcp-server/tests/testutils/requirements"
	"github.com/stretchr/testify/assert"
)

func TestJVM_JVMAvailable_DoesNotSkip(t *testing.T) {
	// Arrange
	mockT := mocks.NewMockTestingT(t)
	defer mockT.AssertExpectations(t)
	mockSession := mocks.NewMockMATLABSession(t)
	defer mockSession.AssertExpectations(t)
	ctx := t.Context()

	mockT.EXPECT().Helper().Once()
	mockSession.EXPECT().EvaluateCode(ctx, requirements.JVMCheckCode).Return("MCP_USEJAVA_JVM=1\n", nil).Once()

	// Act
	requirements.JVM(ctx, mockT, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestJVM_JVMUnavailable_Skips(t *testing.T) {
	// Arrange
	mockT := mocks.NewMockTestingT(t)
	defer mockT.AssertExpectations(t)
	mockSession := mocks.NewMockMATLABSession(t)
	defer mockSession.AssertExpectations(t)
	ctx := t.Context()

	t.Setenv(requirements.JVMForceEnvVar, "")

	mockT.EXPECT().Helper().Twice()
	mockSession.EXPECT().EvaluateCode(ctx, requirements.JVMCheckCode).Return("MCP_USEJAVA_JVM=0\n", nil).Once()
	mockT.EXPECT().Skipf("%s. Set %s=1 to run the test anyway", []any{requirements.JVMUnavailableReason, requirements.JVMForceEnvVar}).Once()

	// Act
	requirements.JVM(ctx, mockT, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestJVM_JVMUnavailable_Forced_RunsAnyway(t *testing.T) {
	// Arrange
	mockT := mocks.NewMockTestingT(t)
	defer mockT.AssertExpectations(t)
	mockSession := mocks.NewMockMATLABSession(t)
	defer mockSession.AssertExpectations(t)
	ctx := t.Context()

	t.Setenv(requirements.JVMForceEnvVar, "1")

	mockT.EXPECT().Helper().Twice()
	mockSession.EXPECT().EvaluateCode(ctx, requirements.JVMCheckCode).Return("MCP_USEJAVA_JVM=0\n", nil).Once()
	mockT.EXPECT().Logf("%s. Running anyway because %s is set", []any{requirements.JVMUnavailableReason, requirements.JVMForceEnvVar}).Once()

	// Act
	requirements.JVM(ctx, mockT, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestJVM_EvaluateCodeFails_Fails(t *testing.T) {
	// Arrange
	mockT := mocks.NewMockTestingT(t)
	defer mockT.AssertExpectations(t)
	mockSession := mocks.NewMockMATLABSession(t)
	defer mockSession.AssertExpectations(t)
	ctx := t.Context()

	mockT.EXPECT().Helper().Once()
	mockSession.EXPECT().EvaluateCode(ctx, requirements.JVMCheckCode).Return("", assert.AnError).Once()
	mockT.EXPECT().Fatalf(`usejava('jvm') should run in the MATLAB under test: %v`, []any{assert.AnError}).Once()

	// Act
	requirements.JVM(ctx, mockT, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestJVM_UnexpectedOutput_Fails(t *testing.T) {
	// Arrange
	mockT := mocks.NewMockTestingT(t)
	defer mockT.AssertExpectations(t)
	mockSession := mocks.NewMockMATLABSession(t)
	defer mockSession.AssertExpectations(t)
	ctx := t.Context()

	mockT.EXPECT().Helper().Once()
	mockSession.EXPECT().EvaluateCode(ctx, requirements.JVMCheckCode).Return("unexpected\n", nil).Once()
	mockT.EXPECT().Fatalf(`usejava('jvm') should print MCP_USEJAVA_JVM=0 or 1, got %q`, []any{"unexpected\n"}).Once()

	// Act
	requirements.JVM(ctx, mockT, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}
