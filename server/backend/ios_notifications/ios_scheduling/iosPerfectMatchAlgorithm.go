package ios_scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
)

type DeviceLectureAttend map[string][]string
type CompareFunc[T comparable] func(T) bool

func FindPerfectDevicesMatch(lectures *[]model.IOSLecture, devices *[]model.IOSDeviceWithAvgResponseTime, devicesLectures *[]model.IOSDeviceLecture, maxAttendedLecturesCount int) *[]model.IOSDeviceWithAvgResponseTime {
	devicesLecturesMap := devicesLecturesToMap(devicesLectures)
	devicesDetailsMap := devicesToMap(devices)

	overlappingDeviceIds := getOverlappingDevices(
		devicesLecturesMap,
		devicesDetailsMap,
		lectures,
		maxAttendedLecturesCount,
		false,
	)

	var overlappingDevices = make([]model.IOSDeviceWithAvgResponseTime, len(*overlappingDeviceIds))
	for i, deviceId := range *overlappingDeviceIds {
		overlappingDevices[i] = (*devicesDetailsMap)[deviceId]
	}

	return &overlappingDevices
}

func getOverlappingDevices(
	devicesLectures *DeviceLectureAttend,
	devicesDetails *map[string]model.IOSDeviceWithAvgResponseTime,
	lectures *[]model.IOSLecture,
	maxAttendedLecturesCount int,
	useAverageResponseTime bool,
) *[]string {
	var overlapped []string
	var overlappingDevices []string
	currentMaxAttended := maxAttendedLecturesCount

	for len(overlapped) < len(*lectures) {
		newMax, overlappingLectures := findBestNextMatch(
			devicesLectures,
			devicesDetails,
			&overlapped,
			currentMaxAttended,
			useAverageResponseTime,
		)

		if newMax == "" {
			break
		}

		delete(*devicesLectures, newMax)

		overlappingDevices = append(overlappingDevices, newMax)

		overlapped = append(overlapped, *overlappingLectures...)

		overlappingLecturesCount := len(*overlappingLectures)

		if currentMaxAttended > overlappingLecturesCount {
			currentMaxAttended = overlappingLecturesCount
		}
	}

	return &overlappingDevices
}

func findBestNextMatch(
	devicesLectures *DeviceLectureAttend,
	devicesDetails *map[string]model.IOSDeviceWithAvgResponseTime,
	overlapped *[]string,
	currentMaxAttended int,
	useAverageResponseTime bool,
) (string, *[]string) {
	maxAttends := 0
	studentWithMaxAttends := ""
	var overlappingLectures []string
	var responseTime float64

	for student, lectures := range *devicesLectures {
		newAttends := filter(&lectures, func(lecture string) bool {
			return !contains(overlapped, lecture)
		})

		newAttendsCount := len(*newAttends)

		if useAverageResponseTime {
			if newAttendsCount > maxAttends || (newAttendsCount == maxAttends && (*devicesDetails)[student].AvgResponseTime < responseTime) {
				maxAttends = len(*newAttends)
				studentWithMaxAttends = student
				overlappingLectures = *newAttends
				responseTime = (*devicesDetails)[student].AvgResponseTime
			}
		} else {
			if newAttendsCount > maxAttends {
				maxAttends = len(*newAttends)
				studentWithMaxAttends = student
				overlappingLectures = *newAttends
				responseTime = (*devicesDetails)[student].AvgResponseTime

				if maxAttends == currentMaxAttended {
					break
				}
			}
		}
	}

	return studentWithMaxAttends, &overlappingLectures
}

func filter[T comparable](s *[]T, fn CompareFunc[T]) *[]T {
	var p []T
	for _, v := range *s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return &p
}

func contains[T comparable](s *[]T, e T) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}

func devicesLecturesToMap(dl *[]model.IOSDeviceLecture) *DeviceLectureAttend {
	devicesLectures := make(DeviceLectureAttend)

	for _, d := range *dl {
		devicesLectures[d.DeviceId] = append(devicesLectures[d.DeviceId], d.LectureId)
	}

	return &devicesLectures
}

func devicesToMap(d *[]model.IOSDeviceWithAvgResponseTime) *map[string]model.IOSDeviceWithAvgResponseTime {
	devices := make(map[string]model.IOSDeviceWithAvgResponseTime)

	for _, device := range *d {
		devices[device.DeviceID] = device
	}

	return &devices
}
