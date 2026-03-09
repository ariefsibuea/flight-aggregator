package usecase

import (
	"cmp"
	"math"
	"slices"
	"strings"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
)

func filterAndSort(flights []model.Flight, req model.SearchRequest) []model.Flight {
	res := flights

	if req.Filter != nil {
		res = filterByPriceRange(res, req.Filter.MinPrice, req.Filter.MaxPrice)
		res = filterByStops(res, req.Filter.MaxStops)
		res = filterByDepartureTimes(res, req.Filter.DepartureTimes)
		res = filterByArrivalTimes(res, req.Filter.ArrivalTimes)
		res = filterByAirlines(res, req.Filter.Airlines)
		res = filterByDurationRange(res, req.Filter.MinDuration, req.Filter.MaxDuration)
	}

	res = calculateScore(res)

	if req.Sort != nil {
		res = sortBy(res, req.Sort.Field, req.Sort.Direction)
	}

	return res
}

func filterByPriceRange(flights []model.Flight, minPrice, maxPrice *int64) []model.Flight {
	if minPrice == nil && maxPrice == nil {
		return flights
	}

	filteredFlights := make([]model.Flight, 0, len(flights))
	for _, f := range flights {
		if minPrice != nil && f.Price.Amount < *minPrice {
			continue
		}
		if maxPrice != nil && f.Price.Amount > *maxPrice {
			continue
		}
		filteredFlights = append(filteredFlights, f)
	}

	return filteredFlights
}

func filterByStops(flights []model.Flight, maxStops *int) []model.Flight {
	if maxStops == nil {
		return flights
	}

	filteredFlights := make([]model.Flight, 0, len(flights))
	for _, f := range flights {
		if f.Stops > *maxStops {
			continue
		}
		filteredFlights = append(filteredFlights, f)
	}

	return filteredFlights
}

func filterByDepartureTimes(flights []model.Flight, times []model.TimeWindow) []model.Flight {
	if len(times) == 0 {
		return flights
	}

	filteredFlights := make([]model.Flight, 0, len(flights))
	for _, f := range flights {
		hour := f.Departure.Datetime.Hour()
		if model.MatchTimeWindows(hour, times) {
			filteredFlights = append(filteredFlights, f)
		}
	}

	return filteredFlights
}

func filterByArrivalTimes(flights []model.Flight, times []model.TimeWindow) []model.Flight {
	if len(times) == 0 {
		return flights
	}

	filteredFlights := make([]model.Flight, 0, len(flights))
	for _, f := range flights {
		hour := f.Arrival.Datetime.Hour()
		if model.MatchTimeWindows(hour, times) {
			filteredFlights = append(filteredFlights, f)
		}
	}

	return filteredFlights
}

func filterByAirlines(flights []model.Flight, airlines []string) []model.Flight {
	if len(airlines) == 0 {
		return flights
	}

	allowedAirlines := make(map[string]bool, len(airlines))
	for _, airline := range airlines {
		allowedAirlines[strings.ToLower(airline)] = true
	}

	filteredFlights := make([]model.Flight, 0, len(flights))
	for _, f := range flights {
		if allowedAirlines[strings.ToLower(f.Airline.Code)] {
			filteredFlights = append(filteredFlights, f)
		}
	}

	return filteredFlights
}

func filterByDurationRange(flights []model.Flight, shortest, longest *int) []model.Flight {
	if shortest == nil && longest == nil {
		return flights
	}

	filteredFlights := make([]model.Flight, 0, len(flights))
	for _, f := range flights {
		if shortest != nil && *shortest > f.Duration.TotalMinutes {
			continue
		}
		if longest != nil && *longest < f.Duration.TotalMinutes {
			continue
		}
		filteredFlights = append(filteredFlights, f)
	}

	return filteredFlights
}

func sortBy(flights []model.Flight, field, direction string) []model.Flight {
	if field == "" {
		return flights
	}

	compare := func(i, j model.Flight) int {
		var result int

		switch field {
		case "price":
			result = cmp.Compare(i.Price.Amount, j.Price.Amount)
		case "duration":
			result = cmp.Compare(i.Duration.TotalMinutes, j.Duration.TotalMinutes)
		case "departure":
			result = cmp.Compare(i.Departure.Timestamp, j.Departure.Timestamp)
		case "arrival":
			result = cmp.Compare(i.Arrival.Timestamp, j.Arrival.Timestamp)
		case "score":
			result = cmp.Compare(i.BestValueScore, j.BestValueScore)
		default:
			return 0
		}

		if direction == "desc" {
			return -result
		}
		return result
	}

	slices.SortFunc(flights, compare)
	return flights
}

func calculateScore(flights []model.Flight) []model.Flight {
	if len(flights) == 0 {
		return flights
	}

	minPrice, maxPrice := minMaxPrice(flights)
	minDuration, maxDuration := minMaxDuration(flights)

	priceRange := maxPrice - minPrice
	durationRange := maxDuration - minDuration

	for i := range flights {
		priceScore := 1.0
		if priceRange > 0 {
			priceScore = 1.0 - float64(flights[i].Price.Amount-minPrice)/float64(priceRange)
		}

		durationScore := 1.0
		if durationRange > 0 {
			durationScore = 1.0 - float64(flights[i].Duration.TotalMinutes-minDuration)/float64(durationRange)
		}

		stopsScore := 1.0
		switch flights[i].Stops {
		case 0:
			stopsScore = 1.0
		case 1:
			stopsScore = 0.5
		default:
			stopsScore = 0.0
		}

		amenityScore := calculateAmenityScore(flights[i].Amenities)

		flights[i].BestValueScore = (priceScore * 0.50) + (durationScore * 0.25) + (stopsScore * 0.15) + (amenityScore * 0.10)
	}

	return flights
}

// minMaxPrice returns the minimum and the maximum price of provided flights.
func minMaxPrice(flights []model.Flight) (int64, int64) {
	minPrice := flights[0].Price.Amount
	maxPrice := flights[0].Price.Amount
	for _, f := range flights[1:] {
		minPrice = min(minPrice, f.Price.Amount)
		maxPrice = max(maxPrice, f.Price.Amount)
	}
	return minPrice, maxPrice
}

// minMaxDuration returns the minimum and the maximum total duration of provided flights.
func minMaxDuration(flights []model.Flight) (int, int) {
	minDuration := flights[0].Duration.TotalMinutes
	maxDuration := flights[0].Duration.TotalMinutes
	for _, f := range flights[1:] {
		minDuration = min(minDuration, f.Duration.TotalMinutes)
		maxDuration = max(maxDuration, f.Duration.TotalMinutes)
	}
	return minDuration, maxDuration
}

func calculateAmenityScore(amenities []string) float64 {
	score := 0.0

	existedAmenities := make(map[string]bool)
	for _, a := range amenities {
		existedAmenities[strings.ToLower(a)] = true
	}

	if _, exist := existedAmenities["wifi"]; exist {
		score += 0.3
	}
	if _, exist := existedAmenities["meal"]; exist {
		score += 0.4
	}
	if _, exist := existedAmenities["entertainment"]; exist {
		score += 0.3
	}

	return math.Min(score, 1.0)
}
