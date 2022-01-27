package backend

import (
	"context"
	"errors"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"net"

	pb "github.com/TUM-Dev/Campus-Backend/api"
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
	db *gorm.DB
}

func New(db *gorm.DB) *CampusServer {
	return &CampusServer{
		db: db,
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

// checkDevice checks if the device is approved (TODO: implement)
func (s *CampusServer) checkDevice(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Internal, "can't extract metadata from request")
	}
	if len(md["x-device-id"]) == 0 && len(md["grpcgateway-referer"]) == 0 && md["x-forwarded-for"][0] != "::1" {
		return status.Errorf(codes.PermissionDenied, "no device id")
	}
	log.WithField("DeviceID", md["x-device-id"]).Info("Request from device")
	return nil
}
