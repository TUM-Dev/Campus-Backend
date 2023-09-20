// Package ios_scheduling provides functionality for updating user information
// and optionally sending notifications to iOS devices.
package ios_scheduling

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_scheduled_update_log"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

const (
	DevicesToCheckPerCronBase = 10
	MaxRoutineCount           = 10
)

var devicesToUpdate = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "ios_scheduling_devices_to_update",
	Help: "The numer of devices that should be updated for a given priority",
})

type Service struct {
	Repository             *Repository
	DevicesRepository      *ios_device.Repository
	SchedulerLogRepository *ios_scheduled_update_log.Repository
	Priority               *model.IOSSchedulingPriority
	APNs                   *ios_apns.Service
}

func (service *Service) HandleScheduledCron() error {
	priorities, err := service.Repository.FindSchedulingPriorities()

	if err != nil {
		return err
	}

	currentPriority := findIOSSchedulingPriorityForNow(priorities)

	devices, err := service.DevicesRepository.GetDevicesThatShouldUpdateGrades()

	if err != nil {
		log.WithError(err).Error("can't get devices")
		return err
	}
	devicesToUpdate.Set(float64(len(devices)))

	if len(devices) == 0 {
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
		if err := service.APNs.RequestGradeUpdateForDevice(device.DeviceID); err != nil {
			log.WithError(err).Error("could not RequestGradeUpdateForDevice")
			continue
		}
		if err := service.LogScheduledUpdate(device.DeviceID); err != nil {
			log.WithError(err).WithField("deviceID", device.DeviceID).Error("could not log scheduled update")
		}
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
