package main

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
	"time"
)

const (
	address = "api-grpc.tum.app:443"
)

func main() {
	// Set up a connection to the server.
	log.Println("Connecting...")
	/*pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")

	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(creds))
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

	/*log.Println("Trying to fetch top news")
	r, err := c.GetTopNews(ctx, &pb.GetTopNewsRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.String())*/
	createCafeteriaRatingSampleData()
}

func createCafeteriaRatingSampleData() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	c := pb.NewCampusClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	y := make([]string, 3)
	for i := 0; i < 3; i++ {
		y[i] = "Spicy" + strconv.Itoa(i)
	}

	_, errRequest := c.NewCafeteriaRating(ctx, &pb.NewRating{
		Rating:        int32(8),
		CafeteriaName: "MENSA_ARCISSTR",
		Comment:       "Alles Hähnchen",
		Tags:          y,
	})

	if errRequest != nil {
		log.Println(err)
	}

	c.NewMealRating(ctx, &pb.NewRating{
		Rating:        int32(8),
		CafeteriaName: "MENSA_ARCISSTR",
		Meal:          "rinder",
		Comment:       "Alles Hähnchen",
		Tags:          y,
	})

	if errRequest != nil {
		log.Println(err)
	}
}
