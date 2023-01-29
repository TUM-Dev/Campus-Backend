// Package influx provides a handy wrapper around the influxdb client
package influx

import (
	"crypto/sha256"
	"encoding/hex"
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

func LogPoint(p *write.Point) {
	w, err := writeAPI()

	if err != nil {
		logClientNotConfigured()
		return
	}

	w.WritePoint(p)

	flushIfDevelop(w)
}

func hashSha256(s string) string {
	h := sha256.New()

	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}

func logClientNotConfigured() {
	log.Warn("could not log because influx client is not configured")
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
