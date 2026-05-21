// Copyright 2026 The MathWorks, Inc.

package retry_test

import (
	"context"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetry_SucceedsOnFirstAttempt(t *testing.T) {
	// Arrange
	expectedResult := "success"
	fn := func() (string, bool, error) {
		return expectedResult, true, nil
	}
	retryStrategy := retry.NewLinearRetryStrategy(10 * time.Millisecond)

	// Act
	result, err := retry.Retry(t.Context(), fn, retryStrategy)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestRetry_SucceedsAfterMultipleAttempts(t *testing.T) {
	// Arrange
	expectedResult := "success"
	attempts := 0
	fn := func() (string, bool, error) {
		attempts++
		if attempts < 3 {
			return "", false, nil
		}
		return expectedResult, true, nil
	}
	retryStrategy := retry.NewLinearRetryStrategy(10 * time.Millisecond)

	// Act
	result, err := retry.Retry(t.Context(), fn, retryStrategy)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 3, attempts)
}

func TestRetry_ReturnsErrorFromFunction(t *testing.T) {
	// Arrange
	expectedErr := assert.AnError
	fn := func() (string, bool, error) {
		return "", false, expectedErr
	}
	retryStrategy := retry.NewLinearRetryStrategy(10 * time.Millisecond)

	// Act
	result, err := retry.Retry(t.Context(), fn, retryStrategy)

	// Assert
	require.ErrorIs(t, err, expectedErr)
	assert.Empty(t, result)
}

func TestRetry_ReturnsContextError(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(t.Context())
	attempts := 0
	fn := func() (string, bool, error) {
		attempts++
		if attempts == 2 {
			cancel()
		}
		return "", false, nil
	}
	retryStrategy := retry.NewLinearRetryStrategy(10 * time.Millisecond)

	// Act
	result, err := retry.Retry(ctx, fn, retryStrategy)

	// Assert
	require.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, result)
}

func TestRetry_ReturnsContextCancelCauseError(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancelCause(t.Context())
	expectedError := assert.AnError

	attempts := 0
	fn := func() (string, bool, error) {
		attempts++
		if attempts == 2 {
			cancel(expectedError)
		}
		return "", false, nil
	}
	retryStrategy := retry.NewLinearRetryStrategy(10 * time.Millisecond)

	// Act
	result, err := retry.Retry(ctx, fn, retryStrategy)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestRetry_ReturnsExhaustedErrorWhenStrategyExhausted(t *testing.T) {
	// Arrange
	fn := func() (string, bool, error) {
		return "", false, nil
	}
	retryStrategy := retry.NewFixedCountRetryStrategy(3)

	// Act
	result, err := retry.Retry(t.Context(), fn, retryStrategy)

	// Assert
	require.ErrorIs(t, err, retry.ErrRetryStrategyExhausted)
	assert.Empty(t, result)
}

func TestRetry_ReturnsErrorForNilRetryStrategy(t *testing.T) {
	// Arrange
	fn := func() (string, bool, error) {
		return "success", true, nil
	}

	// Act
	result, err := retry.Retry(t.Context(), fn, nil)

	// Assert
	require.ErrorIs(t, err, retry.ErrInvalidRetryStrategy)
	assert.Empty(t, result)
}

func TestRetry_ReturnsErrorWhenContextAlreadyCancelled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	fnCalled := false
	fn := func() (string, bool, error) {
		fnCalled = true
		return "success", true, nil
	}
	retryStrategy := retry.NewLinearRetryStrategy(10 * time.Millisecond)

	// Act
	result, err := retry.Retry(ctx, fn, retryStrategy)

	// Assert
	require.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, result)
	assert.False(t, fnCalled, "fn should not be called when context is already cancelled")
}
