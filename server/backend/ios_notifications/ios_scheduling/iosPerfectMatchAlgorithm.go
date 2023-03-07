package ios_scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
)

type LectureOverlapped struct {
	LectureId  string
	Overlapped bool
}

type LectureToOverlapped map[string]*LectureOverlapped
type DeviceToOverlappedLectures map[string][]*LectureOverlapped
type CompareFunc[T comparable] func(T) bool

// FindPerfectDevicesMatch returns the minimum set of devices that are able to attend all lectures
func FindPerfectDevicesMatch(lectures *[]model.IOSLecture, devices *[]model.IOSDevice, devicesLectures *[]model.IOSDeviceLecture, maxAttendedLecturesCount int) *[]model.IOSDevice {
	lecturesToOverlapping := lectureToOverlappedLectureMap(lectures)
	devicesLecturesMap := devicesLecturesToMap(devicesLectures, lecturesToOverlapping)
	devicesDetailsMap := devicesToMap(devices)

	overlappingDeviceIds := getOverlappingDevices(
		devicesLecturesMap,
		lectures,
		maxAttendedLecturesCount,
	)

	var overlappingDevices = make([]model.IOSDevice, len(*overlappingDeviceIds))
	for i, deviceId := range *overlappingDeviceIds {
		overlappingDevices[i] = (*devicesDetailsMap)[deviceId]
	}

	return &overlappingDevices
}

// getOverlappingDevices returns the minimum set of devices that are able to attend all lectures
func getOverlappingDevices(
	devicesLectures *DeviceToOverlappedLectures,
	lectures *[]model.IOSLecture,
	maxAttendedLecturesCount int,
) *[]string {
	var overlapped []*LectureOverlapped
	var overlappingStudents []string
	currentMaxAttended := maxAttendedLecturesCount

	for len(overlapped) < len(*lectures) {
		var newMax string
		var overlappingLectures *[]*LectureOverlapped

		newMax, overlappingLectures = findBestNextMatch(
			devicesLectures,
			currentMaxAttended,
		)

		if newMax == "" {
			break
		}

		delete(*devicesLectures, newMax)

		overlappingStudents = append(overlappingStudents, newMax)

		overlapped = append(overlapped, *overlappingLectures...)

		for _, lecture := range *overlappingLectures {
			lecture.Overlapped = true
		}

		overlappingLecturesCount := len(*overlappingLectures)

		if currentMaxAttended > overlappingLecturesCount {
			currentMaxAttended = overlappingLecturesCount
		}
	}

	return &overlappingStudents
}

// findBestNextMatch returns the device with the most lectures that are not overlapped
func findBestNextMatch(
	devicesLectures *DeviceToOverlappedLectures,
	currentMaxAttended int,
) (string, *[]*LectureOverlapped) {
	maxAttends := 0
	studentWithMaxAttends := ""
	var overlappingLectures []*LectureOverlapped

	for device, lectures := range *devicesLectures {
		newAttends := filter(&lectures, func(lecture *LectureOverlapped) bool {
			return !lecture.Overlapped
		})

		newAttendsCount := len(*newAttends)

		if newAttendsCount > maxAttends {
			maxAttends = newAttendsCount
			studentWithMaxAttends = device
			overlappingLectures = *newAttends

			if maxAttends == currentMaxAttended {
				break
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

func lectureToOverlappedLectureMap(l *[]model.IOSLecture) *LectureToOverlapped {
	lectures := make(LectureToOverlapped)

	for _, lecture := range *l {
		overlapped := LectureOverlapped{
			LectureId:  lecture.Id,
			Overlapped: false,
		}

		lectures[lecture.Id] = &overlapped
	}

	return &lectures
}

func devicesLecturesToMap(dl *[]model.IOSDeviceLecture, lectures *LectureToOverlapped) *DeviceToOverlappedLectures {
	devicesLectures := make(DeviceToOverlappedLectures)

	for _, d := range *dl {
		overlappedLecture := (*lectures)[d.LectureId]

		devicesLectures[d.DeviceId] = append(devicesLectures[d.DeviceId], overlappedLecture)
	}

	return &devicesLectures
}

func devicesToMap(d *[]model.IOSDevice) *map[string]model.IOSDevice {
	devices := make(map[string]model.IOSDevice)

	for _, device := range *d {
		devices[device.DeviceID] = device
	}

	return &devices
}
