package ios_notifications_service

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"gorm.io/gorm"
)

type IOSUsageRepository struct {
	DB *gorm.DB
}

func (service *IOSUsageRepository) AddUsage(usage *model.IOSDeviceUsageLog) (*model.IOSDeviceUsageLog, error) {
	if err := service.DB.Create(&usage).Error; err != nil {
		return nil, err
	}

	return usage, nil
}
