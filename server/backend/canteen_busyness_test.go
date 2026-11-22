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

func TestGetCanteenBusynessWithMySQL(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}

	berlin, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)
	// Tuesday 2026-11-17 at 12:05 Berlin (weekday=1, bucket=48 -> 12:00).
	now := time.Date(2026, time.November, 17, 12, 5, 0, 0, berlin)

	require.NoError(t, db.Create(&model.CanteenHeadCount{
		CanteenId: "mensa-garching",
		Count:     750,
		MaxCount:  1000,
		Percent:   75,
		Timestamp: now.Add(-2 * time.Minute),
	}).Error)

	// Two recent Tuesdays + one stale sample older than 56 days (must be excluded).
	history := []model.CanteenHeadCountHistory{
		{CanteenID: "mensa-garching", SampleDate: "2026-11-10", Bucket: 48, Weekday: 1, AvgCount: 600, AvgPercent: 60, SampleCount: 3},
		{CanteenID: "mensa-garching", SampleDate: "2026-11-03", Bucket: 48, Weekday: 1, AvgCount: 500, AvgPercent: 50, SampleCount: 3},
		{CanteenID: "mensa-garching", SampleDate: "2026-08-01", Bucket: 48, Weekday: 1, AvgCount: 100, AvgPercent: 10, SampleCount: 3},
		// Upcoming lunch slots: 12:30 (bucket 50, 80%) and 13:15 (bucket 53, 35% -> best upcoming).
		{CanteenID: "mensa-garching", SampleDate: "2026-11-10", Bucket: 50, Weekday: 1, AvgCount: 800, AvgPercent: 80, SampleCount: 3},
		{CanteenID: "mensa-garching", SampleDate: "2026-11-10", Bucket: 53, Weekday: 1, AvgCount: 350, AvgPercent: 35, SampleCount: 3},
	}
	for _, h := range history {
		require.NoError(t, db.Create(&h).Error)
	}

	reply, err := server.getCanteenBusynessAt(ctx, &pb.GetCanteenBusynessRequest{CanteenId: "MENSA_GARCHING"}, now)
	require.NoError(t, err)
	require.Equal(t, int32(1), reply.Weekday)
	require.Equal(t, float32(75), reply.Current.Percent)
	require.Len(t, reply.TypicalSlots, 3)

	// Bucket 48 average across the two recent Tuesdays = (60 + 50)/2 = 55% (stale 10% excluded).
	require.Equal(t, int32(12), reply.TypicalSlots[0].Hour)
	require.Equal(t, int32(0), reply.TypicalSlots[0].Minute)
	require.InDelta(t, 55.0, reply.TypicalSlots[0].TypicalPercent, 0.01)
	require.Equal(t, int32(2), reply.TypicalSlots[0].SampleDays)

	// Live 75% vs typical 55% = +20%.
	require.True(t, reply.HasTypicalComparison)
	require.InDelta(t, 20.0, reply.VsTypicalPercent, 0.01)

	// Best upcoming lunch slot should be 13:15 (35%).
	require.NotNil(t, reply.BestUpcomingSlot)
	require.Equal(t, int32(13), reply.BestUpcomingSlot.Hour)
	require.Equal(t, int32(15), reply.BestUpcomingSlot.Minute)
	require.InDelta(t, 35.0, reply.BestUpcomingSlot.TypicalPercent, 0.01)

	// Invalid weekday rejected.
	badDay := int32(7)
	_, err = server.getCanteenBusynessAt(ctx, &pb.GetCanteenBusynessRequest{CanteenId: "mensa-garching", Weekday: &badDay}, now)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
