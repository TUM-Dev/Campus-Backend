package backend

import (
	"context"
	"testing"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListLibrariesAndGetLibraryBusyness(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}

	berlin, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)
	// Wednesday 2026-11-18 14:05 Berlin (weekday=2, bucket=56 -> 14:00).
	now := time.Date(2026, time.November, 18, 14, 5, 0, 0, berlin)

	require.NoError(t, db.Create(&model.LibraryOccupancy{
		LibraryID: "tb-maschinenwesen",
		Name:      "Teilbibliothek Maschinenwesen",
		Campus:    "garching",
		Count:     176,
		MaxCount:  220,
		Percent:   80,
		Timestamp: now.Add(-3 * time.Minute),
	}).Error)

	history := []model.LibraryOccupancyHistory{
		{LibraryID: "tb-maschinenwesen", SampleDate: "2026-11-11", Bucket: 56, Weekday: 2, AvgCount: 132, AvgPercent: 60, SampleCount: 3},
		{LibraryID: "tb-maschinenwesen", SampleDate: "2026-11-04", Bucket: 56, Weekday: 2, AvgCount: 0, AvgPercent: 0, SampleCount: 3},    // closed day -> excluded
		{LibraryID: "tb-maschinenwesen", SampleDate: "2026-11-18", Bucket: 56, Weekday: 2, AvgCount: 198, AvgPercent: 90, SampleCount: 1}, // today -> excluded
		{LibraryID: "tb-maschinenwesen", SampleDate: "2026-11-11", Bucket: 64, Weekday: 2, AvgCount: 88, AvgPercent: 40, SampleCount: 3},  // 16:00
		{LibraryID: "tb-maschinenwesen", SampleDate: "2026-11-11", Bucket: 72, Weekday: 2, AvgCount: 44, AvgPercent: 20, SampleCount: 3},  // 18:00 (best)
	}
	for _, h := range history {
		require.NoError(t, db.Create(&h).Error)
	}

	// ListLibraries filtered by campus="garching" returns all Garching libraries with live data merged.
	list, err := server.ListLibraries(ctx, &pb.ListLibrariesRequest{Campus: "garching"})
	require.NoError(t, err)
	require.Len(t, list.Libraries, 3)
	var mwLib *pb.LibraryOccupancy
	for _, l := range list.Libraries {
		require.Equal(t, "garching", l.Campus)
		if l.LibraryId == "tb-maschinenwesen" {
			mwLib = l
		}
	}
	require.NotNil(t, mwLib)
	require.Equal(t, float32(80), mwLib.Percent)

	// GetLibraryBusyness returns live vs typical (+20%) and best upcoming study slot (18:00, 20%).
	busy, err := server.getLibraryBusynessAt(ctx, &pb.GetLibraryBusynessRequest{LibraryId: "tb-maschinenwesen"}, now)
	require.NoError(t, err)
	require.Equal(t, int32(2), busy.Weekday)
	require.True(t, busy.HasTypicalComparison)
	require.InDelta(t, 20.0, busy.VsTypicalPercent, 0.01)
	require.NotNil(t, busy.BestUpcomingSlot)
	require.Equal(t, int32(18), busy.BestUpcomingSlot.Hour)
	require.Equal(t, int32(0), busy.BestUpcomingSlot.Minute)

	// Unknown library returns NotFound.
	_, err = server.getLibraryBusynessAt(ctx, &pb.GetLibraryBusynessRequest{LibraryId: "does-not-exist"}, now)
	require.Equal(t, codes.NotFound, status.Code(err))
}
