package main

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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

	/*generateCafeteriaRating(c, ctx, "MENSA_GARCHING", 1)
	generateCafeteriaRating(c, ctx, "MENSA_GARCHING", 5)
	generateCafeteriaRating(c, ctx, "MENSA_GARCHING", 5)
	generateCafeteriaRating(c, ctx, "MENSA_GARCHING", 5)
	generateCafeteriaRating(c, ctx, "FMI_BISTRO", 5)
	generateCafeteriaRating(c, ctx, "FMI_BISTRO", 5)*/
	/*generateMealRating(c, ctx, "MENSA_GARCHING", "Bio-Pasta mit Bio-Tomaten-Frischkäse-Sauce", 1)
	generateMealRating(c, ctx, "MENSA_GARCHING", "Bio-Pasta mit Bio-Tomaten-Frischkäse-Sauce", 3)
	generateMealRating(c, ctx, "MENSA_GARCHING", "Bio-Pasta mit Bio-Tomaten-Frischkäse-Sauce", 5)
	generateMealRating(c, ctx, "MENSA_GARCHING", "Bio-Pasta mit Bio-Tomaten-Frischkäse-Sauce", 7)*/

	//generateMealRating(c, ctx, "MENSA_GARCHING", "Pasta all'arrabiata", 2)
	//generateCafeteriaRating(c, ctx, "MENSA_GARCHING", 2)
	res, err := c.GetCafeteriaRatings(ctx, &pb.CafeteriaRatingRequest{

		CafeteriaName: "MENSA_GARCHING",
		Limit:         3,
	})
	println(res.AverageRating)
	if err != nil {
		println(err)
	}
}

func generateCafeteriaRating(c pb.CampusClient, ctx context.Context, cafeteria string, rating int32) {
	y := make([]*pb.TagRating, 2)
	y[0] = &pb.TagRating{
		Rating: 4,
		Tag:    "Sauberkeit",
	}
	y[1] = &pb.TagRating{
		Rating: 2,
		Tag:    "Anzahl an Plätzen",
	}

	_, err := c.NewCafeteriaRating(ctx, &pb.NewCafeteriaRatingRequest{
		Rating:        rating,
		CafeteriaName: cafeteria,
		Comment:       "Alles super, 2 Sterne",
		Tags:          y,
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Request successfully: Cafeteria Rating should be stored")
	}
}

func generateMealRating(c pb.CampusClient, ctx context.Context, cafeteria string, meal string, rating int32) {
	y := make([]*pb.TagRating, 3)
	y[0] = &pb.TagRating{
		Rating: 2,
		Tag:    "Spicy",
	}
	y[1] = &pb.TagRating{
		Rating: 2,
		Tag:    "Salz",
	}
	y[2] = &pb.TagRating{
		Rating: 2,
		Tag:    "Aussehen",
	}

	_, err := c.NewMealRating(ctx, &pb.NewMealRatingRequest{
		Rating:        rating,
		CafeteriaName: cafeteria,
		Meal:          meal,
		Comment:       "Alles Hähnchen",
		Tags:          y,
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Request successfully: Meal Rating should be stored")
	}
}
