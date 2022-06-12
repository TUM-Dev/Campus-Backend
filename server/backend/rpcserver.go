package backend

import (
	"context"
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

/*
Writes all available tags from the json file into tables in order to make them easier to use
*/
func initTagRatingOptions(db *gorm.DB) {
	absPathMeal, _ := filepath.Abs("backend/static_data/mealRatingTags.json")
	absPathCafeteria, _ := filepath.Abs("backend/static_data/cafeteriaRatingTags.json")
	tagsMeal := generateTagListFromFile(absPathMeal)
	tagsCafeteria := generateTagListFromFile(absPathCafeteria)

	//delete all existing values to prevent any inconsistencies with ids
	db.Where("1=1").Delete(&cafeteria_rating_models.CafeteriaRatingsTagsOptions{}) //Remove all meals of the previous week
	db.Where("1=1").Delete(&cafeteria_rating_models.MealRatingsTagsOptions{})      //Remove all meals of the previous week

	for _, v := range tagsMeal.MultiLanguageTags {
		db.Model(&cafeteria_rating_models.MealRatingsTagsOptions{}).
			Create(&cafeteria_rating_models.MealRatingsTagsOptions{
				NameDE: v.TagNameGerman,
				NameEN: v.TagNameEnglish})
	}

	for _, v := range tagsCafeteria.MultiLanguageTags {
		db.Model(&cafeteria_rating_models.CafeteriaRatingsTagsOptions{}).
			Create(&cafeteria_rating_models.CafeteriaRatingsTagsOptions{
				NameDE: v.TagNameGerman,
				NameEN: v.TagNameEnglish})
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
	var result cafeteria_rating_models.MealRatingResult
	err := s.db.Model(&cafeteria_rating_models.MealRatingResult{}).
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

	var parentid = rating.Id
	//Add Tag Ratings for the first cafeteria

	for i := 0; i < len(input.Tags); i++ {
		//todo tag must be included in the tag lists
		//todo add rating once the proto file is fixed
		rating := cafeteria_rating_models.CafeteriaTagRating{ParentRating: int32(parentid), Rating: int32(5), TagID: i}
		s.db.Table("cafeteria_rating_tags").Create(&rating)
	}

	return &emptypb.Empty{}, nil
}

func (s *CampusServer) NewMealRating(ctx context.Context, input *pb.NewRating) (*emptypb.Empty, error) {
	s.db.Where("1=1").Delete(&cafeteria_rating_models.MealRating{})
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

	var dish *cafeteria_rating_models.Dish
	testDish := s.db.Model(cafeteria_rating_models.Dish{Name: input.Meal, Cafeteria: input.CafeteriaName}).First(&dish)

	if testDish.RowsAffected != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "Dish is not offered in this week in this canteen. Rating has not been saved.")
	}

	rating := cafeteria_rating_models.MealRating{
		Comment:   input.Comment,
		Meal:      input.Meal,
		Rating:    input.Rating,
		Timestamp: time.Now()}

	s.db.Model(cafeteria_rating_models.MealRating{}).Create(&rating)

	if len(input.Tags) > 0 {
		//todo tags immer in deutsch speichern -
		//fixme tags in db apspeichern -< beim startup eintragen -> jeder tag hat eine id -> mit dieser tagie abspeichern und zurückgeben
		for i := 0; i < len(input.Tags); i++ {
			//todo add rating to each tag once the proto file is fixed
			//todo check whether the tag is actually included in the json file
			//todo retrieve tag id

			rating := cafeteria_rating_models.CafeteriaTagRating{ParentRating: int32(rating.Id), Rating: int32(5), TagID: i}
			s.db.Table("meal_rating_tags").Create(&rating)
		}
	}

	return &emptypb.Empty{}, nil
}

type MultiLanguageTags struct {
	MultiLanguageTags []Tag `json:"tags"`
}
type Tag struct {
	TagNameEnglish string `json:"tagNameEnglish"`
	TagNameGerman  string `json:"tagNameGerman"`
}

func (s *CampusServer) GetAvailableMealTags(ctx context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	//absPath, _ := filepath.Abs("backend/static_data/mealRatingTags.json")
	//tags := generateTagListFromFile(absPath)

	//todo adapt to query from db
}

func (s *CampusServer) GetAvailableCafeteriaTags(ctx context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	//absPath, _ := filepath.Abs("backend/static_data/cafeteriaRatingTags.json")
	//tags := generateTagListFromFile(absPath)
	//todo adapt to query from db
	/*
		return &pb.GetRatingTagsReply{
			Tags: tags,
		}, nil*/
}

func generateTagListFromFile(path string) MultiLanguageTags {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var tags MultiLanguageTags
	json.Unmarshal(byteValue, &tags)

	return tags
}
