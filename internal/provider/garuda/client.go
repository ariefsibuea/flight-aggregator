package garuda

import (
	"context"
	_ "embed"
	"encoding/json"
	"math/rand/v2"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider"
)

//go:embed mock/garuda_indonesia_search_response.json
var mockData []byte

type Client struct{}

// NewClient creates a new Garuda provider client that implements the FlightFetcher interface.
func NewClient() provider.FlightFetcher {
	return &Client{}
}

// Name returns provider's name.
func (c *Client) Name() string {
	return "Garuda Indonesia"
}

// Fetch returns a list of flights from the Garuda provider.
func (c *Client) Fetch(ctx context.Context, req model.SearchRequest) ([]model.Flight, error) {
	// NOTE: Currently, this method only simulates flight provider by returning mock data after a random delay
	// between 50 - 100 milliseconds.

	min, max := 50, 100
	delay := time.Duration(rand.IntN(max-min+1)+min) * time.Millisecond

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	var response model.GarudaResponse
	if err := json.Unmarshal(mockData, &response); err != nil {
		return nil, err
	}

	return GarudaResponseToFlights(response), nil
}
