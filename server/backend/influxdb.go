package backend

import (
	"context"
	"errors"
	"github.com/TUM-Dev/Campus-Backend/server/backend/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"os"
)

const (
	// InfluxBatchSize is the number of points that are sent to influxdb at once
	// This is to reduce the number of requests to the database
	InfluxBatchSize = 1_000
)

var (
	influxToken = os.Getenv("INFLUXDB_TOKEN")
	influxURL   = os.Getenv("INFLUXDB_URL")
)

func ConnectToInfluxDB() error {
	if influxToken == "" {
		return errors.New("no influxdb token provided")
	}

	if influxURL == "" {
		return errors.New("no influxdb url provided")
	}

	client := influxdb2.NewClientWithOptions(influxURL, influxToken,
		influxdb2.DefaultOptions().SetBatchSize(InfluxBatchSize))

	influx.SetClient(&client)

	_, err := influx.GetClient().Health(context.Background())

	return err
}
