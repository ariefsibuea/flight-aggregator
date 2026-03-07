package garuda_test

import (
	"testing"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/garuda"

	"github.com/stretchr/testify/assert"
)

func TestToFlights_SkipFlightWithArrivalBeforeDeparture(t *testing.T) {
	input := model.GarudaResponse{
		Status: "success",
		Flights: []model.GarudaFlight{
			{
				FlightID:        "GA400",
				Airline:         "Garuda Indonesia",
				AirlineCode:     "GA",
				DurationMinutes: 110,
				Stops:           0,
				FareClass:       "economy",
				Departure: model.GarudaEndpoint{
					Airport: "CGK",
					City:    "Jakarta",
					Time:    "2025-12-15T10:00:00+07:00",
				},
				Arrival: model.GarudaEndpoint{
					Airport: "DPS",
					City:    "Denpasar",
					Time:    "2025-12-15T08:00:00+07:00",
				},
				Price:          model.GarudaPrice{Amount: 1000000, Currency: "IDR"},
				AvailableSeats: 10,
			},
		},
	}

	flights := garuda.ToFlights(input)
	assert.Empty(t, flights)
}

func TestToFlights_SkipFlightWithInvalidDatetimeFormat(t *testing.T) {
	input := model.GarudaResponse{
		Status: "success",
		Flights: []model.GarudaFlight{
			{
				FlightID:    "GA400",
				Airline:     "Garuda Indonesia",
				AirlineCode: "GA",
				FareClass:   "economy",
				Departure: model.GarudaEndpoint{
					Airport: "CGK",
					Time:    "2025-12-15 08:00",
				},
				Arrival: model.GarudaEndpoint{
					Airport: "DPS",
					Time:    "2025-12-15T08:00:00+07:00",
				},
				Price:          model.GarudaPrice{Amount: 1000000, Currency: "IDR"},
				AvailableSeats: 10,
			},
		},
	}

	flights := garuda.ToFlights(input)
	assert.Empty(t, flights)
}
