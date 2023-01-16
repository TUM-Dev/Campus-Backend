package ios_scheduled_update_log

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

func (service *Repository) LogScheduledUpdate(log *model.IOSScheduledUpdateLog) error {
	tx := service.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"created_at",
		}),
	}).Create(log)

	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (service *Repository) GetScheduledUpdateLogForDevice(deviceId string) ([]model.IOSScheduledUpdateLog, error) {
	var logs []model.IOSScheduledUpdateLog

	if err := service.DB.Where("device_id = ?", deviceId).Order("created_at").Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
