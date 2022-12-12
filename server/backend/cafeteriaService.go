package backend

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"image"
	"image/jpeg"
	"math"
	"os"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type modelType int

// Used to differentiate between the type of the model for different queries to reduce duplicated code.
const (
	DISH      modelType = 1
	CAFETERIA modelType = 2
	NAME      modelType = 3
)

// GetCafeteriaRatings RPC Endpoint
// Allows to query ratings for a specific cafeteria.
// It returns the average rating, max/min rating as well as a number of actual ratings and the average ratings for
// all cafeteria rating tags which were used to rate this cafeteria.
// The parameter limit defines how many actual ratings should be returned.
// The optional parameters from and to can define an interval in which the queried ratings have been stored.
// If these aren't specified, the newest ratings will be returned as the default
func (s *CampusServer) GetCafeteriaRatings(_ context.Context, input *pb.CafeteriaRatingRequest) (*pb.CafeteriaRatingReply, error) {
	var result model.CafeteriaRatingAverage //get the average rating for this specific cafeteria
	cafeteriaId := getIDForCafeteriaName(input.CafeteriaId, s.db)
	res := s.db.Model(&model.CafeteriaRatingAverage{}).
		Where("cafeteriaId = ?", cafeteriaId).
		First(&result)

	if res.Error != nil {
		log.WithError(res.Error).Error("Error while querying the cafeteria with Id {}", cafeteriaId)
		return nil, status.Errorf(codes.Internal, "This cafeteria has not yet been rated.")
	}

	if res.RowsAffected > 0 {
		ratings := queryLastCafeteriaRatingsWithLimit(input, cafeteriaId, s)
		cafeteriaTags := queryTags(s.db, cafeteriaId, -1, CAFETERIA)

		return &pb.CafeteriaRatingReply{
			Avg:        float64(result.Average),
			Std:        float64(result.Std),
			Min:        int32(result.Min),
			Max:        int32(result.Max),
			Rating:     ratings,
			RatingTags: cafeteriaTags,
		}, nil
	} else {
		return &pb.CafeteriaRatingReply{
			Avg: -1,
			Std: -1,
			Min: -1,
			Max: -1,
		}, nil
	}
}

// queryLastCafeteriaRatingsWithLimit
// Queries the actual ratings for a cafeteria and attaches the tag ratings which belong to the ratings
func queryLastCafeteriaRatingsWithLimit(input *pb.CafeteriaRatingRequest, cafeteriaID int32, s *CampusServer) []*pb.SingleRatingReply {
	var ratings []model.CafeteriaRating
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
			err = s.db.Model(&model.CafeteriaRating{}).
				Where("cafeteriaID = ? AND timestamp < ? AND timestamp > ?", cafeteriaID, to, from).
				Order("timestamp desc, cafeteriaRating desc").
				Limit(limit).
				Find(&ratings).Error
		} else {
			err = s.db.Model(&model.CafeteriaRating{}).
				Where("cafeteriaID = ?", cafeteriaID).
				Order("timestamp desc, cafeteriaRating desc").
				Limit(limit).
				Find(&ratings).Error
		}

		if err != nil {
			log.WithError(err).Error("Error while querying last cafeteria ratings.")
			return make([]*pb.SingleRatingReply, 0)
		}
		ratingResults := make([]*pb.SingleRatingReply, len(ratings))

		for i, v := range ratings {

			tagRatings := queryTagRatingsOverviewForRating(s, v.CafeteriaRating, CAFETERIA)
			ratingResults[i] = &pb.SingleRatingReply{
				Points:     v.Points,
				Comment:    v.Comment,
				Image:      getImageToBytes(v.Image),
				Visited:    timestamppb.New(v.Timestamp),
				RatingTags: tagRatings,
			}
		}
		return ratingResults
	} else {
		return make([]*pb.SingleRatingReply, 0)
	}
}

// GetDishRatings RPC Endpoint
// Allows to query ratings for a specific dish in a specific cafeteria.
// It returns the average rating, max/min rating as well as a number of actual ratings and the average ratings for
// all dish rating tags which were used to rate this dish in this cafeteria. Additionally, the average, max/min are
// returned for every name tag which matches the name of the dish.
// The parameter limit defines how many actual ratings should be returned.
// The optional parameters from and to can define a interval in which the queried ratings have been stored.
// If these aren't specified, the newest ratings will be returned as the default
func (s *CampusServer) GetDishRatings(_ context.Context, input *pb.DishRatingRequest) (*pb.DishRatingReply, error) {
	var result model.DishRatingAverage //get the average rating for this specific dish
	cafeteriaID := getIDForCafeteriaName(input.CafeteriaId, s.db)
	dishID := getIDForDishName(input.Dish, cafeteriaID, s.db)

	err := s.db.Model(&model.DishRatingAverage{}).
		Where("cafeteriaID = ? AND dishID = ?", cafeteriaID, dishID).
		First(&result)

	if err.Error != nil {
		log.WithError(err.Error).Error("Error while querying the average ratings for the dish {} in the cafeteria {}.", dishID, cafeteriaID)
		return nil, status.Errorf(codes.Internal, "This dish has not yet been rated.")
	}

	if err.RowsAffected > 0 {
		ratings := queryLastDishRatingsWithLimit(input, cafeteriaID, dishID, s)
		dishTags := queryTags(s.db, cafeteriaID, dishID, DISH)
		nameTags := queryTags(s.db, cafeteriaID, dishID, NAME)

		return &pb.DishRatingReply{
			Avg:        float64(result.Average),
			Std:        float64(result.Std),
			Min:        int32(result.Min),
			Max:        int32(result.Max),
			Rating:     ratings,
			RatingTags: dishTags,
			NameTags:   nameTags,
		}, nil
	} else {
		return &pb.DishRatingReply{
			Avg: -1,
			Min: -1,
			Max: -1,
			Std: -1,
		}, nil
	}

}

// queryLastDishRatingsWithLimit
// Queries the actual ratings for a dish in a cafeteria and attaches the tag ratings which belong to the ratings
func queryLastDishRatingsWithLimit(input *pb.DishRatingRequest, cafeteriaID int32, dishID int32, s *CampusServer) []*pb.SingleRatingReply {
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

			err = s.db.Model(&model.DishRating{}).
				Where("cafeteriaID = ? AND dishID = ? AND timestamp < ? AND timestamp > ?", cafeteriaID, dishID, to, from).
				Order("timestamp desc, dishRating desc").
				Limit(limit).
				Find(&ratings).Error
		} else {
			err = s.db.Model(&model.DishRating{}).
				Where("cafeteriaID = ? AND dishID = ?", cafeteriaID, dishID).
				Order("timestamp desc, dishRating desc").
				Limit(limit).
				Find(&ratings).Error
		}

		if err != nil {
			log.WithError(err).Error("Error while querying last dish ratings from Database.")
			return make([]*pb.SingleRatingReply, 0)
		}
		ratingResults := make([]*pb.SingleRatingReply, len(ratings))

		for i, v := range ratings {
			ratingResults[i] = &pb.SingleRatingReply{
				Points:     v.Points,
				Comment:    v.Comment,
				RatingTags: queryTagRatingsOverviewForRating(s, v.DishRating, DISH),
				Image:      getImageToBytes(v.Image),
				Visited:    timestamppb.New(v.Timestamp),
			}
		}
		return ratingResults
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
		log.WithError(err).Error("Error while opening image file with path: ", path)
		return nil
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.WithError(err).Error("Unable to close the file for storing the image.")
		}
	}(file)

	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	imageAsBytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(imageAsBytes)
	if err != nil {
		log.WithError(err).Error("Error while trying to read image as bytes")
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
func queryTags(db *gorm.DB, cafeteriaID int32, dishID int32, ratingType modelType) []*pb.RatingTagResult {
	var results []queryRatingTag
	var err error
	if ratingType == DISH {
		err = db.Table("dish_rating_tag_option options").
			Joins("JOIN dish_rating_tag_average results ON options.dishRatingTagOption = results.tagID").
			Select("options.dishRatingTagOption as tagId, results.average as avg, "+
				"results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ? AND results.dishID = ?", cafeteriaID, dishID).
			Scan(&results).Error
	} else if ratingType == CAFETERIA {
		err = db.Table("cafeteria_rating_tag_option options").
			Joins("JOIN cafeteria_rating_tag_average results ON options.cafeteriaRatingTagOption = results.tagID").
			Select("options.cafeteriaRatingTagOption as tagId, results.average as avg, "+
				"results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ?", cafeteriaID).
			Scan(&results).Error
	} else { //Query for name tags
		err = db.Table("dish_to_dish_name_tag mapping").
			Where("mapping.dishID = ?", dishID).
			Select("mapping.nameTagID as tag").
			Joins("JOIN dish_name_tag_average results ON mapping.nameTagID = results.tagID").
			Joins("JOIN dish_name_tag_option options ON mapping.nameTagID = options.dishNameTagOption").
			Select("mapping.nameTagID as tagId, results.average as avg, " +
				"results.min as min, results.max as max, results.std as std").
			Scan(&results).Error
	}

	if err != nil {
		log.WithError(err).Error("Error while querying the tags for the request.")
	}

	//needed since the gRPC element does not specify column names - cannot be directly queried into the grpc message object.
	elements := make([]*pb.RatingTagResult, len(results))
	for i, v := range results {
		elements[i] = &pb.RatingTagResult{
			TagId: v.TagId,
			Avg:   v.Average,
			Std:   v.Std,
			Min:   v.Min,
			Max:   v.Max,
		}
	}

	return elements
}

// queryTagRatingOverviewForRating
// Query all rating tags which belong to a specific rating given with an ID and return it as TagRatingOverviews
func queryTagRatingsOverviewForRating(s *CampusServer, dishID int32, ratingType modelType) []*pb.RatingTagNewRequest {
	var results []*pb.RatingTagNewRequest
	var err error
	if ratingType == DISH {
		err = s.db.Table("dish_rating_tag_option options").
			Joins("JOIN dish_rating_tag rating ON options.dishRatingTagOption = rating.tagID").
			Select("dishRatingTagOption as tagId, points, parentRating").
			Find(&results, "parentRating = ?", dishID).Error
	} else {
		err = s.db.Table("cafeteria_rating_tag_option options").
			Joins("JOIN cafeteria_rating_tag rating ON options.cafeteriaRatingTagOption = rating.tagID").
			Select("cafeteriaRatingTagOption as tagId, points, correspondingRating").
			Find(&results, "correspondingRating = ?", dishID).Error
	}

	if err != nil {
		log.WithError(err).Error("Error while querying the tag rating overview.")
	}
	return results
}

// NewCafeteriaRating RPC Endpoint
// Allows to store a new cafeteria Rating.
// If one of the parameters is invalid, an error will be returned. Otherwise, the rating will be saved.
// All rating tags which were given with the new rating are stored if they are valid tags, if at least one tag was
// invalid, an error is returned, all valid ratings tags will be stored nevertheless. Either the german or the english name can be returned to successfully store tags
func (s *CampusServer) NewCafeteriaRating(_ context.Context, input *pb.NewCafeteriaRatingRequest) (*emptypb.Empty, error) {
	cafeteriaID, errorRes := inputSanitizationForNewRatingElements(input.Points, input.Comment, input.CafeteriaId, s)
	if errorRes != nil {
		return nil, errorRes
	}

	resPath := imageWrapper(input.Image, "cafeterias", cafeteriaID)
	rating := model.CafeteriaRating{
		Comment:     input.Comment,
		Points:      input.Points,
		CafeteriaID: cafeteriaID,
		Timestamp:   time.Now(),
		Image:       resPath,
	}

	err := s.db.Model(&model.CafeteriaRating{}).Create(&rating).Error
	if err != nil {
		log.WithError(err).Error("Error occurred while creating the new cafeteria rating.")
		return nil, status.Errorf(codes.InvalidArgument, "Error while creating new cafeteria rating. Rating has not been saved.")

	}
	return storeRatingTags(s, rating.CafeteriaRating, input.RatingTags, CAFETERIA)
}

func imageWrapper(image []byte, path string, id int32) string {
	var resPath = ""
	if len(image) > 0 {
		var resError error
		path := fmt.Sprintf("%s%s%s%d%s", "/Storage/rating/", path, "/", id, "/")
		resPath, resError = storeImage(path, image)

		if resError != nil {
			log.WithError(resError).Error("Error occurred while storing the image.")
		}
	}
	return resPath
}

// storeImage
// stores an image and returns teh path to this image.
// if needed, a new directory will be created and the path is extended until it is unique
func storeImage(path string, i []byte) (string, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.WithError(err).Error("Directory with path {} could not be created successfully", path)
			return "", nil
		}
	}

	img, _, _ := image.Decode(bytes.NewReader(i))
	resizedImage := imaging.Resize(img, 1280, 0, imaging.Lanczos)

	var opts jpeg.Options
	maxImageSize := 524288 // 0.55MB
	if len(i) > maxImageSize {
		opts.Quality = (maxImageSize / len(i)) * 100
	} else {
		opts.Quality = 100 // if image small enough use it directly
	}

	var imgPath = fmt.Sprintf("%s%x.jpeg", path, md5.Sum(i))

	out, errFile := os.Create(imgPath)
	if errFile != nil {
		log.WithError(errFile).Error("Error while creating a new file on the path: ", path)
		return imgPath, errFile
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing the file.")
		}
	}(out)

	errFile = jpeg.Encode(out, resizedImage, &opts)
	return imgPath, errFile
}

// NewDishRating RPC Endpoint
// Allows to store a new dish Rating.
// If one of the parameters is invalid, an error will be returned. Otherwise, the rating will be saved.
// The ratingNumber will be saved for each corresponding DishNameTag.
// All rating tags which were given with the new rating are stored if they are valid tags, if at least one tag was
// invalid, an error is returned, all valid ratings tags will be stored nevertheless. Either the german or the english name can be returned to successfully store tags
func (s *CampusServer) NewDishRating(_ context.Context, input *pb.NewDishRatingRequest) (*emptypb.Empty, error) {

	cafeteriaID, errorRes := inputSanitizationForNewRatingElements(input.Points, input.Comment, input.CafeteriaId, s)
	if errorRes != nil {
		return nil, errorRes
	}

	var dish *model.Dish
	errDish := s.db.Model(&model.Dish{}). //Dish must exist in the given mensa
						Where("name LIKE ? AND cafeteriaID = ?", input.Dish, cafeteriaID).
						First(&dish).Error
	if errDish != nil || dish == nil {
		log.WithError(errDish).Error("Error while creating a new dish rating.")
		return nil, status.Errorf(codes.InvalidArgument, "Dish is not offered in this week in this canteen. Rating has not been saved.")
	}

	resPath := imageWrapper(input.Image, "dishes", dish.Dish)

	rating := model.DishRating{
		Comment:     input.Comment,
		CafeteriaID: cafeteriaID,
		DishID:      dish.Dish,
		Points:      input.Points,
		Timestamp:   time.Now(),
		Image:       resPath,
	}

	err := s.db.Model(&model.DishRating{}).Create(&rating).Error
	if err != nil {
		log.WithError(err).Error("Error while creating a new dish rating.")
		return nil, status.Errorf(codes.Internal, "Error while creating the new rating in the database. Rating has not been saved.")
	}

	assignDishNameTag(s, rating, dish.Dish)
	return storeRatingTags(s, rating.DishRating, input.RatingTags, DISH)
}

// assignDishNameTag
// Query all name tags for this specific dish and generate the DishNameTag Ratings ffor each name tag
func assignDishNameTag(s *CampusServer, rating model.DishRating, dishID int32) {
	var result []int
	err := s.db.Model(&model.DishToDishNameTag{}).
		Where("dishID = ? ", dishID).
		Select("nameTagID").
		Scan(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while loading the dishID for the given name.")
	} else {
		for _, tagID := range result {
			err := s.db.Model(&model.DishNameTag{}).Create(&model.DishNameTag{
				CorrespondingRating: rating.DishRating,
				Points:              rating.Points,
				TagNameID:           tagID,
			}).Error
			if err != nil {
				log.WithError(err).Error("Error while creating a new dish name rating.")
			}
		}
	}
}

// inputSanitizationForNewRatingElements Checks parameters of the new rating for all cafeteria and dish ratings.
// Additionally, queries the cafeteria ID, since it checks whether the cafeteria actually exists.
func inputSanitizationForNewRatingElements(rating int32, comment string, cafeteriaName string, s *CampusServer) (int32, error) {
	if rating > 5 || rating < 0 {
		return -1, status.Errorf(codes.InvalidArgument, "Rating must be a positive number not larger than 10. Rating has not been saved.")
	}

	if len(comment) > 256 {
		return -1, status.Errorf(codes.InvalidArgument, "Ratings can only contain up to 256 characters, this is too long. Rating has not been saved.")
	}

	if strings.Contains(comment, "@") {
		return -1, status.Errorf(codes.InvalidArgument, "Comments must not contain @ symbols in order to prevent misuse. Rating has not been saved.")
	}

	var result *model.Cafeteria
	res := s.db.Model(&model.Cafeteria{}).
		Where("name LIKE ?", cafeteriaName).
		First(&result)
	if res.Error == gorm.ErrRecordNotFound || res.RowsAffected == 0 {
		log.WithError(res.Error).Error("Error while querying the cafeteria id by name: {}", cafeteriaName)
		return -1, status.Errorf(codes.InvalidArgument, "Cafeteria does not exist. Rating has not been saved.")
	}

	return result.Cafeteria, nil
}

// storeRatingTags
// Checks whether the rating-tag name is a valid option and if so,
// it will be saved with a reference to the rating
func storeRatingTags(s *CampusServer, parentRatingID int32, tags []*pb.RatingTag, tagType modelType) (*emptypb.Empty, error) {
	var errorOccurred = ""
	var warningOccurred = ""
	if len(tags) > 0 {
		usedTagIds := make(map[int]int)
		insertModel := getModelStoreTag(tagType, s.db)
		for _, currentTag := range tags {
			var err error
			var count int64

			if tagType == DISH {
				err = s.db.Model(&model.DishRatingTagOption{}).
					Where("dishRatingTagOption LIKE ?", currentTag.TagId).
					Count(&count).Error
			} else {
				err = s.db.Model(&model.CafeteriaRatingTagOption{}).
					Where("cafeteriaRatingTagOption LIKE ?", currentTag.TagId).
					Count(&count).Error
			}

			if err == gorm.ErrRecordNotFound || count == 0 {
				log.Info("tag with tagid ", currentTag.TagId, "does not exist")
				errorOccurred = fmt.Sprintf("%s, %d", errorOccurred, currentTag.TagId)
			} else {
				if usedTagIds[int(currentTag.TagId)] == 0 {
					err := insertModel.
						Create(&model.DishRatingTag{
							CorrespondingRating: parentRatingID,
							Points:              int32(currentTag.Points),
							TagID:               int(currentTag.TagId),
						}).Error
					if err != nil {
						log.WithError(err).Error("Error while Creating a currentTag rating for a new rating.")
					}
					usedTagIds[int(currentTag.TagId)] = 1

				} else {
					warningOccurred = fmt.Sprintf("%s, %d", warningOccurred, currentTag.TagId)
					log.Info("Each Rating currentTag must be used at most once in a rating. This currentTag rating was not stored.")
				}
			}

		}
	}

	if len(errorOccurred) > 0 && len(warningOccurred) > 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "The tag(s) "+errorOccurred+" does not exist. Remaining rating was saved without this rating tag. The tag(s) "+warningOccurred+" occurred more than once in this rating.")
	} else if len(errorOccurred) > 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "The tag(s) "+errorOccurred+" does not exist. Remaining rating was saved without this rating tag.")
	} else if len(warningOccurred) > 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "The tag(s) "+warningOccurred+" occurred more than once in this rating.")
	} else {
		return &emptypb.Empty{}, nil
	}

}

func getModelStoreTag(tagType modelType, db *gorm.DB) *gorm.DB {
	if tagType == DISH {
		return db.Model(&model.DishRatingTag{})
	} else {
		return db.Model(&model.CafeteriaRatingTag{})
	}
}

func getIDForCafeteriaName(name string, db *gorm.DB) int32 {
	var result int32 = -1
	err := db.Model(&model.Cafeteria{}).
		Where("name LIKE ?", name).
		Select("cafeteria").
		Scan(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while querying the cafeteria name.")
		result = -1
	}
	return result
}

func getIDForDishName(name string, cafeteriaID int32, db *gorm.DB) int32 {
	var result int32 = -1
	err := db.Model(&model.Dish{}).
		Where("name LIKE ? AND cafeteriaID = ?", name, cafeteriaID).
		Select("dish").
		Scan(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while querying the dish name.")
		result = -1
	}

	return result
}

// GetAvailableDishTags RPC Endpoint
// Returns all valid Tags to quickly rate dishes in english and german with the corresponding Id
func (s *CampusServer) GetAvailableDishTags(_ context.Context, _ *emptypb.Empty) (*pb.GetTagsReply, error) {
	var result []*pb.TagsOverview
	var requestStatus error = nil
	err := s.db.Model(&model.DishRatingTagOption{}).Select("DE as de, EN as en, dishRatingTagOption as TagId").Find(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while loading Cafeterias from database.")
		requestStatus = status.Errorf(codes.Internal, "Available dish tags could not be loaded from the database.")
	}

	return &pb.GetTagsReply{
		RatingTags: result,
	}, requestStatus
}

// GetNameTags RPC Endpoint
// Returns all valid Tags to quickly rate dishes in english and german with the corresponding Id
func (s *CampusServer) GetNameTags(_ context.Context, _ *emptypb.Empty) (*pb.GetTagsReply, error) {
	var result []*pb.TagsOverview
	var requestStatus error = nil
	err := s.db.Model(&model.DishNameTagOption{}).Select("DE as de, EN as en, dishNameTagOption as TagId").Find(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while loading available Name Tags from database.")
		requestStatus = status.Errorf(codes.Internal, "Available dish tags could not be loaded from the database.")
	}

	return &pb.GetTagsReply{
		RatingTags: result,
	}, requestStatus
}

// GetAvailableCafeteriaTags  RPC Endpoint
// Returns all valid Tags to quickly rate dishes in english and german
func (s *CampusServer) GetAvailableCafeteriaTags(_ context.Context, _ *emptypb.Empty) (*pb.GetTagsReply, error) {
	var result []*pb.TagsOverview
	var requestStatus error = nil
	err := s.db.Model(&model.CafeteriaRatingTagOption{}).Select("DE as de, EN as en, cafeteriaRatingsTagOption as TagId").Find(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while loading Cafeterias from database.")
		requestStatus = status.Errorf(codes.Internal, "Available cafeteria tags could not be loaded from the database.")
	}

	return &pb.GetTagsReply{
		RatingTags: result,
	}, requestStatus
}

// GetCafeterias RPC endpoint
// Returns all cafeterias with meta information which are available in the eat-api
func (s *CampusServer) GetCafeterias(_ context.Context, _ *emptypb.Empty) (*pb.GetCafeteriaReply, error) {
	var result []*pb.Cafeteria
	var requestStatus error = nil
	err := s.db.Model(&model.Cafeteria{}).Select("cafeteria as id,address,latitude,longitude").Scan(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while loading Cafeterias from database.")
		requestStatus = status.Errorf(codes.Internal, "Cafeterias could not be loaded from the database.")
	}

	return &pb.GetCafeteriaReply{
		Cafeteria: result,
	}, requestStatus
}

func (s *CampusServer) GetDishes(_ context.Context, request *pb.GetDishesRequest) (*pb.GetDishesReply, error) {
	if request.Year < 2022 {
		return &pb.GetDishesReply{}, status.Errorf(codes.Internal, "Years must be larger or equal to 2022 ") // currently, no previous values have been added
	}
	if request.Week < 1 || request.Week > 53 {
		return &pb.GetDishesReply{}, status.Errorf(codes.Internal, "Weeks must be in the range 1 - 53")
	}
	if request.Day < 0 || request.Day > 4 {
		return &pb.GetDishesReply{}, status.Errorf(codes.Internal, "Days must be in the range 1 (Monday) - 4 (Friday)")
	}

	var requestStatus error = nil
	var results []string
	err := s.db.Table("dishes_of_the_week weekly").
		Where("weekly.day = ? AND weekly.week = ? and weekly.year = ?", request.Day, request.Week, request.Year).
		Select("weekly.dishID").
		Joins("JOIN dish d ON d.dish = weekly.dishID").
		Joins("JOIN cafeteria c ON c.cafeteria = d.cafeteriaID").
		Where("c.name LIKE ?", request.CafeteriaId).
		Select("d.name").
		Find(&results).Error

	if err != nil {
		log.WithError(err).Error("Error while loading Cafeterias from database.")
		requestStatus = status.Errorf(codes.Internal, "Cafeterias could not be loaded from the database.")
	}

	return &pb.GetDishesReply{
		Dish: results,
	}, requestStatus
}
