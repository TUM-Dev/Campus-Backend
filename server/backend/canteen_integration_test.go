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

	// Canteen ratings also enforce same-day per-device deduplication.
	canteenRating1, err := server.CreateCanteenRating(aliceCtx, &pb.CreateCanteenRatingRequest{
		Points: 4, CanteenId: "mensa-garching", Comment: "Clean",
	})
	require.NoError(t, err)
	require.Greater(t, canteenRating1.RatingId, int64(0))
	_, err = server.CreateCanteenRating(aliceCtx, &pb.CreateCanteenRatingRequest{
		Points: 5, CanteenId: "mensa-garching", Comment: "Second rating same day",
	})
	require.Equal(t, codes.AlreadyExists, status.Code(err))
	canteenRating2, err := server.CreateCanteenRating(bobCtx, &pb.CreateCanteenRatingRequest{
		Points: 3, CanteenId: "mensa-garching", Comment: "Busy",
	})
	require.NoError(t, err)
	require.Greater(t, canteenRating2.RatingId, canteenRating1.RatingId)
}

func TestDishRatingSummariesHistogramsWindowingAndTopDishes(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}

	canteen := model.Canteen{Name: "MENSA_GARCHING"}
	require.NoError(t, db.Create(&canteen).Error)
	otherCanteen := model.Canteen{Name: "MENSA_ARCISSTR"}
	require.NoError(t, db.Create(&otherCanteen).Error)

	// Dishes offered on 2026-W45-Day 1 (Tuesday).
	curry := model.Dish{Name: "Thai Curry", Type: "Vegan", CanteenID: canteen.CanteenID}
	newCurry := model.Dish{Name: "Red Curry", Type: "Vegan", CanteenID: canteen.CanteenID}
	oneVoteDish := model.Dish{Name: "Mystery Soup", Type: "Suppe", CanteenID: canteen.CanteenID}
	arcisDish := model.Dish{Name: "Lasagne", Type: "Pasta", CanteenID: otherCanteen.CanteenID}
	for _, d := range []*model.Dish{&curry, &newCurry, &oneVoteDish, &arcisDish} {
		require.NoError(t, db.Create(d).Error)
		require.NoError(t, db.Create(&model.DishesOfTheWeek{
			DishID: d.Dish, Year: 2026, Week: 45, Day: 1,
		}).Error)
	}

	curryNameTag := model.DishNameTagOption{DE: "Curry", EN: "Curry"}
	require.NoError(t, db.Create(&curryNameTag).Error)
	require.NoError(t, db.Create(&model.DishToDishNameTag{DishID: curry.Dish, NameTagID: curryNameTag.DishNameTagOption}).Error)
	require.NoError(t, db.Create(&model.DishToDishNameTag{DishID: newCurry.Dish, NameTagID: curryNameTag.DishNameTagOption}).Error)

	oldTime := time.Date(2025, time.October, 1, 12, 0, 0, 0, time.UTC)
	recentTime := time.Date(2026, time.November, 3, 12, 0, 0, 0, time.UTC)

	// Thai Curry: 1 old rating (1 star) + 6 recent ratings (five 5-star, one 4-star) -> all-time avg ~4.29, recent avg ~4.83.
	ratingsToInsert := []struct {
		dishID int64
		points int32
		ts     time.Time
		img    string
	}{
		{curry.Dish, 1, oldTime, ""},
		{curry.Dish, 5, recentTime, ""},
		{curry.Dish, 5, recentTime.Add(time.Minute), ""},
		{curry.Dish, 5, recentTime.Add(2 * time.Minute), ""},
		{curry.Dish, 5, recentTime.Add(3 * time.Minute), ""},
		{curry.Dish, 5, recentTime.Add(4 * time.Minute), "/tmp/photo.jpeg"},
		{curry.Dish, 4, recentTime.Add(5 * time.Minute), ""},
		// Mystery Soup: only 1 rating of 5 stars (Bayesian score should rank below Thai Curry).
		{oneVoteDish.Dish, 5, recentTime, ""},
		// Lasagne at Arcisstr: 3 ratings of 2 stars.
		{arcisDish.Dish, 2, recentTime, ""},
		{arcisDish.Dish, 2, recentTime.Add(time.Minute), ""},
		{arcisDish.Dish, 2, recentTime.Add(2 * time.Minute), ""},
	}
	var latestCurryPhotoID int64
	for _, r := range ratingsToInsert {
		rec := model.DishRating{
			DishID:    r.dishID,
			Points:    r.points,
			Timestamp: r.ts,
			Image:     r.img,
			Status:    model.RatingStatusVisible,
		}
		require.NoError(t, db.Create(&rec).Error)
		if r.dishID == curry.Dish {
			require.NoError(t, db.Create(&model.DishNameTag{
				RatingID: rec.DishRating, TagNameID: curryNameTag.DishNameTagOption, Points: r.points,
			}).Error)
		}
		if r.img != "" {
			latestCurryPhotoID = rec.DishRating
		}
	}

	// 1. All-time GetDishRatings returns count and star histogram.
	allTime, err := server.GetDishRatings(ctx, &pb.GetDishRatingsRequest{
		CanteenId: "mensa-garching", Dish: "Thai Curry", Limit: 10,
	})
	require.NoError(t, err)
	require.Equal(t, int64(7), allTime.RatingCount)
	require.Equal(t, int64(1), allTime.Distribution.Stars_1)
	require.Equal(t, int64(1), allTime.Distribution.Stars_4)
	require.Equal(t, int64(5), allTime.Distribution.Stars_5)

	// 2. Time-windowed GetDishRatings computes stats strictly inside [from, to].
	recentOnly, err := server.GetDishRatings(ctx, &pb.GetDishRatingsRequest{
		CanteenId: "mensa-garching",
		Dish:      "Thai Curry",
		From:      timestamppb.New(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)),
		Limit:     10,
	})
	require.NoError(t, err)
	require.Equal(t, int64(6), recentOnly.RatingCount)
	require.Equal(t, int32(4), recentOnly.Min)
	require.Equal(t, int64(0), recentOnly.Distribution.Stars_1)
	require.Equal(t, int64(5), recentOnly.Distribution.Stars_5)

	// 3. Batch ListDishRatingSummaries for the day returns all 3 Garching dishes in one call,
	//    including unrated "Red Curry" with its Curry name-tag fallback stats!
	batch, err := server.ListDishRatingSummaries(ctx, &pb.ListDishRatingSummariesRequest{
		CanteenId: "mensa-garching", Year: 2026, Week: 45, Day: 1,
	})
	require.NoError(t, err)
	require.Len(t, batch.Summaries, 3)

	byDish := map[string]*pb.DishRatingSummary{}
	for _, s := range batch.Summaries {
		byDish[s.Dish] = s
	}
	require.Equal(t, int64(7), byDish["Thai Curry"].RatingCount)
	require.Equal(t, latestCurryPhotoID, byDish["Thai Curry"].LatestImageRatingId)
	require.Equal(t, int64(0), byDish["Red Curry"].RatingCount)
	require.Len(t, byDish["Red Curry"].NameTags, 1, "unrated dish should still surface name-tag fallback stats")
	require.Equal(t, int32(curryNameTag.DishNameTagOption), byDish["Red Curry"].NameTags[0].TagId)

	// 4. ListTopRatedDishes ranks 7-vote Thai Curry (avg 4.29) above 1-vote Mystery Soup (avg 5.0) via Bayesian weighting.
	top, err := server.ListTopRatedDishes(ctx, &pb.ListTopRatedDishesRequest{
		Year: 2026, Week: 45, Day: 1, Limit: 5,
	})
	require.NoError(t, err)
	require.Len(t, top.Dishes, 3)
	require.Equal(t, "Thai Curry", top.Dishes[0].Dish)
	require.Equal(t, "mensa-garching", top.Dishes[0].CanteenId)
	require.Equal(t, latestCurryPhotoID, top.Dishes[0].LatestImageRatingId)
	require.Greater(t, top.Dishes[0].BayesianScore, top.Dishes[1].BayesianScore)
	require.Equal(t, "Lasagne", top.Dishes[2].Dish)
	require.Equal(t, "mensa-arcisstr", top.Dishes[2].CanteenId)

	// 5. ListTopRatedDishes supports canteen_id and min_ratings filters.
	filteredTop, err := server.ListTopRatedDishes(ctx, &pb.ListTopRatedDishesRequest{
		Year: 2026, Week: 45, Day: 1, CanteenId: "mensa-garching", MinRatings: 2,
	})
	require.NoError(t, err)
	require.Len(t, filteredTop.Dishes, 1)
	require.Equal(t, "Thai Curry", filteredTop.Dishes[0].Dish)

	// 6. ListDishRatingSummaries works with explicit dish names + time window + enforces batch cap.
	explicitBatch, err := server.ListDishRatingSummaries(ctx, &pb.ListDishRatingSummariesRequest{
		CanteenId: "mensa-garching",
		Dishes:    []string{"Thai Curry", "Nonexistent Dish"},
		From:      timestamppb.New(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)),
	})
	require.NoError(t, err)
	require.Len(t, explicitBatch.Summaries, 2)
	require.Equal(t, int64(6), explicitBatch.Summaries[0].RatingCount)
	require.Equal(t, "Nonexistent Dish", explicitBatch.Summaries[1].Dish)
	require.Equal(t, int64(0), explicitBatch.Summaries[1].RatingCount)

	tooManyDishes := make([]string, maxSummaryDishesBatch+1)
	for i := range tooManyDishes {
		tooManyDishes[i] = fmt.Sprintf("Dish %d", i)
	}
	_, err = server.ListDishRatingSummaries(ctx, &pb.ListDishRatingSummariesRequest{
		CanteenId: "mensa-garching", Dishes: tooManyDishes,
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	// Duplicate dishes_of_the_weeks rows must not yield duplicate menu summaries.
	require.NoError(t, db.Create(&model.DishesOfTheWeek{DishID: curry.Dish, Year: 2026, Week: 45, Day: 1}).Error)
	dedupBatch, err := server.ListDishRatingSummaries(ctx, &pb.ListDishRatingSummariesRequest{
		CanteenId: "mensa-garching", Year: 2026, Week: 45, Day: 1,
	})
	require.NoError(t, err)
	require.Len(t, dedupBatch.Summaries, 3)
}
