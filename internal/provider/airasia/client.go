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

// Fetch returns a list of flights from the AirAsia provider.
func (c *Client) Fetch(ctx context.Context, req model.SearchRequest) ([]model.Flight, error) {
	// NOTE: Currently, this method returns mock data with a random delay between 50 - 150 milliseconds. Below, it
	// generates a 10% failure probability, generates the delay time, pauses execution, and then returns the mock data.
	if rand.IntN(100) < 10 {
		return nil, fmt.Errorf("failed to fetch AirAsia flights")
	}

	min, max := 50, 150
	delay := time.Duration(rand.IntN(max-min+1)+min) * time.Millisecond
	time.Sleep(delay)

	// Unmarshal the embedded JSON data into a model.AirAsiaResponse.
	var response model.AirAsiaResponse
	if err := json.Unmarshal(mockData, &response); err != nil {
		return nil, err
	}

	return AirAsiaResponseToFlights(response), nil
}
