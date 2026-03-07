package airasia_test

import (
	"testing"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/airasia"

	"github.com/stretchr/testify/assert"
)

func TestToFlights_SkipFlightWithArrivalBeforeDeparture(t *testing.T) {
	input := model.AirAsiaResponse{
		Status: "ok",
		Flights: []model.AirAsiaFlight{
			{
				FlightCode:    "QZ520",
				Airline:       "AirAsia",
				FromAirport:   "CGK",
				ToAirport:     "DPS",
				DepartTime:    "2025-12-15T10:00:00+07:00",
				ArriveTime:    "2025-12-15T08:00:00+08:00",
				DurationHours: 1.67,
				DirectFlight:  true,
				PriceIDR:      650000,
				Seats:         67,
				CabinClass:    "economy",
			},
		},
	}

	flights := airasia.ToFlights(input)
	assert.Empty(t, flights)
}

func TestToFlights_SkipFlightWithInvalidDatetimeFormat(t *testing.T) {
	input := model.AirAsiaResponse{
		Status: "ok",
		Flights: []model.AirAsiaFlight{
			{
				FlightCode:    "QZ520",
				Airline:       "AirAsia",
				FromAirport:   "CGK",
				ToAirport:     "DPS",
				DepartTime:    "2025-12-15 04:45",
				ArriveTime:    "2025-12-15T07:25:00+08:00",
				DurationHours: 1.67,
				DirectFlight:  true,
				PriceIDR:      650000,
				Seats:         67,
				CabinClass:    "economy",
			},
		},
	}

	flights := airasia.ToFlights(input)
	assert.Empty(t, flights)
}
