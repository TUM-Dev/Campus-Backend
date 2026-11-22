package cron

import (
	"context"
	"testing"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/stretchr/testify/require"
)

func TestLibraryOccupancyUpdateAndHistoryPruning(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	berlin, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)
	now := time.Date(2026, time.November, 18, 14, 5, 0, 0, berlin)

	lib := TrackedLibraries[0]
	require.NoError(t, updateLibraryOccupancy(db, &lib, 325, now))
	require.NoError(t, recordLibraryOccupancyHistory(db, lib.LibraryID, 325, 50, now))
	require.NoError(t, recordLibraryOccupancyHistory(db, lib.LibraryID, 455, 70, now.Add(5*time.Minute)))
	require.NoError(t, recordLibraryOccupancyHistory(db, lib.LibraryID, 100, 15, now.AddDate(0, 0, -190)))

	require.NoError(t, pruneOldLibraryOccupancyHistory(db, now))

	var live model.LibraryOccupancy
	require.NoError(t, db.First(&live, "library_id = ?", lib.LibraryID).Error)
	require.Equal(t, uint32(325), live.Count)
	require.InDelta(t, 50.0, float64(live.Percent), 0.1)

	var hist []model.LibraryOccupancyHistory
	require.NoError(t, db.Find(&hist, "library_id = ?", lib.LibraryID).Error)
	require.Len(t, hist, 1)
	require.Equal(t, int32(2), hist[0].SampleCount)
	require.InDelta(t, 60.0, hist[0].AvgPercent, 0.01)
}
