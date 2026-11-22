package backend

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRatingImagesDisabledByDefault(t *testing.T) {
	t.Setenv(ratingImagesEnabledEnv, "")
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}

	canteen := model.Canteen{Name: "MENSA_GARCHING"}
	require.NoError(t, db.Create(&canteen).Error)
	dish := model.Dish{Name: "Pasta", Type: "Pasta", CanteenID: canteen.CanteenID}
	require.NoError(t, db.Create(&dish).Error)

	img := readTestImage(t)
	_, err := server.CreateDishRating(deviceCtx("dev-img"), &pb.CreateDishRatingRequest{
		Points: 4, CanteenId: "mensa-garching", Dish: "Pasta", Image: img,
	})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))

	_, err = server.CreateCanteenRating(deviceCtx("dev-img"), &pb.CreateCanteenRatingRequest{
		Points: 4, CanteenId: "mensa-garching", Image: img,
	})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))

	// Text-only ratings still succeed when images are disabled.
	created, err := server.CreateDishRating(deviceCtx("dev-img"), &pb.CreateDishRatingRequest{
		Points: 4, CanteenId: "mensa-garching", Dish: "Pasta", Comment: "Tasty",
	})
	require.NoError(t, err)

	// Trying to attach an image via UpdateDishRating is also rejected.
	_, err = server.UpdateDishRating(deviceCtx("dev-img"), &pb.UpdateDishRatingRequest{
		RatingId: created.RatingId, Points: 5, Image: img,
	})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestModerationServiceBansStatusAndDiscordWebhook(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)

	receivedCh := make(chan discordWebhookPayload, 2)
	webhookSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload discordWebhookPayload
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		receivedCh <- payload
		w.WriteHeader(http.StatusNoContent)
	}))
	defer webhookSrv.Close()

	server := CampusServer{
		db:                   db,
		moderationWebhookURL: webhookSrv.URL,
		moderationHTTPClient: webhookSrv.Client(),
	}

	canteen := model.Canteen{Name: "MENSA_GARCHING"}
	require.NoError(t, db.Create(&canteen).Error)
	dish := model.Dish{Name: "Burger", Type: "Fleisch", CanteenID: canteen.CanteenID}
	require.NoError(t, db.Create(&dish).Error)

	badActorCtx := deviceCtx("bad-actor")
	created, err := server.CreateDishRating(badActorCtx, &pb.CreateDishRatingRequest{
		Points: 1, CanteenId: "mensa-garching", Dish: "Burger", Comment: "Offensive spam",
	})
	require.NoError(t, err)

	// Reporting triggers an async text-only Discord webhook embed.
	_, err = server.ReportDishRating(deviceCtx("reporter-1"), &pb.ReportDishRatingRequest{
		RatingId: created.RatingId, Reason: "abusive language",
	})
	require.NoError(t, err)

	select {
	case payload := <-receivedCh:
		require.Len(t, payload.Embeds, 1)
		require.Contains(t, payload.Embeds[0].Title, "reported")
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for Discord moderation webhook")
	}

	// ListReportedDishRatings surfaces the reported rating for the /admin/ dashboard.
	reported, err := server.ListReportedDishRatings(ctx, 20, 0)
	require.NoError(t, err)
	require.Len(t, reported, 1)
	require.Equal(t, created.RatingId, reported[0].RatingID)
	require.Equal(t, int64(1), reported[0].ReportCount)
	require.Contains(t, reported[0].Reasons, "abusive language")

	// SetDishRatingStatus hides and restores ratings cleanly.
	require.NoError(t, server.SetDishRatingStatus(ctx, created.RatingId, model.RatingStatusHidden))
	_, err = server.GetDishRatings(ctx, &pb.GetDishRatingsRequest{CanteenId: "mensa-garching", Dish: "Burger"})
	require.Equal(t, codes.NotFound, status.Code(err))

	require.NoError(t, server.SetDishRatingStatus(ctx, created.RatingId, model.RatingStatusVisible))
	restored, err := server.GetDishRatings(ctx, &pb.GetDishRatingsRequest{CanteenId: "mensa-garching", Dish: "Burger", Limit: 10})
	require.NoError(t, err)
	require.Len(t, restored.Rating, 1)

	// BanRatingDevice hides existing ratings and blocks future submissions across all endpoints.
	badHash := callerDeviceHash(badActorCtx).String
	require.NoError(t, server.BanRatingDevice(ctx, badHash, "repeated spam", true))
	_, err = server.GetDishRatings(ctx, &pb.GetDishRatingsRequest{CanteenId: "mensa-garching", Dish: "Burger"})
	require.Equal(t, codes.NotFound, status.Code(err))

	_, err = server.CreateDishRating(badActorCtx, &pb.CreateDishRatingRequest{
		Points: 1, CanteenId: "mensa-garching", Dish: "Burger",
	})
	require.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = server.CreateCanteenRating(badActorCtx, &pb.CreateCanteenRatingRequest{
		Points: 1, CanteenId: "mensa-garching",
	})
	require.Equal(t, codes.PermissionDenied, status.Code(err))

	// Unbanning restores device access.
	require.NoError(t, server.UnbanRatingDevice(ctx, badHash))
	_, err = server.CreateCanteenRating(badActorCtx, &pb.CreateCanteenRatingRequest{
		Points: 4, CanteenId: "mensa-garching",
	})
	require.NoError(t, err)
}
