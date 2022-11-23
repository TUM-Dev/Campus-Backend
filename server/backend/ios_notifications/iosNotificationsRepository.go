package ios_notifications

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (service *Repository) RegisterDevice(device *model.IOSDevice) error {
	if err := service.DB.Create(device).Error; err != nil {
		return err
	}

	return nil
}

func (service *Repository) RemoveDevice(deviceId string) error {
	if err := service.DB.Delete(&model.IOSDevice{DeviceID: deviceId}).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
