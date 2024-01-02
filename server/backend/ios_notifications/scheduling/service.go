// Package scheduling provides functionality for updating user information
// and optionally sending notifications to iOS devices.
package scheduling

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/scheduled_update_log"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

var devicesToUpdate = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "ios_scheduling_devices_to_update",
	Help: "The numer of devices that should be updated for a given priority",
})

type Service struct {
	DevicesRepository      *device.Repository
	SchedulerLogRepository *scheduled_update_log.Repository
	APNs                   *apns.Service
}

func (service *Service) HandleScheduledCron() error {
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

	priority := findSchedulingPriority()
	devices = service.selectDevicesToUpdate(devices, priority)

	log.Infof("Updating %d devices", len(devices))

	service.handleDevices(devices)

	return nil
}

func (service *Service) handleDevices(devices []model.IOSDeviceLastUpdated) {
	routinesCount := min(len(devices), 10)
	chunkSize := len(devices) / routinesCount
	perfectlyChunkable := (len(devices) % routinesCount) == 0

	chunksArrSize := routinesCount
	if !perfectlyChunkable {
		chunksArrSize = routinesCount + 1
	}

	chunks := make([][]model.IOSDeviceLastUpdated, chunksArrSize)

	for i := 0; i < routinesCount; i++ {
		chunks[i] = devices[i*chunkSize : (i+1)*chunkSize]
	}

	if !perfectlyChunkable {
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

func (service *Service) LogScheduledUpdate(deviceID string) error {
	scheduleLog := model.IOSScheduledUpdateLog{
		DeviceID: deviceID,
		Type:     model.IOSUpdateTypeGrades,
	}

	return service.SchedulerLogRepository.LogScheduledUpdate(&scheduleLog)
}

// selectDevicesToUpdate selects the first DevicesToCheckPerCronBase * priority devices
func (service *Service) selectDevicesToUpdate(devices []model.IOSDeviceLastUpdated, priority int) []model.IOSDeviceLastUpdated {
	maxDevicesToCheck := 10 * priority

	if len(devices) < maxDevicesToCheck {
		return devices
	}

	return devices[:maxDevicesToCheck]
}
func findSchedulingPriority() int {
	now := time.Now()

	isNight := 1 <= now.Hour() && now.Hour() <= 6
	if isNight {
		return 1
	}

	isDuringSummerSemester := 32 <= now.YearDay() && now.YearDay() <= 106
	isDuringWinterSemester := 152 <= now.YearDay() && now.YearDay() <= 288
	if isDuringWinterSemester || isDuringSummerSemester {
		return 10
	}

	return 5
}

func NewService(
	devicesRepository *device.Repository,
	schedulerRepository *scheduled_update_log.Repository,
	apnsService *apns.Service,
) *Service {
	return &Service{
		DevicesRepository:      devicesRepository,
		SchedulerLogRepository: schedulerRepository,
		APNs:                   apnsService,
	}
}
