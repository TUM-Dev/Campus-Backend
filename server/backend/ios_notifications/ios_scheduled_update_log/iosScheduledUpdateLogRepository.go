package ios_scheduled_update_log

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (service *Repository) AddScheduledUpdateLog(log *model.IOSScheduledUpdateLog) error {
	if err := service.DB.Create(log).Error; err != nil {
		return err
	}

	return nil
}

func (service *Repository) GetScheduledUpdateLogForDevice(deviceId string) ([]model.IOSScheduledUpdateLog, error) {
	var logs []model.IOSScheduledUpdateLog

	if err := service.DB.Where("device_id = ?", deviceId).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
