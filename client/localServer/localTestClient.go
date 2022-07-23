package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"
)

const (
	localAddress = "127.0.0.1:50051"
)

func main() {
	// Set up a connection to the local server.
	log.Println("Connecting...")

	conn, err := grpc.Dial(localAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	c := pb.NewCampusClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cafeteriaRatingTools(c, ctx)

}

func cafeteriaRatingTools(c pb.CampusClient, ctx context.Context) {

	currentCafeteria := "MENSA_GARCHING"
	currentDish := "Levantinischer Bulgur mit roten Linsen, Spinat und Kichererbsen" //must be in the dish table
	generateCafeteriaRating(c, ctx, currentCafeteria, 2)
	generateCafeteriaRating(c, ctx, currentCafeteria, 2)
	generateCafeteriaRating(c, ctx, currentCafeteria, 2)
	generateDishRating(c, ctx, currentCafeteria, currentDish, 1)

	queryCafeteria(currentCafeteria, c, ctx, false)
	queryDish(currentCafeteria, currentDish, c, ctx, false)

}

func queryDish(cafeteria string, dish string, c pb.CampusClient, ctx context.Context, imageShouldBeStored bool) {
	res, err := c.GetDishRatings(ctx, &pb.DishRatingRequest{
		Dish:          dish,
		CafeteriaName: cafeteria,
		Limit:         3,
	})

	if err != nil {
		println(err)
	} else {
		println("Result: ")
		println("averagerating: ", res.AveragePoints)
		println("min", res.MinPoints)
		println("max", res.MaxPoints)
		println("Number of individual Ratings", len(res.Rating))
		path := fmt.Sprintf("%s%d%s", "./testImages/dishes/", time.Now().Unix(), "/")
		for _, v := range res.Rating {
			println("\nRating: ", v.Points)
			println("Cafeteria Name: ", v.CafeteriaName)
			println("Comment ", v.Comment)
			println("Number of Tag Ratings: ", len(v.TagRating))
			println("Timestamp: ", v.CafeteriaVisitedAt)
			println("ImageLength:", len(v.Image))
			if imageShouldBeStored {
				storeImage(path, v.Image)
			}
		}

		for _, v := range res.RatingTags {
			println("\nNameDE: ", v.DE)
			println("NameEN: ", v.EN)
			println("averagerating: ", v.AveragePoints)
			println("min", v.MinPoints)
			println("max", v.MaxPoints)
		}
		log.Println("nameTags: ")
		for _, v := range res.NameTags {
			println("\nNameDE: ", v.DE)
			println("NameEN: ", v.EN)
			println("averagerating: ", v.AveragePoints)
			println("min", v.MinPoints)
			println("max", v.MaxPoints)
		}
	}
}

func queryCafeteria(s string, c pb.CampusClient, ctx context.Context, imageShouldBeStored bool) {
	res, err := c.GetCafeteriaRatings(ctx, &pb.CafeteriaRatingRequest{
		CafeteriaName: s,
		Limit:         3,
		//	From:          timestamppb.New(time.Date(2022, 7, 8, 16, 0, 0, 0, time.Local)),
		//	To:            timestamppb.New(time.Date(2022, 7, 8, 17, 0, 0, 0, time.Local)),
	})

	if err != nil {
		println(err)
	} else {
		println("Result: ")
		println("averagerating: ", res.AveragePoints)
		println("min", res.MinPoints)
		println("max", res.MaxPoints)
		println("Number of individual Ratings", len(res.Rating))
		path := fmt.Sprintf("%s%d%s", "./testImages/cafeteria/", time.Now().Unix(), "/")
		for _, v := range res.Rating {
			println("\nRating: ", v.Points)
			println("Cafeteria Name: ", v.CafeteriaName)
			println("Comment ", v.Comment)
			println("Number of Tag Ratings: ", len(v.TagRating))
			println("Timestamp: ", v.CafeteriaVisitedAt)
			println("ImageLength:", len(v.Image))
			if imageShouldBeStored {
				storeImage(path, v.Image)
			}
		}

		for _, v := range res.RatingTags {
			println("\nNameDE: ", v.DE)
			println("NameEN: ", v.EN)
			println("averagerating: ", v.AveragePoints)
			println("min", v.MinPoints)
			println("max", v.MaxPoints)
		}
	}
}

func generateCafeteriaRating(c pb.CampusClient, ctx context.Context, cafeteria string, rating int32) {
	y := make([]*pb.TagRating, 2)
	y[0] = &pb.TagRating{
		Points: float64(1 + rating),
		Tag:    "Sauberkeit",
	}
	y[1] = &pb.TagRating{
		Points: float64(2 + rating),
		Tag:    "Enough Free Tables",
	}

	_, err := c.NewCafeteriaRating(ctx, &pb.NewCafeteriaRatingRequest{
		Points:        rating,
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

func generateDishRating(c pb.CampusClient, ctx context.Context, cafeteria string, dish string, rating int32) {
	y := make([]*pb.TagRating, 3)
	y[0] = &pb.TagRating{
		Points: float64(1 + rating),
		Tag:    "Spicy",
	}
	y[1] = &pb.TagRating{
		Points: float64(2 + rating),
		Tag:    "Salz",
	}
	y[2] = &pb.TagRating{
		Points: float64(3 + rating),
		Tag:    "Aussehen",
	}

	_, err := c.NewDishRating(ctx, &pb.NewDishRatingRequest{
		Points:        rating,
		CafeteriaName: cafeteria,
		Dish:          dish,
		Comment:       "Alles Hähnchen",
		Tags:          y,
		Image:         getImageToBytes("../images/sampleimage.jpeg"),
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Request successfully: Dish Rating should be stored")
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
	println("Length of the image as bytes: ", len(bytes))
	return bytes
}

func storeImage(path string, i []byte) (string, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	img, _, _ := image.Decode(bytes.NewReader(i))
	var imgPath = fmt.Sprintf("%s%d%s", path, 3, ".jpeg") //time.Now().Unix()		//use three to force file name collisions
	var f, _ = os.Stat(imgPath)
	var counter = 1
	for {
		if f == nil {
			break
		} else {
			imgPath = fmt.Sprintf("%s%d%s%d%s", path, 3, "v", counter, ".jpeg") //time.Now().Unix()
			counter++
			f, _ = os.Stat(imgPath)
		}
	}

	out, errFile := os.Create(imgPath)
	defer out.Close()
	var opts jpeg.Options
	opts.Quality = 100
	errFile = jpeg.Encode(out, img, &opts)
	return imgPath, errFile
}
