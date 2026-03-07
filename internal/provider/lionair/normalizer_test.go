package lionair_test

import (
	"testing"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/lionair"

	"github.com/stretchr/testify/assert"
)

func TestToFlights_SkipFlightWithArrivalBeforeDeparture(t *testing.T) {
	input := model.LionAirResponse{
		Success: true,
		Data: model.LionAirData{
			AvailableFlights: []model.LionAirFlight{
				{
					ID:      "JT740",
					Carrier: model.LionAirCarrier{Name: "Lion Air", IATA: "JT"},
					Route: model.LionAirRoute{
						From: model.LionAirAirport{
							Code: "CGK",
							Name: "Soekarno-Hatta",
							City: "Jakarta",
						},
						To: model.LionAirAirport{
							Code: "DPS",
							Name: "Ngurah Rai",
							City: "Denpasar",
						},
					},
					Schedule:   model.LionAirSchedule{Departure: "2025-12-15T10:00:00", DepartureTimezone: "Asia/Jakarta", Arrival: "2025-12-15T08:15:00", ArrivalTimezone: "Asia/Makassar"},
					FlightTime: 105,
					IsDirect:   true,
					Pricing:    model.LionAirPricing{Total: 950000, Currency: "IDR", FareType: "ECONOMY"},
					SeatsLeft:  45,
					PlaneType:  "Boeing 737-900ER",
				},
			},
		},
	}

	flights := lionair.ToFlights(input)
	assert.Empty(t, flights)
}

func TestToFlights_SkipFlightWithInvalidDatetimeFormat(t *testing.T) {
	input := model.LionAirResponse{
		Success: true,
		Data: model.LionAirData{
			AvailableFlights: []model.LionAirFlight{
				{
					ID:      "JT740",
					Carrier: model.LionAirCarrier{Name: "Lion Air", IATA: "JT"},
					Route: model.LionAirRoute{
						From: model.LionAirAirport{
							Code: "CGK",
							Name: "Soekarno-Hatta",
							City: "Jakarta",
						},
						To: model.LionAirAirport{
							Code: "DPS",
							Name: "Ngurah Rai",
							City: "Denpasar",
						},
					},
					Schedule:   model.LionAirSchedule{Departure: "2025-12-15 05:30", DepartureTimezone: "Asia/Jakarta", Arrival: "2025-12-15T08:15:00", ArrivalTimezone: "Asia/Makassar"},
					FlightTime: 105,
					IsDirect:   true,
					Pricing:    model.LionAirPricing{Total: 950000, Currency: "IDR", FareType: "ECONOMY"},
					SeatsLeft:  45,
					PlaneType:  "Boeing 737-900ER",
				},
			},
		},
	}

	flights := lionair.ToFlights(input)
	assert.Empty(t, flights)
}
