package ratelimit

import (
	"context"

	"golang.org/x/time/rate"
)

// Limiter controls the frequency of events using a token bucket algorithm.
type Limiter struct {
	limiter *rate.Limiter
}

// New creates a new Limiter with the specified requests per second (rps) and burst capacity.
func New(rps float64, burst int) *Limiter {
	return &Limiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

// Do waits for the limiter to allow a request, then executes the provided function. It blocks until
// a token is available or the context is cancelled.
func (l *Limiter) Do(ctx context.Context, fn func() error) error {
	if err := l.limiter.Wait(ctx); err != nil {
		return err
	}
	return fn()
}

// Execute waits for the limiter to allow a request, then executes the provided function.
// It returns the result of the function and any error.
func Execute[T any](ctx context.Context, l *Limiter, fn func() (T, error)) (T, error) {
	if err := l.limiter.Wait(ctx); err != nil {
		var zero T
		return zero, err
	}
	return fn()
}
