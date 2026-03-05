package models

import "time"

type SearchRequest struct {
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	DepartureDate time.Time `json:"departureDate"`
	ReturnDate    time.Time `json:"returnDate"`
	Passengers    int32     `json:"passengers"`
	CabinClass    string    `json:"cabinClass"`
}

type SearchResponse struct {
	SearchCriteria SearchCriteria `json:"search_criteria"`
	SearchMetadata SearchMetadata `json:"metadata"`
	Flights        []Flight       `json:"flights"`
}

type SearchCriteria struct {
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	DepartureDate time.Time `json:"departure_date"`
	Passengers    int32     `json:"passengers"`
	CabinClass    string    `json:"cabin_class"`
}

type SearchMetadata struct {
	TotalResults       int32 `json:"total_results"`
	ProvidersQueried   int32 `json:"providers_queried"`
	ProvidersSucceeded int32 `json:"providers_succeeded"`
	ProvidersFailed    int32 `json:"providers_failed"`
	SearchTimeMS       int32 `json:"search_time_ms"`
	CacheHit           bool  `json:"cache_hit"`
}

type Flight struct {
	ID             string      `json:"id"`
	Provider       string      `json:"provider"`
	Airline        Airline     `json:"airline"`
	FlightNumber   string      `json:"flight_number"`
	Departure      Departure   `json:"departure"`
	Arrival        Arrival     `json:"arrival"`
	Duration       Duration    `json:"duration"`
	Stops          int32       `json:"stops"`
	Price          Price       `json:"price"`
	AvailableSeats int32       `json:"available_seats"`
	CabinClass     string      `json:"cabin_class"`
	Aircraft       interface{} `json:"aircraft"`
	Amenities      []string    `json:"amenities"`
	Baggage        Baggage     `json:"baggage"`
}

type Airline struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Departure struct {
	Airport   string    `json:"airport"`
	City      string    `json:"city"`
	Datetime  time.Time `json:"datetime"`
	Timestamp int64     `json:"timestamp"`
}

type Arrival struct {
	Airport   string    `json:"airport"`
	City      string    `json:"city"`
	Datetime  time.Time `json:"datetime"`
	Timestamp int64     `json:"timestamp"`
}

type Duration struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

type Price struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type Baggage struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}
