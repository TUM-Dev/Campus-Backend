package backend

import (
	"context"
	"fmt"
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

	// Duplicate report from same device is idempotent; reports below threshold do not confirm sold-out.
	for i := int32(1); i <= soldOutConfirmThreshold; i++ {
		dCtx := deviceCtx(fmt.Sprintf("d%d", i))
		rep, err := server.ReportDishSoldOut(dCtx, &pb.ReportDishSoldOutRequest{CanteenId: "mensa-garching", Dish: "Currywurst"})
		require.NoError(t, err)
		require.Equal(t, i, rep.ReportCount)
		require.False(t, rep.AlreadyReported)
		if i < soldOutConfirmThreshold {
			require.False(t, rep.Confirmed)
		} else {
			require.True(t, rep.Confirmed)
		}
		if i == 1 {
			dup, err := server.ReportDishSoldOut(dCtx, &pb.ReportDishSoldOutRequest{CanteenId: "mensa-garching", Dish: "Currywurst"})
			require.NoError(t, err)
			require.Equal(t, int32(1), dup.ReportCount)
			require.True(t, dup.AlreadyReported)
		}
	}

	// ListDishRatingSummaries & GetDishRatings (unrated fallback) surface sold_out=true today.
	summary, err := server.ListDishRatingSummaries(ctx, &pb.ListDishRatingSummariesRequest{
		CanteenId: "mensa-garching", Dishes: []string{"Currywurst"},
	})
	require.NoError(t, err)
	require.Len(t, summary.Summaries, 1)
	require.True(t, summary.Summaries[0].SoldOut)
	require.Equal(t, soldOutConfirmThreshold, summary.Summaries[0].SoldOutReportCount)

	fallback, err := server.GetDishRatings(ctx, &pb.GetDishRatingsRequest{
		CanteenId: "mensa-garching", Dish: "Currywurst", IncludeUnratedFallback: true,
	})
	require.NoError(t, err)
	require.True(t, fallback.SoldOut)
	require.Equal(t, soldOutConfirmThreshold, fallback.SoldOutReportCount)

	// Next day in Europe/Berlin starts fresh at 0 reports.
	tomorrow := time.Now().AddDate(0, 0, 1)
	rTomorrow, err := server.reportDishSoldOutAt(deviceCtx("d1"), &pb.ReportDishSoldOutRequest{CanteenId: "mensa-garching", Dish: "Currywurst"}, tomorrow)
	require.NoError(t, err)
	require.Equal(t, int32(1), rTomorrow.ReportCount)
	require.False(t, rTomorrow.Confirmed)
	require.False(t, rTomorrow.AlreadyReported)
}
