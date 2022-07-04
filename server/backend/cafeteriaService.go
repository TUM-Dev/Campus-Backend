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

func (s *CampusServer) GetCafeteriaRatingLastThree(ctx context.Context, _ *pb.GetCafeteriaRating) (*pb.GetCafeteriaRatingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCafeteriaRatingLastThree not implemented")
}

func (s *CampusServer) GetMealRatingLastThree(ctx context.Context, input *pb.GetMealInCafeteriaRating) (*pb.GetMealInCafeteriaRatingReply, error) {
	var result cafeteria_rating_models.MealRatingsAverage //get the average rating for this specific meal
	res := s.db.Model(&cafeteria_rating_models.MealRatingsAverage{}).
		Where("cafeteria = ? AND meal = ?", input.CafeteriaName, input.Meal).
		First(&result)

	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong while accessing the database")
	}

	if res.RowsAffected > 0 {
		ratings := queryLastRatingsWithLimit(input, s)
		//todo query tagRatings name + meal
		mealTags := queryMealTags(s.db)
		//nameTags := queryNameTags(s.db)

		return &pb.GetMealInCafeteriaRatingReply{
			AverageRating: float64(result.Average),
			Rating:        ratings,
			TagRating:     mealTags,
		}, nil
	} else {
		return &pb.GetMealInCafeteriaRatingReply{
			AverageRating: -1,
		}, nil
	}

}

func queryMealTags(db *gorm.DB) []*pb.TagRating {

	var results []*pb.TagRating

	res := db.Model(&cafeteria_rating_models.MealRatingTagsAverage{}).
		Joins("join meal_rating_tags_options on meal_rating_tags_options.id = meal_rating_tags_results.tagID").
		Select("meal_rating_tags_options.nameDE, meal_rating_tags_results.average"). //+ meal_rating_tags_options.nameDE, meal_rating_tags_results.min, meal_rating_tags_results.max
		Scan(&results).Error
	/*err := c.db.Raw("SELECT  mr.cafeteriaID as cafeteriaID, mnt.tagnameID as tagID, AVG(mnt.rating) as average, MAX(mnt.rating) as max, MIN(mnt.rating) as min" +
	" FROM meal_rating mr" +
	" JOIN meal_name_tags mnt ON mr.id = mnt.parentRating" +
	" GROUP BY mr.cafeteriaID, mnt.tagnameID").Scan(&results).Error
	*/
	if res != nil {
		log.Println(res)
	}
	return results
}

func queryNameTags(db *gorm.DB) []*pb.TagRating {

	var results []*pb.TagRating

	res := db.Model(&cafeteria_rating_models.MealRatingTagsAverage{}).
		Joins("join meal_rating_tags_options on meal_rating_tags_options.id = meal_rating_tags_results.tagID").
		Select("meal_rating_tags_options.nameDE, meal_rating_tags_results.average"). //+ meal_rating_tags_options.nameDE, meal_rating_tags_results.min, meal_rating_tags_results.max
		Scan(&results).Error
	/*err := c.db.Raw("SELECT  mr.cafeteriaID as cafeteriaID, mnt.tagnameID as tagID, AVG(mnt.rating) as average, MAX(mnt.rating) as max, MIN(mnt.rating) as min" +
	" FROM meal_rating mr" +
	" JOIN meal_name_tags mnt ON mr.id = mnt.parentRating" +
	" GROUP BY mr.cafeteriaID, mnt.tagnameID").Scan(&results).Error
	*/
	if res != nil {
		log.Println(res)
	}
	return results
}

func queryLastRatingsWithLimit(input *pb.GetMealInCafeteriaRating, s *CampusServer) []*pb.MealRating {
	var ratings []cafeteria_rating_models.MealRating
	if input.Limit > 0 {
		errRatings := s.db.Model(&cafeteria_rating_models.MealRating{}).
			Where("cafeteria = ? AND meal = ?", input.CafeteriaName, input.Meal).
			First(&ratings).
			Limit(int(input.Limit)).
			Find(ratings)

		if errRatings.Error != nil {
			return make([]*pb.MealRating, 0)
		}
		ratingResults := make([]*pb.MealRating, len(ratings))

		//todo add timestamp
		//todo add meal tags which were added to this rating
		for i, v := range ratings {
			ratingResults[i] = &pb.MealRating{
				Rating:        v.Rating,
				Meal:          getNameForMealID(v.MealID, s.db),
				CafeteriaName: getNameForCafeteriaID(v.CafeteriaID, s.db),
				Comment:       v.Comment,
			}
		}
		return ratingResults
	} else {
		return make([]*pb.MealRating, 0)
	}
}

func (s *CampusServer) NewCafeteriaRating(ctx context.Context, input *pb.NewRating) (*emptypb.Empty, error) {
	validInput := inputSanitization(input, s)
	if validInput != nil {
		return nil, validInput
	}

	rating := cafeteria_rating_models.CafeteriaRating{
		Comment:     input.Comment,
		Rating:      input.Rating,
		CafeteriaID: getIDForCafeteriaName(input.CafeteriaName, s.db),
		Timestamp:   time.Now()}

	s.db.Model(&cafeteria_rating_models.CafeteriaRating{}).Create(&rating)
	storeRatingTags(s, rating.Id, input.Tags, CAFETERIA)

	return &emptypb.Empty{}, nil
}

func (s *CampusServer) NewMealRating(ctx context.Context, input *pb.NewRating) (*emptypb.Empty, error) {

	validInput := inputSanitization(input, s)
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
		MealID:      getIDForMealName(input.Meal, input.CafeteriaName, s.db),
		Rating:      input.Rating,
		Timestamp:   time.Now()}

	s.db.Model(&cafeteria_rating_models.MealRating{}).Create(&rating)

	storeRatingTags(s, rating.Id, input.Tags, MEAL)
	extractAndStoreMealNameTags(s, rating, input.Meal)

	return &emptypb.Empty{}, nil
}

func inputSanitization(input *pb.NewRating, s *CampusServer) error {
	if input.Rating > 10 || input.Rating < 0 {
		return status.Errorf(codes.InvalidArgument, "Rating must be a positive number not larger than 10. Rating has not been saved.")
	}

	if len(input.Image) > 131100 {
		return status.Errorf(codes.InvalidArgument, "Image must not be larger than 1MB. Rating has not been saved.")
	}

	if len(input.Comment) > 256 {
		return status.Errorf(codes.InvalidArgument, "Ratings can only contain up to 256 characters, this is too long. Rating has not been saved.")
	}

	if strings.Contains(input.Comment, "@") {
		return status.Errorf(codes.InvalidArgument, "Comments must not contain @ symbols in order to prevent misuse. Rating has not been saved.")
	}

	var result *cafeteria_rating_models.Cafeteria
	testCanteen := s.db.Model(&cafeteria_rating_models.Cafeteria{}).
		Where("name LIKE ?", input.CafeteriaName).
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
func storeRatingTags(s *CampusServer, parentRatingID int32, tags []string, tagType int) {
	if len(tags) > 0 {
		usedTagIds := make(map[int]int)
		insertModel := getModelStoreTag(tagType, s.db)
		for _, tag := range tags {
			var currentTag int

			exists := getModelStoreTagOption(tagType, s.db).
				Where("nameEN LIKE ? OR nameDE LIKE ?", tag, tag).
				Select("id").
				First(&currentTag)

			if exists.RowsAffected > 0 && usedTagIds[currentTag] == 0 {
				insertModel.
					Create(&cafeteria_rating_models.MealRatingsTags{
						ParentRating: parentRatingID,
						Rating:       int32(5),
						TagID:        currentTag})
				usedTagIds[currentTag] = 1
			} else {
				log.Println("Invalid Tag Name, Tag", tag, "was not saved. Each Rating tag must be used at most once in a rating.")
			}
		}
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

func getIDForMealName(name string, cafeteriaName string, db *gorm.DB) int32 {
	var result int32 = -1
	db.Model(&cafeteria_rating_models.Meal{}).
		Where("name LIKE ?", name).
		Where("cafeteriaID LIKE ?", cafeteriaName).
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

	log.Println(len(includedTags))

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
	var result []string
	s.db.Model(&cafeteria_rating_models.MealRatingsTagsOptions{}).Select("nameDE").Scan(&result)

	return &pb.GetRatingTagsReply{
		Tags: result,
	}, nil
}

func (s *CampusServer) GetAvailableCafeteriaTags(ctx context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	var result []string
	s.db.Model(&cafeteria_rating_models.CafeteriaRatingsTagsOptions{}).Select("nameDE").Scan(&result)

	return &pb.GetRatingTagsReply{
		Tags: result,
	}, nil
}

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
