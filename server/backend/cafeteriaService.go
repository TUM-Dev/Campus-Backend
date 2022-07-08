package backend

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MEAL = 1
const CAFETERIA = 2

type MultiLanguageTags struct {
	MultiLanguageTags []Tag `json:"tags"`
}
type Tag struct {
	TagNameEnglish string `json:"tagNameEnglish"`
	TagNameGerman  string `json:"tagNameGerman"`
}

type MultiLanguageNameTags struct {
	MultiLanguageNameTags []NameTag `json:"tags"`
}
type NameTag struct {
	TagNameEnglish string   `json:"tagNameEnglish"`
	TagNameGerman  string   `json:"tagNameGerman"`
	Notincluded    []string `json:"notincluded"`
	Canbeincluded  []string `json:"canbeincluded"`
}

type QueryRatingTag struct {
	NameEN  string  `gorm:"column:nameEN;type:mediumtext;" json:"nameEN"`
	NameDE  string  `gorm:"column:nameDE;type:mediumtext;" json:"nameDE"`
	Average float64 `json:"average"`
	Min     int32   `json:"min"`
	Max     int32   `json:"max"`
}

/*
Writes all available tags from the json file into tables in order to make them easier to use
*/
func initTagRatingOptions(db *gorm.DB) {
	updateTagTable("backend/static_data/mealRatingTags.json", db, MEAL)
	updateTagTable("backend/static_data/cafeteriaRatingTags.json", db, CAFETERIA)
	updateNameTagOptions(db)
}

/*
Updates the list of mealtags.

*/
func updateNameTagOptions(db *gorm.DB) {
	absPathMealNames, _ := filepath.Abs("backend/static_data/mealNameTags.json")
	tagsNames := generateNameTagListFromFile(absPathMealNames)
	var elementID int32
	for _, v := range tagsNames.MultiLanguageNameTags {
		var parentID int32

		potentialTag := db.Model(&cafeteria_rating_models.MealNameTagOptions{}).
			Where("nameEN LIKE ? AND nameDE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
			Select("id").
			Scan(&parentID)

		if potentialTag.RowsAffected == 0 {
			parent := cafeteria_rating_models.MealRatingsTagsOptions{
				NameDE: v.TagNameGerman,
				NameEN: v.TagNameEnglish}

			db.Model(&cafeteria_rating_models.MealNameTagOptions{}).
				Create(&parent)
			parentID = parent.Id
		}

		for _, u := range v.Canbeincluded {
			resultIncluded := db.Model(&cafeteria_rating_models.MealNameTagOptionsIncluded{}).
				Where("expression LIKE ? AND NameTagID = ?", u, parentID).
				Select("id").
				Scan(&elementID)
			if resultIncluded.RowsAffected == 0 {
				db.Model(&cafeteria_rating_models.MealNameTagOptionsIncluded{}).
					Create(&cafeteria_rating_models.MealNameTagOptionsIncluded{
						Expression: u,
						NameTagID:  parentID})
			}
		}
		for _, u := range v.Notincluded {
			resultIncluded := db.Model(&cafeteria_rating_models.MealNameTagOptionsExcluded{}).
				Where("expression LIKE ? AND NameTagID = ?", u, parentID).
				Select("id").
				Scan(&elementID)
			if resultIncluded.RowsAffected == 0 {
				db.Model(&cafeteria_rating_models.MealNameTagOptionsExcluded{}).
					Create(&cafeteria_rating_models.MealNameTagOptionsExcluded{
						Expression: u,
						NameTagID:  parentID})
			}
		}
	}
}

/*
Reads the json file at the given path and checks whether the values have already been inserted into the corresponding table.
If an entry with the same German and English name exists, the entry won't be added.
The TagType is used to identify the corresponding model
*/
func updateTagTable(path string, db *gorm.DB, tagType int) {
	absPathMeal, _ := filepath.Abs(path)
	tagsMeal := generateRatingTagListFromFile(absPathMeal)
	insertModel := getTagModel(tagType, db)
	for _, v := range tagsMeal.MultiLanguageTags {
		var result int32

		potentialTag := getTagModel(tagType, db).
			Where("nameEN LIKE ? AND nameDE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
			Select("id").
			Scan(&result)

		if potentialTag.RowsAffected == 0 {
			println("New entry inserted to Rating Tag Options")
			element := cafeteria_rating_models.MealRatingsTagsOptions{
				NameDE: v.TagNameGerman,
				NameEN: v.TagNameEnglish}
			insertModel.
				Create(&element)
		}
	}
}

func getTagModel(tagType int, db *gorm.DB) *gorm.DB {
	if tagType == MEAL {
		return db.Model(&cafeteria_rating_models.MealRatingsTagsOptions{})
	} else {
		return db.Model(&cafeteria_rating_models.CafeteriaRatingsTagsOptions{})
	}
}

func (s *CampusServer) GetCafeteriaRatings(ctx context.Context, input *pb.CafeteriaRatingRequest) (*pb.CafeteriaRatingResponse, error) {
	var result cafeteria_rating_models.CafeteriaRatingsAverage //get the average rating for this specific cafeteria
	cafeteriaID := getIDForCafeteriaName(input.CafeteriaName, s.db)
	res := s.db.Model(&cafeteria_rating_models.CafeteriaRatingsAverage{}).
		Where("cafeteriaID = ?", cafeteriaID).
		First(&result)

	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong while accessing the database")
	}

	if res.RowsAffected > 0 {
		ratings := queryLastCafeteriaRatingsWithLimit(input, cafeteriaID, s)
		cafeteriaTags := queryCafeteriaTags(s.db, cafeteriaID)

		return &pb.CafeteriaRatingResponse{
			AverageRating: float64(result.Average),
			MinRating:     int32(result.Min),
			MaxRating:     int32(result.Max),
			Rating:        ratings,
			RatingTags:    cafeteriaTags,
			//NameTags: nameTags,
		}, nil
	} else {
		return &pb.CafeteriaRatingResponse{
			AverageRating: -1,
			MinRating:     -1,
			MaxRating:     -1,
		}, nil
	}
}

func queryLastCafeteriaRatingsWithLimit(input *pb.CafeteriaRatingRequest, cafeteriaID int32, s *CampusServer) []*pb.CafeteriaRating {
	var ratings []cafeteria_rating_models.CafeteriaRating
	if input.Limit > 0 {
		errRatings := s.db.Model(&cafeteria_rating_models.CafeteriaRating{}).
			Where("cafeteriaID = ?", cafeteriaID).
			Limit(int(input.Limit)).
			Find(&ratings)

		if errRatings.Error != nil {
			return make([]*pb.CafeteriaRating, 0)
		}
		ratingResults := make([]*pb.CafeteriaRating, len(ratings))

		//todo add timestamp
		//todo add meal tags which were added to this rating
		for i, v := range ratings {
			ratingResults[i] = &pb.CafeteriaRating{
				Rating:        v.Rating,
				CafeteriaName: input.CafeteriaName,
				Comment:       v.Comment,
				//Image: v.Image,
				//CafeteriaVisitedAt: v.Timestamp,
				//TagRating
			}
		}
		return ratingResults
	} else {
		return make([]*pb.CafeteriaRating, 0)
	}
}

func (s *CampusServer) GetMealRatings(ctx context.Context, input *pb.MealRatingsRequest) (*pb.MealRatingsResponse, error) {
	var result cafeteria_rating_models.MealRatingsAverage //get the average rating for this specific meal
	cafeteriaID := getIDForCafeteriaName(input.CafeteriaName, s.db)
	res := s.db.Model(&cafeteria_rating_models.MealRatingsAverage{}).
		Where("cafeteriaID = ? AND meal = ?", cafeteriaID, input.Meal).
		First(&result)

	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong while accessing the database")
	}

	if res.RowsAffected > 0 {
		ratings := queryLastMealRatingsWithLimit(input, cafeteriaID, s)
		mealTags := queryMealTags(s.db, input.CafeteriaName, input.Meal)
		//nameTags := queryNameTags(s.db)

		return &pb.MealRatingsResponse{
			AverageRating: float64(result.Average),
			MinRating:     int32(result.Min),
			MaxRating:     int32(result.Max),
			Rating:        ratings,
			RatingTags:    mealTags,
			//NameTags: nameTags,
		}, nil
	} else {
		return &pb.MealRatingsResponse{
			AverageRating: -1,
			MinRating:     -1,
			MaxRating:     -1,
		}, nil
	}

}

func queryMealTags(db *gorm.DB, cafeteriaName string, mealName string) []*pb.TagRatingsResult {
	//todo where nach cafeteria ID und meal ID
	cafeteriaID := getIDForCafeteriaName(cafeteriaName, db)
	var results []*pb.TagRatingsResult
	res := db.Model(&cafeteria_rating_models.MealRatingTagsAverage{}).
		Joins("join meal_rating_tags_options on meal_rating_tags_options.id = meal_rating_tags_results.tagID").
		Select("meal_rating_tags_options.nameDE, meal_rating_tags_results.average"+
			"meal_rating_tags_options.nameEN, meal_rating_tags_results.min, meal_rating_tags_results.max"). //+ meal_rating_tags_options.nameDE, meal_rating_tags_results.min, meal_rating_tags_results.max
		Where("meal_rating_tags_results.mealID =? AND meal_rating_tags_results.cafeteriaID = ?", getIDForMealName(mealName, cafeteriaID, db), cafeteriaID).
		Scan(&results).Error

	if res != nil {
		log.Println(res)
	}
	return results
}

func queryCafeteriaTags(db *gorm.DB, cafeteriaID int32) []*pb.TagRatingsResult {

	var results []QueryRatingTag

	/*	res := db.Table("cafeteria_rating_tags_results").
		Joins("join cafeteria_rating_tags_options on cafeteria_rating_tags_options.id = cafeteria_rating_tags_results.tagID").
		Select("cafeteria_rating_tags_options.nameDE as nameDE, cafeteria_rating_tags_results.average as average, "+
			"cafeteria_rating_tags_options.nameEN as nameEN, cafeteria_rating_tags_results.min as min, cafeteria_rating_tags_results.max as max"). //+ meal_rating_tags_options.nameDE, meal_rating_tags_results.min, meal_rating_tags_results.max
		Where("cafeteria_rating_tags_results.cafeteriaID = ?", cafeteriaID).
		Scan(&results).Error
	*/
	res := db.Table("cafeteria_rating_tags_options").
		Joins("JOIN cafeteria_rating_tags_results ON cafeteria_rating_tags_options.id = cafeteria_rating_tags_results.tagID").
		Select("cafeteria_rating_tags_options.nameDE as nameDE, cafeteria_rating_tags_results.average as average, "+
			"cafeteria_rating_tags_options.nameEN as nameEN, cafeteria_rating_tags_results.min as min, cafeteria_rating_tags_results.max as max"). //+ meal_rating_tags_options.nameDE, meal_rating_tags_results.min, meal_rating_tags_results.max
		Where("cafeteria_rating_tags_results.cafeteriaID = ?", cafeteriaID).
		Scan(&results).Error
	/*	res := db.Raw("SELECT  options.nameDE as nameDE, results.average as average, options.nameEN as nameEN, results.min as min, results.max as max, options.id as testID"+
		" FROM cafeteria_rating_tags_results results, "+
		" JOIN cafeteria_rating_tags_options options ON  results.tagID = options.id "+
		" WHERE results.cafeteriaID = ? ", cafeteriaID).Scan(&results).Error
	*/
	//fixme: erstmal als raw anfrage wie beim preprocessing durchführen
	elements := make([]*pb.TagRatingsResult, len(results))
	for i, v := range results {
		elements[i] = &pb.TagRatingsResult{
			NameDE:        v.NameDE,
			NameEN:        v.NameEN,
			AverageRating: v.Average,
			MinRating:     v.Min,
			MaxRating:     v.Max,
		}
	}
	if res != nil {
		log.Println(res)
	}
	return elements
}

func queryNameTags(db *gorm.DB) []*pb.TagRating {

	//todo darf nicht auf der options tabelle ansetzten, sondern auf einer unteranfrage, die nur nochalle übrig lässt, die auch im namen enthalten sind
	var results []*pb.TagRating

	res := db.Model(&cafeteria_rating_models.MealNameTagsAverage{}).
		Joins("join meal_name_tags_options on meal_name_tags_options.id = meal_name_tags_results.tagID").
		Select("meal_name_tags_options.nameDE, meal_name_tags_results.average"). //+ meal_rating_tags_options.nameDE, meal_rating_tags_results.min, meal_rating_tags_results.max
		Scan(&results).Error
	//todo only nametags for the current meal -> wieder excluded und included mappen
	if res != nil {
		log.Println(res)
	}
	return results
}

func queryLastMealRatingsWithLimit(input *pb.MealRatingsRequest, cafeteriaID int32, s *CampusServer) []*pb.MealRating {
	var ratings []cafeteria_rating_models.MealRating
	mealID := getIDForMealName(input.Meal, cafeteriaID, s.db)
	if input.Limit > 0 {
		errRatings := s.db.Model(&cafeteria_rating_models.MealRating{}).
			Where("cafeteriaID = ? AND mealID = ?", cafeteriaID, mealID).
			Limit(int(input.Limit)).
			Find(&ratings)

		if errRatings.Error != nil {
			return make([]*pb.MealRating, 0)
		}
		ratingResults := make([]*pb.MealRating, len(ratings))

		//todo add timestamp
		//todo add meal tags which were added to this rating
		for i, v := range ratings {
			ratingResults[i] = &pb.MealRating{
				Rating:        v.Rating,
				Meal:          input.Meal,
				CafeteriaName: input.CafeteriaName,
				Comment:       v.Comment,
			}
		}
		return ratingResults
	} else {
		return make([]*pb.MealRating, 0)
	}
}

func (s *CampusServer) NewCafeteriaRating(ctx context.Context, input *pb.NewCafeteriaRatingRequest) (*emptypb.Empty, error) {
	validInput := inputSanitization(input.Rating, input.Image, input.Comment, input.CafeteriaName, s)
	if validInput != nil {
		return nil, validInput
	}

	rating := cafeteria_rating_models.CafeteriaRating{
		Comment:     input.Comment,
		Rating:      input.Rating,
		CafeteriaID: getIDForCafeteriaName(input.CafeteriaName, s.db),
		Timestamp:   time.Now()}

	s.db.Model(&cafeteria_rating_models.CafeteriaRating{}).Create(&rating)

	return storeRatingTags(s, rating.Id, input.Tags, CAFETERIA)
}

func (s *CampusServer) NewMealRating(ctx context.Context, input *pb.NewMealRatingRequest) (*emptypb.Empty, error) {

	validInput := inputSanitization(input.Rating, input.Image, input.Comment, input.CafeteriaName, s)
	if validInput != nil {
		return nil, validInput
	}

	var meal *cafeteria_rating_models.Meal
	cafeteriaID := getIDForCafeteriaName(input.CafeteriaName, s.db)
	testDish := s.db.Model(&cafeteria_rating_models.Meal{}).
		Where("name LIKE ? AND cafeteriaID = ?", input.Meal, cafeteriaID).
		First(&meal)
	if testDish.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Meal is not offered in this week in this canteen. Rating has not been saved.")
	}

	rating := cafeteria_rating_models.MealRating{
		Comment:     input.Comment,
		CafeteriaID: cafeteriaID,
		MealID:      getIDForMealName(input.Meal, cafeteriaID, s.db),
		Rating:      input.Rating,
		Timestamp:   time.Now()}

	s.db.Model(&cafeteria_rating_models.MealRating{}).Create(&rating)

	extractAndStoreMealNameTags(s, rating, input.Meal)

	return storeRatingTags(s, rating.Id, input.Tags, MEAL)
}

func inputSanitization(rating int32, image []byte, comment string, cafeteriaName string, s *CampusServer) error {
	if rating > 10 || rating < 0 {
		return status.Errorf(codes.InvalidArgument, "Rating must be a positive number not larger than 10. Rating has not been saved.")
	}

	if len(image) > 131100 {
		return status.Errorf(codes.InvalidArgument, "Image must not be larger than 1MB. Rating has not been saved.")
	}

	if len(comment) > 256 {
		return status.Errorf(codes.InvalidArgument, "Ratings can only contain up to 256 characters, this is too long. Rating has not been saved.")
	}

	if strings.Contains(comment, "@") {
		return status.Errorf(codes.InvalidArgument, "Comments must not contain @ symbols in order to prevent misuse. Rating has not been saved.")
	}

	var result *cafeteria_rating_models.Cafeteria
	testCanteen := s.db.Model(&cafeteria_rating_models.Cafeteria{}).
		Where("name LIKE ?", cafeteriaName).
		First(&result)
	if testCanteen.RowsAffected == 0 {
		return status.Errorf(codes.InvalidArgument, "Cafeteria does not exist. Rating has not been saved.")
	}

	return nil
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
		for _, tag := range tags { //todo adapt to new version of tags
			var currentTag int

			exists := getModelStoreTagOption(tagType, s.db).
				Where("nameEN LIKE ? OR nameDE LIKE ?", tag.Tag, tag.Tag).
				Select("id").
				First(&currentTag)

			if exists.Error != nil || exists.RowsAffected == 0 {
				log.Println("Tag with tagname ", tag.Tag, "does not exist")
				errorOccured = errorOccured + ", " + tag.Tag
			} else {
				if usedTagIds[currentTag] == 0 {
					insertModel.
						Create(&cafeteria_rating_models.MealRatingsTags{
							ParentRating: parentRatingID,
							Rating:       int32(tag.Rating),
							TagID:        currentTag})
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

func getModelStoreTagOption(tagType int, db *gorm.DB) *gorm.DB {
	if tagType == MEAL {
		return db.Model(&cafeteria_rating_models.MealRatingsTagsOptions{})
	} else {
		return db.Model(&cafeteria_rating_models.CafeteriaRatingsTagsOptions{})
	}
}

func getModelStoreTag(tagType int, db *gorm.DB) *gorm.DB {
	if tagType == MEAL {
		return db.Model(&cafeteria_rating_models.MealRatingsTags{})
	} else {
		return db.Model(&cafeteria_rating_models.CafeteriaRatingTags{})
	}
}

func getNameForCafeteriaID(id int32, db *gorm.DB) string {
	var result string
	db.Model(&cafeteria_rating_models.Cafeteria{}).
		Where("id = ?", id).
		Select("name").
		First(&result)
	return result
}

func getNameForMealID(id int32, db *gorm.DB) string {
	var result string
	db.Model(&cafeteria_rating_models.Meal{}).
		Where("id = ?", id).
		Select("name").
		First(&result) //Scan(&result)
	return result
}

func getIDForCafeteriaName(name string, db *gorm.DB) int32 {
	var result int32
	db.Model(&cafeteria_rating_models.Cafeteria{}).
		Where("name LIKE ?", name).
		Select("id").
		Scan(&result)
	return result
}

func getIDForMealName(name string, cafeteriaID int32, db *gorm.DB) int32 {
	var result int32 = -1
	db.Model(&cafeteria_rating_models.Meal{}).
		Where("name LIKE ? AND cafeteriaID = ?", name, cafeteriaID). //todo darf nicht auf den namen vergleichen, sondern nur auf der id
		Select("id").
		Scan(&result)

	return result
}

/*
Checks whether the meal name includes one of the expressions for the excluded tas as well as the included tags.
The corresponding tags for all identified MealNames will be saved in the table MealNameTags.
*/
func extractAndStoreMealNameTags(s *CampusServer, rating cafeteria_rating_models.MealRating, meal string) {
	lowercaseMeal := strings.ToLower(meal)
	var includedTags []int
	s.db.Model(&cafeteria_rating_models.MealNameTagOptionsIncluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseMeal).
		Select("nameTagID").
		Scan(&includedTags)

	var excludedTags []int
	s.db.Model(&cafeteria_rating_models.MealNameTagOptionsExcluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseMeal).
		Select("nameTagID").
		Scan(&excludedTags)

	log.Println("Number of included tags: ", len(includedTags))

	//set all entries in included to -1 if the excluded tag was recognised ffor this tag rating.
	if len(excludedTags) > 0 {
		for _, a := range excludedTags {
			i := contains(includedTags, a)
			if i != -1 {
				includedTags[i] = -1
			}
		}
	}

	for _, a := range includedTags {
		if a != -1 {
			s.db.Model(&cafeteria_rating_models.MealNameTags{}).
				Create(&cafeteria_rating_models.MealNameTags{
					ParentRating: rating.Id,
					Rating:       rating.Rating,
					TagNameID:    a,
				})
		}
	}
}

func contains(s []int, e int) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func (s *CampusServer) GetAvailableMealTags(ctx context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	var result []*cafeteria_rating_models.MealRatingsTagsOptions
	s.db.Model(&cafeteria_rating_models.MealRatingsTagsOptions{}).Select("nameDE, nameEN").Scan(&result)

	elements := make([]*pb.TagRatingOverview, len(result))
	for i, a := range result {
		elements[i] = &pb.TagRatingOverview{NameEN: a.NameEN, NameDE: a.NameDE}
	}

	return &pb.GetRatingTagsReply{
		Tags: elements,
	}, nil
}

func (s *CampusServer) GetAvailableCafeteriaTags(ctx context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	var result []*cafeteria_rating_models.CafeteriaRatingsTagsOptions
	s.db.Model(&cafeteria_rating_models.CafeteriaRatingsTagsOptions{}).Select("nameDE,nameEN").Scan(&result)

	elements := make([]*pb.TagRatingOverview, len(result))
	for i, a := range result {
		elements[i] = &pb.TagRatingOverview{NameEN: a.NameEN, NameDE: a.NameDE}
	}

	return &pb.GetRatingTagsReply{
		Tags: elements,
	}, nil
}

//fixme add repeated to the proto File
func (s *CampusServer) GetCafeterias(ctx context.Context, _ *emptypb.Empty) (*pb.GetCafeteriaResponse, error) {
	var result []*pb.GetCafeteriaResponse
	s.db.Model(&cafeteria_rating_models.Cafeteria{}).Select("name,address,latitude,longitude").Scan(&result)

	return &pb.GetCafeteriaResponse{
		Name:      result[0].Name,
		Adress:    result[0].Adress,
		Latitude:  result[0].Latitude,
		Longitude: result[0].Longitude,
	}, nil
}

func generateRatingTagListFromFile(path string) MultiLanguageTags {
	byteValue := readFromFile(path)

	var tags MultiLanguageTags
	errorUnmarshal := json.Unmarshal(byteValue, &tags)
	if errorUnmarshal != nil {
		log.Error("Error in parsing json:", errorUnmarshal)
	}
	return tags
}

func generateNameTagListFromFile(path string) MultiLanguageNameTags {
	byteValue := readFromFile(path)

	var tags MultiLanguageNameTags
	errorUnmarshal := json.Unmarshal(byteValue, &tags)
	if errorUnmarshal != nil {
		log.Error("Error in parsing json:", errorUnmarshal)
	}
	return tags
}

func readFromFile(path string) []byte {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Error("Error in parsing json:", err)
		}
	}(jsonFile)

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	return byteValue
}
