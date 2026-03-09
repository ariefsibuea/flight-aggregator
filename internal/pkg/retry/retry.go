package retry

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"
)

// Retrier handles retry logic with exponential backoff and jitter.
type Retrier struct {
	maxRetries int
	baseDelay  time.Duration
}

// New creates a new Retrier.
func New(maxRetries int, baseDelay time.Duration) *Retrier {
	if baseDelay <= 0 {
		baseDelay = 100 * time.Millisecond // Default fallback
	}

	return &Retrier{
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
	}
}

// Do executes fn with retry logic. It stops retrying if the context is cancelled or if the error is a context error.
func (r *Retrier) Do(ctx context.Context, fn func() error) error {
	var err error
	delay := r.baseDelay

	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if attempt > 0 {
			// Full jitter: sleep random duration between [0, delay)
			jitter := time.Duration(rand.N(int64(delay)))
			timer := time.NewTimer(jitter)

			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
				delay *= 2
			}
		}

		err = fn()
		if err == nil {
			return nil
		}

		// Don't retry if context is already dead or error is unrecoverable
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
	}

	return err
}

// Execute is a generic wrapper around Do for functions that return a value.
func Execute[T any](ctx context.Context, r *Retrier, fn func() (T, error)) (T, error) {
	var result T
	err := r.Do(ctx, func() error {
		var e error
		result, e = fn()
		return e
	})
	return result, err
}
