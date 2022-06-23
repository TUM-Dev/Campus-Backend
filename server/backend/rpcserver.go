package backend

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func (s *CampusServer) GRPCServe(l net.Listener) error {
	grpcServer := grpc.NewServer()
	pb.RegisterCampusServer(grpcServer, s)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return grpcServer.Serve(l)
}

type CampusServer struct {
	pb.UnimplementedCampusServer
	db        *gorm.DB
	deviceBuf *deviceBuffer // deviceBuf stores all devices from recent request and flushes them to db
}

// Verify that CampusServer implements the pb.CampusServer interface
var _ pb.CampusServer = (*CampusServer)(nil)

func New(db *gorm.DB) *CampusServer {
	log.Println("Server starting up")
	initTagRatingOptions(db)

	return &CampusServer{
		db: db,
		deviceBuf: &deviceBuffer{
			lock:     sync.Mutex{},
			devices:  make(map[string]*model.Devices),
			interval: time.Minute,
		},
	}
}

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

		potentialTag := db.Model(cafeteria_rating_models.MealNameTagOptions{}).
			Where("nameEN LIKE ?", v.TagNameEnglish).
			Where("nameDE LIKE ?", v.TagNameGerman).
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
				Where("expression LIKE ?", v.TagNameEnglish).
				Where("NameTagID = ?", parentID).
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
				Where("expression LIKE ?", v.TagNameEnglish).
				Where("NameTagID = ?", parentID).
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
			Where("nameEN LIKE ?", v.TagNameEnglish).
			Where("nameDE LIKE ?", v.TagNameGerman).
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
		return db.Model(cafeteria_rating_models.MealRatingsTagsOptions{})
	} else {
		return db.Model(cafeteria_rating_models.CafeteriaRatingsTagsOptions{})
	}
}

func (s *CampusServer) GetNewsSources(ctx context.Context, _ *emptypb.Empty) (newsSources *pb.NewsSourceArray, err error) {
	if err = s.checkDevice(ctx); err != nil {
		return
	}

	var sources []model.NewsSource
	if err := s.db.Find(&sources).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var resp []*pb.NewsSource
	for _, source := range sources {
		var icon model.Files
		if err := s.db.Where("file = ?", source.Icon).First(&icon).Error; err != nil {
			icon = model.Files{File: 0}
		}
		log.Info("sending news source", source.Title)
		resp = append(resp, &pb.NewsSource{
			Source: fmt.Sprintf("%d", source.Source),
			Title:  source.Title,
			Icon:   icon.URL.String,
		})
	}
	return &pb.NewsSourceArray{Sources: resp}, nil
}

// SearchRooms returns all rooms that match the given search query.
func (s *CampusServer) SearchRooms(ctx context.Context, req *pb.SearchRoomsRequest) (*pb.SearchRoomsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}
	if req.Query == "" {
		return &pb.SearchRoomsReply{Rooms: make([]*pb.Room, 0)}, nil
	}
	var res []struct { // struct to scan database query into
		model.RoomfinderRooms
		Campus string
		Name   string
	}
	err := s.db.Raw("SELECT r.*, a.campus, a.name "+
		"FROM roomfinder_rooms r "+
		"LEFT JOIN roomfinder_building2area a ON a.building_nr = r.building_nr "+
		"WHERE MATCH(room_code, info, address) AGAINST(?)", req.Query).Scan(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &pb.SearchRoomsReply{Rooms: make([]*pb.Room, 0)}, nil
	}
	if err != nil {
		log.WithError(err).Error("failed to search rooms")
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := &pb.SearchRoomsReply{
		Rooms: make([]*pb.Room, len(res)),
	}
	for i, row := range res {
		response.Rooms[i] = &pb.Room{
			RoomId:     row.RoomID,
			RoomCode:   row.RoomCode.String,
			BuildingNr: row.BuildingNr.String,
			ArchId:     row.ArchID.String,
			Info:       row.Info.String,
			Address:    row.Address.String,
			Purpose:    row.Purpose.String,
			Campus:     row.Campus,
			Name:       row.Name,
		}
	}
	return response, nil
}

func (s *CampusServer) GetTopNews(ctx context.Context, _ *emptypb.Empty) (*pb.GetTopNewsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}
	log.Printf("Received: get top news")
	var res *model.NewsAlert
	err := s.db.Joins("Company").Where("NOW() between `from` and `to`").Limit(1).First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("Failed to fetch top news: %w", err)
	} else if res != nil {
		return &pb.GetTopNewsReply{
			//ImageUrl: res.Name,
			Link: res.Link.String,
			To:   timestamppb.New(res.To),
		}, nil
	}
	return &pb.GetTopNewsReply{}, nil
}

func (s *CampusServer) GetCafeteriaRatingLastThree(ctx context.Context, _ *pb.GetCafeteriaRating) (*pb.GetCafeteriaRatingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCafeteriaRatingLastThree not implemented")
}

func (s *CampusServer) GetMealRatingLastThree(ctx context.Context, input *pb.GetMealInCafeteriaRating) (*pb.GetMealInCafeteriaRatingReply, error) {
	var result cafeteria_rating_models.MealRatingsAverage
	err := s.db.Model(&cafeteria_rating_models.MealRatingsAverage{}).
		Where("cafeteria = ?", input.CafeteriaName).
		Where("meal = ?", input.Meal).First(&result)

	if err.Error != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong while accessing the database")
	}

	//todo add nametag ratings to the reply
	if err.RowsAffected > 0 {
		ratings := queryLastRatings(input, s)

		return &pb.GetMealInCafeteriaRatingReply{
			AverageRating: float64(result.Average),
			Rating:        ratings,
		}, nil
	} else {
		return &pb.GetMealInCafeteriaRatingReply{
			AverageRating: -1,
		}, nil
	}

}

func queryLastRatings(input *pb.GetMealInCafeteriaRating, s *CampusServer) []*pb.MealRating {
	var ratings []cafeteria_rating_models.MealRating
	if input.Limit > 0 {
		errRatings := s.db.Model(&cafeteria_rating_models.MealRating{}).
			Where("cafeteria = ?", input.CafeteriaName).
			Where("meal = ?", input.Meal).First(&ratings).
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
				Meal:          v.Meal,
				CafeteriaName: v.Cafeteria,
				Comment:       v.Comment,
			}
		}
		return ratingResults
	} else {
		return make([]*pb.MealRating, 0)
	}
}

func (s *CampusServer) NewCafeteriaRating(ctx context.Context, input *pb.NewRating) (*emptypb.Empty, error) {
	//Add cafeteriaRating
	if input.Rating > 10 || input.Rating < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Rating must be a positive number not larger than 10. Rating has not been saved.")
	}

	if len(input.Image) > 131100 {
		return nil, status.Errorf(codes.InvalidArgument, "Image must not be larger than 1MB. Rating has not been saved.")
	}

	var result *cafeteria_rating_models.Cafeteria
	testCanteen := s.db.Model(cafeteria_rating_models.Cafeteria{Name: input.CafeteriaName}).First(&result)

	if testCanteen.RowsAffected != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "Mensa does not exist. Rating has not been saved.")
	}

	rating := cafeteria_rating_models.CafeteriaRating{
		Comment:   input.Comment,
		Rating:    input.Rating,
		Cafeteria: input.CafeteriaName,
		Timestamp: time.Now()}

	s.db.Model(cafeteria_rating_models.CafeteriaRating{}).Create(&rating)
	storeRatingTags(s, rating.Id, input.Tags, CAFETERIA)

	return &emptypb.Empty{}, nil
}

func (s *CampusServer) NewMealRating(ctx context.Context, input *pb.NewRating) (*emptypb.Empty, error) {

	if input.Rating > 10 || input.Rating < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Rating must be a positive number not larger than 10. Rating has not been saved.")
	}

	if len(input.Image) > 131100 {
		return nil, status.Errorf(codes.InvalidArgument, "Image must not be larger than 1MB. Rating has not been saved.")
	}

	var result *cafeteria_rating_models.Cafeteria
	testCanteen := s.db.Model(cafeteria_rating_models.Cafeteria{Name: input.CafeteriaName}).First(&result)

	if testCanteen.RowsAffected != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "Mensa does not exist. Rating has not been saved.")
	}

	var dish *cafeteria_rating_models.Meal
	testDish := s.db.Model(cafeteria_rating_models.Meal{Name: input.Meal, Cafeteria: input.CafeteriaName}).First(&dish)

	if testDish.RowsAffected != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "Meal is not offered in this week in this canteen. Rating has not been saved.")
	}

	rating := cafeteria_rating_models.MealRating{
		Comment:   input.Comment,
		Meal:      input.Meal,
		Rating:    input.Rating,
		Timestamp: time.Now()}

	s.db.Model(cafeteria_rating_models.MealRating{}).Create(&rating)

	storeRatingTags(s, rating.Id, input.Tags, MEAL)
	extractAndStoreMealNameTags(s, rating)

	return &emptypb.Empty{}, nil
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
				Where("nameEN = @name OR nameDE = @name", sql.Named("name", tag)).
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
		return db.Model(cafeteria_rating_models.MealRatingsTagsOptions{})
	} else {
		return db.Model(cafeteria_rating_models.CafeteriaRatingsTagsOptions{})
	}
}

func getModelStoreTag(tagType int, db *gorm.DB) *gorm.DB {
	if tagType == MEAL {
		return db.Model(cafeteria_rating_models.MealRatingsTags{})
	} else {
		return db.Model(cafeteria_rating_models.CafeteriaRatingTags{})
	}
}

/*
Checks whether the meal name includes one of the expressions for the excluded tas as well as the included tags.
The corresponding tags for all identified MealNames will be saved in the table MealNameTags.
*/
func extractAndStoreMealNameTags(s *CampusServer, rating cafeteria_rating_models.MealRating) {
	lowercaseMeal := strings.ToLower(rating.Meal)
	var includedTags []int
	s.db.Model(cafeteria_rating_models.MealNameTagOptionsIncluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseMeal).
		Select("nameTagID").
		Scan(&includedTags)

	var excludedTags []int
	s.db.Model(cafeteria_rating_models.MealNameTagOptionsExcluded{}).
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
	s.db.Model(cafeteria_rating_models.MealRatingsTagsOptions{}).Select("nameDE").Scan(&result)

	return &pb.GetRatingTagsReply{
		Tags: result,
	}, nil
}

func (s *CampusServer) GetAvailableCafeteriaTags(ctx context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	var result []string
	s.db.Model(cafeteria_rating_models.CafeteriaRatingsTagsOptions{}).Select("nameDE").Scan(&result)

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
