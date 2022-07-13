package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
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
	generateCafeteriaRating(c, ctx, "MENSA_GARCHING", 3)
	generateCafeteriaRating(c, ctx, "MENSA_GARCHING", 5)
	generateCafeteriaRating(c, ctx, "MENSA_GARCHING", 7)
	generateCafeteriaRating(c, ctx, "FMI_BISTRO", 4)
	generateCafeteriaRating(c, ctx, "FMI_BISTRO", 6)*/
	/*	generateMealRating(c, ctx, "MENSA_GARCHING", "Levantinischer Bulgur mit roten Linsen, Spinat und Kichererbsen", 1)
		generateMealRating(c, ctx, "MENSA_GARCHING", "Levantinischer Bulgur mit roten Linsen, Spinat und Kichererbsen", 3)
		generateMealRating(c, ctx, "MENSA_GARCHING", "Levantinischer Bulgur mit roten Linsen, Spinat und Kichererbsen", 5)
	*/
	//	generateMealRating(c, ctx, "MENSA_GARCHING", "Levantinischer Bulgur mit roten Linsen, Spinat und Kichererbsen", 1)

	//generateMealRating(c, ctx, "MENSA_GARCHING", "Pasta all'arrabiata", 2)
	generateCafeteriaRating(c, ctx, "MENSA_GARCHING", 2)

	queryCafeteria("MENSA_GARCHING", c, ctx)
	//queryMeal("MENSA_GARCHING", "Levantinischer Bulgur mit roten Linsen, Spinat und Kichererbsen", c, ctx)

}

func queryMeal(cafeteria string, meal string, c pb.CampusClient, ctx context.Context) {
	res, err := c.GetMealRatings(ctx, &pb.MealRatingsRequest{
		Meal:          meal,
		CafeteriaName: cafeteria,
		Limit:         3,
	})

	println("Result: ")
	println("averagerating: ", res.AverageRating)
	println("min", res.MinRating)
	println("max", res.MaxRating)
	println("Number of individual Ratings", len(res.Rating))
	for _, v := range res.Rating {
		println("\nRating: ", v.Rating)
		println("Cafeteria Name: ", v.CafeteriaName)
		println("Comment ", v.Comment)
		println("Number of Tag Ratings: ", len(v.TagRating))
		println("Timestamp: ", v.CafeteriaVisitedAt)
	}

	for _, v := range res.RatingTags {
		println("\nNameDE: ", v.NameDE)
		println("NameEN: ", v.NameEN)
		println("averagerating: ", v.AverageRating)
		println("min", v.MinRating)
		println("max", v.MaxRating)
	}
	log.Println("nameTags: ")
	for _, v := range res.NameTags {
		println("\nNameDE: ", v.NameDE)
		println("NameEN: ", v.NameEN)
		println("averagerating: ", v.AverageRating)
		println("min", v.MinRating)
		println("max", v.MaxRating)
	}
	if err != nil {
		println(err)
	}
}

func queryCafeteria(s string, c pb.CampusClient, ctx context.Context) {
	res, err := c.GetCafeteriaRatings(ctx, &pb.CafeteriaRatingRequest{

		CafeteriaName: s,
		Limit:         3,
		From:          timestamppb.New(time.Date(2022, 7, 8, 16, 0, 0, 0, time.Local)),
		To:            timestamppb.New(time.Date(2022, 7, 8, 17, 0, 0, 0, time.Local)),
		//From:          timestamppb.New(time.Date(2021, 7, 8, 16, 0, 0, 0, time.Local)),
		//To:            timestamppb.New(time.Date(2021, 7, 8, 17, 0, 0, 0, time.Local)),
	})

	println("Result: ")
	println("averagerating: ", res.AverageRating)
	println("min", res.MinRating)
	println("max", res.MaxRating)
	println("Number of individual Ratings", len(res.Rating))
	for _, v := range res.Rating {
		println("\nRating: ", v.Rating)
		println("Cafeteria Name: ", v.CafeteriaName)
		println("Comment ", v.Comment)
		println("Number of Tag Ratings: ", len(v.TagRating))
		println("Timestamp: ", v.CafeteriaVisitedAt)
	}

	for _, v := range res.RatingTags {
		println("\nNameDE: ", v.NameDE)
		println("NameEN: ", v.NameEN)
		println("averagerating: ", v.AverageRating)
		println("min", v.MinRating)
		println("max", v.MaxRating)
	}
	if err != nil {
		println(err)
	}
}

func generateCafeteriaRating(c pb.CampusClient, ctx context.Context, cafeteria string, rating int32) {
	y := make([]*pb.TagRating, 2)
	y[0] = &pb.TagRating{
		Rating: float64(1 + rating),
		Tag:    "Sauberkeit",
	}
	y[1] = &pb.TagRating{
		Rating: float64(2 + rating),
		Tag:    "Enough Free Tables",
	}

	_, err := c.NewCafeteriaRating(ctx, &pb.NewCafeteriaRatingRequest{
		Rating:        rating,
		CafeteriaName: cafeteria,
		Comment:       "Alles super, 2 Sterne",
		Tags:          y,
		Image:         getImageToBytes("../images/sampleimage.jpeg"),
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
		Rating: float64(1 + rating),
		Tag:    "Spicy",
	}
	y[1] = &pb.TagRating{
		Rating: float64(2 + rating),
		Tag:    "Salz",
	}
	y[2] = &pb.TagRating{
		Rating: float64(3 + rating),
		Tag:    "Aussehen",
	}

	_, err := c.NewMealRating(ctx, &pb.NewMealRatingRequest{
		Rating:        rating,
		CafeteriaName: cafeteria,
		Meal:          meal,
		Comment:       "Alles Hähnchen",
		Tags:          y,
		Image:         getImageToBytes("../images/sampleimage.jpeg"),
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Request successfully: Meal Rating should be stored")
	}
}

func getImageToBytes(path string) []byte {

	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	return bytes
}
