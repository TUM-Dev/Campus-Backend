package backend

import (
	"context"
	"errors"
	"os"

	"github.com/TUM-Dev/Campus-Backend/server/backend/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	// InfluxBatchSize is the number of points that are sent to influxdb at once
	// This is to reduce the number of requests to the database
	InfluxBatchSize = 500
)

var (
	influxToken = os.Getenv("INFLUXDB_TOKEN")
	influxURL   = os.Getenv("INFLUXDB_URL")

	ErrInfluxTokenNotConfigured = errors.New("influxdb token not configured")
	ErrInfluxURLNotConfigured   = errors.New("influxdb url not configured")
)

func ConnectToInfluxDB() error {
	if influxToken == "" {
		return ErrInfluxTokenNotConfigured
	}

	if influxURL == "" {
		return ErrInfluxURLNotConfigured
	}

	client := influxdb2.NewClientWithOptions(influxURL, influxToken,
		influxdb2.DefaultOptions().SetBatchSize(InfluxBatchSize))

	influx.SetClient(&client)

	_, err := influx.GetClient().Health(context.Background())

	return err
}
