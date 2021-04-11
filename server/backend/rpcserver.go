package backend

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "github.com/TUM-Dev/Campus-Backend/api"
)

func GRPCServe(l net.Listener) error {
	s := grpc.NewServer()
	pb.RegisterCampusServer(s, &CampusServer{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return s.Serve(l)
}

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
