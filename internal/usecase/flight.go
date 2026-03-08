package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider"

	"golang.org/x/sync/errgroup"
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
		outboundFlights  []model.Flight
		outboundMetadata model.SearchMetadata
		inboundFlights   []model.Flight
		inboundMetadata  model.SearchMetadata
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		outboundFlights, outboundMetadata = f.aggregate(ctx, req)
		return nil
	})

	if req.ReturnDate != nil {
		g.Go(func() error {
			reqReturn := model.SearchRequest{
				Origin:        req.Destination,
				Destination:   req.Origin,
				DepartureDate: *req.ReturnDate,
				Passengers:    req.Passengers,
				CabinClass:    req.CabinClass,
			}
			inboundFlights, inboundMetadata = f.aggregate(ctx, reqReturn)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return model.SearchResponse{}, fmt.Errorf("found an error while executing concurrent flight search: %w", err)
	}

	// TODO: For round-trip, expose outbound and inbound as separate arrays in the response.
	flights := append(outboundFlights, inboundFlights...)
	flights = filterAndSort(flights, req)

	metadata := mergeMetadata(outboundMetadata, inboundMetadata, req.ReturnDate != nil, time.Since(startTime))
	metadata.TotalResults = len(flights)

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
		Flights:        flights,
	}, nil
}

func mergeMetadata(outboundMeta, inboundMeta model.SearchMetadata, isRoundTrip bool, dur time.Duration) model.SearchMetadata {
	meta := outboundMeta
	meta.SearchTimeMS = int(dur.Milliseconds())

	if isRoundTrip {
		meta.TotalResults += inboundMeta.TotalResults
		meta.ProvidersSucceeded = min(outboundMeta.ProvidersSucceeded, inboundMeta.ProvidersSucceeded)
		meta.ProvidersFailed = max(outboundMeta.ProvidersFailed, inboundMeta.ProvidersFailed)
	}

	return meta
}
