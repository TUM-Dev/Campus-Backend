// Package scheduled_update_log provides functionality for logging scheduler updates
// E.g. when a device updated its grades, the scheduler will log the update
package scheduled_update_log

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

func (service *Repository) LogScheduledUpdate(log *model.IOSScheduledUpdateLog) error {
	return service.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"created_at",
		}),
	}).Create(log).Error
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
