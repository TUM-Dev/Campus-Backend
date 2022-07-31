package backend

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"math"
	"os"
	"strings"
	"time"
)

type ModelType int

// Used to differentiate between the type of the model for different queries to reduce duplicated code.
const (
	DISH      ModelType = 1
	CAFETERIA ModelType = 2
	NAME      ModelType = 3
)

type QueryRatingTag struct {
	En      string  `gorm:"column:en;type:text;" json:"en"`
	De      string  `gorm:"column:De;type:text;" json:"De"`
	Average float64 `json:"average"`
	Std     float64 `json:"std"`
	Min     int32   `json:"min"`
	Max     int32   `json:"max"`
}

type QueryOverviewRatingTag struct {
	En     string `gorm:"column:En;type:mediumtext;" json:"En"`
	De     string `gorm:"column:De;type:mediumtext;" json:"De"`
	Points int32  `gorm:"column:points;type:text;"  json:"rating"`
}

// GetCafeteriaRatings
// RPC Endpoint
// Allows to query ratings for a specific cafeteria.
// It returns the average rating, max/min rating as well as a number of actual ratings and the average ratings for
// all cafeteria rating tags which were used to rate this cafeteria.
// The parameter limit defines how many actual ratings should be returned.
// The optional parameters from and to can define an interval in which the queried ratings have been stored.
// If these aren't specified, the newest ratings will be returned as the default
func (s *CampusServer) GetCafeteriaRatings(_ context.Context, input *pb.CafeteriaRatingRequest) (*pb.CafeteriaRatingResponse, error) {
	var result model.CafeteriaRatingAverage //get the average rating for this specific cafeteria
	cafeteriaId := getIDForCafeteriaName(input.CafeteriaName, s.db)
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

		return &pb.CafeteriaRatingResponse{
			AveragePoints:     float64(result.Average),
			StandardDeviation: float64(result.Std),
			MinPoints:         int32(result.Min),
			MaxPoints:         int32(result.Max),
			Rating:            ratings,
			RatingTags:        cafeteriaTags,
		}, nil
	} else {
		return &pb.CafeteriaRatingResponse{
			AveragePoints:     -1,
			StandardDeviation: -1,
			MinPoints:         -1,
			MaxPoints:         -1,
		}, nil
	}
}

// queryLastCafeteriaRatingsWithLimit
// Queries the actual ratings for a cafeteria and attaches the tag ratings which belong to the ratings
func queryLastCafeteriaRatingsWithLimit(input *pb.CafeteriaRatingRequest, cafeteriaID int32, s *CampusServer) []*pb.CafeteriaRating {
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
				Order("timestamp desc, id desc").
				Limit(limit).
				Find(&ratings).Error
		} else {
			err = s.db.Model(&model.CafeteriaRating{}).
				Where("cafeteriaID = ?", cafeteriaID).
				Order("timestamp desc, id desc").
				Limit(limit).
				Find(&ratings).Error
		}

		if err != nil {
			log.WithError(err).Error("Error while querying last cafeteria ratings.")
			return make([]*pb.CafeteriaRating, 0)
		}
		ratingResults := make([]*pb.CafeteriaRating, len(ratings))

		for i, v := range ratings {

			tagRatings := queryTagRatingsOverviewForRating(s, v.CafeteriaRating, CAFETERIA)
			ratingResults[i] = &pb.CafeteriaRating{
				Points:             v.Points,
				CafeteriaName:      input.CafeteriaName,
				Comment:            v.Comment,
				Image:              getImageToBytes(v.Image),
				CafeteriaVisitedAt: timestamppb.New(v.Timestamp),
				TagRating:          tagRatings,
			}
		}
		return ratingResults
	} else {
		return make([]*pb.CafeteriaRating, 0)
	}
}

// GetDishRatings
// RPC Endpoint
// Allows to query ratings for a specific dish in a specific cafeteria.
// It returns the average rating, max/min rating as well as a number of actual ratings and the average ratings for
// all dish rating tags which were used to rate this dish in this cafeteria. Additionally, the average, max/min are
// returned for every name tag which matches the name of the dish.
// The parameter limit defines how many actual ratings should be returned.
// The optional parameters from and to can define a interval in which the queried ratings have been stored.
// If these aren't specified, the newest ratings will be returned as the default
func (s *CampusServer) GetDishRatings(_ context.Context, input *pb.DishRatingRequest) (*pb.DishRatingResponse, error) {
	var result model.DishRatingAverage //get the average rating for this specific dish
	cafeteriaID := getIDForCafeteriaName(input.CafeteriaName, s.db)
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

		return &pb.DishRatingResponse{
			AveragePoints:     float64(result.Average),
			StandardDeviation: float64(result.Std),
			MinPoints:         int32(result.Min),
			MaxPoints:         int32(result.Max),
			Rating:            ratings,
			RatingTags:        dishTags,
			NameTags:          nameTags,
		}, nil
	} else {
		return &pb.DishRatingResponse{
			AveragePoints:     -1,
			MinPoints:         -1,
			MaxPoints:         -1,
			StandardDeviation: -1,
		}, nil
	}

}

// queryLastDishRatingsWithLimit
// Queries the actual ratings for a dish in a cafeteria and attaches the tag ratings which belong to the ratings
func queryLastDishRatingsWithLimit(input *pb.DishRatingRequest, cafeteriaID int32, dishID int32, s *CampusServer) []*pb.DishRating {
	var ratings []model.DishRating
	var err error
	var limit = int(input.Limit)
	if limit > 100 {
		limit = 100
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
				Order("timestamp desc, id desc").
				Limit(limit).
				Find(&ratings).Error
		} else {
			err = s.db.Model(&model.DishRating{}).
				Where("cafeteriaID = ? AND dishID = ?", cafeteriaID, dishID).
				Order("timestamp desc, id desc").
				Limit(limit).
				Find(&ratings).Error
		}

		if err != nil {
			log.WithError(err).Error("Error while querying last dish ratings from Database.")
			return make([]*pb.DishRating, 0)
		}
		ratingResults := make([]*pb.DishRating, len(ratings))

		for i, v := range ratings {

			tagRatings := queryTagRatingsOverviewForRating(s, v.DishRating, DISH)
			ratingResults[i] = &pb.DishRating{
				Points:             v.Points,
				Dish:               input.Dish,
				CafeteriaName:      input.CafeteriaName,
				Comment:            v.Comment,
				TagRating:          tagRatings,
				Image:              getImageToBytes(v.Image),
				CafeteriaVisitedAt: timestamppb.New(v.Timestamp),
			}
		}
		return ratingResults
	} else {
		return make([]*pb.DishRating, 0)
	}
}

func getImageToBytes(path string) []byte {
	file, err := os.Open(path)

	if err != nil {
		log.WithError(err).Error("Error while opening image ffile with path: {}.", path)
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
	return imageAsBytes
}

//queryTags
// Queries the average ratings for either cafeteriaRatingTags, dishRatingTags or NameTags.
// Since the db only stores IDs in the results, the tags must be joined to retrieve their names form the rating_options tables.
func queryTags(db *gorm.DB, cafeteriaID int32, dishID int32, ratingType ModelType) []*pb.TagRatingsResult {
	var results []QueryRatingTag
	var err error
	if ratingType == DISH {
		err = db.Table("dish_rating_tag_option options").
			Joins("JOIN dish_rating_tag_result results ON options.id = results.tagID").
			Select("options.De as De, results.average as average, "+
				"options.En as En, results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ? AND results.dishID = ?", cafeteriaID, dishID).
			Scan(&results).Error
	} else if ratingType == CAFETERIA {
		err = db.Table("cafeteria_rating_tag_option options").
			Joins("JOIN cafeteria_rating_tag_result results ON options.id = results.tagID").
			Select("options.De as De, results.average as average, "+
				"options.En as En, results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ?", cafeteriaID).
			Scan(&results).Error
	} else { //Query for name tags
		err = db.Table("dish_to_dish_name_tag mapping").
			Where("mapping.dishID = ?", dishID).
			Select("mapping.nameTagID as tag").
			Joins("JOIN dish_name_tag_result results ON mapping.nameTagID = results.tagID").
			Joins("JOIN dish_name_tag_option options ON mapping.nameTagID = options.id").
			Select("options.De as De, results.average as average, " +
				"options.En as En, results.min as min, results.max as max, results.std as std").
			Scan(&results).Error
	}

	if err != nil {
		log.WithError(err).Error("Error while querying the tags for the request.")
	}

	elements := make([]*pb.TagRatingsResult, len(results)) //needed since the gRPC element does not specify column names - cannot be directly queried into the grpc message object.
	for i, v := range results {
		elements[i] = &pb.TagRatingsResult{
			de:                v.De,
			en:                v.En,
			AveragePoints:     v.Average,
			StandardDeviation: v.Std,
			MinPoints:         v.Min,
			MaxPoints:         v.Max,
		}
	}

	return elements
}

// queryTagRatingOverviewForRating
// Query all rating tags which belong to a specific rating given with an ID and return it as TagRatingOverviews
func queryTagRatingsOverviewForRating(s *CampusServer, dishID int32, ratingType int32) []*pb.TagRatingResult {
	var results []QueryOverviewRatingTag
	var err error
	if ratingType == DISH {
		err = s.db.Table("dish_rating_tag_option options").
			Joins("JOIN dish_rating_tag rating ON options.id = rating.tagID").
			Where("rating.parentRating = ?", dishID).
			Select("options.De as De, options.En as En, rating.rating as rating").
			Scan(&results).Error
	} else {
		err = s.db.Table("cafeteria_rating_tag_option options").
			Joins("JOIN cafeteria_rating_tag rating ON options.id = rating.tagID").
			Where("rating.parentRating = ?", dishID).
			Select("options.De as De, options.En as En, rating.rating as rating").
			Scan(&results).Error
	}

	if err != nil {
		log.WithError(err).Error("Error while querying th tag rating overview.")
	}
	elements := make([]*pb.TagRatingResult, len(results))
	for i, a := range results {
		elements[i] = &pb.TagRatingResult{
			en:     a.en,
			de:     a.de,
			Points: a.Points,
		}
	}
	return elements
}

//NewCafeteriaRating
// RPC Endpoint
// Allows to store a new cafeteria Rating.
// If one of the parameters is invalid, an error will be returned. Otherwise, the rating will be saved.
// All rating tags which were given with the new rating are stored if they are valid tags, if at least one tag was
// invalid, an error is returned, all valid ratings tags will be stored nevertheless. Either the german or the english name can be returned to successfully store tags
func (s *CampusServer) NewCafeteriaRating(_ context.Context, input *pb.NewCafeteriaRatingRequest) (*emptypb.Empty, error) {
	cafeteriaID, errorRes := inputSanitizationForNewRatingElements(input.Points, input.Image, input.Comment, input.CafeteriaName, s)
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
	return storeRatingTags(s, rating.CafeteriaRating, input.Tags, CAFETERIA)
}

func imageWrapper(image []byte, path string, id int32) string {
	var resPath = ""
	if image != nil && len(image) > 0 {
		var resError error
		path := fmt.Sprintf("%s%s%s%d%s", "../images/", path, "/", id, "/")
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
	currentTime := time.Now().Unix()
	var imgPath = fmt.Sprintf("%s%d%s", path, currentTime, ".jpeg")
	var f, _ = os.Stat(imgPath)
	var counter = 1
	for {
		if f == nil {
			break
		} else {
			imgPath = fmt.Sprintf("%s%d%s%d%s", path, currentTime, "v", counter, ".jpeg")
			counter++
			f, _ = os.Stat(imgPath)
		}
	}

	out, errFile := os.Create(imgPath)
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing the file.")
		}
	}(out)
	var opts jpeg.Options
	opts.Quality = 100
	errFile = jpeg.Encode(out, img, &opts)
	return imgPath, errFile
}

// NewDishRating
// RPC Endpoint
// Allows to store a new dish Rating.
// If one of the parameters is invalid, an error will be returned. Otherwise, the rating will be saved.
// The ratingNumber will be saved for each corresponding DishNameTag.
// All rating tags which were given with the new rating are stored if they are valid tags, if at least one tag was
// invalid, an error is returned, all valid ratings tags will be stored nevertheless. Either the german or the english name can be returned to successfully store tags
func (s *CampusServer) NewDishRating(_ context.Context, input *pb.NewDishRatingRequest) (*emptypb.Empty, error) {

	cafeteriaID, errorRes := inputSanitizationForNewRatingElements(input.Points, input.Image, input.Comment, input.CafeteriaName, s)
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
	return storeRatingTags(s, rating.DishRating, input.Tags, DISH)
}

// assignDishNameTag
// Query all name tags for this specific dish and generate the DishNameTag Ratings ffor each name tag
func assignDishNameTag(s *CampusServer, rating model.DishRating, dishID int32) {
	var result []int
	err := s.db.Model(&model.DishToDishNameTag{}).Where("dishID = ? ", dishID).Select("nameTagID").Scan(&result).Error
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

// inputSanitizationForNewRatingElements
// Checks parameters of the new rating for all cafeteria and dish ratings.
// Additionally, queries the cafeteria ID, since it checks whether the cafeteria actually exists.
func inputSanitizationForNewRatingElements(rating int32, image []byte, comment string, cafeteriaName string, s *CampusServer) (int32, error) {
	if rating > 5 || rating < 0 {
		return -1, status.Errorf(codes.InvalidArgument, "Rating must be a positive number not larger than 10. Rating has not been saved.")
	}

	if len(image) > 131100 {
		return -1, status.Errorf(codes.InvalidArgument, "Image must not be larger than 1MB. Rating has not been saved.")
	}

	if len(comment) > 256 {
		return -1, status.Errorf(codes.InvalidArgument, "Ratings can only contain up to 256 characters, this is too long. Rating has not been saved.")
	}

	if strings.Contains(comment, "@") {
		return -1, status.Errorf(codes.InvalidArgument, "Comments must not contain @ symbols in order to prevent misuse. Rating has not been saved.")
	}

	var result *model.Cafeteria
	err := s.db.Model(&model.Cafeteria{}).
		Where("name LIKE ?", cafeteriaName).
		First(&result).Error
	if err.Error != nil && result != nil {
		log.WithError(err).Error("Error while querying the cafeteria id by name: {}", cafeteriaName)
		return -1, status.Errorf(codes.InvalidArgument, "Cafeteria does not exist. Rating has not been saved.")
	}

	return result.Cafeteria, nil
}

// storeRatingTags
// Checks whether the rating-tag name is a valid option and if so,
// it will be saved with a reference to the rating
func storeRatingTags(s *CampusServer, parentRatingID int32, tags []*pb.TagRating, tagType int) (*emptypb.Empty, error) {
	var errorOccurred = ""
	var warningOccurred = ""
	if len(tags) > 0 {
		usedTagIds := make(map[int]int)
		insertModel := getModelStoreTag(tagType, s.db)
		for _, tag := range tags {
			var currentTag int

			exists := getModelStoreTagOption(tagType, s.db).
				Where("En LIKE ? OR De LIKE ?", tag.Tag, tag.Tag).
				Select("id").
				First(&currentTag)

			if exists.Error != nil {
				log.WithError(exists.Error).Error("Error while querying the cafeteria name.")
			} else if exists.RowsAffected == 0 {
				log.Info("Tag with tag name ", tag.Tag, "does not exist")
				errorOccurred = errorOccurred + ", " + tag.Tag
			} else {
				if usedTagIds[currentTag] == 0 {
					err := insertModel.
						Create(&model.DishRatingTag{
							CorrespondingRating: parentRatingID,
							Points:              int32(tag.Points),
							TagID:               currentTag}).Error
					if err != nil {
						log.WithError(err).Error("Error while Creating a tag rating for a new rating.")
					}

					usedTagIds[currentTag] = 1
				} else {
					warningOccurred = warningOccurred + ", " + tag.Tag
					log.Info("Each Rating tag must be used at most once in a rating. This tag rating was not stored.")
				}
			}
		}
	}

	if len(errorOccurred) > 0 && len(warningOccurred) > 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "The Tag(s) "+errorOccurred+" does not exist. Remaining rating was saved without this rating tag. The Tag(s) "+warningOccurred+" occurred more than once in this rating.")
	} else if len(errorOccurred) > 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "The Tag(s) "+errorOccurred+" does not exist. Remaining rating was saved without this rating tag.")
	} else if len(warningOccurred) > 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "The Tag(s) "+warningOccurred+" occurred more than once in this rating.")
	} else {
		return &emptypb.Empty{}, nil
	}

}

// getModelStoreTagOption
// Returns the db model of the option table to reduce code duplicates
func getModelStoreTagOption(tagType int, db *gorm.DB) *gorm.DB {
	if tagType == DISH {
		return db.Model(&model.DishRatingTagOption{})
	} else {
		return db.Model(&model.CafeteriaRatingTagOption{})
	}
}

func getModelStoreTag(tagType int, db *gorm.DB) *gorm.DB {
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

//GetAvailableDishTags
// RPC Endpoint
// Returns all valid Tags to quickly rate dishes in english and german with the corresponding Id
func (s *CampusServer) GetAvailableDishTags(_ context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	var result []*model.DishRatingTagOption
	var requestStatus error = nil
	err := s.db.Model(&model.DishRatingTagOption{}).Select("De, En").Find(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while loading Cafeterias from database.")
		requestStatus = status.Errorf(codes.Internal, "Available dish tags could not be loaded from the database.")
	}
	elements := make([]*pb.TagRatingOverview, len(result))
	for i, a := range result {
		elements[i] = &pb.TagRatingOverview{en: a.en, de: a.de}
	}

	return &pb.GetRatingTagsReply{
		Tags: elements,
	}, requestStatus
}

//GetAvailableCafeteriaTags
// RPC Endpoint
// Returns all valid Tags to quickly rate dishes in english and german
func (s *CampusServer) GetAvailableCafeteriaTags(_ context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	var result []*model.CafeteriaRatingTagOption
	var requestStatus error = nil
	err := s.db.Model(&model.CafeteriaRatingTagOption{}).Select("De,En").Find(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while loading Cafeterias from database.")
		requestStatus = status.Errorf(codes.Internal, "Available cafeteria tags could not be loaded from the database.")
	}

	elements := make([]*pb.TagRatingOverview, len(result))
	for i, a := range result {
		elements[i] = &pb.TagRatingOverview{en: a.en, de: a.de}
	}

	return &pb.GetRatingTagsReply{
		Tags: elements,
	}, requestStatus
}

// GetCafeterias
// RPC endpoint
// Returns all cafeterias with meta information which are available in the eat-api
func (s *CampusServer) GetCafeterias(_ context.Context, _ *emptypb.Empty) (*pb.GetCafeteriaResponse, error) {
	var result []*pb.Cafeteria
	var requestStatus error = nil
	err := s.db.Model(&model.Cafeteria{}).Select("name,address,latitude,longitude").Scan(&result).Error
	if err != nil {
		log.WithError(err).Error("Error while loading Cafeterias from database.")
		requestStatus = status.Errorf(codes.Internal, "Cafeterias could not be loaded from the database.")
	}
	ratingResults := make([]*pb.Cafeteria, len(result))

	for i, v := range result {
		ratingResults[i] = &pb.Cafeteria{
			Name:      v.Name,
			Address:   v.Address,
			Latitude:  v.Latitude,
			Longitude: v.Longitude,
		}
	}
	return &pb.GetCafeteriaResponse{
		Cafeteria: ratingResults,
	}, requestStatus
}
