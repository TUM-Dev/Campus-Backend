package main

import (
	"context"
	"crypto/x509"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	address = "api-grpc.tum.app:443"
)

func main() {
	// Set up a connection to the server.
	log.Info("Connecting...")
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.WithError(err).Fatal("did not connect")
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.WithError(err).Error("did not close connection")
		}
	}(conn)
	c := pb.NewCampusClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Add header metadata
	md := metadata.New(map[string]string{"x-device-id": "grpc-tests"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	log.Info("Trying to fetch top news")
	if r, err := c.GetTopNews(ctx, &pb.GetNewsRequest{}); err != nil {
		log.WithError(err).Fatal("could not greet")
	} else {
		log.WithField("topNewsResponse", r.String()).Info("fetched top news successfully")
	}
}
