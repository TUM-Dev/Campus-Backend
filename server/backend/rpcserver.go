package backend

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"net"
	"sync"
	"time"

	"strconv"
	"strings"
	"github.com/TUM-Dev/Campus-Backend/model/heatmap/dbservice"
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

const (
	heatmapDB = "./data/sqlite/heatmap.db"
)

func (s *CampusServer) GetAccessPoint(ctx context.Context, in *pb.APRequest) (*pb.AccessPoint, error) {
	name := in.Name
	ts := in.Timestamp
	log.Printf("Received request for AP with name: %s and timestamp: %s", name, in.Timestamp)

	db := dbservice.InitDB(heatmapDB)

	day, hr := getDayAndHourFromTimestamp(ts)
	ap := dbservice.GetHistoryForSingleAP(name, day, hr)
	db.Close()

	load, _ := strconv.Atoi(ap.Load)

	return &pb.AccessPoint{
		Name:      ap.Name,
		Lat:       ap.Lat,
		Long:      ap.Long,
		Intensity: int64(load),
		Max:       int64(ap.Max),
		Min:       int64(ap.Min),
	}, nil
}

func getDayAndHourFromTimestamp(timestamp string) (int, int) {
	ts := strings.Split(timestamp, " ")
	date := ts[0]
	yearMonthDay := strings.Split(date, "-")
	day, err := strconv.Atoi(yearMonthDay[2])
	if err != nil {
		day = 0
	}
	hr, err := strconv.Atoi(ts[1])
	if err != nil {
		hr = 0
	}
	today := time.Now().Day()
	day = day - today
	return day, hr
}

type location struct {
	lat  string
	long string
}

var locations map[string]location

func (s *CampusServer) ListAccessPoints(in *pb.APRequest, stream pb.Campus_ListAccessPointsServer) error {
	ts := in.Timestamp
	day, hr := getDayAndHourFromTimestamp(ts)

	apList := dbservice.GetHistoryForAllAPs(day, hr)

	log.Printf("Sending %d APs ...", len(apList))

	locations = make(map[string]location)
	accessPoints := dbservice.RetrieveAPsOfTUM(true)
	for _, ap := range accessPoints {
		locations[ap.Name] = location{ap.Lat, ap.Long}
	}

	for _, ap := range apList {
		location := locations[ap.Name]

		load, _ := strconv.Atoi(ap.Load)

		if err := stream.Send(
			&pb.APResponse{
				Accesspoint: &pb.AccessPoint{
					Name:      ap.Name,
					Lat:       location.lat,
					Long:      location.long,
					Intensity: int64(load),
					Max:       int64(ap.Max),
					Min:       int64(ap.Min),
				},
			}); err != nil {
			return err
		}
	}

	return nil
}

func (s *CampusServer) ListAllAPNames(in *emptypb.Empty, stream pb.Campus_ListAllAPNamesServer) error {
	names := dbservice.GetAllNames()
	for _, name := range names {
		if err := stream.Send(
			&pb.APName{
				Name: name,
			}); err != nil {
			return err
		}
	}
	return nil
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
