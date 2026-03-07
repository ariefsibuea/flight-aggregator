package lionair

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider"
)

//go:embed mock/lion_air_search_response.json
var mockData []byte

type Client struct{}

// NewClient creates a new LionAir provider client that implements the FlightFetcher interface.
func NewClient() provider.FlightFetcher {
	return &Client{}
}

// Name returns provider's name.
func (c *Client) Name() string {
	return "Lion Air"
}

// Fetch returns a list of flights from the LionAir provider.
func (c *Client) Fetch(ctx context.Context, req model.SearchRequest) ([]model.Flight, error) {
	// NOTE: Currently, this method only simulates flight provider by returning mock data after a random delay
	// between 100 - 200 milliseconds.

	min, max := 100, 200
	delay := time.Duration(rand.IntN(max-min+1)+min) * time.Millisecond

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	var res model.LionAirResponse
	if err := json.Unmarshal(mockData, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return ToFlights(res), nil
}
