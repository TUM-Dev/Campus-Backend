package main

import (
	"github.com/TUM-Dev/Campus-Backend/backend"
	"github.com/TUM-Dev/Campus-Backend/web"
	"log"
	"net"
	"os"

	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	port = ":50051"
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

	// Listen to our configured port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Use cmux so we can server http on the same port
	m := cmux.New(lis)
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	// Start each server in its own go routine and logs any errors
	g := new(errgroup.Group)
	g.Go(func() error { return backend.GRPCServe(grpcListener) })
	g.Go(func() error { return web.HTTPServe(httpListener) })
	g.Go(func() error { return m.Serve() })
	log.Println("run server: ", g.Wait())
}
