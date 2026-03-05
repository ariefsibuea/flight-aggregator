package batikair

import (
	"context"
	_ "embed"
	"encoding/json"
	"math/rand/v2"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider"
)

//go:embed mock/batik_air_search_response.json
var mockData []byte

type Client struct{}

// NewClient creates a new BatikAir provider client that implements the FlightFetcher interface.
func NewClient() provider.FlightFetcher {
	return &Client{}
}

// Fetch returns a list of flights from the BatikAir provider.
func (c *Client) Fetch(ctx context.Context, req model.SearchRequest) ([]model.Flight, error) {
	// NOTE: Currently, this method only simulates flight provider by returning mock data after a random delay
	// between 200 - 400 milliseconds.

	min, max := 200, 400
	delay := time.Duration(rand.IntN(max-min+1)+min) * time.Millisecond

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	var response model.BatikAirResponse
	if err := json.Unmarshal(mockData, &response); err != nil {
		return nil, err
	}

	return BatikAirResponseToFlights(response), nil
}
