package provider

import (
	"context"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/retry"
)

// flightFetcherRetrier wraps a FlightFetcher with retry logic via pkg/retry.
type flightFetcherRetrier struct {
	flightFetcher FlightFetcher
	retrier       *retry.Retrier
}

func NewFlightFetcherRetrier(flightFetcher FlightFetcher, maxRetries int, baseDelay time.Duration) FlightFetcher {
	return &flightFetcherRetrier{
		flightFetcher: flightFetcher,
		retrier:       retry.New(maxRetries, baseDelay),
	}
}

func (r *flightFetcherRetrier) Name() string {
	return r.flightFetcher.Name()
}

func (r *flightFetcherRetrier) Fetch(ctx context.Context, req model.SearchRequest) ([]model.Flight, error) {
	return retry.Execute(ctx, r.retrier, func() ([]model.Flight, error) {
		return r.flightFetcher.Fetch(ctx, req)
	})
}
