// Package influx provides a handy wrapper around the influxdb client
package influx

import (
	"github.com/TUM-Dev/Campus-Backend/server/env"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"os"
)

var (
	influxOrg    = os.Getenv("INFLUXDB_ORG")
	influxBucket = os.Getenv("INFLUXDB_BUCKET")
)

var Client *influxdb2.Client

func GetClient() influxdb2.Client {
	return *Client
}

func SetClient(client *influxdb2.Client) {
	Client = client
}

func HasClient() bool {
	return Client != nil
}

// Example of how to use the influx client
/* func LogFileDownload() {
	if !HasClient() {
		return
	}

	write := writeAPI()

	p := influxdb2.NewPointWithMeasurement("file_download").
		AddTag("user", "test").
		AddField("file", "test")

	write.WritePoint(p)

	FlushIfDevelop(write)
}
*/

func FlushIfDevelop(write api.WriteAPI) {
	if env.IsDev() {
		write.Flush()
	}
}

func writeAPI() api.WriteAPI {
	return GetClient().WriteAPI(influxOrg, influxBucket)
}
