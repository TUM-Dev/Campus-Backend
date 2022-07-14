package backend

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"os"
	"strings"
	"time"
)

const MEAL = 1
const CAFETERIA = 2
const NAME = 3

type QueryRatingTag struct {
	EN      string  `gorm:"column:EN;type:text;" json:"EN"`
	DE      string  `gorm:"column:DE;type:text;" json:"DE"`
	Average float64 `json:"average"`
	Std     float64 `json:"std"`
	Min     int32   `json:"min"`
	Max     int32   `json:"max"`
}

type QueryOverviewRatingTag struct {
	EN     string `gorm:"column:EN;type:mediumtext;" json:"EN"`
	DE     string `gorm:"column:DE;type:mediumtext;" json:"DE"`
	Points int32  `gorm:"column:points;type:text;"  json:"rating"`
}

/*
RPC Endpoint
Allows to query ratings for a specific cafeteria.
It returns the average rating, max/min rating as well as a number of actual ratings and the average ratings for
all cafeteria rating tags which were used to rate this cafeteria.

The parameter limit defines how many actual ratings should be returned.

The optional parameters from and to can define a interval in which the queried ratings have been stored.
If these aren't specified, the newest ratings will be returnes as the default
*/
func (s *CampusServer) GetCafeteriaRatings(_ context.Context, input *pb.CafeteriaRatingRequest) (*pb.CafeteriaRatingResponse, error) {
	var result cafeteria_rating_models.CafeteriaRatingAverage //get the average rating for this specific cafeteria
	cafeteriaID := getIDForCafeteriaName(input.CafeteriaName, s.db)
	res := s.db.Model(&cafeteria_rating_models.CafeteriaRatingAverage{}).
		Where("cafeteriaID = ?", cafeteriaID).
		First(&result)

	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "This cafeteria has not yet been rated.")
	}

	if res.RowsAffected > 0 {
		ratings := queryLastCafeteriaRatingsWithLimit(input, cafeteriaID, s)
		cafeteriaTags := queryTags(s.db, cafeteriaID, -1, CAFETERIA)

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

/*
Queries the actual ratings for a cafeteria and attaches the tag ratings which belong to the ratings
*/
func queryLastCafeteriaRatingsWithLimit(input *pb.CafeteriaRatingRequest, cafeteriaID int32, s *CampusServer) []*pb.CafeteriaRating {
	var ratings []cafeteria_rating_models.CafeteriaRating
	var errRatings error
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
			errRatings = s.db.Model(&cafeteria_rating_models.CafeteriaRating{}).
				Where("cafeteriaID = ? AND timestamp < ? AND timestamp > ?", cafeteriaID, to, from).
				Order("timestamp desc, id desc").
				Limit(limit).
				Find(&ratings).Error
		} else {
			errRatings = s.db.Model(&cafeteria_rating_models.CafeteriaRating{}).
				Where("cafeteriaID = ?", cafeteriaID).
				Order("timestamp desc, id desc").
				Limit(limit).
				Find(&ratings).Error
		}

		if errRatings != nil {
			log.Error(errRatings)
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

/*
RPC Endpoint
Allows to query ratings for a specific meal in a specific cafeteria.
It returns the average rating, max/min rating as well as a number of actual ratings and the average ratings for
all meal rating tags which were used to rate this meal in this cafeteria. Additionally the average, max/min are
returned for every name tag which matches the nemae of the meal.

The parameter limit defines how many actual ratings should be returned.

The optional parameters from and to can define a interval in which the queried ratings have been stored.
If these aren't specified, the newest ratings will be returnes as the default
*/
func (s *CampusServer) GetMealRatings(_ context.Context, input *pb.MealRatingRequest) (*pb.MealRatingResponse, error) {
	var result cafeteria_rating_models.MealRatingAverage //get the average rating for this specific meal
	cafeteriaID := getIDForCafeteriaName(input.CafeteriaName, s.db)
	mealID := getIDForMealName(input.Meal, cafeteriaID, s.db)

	res := s.db.Model(&cafeteria_rating_models.MealRatingAverage{}).
		Where("cafeteriaID = ? AND mealID = ?", cafeteriaID, mealID).
		First(&result)

	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "This meal has not yet been rated.")
	}

	if res.RowsAffected > 0 {
		ratings := queryLastMealRatingsWithLimit(input, cafeteriaID, mealID, s)
		mealTags := queryTags(s.db, cafeteriaID, mealID, MEAL)
		nameTags := queryTags(s.db, cafeteriaID, mealID, NAME)

		return &pb.MealRatingResponse{
			AveragePoints:     float64(result.Average),
			StandardDeviation: float64(result.Std),
			MinPoints:         int32(result.Min),
			MaxPoints:         int32(result.Max),
			Rating:            ratings,
			RatingTags:        mealTags,
			NameTags:          nameTags,
		}, nil
	} else {
		return &pb.MealRatingResponse{
			AveragePoints:     -1,
			MinPoints:         -1,
			MaxPoints:         -1,
			StandardDeviation: -1,
		}, nil
	}

}

/*
Queries the actual ratings for a meal in a cafeteria and attaches the tag ratings which belong to the ratings
*/
func queryLastMealRatingsWithLimit(input *pb.MealRatingRequest, cafeteriaID int32, mealID int32, s *CampusServer) []*pb.MealRating {
	var ratings []cafeteria_rating_models.MealRating
	var errRatings error
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

			errRatings = s.db.Model(&cafeteria_rating_models.MealRating{}).
				Where("cafeteriaID = ? AND mealID = ? AND timestamp < ? AND timestamp > ?", cafeteriaID, mealID, to, from).
				Order("timestamp desc, id desc").
				Limit(limit).
				Find(&ratings).Error
		} else {
			errRatings = s.db.Model(&cafeteria_rating_models.MealRating{}).
				Where("cafeteriaID = ? AND mealID = ?", cafeteriaID, mealID).
				Order("timestamp desc, id desc").
				Limit(limit).
				Find(&ratings).Error
		}

		if errRatings != nil {
			return make([]*pb.MealRating, 0)
		}
		ratingResults := make([]*pb.MealRating, len(ratings))

		for i, v := range ratings {

			tagRatings := queryTagRatingsOverviewForRating(s, v.MealRating, MEAL)
			ratingResults[i] = &pb.MealRating{
				Points:             v.Points,
				Meal:               input.Meal,
				CafeteriaName:      input.CafeteriaName,
				Comment:            v.Comment,
				TagRating:          tagRatings,
				Image:              getImageToBytes(v.Image),
				CafeteriaVisitedAt: timestamppb.New(v.Timestamp),
			}
		}
		return ratingResults
	} else {
		return make([]*pb.MealRating, 0)
	}
}

func getImageToBytes(path string) []byte {

	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	return bytes
}

/*
Queries the average ratings for either cafeteriaRatingTags, mealRatingTags or NameTags.
Since the db only stores IDs in the results, the tags must be joined to retrieve their names form the rating_options tables.
*/
func queryTags(db *gorm.DB, cafeteriaID int32, mealID int32, ratingType int32) []*pb.TagRatingsResult {
	var results []QueryRatingTag
	var res error
	if ratingType == MEAL {
		res = db.Table("meal_rating_tag_option options").
			Joins("JOIN meal_rating_tag_result results ON options.id = results.tagID").
			Select("options.DE as DE, results.average as average, "+
				"options.EN as EN, results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ? AND results.mealID = ?", cafeteriaID, mealID).
			Scan(&results).Error
	} else if ratingType == CAFETERIA {
		res = db.Table("cafeteria_rating_tag_option options").
			Joins("JOIN cafeteria_rating_tag_result results ON options.id = results.tagID").
			Select("options.DE as DE, results.average as average, "+
				"options.EN as EN, results.min as min, results.max as max, results.std as std").
			Where("results.cafeteriaID = ?", cafeteriaID).
			Scan(&results).Error
	} else { //Query for name tags
		res = db.Table("meal_to_meal_name_tag mapping").
			Where("mapping.mealID = ?", mealID).
			Select("mapping.nameTagID as tag").
			Joins("JOIN meal_name_tag_result results ON mapping.nameTagID = results.tagID").
			Joins("JOIN meal_name_tag_option options ON mapping.nameTagID = options.id").
			Select("options.DE as DE, results.average as average, " +
				"options.EN as EN, results.min as min, results.max as max, results.std as std").
			Scan(&results).Error
	}

	if res != nil {
		log.Println(res)
	}

	elements := make([]*pb.TagRatingsResult, len(results)) //needed since the gRPC element does not specify column names - cannot be directly queried into the grpc message object.
	for i, v := range results {
		elements[i] = &pb.TagRatingsResult{
			DE:                v.DE,
			EN:                v.EN,
			AveragePoints:     v.Average,
			StandardDeviation: v.Std,
			MinPoints:         v.Min,
			MaxPoints:         v.Max,
		}
	}

	return elements
}

/*
Query all rating tags which belong to a specific rating given with an ID and return it as TagratingOverviews
*/
func queryTagRatingsOverviewForRating(s *CampusServer, mealID int32, ratingType int32) []*pb.TagRatingResult {
	var results []QueryOverviewRatingTag
	var res error
	if ratingType == MEAL {
		res = s.db.Table("meal_rating_tag_option options").
			Joins("JOIN meal_rating_tag rating ON options.id = rating.tagID").
			Where("rating.parentRating = ?", mealID).
			Select("options.DE as DE, options.EN as EN, rating.rating as rating").
			Scan(&results).Error
	} else {
		res = s.db.Table("cafeteria_rating_tag_option options").
			Joins("JOIN cafeteria_rating_tag rating ON options.id = rating.tagID").
			Where("rating.parentRating = ?", mealID).
			Select("options.DE as DE, options.EN as EN, rating.rating as rating").
			Scan(&results).Error
	}

	if res != nil {
		log.Error(res)
	}
	elements := make([]*pb.TagRatingResult, len(results))
	for i, a := range results {
		elements[i] = &pb.TagRatingResult{
			EN:     a.EN,
			DE:     a.DE,
			Points: a.Points,
		}
	}
	return elements
}

/*
RPC Endpoint
Allows to store a new cafeteria Rating.
If one of the parameters is invalid, an error will be returned. Otherwise the rating will be saved.
All rating tags which were given with the new rating are stored if they are valid tags, if at least one tag was
invalid, an error is returned, all valid ratings tags will be stored nevertheless. Either the german or the english name can be returned to sucessfully store tags
*/
func (s *CampusServer) NewCafeteriaRating(_ context.Context, input *pb.NewCafeteriaRatingRequest) (*emptypb.Empty, error) {
	cafeteriaID, errorRes := inputSanitization(input.Points, input.Image, input.Comment, input.CafeteriaName, s)
	if errorRes != nil {
		return nil, errorRes
	}

	path := fmt.Sprintf("%s%d%s", "../images/cafeterias/", cafeteriaID, "/")
	respath, reserror := storeImage(path, input.Image)

	if reserror != nil {
		println("storing an image did not succeed")
	}

	rating := cafeteria_rating_models.CafeteriaRating{
		Comment:     input.Comment,
		Points:      input.Points,
		CafeteriaID: cafeteriaID,
		Timestamp:   time.Now(),
		Image:       respath,
	}

	s.db.Model(&cafeteria_rating_models.CafeteriaRating{}).Create(&rating)

	return storeRatingTags(s, rating.CafeteriaRating, input.Tags, CAFETERIA)
}

/*
stores an image and returns teh path to this image.
if needed, a new directory will be created and the path is extended until it is unique
*/
func storeImage(path string, i []byte) (string, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
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
	defer out.Close()
	var opts jpeg.Options
	opts.Quality = 100
	errFile = jpeg.Encode(out, img, &opts)
	return imgPath, errFile
}

/*
RPC Endpoint
Allows to store a new meal Rating.
If one of the parameters is invalid, an error will be returned. Otherwise the rating will be saved.
The ratingNumber will be saved for each corresponding Mealnametag.
All rating tags which were given with the new rating are stored if they are valid tags, if at least one tag was
invalid, an error is returned, all valid ratings tags will be stored nevertheless. Either the german or the english name can be returned to sucessfully store tags
*/
func (s *CampusServer) NewMealRating(_ context.Context, input *pb.NewMealRatingRequest) (*emptypb.Empty, error) {

	cafeteriaID, errorRes := inputSanitization(input.Points, input.Image, input.Comment, input.CafeteriaName, s)
	if errorRes != nil {
		return nil, errorRes
	}

	var meal *cafeteria_rating_models.Meal
	//cafeteriaID := getIDForCafeteriaName(input.CafeteriaName, s.db)
	mealExists := s.db.Model(&cafeteria_rating_models.Meal{}). //Meal must exist in the given mensa
									Where("name LIKE ? AND cafeteriaID = ?", input.Meal, cafeteriaID).
									First(&meal).RowsAffected

	if mealExists == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Meal is not offered in this week in this canteen. Rating has not been saved.")
	}

	path := fmt.Sprintf("%s%d%s%d%s", "../images/meals/", cafeteriaID, "/", meal.Meal, "/")
	respath, reserror := storeImage(path, input.Image)

	if reserror != nil {
		println("storing an image did not succeed")
	}

	rating := cafeteria_rating_models.MealRating{
		Comment:     input.Comment,
		CafeteriaID: cafeteriaID,
		MealID:      meal.Meal,
		Points:      input.Points,
		Timestamp:   time.Now(),
		Image:       respath,
	}

	s.db.Model(&cafeteria_rating_models.MealRating{}).Create(&rating)

	assignMealNameTag(s, rating, meal.Meal)
	return storeRatingTags(s, rating.MealRating, input.Tags, MEAL)
}

/*
Query all name tags for this specific meal and generate the MealNameTag Ratings ffor each name tag
*/
func assignMealNameTag(s *CampusServer, rating cafeteria_rating_models.MealRating, mealID int32) {
	var result []int
	err := s.db.Model(&cafeteria_rating_models.MealToMealNameTag{}).Where("mealID = ? ", mealID).Select("nameTagID").Scan(&result).Error
	if err != nil {
		log.Error(err)
	} else {
		for _, tagID := range result {
			s.db.Model(&cafeteria_rating_models.MealNameTag{}).Create(&cafeteria_rating_models.MealNameTag{
				CorrespondingRating: rating.MealRating,
				Points:              rating.Points,
				TagNameID:           tagID,
			})
		}
	}
}

/*
Checks parameters of the new rating for all cafeteria and meal ratings.
Additionally, queries the cafeteria ID, since it checks whether the cafeteria actually exists.
*/
func inputSanitization(rating int32, image []byte, comment string, cafeteriaName string, s *CampusServer) (int32, error) {
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

	var result *cafeteria_rating_models.Cafeteria
	testCanteen := s.db.Model(&cafeteria_rating_models.Cafeteria{}).
		Where("name LIKE ?", cafeteriaName).
		First(&result)
	if testCanteen.Error != nil || testCanteen.RowsAffected == 0 {
		return -1, status.Errorf(codes.InvalidArgument, "Cafeteria does not exist. Rating has not been saved.")
	}

	return result.Cafeteria, nil
}

/*
Checks whether the rating-tag name is a valid option and if so,
it will be saved with a reference to the rating
*/
func storeRatingTags(s *CampusServer, parentRatingID int32, tags []*pb.TagRating, tagType int) (*emptypb.Empty, error) {
	var errorOccured = ""
	if len(tags) > 0 {
		usedTagIds := make(map[int]int)
		insertModel := getModelStoreTag(tagType, s.db)
		for _, tag := range tags {
			var currentTag int

			exists := getModelStoreTagOption(tagType, s.db).
				Where("EN LIKE ? OR DE LIKE ?", tag.Tag, tag.Tag).
				Select("id").
				First(&currentTag)

			if exists.Error != nil || exists.RowsAffected == 0 {
				log.Println("Tag with tagname ", tag.Tag, "does not exist")
				errorOccured = errorOccured + ", " + tag.Tag
			} else {
				if usedTagIds[currentTag] == 0 {
					insertModel.
						Create(&cafeteria_rating_models.MealRatingTag{
							CorrespondingRating: parentRatingID,
							Points:              int32(tag.Points),
							TagID:               currentTag})
					usedTagIds[currentTag] = 1
				} else {
					log.Println("Each Rating tag must be used at most once in a rating.")
				}
			}
		}
	}

	if len(errorOccured) > 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "The Tag "+errorOccured+" does not exist. Remaining rating was saved without this rating tag")
	} else {
		return &emptypb.Empty{}, nil
	}

}

/*
Returns the db model of the option toable to reduce code duplicates
*/
func getModelStoreTagOption(tagType int, db *gorm.DB) *gorm.DB {
	if tagType == MEAL {
		return db.Model(&cafeteria_rating_models.MealRatingTagOption{})
	} else {
		return db.Model(&cafeteria_rating_models.CafeteriaRatingTagOption{})
	}
}

func getModelStoreTag(tagType int, db *gorm.DB) *gorm.DB {
	if tagType == MEAL {
		return db.Model(&cafeteria_rating_models.MealRatingTag{})
	} else {
		return db.Model(&cafeteria_rating_models.CafeteriaRatingTag{})
	}
}

func getIDForCafeteriaName(name string, db *gorm.DB) int32 {
	var result int32
	db.Model(&cafeteria_rating_models.Cafeteria{}).
		Where("name LIKE ?", name).
		Select("cafeteria").
		Scan(&result)
	return result
}

func getIDForMealName(name string, cafeteriaID int32, db *gorm.DB) int32 {
	var result int32 = -1
	db.Model(&cafeteria_rating_models.Meal{}).
		Where("name LIKE ? AND cafeteriaID = ?", name, cafeteriaID).
		Select("meal").
		Scan(&result)

	return result
}

/*
RPC Endpoint
Retunrs all valid Tags to quickly rate meals in english and german
*/
func (s *CampusServer) GetAvailableMealTags(_ context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	var result []*cafeteria_rating_models.MealRatingTagOption
	s.db.Model(&cafeteria_rating_models.MealRatingTagOption{}).Select("DE, EN").Scan(&result)

	elements := make([]*pb.TagRatingOverview, len(result))
	for i, a := range result {
		elements[i] = &pb.TagRatingOverview{EN: a.EN, DE: a.DE}
	}

	return &pb.GetRatingTagsReply{
		Tags: elements,
	}, nil
}

/*
RPC Endpoint
Retunrs all valid Tags to quickly rate meals in english and german
*/
func (s *CampusServer) GetAvailableCafeteriaTags(_ context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	var result []*cafeteria_rating_models.CafeteriaRatingTagOption
	s.db.Model(&cafeteria_rating_models.CafeteriaRatingTagOption{}).Select("DE,EN").Scan(&result)

	elements := make([]*pb.TagRatingOverview, len(result))
	for i, a := range result {
		elements[i] = &pb.TagRatingOverview{EN: a.EN, DE: a.DE}
	}

	return &pb.GetRatingTagsReply{
		Tags: elements,
	}, nil
}

func (s *CampusServer) GetCafeterias(_ context.Context, _ *emptypb.Empty) (*pb.GetCafeteriaResponse, error) {
	var result []*pb.Cafeteria
	s.db.Model(&cafeteria_rating_models.Cafeteria{}).Select("name,address,latitude,longitude").Scan(&result)

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
	}, nil
}
