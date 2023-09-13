package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	localAddress = "127.0.0.1:50051"
	testImage    = "./localServer/images/sampleimage.jpeg"
)

// main connects to a seperatly started local server and creates ratings for both, cafeterias and dishes.
// Afterwards, they are queried and displayed on the console
func main() {
	// Set up a connection to the local server.
	log.Info("Connecting...")

	conn, err := grpc.Dial(localAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithError(err).Error("could not dial localAddress")
	}
	c := pb.NewCampusClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	canteenHeadCount(c, ctx)

	cafeteriaRatingTools(c, ctx)

}

func canteenHeadCount(c pb.CampusClient, ctx context.Context) {
	res, err := c.GetCanteenHeadCount(ctx, &pb.GetCanteenHeadCountRequest{
		CanteenId: "mensa-garching",
	})

	if err != nil {
		log.WithError(err).Error("Canteen HeadCount data request failed.")
	} else {
		log.WithField("res", res).Info("Canteen HeadCount data request successful.")
	}
}

func cafeteriaRatingTools(c pb.CampusClient, ctx context.Context) {

	currentCafeteria := "MENSA_GARCHING"
	currentDish := "Vegane rote Grütze mit Soja-Vanillesauce" //must be in the dish table
	generateDishRating(c, ctx, currentCafeteria, currentDish, 3)
	generateCafeteriaRating(c, ctx, currentCafeteria, 2)
	queryCafeteria(currentCafeteria, c, ctx, true)
	queryDish(currentCafeteria, currentDish, c, ctx, false)
	generateCafeteriaRating(c, ctx, currentCafeteria, 2)
	generateCafeteriaRating(c, ctx, currentCafeteria, 2)
	generateDishRating(c, ctx, currentCafeteria, currentDish, 1)

	queryCafeteria(currentCafeteria, c, ctx, false)
	queryDish(currentCafeteria, currentDish, c, ctx, false)

}

func queryDish(cafeteria string, dish string, c pb.CampusClient, ctx context.Context, imageShouldBeStored bool) {
	res, err := c.GetDishRatings(ctx, &pb.DishRatingRequest{
		Dish:        dish,
		CafeteriaId: cafeteria,
		Limit:       3,
	})

	if err != nil {
		log.WithError(err).Info("failed to query dish")
	} else {
		fields := log.Fields{
			"avg":               res.Avg,
			"min":               res.Min,
			"max":               res.Max,
			"std":               res.Std,
			"individualRatings": len(res.Rating),
		}
		log.WithFields(fields).Info("succeeded to query dish")
		path := fmt.Sprintf("%s%d%s", "./testImages/dishes/", time.Now().Unix(), "/")
		for _, v := range res.Rating {
			fields := log.Fields{
				"Rating":                v.Points,
				"Comment":               v.Comment,
				"Number of Tag Ratings": len(v.RatingTags),
				"Timestamp":             v.Visited,
				"ImageLength":           len(v.Image),
			}
			log.WithFields(fields).Info("storing image")
			if imageShouldBeStored {
				_, err := storeImage(path, v.Image)
				if err != nil {
					log.WithError(err).Error("image was not saved successfully")
				}
			}
		}
		log.Info("Rating Tags: ")
		for _, v := range res.RatingTags {
			fields := log.Fields{
				"avg": v.Avg,
				"min": v.Min,
				"max": v.Max,
				"std": v.Std,
			}
			log.WithFields(fields).Info(v.TagId)
		}
		log.Info("nameTags: ")
		for _, v := range res.NameTags {
			fields := log.Fields{
				"avg": v.Avg,
				"min": v.Min,
				"max": v.Max,
				"std": v.Std,
			}
			log.WithFields(fields).Info(v.TagId)
		}
	}
}

func queryCafeteria(s string, c pb.CampusClient, ctx context.Context, imageShouldBeStored bool) {
	res, err := c.GetCafeteriaRatings(ctx, &pb.CafeteriaRatingRequest{
		CafeteriaId: s,
		Limit:       3,
		//	From:          timestamppb.New(time.Date(2022, 7, 8, 16, 0, 0, 0, time.Local)),
		//	To:            timestamppb.New(time.Date(2022, 7, 8, 17, 0, 0, 0, time.Local)),
	})

	if err != nil {
		log.WithError(err).Error("failed to query cafeteria")
	} else {
		fields := log.Fields{
			"avg":                          res.Avg,
			"min":                          res.Min,
			"max":                          res.Max,
			"std":                          res.Std,
			"Number of individual Ratings": len(res.Rating),
		}
		log.WithFields(fields).Info("succeeded to query cafeteria")
		path := fmt.Sprintf("%s%d%s", "./testImages/cafeteria/", time.Now().Unix(), "/")
		for i, v := range res.Rating {
			fields := log.Fields{
				"Rating":                v.Points,
				"Comment":               v.Comment,
				"Number of Tag Ratings": len(v.RatingTags),
				"Timestamp":             v.Visited,
				"ImageLength":           len(v.Image),
			}
			log.WithFields(fields).Infof("Rating %d", i)
			if imageShouldBeStored {
				_, err := storeImage(path, v.Image)
				if err != nil {
					log.WithError(err).Error("image was not saved successfully")
				}
			}
		}

		for _, v := range res.RatingTags {
			fields := log.Fields{
				"avg": v.Avg,
				"min": v.Min,
				"max": v.Max,
				"std": v.Std,
			}
			log.WithFields(fields).Info(v.TagId)
		}
	}
}

func generateCafeteriaRating(c pb.CampusClient, ctx context.Context, cafeteria string, rating int32) {
	y := make([]*pb.RatingTag, 2)
	y[0] = &pb.RatingTag{
		Points: float64(1 + rating),
		TagId:  1,
	}
	y[1] = &pb.RatingTag{
		Points: float64(2 + rating),
		TagId:  2,
	}

	_, err := c.NewCafeteriaRating(ctx, &pb.NewCafeteriaRatingRequest{
		Points:      rating,
		CafeteriaId: cafeteria,
		Comment:     "Alles super, 2 Sterne",
		RatingTags:  y,
		Image:       getImageToBytes(testImage),
	})

	if err != nil {
		log.WithError(err).Error("could not store new Cafeteria Rating")
	} else {
		log.Info("Cafeteria Rating successfully be stored")
	}
}

func generateDishRating(c pb.CampusClient, ctx context.Context, cafeteria string, dish string, rating int32) {
	y := make([]*pb.RatingTag, 3)
	y[0] = &pb.RatingTag{
		Points: float64(1 + rating),
		TagId:  1,
	}
	y[1] = &pb.RatingTag{
		Points: float64(2 + rating),
		TagId:  2,
	}
	y[2] = &pb.RatingTag{
		Points: float64(3 + rating),
		TagId:  3,
	}

	_, err := c.NewDishRating(ctx, &pb.NewDishRatingRequest{
		Points:      rating,
		CafeteriaId: cafeteria,
		Dish:        dish,
		Comment:     "Alles Hähnchen",
		RatingTags:  y,
		Image:       getImageToBytes(testImage),
	})

	if err != nil {
		log.WithError(err).Error("failed to store dish rating")
	} else {
		log.Info("Dish Rating successfully stored")
	}
}

func getImageToBytes(path string) []byte {

	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return make([]byte, 0)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.WithError(err).Error("could not close file")
		}
	}(file)

	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	byteArray := make([]byte, size)

	n, err := bufio.NewReader(file).Read(byteArray)
	fields := log.Fields{"readBytes": n, "len(byteArray)": len(byteArray)}
	if err != nil {
		log.WithError(err).Error("Unable to read the byteArray")
	} else {
		log.WithFields(fields).Info("read image to byteArray")
	}
	return byteArray
}

func storeImage(path string, i []byte) (string, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.WithError(err).Error("could not make dir")
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
	if errFile != nil {
		log.WithError(errFile).Error("Unable to create the new testfile")
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.WithError(err).Error("File was not closed successfully")
		}
	}(out)
	var opts jpeg.Options
	opts.Quality = 100
	errFile = jpeg.Encode(out, img, &opts)
	return imgPath, errFile
}
