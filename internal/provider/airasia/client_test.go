package airasia_test

import (
	"context"
	"testing"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/airasia"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch_Success(t *testing.T) {
	client := airasia.NewClient()

	for i := 0; i < 20; i++ {
		flights, err := client.Fetch(context.Background(), model.SearchRequest{})
		if err == nil {
			assert.Len(t, flights, 4)

			expectIDs := []string{
				"QZ520_AirAsia",
				"QZ524_AirAsia",
				"QZ532_AirAsia",
				"QZ7250_AirAsia",
			}
			for i, expect := range expectIDs {
				assert.Equal(t, expect, flights[i].ID)
			}

			return
		}
	}

	t.Skip("AirAsia has 10% random failure rate, could not get successful response after 20 retries")
}

func TestFetch_NormalizedFields(t *testing.T) {
	client := airasia.NewClient()

	var flights []model.Flight
	for i := 0; i < 20; i++ {
		result, err := client.Fetch(context.Background(), model.SearchRequest{})
		if err == nil {
			flights = result
			break
		}
	}
	if len(flights) == 0 {
		t.Skip("AirAsia has 10% random failure rate, could not get successful response after 20 retries")
	}

	t.Run("duration converted from float hours to minutes", func(t *testing.T) {
		for _, f := range flights {
			// duration_hours is a float64 (e.g. 1.67); rounded to int minutes (e.g. 100).
			assert.Positive(t, f.Duration.TotalMinutes, "flight %s has zero duration", f.ID)
			assert.Regexp(t, `^\d+h \d+m$`, f.Duration.Formatted, "flight %s formatted duration %q is not canonical", f.ID, f.Duration.Formatted)
		}
	})

	t.Run("baggage carry-on set from baggage_note", func(t *testing.T) {
		for _, f := range flights {
			assert.NotEmpty(t, f.Baggage.CarryOn, "flight %s missing carry-on", f.ID)
		}
	})

	t.Run("stops derived from direct_flight field", func(t *testing.T) {
		// QZ520, QZ524, QZ532 are direct; QZ7250 has 1 stop.
		for _, f := range flights {
			assert.GreaterOrEqual(t, f.Stops, 0, "flight %s has negative stops", f.ID)
		}

		var qz7250 model.Flight
		for _, f := range flights {
			if f.ID == "QZ7250_AirAsia" {
				qz7250 = f
				break
			}
		}
		require.NotEmpty(t, qz7250.ID, "QZ7250 not found")
		assert.Equal(t, 1, qz7250.Stops, "QZ7250 has 1 stop")
	})
}

func TestFetch_ContextCancellation(t *testing.T) {
	client := airasia.NewClient()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Fetch(ctx, model.SearchRequest{})
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestFetch_ContextTimeout(t *testing.T) {
	client := airasia.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Fetch(ctx, model.SearchRequest{})
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
