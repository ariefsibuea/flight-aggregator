package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/pkg/timeutil"
)

type TimeWindow string

const (
	TimeWindowEarlyMorning TimeWindow = "early_morning"
	TimeWindowMorning      TimeWindow = "morning"
	TimeWindowAfternoon    TimeWindow = "afternoon"
	TimeWindowEvening      TimeWindow = "evening"
)

var ValidSortFields = map[string]bool{
	"price":     true,
	"duration":  true,
	"departure": true,
	"arrival":   true,
	"score":     true,
}

var ValidSortDirections = map[string]bool{
	"asc":  true,
	"desc": true,
}

func MatchTimeWindows(hour int, timeWindows []TimeWindow) bool {
	for _, tw := range timeWindows {
		switch tw {
		case TimeWindowEarlyMorning:
			if hour >= 0 && hour < 6 {
				return true
			}
		case TimeWindowMorning:
			if hour >= 6 && hour < 12 {
				return true
			}
		case TimeWindowAfternoon:
			if hour >= 12 && hour < 18 {
				return true
			}
		case TimeWindowEvening:
			if hour >= 18 && hour < 24 {
				return true
			}
		}
	}
	return false
}

type SearchRequest struct {
	// initial request parameters
	Origin        string         `json:"origin"`
	Destination   string         `json:"destination"`
	DepartureDate timeutil.Date  `json:"departure_date"`
	ReturnDate    *timeutil.Date `json:"return_date"`
	Passengers    int            `json:"passengers"`
	CabinClass    string         `json:"cabin_class"`
	// filter parameter
	Filter *SearchFilter `json:"filter"`
	// sort parameter
	Sort *SearchSort `json:"sort"`
}

func (r *SearchRequest) Validate() error {
	if strings.TrimSpace(r.Origin) == "" {
		return fmt.Errorf("origin is empty")
	}
	if strings.TrimSpace(r.Destination) == "" {
		return fmt.Errorf("destination is empty")
	}
	if r.DepartureDate.IsZero() {
		return fmt.Errorf("departure date is empty")
	}
	if r.ReturnDate != nil && !r.ReturnDate.After(r.DepartureDate.Time) {
		return fmt.Errorf("return date must be after departure date")
	}
	if r.Passengers < 1 {
		return fmt.Errorf("minimum passengers is 1")
	}
	if strings.TrimSpace(r.CabinClass) == "" {
		return fmt.Errorf("cabin class is empty")
	}

	if r.Sort != nil {
		if r.Sort.Field != "" && !ValidSortFields[r.Sort.Field] {
			return fmt.Errorf("sort.field must be one of price, duration, departure, arrival, score")
		}
		if r.Sort.Direction != "" && !ValidSortDirections[r.Sort.Direction] {
			return fmt.Errorf("sort.direction must be asc or desc")
		}
	}

	return nil
}

type SearchFilter struct {
	MinPrice       *int64       `json:"min_price,omitempty"`
	MaxPrice       *int64       `json:"max_price,omitempty"`
	MaxStops       *int         `json:"max_stops,omitempty"`
	DepartureTimes []TimeWindow `json:"departure_times,omitempty"`
	ArrivalTimes   []TimeWindow `json:"arrival_times,omitempty"`
	Airlines       []string     `json:"airlines,omitempty"`
	MinDuration    *int         `json:"min_duration,omitempty"`
	MaxDuration    *int         `json:"max_duration,omitempty"`
}

type SearchSort struct {
	Field     string `json:"field,omitempty"`
	Direction string `json:"direction,omitempty"`
}

type SearchResponse struct {
	SearchCriteria  SearchCriteria `json:"search_criteria"`
	SearchMetadata  SearchMetadata `json:"metadata"`
	Flights         []Flight       `json:"flights,omitempty"`
	OutboundFlights []Flight       `json:"outbound_flights,omitempty"`
	InboundFlights  []Flight       `json:"inbound_flights,omitempty"`
}

type SearchCriteria struct {
	Origin        string         `json:"origin"`
	Destination   string         `json:"destination"`
	DepartureDate timeutil.Date  `json:"departure_date"`
	ReturnDate    *timeutil.Date `json:"return_date,omitempty"`
	Passengers    int            `json:"passengers"`
	CabinClass    string         `json:"cabin_class"`
}

type SearchMetadata struct {
	TotalResults       int  `json:"total_results"`
	ProvidersQueried   int  `json:"providers_queried"`
	ProvidersSucceeded int  `json:"providers_succeeded"`
	ProvidersFailed    int  `json:"providers_failed"`
	SearchTimeMS       int  `json:"search_time_ms"`
	CacheHit           bool `json:"cache_hit"`
}

type Flight struct {
	ID             string         `json:"id"`
	Provider       string         `json:"provider"`
	Airline        Airline        `json:"airline"`
	FlightNumber   string         `json:"flight_number"`
	Departure      FlightEndpoint `json:"departure"`
	Arrival        FlightEndpoint `json:"arrival"`
	Duration       Duration       `json:"duration"`
	Stops          int            `json:"stops"`
	Price          Price          `json:"price"`
	AvailableSeats int            `json:"available_seats"`
	CabinClass     string         `json:"cabin_class"`
	Aircraft       *string        `json:"aircraft"`
	Amenities      []string       `json:"amenities"`
	Baggage        Baggage        `json:"baggage"`
	BestValueScore float64        `json:"-"`
}

type Airline struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type FlightEndpoint struct {
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
	Display  string `json:"display,omitempty"`
}

type Baggage struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}
