package main

import (
	"github.com/TUM-Dev/Campus-Backend/server"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/TUM-Dev/Campus-Backend/api"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCampusServer(s, &server.CampusServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
