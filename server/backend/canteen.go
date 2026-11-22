package backend

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
	"unicode/utf8"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/disintegration/imaging"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	maxRatingCommentRunes         = 256
	maxRatingImageWidth           = 1280
	targetMaxImageBytes           = 524288 // 512 KiB
	defaultJPEGQuality            = 95
	minJPEGQuality                = 40
	ratingAutoHideReportThreshold = 5
	dishRatingCooldownDays        = 7
	canteenRatingCooldownDays     = 7
	maxDishRatingsPerDeviceDay    = 5
	maxCanteenRatingsPerDeviceDay = 2
)

func (s *CampusServer) recordDeviceIfPresent(ctx context.Context) {
	if s.deviceBuf == nil {
		return
	}
	if _, ok := metadata.FromIncomingContext(ctx); !ok {
		return
	}
	_ = s.checkDevice(ctx)
}

func callerDeviceHash(ctx context.Context) null.String {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return null.String{}
	}
	vals := md.Get("x-device-id")
	if len(vals) == 0 {
		return null.String{}
	}
	deviceID := strings.TrimSpace(vals[0])
	if deviceID == "" || strings.EqualFold(deviceID, "unknown") {
		return null.String{}
	}
	sum := sha256.Sum256([]byte("campus-rating:" + deviceID))
	return null.StringFrom(hex.EncodeToString(sum[:]))
}

func berlinDayBounds(now time.Time) (time.Time, time.Time) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		loc = time.UTC
	}
	local := now.In(loc)
	start := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	return start, start.AddDate(0, 0, 1)
}

// ListCanteenRatings RPC Endpoint
// Allows to query ratings for a specific canteen.
// It returns the average rating, max/min rating as well as a number of actual ratings and the average ratings for
// all canteen rating tags which were used to rate this canteen.
// The parameter limit defines how many actual ratings should be returned.
// The optional parameters from and to can define an interval in which the queried ratings have been stored.
// If these aren't specified, the newest ratings will be returned as the default
func (s *CampusServer) ListCanteenRatings(ctx context.Context, input *pb.ListCanteenRatingsRequest) (*pb.ListCanteenRatingsReply, error) {
	s.recordDeviceIfPresent(ctx)
	if input.GetCanteenId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Canteen Id is required.")
	}
	var statsForCanteen model.CanteenRatingStatistic
	tx := s.db.WithContext(ctx)
	canteenId, err := getIDForCanteenName(input.CanteenId, tx)
	if err != nil {
		return nil, err
	}
	err = tx.First(&statsForCanteen, "cafeteriaId = ?", canteenId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "No canteen with this Id found.")
	}
	if err != nil {
		log.WithError(err).Error("Error while querying the canteen with Id ", canteenId)
		return nil, status.Error(codes.Internal, "could not query the canteen with the given Id")
	}

	callerHash := callerDeviceHash(ctx)
	return &pb.ListCanteenRatingsReply{
		Avg:        statsForCanteen.Average,
		Std:        statsForCanteen.Std,
		Min:        statsForCanteen.Min,
		Max:        statsForCanteen.Max,
		Rating:     queryLastCanteenRatingsWithLimit(input, canteenId, callerHash, tx),
		RatingTags: queryTags(canteenId, -1, model.CANTEEN, tx),
	}, nil
}

// queryLastCanteenRatingsWithLimit
// Queries the actual ratings for a canteen and attaches the tag ratings which belong to the ratings
func queryLastCanteenRatingsWithLimit(input *pb.ListCanteenRatingsRequest, canteenID int32, callerHash null.String, tx *gorm.DB) []*pb.SingleRatingReply {
	var ratings []model.CanteenRating
	var err error

	var limit = int(input.Limit)
	if limit == -1 {
		limit = math.MaxInt32
	}
	if limit > 0 {
		if input.From != nil || input.To != nil {

			var from time.Time
			var to time.Time
			if input.From == nil {
				from = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
			} else {
				from = input.From.AsTime()
			}

			if input.To == nil {
				to = time.Now()
			} else {
				to = input.To.AsTime()
			}
			err = tx.
				Order("timestamp desc, cafeteriaRating desc").
				Limit(limit).
				Find(&ratings, "cafeteriaID = ? AND status = ? AND timestamp < ? AND timestamp > ?", canteenID, model.RatingStatusVisible, to, from).Error
		} else {
			err = tx.Order("timestamp desc, cafeteriaRating desc").
				Limit(limit).
				Find(&ratings, "cafeteriaID = ? AND status = ?", canteenID, model.RatingStatusVisible).Error
		}

		if err != nil {
			log.WithError(err).Error("while querying last canteen ratings.")
			return make([]*pb.SingleRatingReply, 0)
		}
		ratingIDs := make([]int64, len(ratings))
		for i, v := range ratings {
			ratingIDs[i] = v.CanteenRating
		}
		tagsByRating := queryTagRatingsOverviewForRatings(ratingIDs, model.CANTEEN, tx)
		resp := make([]*pb.SingleRatingReply, 0, len(ratings))
		for _, v := range ratings {
			entry := &pb.SingleRatingReply{
				RatingId:    v.CanteenRating,
				Points:      v.Points,
				Comment:     v.Comment,
				Image:       getImageToBytes(v.Image),
				Visited:     timestamppb.New(v.Timestamp),
				RatingTags:  tagsByRating[v.CanteenRating],
				IsOwnRating: callerHash.Valid && v.DeviceHash.Valid && callerHash.String == v.DeviceHash.String,
			}
			if v.UpdatedAt.Valid {
				entry.UpdatedAt = timestamppb.New(v.UpdatedAt.Time)
			}
			resp = append(resp, entry)
		}
		return resp
	} else {
		return make([]*pb.SingleRatingReply, 0)
	}
}

func (s *CampusServer) GetDishRatings(ctx context.Context, input *pb.GetDishRatingsRequest) (*pb.GetDishRatingsReply, error) {
	s.recordDeviceIfPresent(ctx)
	if input.GetCanteenId() == "" || input.GetDish() == "" {
		return nil, status.Error(codes.NotFound, "No ratings found for this dish.")
	}
	tx := s.db.WithContext(ctx)
	canteenID, err := getIDForCanteenName(input.CanteenId, tx)
	if err != nil {
		return nil, err
	}
	dishID, err := getIDForDishName(input.Dish, canteenID, tx)
	if err != nil {
		return nil, err
	}

	var statsForDish model.DishRatingStatistic
	err = tx.First(&statsForDish, "cafeteriaID = ? AND dishID = ?", canteenID, dishID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "This dish has not yet been rated.")
	}
	if err != nil {
		fields := log.Fields{"dishID": dishID, "cafeteriaID": canteenID}
		log.WithError(err).WithFields(fields).Error("Error while querying the average ratings")
		return nil, status.Error(codes.Internal, "Could not query dish ratings.")
	}

	callerHash := callerDeviceHash(ctx)
	return &pb.GetDishRatingsReply{
		Avg:        statsForDish.Average,
		Std:        statsForDish.Std,
		Min:        statsForDish.Min,
		Max:        statsForDish.Max,
		Rating:     queryLastDishRatingsWithLimit(input, dishID, callerHash, tx),
		RatingTags: queryTags(canteenID, dishID, model.DISH, tx),
		NameTags:   queryTags(canteenID, dishID, model.NAME, tx),
	}, nil
}

// queryLastDishRatingsWithLimit
// Queries the actual ratings for a dish in a canteen and attaches the tag ratings which belong to the ratings
func queryLastDishRatingsWithLimit(input *pb.GetDishRatingsRequest, dishID int32, callerHash null.String, tx *gorm.DB) []*pb.SingleRatingReply {
	var ratings []model.DishRating
	var err error
	var limit = int(input.Limit)
	if limit == -1 {
		limit = math.MaxInt32
	}
	if limit > 0 {
		if input.From != nil || input.To != nil {
			var from time.Time
			var to time.Time
			if input.From == nil {
				from = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
			} else {
				from = input.From.AsTime()
			}

			if input.To == nil {
				to = time.Now()
			} else {
				to = input.To.AsTime()
			}

			err = tx.Order("timestamp desc, dishRating desc").
				Limit(limit).
				Find(&ratings, "dishID = ? AND status = ? AND timestamp < ? AND timestamp > ?", dishID, model.RatingStatusVisible, to, from).Error
		} else {
			err = tx.Order("timestamp desc, dishRating desc").
				Limit(limit).
				Find(&ratings, "dishID = ? AND status = ?", dishID, model.RatingStatusVisible).Error
		}

		if err != nil {
			log.WithError(err).Error("while querying last dish ratings from Database.")
			return make([]*pb.SingleRatingReply, 0)
		}
		ratingIDs := make([]int64, len(ratings))
		for i, v := range ratings {
			ratingIDs[i] = v.DishRating
		}
		tagsByRating := queryTagRatingsOverviewForRatings(ratingIDs, model.DISH, tx)
		resp := make([]*pb.SingleRatingReply, 0, len(ratings))
		for _, v := range ratings {
			entry := &pb.SingleRatingReply{
				RatingId:    v.DishRating,
				Points:      v.Points,
				Comment:     v.Comment,
				RatingTags:  tagsByRating[v.DishRating],
				Image:       getImageToBytes(v.Image),
				Visited:     timestamppb.New(v.Timestamp),
				IsOwnRating: callerHash.Valid && v.DeviceHash.Valid && callerHash.String == v.DeviceHash.String,
			}
			if v.UpdatedAt.Valid {
				entry.UpdatedAt = timestamppb.New(v.UpdatedAt.Time)
			}
			resp = append(resp, entry)
		}
		return resp
	} else {
		return make([]*pb.SingleRatingReply, 0)
	}
}

func getImageToBytes(imagePath string) []byte {
	if len(imagePath) == 0 {
		return make([]byte, 0)
	}
	imageAsBytes, err := os.ReadFile(imagePath)
	if err != nil {
		log.WithError(err).Error("while reading image file with path: ", imagePath)
		return nil
	}
	return imageAsBytes
}

type queryRatingTag struct {
	TagId   int32   `gorm:"column:tagId;type:int32;" json:"tagId"`
	Average float64 `gorm:"column:avg" json:"avg"`
	Std     float64 `json:"std"`
	Min     int32   `json:"min"`
	Max     int32   `json:"max"`
}

// queryTags
// Queries the average ratings for either canteenRatingTags, dishRatingTags or NameTags.
// Since the db only stores IDs in the results, the tags must be joined to retrieve their names form the rating_options tables.
func queryTags(canteenID int32, dishID int32, ratingType model.ModelType, tx *gorm.DB) []*pb.RatingTagResult {
	var results []queryRatingTag
	var err error
	switch ratingType {
	case model.DISH:
		err = tx.Table("dish_rating_tag_options options").
			Joins("JOIN dish_rating_tag_statistics results ON options.dishRatingTagOption = results.tagID").
			Select("options.dishRatingTagOption as tagId, results.average as avg, "+
				"results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ? AND results.dishID = ?", canteenID, dishID).
			Scan(&results).Error
	case model.CANTEEN:
		err = tx.Table("canteen_rating_tag_options options").
			Joins("JOIN cafeteria_rating_tag_statistics results ON options.cafeteriaRatingTagOption = results.tagID").
			Select("options.cafeteriaRatingTagOption as tagId, results.average as avg, "+
				"results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ?", canteenID).
			Scan(&results).Error
	default: //Query for name tags
		err = tx.Table("dish_to_dish_name_tags mapping").
			Where("mapping.dishID = ? AND results.cafeteriaID = ?", dishID, canteenID).
			Select("mapping.nameTagID as tag").
			Joins("JOIN dish_name_tag_statistics results ON mapping.nameTagID = results.tagID").
			Joins("JOIN dish_name_tag_options options ON mapping.nameTagID = options.dishNameTagOption").
			Select("mapping.nameTagID as tagId, results.average as avg, " +
				"results.min as min, results.max as max, results.std as std").
			Scan(&results).Error
	}

	if err != nil {
		log.WithError(err).Error("while querying the tags for the request")
	}

	//needed since the gRPC element does not specify column names - cannot be directly queried into the grpc message object.
	var resp []*pb.RatingTagResult
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

type ratingTagRow struct {
	ParentID int64 `gorm:"column:parentId"`
	TagID    int32 `gorm:"column:tagId"`
	Points   int32 `gorm:"column:points"`
}

// queryTagRatingsOverviewForRatings batch-loads all rating tags for a set of rating IDs in one query.
func queryTagRatingsOverviewForRatings(ratingIDs []int64, ratingType model.ModelType, tx *gorm.DB) map[int64][]*pb.RatingTagNewRequest {
	byRating := make(map[int64][]*pb.RatingTagNewRequest, len(ratingIDs))
	if len(ratingIDs) == 0 {
		return byRating
	}
	var rows []ratingTagRow
	var err error
	if ratingType == model.DISH {
		err = tx.Table("dish_rating_tag_options options").
			Joins("JOIN dish_rating_tags rating ON options.dishRatingTagOption = rating.tagID").
			Select("rating.parentRating as parentId, options.dishRatingTagOption as tagId, rating.points as points").
			Where("rating.parentRating IN ?", ratingIDs).
			Scan(&rows).Error
	} else {
		err = tx.Table("canteen_rating_tag_options options").
			Joins("JOIN canteen_rating_tags rating ON options.cafeteriaRatingTagOption = rating.tagID").
			Select("rating.correspondingRating as parentId, options.cafeteriaRatingTagOption as tagId, rating.points as points").
			Where("rating.correspondingRating IN ?", ratingIDs).
			Scan(&rows).Error
	}
	if err != nil {
		log.WithError(err).Error("while batch querying tag rating overviews.")
		return byRating
	}
	for _, r := range rows {
		byRating[r.ParentID] = append(byRating[r.ParentID], &pb.RatingTagNewRequest{
			TagId:  r.TagID,
			Points: r.Points,
		})
	}
	return byRating
}

// queryTagRatingsOverviewForRating
// Query all rating tags which belong to a specific rating given with an ID and return it as TagRatingOverviews
func queryTagRatingsOverviewForRating(dishID int64, ratingType model.ModelType, tx *gorm.DB) []*pb.RatingTagNewRequest {
	var results []*pb.RatingTagNewRequest
	var err error
	if ratingType == model.DISH {
		err = tx.Table("dish_rating_tag_options options").
			Joins("JOIN dish_rating_tags rating ON options.dishRatingTagOption = rating.tagID").
			Select("dishRatingTagOption as tagId, points, parentRating").
			Find(&results, "parentRating = ?", dishID).Error
	} else {
		err = tx.Table("canteen_rating_tag_options options").
			Joins("JOIN canteen_rating_tags rating ON options.cafeteriaRatingTagOption = rating.tagID").
			Select("cafeteriaRatingTagOption as tagId, points, correspondingRating").
			Find(&results, "correspondingRating = ?", dishID).Error
	}

	if err != nil {
		log.WithError(err).Error("while querying the tag rating overview.")
	}
	return results
}

// CreateCanteenRating RPC Endpoint
// Allows to store a new canteen rating atomically after validating all fields, image bytes, and tags.
func (s *CampusServer) CreateCanteenRating(ctx context.Context, input *pb.CreateCanteenRatingRequest) (*pb.CreateCanteenRatingReply, error) {
	s.recordDeviceIfPresent(ctx)
	if len(input.Image) > 0 && !ratingImagesEnabled() {
		return nil, status.Error(codes.FailedPrecondition, "Rating image uploads are currently disabled.")
	}
	if err := validateRatingTagsInput(input.RatingTags); err != nil {
		return nil, err
	}
	tx := s.db.WithContext(ctx)
	deviceHash := callerDeviceHash(ctx)
	if err := checkDeviceNotBanned(tx, deviceHash); err != nil {
		return nil, err
	}
	canteenID, errorRes := inputSanitizationForNewRatingElements(input.Points, input.Comment, input.CanteenId, tx)
	if errorRes != nil {
		return nil, errorRes
	}
	if err := verifyRatingTagsExist(input.RatingTags, model.CANTEEN, tx); err != nil {
		return nil, err
	}
	resPath, err := imageWrapper(input.Image, "cafeterias", canteenID)
	if err != nil {
		return nil, err
	}

	var ratingID int64
	if err := tx.Transaction(func(dbTx *gorm.DB) error {
		if deviceHash.Valid {
			now := time.Now()
			cooldownCutoff := now.AddDate(0, 0, -canteenRatingCooldownDays)
			var existing int64
			if err := dbTx.Model(&model.CanteenRating{}).
				Where("cafeteriaID = ? AND device_hash = ? AND timestamp >= ?",
					canteenID, deviceHash.String, cooldownCutoff).
				Count(&existing).Error; err != nil {
				log.WithError(err).Error("Error checking duplicate canteen rating.")
				return status.Error(codes.Internal, "Could not validate rating uniqueness. Rating has not been saved.")
			}
			if existing > 0 {
				return status.Errorf(codes.AlreadyExists, "You have already rated this canteen within the last %d days.", canteenRatingCooldownDays)
			}

			dayStart, dayEnd := berlinDayBounds(now)
			var dailyCount int64
			if err := dbTx.Model(&model.CanteenRating{}).
				Where("device_hash = ? AND timestamp >= ? AND timestamp < ?",
					deviceHash.String, dayStart, dayEnd).
				Count(&dailyCount).Error; err != nil {
				log.WithError(err).Error("Error checking daily canteen rating limit.")
				return status.Error(codes.Internal, "Could not validate daily rating limit. Rating has not been saved.")
			}
			if dailyCount >= maxCanteenRatingsPerDeviceDay {
				return status.Errorf(codes.ResourceExhausted, "You have reached the daily limit of %d canteen ratings.", maxCanteenRatingsPerDeviceDay)
			}
		}

		rating := model.CanteenRating{
			Comment:    input.Comment,
			Points:     input.Points,
			CanteenID:  canteenID,
			Timestamp:  time.Now(),
			Image:      resPath,
			DeviceHash: deviceHash,
			Status:     model.RatingStatusVisible,
		}
		if err := dbTx.Create(&rating).Error; err != nil {
			log.WithError(err).Error("Error occurred while creating the new canteen rating.")
			return status.Error(codes.Internal, "Error while creating new canteen rating. Rating has not been saved.")
		}
		ratingID = rating.CanteenRating
		return insertRatingTags(rating.CanteenRating, input.RatingTags, model.CANTEEN, dbTx)
	}); err != nil {
		cleanupUnusedCanteenRatingImage(tx, resPath)
		return nil, err
	}
	return &pb.CreateCanteenRatingReply{RatingId: ratingID}, nil
}

func isDuplicateKeyError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func imageWrapper(imageBytes []byte, subDir string, id int64) (string, error) {
	if len(imageBytes) == 0 {
		return "", nil
	}
	dir := path.Join(cron.StorageDir, "rating", subDir, strconv.FormatInt(id, 10)) + "/"
	return storeImage(dir, imageBytes)
}

func cleanupImageFile(imagePath string) {
	if imagePath == "" {
		return
	}
	if err := os.Remove(imagePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.WithError(err).WithField("path", imagePath).Warn("Failed to remove rating image")
	}
}

func cleanupUnusedCanteenRatingImage(db *gorm.DB, imagePath string) {
	if imagePath == "" {
		return
	}
	var refCount int64
	if err := db.Model(&model.CanteenRating{}).Where("image = ?", imagePath).Count(&refCount).Error; err != nil {
		log.WithError(err).WithField("path", imagePath).Warn("Could not verify canteen rating image references before cleanup")
		return
	}
	if refCount == 0 {
		cleanupImageFile(imagePath)
	}
}

func cleanupUnusedDishRatingImage(db *gorm.DB, imagePath string) {
	if imagePath == "" {
		return
	}
	var refCount int64
	if err := db.Model(&model.DishRating{}).Where("image = ?", imagePath).Count(&refCount).Error; err != nil {
		log.WithError(err).WithField("path", imagePath).Warn("Could not verify dish rating image references before cleanup")
		return
	}
	if refCount == 0 {
		cleanupImageFile(imagePath)
	}
}

// storeImage validates, resizes, and encodes an uploaded JPEG/PNG image as JPEG on disk.
func storeImage(dir string, raw []byte) (string, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil {
		return "", status.Error(codes.InvalidArgument, "Invalid image format. Only JPEG and PNG images are supported.")
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > maxImageDecodeDimension || cfg.Height > maxImageDecodeDimension {
		return "", status.Error(codes.InvalidArgument, "Image dimensions are out of allowed bounds.")
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return "", status.Error(codes.InvalidArgument, "Invalid image format. Only JPEG and PNG images are supported.")
	}

	resizedImage := img
	if img.Bounds().Dx() > maxRatingImageWidth {
		resizedImage = imaging.Resize(img, maxRatingImageWidth, 0, imaging.Lanczos)
	}

	quality := defaultJPEGQuality
	if len(raw) > targetMaxImageBytes {
		scaled := int(math.Round((float64(targetMaxImageBytes) / float64(len(raw))) * 100))
		quality = max(minJPEGQuality, min(defaultJPEGQuality, scaled))
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.WithError(err).WithField("path", dir).Error("Directory could not be created successfully")
		return "", status.Error(codes.Internal, "Could not store rating image.")
	}

	imgPath := path.Join(dir, fmt.Sprintf("%x.jpeg", md5.Sum(raw)))
	out, err := os.Create(imgPath)
	if err != nil {
		log.WithError(err).WithField("path", imgPath).Error("Error while creating rating image file")
		return "", status.Error(codes.Internal, "Could not store rating image.")
	}

	encodeErr := jpeg.Encode(out, resizedImage, &jpeg.Options{Quality: quality})
	closeErr := out.Close()
	if encodeErr != nil || closeErr != nil {
		_ = os.Remove(imgPath)
		log.WithError(errors.Join(encodeErr, closeErr)).WithField("path", imgPath).Error("Error while writing rating image")
		return "", status.Error(codes.Internal, "Could not store rating image.")
	}
	return imgPath, nil
}

// CreateDishRating RPC Endpoint
// Allows to store a new dish rating atomically after validating all fields, image bytes, and tags.
func (s *CampusServer) CreateDishRating(ctx context.Context, input *pb.CreateDishRatingRequest) (*pb.CreateDishRatingReply, error) {
	s.recordDeviceIfPresent(ctx)
	if len(input.Image) > 0 && !ratingImagesEnabled() {
		return nil, status.Error(codes.FailedPrecondition, "Rating image uploads are currently disabled.")
	}
	if err := validateRatingTagsInput(input.RatingTags); err != nil {
		return nil, err
	}
	tx := s.db.WithContext(ctx)
	deviceHash := callerDeviceHash(ctx)
	if err := checkDeviceNotBanned(tx, deviceHash); err != nil {
		return nil, err
	}
	canteenID, errorRes := inputSanitizationForNewRatingElements(input.Points, input.Comment, input.CanteenId, tx)
	if errorRes != nil {
		return nil, errorRes
	}

	var dishInCanteen model.Dish
	if err := tx.First(&dishInCanteen, "name = ? AND cafeteriaID = ?", input.Dish, canteenID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.InvalidArgument, "Dish is not offered in this week in this canteen. Rating has not been saved.")
		}
		log.WithError(err).Error("Error while looking up dish for new rating.")
		return nil, status.Error(codes.Internal, "Could not look up dish. Rating has not been saved.")
	}

	if err := verifyRatingTagsExist(input.RatingTags, model.DISH, tx); err != nil {
		return nil, err
	}

	resPath, err := imageWrapper(input.Image, "dishes", dishInCanteen.Dish)
	if err != nil {
		return nil, err
	}

	var ratingID int64
	if err := tx.Transaction(func(dbTx *gorm.DB) error {
		if deviceHash.Valid {
			now := time.Now()
			cooldownCutoff := now.AddDate(0, 0, -dishRatingCooldownDays)
			var existing int64
			if err := dbTx.Model(&model.DishRating{}).
				Where("dishID = ? AND device_hash = ? AND timestamp >= ?",
					dishInCanteen.Dish, deviceHash.String, cooldownCutoff).
				Count(&existing).Error; err != nil {
				log.WithError(err).Error("Error checking duplicate dish rating.")
				return status.Error(codes.Internal, "Could not validate rating uniqueness. Rating has not been saved.")
			}
			if existing > 0 {
				return status.Errorf(codes.AlreadyExists, "You have already rated this dish within the last %d days. You can update your existing rating instead.", dishRatingCooldownDays)
			}

			dayStart, dayEnd := berlinDayBounds(now)
			var dailyCount int64
			if err := dbTx.Model(&model.DishRating{}).
				Where("device_hash = ? AND timestamp >= ? AND timestamp < ?",
					deviceHash.String, dayStart, dayEnd).
				Count(&dailyCount).Error; err != nil {
				log.WithError(err).Error("Error checking daily dish rating limit.")
				return status.Error(codes.Internal, "Could not validate daily rating limit. Rating has not been saved.")
			}
			if dailyCount >= maxDishRatingsPerDeviceDay {
				return status.Errorf(codes.ResourceExhausted, "You have reached the daily limit of %d dish ratings.", maxDishRatingsPerDeviceDay)
			}
		}

		rating := model.DishRating{
			Comment:    input.Comment,
			DishID:     dishInCanteen.Dish,
			Points:     input.Points,
			Timestamp:  time.Now(),
			Image:      resPath,
			DeviceHash: deviceHash,
			Status:     model.RatingStatusVisible,
		}
		if err := dbTx.Create(&rating).Error; err != nil {
			log.WithError(err).Error("while creating a new dish rating.")
			return status.Error(codes.Internal, "Error while creating the new rating in the database. Rating has not been saved.")
		}
		ratingID = rating.DishRating

		if err := assignDishNameTag(&rating, dishInCanteen.Dish, dbTx); err != nil {
			return err
		}
		return insertRatingTags(rating.DishRating, input.RatingTags, model.DISH, dbTx)
	}); err != nil {
		cleanupUnusedDishRatingImage(tx, resPath)
		return nil, err
	}
	return &pb.CreateDishRatingReply{RatingId: ratingID}, nil
}

// UpdateDishRating allows the device that created a dish rating to modify its score, comment, tags, or photo.
func (s *CampusServer) UpdateDishRating(ctx context.Context, input *pb.UpdateDishRatingRequest) (*pb.UpdateDishRatingReply, error) {
	s.recordDeviceIfPresent(ctx)
	deviceHash := callerDeviceHash(ctx)
	if !deviceHash.Valid {
		return nil, status.Error(codes.Unauthenticated, "Device identification (x-device-id) is required to update a rating.")
	}
	if input.RatingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "A valid rating_id is required.")
	}
	if len(input.Image) > 0 && input.RemoveImage {
		return nil, status.Error(codes.InvalidArgument, "Cannot provide both image and remove_image.")
	}
	if len(input.Image) > 0 && !ratingImagesEnabled() {
		return nil, status.Error(codes.FailedPrecondition, "Rating image uploads are currently disabled.")
	}
	if err := validateRatingPointsAndComment(input.Points, input.Comment); err != nil {
		return nil, err
	}
	if err := validateRatingTagsInput(input.RatingTags); err != nil {
		return nil, err
	}

	tx := s.db.WithContext(ctx)
	if err := checkDeviceNotBanned(tx, deviceHash); err != nil {
		return nil, err
	}
	var existing model.DishRating
	if err := tx.First(&existing, "dishRating = ? AND status = ?", input.RatingId, model.RatingStatusVisible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Rating not found.")
		}
		log.WithError(err).Error("while loading dish rating for update")
		return nil, status.Error(codes.Internal, "Could not load rating.")
	}
	if !existing.DeviceHash.Valid || existing.DeviceHash.String != deviceHash.String {
		return nil, status.Error(codes.PermissionDenied, "You can only update your own ratings.")
	}
	if err := verifyRatingTagsExist(input.RatingTags, model.DISH, tx); err != nil {
		return nil, err
	}

	newImagePath := existing.Image
	oldImagePath := existing.Image
	wroteNewImage := false
	switch {
	case len(input.Image) > 0:
		stored, err := imageWrapper(input.Image, "dishes", existing.DishID)
		if err != nil {
			return nil, err
		}
		newImagePath = stored
		wroteNewImage = stored != oldImagePath
	case input.RemoveImage:
		newImagePath = ""
	}

	now := time.Now()
	if err := tx.Transaction(func(dbTx *gorm.DB) error {
		// Preserve original timestamp explicitly so MySQL ON UPDATE CURRENT_TIMESTAMP doesn't overwrite it.
		res := dbTx.Model(&model.DishRating{}).
			Where("dishRating = ? AND status = ?", existing.DishRating, model.RatingStatusVisible).
			Updates(map[string]any{
				"points":     input.Points,
				"comment":    input.Comment,
				"image":      newImagePath,
				"updated_at": now,
				"timestamp":  gorm.Expr("timestamp"),
			})
		if res.Error != nil {
			log.WithError(res.Error).Error("while updating dish rating")
			return status.Error(codes.Internal, "Could not update rating.")
		}
		if res.RowsAffected == 0 {
			return status.Error(codes.NotFound, "Rating not found.")
		}
		if err := dbTx.Where("parentRating = ?", existing.DishRating).Delete(&model.DishRatingTag{}).Error; err != nil {
			log.WithError(err).Error("while replacing dish rating tags")
			return status.Error(codes.Internal, "Could not update rating tags.")
		}
		if err := insertRatingTags(existing.DishRating, input.RatingTags, model.DISH, dbTx); err != nil {
			return err
		}
		if err := dbTx.Model(&model.DishNameTag{}).
			Where("correspondingRating = ?", existing.DishRating).
			Update("points", input.Points).Error; err != nil {
			log.WithError(err).Error("while updating dish name tag scores")
			return status.Error(codes.Internal, "Could not update dish name tag scores.")
		}
		return nil
	}); err != nil {
		if wroteNewImage {
			cleanupUnusedDishRatingImage(tx, newImagePath)
		}
		return nil, err
	}

	if oldImagePath != "" && oldImagePath != newImagePath {
		cleanupUnusedDishRatingImage(tx, oldImagePath)
	}
	return &pb.UpdateDishRatingReply{RatingId: existing.DishRating}, nil
}

// DeleteDishRating allows the device that created a visible dish rating to permanently remove it and its disk photo.
func (s *CampusServer) DeleteDishRating(ctx context.Context, input *pb.DeleteDishRatingRequest) (*pb.DeleteDishRatingReply, error) {
	s.recordDeviceIfPresent(ctx)
	deviceHash := callerDeviceHash(ctx)
	if !deviceHash.Valid {
		return nil, status.Error(codes.Unauthenticated, "Device identification (x-device-id) is required to delete a rating.")
	}
	if input.RatingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "A valid rating_id is required.")
	}

	tx := s.db.WithContext(ctx)
	var existing model.DishRating
	if err := tx.First(&existing, "dishRating = ? AND status = ?", input.RatingId, model.RatingStatusVisible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Rating not found.")
		}
		log.WithError(err).Error("while loading dish rating for deletion")
		return nil, status.Error(codes.Internal, "Could not load rating.")
	}
	if !existing.DeviceHash.Valid || existing.DeviceHash.String != deviceHash.String {
		return nil, status.Error(codes.PermissionDenied, "You can only delete your own ratings.")
	}

	if err := tx.Transaction(func(dbTx *gorm.DB) error {
		res := dbTx.Where("dishRating = ? AND status = ?", existing.DishRating, model.RatingStatusVisible).Delete(&model.DishRating{})
		if res.Error != nil {
			log.WithError(res.Error).Error("while deleting dish rating")
			return status.Error(codes.Internal, "Could not delete rating.")
		}
		if res.RowsAffected == 0 {
			return status.Error(codes.NotFound, "Rating not found.")
		}
		if err := dbTx.Where("correspondingRating = ?", existing.DishRating).Delete(&model.DishNameTag{}).Error; err != nil {
			log.WithError(err).Error("while deleting dish rating name tags")
			return status.Error(codes.Internal, "Could not delete rating name tags.")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	cleanupUnusedDishRatingImage(tx, existing.Image)
	return &pb.DeleteDishRatingReply{}, nil
}

// ReportDishRating lets a device report an abusive dish rating; once enough distinct devices report it, the rating is auto-hidden.
func (s *CampusServer) ReportDishRating(ctx context.Context, input *pb.ReportDishRatingRequest) (*pb.ReportDishRatingReply, error) {
	s.recordDeviceIfPresent(ctx)
	deviceHash := callerDeviceHash(ctx)
	if !deviceHash.Valid {
		return nil, status.Error(codes.Unauthenticated, "Device identification (x-device-id) is required to report a rating.")
	}
	if input.RatingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "A valid rating_id is required.")
	}
	reason := strings.TrimSpace(input.Reason)
	if utf8.RuneCountInString(reason) > maxRatingCommentRunes {
		return nil, status.Error(codes.InvalidArgument, "Report reason can only contain up to 256 characters.")
	}

	tx := s.db.WithContext(ctx)
	if err := checkDeviceNotBanned(tx, deviceHash); err != nil {
		return nil, err
	}
	var hidden bool
	var alert moderationReportAlert
	if err := tx.Transaction(func(dbTx *gorm.DB) error {
		var rating model.DishRating
		if err := dbTx.First(&rating, "dishRating = ? AND status = ?", input.RatingId, model.RatingStatusVisible).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Error(codes.NotFound, "Rating not found.")
			}
			log.WithError(err).Error("while loading dish rating for report")
			return status.Error(codes.Internal, "Could not load rating.")
		}
		if rating.DeviceHash.Valid && rating.DeviceHash.String == deviceHash.String {
			return status.Error(codes.InvalidArgument, "You cannot report your own rating.")
		}
		var existingReportCount int64
		if err := dbTx.Model(&model.DishRatingReport{}).
			Where("rating_id = ? AND device_hash = ?", rating.DishRating, deviceHash.String).
			Count(&existingReportCount).Error; err != nil {
			log.WithError(err).Error("while checking existing dish rating report")
			return status.Error(codes.Internal, "Could not check existing report.")
		}
		if existingReportCount > 0 {
			return status.Error(codes.AlreadyExists, "You have already reported this rating.")
		}

		report := model.DishRatingReport{
			RatingID:   rating.DishRating,
			DeviceHash: deviceHash.String,
			Reason:     reason,
			CreatedAt:  time.Now(),
		}
		if err := dbTx.Create(&report).Error; err != nil {
			if isDuplicateKeyError(err) {
				return status.Error(codes.AlreadyExists, "You have already reported this rating.")
			}
			log.WithError(err).Error("while creating dish rating report")
			return status.Error(codes.Internal, "Could not submit report.")
		}

		var totalReports int64
		if err := dbTx.Model(&model.DishRatingReport{}).
			Where("rating_id = ?", rating.DishRating).
			Count(&totalReports).Error; err != nil {
			log.WithError(err).Error("while counting dish rating reports")
			return status.Error(codes.Internal, "Could not count reports.")
		}
		if totalReports >= ratingAutoHideReportThreshold {
			if err := dbTx.Model(&model.DishRating{}).
				Where("dishRating = ?", rating.DishRating).
				Updates(map[string]any{
					"status":    model.RatingStatusHidden,
					"timestamp": gorm.Expr("timestamp"),
				}).Error; err != nil {
				log.WithError(err).Error("while auto-hiding reported dish rating")
				return status.Error(codes.Internal, "Could not update rating moderation status.")
			}
			hidden = true
		}

		var dishMeta struct {
			DishName    string `gorm:"column:dish_name"`
			CanteenName string `gorm:"column:canteen_name"`
		}
		_ = dbTx.Table("dishes d").
			Joins("JOIN canteens c ON c.cafeteria = d.cafeteriaID").
			Select("d.name AS dish_name, c.name AS canteen_name").
			Where("d.dish = ?", rating.DishID).
			Scan(&dishMeta).Error
		alert = moderationReportAlert{
			RatingID:     rating.DishRating,
			DishName:     dishMeta.DishName,
			CanteenID:    model.CanteenEatAPIID(dishMeta.CanteenName),
			Points:       rating.Points,
			Comment:      rating.Comment,
			Reason:       reason,
			TotalReports: totalReports,
			AutoHidden:   hidden,
			DeviceHash:   rating.DeviceHash.String,
		}
		return nil
	}); err != nil {
		return nil, err
	}
	s.notifyModerationDiscordAsync(alert)
	return &pb.ReportDishRatingReply{Hidden: hidden}, nil
}

// assignDishNameTag
// Query all name tags for this specific dish and generate the DishNameTag Ratings for each name tag
func assignDishNameTag(rating *model.DishRating, dishID int64, tx *gorm.DB) error {
	var nameTagIDs []int64
	if err := tx.Model(&model.DishToDishNameTag{}).
		Where("dishID = ?", dishID).
		Select("nameTagID").
		Scan(&nameTagIDs).Error; err != nil {
		log.WithError(err).Error("while loading dish name tags for dish ", dishID)
		return status.Error(codes.Internal, "Error while saving dish name tags. Rating has not been saved.")
	}
	for _, tagID := range nameTagIDs {
		if err := tx.Create(&model.DishNameTag{
			RatingID:  rating.DishRating,
			Points:    rating.Points,
			TagNameID: tagID,
		}).Error; err != nil {
			log.WithError(err).Error("while creating a new dish name rating.")
			return status.Error(codes.Internal, "Error while saving dish name tags. Rating has not been saved.")
		}
	}
	return nil
}

func validateRatingPointsAndComment(rating int32, comment string) error {
	if rating > 5 || rating < 1 {
		return status.Error(codes.InvalidArgument, "Rating must be a positive number not larger than 5. Rating has not been saved.")
	}

	if utf8.RuneCountInString(comment) > maxRatingCommentRunes {
		return status.Error(codes.InvalidArgument, "Ratings can only contain up to 256 characters, this is too long. Rating has not been saved.")
	}

	if strings.Contains(comment, "@") {
		return status.Error(codes.InvalidArgument, "Comments must not contain @ symbols in order to prevent misuse. Rating has not been saved.")
	}
	return nil
}

// inputSanitizationForNewRatingElements Checks parameters of the new rating for all canteen and dish ratings.
// Additionally, queries the canteen ID, since it checks whether the canteen actually exists.
func inputSanitizationForNewRatingElements(rating int32, comment string, canteenName string, tx *gorm.DB) (int64, error) {
	if err := validateRatingPointsAndComment(rating, comment); err != nil {
		return -1, err
	}

	canteenName = model.CanteenEnumName(canteenName)
	var result model.Canteen
	if res := tx.First(&result, "name = ?", canteenName); errors.Is(res.Error, gorm.ErrRecordNotFound) || res.RowsAffected == 0 {
		log.WithError(res.Error).Error("Error while querying the canteen id by name: ", canteenName)
		return -1, status.Error(codes.InvalidArgument, "Canteen does not exist. Rating has not been saved.")
	} else if res.Error != nil {
		log.WithError(res.Error).Error("Error while querying the canteen id by name: ", canteenName)
		return -1, status.Error(codes.Internal, "Could not verify canteen. Rating has not been saved.")
	}

	return result.CanteenID, nil
}

func validateRatingTagsInput(tags []*pb.RatingTag) error {
	if len(tags) == 0 {
		return nil
	}
	seen := make(map[int64]bool, len(tags))
	for _, t := range tags {
		if math.IsNaN(t.Points) || math.IsInf(t.Points, 0) || t.Points < 1 || t.Points > 5 {
			return status.Errorf(codes.InvalidArgument, "Tag rating points for tag %d must be between 1 and 5. Rating has not been saved.", t.TagId)
		}
		if seen[t.TagId] {
			return status.Errorf(codes.InvalidArgument, "The tag %d occurred more than once in this rating. Rating has not been saved.", t.TagId)
		}
		seen[t.TagId] = true
	}
	return nil
}

func verifyRatingTagsExist(tags []*pb.RatingTag, tagType model.ModelType, tx *gorm.DB) error {
	if len(tags) == 0 {
		return nil
	}
	tagIDs := make([]int64, len(tags))
	for i, t := range tags {
		tagIDs[i] = t.TagId
	}
	var existingIDs []int64
	var err error
	if tagType == model.DISH {
		err = tx.Model(&model.DishRatingTagOption{}).
			Where("dishRatingTagOption IN ?", tagIDs).
			Pluck("dishRatingTagOption", &existingIDs).Error
	} else {
		err = tx.Model(&model.CanteenRatingTagOption{}).
			Where("cafeteriaRatingTagOption IN ?", tagIDs).
			Pluck("cafeteriaRatingTagOption", &existingIDs).Error
	}
	if err != nil {
		log.WithError(err).Error("while verifying rating tags exist")
		return status.Error(codes.Internal, "Could not validate rating tags. Rating has not been saved.")
	}
	existingSet := make(map[int64]bool, len(existingIDs))
	for _, id := range existingIDs {
		existingSet[id] = true
	}
	var missing []string
	for _, id := range tagIDs {
		if !existingSet[id] {
			missing = append(missing, strconv.FormatInt(id, 10))
		}
	}
	if len(missing) > 0 {
		return status.Errorf(codes.InvalidArgument, "The tag(s) %s do not exist. Rating has not been saved.", strings.Join(missing, ", "))
	}
	return nil
}

func insertRatingTags(parentRatingID int64, tags []*pb.RatingTag, tagType model.ModelType, tx *gorm.DB) error {
	for _, currentTag := range tags {
		if tagType == model.DISH {
			if err := tx.Create(&model.DishRatingTag{
				RatingID: parentRatingID,
				Points:   int32(math.Round(currentTag.Points)),
				TagID:    currentTag.TagId,
			}).Error; err != nil {
				log.WithError(err).Error("while creating a dish tag rating")
				return status.Error(codes.Internal, "Error while saving rating tags. Rating has not been saved.")
			}
		} else {
			if err := tx.Create(&model.CanteenRatingTag{
				RatingID: parentRatingID,
				Points:   int32(math.Round(currentTag.Points)),
				TagID:    currentTag.TagId,
			}).Error; err != nil {
				log.WithError(err).Error("while creating a canteen tag rating")
				return status.Error(codes.Internal, "Error while saving rating tags. Rating has not been saved.")
			}
		}
	}
	return nil
}

func getIDForCanteenName(name string, tx *gorm.DB) (int32, error) {
	var result int32 = -1
	res := tx.Model(&model.Canteen{}).
		Where("name = ?", model.CanteenEnumName(name)).
		Select("cafeteria").
		Scan(&result)
	if res.Error != nil {
		log.WithError(res.Error).Error("while querying the canteen name.")
		return -1, status.Error(codes.Internal, "Could not query the canteen.")
	}
	if res.RowsAffected == 0 {
		return -1, status.Error(codes.NotFound, "No canteen with this Id found.")
	}
	return result, nil
}

func getIDForDishName(name string, canteenID int32, tx *gorm.DB) (int32, error) {
	var result int32 = -1
	res := tx.Model(&model.Dish{}).
		Where("name = ? AND cafeteriaID = ?", name, canteenID).
		Select("dish").
		Scan(&result)
	if res.Error != nil {
		log.WithError(res.Error).Error("while querying the dish name.")
		return -1, status.Error(codes.Internal, "Could not query the dish.")
	}
	if res.RowsAffected == 0 {
		return -1, status.Error(codes.NotFound, "No dish with this name found.")
	}

	return result, nil
}

// ListAvailableDishTags RPC Endpoint
// Returns all valid Tags to quickly rate dishes in english and german with the corresponding Id
func (s *CampusServer) ListAvailableDishTags(ctx context.Context, _ *pb.ListAvailableDishTagsRequest) (*pb.ListAvailableDishTagsReply, error) {
	var result []*pb.TagsOverview
	var requestStatus error = nil
	err := s.db.WithContext(ctx).
		Model(&model.DishRatingTagOption{}).
		Select("DE as de, EN as en, dishRatingTagOption as TagId").
		Find(&result).Error
	if err != nil {
		log.WithError(err).Error("while loading available dish tags from database.")
		requestStatus = status.Error(codes.Internal, "Available dish tags could not be loaded from the database.")
	}

	return &pb.ListAvailableDishTagsReply{
		RatingTags: result,
	}, requestStatus
}

// ListNameTags RPC Endpoint
// Returns all valid Tags to quickly rate dishes in english and german with the corresponding Id
func (s *CampusServer) ListNameTags(ctx context.Context, _ *pb.ListNameTagsRequest) (*pb.ListNameTagsReply, error) {
	var result []*pb.TagsOverview
	var requestStatus error = nil
	err := s.db.WithContext(ctx).
		Model(&model.DishNameTagOption{}).
		Select("DE as de, EN as en, dishNameTagOption as TagId").
		Find(&result).Error
	if err != nil {
		log.WithError(err).Error("while loading available Name Tags from database.")
		requestStatus = status.Error(codes.Internal, "Available dish tags could not be loaded from the database.")
	}

	return &pb.ListNameTagsReply{
		RatingTags: result,
	}, requestStatus
}

// ListAvailableCanteenTags RPC Endpoint
// Returns all valid Tags to quickly rate dishes in english and german
func (s *CampusServer) ListAvailableCanteenTags(ctx context.Context, _ *pb.ListAvailableCanteenTagsRequest) (*pb.ListAvailableCanteenTagsReply, error) {
	var result []*pb.TagsOverview
	var requestStatus error = nil
	err := s.db.WithContext(ctx).
		Model(&model.CanteenRatingTagOption{}).
		Select("DE as de, EN as en, cafeteriaRatingTagOption as TagId").
		Find(&result).Error
	if err != nil {
		log.WithError(err).Error("while loading available canteen tags from database.")
		requestStatus = status.Error(codes.Internal, "Available canteen tags could not be loaded from the database.")
	}

	return &pb.ListAvailableCanteenTagsReply{
		RatingTags: result,
	}, requestStatus
}

// ListCanteens RPC endpoint
// Returns all canteens with meta information which are available in the eat-api.
// The id is the eat-api canteen_id (e.g. "mensa-garching").
func (s *CampusServer) ListCanteens(ctx context.Context, _ *pb.ListCanteensRequest) (*pb.ListCanteensReply, error) {
	var canteens []model.Canteen
	if err := s.db.WithContext(ctx).Order("name").Find(&canteens).Error; err != nil {
		log.WithError(err).Error("while loading canteens from database.")
		return nil, status.Error(codes.Internal, "Canteens could not be loaded from the database.")
	}

	result := make([]*pb.Canteen, 0, len(canteens))
	for _, c := range canteens {
		result = append(result, &pb.Canteen{
			Id:        model.CanteenEatAPIID(c.Name),
			Address:   c.Address,
			Latitude:  c.Latitude,
			Longitude: c.Longitude,
		})
	}
	return &pb.ListCanteensReply{Canteen: result}, nil
}

func (s *CampusServer) ListDishes(ctx context.Context, req *pb.ListDishesRequest) (*pb.ListDishesReply, error) {
	s.recordDeviceIfPresent(ctx)
	if req.Year < 2022 {
		return &pb.ListDishesReply{}, status.Error(codes.InvalidArgument, "Years must be larger or equal to 2022 ") // currently, no previous values have been added
	}
	_, lastWeek := time.Date(int(req.Year), time.December, 28, 0, 0, 0, 0, time.UTC).ISOWeek()
	if req.Week < 1 || req.Week > int32(lastWeek) {
		return &pb.ListDishesReply{}, status.Errorf(codes.InvalidArgument, "Weeks must be in the range 1 - %d for ISO year %d", lastWeek, req.Year)
	}
	if req.Day < 0 || req.Day > 4 {
		return &pb.ListDishesReply{}, status.Error(codes.InvalidArgument, "Days must be in the range 0 (Monday) - 4 (Friday)")
	}

	var requestStatus error = nil
	var results []string
	canteenName := model.CanteenEnumName(req.CanteenId)

	err := s.db.WithContext(ctx).
		Table("dishes_of_the_weeks weekly").
		Where("weekly.day = ? AND weekly.week = ? and weekly.year = ?", req.Day, req.Week, req.Year).
		Select("weekly.dishID").
		Joins("JOIN dishes d ON d.dish = weekly.dishID").
		Joins("JOIN canteens c ON c.cafeteria = d.cafeteriaID").
		Where("c.name = ?", canteenName).
		Select("d.name").
		Find(&results).Error

	if err != nil {
		log.WithError(err).Error("while loading dishes from database.")
		requestStatus = status.Error(codes.Internal, "Dishes could not be loaded from the database.")
	}

	return &pb.ListDishesReply{
		Dish: results,
	}, requestStatus
}

// GetCanteenHeadCount RPC Endpoint
func (s *CampusServer) GetCanteenHeadCount(ctx context.Context, input *pb.GetCanteenHeadCountRequest) (*pb.GetCanteenHeadCountReply, error) {
	data := model.CanteenHeadCount{Count: 0, MaxCount: 0, Percent: -1} // Initialize with an empty (not found) value
	err := s.db.WithContext(ctx).
		Where(model.CanteenHeadCount{CanteenId: model.CanteenEatAPIID(input.CanteenId)}).
		FirstOrInit(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error("while querying the canteen head count for: ", input.CanteenId)
		return nil, status.Error(codes.Internal, "failed to query head count")
	}

	return &pb.GetCanteenHeadCountReply{
		Count:     data.Count,
		MaxCount:  data.MaxCount,
		Percent:   data.Percent,
		Timestamp: timestamppb.New(data.Timestamp),
	}, nil
}
