// Package influx provides a handy wrapper around the influxdb client
package influx

import (
	"errors"
	"github.com/TUM-Dev/Campus-Backend/server/env"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	influxOrg    = os.Getenv("INFLUXDB_ORG")
	influxBucket = os.Getenv("INFLUXDB_BUCKET")

	ErrInfluxClientNotConfigured = errors.New("influx client not configured")
)

var Client *influxdb2.Client

func GetClient() influxdb2.Client {
	return *Client
}

func SetClient(client *influxdb2.Client) {
	Client = client
}

/* Example of how to use the influx client
func LogFileDownload() {
	p := influxdb2.NewPointWithMeasurement("file_download").
		AddTag("user", "test").
		AddField("file", "test")

	LogPoint(p)
}
*/

func LogPoint(p *write.Point) {
	w, err := writeAPI()

	if err != nil {
		log.Warn("could not log because influx client is not configured")
		return
	}

	w.WritePoint(p)

	flushIfDevelop(w)
}

func flushIfDevelop(write api.WriteAPI) {
	if env.IsDev() {
		write.Flush()
	}
}

func writeAPI() (api.WriteAPI, error) {
	if Client != nil {
		return GetClient().WriteAPI(influxOrg, influxBucket), nil
	}

	return nil, ErrInfluxClientNotConfigured
}
