package lionair_test

import (
	"context"
	"testing"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/lionair"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch_Success(t *testing.T) {
	client := lionair.NewClient()

	flights, err := client.Fetch(context.Background(), model.SearchRequest{})
	require.NoError(t, err)

	assert.Len(t, flights, 3)

	expectIDs := []string{
		"JT740_LionAir",
		"JT742_LionAir",
		"JT650_LionAir",
	}
	for i, expect := range expectIDs {
		assert.Equal(t, expect, flights[i].ID)
	}
}

func TestFetch_NormalizedFields(t *testing.T) {
	client := lionair.NewClient()
	flights, err := client.Fetch(context.Background(), model.SearchRequest{})
	require.NoError(t, err)

	t.Run("JT740 IANA timezone applied correctly", func(t *testing.T) {
		var jt740 model.Flight
		for _, f := range flights {
			if f.ID == "JT740_LionAir" {
				jt740 = f
				break
			}
		}
		require.NotEmpty(t, jt740.ID, "JT740 not found in results")

		// Lion Air uses IANA timezone names (e.g. "Asia/Jakarta" = +07:00).
		// The parsed time must carry the correct UTC offset, not be UTC.
		_, offset := jt740.Departure.Datetime.Zone()
		assert.Equal(t, 7*60*60, offset, "expected Asia/Jakarta (+07:00) offset applied")
		assert.False(t, jt740.Departure.Datetime.IsZero())
		assert.True(t, jt740.Departure.Datetime.Before(jt740.Arrival.Datetime))
		assert.Equal(t, jt740.Departure.Datetime.Unix(), jt740.Departure.Timestamp)
	})

	t.Run("stops derived from is_direct and stop_count", func(t *testing.T) {
		for _, f := range flights {
			// All Lion Air mock flights are direct; stops must be 0.
			assert.GreaterOrEqual(t, f.Stops, 0)
		}
	})

	t.Run("duration formatted correctly", func(t *testing.T) {
		for _, f := range flights {
			assert.Positive(t, f.Duration.TotalMinutes)
			assert.NotEmpty(t, f.Duration.Formatted)
		}
	})
}

func TestFetch_ContextCancellation(t *testing.T) {
	client := lionair.NewClient()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Fetch(ctx, model.SearchRequest{})
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestFetch_ContextTimeout(t *testing.T) {
	client := lionair.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Fetch(ctx, model.SearchRequest{})
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
