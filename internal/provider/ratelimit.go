package provider

import (
	"context"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/ratelimit"
)

// flightFetcherLimiter wraps a FlightFetcher with a token bucket rate limiter via pkg/ratelimit.
type flightFetcherLimiter struct {
	flightFetcher FlightFetcher
	limiter       *ratelimit.Limiter
}

func NewFlightFetcherLimiter(flightFetcher FlightFetcher, rps float64) FlightFetcher {
	return &flightFetcherLimiter{
		flightFetcher: flightFetcher,
		limiter:       ratelimit.New(rps, 2),
	}
}

func (r *flightFetcherLimiter) Name() string {
	return r.flightFetcher.Name()
}

func (r *flightFetcherLimiter) Fetch(ctx context.Context, req model.SearchRequest) ([]model.Flight, error) {
	return ratelimit.Execute(ctx, r.limiter, func() ([]model.Flight, error) {
		return r.flightFetcher.Fetch(ctx, req)
	})
}
