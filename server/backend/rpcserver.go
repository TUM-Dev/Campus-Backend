package backend

import (
	"context"
	"github.com/TUM-Dev/Campus-Backend/model"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

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

func (s *CampusServer) GetTopNews(ctx context.Context, in *pb.GetTopNewsRequest) (*pb.GetTopNewsReply, error) {
	log.Printf("Received: get top news")
	var res *model.NewsAlert
	err := s.db.Joins("Company").Where("NOW() between `from` and `to`").Limit(1).First(&res).Error
	if err != nil {
		log.Error(err)
	} else if res != nil {
		return &pb.GetTopNewsReply{
			//ImageUrl: res.Name,
			Link: res.Link.String,
			To:   timestamppb.New(res.To),
		}, nil
	}

	now := timestamppb.New(time.Now())
	return &pb.GetTopNewsReply{
		Name:    "Test Top News",
		Link:    "https://google.com",
		Created: now,
		From:    nil,
		To:      nil,
	}, nil
}
