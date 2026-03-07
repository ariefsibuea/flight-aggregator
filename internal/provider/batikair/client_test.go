package batikair_test

import (
	"context"
	"testing"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/batikair"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch_Success(t *testing.T) {
	client := batikair.NewClient()

	flights, err := client.Fetch(context.Background(), model.SearchRequest{})
	require.NoError(t, err)

	assert.Len(t, flights, 3)

	expectIDs := []string{
		"ID6514_BatikAir",
		"ID6520_BatikAir",
		"ID7042_BatikAir",
	}
	for i, expect := range expectIDs {
		assert.Equal(t, expect, flights[i].ID)
	}
}

func TestFetch_NormalizedFields(t *testing.T) {
	client := batikair.NewClient()
	flights, err := client.Fetch(context.Background(), model.SearchRequest{})
	require.NoError(t, err)

	t.Run("duration parsed from string and formatted canonically", func(t *testing.T) {
		for _, f := range flights {
			assert.Positive(t, f.Duration.TotalMinutes, "flight %s has zero duration", f.ID)
			// Formatted must be canonical "Xh Ym", not the raw provider string.
			assert.Regexp(t, `^\d+h \d+m$`, f.Duration.Formatted, "flight %s formatted duration %q is not canonical", f.ID, f.Duration.Formatted)
		}
	})

	t.Run("baggage parsed from comma-separated string", func(t *testing.T) {
		for _, f := range flights {
			assert.NotEmpty(t, f.Baggage.CarryOn, "flight %s missing carry-on baggage", f.ID)
		}
	})

	t.Run("datetime has timezone offset from +0700 format", func(t *testing.T) {
		for _, f := range flights {
			_, offset := f.Departure.Datetime.Zone()
			assert.Equal(t, 7*60*60, offset, "flight %s expected +07:00 offset", f.ID)
			assert.True(t, f.Departure.Datetime.Before(f.Arrival.Datetime), "flight %s arrival before departure", f.ID)
		}
	})
}

func TestFetch_ContextCancellation(t *testing.T) {
	client := batikair.NewClient()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Fetch(ctx, model.SearchRequest{})
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestFetch_ContextTimeout(t *testing.T) {
	client := batikair.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Fetch(ctx, model.SearchRequest{})
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
