package server

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
)

type CampusServer struct {
	pb.UnimplementedCampusServer
}

func (s *CampusServer) GetTopNews(ctx context.Context, in *pb.GetTopNewsRequest) (*pb.GetTopNewsReply, error) {
	log.Printf("Received: get top news")

	now, _ := ptypes.TimestampProto(time.Now())
	return &pb.GetTopNewsReply{
		Name:    "Test Top News",
		Link:    "https://google.com",
		Created: now,
		From:    nil,
		To:      nil,
	}, nil
}
