package backend

import (
	"context"
	"github.com/TUM-Dev/Campus-Backend/model"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"

	pb "github.com/TUM-Dev/Campus-Backend/api"
)

func (s *CampusServer) GRPCServe(l net.Listener) error {
	grpc := grpc.NewServer()
	pb.RegisterCampusServer(grpc, s)
	if err := grpc.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return grpc.Serve(l)
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
	var res *model.TopNews
	err := s.db.Joins("Company").Where("NOW() between from and to ").Limit(1).First(&res).Error
	if err != nil {
		log.Error(err)
	}

	now, _ := ptypes.TimestampProto(time.Now())
	if res != nil {
		return &pb.GetTopNewsReply{
			Name:    "Test Top News",
			Link:    "https://google.com",
			Created: now,
			From:    nil,
			To:      nil,
		}, nil
	}

	return &pb.GetTopNewsReply{
		Name:    "Test Top News",
		Link:    "https://google.com",
		Created: now,
		From:    nil,
		To:      nil,
	}, nil
}
