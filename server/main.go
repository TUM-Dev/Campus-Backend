package main

import (
	"github.com/TUM-Dev/Campus-Backend/backend"
	"github.com/TUM-Dev/Campus-Backend/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/TUM-Dev/Campus-Backend/web"
	"github.com/getsentry/sentry-go"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

const (
	httpPort = ":50051"
	grpcPort = ":50052"
)
var Version = "dev"

func main() {
	// Connect to DB
	var conn gorm.Dialector
	shouldAutoMigrate := false
	if dbHost := os.Getenv("DB_DSN"); dbHost != "" {
		log.Printf("Connecting to dsn")
		conn = mysql.Open(dbHost)
	} else {
		conn = sqlite.Open("test.db")
		shouldAutoMigrate = true
	}
	if sentryDSN := os.Getenv("SENTRY_DSN"); sentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: os.Getenv("SENTRY_DSN"),
			Release: Version,
		}); err != nil {
			log.Printf("Sentry initialization failed: %v\n", err)
		}
	} else {
		log.Println("continuing without sentry")
	}
	db, err := gorm.Open(conn, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema only in local development mode
	if shouldAutoMigrate {
		log.Println("Running auto migrations")
		err = db.AutoMigrate(
			&model.TopNews{},
			&model.Crontab{},
			&model.Files{},
			&model.NewsSource{},
			&model.NewsAlert{},
			&model.News{},
		)
		if err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
	}

	// Create any other background services (these shouldn't do any long running work here)
	cronService := cron.New(db)
	campusService := backend.New(db)

	// Listen to our configured ports
	httpListener, err := net.Listen("tcp", httpPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcListener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	g := errgroup.Group{}
	// Start each server in its own go routine and logs any errors
	g.Go(func() error { return web.HTTPServe(httpListener) })
	g.Go(func() error { return campusService.GRPCServe(grpcListener) })

	// Setup cron jobs
	g.Go(func() error { return cronService.Run() })

	log.Println("run server: ", g.Wait())
}
