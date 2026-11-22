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

func TestReportDishSoldOutFlow(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}

	canteen := model.Canteen{Name: "MENSA_GARCHING"}
	require.NoError(t, db.Create(&canteen).Error)
	dish := model.Dish{Name: "Currywurst", Type: "Fleisch", CanteenID: canteen.CanteenID}
	require.NoError(t, db.Create(&dish).Error)

	// Unauthenticated request rejected.
	_, err := server.ReportDishSoldOut(ctx, &pb.ReportDishSoldOutRequest{CanteenId: "mensa-garching", Dish: "Currywurst"})
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	// First two reports do not yet confirm sold-out; duplicate from same device is idempotent.
	r1, err := server.ReportDishSoldOut(deviceCtx("d1"), &pb.ReportDishSoldOutRequest{CanteenId: "mensa-garching", Dish: "Currywurst"})
	require.NoError(t, err)
	require.Equal(t, int32(1), r1.ReportCount)
	require.False(t, r1.Confirmed)
	require.False(t, r1.AlreadyReported)

	r1Dup, err := server.ReportDishSoldOut(deviceCtx("d1"), &pb.ReportDishSoldOutRequest{CanteenId: "mensa-garching", Dish: "Currywurst"})
	require.NoError(t, err)
	require.Equal(t, int32(1), r1Dup.ReportCount)
	require.True(t, r1Dup.AlreadyReported)

	r2, err := server.ReportDishSoldOut(deviceCtx("d2"), &pb.ReportDishSoldOutRequest{CanteenId: "mensa-garching", Dish: "Currywurst"})
	require.NoError(t, err)
	require.Equal(t, int32(2), r2.ReportCount)
	require.False(t, r2.Confirmed)

	// Third distinct device confirms sold-out today.
	r3, err := server.ReportDishSoldOut(deviceCtx("d3"), &pb.ReportDishSoldOutRequest{CanteenId: "mensa-garching", Dish: "Currywurst"})
	require.NoError(t, err)
	require.Equal(t, int32(3), r3.ReportCount)
	require.True(t, r3.Confirmed)

	// ListDishRatingSummaries & GetDishRatings surface sold_out=true today.
	summary, err := server.ListDishRatingSummaries(ctx, &pb.ListDishRatingSummariesRequest{
		CanteenId: "mensa-garching", Dishes: []string{"Currywurst"},
	})
	require.NoError(t, err)
	require.Len(t, summary.Summaries, 1)
	require.True(t, summary.Summaries[0].SoldOut)
	require.Equal(t, int32(3), summary.Summaries[0].SoldOutReportCount)

	// Next day in Europe/Berlin starts fresh at 0 reports.
	tomorrow := time.Now().AddDate(0, 0, 1)
	rTomorrow, err := server.reportDishSoldOutAt(deviceCtx("d1"), &pb.ReportDishSoldOutRequest{CanteenId: "mensa-garching", Dish: "Currywurst"}, tomorrow)
	require.NoError(t, err)
	require.Equal(t, int32(1), rTomorrow.ReportCount)
	require.False(t, rTomorrow.Confirmed)
	require.False(t, rTomorrow.AlreadyReported)
}
