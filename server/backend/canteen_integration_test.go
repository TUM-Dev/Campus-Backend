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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCanteenQueriesWithMySQL(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}
	visited := time.Date(2026, time.November, 2, 12, 0, 0, 0, time.UTC)

	canteenTag := model.CanteenRatingTagOption{DE: "Sauberkeit", EN: "Cleanliness"}
	dishTag := model.DishRatingTagOption{DE: "Geschmack", EN: "Taste"}
	require.NoError(t, db.Create(&canteenTag).Error)
	require.NoError(t, db.Create(&dishTag).Error)

	for _, id := range []string{"mensa-garching", "mensa-arcisstr", "mensa-leopoldstr"} {
		t.Run(id, func(t *testing.T) {
			canteen := model.Canteen{Name: model.CanteenEnumName(id)}
			require.NoError(t, db.Create(&canteen).Error)
			dish := model.Dish{Name: "Pasta", Type: "Pasta", CanteenID: canteen.CanteenID}
			require.NoError(t, db.Create(&dish).Error)
			for day := int32(0); day < 5; day++ {
				require.NoError(t, db.Create(&model.DishesOfTheWeek{
					DishID: dish.Dish, Year: 2026, Week: 45, Day: day,
				}).Error)
			}
			canteenRating := model.CanteenRating{
				CanteenID: canteen.CanteenID, Points: 4, Comment: id, Timestamp: visited,
			}
			require.NoError(t, db.Create(&canteenRating).Error)
			require.NoError(t, db.Create(&model.CanteenRatingTag{
				RatingID: canteenRating.CanteenRating, TagID: canteenTag.CanteenRatingTagOption, Points: 3,
			}).Error)
			dishRating := model.DishRating{
				DishID: dish.Dish, Points: 5, Comment: id, Timestamp: visited,
			}
			require.NoError(t, db.Create(&dishRating).Error)
			require.NoError(t, db.Create(&model.DishRatingTag{
				RatingID: dishRating.DishRating, TagID: dishTag.DishRatingTagOption, Points: 2,
			}).Error)

			for _, canteenID := range []string{id, canteen.Name} {
				for day := int32(0); day < 5; day++ {
					dishes, err := server.ListDishes(ctx, &pb.ListDishesRequest{
						CanteenId: canteenID, Year: 2026, Week: 45, Day: day,
					})
					require.NoError(t, err)
					require.Equal(t, []string{dish.Name}, dishes.Dish)
				}

				for _, filtered := range []bool{false, true} {
					canteenRequest := &pb.ListCanteenRatingsRequest{CanteenId: canteenID, Limit: -1}
					dishRequest := &pb.GetDishRatingsRequest{CanteenId: canteenID, Dish: dish.Name, Limit: -1}
					if filtered {
						from := timestamppb.New(visited.Add(-time.Hour))
						to := timestamppb.New(visited.Add(time.Hour))
						canteenRequest.From, canteenRequest.To = from, to
						dishRequest.From, dishRequest.To = from, to
					}

					canteenRatings, err := server.ListCanteenRatings(ctx, canteenRequest)
					require.NoError(t, err)
					require.Equal(t, float64(4), canteenRatings.Avg)
					require.Equal(t, int32(4), canteenRatings.Min)
					require.Equal(t, int32(4), canteenRatings.Max)
					require.Zero(t, canteenRatings.Std)
					require.Len(t, canteenRatings.Rating, 1)
					require.Equal(t, id, canteenRatings.Rating[0].Comment)
					require.Len(t, canteenRatings.Rating[0].RatingTags, 1)
					require.Equal(t, int32(3), canteenRatings.Rating[0].RatingTags[0].Points)
					require.Len(t, canteenRatings.RatingTags, 1)
					require.Equal(t, int32(3), canteenRatings.RatingTags[0].Min)
					require.Equal(t, int32(3), canteenRatings.RatingTags[0].Max)

					dishRatings, err := server.GetDishRatings(ctx, dishRequest)
					require.NoError(t, err)
					require.Equal(t, float64(5), dishRatings.Avg)
					require.Equal(t, int32(5), dishRatings.Min)
					require.Equal(t, int32(5), dishRatings.Max)
					require.Zero(t, dishRatings.Std)
					require.Len(t, dishRatings.Rating, 1)
					require.Equal(t, id, dishRatings.Rating[0].Comment)
					require.Len(t, dishRatings.Rating[0].RatingTags, 1)
					require.Equal(t, int32(2), dishRatings.Rating[0].RatingTags[0].Points)
					require.Len(t, dishRatings.RatingTags, 1)
					require.Equal(t, int32(2), dishRatings.RatingTags[0].Min)
					require.Equal(t, int32(2), dishRatings.RatingTags[0].Max)
				}
			}
		})
	}

	dishes, err := server.ListDishes(ctx, &pb.ListDishesRequest{
		CanteenId: "mensa-garching", Year: 2026, Week: 46, Day: 0,
	})
	require.NoError(t, err)
	require.Empty(t, dishes.Dish)

	_, err = server.ListCanteenRatings(ctx, &pb.ListCanteenRatingsRequest{CanteenId: "unknown"})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = server.GetDishRatings(ctx, &pb.GetDishRatingsRequest{CanteenId: "mensa-garching", Dish: "unknown"})
	require.Equal(t, codes.NotFound, status.Code(err))
}
