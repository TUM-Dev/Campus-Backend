// Package ios_scheduling provides functionality for updating user information
// and optionally sending notifications to iOS devices.
package ios_scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_lectures"
	"github.com/TUM-Dev/Campus-Backend/server/backend/utils"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

const (
	MaxRoutineCount = 10
)

type Service struct {
	Repository        *Repository
	DevicesRepository *ios_device.Repository
	Priority          *model.IOSSchedulingPriority
	APNs              *ios_apns.Service
}

// HandleScheduledCron will be called by the cron job every minute. It will calculate a perfect set
// of devices that together attend all lectures. Next, it will request a lecture update for these devices.
func (service *Service) HandleScheduledCron() error {
	priorities, err := service.Repository.FindSchedulingPriorities()
	if err != nil {
		log.WithError(err).Error("Error while getting priorities")
		return err
	}

	currentPriority := findIOSSchedulingPriorityForNow(priorities)

	// TODO: implement currentPriority
	log.Infof("Current priority: %s", currentPriority)

	lecturesRepo := ios_lectures.NewRepository(service.Repository.DB)

	lectures, err := lecturesRepo.GetLectures()
	if err != nil {
		log.WithError(err).Error("Error while getting lectures to check")
		return err
	}

	devicesRepo := ios_device.NewRepository(service.Repository.DB)

	maxAttendedLecturesCount, err := devicesRepo.GetMaxAttendedLecturesCount()
	if err != nil {
		log.WithError(err).Error("Error while getting max attended lectures count")
		return err
	}

	deviceLectures, err := lecturesRepo.GetDeviceLectures()
	if err != nil {
		log.WithError(err).Error("Error while getting device lectures")
		return err
	}

	devices, err := devicesRepo.GetReadyDevices()
	if err != nil {
		log.WithError(err).Error("Error while getting devices")
		return err
	}

	if len(*devices) == 0 {
		log.Info("No devices to check")
		return nil
	}

	devicesMatch := FindPerfectDevicesMatch(lectures, devices, deviceLectures, maxAttendedLecturesCount)

	log.Infof("Found %d perfect matches", len(*devicesMatch))

	service.requestUpdateForDevices(devicesMatch)

	return nil
}

func (service *Service) requestUpdateForDevices(devices *[]model.IOSDevice) {
	routineCount := routineCount(devices)

	utils.RunTasksInRoutines(devices, func(device model.IOSDevice) {
		err := service.APNs.RequestLectureUpdateForDevice(device.DeviceID)
		if err != nil {
			log.WithError(err).Error("Error while requesting grades update")
		}
	}, routineCount)
}

func routineCount(devices *[]model.IOSDevice) int {
	if len(*devices) < MaxRoutineCount {
		return 1
	}

	return MaxRoutineCount
}

func findIOSSchedulingPriorityForNow(priorities []model.IOSSchedulingPriority) *model.IOSSchedulingPriority {
	var prioritiesThatAreInRange []model.IOSSchedulingPriority

	for _, priority := range priorities {
		if priority.IsCurrentlyInRange() {
			prioritiesThatAreInRange = append(prioritiesThatAreInRange, priority)
		}
	}

	if len(prioritiesThatAreInRange) == 0 {
		return model.DefaultIOSSchedulingPriority()
	}

	return mergeIOSSchedulingPriorities(prioritiesThatAreInRange)
}

func mergeIOSSchedulingPriorities(priorities []model.IOSSchedulingPriority) *model.IOSSchedulingPriority {
	mergedPriority := model.DefaultIOSSchedulingPriority()
	prioritiesSum := 0

	for _, priority := range priorities {
		if priority.IsMorePreciseThan(mergedPriority) {
			mergedPriority = &priority
		}

		prioritiesSum += priority.Priority
	}

	mergedPriority.Priority = prioritiesSum / len(priorities)

	return mergedPriority
}

func NewService(repository *Repository,
	devicesRepository *ios_device.Repository,
	apnsService *ios_apns.Service,
) *Service {
	return &Service{
		Repository:        repository,
		DevicesRepository: devicesRepository,
		Priority:          model.DefaultIOSSchedulingPriority(),
		APNs:              apnsService,
	}
}
