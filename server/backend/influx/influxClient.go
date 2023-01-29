// Package influx provides a handy wrapper around the influxdb client
package influx

import (
	"crypto/sha256"
	"encoding/hex"
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

func LogRegisterDevice(deviceId string) {
	write := writeAPI()

	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_register_device").
		AddTag("device_id", hashedDeviceId)

	write.WritePoint(p)

	FlushIfDevelop(write)
}

func LogRemoveDevice(deviceId string) {
	write := writeAPI()

	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_remove_device").
		AddTag("device_id", hashedDeviceId)

	write.WritePoint(p)

	FlushIfDevelop(write)
}

func LogIOSNewGrades(deviceId string, gradesCount int) {
	write := writeAPI()

	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_new_grades").
		AddTag("device_id", hashedDeviceId).
		AddField("new_grades_count", gradesCount)

	write.WritePoint(p)

	FlushIfDevelop(write)
}

func LogIOSSchedulingDevicesToUpdate(devicesToUpdateCount int, priority int) {
	write := writeAPI()

	p := influxdb2.NewPointWithMeasurement("ios_scheduling_devices_to_update").
		AddTag("priority", string(rune(priority))).
		AddField("devices_to_update", devicesToUpdateCount)

	write.WritePoint(p)

	FlushIfDevelop(write)
}

func LogIOSBackgroundRequest(deviceId, requestType, reason string) {
	write := writeAPI()

	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_campus_token_request").
		AddTag("device_id", hashedDeviceId).
		AddField("request_type", requestType).
		AddField("notification_reason_response", reason)

	write.WritePoint(p)

	FlushIfDevelop(write)
}

func LogIOSBackgroundRequestResponse(deviceId, requestType string) {
	write := writeAPI()

	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_campus_token_response").
		AddTag("device_id", hashedDeviceId).
		AddField("request_type", requestType)

	write.WritePoint(p)

	FlushIfDevelop(write)
}

func hashSha256(s string) string {
	h := sha256.New()

	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}

func FlushIfDevelop(write api.WriteAPI) {
	if env.IsDev() {
		write.Flush()
	}
}

func writeAPI() api.WriteAPI {
	return GetClient().WriteAPI(influxOrg, influxBucket)
}
