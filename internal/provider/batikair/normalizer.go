package batikair

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/airport"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/strutil"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/timeutil"
)

func ToFlights(res model.BatikAirResponse) []model.Flight {
	flights := make([]model.Flight, 0, len(res.Results))

	for _, f := range res.Results {
		flight := model.Flight{
			ID:             fmt.Sprintf("%s_BatikAir", f.FlightNumber),
			Provider:       "Batik Air",
			FlightNumber:   f.FlightNumber,
			Stops:          f.NumberOfStops,
			AvailableSeats: f.SeatsAvailable,
			CabinClass:     toCabinClass(f.Fare.Class),
		}

		flight.Airline = model.Airline{
			Name: f.AirlineName,
			Code: strings.ToUpper(f.AirlineIATA),
		}

		if f.AircraftModel != "" {
			aircraft := f.AircraftModel
			flight.Aircraft = &aircraft
		}

		departureDatetime, err := timeutil.ParseDateTime(f.DepartureDateTime)
		if err != nil {
			slog.Warn("skip flight: cannot parse departure datetime", "flightNumber", f.FlightNumber, "error", err)
			continue
		}
		flight.Departure = model.FlightEndpoint{
			Airport:   f.Origin,
			City:      airport.City(f.Origin),
			Datetime:  departureDatetime,
			Timestamp: departureDatetime.Unix(),
		}

		arrivalDatetime, err := timeutil.ParseDateTime(f.ArrivalDateTime)
		if err != nil {
			slog.Warn("skip flight: cannot parse arrival datetime", "flightNumber", f.FlightNumber, "error", err)
			continue
		}
		flight.Arrival = model.FlightEndpoint{
			Airport:   f.Destination,
			City:      airport.City(f.Destination),
			Datetime:  arrivalDatetime,
			Timestamp: arrivalDatetime.Unix(),
		}

		if !departureDatetime.Before(arrivalDatetime) {
			slog.Warn("skip flight: arrival at the same time or before departure", "flightNumber", f.FlightNumber)
			continue
		}

		durationMinutes, err := timeutil.ParseDuration(f.TravelTime)
		if err != nil {
			slog.Warn("skip flight: cannot parse travel time", "flightNumber", f.FlightNumber, "travelTime", f.TravelTime, "error", err)
			continue
		}
		flight.Duration = model.Duration{
			TotalMinutes: durationMinutes,
			Formatted:    timeutil.FormatDuration(durationMinutes),
		}

		flight.Price = model.Price{
			Amount:   int64(f.Fare.TotalPrice),
			Currency: f.Fare.CurrencyCode,
		}

		flight.Amenities = f.OnboardServices
		flight.Baggage = toBaggage(f.BaggageInfo)

		flights = append(flights, flight)
	}

	return flights
}

func toCabinClass(fareClass string) string {
	switch fareClass {
	case "J", "C", "D", "I", "Z":
		return "business"
	default:
		return "economy"
	}
}

func toBaggage(baggageInfo string) model.Baggage {
	baggage := model.Baggage{}
	baggageInfo = strings.TrimSpace(baggageInfo)

	notes := strings.Split(baggageInfo, ",")
	for i, note := range notes {
		note = strings.TrimSpace(note)
		if i == 0 {
			baggage.CarryOn = strutil.CapitalizeFirst(note)
		}
		if i == 1 {
			baggage.Checked = strutil.CapitalizeFirst(note)
		}
	}

	return baggage
}
