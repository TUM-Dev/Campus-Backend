package main

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/TUM-Dev/Campus-Backend/server"
	"github.com/TUM-Dev/Campus-Backend/server/cron"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/api"
)

const (
	port = ":50051"
)

func main() {
	// Connect to DB
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: os.Getenv("SentryDSN"),
	})
	// Migrate the schema
	err = db.AutoMigrate(
		&model.TopNews{},
		&model.RoomfinderRooms{},
		&model.RoomfinderMaps{},
	)
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	if os.Getenv("SentryDSN") == "" {
		log.Fatalln("couldn't start without env variable \"SentryDSN\"")
	}
	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SentryDSN"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")

	// initialize cron service
	cronService := cron.ServiceCron{DB: db}
	cronService.Init()

	// Start Server
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
