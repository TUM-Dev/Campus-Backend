package backend

import (
	"context"
	"sort"
	"strings"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	defaultTopDishesLimit = 10
	maxTopDishesLimit     = 50
	maxSummaryDishesBatch = 100
	bayesianPriorWeight   = 5.0
	bayesianDefaultPrior  = 3.0
)

type aggregateStatsRow struct {
	DishID      int64   `gorm:"column:dishID"`
	Average     float64 `gorm:"column:average"`
	Std         float64 `gorm:"column:std"`
	Min         int32   `gorm:"column:min"`
	Max         int32   `gorm:"column:max"`
	RatingCount int64   `gorm:"column:rating_count"`
	Stars1      int64   `gorm:"column:stars_1"`
	Stars2      int64   `gorm:"column:stars_2"`
	Stars3      int64   `gorm:"column:stars_3"`
	Stars4      int64   `gorm:"column:stars_4"`
	Stars5      int64   `gorm:"column:stars_5"`
}

func buildRatingDistribution(s1, s2, s3, s4, s5 int64) *pb.RatingDistribution {
	return &pb.RatingDistribution{
		Stars_1: s1,
		Stars_2: s2,
		Stars_3: s3,
		Stars_4: s4,
		Stars_5: s5,
	}
}

func resolveTimeWindow(fromPB, toPB *timestamppb.Timestamp) (time.Time, time.Time, bool) {
	if fromPB == nil && toPB == nil {
		return time.Time{}, time.Time{}, false
	}
	from := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
	if fromPB != nil {
		from = fromPB.AsTime()
	}
	to := time.Now()
	if toPB != nil {
		to = toPB.AsTime()
	}
	return from, to, true
}

func validateISODateSelector(year, week, day int32) error {
	if year < 2022 {
		return status.Error(codes.InvalidArgument, "Years must be larger or equal to 2022 ")
	}
	_, lastWeek := time.Date(int(year), time.December, 28, 0, 0, 0, 0, time.UTC).ISOWeek()
	if week < 1 || week > int32(lastWeek) {
		return status.Errorf(codes.InvalidArgument, "Weeks must be in the range 1 - %d for ISO year %d", lastWeek, year)
	}
	if day < 0 || day > 4 {
		return status.Error(codes.InvalidArgument, "Days must be in the range 0 (Monday) - 4 (Friday)")
	}
	return nil
}

func queryWindowedCanteenStats(canteenID int32, from, to time.Time, tx *gorm.DB) (aggregateStatsRow, error) {
	var row aggregateStatsRow
	err := tx.Table("canteen_ratings").
		Select(`COALESCE(AVG(points), 0) AS average,
			COALESCE(STD(points), 0) AS std,
			COALESCE(MIN(points), 0) AS min,
			COALESCE(MAX(points), 0) AS max,
			COUNT(*) AS rating_count,
			COALESCE(SUM(CASE WHEN points = 1 THEN 1 ELSE 0 END), 0) AS stars_1,
			COALESCE(SUM(CASE WHEN points = 2 THEN 1 ELSE 0 END), 0) AS stars_2,
			COALESCE(SUM(CASE WHEN points = 3 THEN 1 ELSE 0 END), 0) AS stars_3,
			COALESCE(SUM(CASE WHEN points = 4 THEN 1 ELSE 0 END), 0) AS stars_4,
			COALESCE(SUM(CASE WHEN points = 5 THEN 1 ELSE 0 END), 0) AS stars_5`).
		Where("cafeteriaID = ? AND status = ? AND timestamp > ? AND timestamp < ?", canteenID, model.RatingStatusVisible, from, to).
		Scan(&row).Error
	return row, err
}

func queryWindowedDishStats(dishID int32, from, to time.Time, tx *gorm.DB) (aggregateStatsRow, error) {
	var row aggregateStatsRow
	err := tx.Table("dish_ratings").
		Select(`COALESCE(AVG(points), 0) AS average,
			COALESCE(STD(points), 0) AS std,
			COALESCE(MIN(points), 0) AS min,
			COALESCE(MAX(points), 0) AS max,
			COUNT(*) AS rating_count,
			COALESCE(SUM(CASE WHEN points = 1 THEN 1 ELSE 0 END), 0) AS stars_1,
			COALESCE(SUM(CASE WHEN points = 2 THEN 1 ELSE 0 END), 0) AS stars_2,
			COALESCE(SUM(CASE WHEN points = 3 THEN 1 ELSE 0 END), 0) AS stars_3,
			COALESCE(SUM(CASE WHEN points = 4 THEN 1 ELSE 0 END), 0) AS stars_4,
			COALESCE(SUM(CASE WHEN points = 5 THEN 1 ELSE 0 END), 0) AS stars_5`).
		Where("dishID = ? AND status = ? AND timestamp > ? AND timestamp < ?", dishID, model.RatingStatusVisible, from, to).
		Scan(&row).Error
	return row, err
}

func queryWindowedTags(canteenID, dishID int32, ratingType model.ModelType, from, to time.Time, tx *gorm.DB) []*pb.RatingTagResult {
	var results []queryRatingTag
	var err error
	switch ratingType {
	case model.DISH:
		err = tx.Table("dish_rating_tags mrt").
			Joins("JOIN dish_ratings mr ON mr.dishRating = mrt.parentRating").
			Select("mrt.tagID AS tagId, AVG(mrt.points) AS avg, MIN(mrt.points) AS min, MAX(mrt.points) AS max, STD(mrt.points) AS std").
			Where("mr.dishID = ? AND mr.status = ? AND mr.timestamp > ? AND mr.timestamp < ?", dishID, model.RatingStatusVisible, from, to).
			Group("mrt.tagID").
			Scan(&results).Error
	case model.CANTEEN:
		err = tx.Table("canteen_rating_tags crt").
			Joins("JOIN canteen_ratings cr ON cr.cafeteriaRating = crt.correspondingRating").
			Select("crt.tagID AS tagId, AVG(crt.points) AS avg, MIN(crt.points) AS min, MAX(crt.points) AS max, STD(crt.points) AS std").
			Where("cr.cafeteriaID = ? AND cr.status = ? AND cr.timestamp > ? AND cr.timestamp < ?", canteenID, model.RatingStatusVisible, from, to).
			Group("crt.tagID").
			Scan(&results).Error
	default:
		return queryTags(canteenID, dishID, ratingType, tx)
	}
	if err != nil {
		log.WithError(err).Error("while querying windowed rating tags")
	}
	resp := make([]*pb.RatingTagResult, 0, len(results))
	for _, v := range results {
		resp = append(resp, &pb.RatingTagResult{
			TagId: v.TagId,
			Avg:   v.Average,
			Std:   v.Std,
			Min:   v.Min,
			Max:   v.Max,
		})
	}
	return resp
}

func parseDishLabels(raw null.String) []string {
	if !raw.Valid || strings.TrimSpace(raw.String) == "" {
		return nil
	}
	parts := strings.Split(raw.String, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func dishHasAllLabels(dishLabels []string, required []string) bool {
	if len(required) == 0 {
		return true
	}
	have := make(map[string]bool, len(dishLabels))
	for _, l := range dishLabels {
		have[strings.ToUpper(l)] = true
	}
	for _, req := range required {
		r := strings.ToUpper(strings.TrimSpace(req))
		if r != "" && !have[r] {
			return false
		}
	}
	return true
}

// ListDishRatingSummaries batch-loads rating statistics, tag averages, name-tag fallbacks,
// and latest photo rating IDs for an entire day's menu or an explicit list of dish names.
func (s *CampusServer) ListDishRatingSummaries(ctx context.Context, req *pb.ListDishRatingSummariesRequest) (*pb.ListDishRatingSummariesReply, error) {
	s.recordDeviceIfPresent(ctx)
	if req.GetCanteenId() == "" {
		return nil, status.Error(codes.InvalidArgument, "canteen_id is required.")
	}
	useMenuDate := req.Year > 0
	if useMenuDate {
		if err := validateISODateSelector(req.Year, req.Week, req.Day); err != nil {
			return nil, err
		}
	} else {
		if len(req.Dishes) == 0 {
			return nil, status.Error(codes.InvalidArgument, "Either (year, week, day) or dishes must be provided.")
		}
		if len(req.Dishes) > maxSummaryDishesBatch {
			return nil, status.Errorf(codes.InvalidArgument, "At most %d dishes can be summarized per request.", maxSummaryDishesBatch)
		}
	}

	tx := s.db.WithContext(ctx)
	canteenID, err := getIDForCanteenName(req.CanteenId, tx)
	if err != nil {
		return nil, err
	}

	var dishes []model.Dish
	var orderedNames []string
	if useMenuDate {
		err = tx.Table("dishes_of_the_weeks weekly").
			Joins("JOIN dishes d ON d.dish = weekly.dishID").
			Where("weekly.day = ? AND weekly.week = ? AND weekly.year = ? AND d.cafeteriaID = ?", req.Day, req.Week, req.Year, canteenID).
			Select("DISTINCT d.dish, d.name, d.type, d.cafeteriaID, d.labels, d.price_student_cents").
			Order("d.name ASC").
			Find(&dishes).Error
		if err != nil {
			log.WithError(err).Error("while loading daily dishes for rating summaries")
			return nil, status.Error(codes.Internal, "Could not load dishes.")
		}
		for _, d := range dishes {
			orderedNames = append(orderedNames, d.Name)
		}
	} else {
		orderedNames = req.Dishes
		uniqueNames := make([]string, 0, len(req.Dishes))
		seenNames := make(map[string]bool, len(req.Dishes))
		for _, n := range req.Dishes {
			if !seenNames[n] {
				seenNames[n] = true
				uniqueNames = append(uniqueNames, n)
			}
		}
		err = tx.Where("cafeteriaID = ? AND name IN ?", canteenID, uniqueNames).Find(&dishes).Error
		if err != nil {
			log.WithError(err).Error("while loading requested dishes for rating summaries")
			return nil, status.Error(codes.Internal, "Could not load dishes.")
		}
	}

	dishByName := make(map[string]model.Dish, len(dishes))
	dishIDs := make([]int64, 0, len(dishes))
	for _, d := range dishes {
		dishByName[d.Name] = d
		dishIDs = append(dishIDs, d.Dish)
	}

	if len(dishIDs) == 0 {
		summaries := make([]*pb.DishRatingSummary, 0, len(orderedNames))
		for _, name := range orderedNames {
			summaries = append(summaries, &pb.DishRatingSummary{
				Dish:         name,
				Distribution: buildRatingDistribution(0, 0, 0, 0, 0),
			})
		}
		return &pb.ListDishRatingSummariesReply{Summaries: summaries}, nil
	}

	from, to, windowed := resolveTimeWindow(req.From, req.To)

	statsByDish, err := batchLoadDishStats(canteenID, dishIDs, from, to, windowed, tx)
	if err != nil {
		log.WithError(err).Error("while batch loading dish rating statistics")
		return nil, status.Error(codes.Internal, "Could not load dish rating summaries.")
	}
	ratingTagsByDish := batchLoadDishRatingTags(canteenID, dishIDs, from, to, windowed, tx)
	nameTagsByDish := batchLoadDishNameTags(canteenID, dishIDs, tx)
	latestImgByDish := batchLoadLatestImageRatingIDs(dishIDs, tx)

	summaries := make([]*pb.DishRatingSummary, 0, len(orderedNames))
	for _, name := range orderedNames {
		d, exists := dishByName[name]
		if !exists {
			summaries = append(summaries, &pb.DishRatingSummary{
				Dish:         name,
				Distribution: buildRatingDistribution(0, 0, 0, 0, 0),
			})
			continue
		}
		st := statsByDish[d.Dish]
		imgRatingID := latestImgByDish[d.Dish]
		var priceCents int32
		if d.PriceStudentCents.Valid {
			priceCents = int32(d.PriceStudentCents.Int64)
		}
		summaries = append(summaries, &pb.DishRatingSummary{
			Dish:                d.Name,
			Avg:                 st.Average,
			Std:                 st.Std,
			Min:                 st.Min,
			Max:                 st.Max,
			RatingCount:         st.RatingCount,
			Distribution:        buildRatingDistribution(st.Stars1, st.Stars2, st.Stars3, st.Stars4, st.Stars5),
			RatingTags:          ratingTagsByDish[d.Dish],
			NameTags:            nameTagsByDish[d.Dish],
			LatestImageRatingId: imgRatingID,
			Labels:              parseDishLabels(d.Labels),
			PriceStudentCents:   priceCents,
		})
	}
	return &pb.ListDishRatingSummariesReply{Summaries: summaries}, nil
}

func batchLoadDishStats(canteenID int32, dishIDs []int64, from, to time.Time, windowed bool, tx *gorm.DB) (map[int64]aggregateStatsRow, error) {
	var rows []aggregateStatsRow
	var err error
	if windowed {
		err = tx.Table("dish_ratings").
			Select(`dishID,
				AVG(points) AS average,
				STD(points) AS std,
				MIN(points) AS min,
				MAX(points) AS max,
				COUNT(*) AS rating_count,
				SUM(CASE WHEN points = 1 THEN 1 ELSE 0 END) AS stars_1,
				SUM(CASE WHEN points = 2 THEN 1 ELSE 0 END) AS stars_2,
				SUM(CASE WHEN points = 3 THEN 1 ELSE 0 END) AS stars_3,
				SUM(CASE WHEN points = 4 THEN 1 ELSE 0 END) AS stars_4,
				SUM(CASE WHEN points = 5 THEN 1 ELSE 0 END) AS stars_5`).
			Where("dishID IN ? AND status = ? AND timestamp > ? AND timestamp < ?", dishIDs, model.RatingStatusVisible, from, to).
			Group("dishID").
			Scan(&rows).Error
	} else {
		err = tx.Table("dish_rating_statistics").
			Select("dishID, average, std, min, max, rating_count, stars_1, stars_2, stars_3, stars_4, stars_5").
			Where("cafeteriaID = ? AND dishID IN ?", canteenID, dishIDs).
			Scan(&rows).Error
	}
	if err != nil {
		return nil, err
	}
	out := make(map[int64]aggregateStatsRow, len(rows))
	for _, r := range rows {
		out[r.DishID] = r
	}
	return out, nil
}

type batchTagRow struct {
	DishID  int64   `gorm:"column:dishID"`
	TagID   int32   `gorm:"column:tagId"`
	Average float64 `gorm:"column:avg"`
	Std     float64 `gorm:"column:std"`
	Min     int32   `gorm:"column:min"`
	Max     int32   `gorm:"column:max"`
}

func batchLoadDishRatingTags(canteenID int32, dishIDs []int64, from, to time.Time, windowed bool, tx *gorm.DB) map[int64][]*pb.RatingTagResult {
	var rows []batchTagRow
	var err error
	if windowed {
		err = tx.Table("dish_rating_tags mrt").
			Joins("JOIN dish_ratings mr ON mr.dishRating = mrt.parentRating").
			Select("mr.dishID AS dishID, mrt.tagID AS tagId, AVG(mrt.points) AS avg, MIN(mrt.points) AS min, MAX(mrt.points) AS max, STD(mrt.points) AS std").
			Where("mr.dishID IN ? AND mr.status = ? AND mr.timestamp > ? AND mr.timestamp < ?", dishIDs, model.RatingStatusVisible, from, to).
			Group("mr.dishID, mrt.tagID").
			Scan(&rows).Error
	} else {
		err = tx.Table("dish_rating_tag_statistics").
			Select("dishID, tagID AS tagId, average AS avg, min, max, std").
			Where("cafeteriaID = ? AND dishID IN ?", canteenID, dishIDs).
			Scan(&rows).Error
	}
	if err != nil {
		log.WithError(err).Error("while batch loading dish rating tag statistics")
	}
	return groupBatchTagRows(rows)
}

func batchLoadDishNameTags(canteenID int32, dishIDs []int64, tx *gorm.DB) map[int64][]*pb.RatingTagResult {
	var rows []batchTagRow
	err := tx.Table("dish_to_dish_name_tags mapping").
		Joins("JOIN dish_name_tag_statistics results ON mapping.nameTagID = results.tagID").
		Select("mapping.dishID AS dishID, mapping.nameTagID AS tagId, results.average AS avg, results.min AS min, results.max AS max, results.std AS std").
		Where("mapping.dishID IN ? AND results.cafeteriaID = ?", dishIDs, canteenID).
		Scan(&rows).Error
	if err != nil {
		log.WithError(err).Error("while batch loading dish name tag statistics")
	}
	return groupBatchTagRows(rows)
}

func groupBatchTagRows(rows []batchTagRow) map[int64][]*pb.RatingTagResult {
	out := make(map[int64][]*pb.RatingTagResult)
	for _, r := range rows {
		out[r.DishID] = append(out[r.DishID], &pb.RatingTagResult{
			TagId: r.TagID,
			Avg:   r.Average,
			Std:   r.Std,
			Min:   r.Min,
			Max:   r.Max,
		})
	}
	return out
}

type latestImageRow struct {
	DishID   int64 `gorm:"column:dishID"`
	RatingID int64 `gorm:"column:ratingID"`
}

func batchLoadLatestImageRatingIDs(dishIDs []int64, tx *gorm.DB) map[int64]int64 {
	out := make(map[int64]int64)
	if len(dishIDs) == 0 {
		return out
	}
	var rows []latestImageRow
	err := tx.Table("dish_ratings").
		Select("dishID, MAX(dishRating) AS ratingID").
		Where("dishID IN ? AND status = ? AND image <> ''", dishIDs, model.RatingStatusVisible).
		Group("dishID").
		Scan(&rows).Error
	if err != nil {
		log.WithError(err).Error("while batch loading latest rating images")
		return out
	}
	for _, r := range rows {
		out[r.DishID] = r.RatingID
	}
	return out
}

type topRatedDishRow struct {
	DishID            int64       `gorm:"column:dishID"`
	CanteenName       string      `gorm:"column:canteenName"`
	DishName          string      `gorm:"column:dishName"`
	Labels            null.String `gorm:"column:labels"`
	PriceStudentCents null.Int    `gorm:"column:price_student_cents"`
	Average           float64     `gorm:"column:average"`
	RatingCount       int64       `gorm:"column:rating_count"`
}

// ListTopRatedDishes returns the highest-rated dishes served on a given ISO date,
// ranked by Bayesian weighted average (v*R + m*C) / (v + m).
func (s *CampusServer) ListTopRatedDishes(ctx context.Context, req *pb.ListTopRatedDishesRequest) (*pb.ListTopRatedDishesReply, error) {
	s.recordDeviceIfPresent(ctx)
	if err := validateISODateSelector(req.Year, req.Week, req.Day); err != nil {
		return nil, err
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = defaultTopDishesLimit
	}
	if limit > maxTopDishesLimit {
		limit = maxTopDishesLimit
	}
	minRatings := int64(req.MinRatings)
	if minRatings <= 0 {
		minRatings = 1
	}

	tx := s.db.WithContext(ctx)
	var globalAvg struct {
		Avg float64 `gorm:"column:avg"`
	}
	if err := tx.Table("dish_ratings").
		Select("COALESCE(AVG(points), ?) AS avg", bayesianDefaultPrior).
		Where("status = ?", model.RatingStatusVisible).
		Scan(&globalAvg).Error; err != nil {
		log.WithError(err).Error("while loading global average dish rating")
		return nil, status.Error(codes.Internal, "Could not compute top rated dishes.")
	}
	priorMean := globalAvg.Avg
	if priorMean <= 0 {
		priorMean = bayesianDefaultPrior
	}

	q := tx.Table("dishes_of_the_weeks weekly").
		Joins("JOIN dishes d ON d.dish = weekly.dishID").
		Joins("JOIN canteens c ON c.cafeteria = d.cafeteriaID").
		Joins("JOIN dish_rating_statistics s ON s.dishID = d.dish AND s.cafeteriaID = d.cafeteriaID").
		Select("DISTINCT d.dish AS dishID, c.name AS canteenName, d.name AS dishName, d.labels AS labels, d.price_student_cents AS price_student_cents, s.average AS average, s.rating_count AS rating_count").
		Where("weekly.year = ? AND weekly.week = ? AND weekly.day = ? AND s.rating_count >= ?", req.Year, req.Week, req.Day, minRatings)

	if req.GetCanteenId() != "" {
		q = q.Where("c.name = ?", model.CanteenEnumName(req.CanteenId))
	}
	if req.MaxPriceStudentCents > 0 {
		q = q.Where("d.price_student_cents IS NOT NULL AND d.price_student_cents <= ?", req.MaxPriceStudentCents)
	}

	var rows []topRatedDishRow
	if err := q.Scan(&rows).Error; err != nil {
		log.WithError(err).Error("while loading top rated dishes")
		return nil, status.Error(codes.Internal, "Could not load top rated dishes.")
	}
	if len(rows) == 0 {
		return &pb.ListTopRatedDishesReply{Dishes: []*pb.TopRatedDish{}}, nil
	}

	type scoredDish struct {
		row    topRatedDishRow
		labels []string
		score  float64
	}
	scored := make([]scoredDish, 0, len(rows))
	for _, r := range rows {
		labels := parseDishLabels(r.Labels)
		if !dishHasAllLabels(labels, req.RequiredLabels) {
			continue
		}
		v := float64(r.RatingCount)
		score := (v*r.Average + bayesianPriorWeight*priorMean) / (v + bayesianPriorWeight)
		scored = append(scored, scoredDish{row: r, labels: labels, score: score})
	}
	if len(scored) == 0 {
		return &pb.ListTopRatedDishesReply{Dishes: []*pb.TopRatedDish{}}, nil
	}
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score > scored[j].score
		}
		if scored[i].row.RatingCount != scored[j].row.RatingCount {
			return scored[i].row.RatingCount > scored[j].row.RatingCount
		}
		return scored[i].row.DishName < scored[j].row.DishName
	})
	if len(scored) > limit {
		scored = scored[:limit]
	}

	dishIDs := make([]int64, len(scored))
	for i, s := range scored {
		dishIDs[i] = s.row.DishID
	}
	latestImgByDish := batchLoadLatestImageRatingIDs(dishIDs, tx)

	result := make([]*pb.TopRatedDish, 0, len(scored))
	for _, item := range scored {
		imgRatingID := latestImgByDish[item.row.DishID]
		var priceCents int32
		if item.row.PriceStudentCents.Valid {
			priceCents = int32(item.row.PriceStudentCents.Int64)
		}
		result = append(result, &pb.TopRatedDish{
			CanteenId:           model.CanteenEatAPIID(item.row.CanteenName),
			Dish:                item.row.DishName,
			Avg:                 item.row.Average,
			BayesianScore:       item.score,
			RatingCount:         item.row.RatingCount,
			LatestImageRatingId: imgRatingID,
			Labels:              item.labels,
			PriceStudentCents:   priceCents,
		})
	}
	return &pb.ListTopRatedDishesReply{Dishes: result}, nil
}
