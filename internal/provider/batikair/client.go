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
	// NOTE: Currently, this method returns mock data with a random delay between 100 - 200 milliseconds. Below, it
	// generates the delay, pauses execution, and then returns the mock data.
	min, max := 100, 200
	delay := time.Duration(rand.IntN(max-min+1)+min) * time.Millisecond
	time.Sleep(delay)

	// Unmarshal the embedded JSON data into a model.BatikAirResponse.
	var response model.BatikAirResponse
	if err := json.Unmarshal(mockData, &response); err != nil {
		return nil, err
	}

	return BatikAirResponseToFlights(response), nil
}
