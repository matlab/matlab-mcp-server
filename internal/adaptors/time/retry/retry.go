// Copyright 2026 The MathWorks, Inc.

package retry

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidRetryStrategy   = errors.New("invalid retry strategy")
	ErrRetryStrategyExhausted = errors.New("exhausted retry strategy")
)

type RetryStrategy interface {
	C(ctx context.Context) <-chan time.Time
	lock() // unexported to prevent external implementations; add new strategies to this package
}

// Retry repeatedly calls fn using the given retry strategy until one of the following conditions is met:
//   - fn returns a non-nil error: stop immediately and return the error.
//   - fn returns (output, true, nil): stop retrying and return output.
//   - fn returns (_, false, nil): keep retrying.
//   - The retry strategy is exhausted: stop and return ErrRetryStrategyExhausted.
//   - ctx is canceled: stop and return ctx.Err() (or context.Cause if set).
func Retry[OutputType any](ctx context.Context, fn func() (OutputType, bool, error), retryStrategy RetryStrategy) (OutputType, error) {
	var zeroValue OutputType

	if retryStrategy == nil {
		return zeroValue, ErrInvalidRetryStrategy
	}

	if err := context.Cause(ctx); err != nil {
		return zeroValue, err
	}

	retryCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	retryC := retryStrategy.C(retryCtx)

	for {
		output, ok, err := fn()
		if err != nil {
			return zeroValue, err
		}

		if ok {
			return output, nil
		}

		select {
		case _, ok := <-retryC:
			if !ok {
				// Channel closed, the strategy is exhausted
				return zeroValue, ErrRetryStrategyExhausted
			}
			// Try again
		case <-retryCtx.Done():
			return zeroValue, context.Cause(retryCtx)
		}
	}
}
