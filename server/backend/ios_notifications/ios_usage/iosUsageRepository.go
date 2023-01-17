package ios_usage

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (service *Repository) AddUsage(usage *model.IOSDeviceUsageLog) (*model.IOSDeviceUsageLog, error) {
	if err := service.DB.Create(&usage).Error; err != nil {
		return nil, err
	}

	return usage, nil
}

func (service *Repository) GetUsage(deviceID string) ([]model.IOSDeviceUsageLog, error) {
	var usage []model.IOSDeviceUsageLog

	if err := service.DB.Where("device_id = ?", deviceID).Find(&usage).Error; err != nil {
		return nil, err
	}

	return usage, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
