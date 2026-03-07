package lionair

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/timeutil"
)

func ToFlights(res model.LionAirResponse) []model.Flight {
	flights := make([]model.Flight, 0, len(res.Data.AvailableFlights))

	for _, f := range res.Data.AvailableFlights {
		flight := model.Flight{
			ID:             fmt.Sprintf("%s_LionAir", f.ID),
			Provider:       "Lion Air",
			FlightNumber:   f.ID,
			Stops:          0,
			AvailableSeats: f.SeatsLeft,
			CabinClass:     strings.ToLower(f.Pricing.FareType),
		}

		flight.Airline = model.Airline{
			Name: f.Carrier.Name,
			Code: f.Carrier.IATA,
		}

		if f.PlaneType != "" {
			aircraft := f.PlaneType
			flight.Aircraft = &aircraft
		}

		departureDatetime, err := timeutil.ParseDateTimeInLocation(f.Schedule.Departure, f.Schedule.DepartureTimezone)
		if err != nil {
			slog.Warn("skip flight: cannot parse departure datetime", "flight_id", f.ID, "error", err)
			continue
		}
		flight.Departure = model.FlightEndpoint{
			Airport:   f.Route.From.Code,
			City:      f.Route.From.City,
			Datetime:  departureDatetime,
			Timestamp: departureDatetime.Unix(),
		}

		arrivalDatetime, err := timeutil.ParseDateTimeInLocation(f.Schedule.Arrival, f.Schedule.ArrivalTimezone)
		if err != nil {
			slog.Warn("skip flight: cannot parse arrival datetime", "flight_id", f.ID, "error", err)
			continue
		}
		flight.Arrival = model.FlightEndpoint{
			Airport:   f.Route.To.Code,
			City:      f.Route.To.City,
			Datetime:  arrivalDatetime,
			Timestamp: arrivalDatetime.Unix(),
		}

		if !departureDatetime.Before(arrivalDatetime) {
			slog.Warn("skip flight: arrival at the same time or before departure", "flight_id", f.ID)
			continue
		}

		flight.Duration = model.Duration{
			TotalMinutes: f.FlightTime,
			Formatted:    timeutil.FormatDuration(f.FlightTime),
		}

		if !f.IsDirect {
			flight.Stops = f.StopCount
		}

		flight.Price = model.Price{
			Amount:   int64(f.Pricing.Total),
			Currency: f.Pricing.Currency,
		}

		amenities := make([]string, 0)
		if f.Services.WifiAvailable {
			amenities = append(amenities, "wifi")
		}
		if f.Services.MealsIncluded {
			amenities = append(amenities, "meal")
		}
		flight.Amenities = amenities

		flight.Baggage = model.Baggage{
			CarryOn: f.Services.BaggageAllowance.Cabin,
			Checked: f.Services.BaggageAllowance.Hold,
		}

		flights = append(flights, flight)
	}

	return flights
}
