package main

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

const (
	address = "api.tum.app:50052"
)

func main() {
	// Set up a connection to the server.
	log.Println("Connecting...")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCampusClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Add header metadata
	md := metadata.New(map[string]string{"x-device-id": "grpc-tests"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	log.Println("Trying to fetch top news")
	r, err := c.GetTopNews(ctx, &pb.GetTopNewsRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.String())
}
