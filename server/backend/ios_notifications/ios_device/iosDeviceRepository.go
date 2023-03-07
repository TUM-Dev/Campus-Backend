package ios_device

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

// RegisterDevice registers a new ios device or updates an existing one.
// Returns true if the device already existed, false otherwise.
// If the device already existed, the activity counters are updated.
func (repository *Repository) RegisterDevice(device *model.IOSDevice) (bool, error) {
	exists, err := repository.CheckIfDeviceExists(device.DeviceID)

	if err != nil {
		return false, err
	}

	err = repository.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where("device_id = ?", device.DeviceID).
			FirstOrCreate(&device).Error

		if err != nil {
			return err
		}

		newDevice := model.IOSDevice{
			DeviceID:          device.DeviceID,
			PublicKey:         device.PublicKey,
			ActivityToday:     device.ActivityToday + 1,
			ActivityThisWeek:  device.ActivityThisWeek + 1,
			ActivityThisMonth: device.ActivityThisMonth + 1,
			ActivityThisYear:  device.ActivityThisYear + 1,
		}

		return tx.Save(&newDevice).Error
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repository *Repository) CheckIfDeviceExists(deviceId string) (bool, error) {
	var devices []model.IOSDevice
	if err := repository.DB.Limit(1).Find(&devices, "device_id = ?", deviceId).Error; err != nil {
		return false, err
	}

	return len(devices) > 0, nil
}

func (repository *Repository) RemoveDevice(deviceId string) error {
	if err := repository.DB.Delete(&model.IOSDevice{DeviceID: deviceId}).Error; err != nil {
		return err
	}

	return nil
}

func (repository *Repository) GetDevices() ([]model.IOSDevice, error) {
	var devices []model.IOSDevice
	if err := repository.DB.Find(&devices).Error; err != nil {
		return []model.IOSDevice{}, err
	}

	return devices, nil
}

func (repository *Repository) GetDevice(id string) (*model.IOSDevice, error) {
	var device *model.IOSDevice
	if err := repository.DB.First(&device, "device_id = ?", id).Error; err != nil {
		return nil, err
	}

	return device, nil
}

func (repository *Repository) GetMaxAttendedLecturesCount() (int, error) {
	var maxCount int

	tx := repository.DB.Raw("select coalesce(max(lecture_count), 0) from (select count(*) as lecture_count from ios_device_lectures group by device_id) as t;").Scan(&maxCount)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return maxCount, nil
}

// GetReadyDevices returns a list of devices that can be used to update their grades.
// A device that can be updated is either a new device or a device that has not
// been updated in the last model.IOSMinimumUpdateInterval minutes
func (repository *Repository) GetReadyDevices() (*[]model.IOSDevice, error) {
	var devices []model.IOSDevice

	tx := repository.DB.Raw(buildReadyDevicesQuery(), model.IOSMinimumDeviceUpdateInterval).Scan(&devices)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &devices, nil
}

func buildReadyDevicesQuery() string {
	return `
	select d.*
	from ios_devices d
	where not exists(select drl.device_id
					 from ios_device_request_logs drl
					 where drl.device_id = d.device_id
					   and drl.created_at
						 > subdate(now()
							   , interval ? minute)
					 limit 1)
	group by d.device_id;
`
}

func (repository *Repository) ResetDevicesDailyActivity() error {
	return repository.DB.Model(model.IOSDevice{}).Where("activity_today != ?", 0).Update("activity_today", 0).Error
}

func (repository *Repository) ResetDevicesWeeklyActivity() error {
	return repository.DB.Model(model.IOSDevice{}).Where("activity_this_week != ?", 0).Update("activity_this_week", 0).Error
}

func (repository *Repository) ResetDevicesMonthlyActivity() error {
	return repository.DB.Model(model.IOSDevice{}).Where("activity_this_month != ?", 0).Update("activity_this_month", 0).Error
}

func (repository *Repository) ResetDevicesYearlyActivity() error {
	return repository.DB.Model(model.IOSDevice{}).Where("activity_this_year != ?", 0).Update("activity_this_year", 0).Error
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
