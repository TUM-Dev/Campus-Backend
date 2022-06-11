package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/TUM-Dev/Campus-Backend/model/mensa_rating_models"
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
	return &CampusServer{
		db: db,
		deviceBuf: &deviceBuffer{
			lock:     sync.Mutex{},
			devices:  make(map[string]*model.Devices),
			interval: time.Minute,
		},
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
	return nil, status.Errorf(codes.Unimplemented, "method GetCafeteriaRatingLastThree not implemented but I am working on it")
}
func (s *CampusServer) GetMealRatingLastThree(ctx context.Context, _ *pb.GetMealInCafeteriaRating) (*pb.GetMealInCafeteriaRatingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMealRatingLastThree not implemented but I am working on it")
}
func (s *CampusServer) NewCafeteriaRating(ctx context.Context, _ *pb.NewRating) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewCafeteriaRating not implemented but I am working on it")
}

func (s *CampusServer) NewMealRating(ctx context.Context, input *pb.NewRating) (*emptypb.Empty, error) {

	//Add cafeteriaRating
	if input.Rating > 10 || input.Rating < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Rating must be a positive number not larger than 10. Rating has not been saved.")
	}

	if len(input.Image) > 131100 {
		return nil, status.Errorf(codes.InvalidArgument, "Image must not be larger than 1MB. Rating has not been saved.")
	}

	if input.Meal == "" || len(input.Meal) > 128 { //todo check if it actually exists in the daily meal names
		return nil, status.Errorf(codes.InvalidArgument, "Image must not be larger than 1MB. Rating has not been saved.")
	}

	var result *model.Mensa
	test := s.db.Model(model.Mensa{Name: input.CafeteriaName}).First(&result)

	if test.RowsAffected != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "Mensa does not exist. Rating has not been saved.")
	}

	rating := mensa_rating_models.CafeteriaRating{
		Comment:   input.Comment,
		Meal:      input.Meal,
		Rating:    input.Rating,
		Timestamp: time.Now()}

	s.db.Table("mensa_garching_rating").Create(&rating)

	var parentid = rating.Id
	//Add Tag Ratings for the first cafeteria

	for i := 0; i < len(input.Tags); i++ {

		//todo add rating once the proto file is fixed
		rating := mensa_rating_models.TagRating{ParentRating: int32(parentid), Rating: int32(5), Tagname: input.Tags[i]}
		s.db.Table("mensa_garching_tags").Create(&rating)
	}

	/*	var retrieved *mensa_rating_models.CafeteriaRating
		s.db.Table("mensa_garching_rating").First(&retrieved)
		log.Println("First comment: ")
		log.Println(retrieved.Comment)*/
	//	s.db.Table("mensa_garching_rating").Raw("INSERT INTO mensa_garching_rating (rating, comment)  VALUES (`ratingFirst`,`comment);").Scan(&result)
	//result := s.db.Table("mensa_garching_rating").Create()
	return &emptypb.Empty{}, nil //nil, status.Errorf(codes.Unimplemented, "method NewMealRating not implemented but I am working on it")
}

type MultiLanguageTags struct {
	MultiLanguageTags []Tag `json:"tags"`
}
type Tag struct {
	TagNameEnglish string `json:"tagNameEnglish"`
	TagNameGerman  string `json:"tagNameGerman"`
}

func (s *CampusServer) GetAvailableMealTags(ctx context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	absPath, _ := filepath.Abs("backend/static_data/mealRatingTags.json")
	tags := generateTagListFromFile(absPath)

	return &pb.GetRatingTagsReply{
		Tags: tags,
	}, nil
}

func (s *CampusServer) GetAvailableCafeteriaTags(ctx context.Context, _ *emptypb.Empty) (*pb.GetRatingTagsReply, error) {
	absPath, _ := filepath.Abs("backend/static_data/cafeteriaRatingTags.json")
	tags := generateTagListFromFile(absPath)

	return &pb.GetRatingTagsReply{
		Tags: tags,
	}, nil
}

func generateTagListFromFile(path string) []string {
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

	var helper = len(tags.MultiLanguageTags)
	y := make([]string, helper)
	for i := 0; i < len(tags.MultiLanguageTags); i++ {
		y[i] = tags.MultiLanguageTags[i].TagNameEnglish
	}
	return y
}

/*
func (s *CampusServer) GetTopNews(ctx context.Context, _ *emptypb.Empty) (*pb.GetTopNewsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}
	log.Printf("Received: get top news adaption")
	var res *model.NewsAlert
	//s.db.Table("roles")
	s.db.Table("mensa_garching_rating").Raw("INSERT INTO mensa_garching_rating (rating, comment)  VALUES (`ratingFirst`,`comment);")

	test := s.db.Table("roles").First("roles")
	log.Println(test.Error)
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

	s.db.Table("mensa_garching_rating").Raw("INSERT INTO mensa_garching_rating (rating, comment)  VALUES (`ratingFirst`,`comment);")
	return &pb.GetTopNewsReply{}, nil
}*/

/*
type GetMensaRatingReply struct {
	arch_id string
}
*/
/*func (s *CampusServer) GetMensaRating(ctx context.Context, _ *emptypb.Empty) (*pb.GetRoomCoordinatesRequest, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}
	log.Printf("Received: mensa rating")
	var res *model.NewsAlert
	err := s.db.Joins("Company").Where("NOW() between `from` and `to`").Limit(1).First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("Failed to fetch top news: %w", err)
	} else if res != nil {
		return &pb.GetRoomCoordinatesRequest{
			//ImageUrl: res.Name,
			ArchId: string("abcdefg"),
		}, nil
	}
	return &pb.GetRoomCoordinatesRequest{}, nil
}*/
/*
func (s *CampusServer) GetMensaRating(context.Context, *emptypb.Empty) (*pb.GetRoomCoordinatesRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMensaRating not implemented ahhh")
}

func (s *CampusServer) GetRoomSchedule(context.Context, *pb.GetRoomScheduleRequest) (*pb.GetRoomScheduleReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomSchedule not implemented but I am working on it")
}*/
