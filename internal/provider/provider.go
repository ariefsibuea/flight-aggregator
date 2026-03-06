package provider

import (
	"context"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
)

type FlightFetcher interface {
	Name() string
	Fetch(ctx context.Context, req model.SearchRequest) ([]model.Flight, error)
}
