package main

import (
	"log"
	"net"

	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/TUM-Dev/Campus-Backend/server"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

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

	// Migrate the schema
	err = db.AutoMigrate(&model.TopNews{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	cronService := server.CronService{DB: db}
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
