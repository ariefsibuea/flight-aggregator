package batikair_test

import (
	"testing"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/batikair"

	"github.com/stretchr/testify/assert"
)

func TestToFlights_SkipFlightWithArrivalBeforeDeparture(t *testing.T) {
	input := model.BatikAirResponse{
		Code:    200,
		Message: "OK",
		Results: []model.BatikAirFlight{
			{
				FlightNumber:      "ID6514",
				AirlineName:       "Batik Air",
				AirlineIATA:       "ID",
				Origin:            "CGK",
				Destination:       "DPS",
				DepartureDateTime: "2025-12-15T10:00:00+0700",
				ArrivalDateTime:   "2025-12-15T08:00:00+0800",
				TravelTime:        "1h 45m",
				NumberOfStops:     0,
				Fare: model.BatikAirFare{
					BasePrice:    980000,
					Taxes:        120000,
					TotalPrice:   1100000,
					CurrencyCode: "IDR",
					Class:        "Y",
				},
				SeatsAvailable: 32,
			},
		},
	}

	flights := batikair.ToFlights(input)
	assert.Empty(t, flights)
}

func TestToFlights_SkipFlightWithInvalidDatetimeFormat(t *testing.T) {
	input := model.BatikAirResponse{
		Code:    200,
		Message: "OK",
		Results: []model.BatikAirFlight{
			{
				FlightNumber:      "ID6514",
				AirlineName:       "Batik Air",
				AirlineIATA:       "ID",
				Origin:            "CGK",
				Destination:       "DPS",
				DepartureDateTime: "2025-12-15 07:15",
				ArrivalDateTime:   "2025-12-15T10:00:00+0800",
				TravelTime:        "1h 45m",
				NumberOfStops:     0,
				Fare: model.BatikAirFare{
					BasePrice:    980000,
					Taxes:        120000,
					TotalPrice:   1100000,
					CurrencyCode: "IDR",
					Class:        "Y",
				},
				SeatsAvailable: 32,
			},
		},
	}

	flights := batikair.ToFlights(input)
	assert.Empty(t, flights)
}
