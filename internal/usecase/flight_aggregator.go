package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/cache"
	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider"
)

type providerResult struct {
	provider string
	flights  []model.Flight
	err      error
}

func (f *flightUsecase) aggregate(ctx context.Context, req model.SearchRequest) ([]model.Flight, model.SearchMetadata) {
	var flights []model.Flight

	cacheKey := generateCacheKey(req)
	err := f.cache.Get(ctx, cacheKey, &flights)
	if err == nil {
		// no providers were queried
		return flights, model.SearchMetadata{
			TotalResults: len(flights),
			CacheHit:     true,
		}
	}
	if err != cache.ErrCacheKeyNotFound {
		slog.Warn("failed to get flights from cache", "key", cacheKey, "err", err)
	}

	startTime := time.Now()

	// NOTE: Implement fan-out,fan-in here, inspired by https://go.dev/blog/pipelines.
	numProviders := len(f.providers)
	resultsChan := make(chan providerResult, numProviders)
	var wg sync.WaitGroup

	for _, p := range f.providers {
		wg.Add(1)

		go func(fetcher provider.FlightFetcher) {
			defer wg.Done()

			res, err := fetcher.Fetch(ctx, req)
			resultsChan <- providerResult{
				provider: fetcher.Name(),
				flights:  res,
				err:      err,
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	providersSucceeded := 0
	providersFailed := 0

	for result := range resultsChan {
		if result.err != nil {
			providersFailed++
			slog.Warn("provider fetch failed", "provider", result.provider, "err", result.err)
			continue
		}
		providersSucceeded++
		flights = append(flights, result.flights...)
	}

	if providersSucceeded > 0 {
		if err := f.cache.Set(ctx, cacheKey, flights, f.defaultCacheTTL); err != nil {
			slog.Warn("failed to cache flight results", "key", cacheKey, "err", err)
		}
	}

	searchDuration := time.Since(startTime)
	metadata := model.SearchMetadata{
		TotalResults:       len(flights),
		ProvidersQueried:   numProviders,
		ProvidersSucceeded: providersSucceeded,
		ProvidersFailed:    providersFailed,
		SearchTimeMS:       int(searchDuration.Milliseconds()),
		CacheHit:           false,
	}

	return flights, metadata
}

func generateCacheKey(req model.SearchRequest) string {
	return fmt.Sprintf("flights:%s:%s:%s:%d:%s",
		strings.ToUpper(req.Origin),
		strings.ToUpper(req.Destination),
		req.DepartureDate.Format(time.DateOnly),
		req.Passengers,
		strings.ToLower(req.CabinClass),
	)
}
