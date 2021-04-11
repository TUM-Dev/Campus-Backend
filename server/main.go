package main

import (
	"github.com/TUM-Dev/Campus-Backend/backend"
	"log"
	"net"
	"os"

	"github.com/TUM-Dev/Campus-Backend/model"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/TUM-Dev/Campus-Backend/api"
)

const (
	port = ":50051"
)

func main() {
	// Connect to DB
	var conn gorm.Dialector
	shouldAutoMigrate := false
	if dbHost := os.Getenv("DB_DSN"); dbHost != "" {
		log.Printf("Connecting to dsn: %s", dbHost)
		conn = mysql.Open(dbHost)
	} else {
		conn = sqlite.Open("test.db")
		shouldAutoMigrate = true
	}
	db, err := gorm.Open(conn, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema only in local development mode
	if shouldAutoMigrate {
		err = db.AutoMigrate(&model.TopNews{})
		if err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
	}

	// Start Server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCampusServer(s, &backend.CampusServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
