package ios_device

import (
	"errors"
	"github.com/TUM-Dev/Campus-Backend/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) RegisterDevice(device *model.IOSDevice) error {
	if err := repository.DB.Create(device).Error; err != nil {
		errors.Is(err, gorm.ErrEmptySlice)

		return err
	}

	return nil
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

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
