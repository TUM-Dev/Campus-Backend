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

func LogIOSBackgroundRequest(deviceId, requestType, reason string) {
	write := writeAPI()

	p := influxdb2.NewPointWithMeasurement("ios_campus_token_request").
		AddTag("device_id", deviceId).
		AddField("request_type", requestType).
		AddField("notification_reason_response", reason)

	write.WritePoint(p)

	FlushIfDevelop(write)
}

func LogIOSBackgroundRequestResponse(deviceId, requestType string) {
	write := writeAPI()

	p := influxdb2.NewPointWithMeasurement("ios_campus_token_response").
		AddTag("device_id", deviceId).
		AddField("request_type", requestType)

	write.WritePoint(p)

	FlushIfDevelop(write)
}

func FlushIfDevelop(write api.WriteAPI) {
	if env.IsDev() {
		write.Flush()
	}
}

func writeAPI() api.WriteAPI {
	return GetClient().WriteAPI(influxOrg, influxBucket)
}
