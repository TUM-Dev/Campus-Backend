package backend

import (
	"context"
	"errors"
	"github.com/TUM-Dev/Campus-Backend/server/backend/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"os"
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

	client := influxdb2.NewClient(influxURL, influxToken)

	influx.SetClient(&client)

	_, err := influx.GetClient().Health(context.Background())

	return err
}
