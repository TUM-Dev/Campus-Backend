package main

import (
	"context"
	"crypto/x509"
	"log"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	// "google.golang.org/protobuf/types/known/emptypb"
)

const (
	address = "api-grpc.tum.app:443"
)

func main() {
	// Set up a connection to the server.
	log.Println("Connecting...")
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
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
	request := pb.APRequest{Name: "apa05-0mg",Timestamp: "2022-06-15 22"}
	r, err := c.GetAccessPoint(ctx, &request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.Name)
}
