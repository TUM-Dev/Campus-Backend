package influx

import influxdb2 "github.com/influxdata/influxdb-client-go/v2"

func LogIOSRegisterDevice(deviceId string) {
	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_register_device").
		AddTag("device_id", hashedDeviceId)

	LogPoint(p)
}

func LogIOSRemoveDevice(deviceId string) {
	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_remove_device").
		AddTag("device_id", hashedDeviceId)

	LogPoint(p)
}

func LogIOSNewGrades(deviceId string, gradesCount int) {
	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_new_grades").
		AddTag("device_id", hashedDeviceId).
		AddField("new_grades_count", gradesCount)

	LogPoint(p)
}

func LogIOSBackgroundRequest(deviceId, requestType, reason string) {
	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_campus_token_request").
		AddTag("device_id", hashedDeviceId).
		AddField("request_type", requestType).
		AddField("notification_reason_response", reason)

	LogPoint(p)
}

func LogIOSBackgroundRequestResponse(deviceId, requestType string) {
	hashedDeviceId := hashSha256(deviceId)

	p := influxdb2.NewPointWithMeasurement("ios_campus_token_response").
		AddTag("device_id", hashedDeviceId).
		AddField("request_type", requestType)

	LogPoint(p)
}
