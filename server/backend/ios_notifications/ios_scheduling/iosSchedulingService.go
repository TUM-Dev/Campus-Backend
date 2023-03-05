// Package ios_scheduling provides functionality for updating user information
// and optionally sending notifications to iOS devices.
package ios_scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/influx"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_lectures"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_scheduled_update_log"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	DevicesToCheckPerCronBase = 10
	MaxRoutineCount           = 10
)

type Service struct {
	Repository             *Repository
	DevicesRepository      *ios_device.Repository
	SchedulerLogRepository *ios_scheduled_update_log.Repository
	Priority               *model.IOSSchedulingPriority
	APNs                   *ios_apns.Service
}

func (service *Service) NewHandleScheduledCron() error {
	priorities, err := service.Repository.FindSchedulingPriorities()
	if err != nil {
		log.WithError(err).Error("Error while getting priorities")
		return err
	}

	currentPriority := findIOSSchedulingPriorityForNow(priorities)

	log.Infof("New HandleScheduledCron: %d", currentPriority.Priority)

	lecturesRepo := ios_lectures.NewRepository(service.Repository.DB)

	lecturesToCheck, err := lecturesRepo.GetLecturesToUpdate()
	if err != nil {
		log.WithError(err).Error("Error while getting lectures to check")
		return err
	}

	err = service.checkLectures(lecturesToCheck)
	if err != nil {
		log.WithError(err).Error("Error while checking lectures")
		return err
	}

	return nil
}

func (service *Service) checkLectures(lectures []model.IOSLecture) error {
	log.Infof("Checking %d lectures", len(lectures))

	for _, lecture := range lectures {
		if err := service.updateLecture(lecture); err != nil {
			log.WithError(err).Error("Error while updating lecture")
			continue
		}
	}

	return nil
}

func (service *Service) updateLecture(lecture model.IOSLecture) error {
	log.Infof("Updating lecture %s", lecture.Id)

	devices, err := service.selectDevicesToUpdateLecture(&lecture)
	if err != nil {
		log.WithError(err).Error("Error while selecting devices")
		return err
	}

	log.Infof("Found %d devices to update lecture %s", len(devices), lecture.Id)

	for _, device := range devices {
		err := service.APNs.RequestLectureUpdateForDevice(device.DeviceID)
		if err != nil {
			log.WithError(err).Error("Error while requesting lecture update")
			continue
		}
	}

	return nil
}

func (service *Service) selectDevicesToUpdateLecture(lecture *model.IOSLecture) ([]model.IOSDevice, error) {
	devices, err := service.DevicesRepository.GetDevicesThatCouldUpdateLecture(lecture)
	if err != nil {
		log.WithError(err).Error("Error while getting devices")
		return nil, err
	}

	return devices, nil
}

func (service *Service) HandleScheduledCron() error {
	priorities, err := service.Repository.FindSchedulingPriorities()

	if err != nil {
		return err
	}

	currentPriority := findIOSSchedulingPriorityForNow(priorities)

	devices, err := service.DevicesRepository.GetDevicesThatShouldUpdateGrades()

	if err != nil {
		log.Errorf("Error while getting devices: %s", err)
		return err
	}

	devicesLen := len(devices)

	influx.LogIOSSchedulingDevicesToUpdate(devicesLen, currentPriority.Priority)

	if devicesLen == 0 {
		log.Info("No devices to update")
		return nil
	}

	devices = service.selectDevicesToUpdate(devices, currentPriority.Priority)

	log.Infof("Updating %d devices", len(devices))

	service.handleDevices(devices)

	return nil
}

func (service *Service) handleDevices(devices []model.IOSDeviceLastUpdated) {
	routinesCount := routineCount(devices)

	chunkSize := len(devices) / routinesCount

	perfectChunkable := (len(devices) % routinesCount) == 0

	chunksArrSize := routinesCount

	if !perfectChunkable {
		chunksArrSize = routinesCount + 1
	}

	chunks := make([][]model.IOSDeviceLastUpdated, chunksArrSize)

	for i := 0; i < routinesCount; i++ {
		chunks[i] = devices[i*chunkSize : (i+1)*chunkSize]
	}

	if !perfectChunkable {
		chunks[routinesCount] = devices[routinesCount*chunkSize:]
	}

	var group sync.WaitGroup

	for _, chunk := range chunks {
		group.Add(1)
		go func(chunk []model.IOSDeviceLastUpdated) {
			defer group.Done()
			service.handleDevicesChunk(chunk)
		}(chunk)
	}

	group.Wait()
}

func (service *Service) handleDevicesChunk(devices []model.IOSDeviceLastUpdated) {
	for _, device := range devices {
		err := service.APNs.RequestGradeUpdateForDevice(device.DeviceID)

		if err != nil {
			log.Errorf("Error while handling device: %s", err)
			continue
		}

		service.LogScheduledUpdate(device.DeviceID)
	}
}

func routineCount(devices []model.IOSDeviceLastUpdated) int {
	if len(devices) < MaxRoutineCount {
		return len(devices)
	}

	return MaxRoutineCount
}

func (service *Service) LogScheduledUpdate(deviceID string) error {
	scheduleLog := model.IOSScheduledUpdateLog{
		DeviceID: deviceID,
		Type:     model.IOSUpdateTypeGrades,
	}

	return service.SchedulerLogRepository.LogScheduledUpdate(&scheduleLog)
}

// selectDevicesToUpdate selects max DevicesToCheckPerCronBase devices to update
// based on the priority.
func (service *Service) selectDevicesToUpdate(devices []model.IOSDeviceLastUpdated, priority int) []model.IOSDeviceLastUpdated {
	maxDevicesToCheck := DevicesToCheckPerCronBase * priority

	if len(devices) < maxDevicesToCheck {
		return devices
	}

	return devices[:maxDevicesToCheck]
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
	schedulerRepository *ios_scheduled_update_log.Repository,
	apnsService *ios_apns.Service,
) *Service {
	return &Service{
		Repository:             repository,
		DevicesRepository:      devicesRepository,
		SchedulerLogRepository: schedulerRepository,
		Priority:               model.DefaultIOSSchedulingPriority(),
		APNs:                   apnsService,
	}
}
