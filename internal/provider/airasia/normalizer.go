package airasia

import (
	"fmt"
	"log/slog"
	"math"
	"strings"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/airport"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/strutil"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/timeutil"
)

func ToFlights(res model.AirAsiaResponse) []model.Flight {
	flights := make([]model.Flight, 0, len(res.Flights))

	for _, f := range res.Flights {
		flight := model.Flight{
			ID:             fmt.Sprintf("%s_AirAsia", f.FlightCode),
			Provider:       "AirAsia",
			FlightNumber:   f.FlightCode,
			AvailableSeats: f.Seats,
			CabinClass:     strings.ToLower(f.CabinClass),
		}

		flight.Airline = model.Airline{
			Name: f.Airline,
			Code: strings.ToUpper(getAirlineCode(f.FlightCode)),
		}

		departureDatetime, err := timeutil.ParseDateTime(f.DepartTime)
		if err != nil {
			slog.Warn("skip flight: cannot parse departure datetime", "flight_code", f.FlightCode, "error", err)
			continue
		}
		flight.Departure = model.FlightEndpoint{
			Airport:   f.FromAirport,
			City:      airport.City(f.FromAirport),
			Datetime:  departureDatetime,
			Timestamp: departureDatetime.Unix(),
		}

		arrivalDatetime, err := timeutil.ParseDateTime(f.ArriveTime)
		if err != nil {
			slog.Warn("skip flight: cannot parse arrival datetime", "flight_code", f.FlightCode, "error", err)
			continue
		}
		flight.Arrival = model.FlightEndpoint{
			Airport:   f.ToAirport,
			City:      airport.City(f.ToAirport),
			Datetime:  arrivalDatetime,
			Timestamp: arrivalDatetime.Unix(),
		}

		if !departureDatetime.Before(arrivalDatetime) {
			slog.Warn("skip flight: arrival at the same time or before departure", "flight_code", f.FlightCode)
			continue
		}

		durationMinutes := int(math.Round(f.DurationHours * 60))
		flight.Duration = model.Duration{
			TotalMinutes: durationMinutes,
			Formatted:    timeutil.FormatDuration(durationMinutes),
		}

		if !f.DirectFlight {
			flight.Stops = len(f.Stops)
		}

		flight.Price = model.Price{
			Amount:   int64(f.PriceIDR),
			Currency: "IDR",
		}

		flight.Amenities = make([]string, 0)
		flight.Baggage = toBaggage(f.BaggageNote)

		flights = append(flights, flight)
	}

	return flights
}

// getAirlineCode extracts the airline code from a flight code. IATA flight codes begin with a 2-character
// airline designator followed by a 1-4 digit flight number, e.g.: QZ520 → QZ.
func getAirlineCode(flightCode string) string {
	if len(flightCode) >= 2 {
		return flightCode[:2]
	}
	return flightCode
}

func toBaggage(baggageNote string) model.Baggage {
	baggage := model.Baggage{}
	baggageNote = strings.TrimSpace(baggageNote)

	notes := strings.Split(baggageNote, ",")
	for i, note := range notes {
		note = strings.TrimSpace(note)
		if i == 0 {
			baggage.CarryOn = note
		}
		if i == 1 {
			baggage.Checked = strutil.CapitalizeFirst(note)
		}
	}

	return baggage
}
