package main

import (
	"context"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	localAddress = "127.0.0.1:50051"
	testImage    = "./localServer/images/sampleimage.jpeg"
)

// main connects to a separately started local server an allows to test things
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

	// add your test here
	feedback, err := c.SendFeedback(ctx, &pb.SendFeedbackRequest{
		Topic:      "tca",
		Email:      "frank@elsinga.de",
		EmailId:    "magic id",
		Message:    "hi, comrades",
		ImageCount: 1,
		OsVersion:  "manjaro",
		AppVersion: "v2",
	})
	if err != nil {
		log.WithError(err).Error("error sending feedback occurred")
	} else {
		log.Info("Success ", feedback.Status)
	}
}
