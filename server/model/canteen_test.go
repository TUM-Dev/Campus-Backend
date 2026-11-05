package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanteenIdentifierConversion(t *testing.T) {
	for _, in := range []string{"mensa-garching", "MENSA_GARCHING", " Mensa-Garching "} {
		require.Equal(t, "MENSA_GARCHING", CanteenEnumName(in), in)
		require.Equal(t, "mensa-garching", CanteenEatAPIID(in), in)
	}
	require.Equal(t, "STUCAFE_AKADEMIE_WEIHENSTEPHAN", CanteenEnumName("stucafe-akademie-weihenstephan"))
}
