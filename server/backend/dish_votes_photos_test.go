package backend

import (
	"context"
	"testing"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestVoteDishRatingAndListDishPhotos(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	prevStorage := cron.StorageDir
	cron.StorageDir = t.TempDir()
	t.Cleanup(func() { cron.StorageDir = prevStorage })

	server := CampusServer{db: db}
	canteen := model.Canteen{Name: "MENSA_GARCHING"}
	require.NoError(t, db.Create(&canteen).Error)
	dish := model.Dish{Name: "Pizza Margherita", Type: "Pizza", CanteenID: canteen.CanteenID}
	require.NoError(t, db.Create(&dish).Error)

	img := readTestImage(t)
	r1, err := server.CreateDishRating(deviceCtx("u1"), &pb.CreateDishRatingRequest{
		Points: 4, CanteenId: "mensa-garching", Dish: "Pizza Margherita", Comment: "Good crust", Image: img,
	})
	require.NoError(t, err)
	r2, err := server.CreateDishRating(deviceCtx("u2"), &pb.CreateDishRatingRequest{
		Points: 5, CanteenId: "mensa-garching", Dish: "Pizza Margherita", Comment: "Best slice", Image: img,
	})
	require.NoError(t, err)
	r3, err := server.CreateDishRating(deviceCtx("u3"), &pb.CreateDishRatingRequest{
		Points: 3, CanteenId: "mensa-garching", Dish: "Pizza Margherita", Comment: "No photo",
	})
	require.NoError(t, err)
	_ = r3

	// Unauthenticated vote rejected.
	_, err = server.VoteDishRating(ctx, &pb.VoteDishRatingRequest{RatingId: r1.RatingId, Helpful: true})
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	// Two devices upvote r1; one device upvotes r2. Duplicate upvote is idempotent.
	v1, err := server.VoteDishRating(deviceCtx("voter-1"), &pb.VoteDishRatingRequest{RatingId: r1.RatingId, Helpful: true})
	require.NoError(t, err)
	require.Equal(t, int64(1), v1.HelpfulCount)
	require.True(t, v1.UserVoted)

	v1Dup, err := server.VoteDishRating(deviceCtx("voter-1"), &pb.VoteDishRatingRequest{RatingId: r1.RatingId, Helpful: true})
	require.NoError(t, err)
	require.Equal(t, int64(1), v1Dup.HelpfulCount)

	v2, err := server.VoteDishRating(deviceCtx("voter-2"), &pb.VoteDishRatingRequest{RatingId: r1.RatingId, Helpful: true})
	require.NoError(t, err)
	require.Equal(t, int64(2), v2.HelpfulCount)

	_, err = server.VoteDishRating(deviceCtx("voter-2"), &pb.VoteDishRatingRequest{RatingId: r2.RatingId, Helpful: true})
	require.NoError(t, err)

	// RATING_SORT_HELPFUL puts r1 (2 votes) ahead of r2 (1 vote) and r3 (0 votes), and sets user_voted per caller.
	sorted, err := server.GetDishRatings(deviceCtx("voter-1"), &pb.GetDishRatingsRequest{
		CanteenId:        "mensa-garching",
		Dish:             "Pizza Margherita",
		Limit:            10,
		SortBy:           pb.RatingSortOrder_RATING_SORT_HELPFUL,
		OmitInlineImages: true,
	})
	require.NoError(t, err)
	require.Len(t, sorted.Rating, 3)
	require.Equal(t, r1.RatingId, sorted.Rating[0].RatingId)
	require.Equal(t, int64(2), sorted.Rating[0].HelpfulCount)
	require.True(t, sorted.Rating[0].UserVoted)
	require.Equal(t, r2.RatingId, sorted.Rating[1].RatingId)
	require.Equal(t, int64(1), sorted.Rating[1].HelpfulCount)
	require.False(t, sorted.Rating[1].UserVoted)

	// ListDishPhotos returns only the 2 ratings with photos, ordered by helpful_count desc.
	gallery, err := server.ListDishPhotos(ctx, &pb.ListDishPhotosRequest{
		CanteenId: "mensa-garching", Dish: "Pizza Margherita", Limit: 1, Offset: 0,
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), gallery.Total)
	require.True(t, gallery.HasMore)
	require.Len(t, gallery.Photos, 1)
	require.Equal(t, r1.RatingId, gallery.Photos[0].RatingId)
	require.Equal(t, dishRatingImageURL(r1.RatingId), gallery.Photos[0].ImageUrl)
	require.Equal(t, int64(2), gallery.Photos[0].HelpfulCount)

	// Removing vote decrements helpful_count.
	unvoted, err := server.VoteDishRating(deviceCtx("voter-1"), &pb.VoteDishRatingRequest{RatingId: r1.RatingId, Helpful: false})
	require.NoError(t, err)
	require.Equal(t, int64(1), unvoted.HelpfulCount)
	require.False(t, unvoted.UserVoted)
}
