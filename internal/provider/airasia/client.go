package airasia

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

//go:embed mock/airasia_search_response.json
var mockData []byte

type Client struct{}

// NewClient creates a new AirAsia provider client that implements the FlightFetcher interface.
func NewClient() provider.FlightFetcher {
	return &Client{}
}

// Name returns provider's name.
func (c *Client) Name() string {
	return "AirAsia"
}

// Fetch returns a list of flights from the AirAsia provider.
func (c *Client) Fetch(ctx context.Context, req model.SearchRequest) ([]model.Flight, error) {
	// NOTE: Currently, this method only simulates flight provider by returning mock data after a random delay
	// between 50 - 150 milliseconds and the success rate is 90%.

	min, max := 50, 150
	delay := time.Duration(rand.IntN(max-min+1)+min) * time.Millisecond

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if rand.IntN(100) < 10 {
		return nil, fmt.Errorf("failed to fetch AirAsia flights")
	}

	var res model.AirAsiaResponse
	if err := json.Unmarshal(mockData, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return ToFlights(res), nil
}
