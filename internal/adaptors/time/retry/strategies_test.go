// Copyright 2026 The MathWorks, Inc.

package retry_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/stretchr/testify/assert"
)

func TestNewLinearRetryStrategy_HappyPath(t *testing.T) {
	// Arrange
	retryInterval := 10 * time.Millisecond

	// Act
	strategy := retry.NewLinearRetryStrategy(retryInterval)

	// Assert
	assert.NotNil(t, strategy)
}

func TestLinearRetryStrategy_C_ReturnsChannel(t *testing.T) {
	// Arrange
	retryInterval := 10 * time.Millisecond
	strategy := retry.NewLinearRetryStrategy(retryInterval)

	// Act
	ch := strategy.C(t.Context())

	// Assert
	assert.NotNil(t, ch)
}

func TestLinearRetryStrategy_C_TicksAtInterval(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		retryInterval := 50 * time.Millisecond
		strategy := retry.NewLinearRetryStrategy(retryInterval)
		ch := strategy.C(t.Context())

		// Assert: not ready before interval
		time.Sleep(retryInterval - time.Nanosecond)
		synctest.Wait()
		select {
		case <-ch:
			t.Fatal("channel ticked before interval elapsed")
		default:
		}

		// Assert: ready after interval
		time.Sleep(time.Nanosecond)
		synctest.Wait()
		select {
		case <-ch:
		default:
			t.Fatal("channel did not tick after interval elapsed")
		}
	})
}

func TestNewLinearRetryStrategy_UsesDefaultForNegativeDuration(t *testing.T) {
	// Arrange
	negativeDuration := -10 * time.Millisecond

	// Act
	strategy := retry.NewLinearRetryStrategy(negativeDuration)

	// Assert
	assert.NotNil(t, strategy)
	assert.NotNil(t, strategy.C(t.Context()))
}

func TestLinearRetryStrategy_CanBeReusedConcurrently(t *testing.T) {
	// Arrange
	retryInterval := 10 * time.Millisecond
	strategy := retry.NewLinearRetryStrategy(retryInterval)

	// Act
	ch1 := strategy.C(t.Context())
	ch2 := strategy.C(t.Context())

	// Assert
	assert.NotNil(t, ch1)
	assert.NotNil(t, ch2)
	assert.NotEqual(t, ch1, ch2, "each call to C() should return a new channel")
}

func TestNewFixedCountRetryStrategy_HappyPath(t *testing.T) {
	// Arrange
	count := 3

	// Act
	strategy := retry.NewFixedCountRetryStrategy(count)

	// Assert
	assert.NotNil(t, strategy)
}

func TestFixedCountRetryStrategy_C_ClosesImmediately_ForCount1_And_InvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{"CountOne", 1},
		{"DefaultForZeroCount", 0},
		{"DefaultForNegativeCount", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			strategy := retry.NewFixedCountRetryStrategy(tt.count)

			// Act
			ch := strategy.C(t.Context())

			// Assert: count=1 means no sends and immediate close
			_, ok := <-ch
			assert.False(t, ok, "channel should be closed immediately")
		})
	}
}

func TestFixedCountRetryStrategy_C_ReturnsChannel(t *testing.T) {
	// Arrange
	strategy := retry.NewFixedCountRetryStrategy(1)

	// Act
	ch := strategy.C(t.Context())

	// Assert
	assert.NotNil(t, ch)
}

func TestFixedCountRetryStrategy_C_SendsCountMinusOneTimesAndCloses(t *testing.T) {
	// Arrange
	count := 3
	strategy := retry.NewFixedCountRetryStrategy(count)

	// Act
	ch := strategy.C(t.Context())
	for range count - 1 {
		<-ch
	}

	// Assert: channel is closed after count-1 sends
	_, ok := <-ch
	assert.False(t, ok, "channel should be closed after count-1 sends")
}

func TestFixedCountRetryStrategy_C_ClosesChannelWhenContextCancelled(t *testing.T) {
	// Arrange
	count := 5
	strategy := retry.NewFixedCountRetryStrategy(count)
	ctx, ctxCancel := context.WithCancel(t.Context())

	// Act
	ctxCancel()
	ch := strategy.C(ctx)

	// Assert: channel is closed when context is cancelled
	_, ok := <-ch
	assert.False(t, ok, "channel should be closed when context is cancelled")
}
