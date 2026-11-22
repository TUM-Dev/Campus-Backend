package backend

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

func TestCreateDishRatingAtomicityAndWildcardSafety(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}

	canteen := model.Canteen{Name: "MENSA_GARCHING"}
	require.NoError(t, db.Create(&canteen).Error)
	dish := model.Dish{Name: "Pasta", Type: "Pasta", CanteenID: canteen.CanteenID}
	require.NoError(t, db.Create(&dish).Error)
	tag := model.DishRatingTagOption{DE: "Salz", EN: "Saltiness"}
	require.NoError(t, db.Create(&tag).Error)

	// SQL wildcard % must not match "Pasta".
	_, err := server.CreateDishRating(ctx, &pb.CreateDishRatingRequest{
		Points: 4, CanteenId: "mensa-garching", Dish: "%",
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	// Invalid tag must not leave a partial dish_ratings row behind.
	_, err = server.CreateDishRating(ctx, &pb.CreateDishRatingRequest{
		Points:     5,
		CanteenId:  "mensa-garching",
		Dish:       "Pasta",
		RatingTags: []*pb.RatingTag{{TagId: tag.DishRatingTagOption, Points: 4}, {TagId: 99999, Points: 3}},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	var ratingCount int64
	require.NoError(t, db.Model(&model.DishRating{}).Count(&ratingCount).Error)
	require.Zero(t, ratingCount)
}

func deviceCtx(deviceID string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-device-id", deviceID))
}

func TestDishRatingOwnershipDeduplicationAndModeration(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	prevStorage := cron.StorageDir
	cron.StorageDir = t.TempDir()
	t.Cleanup(func() { cron.StorageDir = prevStorage })

	server := CampusServer{db: db}

	canteen := model.Canteen{Name: "MENSA_GARCHING"}
	require.NoError(t, db.Create(&canteen).Error)
	dish := model.Dish{Name: "Schnitzel", Type: "Fleisch", CanteenID: canteen.CanteenID}
	require.NoError(t, db.Create(&dish).Error)
	tag1 := model.DishRatingTagOption{DE: "Salz", EN: "Saltiness"}
	tag2 := model.DishRatingTagOption{DE: "Portionsgröße", EN: "Portion Size"}
	require.NoError(t, db.Create(&tag1).Error)
	require.NoError(t, db.Create(&tag2).Error)
	nameTag := model.DishNameTagOption{DE: "Schnitzel", EN: "Schnitzel"}
	require.NoError(t, db.Create(&nameTag).Error)
	require.NoError(t, db.Create(&model.DishToDishNameTag{DishID: dish.Dish, NameTagID: nameTag.DishNameTagOption}).Error)

	img := readTestImage(t)
	aliceCtx := deviceCtx("alice-device")
	bobCtx := deviceCtx("bob-device")

	// Alice submits a rating with image and tag1.
	created, err := server.CreateDishRating(aliceCtx, &pb.CreateDishRatingRequest{
		Points:     4,
		CanteenId:  "mensa-garching",
		Dish:       "Schnitzel",
		Comment:    "Pretty good",
		Image:      img,
		RatingTags: []*pb.RatingTag{{TagId: tag1.DishRatingTagOption, Points: 3}},
	})
	require.NoError(t, err)
	require.Greater(t, created.RatingId, int64(0))

	// Alice cannot submit a second rating for the same dish on the same day.
	_, err = server.CreateDishRating(aliceCtx, &pb.CreateDishRatingRequest{
		Points:    5,
		CanteenId: "mensa-garching",
		Dish:      "Schnitzel",
	})
	require.Equal(t, codes.AlreadyExists, status.Code(err))

	// Bob CAN submit a rating on the same day.
	bobCreated, err := server.CreateDishRating(bobCtx, &pb.CreateDishRatingRequest{
		Points:    2,
		CanteenId: "mensa-garching",
		Dish:      "Schnitzel",
		Comment:   "Too dry",
	})
	require.NoError(t, err)

	// GetDishRatings marks is_own_rating accurately per caller.
	aliceView, err := server.GetDishRatings(aliceCtx, &pb.GetDishRatingsRequest{
		CanteenId: "mensa-garching", Dish: "Schnitzel", Limit: 10,
	})
	require.NoError(t, err)
	require.Len(t, aliceView.Rating, 2)
	ownCount := 0
	for _, r := range aliceView.Rating {
		if r.RatingId == created.RatingId {
			require.True(t, r.IsOwnRating)
			ownCount++
		} else {
			require.False(t, r.IsOwnRating)
		}
	}
	require.Equal(t, 1, ownCount)

	// Bob cannot update or delete Alice's rating.
	_, err = server.UpdateDishRating(bobCtx, &pb.UpdateDishRatingRequest{
		RatingId: created.RatingId, Points: 1, Comment: "Hacked",
	})
	require.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = server.DeleteDishRating(bobCtx, &pb.DeleteDishRatingRequest{RatingId: created.RatingId})
	require.Equal(t, codes.PermissionDenied, status.Code(err))

	// Alice updates her rating (changes points, removes image, switches tag).
	var before model.DishRating
	require.NoError(t, db.First(&before, created.RatingId).Error)
	require.FileExists(t, before.Image)

	_, err = server.UpdateDishRating(aliceCtx, &pb.UpdateDishRatingRequest{
		RatingId:    created.RatingId,
		Points:      5,
		Comment:     "Actually amazing",
		RemoveImage: true,
		RatingTags:  []*pb.RatingTag{{TagId: tag2.DishRatingTagOption, Points: 5}},
	})
	require.NoError(t, err)

	var after model.DishRating
	require.NoError(t, db.First(&after, created.RatingId).Error)
	require.Equal(t, int32(5), after.Points)
	require.Equal(t, "Actually amazing", after.Comment)
	require.Empty(t, after.Image)
	require.True(t, after.UpdatedAt.Valid)
	require.Equal(t, before.Timestamp.Unix(), after.Timestamp.Unix(), "visited timestamp must be preserved on edit")
	_, statErr := os.Stat(before.Image)
	require.True(t, os.IsNotExist(statErr), "removed image file should be cleaned up")

	// Bob cannot report his own rating.
	_, err = server.ReportDishRating(bobCtx, &pb.ReportDishRatingRequest{
		RatingId: bobCreated.RatingId, Reason: "self report",
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	// Reporting Bob's rating: duplicate report from same device rejected, 5 distinct reports auto-hide it.
	for i := 1; i <= ratingAutoHideReportThreshold; i++ {
		reporterCtx := deviceCtx(fmt.Sprintf("reporter-%d", i))
		rep, err := server.ReportDishRating(reporterCtx, &pb.ReportDishRatingRequest{
			RatingId: bobCreated.RatingId,
			Reason:   "spam",
		})
		require.NoError(t, err)
		if i < ratingAutoHideReportThreshold {
			require.False(t, rep.Hidden)
		} else {
			require.True(t, rep.Hidden)
		}
		if i == 1 {
			_, dupErr := server.ReportDishRating(reporterCtx, &pb.ReportDishRatingRequest{
				RatingId: bobCreated.RatingId, Reason: "spam again",
			})
			require.Equal(t, codes.AlreadyExists, status.Code(dupErr))
		}
	}

	// Bob cannot delete his hidden rating (preserving moderation reports) or re-rate the same dish today.
	_, err = server.DeleteDishRating(bobCtx, &pb.DeleteDishRatingRequest{RatingId: bobCreated.RatingId})
	require.Equal(t, codes.NotFound, status.Code(err))
	var reportCount int64
	require.NoError(t, db.Model(&model.DishRatingReport{}).Where("rating_id = ?", bobCreated.RatingId).Count(&reportCount).Error)
	require.Equal(t, int64(ratingAutoHideReportThreshold), reportCount)
	_, err = server.CreateDishRating(bobCtx, &pb.CreateDishRatingRequest{
		Points: 1, CanteenId: "mensa-garching", Dish: "Schnitzel",
	})
	require.Equal(t, codes.AlreadyExists, status.Code(err))

	// Hidden rating must be excluded from both rating list AND dish_rating_statistics view.
	afterHide, err := server.GetDishRatings(aliceCtx, &pb.GetDishRatingsRequest{
		CanteenId: "mensa-garching", Dish: "Schnitzel", Limit: 10,
	})
	require.NoError(t, err)
	require.Len(t, afterHide.Rating, 1)
	require.Equal(t, created.RatingId, afterHide.Rating[0].RatingId)
	require.Equal(t, float64(5), afterHide.Avg)

	// Alice deletes her own rating; dish now has no visible ratings.
	_, err = server.DeleteDishRating(aliceCtx, &pb.DeleteDishRatingRequest{RatingId: created.RatingId})
	require.NoError(t, err)
	_, err = server.GetDishRatings(aliceCtx, &pb.GetDishRatingsRequest{
		CanteenId: "mensa-garching", Dish: "Schnitzel", Limit: 10,
	})
	require.Equal(t, codes.NotFound, status.Code(err))

	// Canteen ratings also enforce per-canteen cooldown.
	canteenRating1, err := server.CreateCanteenRating(aliceCtx, &pb.CreateCanteenRatingRequest{
		Points: 4, CanteenId: "mensa-garching", Comment: "Clean",
	})
	require.NoError(t, err)
	require.Greater(t, canteenRating1.RatingId, int64(0))
	_, err = server.CreateCanteenRating(aliceCtx, &pb.CreateCanteenRatingRequest{
		Points: 5, CanteenId: "mensa-garching", Comment: "Second rating within cooldown",
	})
	require.Equal(t, codes.AlreadyExists, status.Code(err))
	canteenRating2, err := server.CreateCanteenRating(bobCtx, &pb.CreateCanteenRatingRequest{
		Points: 3, CanteenId: "mensa-garching", Comment: "Busy",
	})
	require.NoError(t, err)
	require.Greater(t, canteenRating2.RatingId, canteenRating1.RatingId)
}

func TestDishAndCanteenRatingCooldownAndDailyLimits(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}

	canteens := []model.Canteen{
		{Name: "MENSA_GARCHING"},
		{Name: "MENSA_ARCISSTR"},
		{Name: "MENSA_LEOPOLDSTR"},
	}
	for i := range canteens {
		require.NoError(t, db.Create(&canteens[i]).Error)
	}

	dishNames := []string{"Soup", "Salad", "Main", "Side", "Dessert", "ExtraSnack"}
	dishes := make([]model.Dish, len(dishNames))
	for i, name := range dishNames {
		dishes[i] = model.Dish{Name: name, Type: "Tagesgericht", CanteenID: canteens[0].CanteenID}
		require.NoError(t, db.Create(&dishes[i]).Error)
	}

	userCtx := deviceCtx("daily-limit-device")

	// A user can rate up to 5 different dishes in a single day (e.g. soup, salad, main, side, dessert).
	for i := 0; i < maxDishRatingsPerDeviceDay; i++ {
		_, err := server.CreateDishRating(userCtx, &pb.CreateDishRatingRequest{
			Points: 4, CanteenId: "mensa-garching", Dish: dishNames[i],
		})
		require.NoError(t, err)
	}
	// A 6th dish rating on the same day hits the daily cap.
	_, err := server.CreateDishRating(userCtx, &pb.CreateDishRatingRequest{
		Points: 4, CanteenId: "mensa-garching", Dish: dishNames[5],
	})
	require.Equal(t, codes.ResourceExhausted, status.Code(err))

	// A recurring dish rated more than 7 days ago can be rated again by the same device.
	recurringCtx := deviceCtx("recurring-device")
	require.NoError(t, db.Create(&model.DishRating{
		DishID:     dishes[0].Dish,
		Points:     3,
		Comment:    "Last week",
		Timestamp:  time.Now().AddDate(0, 0, -(dishRatingCooldownDays + 1)),
		DeviceHash: callerDeviceHash(recurringCtx),
		Status:     model.RatingStatusVisible,
	}).Error)
	_, err = server.CreateDishRating(recurringCtx, &pb.CreateDishRatingRequest{
		Points: 5, CanteenId: "mensa-garching", Dish: dishNames[0], Comment: "Better this week",
	})
	require.NoError(t, err)
	// Rating it again within the 7-day cooldown is rejected.
	_, err = server.CreateDishRating(recurringCtx, &pb.CreateDishRatingRequest{
		Points: 4, CanteenId: "mensa-garching", Dish: dishNames[0],
	})
	require.Equal(t, codes.AlreadyExists, status.Code(err))

	// Canteen ratings allow up to 2 different canteens per day; a 3rd hits ResourceExhausted.
	_, err = server.CreateCanteenRating(userCtx, &pb.CreateCanteenRatingRequest{Points: 4, CanteenId: "mensa-garching"})
	require.NoError(t, err)
	_, err = server.CreateCanteenRating(userCtx, &pb.CreateCanteenRatingRequest{Points: 4, CanteenId: "mensa-arcisstr"})
	require.NoError(t, err)
	_, err = server.CreateCanteenRating(userCtx, &pb.CreateCanteenRatingRequest{Points: 4, CanteenId: "mensa-leopoldstr"})
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}
