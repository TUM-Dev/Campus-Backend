package main

import (
	"github.com/TUM-Dev/Campus-Backend/backend"
	"github.com/TUM-Dev/Campus-Backend/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/web"
	"log"
	"net"
	"os"

	"github.com/TUM-Dev/Campus-Backend/model"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	httpPort = ":50051"
	grpcPort = ":50052"
)

func main() {
	// Connect to DB
	var conn gorm.Dialector
	shouldAutoMigrate := false
	if dbHost := os.Getenv("DB_DSN"); dbHost != "" {
		log.Printf("Connecting to dsn: %grpcServer", dbHost)
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
		log.Println("Running auto migrations")
		err = db.AutoMigrate(&model.TopNews{})
		if err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
	}

	// Create any other background services (these shouldn't do any long running work here)
	cronService := cron.New(db)

	// Listen to our configured ports
	httpListener, err := net.Listen("tcp", httpPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcListener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Start each server in its own go routine and logs any errors
	g := new(errgroup.Group)
	g.Go(func() error { return web.HTTPServe(httpListener) })
	g.Go(func() error { return backend.GRPCServe(grpcListener) })
	log.Println("run server: ", g.Wait())

	// Setup cron jobs
	g.Go(func() error { return cronService.Run() })
}
