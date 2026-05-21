// Copyright 2026 The MathWorks, Inc.

package retry

import (
	"context"
	"time"
)

const (
	defaultLinearRetryInterval     = 100 * time.Millisecond
	defaultFixedCount          int = 1
)

type linearRetryStrategy struct {
	retryInterval time.Duration
}

func NewLinearRetryStrategy(retryInterval time.Duration) RetryStrategy {
	if retryInterval <= 0 {
		retryInterval = defaultLinearRetryInterval
	}

	return &linearRetryStrategy{
		retryInterval: retryInterval,
	}
}

func (l *linearRetryStrategy) C(_ context.Context) <-chan time.Time {
	return time.Tick(l.retryInterval)
}

func (l *linearRetryStrategy) lock() {}

type fixedCountRetryStrategy struct {
	count int
}

func NewFixedCountRetryStrategy(count int) RetryStrategy {
	if count <= 0 {
		count = defaultFixedCount
	}

	return &fixedCountRetryStrategy{
		count: count,
	}
}

func (f *fixedCountRetryStrategy) C(ctx context.Context) <-chan time.Time {
	c := make(chan time.Time)

	go func() {
		defer close(c)
		for i := 0; i < (f.count - 1); i++ {
			// If the context is cancelled, we're done
			if ctx.Err() != nil {
				return
			}
			// Otherwise, enqueue
			select {
			case c <- time.Now():
			case <-ctx.Done():
				// If the context is cancelled, we're done
				return
			}

		}
	}()

	return c
}

func (f *fixedCountRetryStrategy) lock() {}
