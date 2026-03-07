package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider"
)

type providerResult struct {
	provider string
	flights  []model.Flight
	err      error
}

func (f *flightUsecase) aggregate(ctx context.Context, req model.SearchRequest) ([]model.Flight, model.SearchMetadata) {
	startTime := time.Now()

	numProviders := len(f.providers)
	resultsChan := make(chan providerResult, numProviders)

	// NOTE: The pattern of fan-out,fan-in inspired by https://go.dev/blog/pipelines.
	var wg sync.WaitGroup

	for _, p := range f.providers {
		wg.Add(1)

		go func(fetcher provider.FlightFetcher) {
			defer wg.Done()

			flights, err := fetcher.Fetch(ctx, req)
			resultsChan <- providerResult{
				provider: fetcher.Name(),
				flights:  flights,
				err:      err,
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	flights := make([]model.Flight, 0)
	providersSucceeded := 0
	providersFailed := 0

	for result := range resultsChan {
		if result.err != nil {
			providersFailed++
			continue
		}
		providersSucceeded++
		flights = append(flights, result.flights...)
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
