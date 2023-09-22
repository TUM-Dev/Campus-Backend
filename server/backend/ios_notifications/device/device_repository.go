package device

import (
	"errors"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) RegisterDevice(device *model.IOSDevice) error {

	return repository.DB.Transaction(func(tx *gorm.DB) error {

		var foundDevice model.IOSDevice

		res := tx.First(&foundDevice, "device_id = ?", device.DeviceID)

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return tx.Create(device).Error
		}

		if res.Error != nil {
			return res.Error
		}

		newDevice := model.IOSDevice{
			DeviceID:          device.DeviceID,
			PublicKey:         device.PublicKey,
			ActivityToday:     foundDevice.ActivityToday + 1,
			ActivityThisWeek:  foundDevice.ActivityThisWeek + 1,
			ActivityThisMonth: foundDevice.ActivityThisMonth + 1,
			ActivityThisYear:  foundDevice.ActivityThisYear + 1,
		}

		return tx.Save(&newDevice).Error
	})
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

// GetDevicesThatShouldUpdateGrades returns a list of devices that should be updated
// A device that needs to be updated is either a new device or a device that has not
// been updated in the last model.IOSMinimumUpdateInterval minutes
func (repository *Repository) GetDevicesThatShouldUpdateGrades() ([]model.IOSDeviceLastUpdated, error) {
	var devices []model.IOSDeviceLastUpdated

	tx := repository.DB.Raw(
		buildDevicesThatShouldUpdateGradesQuery(),
		model.IOSUpdateTypeGrades,
		model.IOSMinimumUpdateInterval,
	).Scan(&devices)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return devices, nil
}

func buildDevicesThatShouldUpdateGradesQuery() string {
	return `
		select d.device_id, ul.created_at as last_updated, d.public_key
		from ios_devices d
				 left join ios_scheduled_update_logs ul on d.device_id = ul.device_id
		where ul.created_at is null
		   or (ul.type = ?
			and ul.created_at < date_sub(now(), interval ? minute))
		group by d.device_id, ul.created_at
		order by ul.created_at;
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
