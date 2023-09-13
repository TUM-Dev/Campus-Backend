package main

import (
	"context"
	"crypto/x509"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
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
		err := conn.Close()
		if err != nil {
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
	if r, err := c.GetTopNews(ctx, &emptypb.Empty{}); err != nil {
		log.WithError(err).Fatal("could not greet")
	} else {
		log.WithField("topNewsResponse", r.String()).Info("fetched top news successfully")
	}
}
