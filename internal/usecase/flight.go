package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/cache"
	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider"

	"golang.org/x/sync/errgroup"
)

type FlightUsecase interface {
	SearchFlights(ctx context.Context, req model.SearchRequest) (model.SearchResponse, error)
}

type flightUsecase struct {
	providers       []provider.FlightFetcher
	cache           cache.Cache
	defaultCacheTTL time.Duration
}

func NewFlightUsecase(providers []provider.FlightFetcher, c cache.Cache, ttl time.Duration) FlightUsecase {
	return &flightUsecase{
		providers:       providers,
		cache:           c,
		defaultCacheTTL: ttl,
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

	criteria := model.SearchCriteria{
		Origin:        req.Origin,
		Destination:   req.Destination,
		DepartureDate: req.DepartureDate,
		ReturnDate:    req.ReturnDate,
		Passengers:    req.Passengers,
		CabinClass:    req.CabinClass,
	}

	if req.ReturnDate != nil {
		// Round-trip: filter and sort each leg independently, expose as separate arrays.
		outboundFlights = filterAndSort(outboundFlights, req)
		inboundFlights = filterAndSort(inboundFlights, req)

		metadata := mergeMetadata(outboundMetadata, inboundMetadata, true, time.Since(startTime))
		metadata.TotalResults = len(outboundFlights) + len(inboundFlights)

		return model.SearchResponse{
			SearchCriteria:  criteria,
			SearchMetadata:  metadata,
			OutboundFlights: outboundFlights,
			InboundFlights:  inboundFlights,
		}, nil
	}

	// one-way search.
	outboundFlights = filterAndSort(outboundFlights, req)

	metadata := mergeMetadata(outboundMetadata, model.SearchMetadata{}, false, time.Since(startTime))
	metadata.TotalResults = len(outboundFlights)

	return model.SearchResponse{
		SearchCriteria: criteria,
		SearchMetadata: metadata,
		Flights:        outboundFlights,
	}, nil
}

func mergeMetadata(outboundMeta, inboundMeta model.SearchMetadata, isRoundTrip bool, dur time.Duration) model.SearchMetadata {
	meta := outboundMeta
	meta.SearchTimeMS = int(dur.Milliseconds())

	if isRoundTrip {
		meta.ProvidersSucceeded = min(outboundMeta.ProvidersSucceeded, inboundMeta.ProvidersSucceeded)
		meta.ProvidersFailed = max(outboundMeta.ProvidersFailed, inboundMeta.ProvidersFailed)
	}

	return meta
}
