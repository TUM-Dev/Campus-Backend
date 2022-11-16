package ios_notifications_service

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"gorm.io/gorm"
)

type IOSNotificationsRepository struct {
	DB *gorm.DB
}

func (service *IOSNotificationsRepository) RegisterDevice(device *model.IOSDevice) error {
	if err := service.DB.Create(device).Error; err != nil {
		return err
	}

	return nil
}

func (service *IOSNotificationsRepository) RemoveDevice(deviceId string) error {
	if err := service.DB.Delete(&model.IOSDevice{DeviceID: deviceId}).Error; err != nil {
		return err
	}

	return nil
}
