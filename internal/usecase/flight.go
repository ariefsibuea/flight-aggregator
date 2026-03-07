package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider"
)

type FlightUsecase interface {
	SearchFlights(ctx context.Context, req model.SearchRequest) (model.SearchResponse, error)
}

type flightUsecase struct {
	providers []provider.FlightFetcher
}

func NewFlightUsecase(providers []provider.FlightFetcher) FlightUsecase {
	return &flightUsecase{
		providers: providers,
	}
}

func (f *flightUsecase) SearchFlights(ctx context.Context, req model.SearchRequest) (model.SearchResponse, error) {
	startTime := time.Now()

	var (
		goFlights    []model.Flight
		goMetadata   model.SearchMetadata
		backFlights  []model.Flight
		backMetadata model.SearchMetadata
	)

	var wg sync.WaitGroup

	wg.Go(func() {
		goFlights, goMetadata = f.aggregate(ctx, req)
	})

	if req.ReturnDate != nil {
		wg.Go(func() {
			reqReturn := model.SearchRequest{
				Origin:        req.Destination,
				Destination:   req.Origin,
				DepartureDate: *req.ReturnDate,
				Passengers:    req.Passengers,
				CabinClass:    req.CabinClass,
			}
			backFlights, backMetadata = f.aggregate(ctx, reqReturn)
		})
	}

	wg.Wait()

	metadata := mergeMetadata(goMetadata, backMetadata, req.ReturnDate != nil, time.Since(startTime))

	return model.SearchResponse{
		SearchCriteria: model.SearchCriteria{
			Origin:        req.Origin,
			Destination:   req.Destination,
			DepartureDate: req.DepartureDate,
			ReturnDate:    req.ReturnDate,
			Passengers:    req.Passengers,
			CabinClass:    req.CabinClass,
		},
		SearchMetadata: metadata,
		Flights:        append(goFlights, backFlights...),
	}, nil
}

func mergeMetadata(goMeta, backMeta model.SearchMetadata, isRoundTrip bool, dur time.Duration) model.SearchMetadata {
	meta := goMeta
	meta.SearchTimeMS = int(dur.Milliseconds())

	if isRoundTrip {
		meta.TotalResults += backMeta.TotalResults
		meta.ProvidersSucceeded = min(goMeta.ProvidersSucceeded, backMeta.ProvidersSucceeded)
		meta.ProvidersFailed = max(goMeta.ProvidersFailed, backMeta.ProvidersFailed)
	}

	return meta
}
