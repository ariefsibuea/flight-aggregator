package garuda_test

import (
	"context"
	"testing"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/garuda"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch_NormalizedResponse(t *testing.T) {
	client := garuda.NewClient()

	flights, err := client.Fetch(context.Background(), model.SearchRequest{})
	require.NoError(t, err)

	assert.Len(t, flights, 3)

	t.Run("multi-segment-flight", func(t *testing.T) {
		flight := flightByID(flights, "GA315_GarudaIndonesia")
		require.NotNil(t, flight)

		assert.Equal(t, "DPS", flight.Arrival.Airport)
		assert.Equal(t, 1, flight.Stops)
		assert.Equal(t, 285, flight.Duration.TotalMinutes) // 90 (CGK→SUB) + 105 layover + 90 (SUB→DPS)
		assert.Equal(t, "4h 45m", flight.Duration.Formatted)
	})

	t.Run("valid-timezone-offset", func(t *testing.T) {
		flight := flightByID(flights, "GA400_GarudaIndonesia")
		require.NotNil(t, flight)

		assert.False(t, flight.Departure.Datetime.IsZero())
		assert.Equal(t, "2025-12-15T06:00:00+07:00", flight.Departure.Datetime.Format(time.RFC3339))

		_, offset := flight.Departure.Datetime.Zone()
		assert.Equal(t, 7*60*60, offset)
		assert.True(t, flight.Departure.Datetime.Before(flight.Arrival.Datetime))
		assert.Equal(t, flight.Departure.Datetime.Unix(), flight.Departure.Timestamp)
	})

	t.Run("baggage-is-formatted-as-string", func(t *testing.T) {
		flight := flightByID(flights, "GA400_GarudaIndonesia")
		require.NotNil(t, flight)

		assert.NotEmpty(t, flight.Baggage.CarryOn)
		assert.NotEmpty(t, flight.Baggage.Checked)
	})
}

func TestFetch_ContextCancellation(t *testing.T) {
	client := garuda.NewClient()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Fetch(ctx, model.SearchRequest{})
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestFetch_ContextTimeout(t *testing.T) {
	client := garuda.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Fetch(ctx, model.SearchRequest{})
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

// flightByID is a test helper that finds a flight by ID in a slice.
func flightByID(flights []model.Flight, id string) *model.Flight {
	for _, f := range flights {
		if f.ID == id {
			return &f
		}
	}
	return nil
}
