package backend

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"os"
	"strings"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// ListCanteenRatings RPC Endpoint
// Allows to query ratings for a specific cafeteria.
// It returns the average rating, max/min rating as well as a number of actual ratings and the average ratings for
// all cafeteria rating tags which were used to rate this cafeteria.
// The parameter limit defines how many actual ratings should be returned.
// The optional parameters from and to can define an interval in which the queried ratings have been stored.
// If these aren't specified, the newest ratings will be returned as the default
func (s *CampusServer) ListCanteenRatings(ctx context.Context, input *pb.ListCanteenRatingsRequest) (*pb.ListCanteenRatingsReply, error) {
	var statsForCanteen model.CafeteriaRatingStatistic
	tx := s.db.WithContext(ctx)
	cafeteriaId := getIDForCafeteriaName(input.CanteenId, tx)
	err := tx.First(&statsForCanteen, "cafeteriaId = ?", cafeteriaId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "No cafeteria with this Id found.")
	}
	if err != nil {
		log.WithError(err).Error("Error while querying the cafeteria with Id ", cafeteriaId)
		return nil, status.Error(codes.Internal, "could not query the cafeteria with the given Id")
	}

	return &pb.ListCanteenRatingsReply{
		Avg:        statsForCanteen.Average,
		Std:        statsForCanteen.Std,
		Min:        statsForCanteen.Min,
		Max:        statsForCanteen.Max,
		Rating:     queryLastCafeteriaRatingsWithLimit(input, cafeteriaId, tx),
		RatingTags: queryTags(cafeteriaId, -1, model.CAFETERIA, tx),
	}, nil
}

// queryLastCafeteriaRatingsWithLimit
// Queries the actual ratings for a cafeteria and attaches the tag ratings which belong to the ratings
func queryLastCafeteriaRatingsWithLimit(input *pb.ListCanteenRatingsRequest, cafeteriaID int32, tx *gorm.DB) []*pb.SingleRatingReply {
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
				Find(&ratings, "cafeteriaID = ? AND timestamp < ? AND timestamp > ?", cafeteriaID, to, from).Error
		} else {
			err = tx.Order("timestamp desc, cafeteriaRating desc").
				Limit(limit).
				Find(&ratings, "cafeteriaID = ?", cafeteriaID).Error
		}

		if err != nil {
			log.WithError(err).Error("while querying last cafeteria ratings.")
			return make([]*pb.SingleRatingReply, 0)
		}
		var resp []*pb.SingleRatingReply
		for _, v := range ratings {
			resp = append(resp, &pb.SingleRatingReply{
				Points:     v.Points,
				Comment:    v.Comment,
				Image:      getImageToBytes(v.Image),
				Visited:    timestamppb.New(v.Timestamp),
				RatingTags: queryTagRatingsOverviewForRating(v.CafeteriaRating, model.CAFETERIA, tx),
			})
		}
		return resp
	} else {
		return make([]*pb.SingleRatingReply, 0)
	}
}

func (s *CampusServer) GetDishRatings(ctx context.Context, input *pb.GetDishRatingsRequest) (*pb.GetDishRatingsReply, error) {
	tx := s.db.WithContext(ctx)
	cafeteriaID := getIDForCafeteriaName(input.CanteenId, tx)
	dishID := getIDForDishName(input.Dish, cafeteriaID, tx)

	var statsForDish model.DishRatingStatistic
	err := tx.First(&statsForDish, "cafeteriaID = ? AND dishID = ?", cafeteriaID, dishID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "No cafeteria with this Id found.")
	}
	if err != nil {
		fields := log.Fields{"dishID": dishID, "cafeteriaID": cafeteriaID}
		log.WithError(err).WithFields(fields).Error("Error while querying the average ratings")
		return nil, status.Error(codes.Internal, "This dish has not yet been rated.")
	}

	return &pb.GetDishRatingsReply{
		Avg:        statsForDish.Average,
		Std:        statsForDish.Std,
		Min:        statsForDish.Min,
		Max:        statsForDish.Max,
		Rating:     queryLastDishRatingsWithLimit(input, cafeteriaID, dishID, tx),
		RatingTags: queryTags(cafeteriaID, dishID, model.DISH, tx),
		NameTags:   queryTags(cafeteriaID, dishID, model.NAME, tx),
	}, nil
}

// queryLastDishRatingsWithLimit
// Queries the actual ratings for a dish in a cafeteria and attaches the tag ratings which belong to the ratings
func queryLastDishRatingsWithLimit(input *pb.GetDishRatingsRequest, cafeteriaID int32, dishID int32, tx *gorm.DB) []*pb.SingleRatingReply {
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
				Find(&ratings, "cafeteriaID = ? AND dishID = ? AND timestamp < ? AND timestamp > ?", cafeteriaID, dishID, to, from).Error
		} else {
			err = tx.Order("timestamp desc, dishRating desc").
				Limit(limit).
				Find(&ratings, "cafeteriaID = ? AND dishID = ?", cafeteriaID, dishID).Error
		}

		if err != nil {
			log.WithError(err).Error("while querying last dish ratings from Database.")
			return make([]*pb.SingleRatingReply, 0)
		}
		var resp []*pb.SingleRatingReply
		for _, v := range ratings {
			resp = append(resp, &pb.SingleRatingReply{
				Points:     v.Points,
				Comment:    v.Comment,
				RatingTags: queryTagRatingsOverviewForRating(v.DishRating, model.DISH, tx),
				Image:      getImageToBytes(v.Image),
				Visited:    timestamppb.New(v.Timestamp),
			})
		}
		return resp
	} else {
		return make([]*pb.SingleRatingReply, 0)
	}
}

func getImageToBytes(path string) []byte {
	if len(path) == 0 {
		return make([]byte, 0)
	}
	file, err := os.Open(path)
	if err != nil {
		log.WithError(err).Error("while opening image file with path: ", path)
		return nil
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.WithError(err).Error("Unable to close the file for storing the image.")
		}
	}(file)

	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	imageAsBytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	if _, err = buffer.Read(imageAsBytes); err != nil {
		log.WithError(err).Error("while trying to read image as bytes")
		return nil
	}
	return imageAsBytes
}

type queryRatingTag struct {
	TagId   int32   `gorm:"column:tagId;type:int32;" json:"tagId"`
	Average float64 `json:"avg"`
	Std     float64 `json:"std"`
	Min     int32   `json:"min"`
	Max     int32   `json:"max"`
}

// queryTags
// Queries the average ratings for either cafeteriaRatingTags, dishRatingTags or NameTags.
// Since the db only stores IDs in the results, the tags must be joined to retrieve their names form the rating_options tables.
func queryTags(cafeteriaID int32, dishID int32, ratingType model.ModelType, tx *gorm.DB) []*pb.RatingTagResult {
	var results []queryRatingTag
	var err error
	if ratingType == model.DISH {
		err = tx.Table("dish_rating_tag_options options").
			Joins("JOIN dish_rating_tag_statistics results ON options.dishRatingTagOption = results.tagID").
			Select("options.dishRatingTagOption as tagId, results.average as avg, "+
				"results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ? AND results.dishID = ?", cafeteriaID, dishID).
			Scan(&results).Error
	} else if ratingType == model.CAFETERIA {
		err = tx.Table("cafeteria_rating_tag_options options").
			Joins("JOIN cafeteria_rating_tag_statistics results ON options.cafeteriaRatingTagOption = results.tagID").
			Select("options.cafeteriaRatingTagOption as tagId, results.average as avg, "+
				"results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ?", cafeteriaID).
			Scan(&results).Error
	} else { //Query for name tags
		err = tx.Table("dish_to_dish_name_tags mapping").
			Where("mapping.dishID = ?", dishID).
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

// queryTagRatingOverviewForRating
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
		err = tx.Table("cafeteria_rating_tag_options options").
			Joins("JOIN cafeteria_rating_tags rating ON options.cafeteriaRatingTagOption = rating.tagID").
			Select("cafeteriaRatingTagOption as tagId, points, correspondingRating").
			Find(&results, "correspondingRating = ?", dishID).Error
	}

	if err != nil {
		log.WithError(err).Error("while querying the tag rating overview.")
	}
	return results
}

// CreateCanteenRating RPC Endpoint
// Allows to store a new cafeteria Rating.
// If one of the parameters is invalid, an error will be returned. Otherwise, the rating will be saved.
// All rating tags which were given with the new rating are stored if they are valid tags, if at least one tag was
// invalid, an error is returned, all valid ratings tags will be stored nevertheless. Either the german or the english name can be returned to successfully store tags
func (s *CampusServer) CreateCanteenRating(ctx context.Context, input *pb.CreateCanteenRatingRequest) (*pb.CreateCanteenRatingReply, error) {
	tx := s.db.WithContext(ctx)
	cafeteriaID, errorRes := inputSanitizationForNewRatingElements(input.Points, input.Comment, input.CanteenId, tx)
	if errorRes != nil {
		return nil, errorRes
	}

	resPath := imageWrapper(input.Image, "cafeterias", cafeteriaID)
	rating := model.CanteenRating{
		Comment:     input.Comment,
		Points:      input.Points,
		CafeteriaID: cafeteriaID,
		Timestamp:   time.Now(),
		Image:       resPath,
	}
	if err := tx.Create(&rating).Error; err != nil {
		log.WithError(err).Error("Error occurred while creating the new cafeteria rating.")
		return nil, status.Error(codes.InvalidArgument, "Error while creating new cafeteria rating. Rating has not been saved.")

	}
	if err := storeRatingTags(rating.CafeteriaRating, input.RatingTags, model.CAFETERIA, tx); err != nil {
		return &pb.CreateCanteenRatingReply{}, err
	}
	return &pb.CreateCanteenRatingReply{}, nil
}

func imageWrapper(image []byte, path string, id int64) string {
	if len(image) == 0 {
		return ""
	}
	path = fmt.Sprintf("/Storage/rating/%s/%d/", path, id)
	resPath, err := storeImage(path, image)
	if err != nil {
		log.WithError(err).Error("Error occurred while storing the image.")
	}
	return resPath
}

// storeImage
// stores an image and returns the path to this image.
// if needed, a new directory will be created and the path is extended until it is unique
func storeImage(path string, i []byte) (string, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.WithError(err).WithField("path", path).Error("Directory could not be created successfully")
			return "", nil
		}
	}

	img, _, _ := image.Decode(bytes.NewReader(i))
	resizedImage := imaging.Resize(img, 1280, 0, imaging.Lanczos)

	quality := 100         // if image small enough use it directly
	maxImageSize := 524288 // 0.55MB
	if len(i) > maxImageSize {
		quality = (maxImageSize / len(i)) * 100
	}

	var imgPath = fmt.Sprintf("%s%x.jpeg", path, md5.Sum(i))

	out, err := os.Create(imgPath)
	if err != nil {
		log.WithError(err).Error("Error while creating a new file on the path: ", path)
		return imgPath, err
	}
	defer func(out *os.File) {
		if err := out.Close(); err != nil {
			log.WithError(err).Error("while closing the file.")
		}
	}(out)

	return imgPath, jpeg.Encode(out, resizedImage, &jpeg.Options{Quality: quality})
}

// CreateDishRating RPC Endpoint
// Allows to store a new dish Rating.
// If one of the parameters is invalid, an error will be returned. Otherwise, the rating will be saved.
// The ratingNumber will be saved for each corresponding DishNameTag.
// All rating tags which were given with the new rating are stored if they are valid tags, if at least one tag was
// invalid, an error is returned, all valid ratings tags will be stored nevertheless. Either the german or the english name can be returned to successfully store tags
func (s *CampusServer) CreateDishRating(ctx context.Context, input *pb.CreateDishRatingRequest) (*pb.CreateDishRatingReply, error) {
	tx := s.db.WithContext(ctx)
	cafeteriaID, errorRes := inputSanitizationForNewRatingElements(input.Points, input.Comment, input.CanteenId, tx)
	if errorRes != nil {
		return nil, errorRes
	}

	var dishInCafeteria *model.Dish
	if err := tx.First(&dishInCafeteria, "name LIKE ? AND cafeteriaID = ?", input.Dish, cafeteriaID).Error; err != nil || dishInCafeteria == nil {
		log.WithError(err).Error("Error while creating a new dishInCafeteria rating.")
		return nil, status.Error(codes.InvalidArgument, "Dish is not offered in this week in this canteen. Rating has not been saved.")
	}

	resPath := imageWrapper(input.Image, "dishes", dishInCafeteria.Dish)

	rating := model.DishRating{
		Comment:   input.Comment,
		DishID:    dishInCafeteria.Dish,
		Points:    input.Points,
		Timestamp: time.Now(),
		Image:     resPath,
	}
	if err := tx.Create(&rating).Error; err != nil {
		log.WithError(err).Error("while creating a new dishInCafeteria rating.")
		return nil, status.Error(codes.Internal, "Error while creating the new rating in the database. Rating has not been saved.")
	}

	assignDishNameTag(&rating, dishInCafeteria.Dish, tx)

	if err := storeRatingTags(rating.DishRating, input.RatingTags, model.DISH, tx); err != nil {
		return &pb.CreateDishRatingReply{}, err
	}
	return &pb.CreateDishRatingReply{}, nil
}

// assignDishNameTag
// Query all name tags for this specific dish and generate the DishNameTag Ratings ffor each name tag
func assignDishNameTag(rating *model.DishRating, dishID int64, tx *gorm.DB) {
	var nameTagIDs []int64
	err := tx.Model(&model.DishToDishNameTag{}).
		Where("dishID = ? ", dishID).
		Select("nameTagID").
		Scan(&nameTagIDs).Error
	if err != nil {
		log.WithError(err).Error("while loading the dishID for the given name.")
	} else {
		for _, tagID := range nameTagIDs {
			if err := tx.Create(&model.DishNameTag{
				RatingID:  rating.DishRating,
				Points:    rating.Points,
				TagNameID: tagID,
			}).Error; err != nil {
				log.WithError(err).Error("while creating a new dish name rating.")
			}
		}
	}
}

// inputSanitizationForNewRatingElements Checks parameters of the new rating for all cafeteria and dish ratings.
// Additionally, queries the cafeteria ID, since it checks whether the cafeteria actually exists.
func inputSanitizationForNewRatingElements(rating int32, comment string, cafeteriaName string, tx *gorm.DB) (int64, error) {
	if rating > 5 || rating < 0 {
		return -1, status.Error(codes.InvalidArgument, "Rating must be a positive number not larger than 5. Rating has not been saved.")
	}

	if len(comment) > 256 {
		return -1, status.Error(codes.InvalidArgument, "Ratings can only contain up to 256 characters, this is too long. Rating has not been saved.")
	}

	if strings.Contains(comment, "@") {
		return -1, status.Error(codes.InvalidArgument, "Comments must not contain @ symbols in order to prevent misuse. Rating has not been saved.")
	}

	var result *model.Canteen
	if res := tx.First(&result, "name LIKE ?", cafeteriaName); errors.Is(res.Error, gorm.ErrRecordNotFound) || res.RowsAffected == 0 {
		log.WithError(res.Error).Error("Error while querying the cafeteria id by name: ", cafeteriaName)
		return -1, status.Error(codes.InvalidArgument, "Canteen does not exist. Rating has not been saved.")
	}

	return result.Cafeteria, nil
}

// storeRatingTags
// Checks whether the rating-tag name is a valid option and if so,
// it will be saved with a reference to the rating
func storeRatingTags(parentRatingID int64, tags []*pb.RatingTag, tagType model.ModelType, tx *gorm.DB) error {
	var errorOccurred = ""
	var warningOccurred = ""
	if len(tags) > 0 {
		usedTagIds := make(map[int64]bool)
		for _, currentTag := range tags {
			var err error
			var count int64

			if tagType == model.DISH {
				err = tx.Model(&model.DishRatingTagOption{}).
					Where("dishRatingTagOption LIKE ?", currentTag.TagId).
					Count(&count).Error
			} else {
				err = tx.Model(&model.CanteenRatingTagOption{}).
					Where("cafeteriaRatingTagOption LIKE ?", currentTag.TagId).
					Count(&count).Error
			}

			if errors.Is(err, gorm.ErrRecordNotFound) || count == 0 {
				fields := log.Fields{
					"tagid": currentTag.TagId,
					"count": count,
				}
				log.WithFields(fields).Info("tag does not exist")
				errorOccurred = fmt.Sprintf("%s, %d", errorOccurred, currentTag.TagId)
			} else {
				if !usedTagIds[currentTag.TagId] {
					if tagType == model.DISH {
						if err := tx.Create(&model.DishRatingTag{
							RatingID: parentRatingID,
							Points:   int32(currentTag.Points),
							TagID:    currentTag.TagId,
						}).Error; err != nil {
							log.WithError(err).Error("while Creating a currentTag rating for a new rating.")
						}
					} else {
						if err := tx.Create(&model.CanteenRatingTag{
							RatingID: parentRatingID,
							Points:   int32(currentTag.Points),
							TagID:    currentTag.TagId,
						}).Error; err != nil {
							log.WithError(err).Error("while Creating a currentTag rating for a new rating.")
						}
					}
					usedTagIds[currentTag.TagId] = true

				} else {
					warningOccurred = fmt.Sprintf("%s, %d", warningOccurred, currentTag.TagId)
					log.Info("Each Rating currentTag must be used at most once in a rating. This currentTag rating was not stored.")
				}
			}

		}
	}

	if len(errorOccurred) > 0 && len(warningOccurred) > 0 {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("The tag(s) %s does not exist. Remaining rating was saved without this rating tag. The tag(s) %s occurred more than once in this rating.", errorOccurred, warningOccurred))
	} else if len(errorOccurred) > 0 {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("The tag(s) %s does not exist. Remaining rating was saved without this rating tag.", errorOccurred))
	} else if len(warningOccurred) > 0 {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("The tag(s) %s occurred more than once in this rating.", warningOccurred))
	} else {
		return nil
	}
}

func getIDForCafeteriaName(name string, tx *gorm.DB) int32 {
	var result int32 = -1
	err := tx.Model(&model.Canteen{}).
		Where("name LIKE ?", name).
		Select("cafeteria").
		Scan(&result).Error
	if err != nil {
		log.WithError(err).Error("while querying the cafeteria name.")
		result = -1
	}
	return result
}

func getIDForDishName(name string, cafeteriaID int32, tx *gorm.DB) int32 {
	var result int32 = -1
	err := tx.Model(&model.Dish{}).
		Where("name LIKE ? AND cafeteriaID = ?", name, cafeteriaID).
		Select("dish").
		Scan(&result).Error
	if err != nil {
		log.WithError(err).Error("while querying the dish name.")
		result = -1
	}

	return result
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
		log.WithError(err).Error("while loading Cafeterias from database.")
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

// GetAvailableCafeteriaTags  RPC Endpoint
// Returns all valid Tags to quickly rate dishes in english and german
func (s *CampusServer) GetAvailableCafeteriaTags(ctx context.Context, _ *pb.ListAvailableCanteenTagsRequest) (*pb.ListAvailableCanteenTagsReply, error) {
	var result []*pb.TagsOverview
	var requestStatus error = nil
	err := s.db.WithContext(ctx).
		Model(&model.CanteenRatingTagOption{}).
		Select("DE as de, EN as en, cafeteriaRatingsTagOption as TagId").
		Find(&result).Error
	if err != nil {
		log.WithError(err).Error("while loading Cafeterias from database.")
		requestStatus = status.Error(codes.Internal, "Available cafeteria tags could not be loaded from the database.")
	}

	return &pb.ListAvailableCanteenTagsReply{
		RatingTags: result,
	}, requestStatus
}

// GetCafeterias RPC endpoint
// Returns all cafeterias with meta information which are available in the eat-api
func (s *CampusServer) GetCafeterias(ctx context.Context, _ *pb.ListCanteensRequest) (*pb.ListCanteensReply, error) {
	var result []*pb.Canteen
	var requestStatus error = nil
	if err := s.db.WithContext(ctx).
		Model(&model.Canteen{}).
		Select("cafeteria as id,address,latitude,longitude").
		Scan(&result).Error; err != nil {
		log.WithError(err).Error("while loading Cafeterias from database.")
		requestStatus = status.Error(codes.Internal, "Cafeterias could not be loaded from the database.")
	}

	return &pb.ListCanteensReply{
		Canteen: result,
	}, requestStatus
}

func (s *CampusServer) ListDishes(ctx context.Context, req *pb.ListDishesRequest) (*pb.ListDishesReply, error) {
	if req.Year < 2022 {
		return &pb.ListDishesReply{}, status.Error(codes.Internal, "Years must be larger or equal to 2022 ") // currently, no previous values have been added
	}
	if req.Week < 1 || req.Week > 52 {
		return &pb.ListDishesReply{}, status.Error(codes.Internal, "Weeks must be in the range 1 - 52")
	}
	if req.Day < 0 || req.Day > 4 {
		return &pb.ListDishesReply{}, status.Error(codes.Internal, "Days must be in the range 0 (Monday) - 4 (Friday)")
	}

	var requestStatus error = nil
	var results []string
	// the eat api has two types of ids, the enum ids (uppercase, with `_`) and the ids (lowercase, with `-`)
	cafeteriaName := strings.ReplaceAll(strings.ToUpper(req.CanteenId), "-", "_")

	err := s.db.WithContext(ctx).
		Table("dishes_of_the_weeks weekly").
		Where("weekly.day = ? AND weekly.week = ? and weekly.year = ?", req.Day, req.Week, req.Year).
		Select("weekly.dishID").
		Joins("JOIN dish d ON d.dish = weekly.dishID").
		Joins("JOIN cafeteria c ON c.cafeteria = d.cafeteriaID").
		Where("c.name LIKE ?", cafeteriaName).
		Select("d.name").
		Find(&results).Error

	if err != nil {
		log.WithError(err).Error("while loading Cafeterias from database.")
		requestStatus = status.Error(codes.Internal, "Cafeterias could not be loaded from the database.")
	}

	return &pb.ListDishesReply{
		Dish: results,
	}, requestStatus
}

// GetCanteenHeadCount RPC Endpoint
func (s *CampusServer) GetCanteenHeadCount(ctx context.Context, input *pb.GetCanteenHeadCountRequest) (*pb.GetCanteenHeadCountReply, error) {
	data := model.CanteenHeadCount{Count: 0, MaxCount: 0, Percent: -1} // Initialize with an empty (not found) value
	err := s.db.WithContext(ctx).
		Where(model.CanteenHeadCount{CanteenId: input.CanteenId}).
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
