package device

import (
	"errors"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) CreateDevice(device *model.IOSDevice) (bool, error) {

	exists, err := repository.CheckIfDeviceExists(device.DeviceID)
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, repository.DB.Create(device).Error
}
func (repository *Repository) CheckIfDeviceExists(deviceId string) (bool, error) {
	var devices []model.IOSDevice
	if err := repository.DB.Limit(1).Find(&devices, "device_id = ?", deviceId).Error; err != nil {
		return false, err
	}
	return len(devices) > 0, nil
}

func (repository *Repository) DeleteDevice(deviceId string) error {
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

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
