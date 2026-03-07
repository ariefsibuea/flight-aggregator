package garuda

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/pkg/timeutil"
)

func ToFlights(res model.GarudaResponse) []model.Flight {
	flights := make([]model.Flight, 0, len(res.Flights))

	for _, f := range res.Flights {
		flight := model.Flight{
			ID:             fmt.Sprintf("%s_GarudaIndonesia", f.FlightID),
			Provider:       "Garuda Indonesia",
			FlightNumber:   f.FlightID,
			Stops:          f.Stops,
			AvailableSeats: f.AvailableSeats,
			CabinClass:     strings.ToLower(f.FareClass),
		}

		flight.Airline = model.Airline{
			Name: f.Airline,
			Code: f.AirlineCode,
		}

		if f.Aircraft != "" {
			aircraft := f.Aircraft
			flight.Aircraft = &aircraft
		}

		departureDatetime, err := timeutil.ParseDateTime(f.Departure.Time)
		if err != nil {
			slog.Warn("skip flight: cannot parse departure datetime", "flight_id", f.FlightID, "error", err)
			continue
		}
		flight.Departure = model.FlightEndpoint{
			Airport:   f.Departure.Airport,
			City:      f.Departure.City,
			Datetime:  departureDatetime,
			Timestamp: departureDatetime.Unix(),
		}

		arrivalDatetime, err := timeutil.ParseDateTime(f.Arrival.Time)
		if err != nil {
			slog.Warn("skip flight: cannot parse arrival datetime", "flight_id", f.FlightID, "error", err)
			continue
		}
		flight.Arrival = model.FlightEndpoint{
			Airport:   f.Arrival.Airport,
			City:      f.Arrival.City,
			Datetime:  arrivalDatetime,
			Timestamp: arrivalDatetime.Unix(),
		}

		if !departureDatetime.Before(arrivalDatetime) {
			slog.Warn("skip flight: arrival at the same time or before departure", "flight_id", f.FlightID)
			continue
		}

		// Handle flight with multiple segments, use last segment as arrival datetime and recompute duration.
		durationMinutes := f.DurationMinutes
		if len(f.Segments) > 0 {
			lastSegment := f.Segments[len(f.Segments)-1]
			lastArrival, err := timeutil.ParseDateTime(lastSegment.Arrival.Time)
			if err != nil {
				slog.Warn("skip flight: cannot parse last segment arrival datetime", "flight_id", f.FlightID, "error", err)
				continue
			}

			flight.Arrival.Airport = lastSegment.Arrival.Airport
			flight.Arrival.City = lastSegment.Arrival.City
			flight.Arrival.Datetime = lastArrival
			flight.Arrival.Timestamp = lastArrival.Unix()

			totalDuration := 0
			for _, segment := range f.Segments {
				totalDuration += segment.DurationMinutes + segment.LayoverMinutes
			}

			durationMinutes = totalDuration
			flight.Stops = len(f.Segments) - 1
		}

		flight.Duration = model.Duration{
			TotalMinutes: durationMinutes,
			Formatted:    timeutil.FormatDuration(durationMinutes),
		}

		flight.Price = model.Price{
			Amount:   int64(f.Price.Amount),
			Currency: f.Price.Currency,
		}

		// Define empty slice of string to prevent nil value.
		amenities := make([]string, 0, len(f.Amenities))
		amenities = append(amenities, f.Amenities...)
		flight.Amenities = amenities

		flight.Baggage = model.Baggage{
			CarryOn: countBaggage(f.Baggage.CarryOn),
			Checked: countBaggage(f.Baggage.Checked),
		}

		flights = append(flights, flight)
	}

	return flights
}

func countBaggage(count int) string {
	switch count {
	case 0:
		return ""
	case 1:
		return "1 pc"
	default:
		return fmt.Sprintf("%d pcs", count)
	}
}
