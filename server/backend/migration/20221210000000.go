package migration

import (
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// InitialCanteenHeadCount stores all available people counts for available canteens. The CanteenId represents the same ID, as for the canteen inside the eat-api.
type InitialCanteenHeadCount struct {
	CanteenId string    `gorm:"primary_key;column:canteen_id;type:varchar(64);not null;" json:"canteen_id"`
	Count     uint32    `gorm:"column:count;type:int;not null;" json:"count"`
	MaxCount  uint32    `gorm:"column:max_count;type:int;not null;" json:"max_count"`
	Percent   float32   `gorm:"column:percent;type:float;not null;" json:"percent"`
	Timestamp time.Time `gorm:"column:timestamp;type:timestamp;not null;default:current_timestamp();OnUpdate:current_timestamp();" json:"timestamp" `
}

// TableName sets the insert table name for this struct type
func (n *InitialCanteenHeadCount) TableName() string {
	return "canteen_head_count"
}

// migrate20221210000000
// - adds InitialCanteenHeadCount table
// - adds a "canteenHeadCount" cron job that runs every 5 minutes.
func migrate20221210000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221210000000",
		Migrate: func(tx *gorm.DB) error {
			// table
			if err := tx.AutoMigrate(&InitialCanteenHeadCount{}); err != nil {
				return err
			}

			// cron
			if err := SafeEnumAdd(tx, model.Crontab{}, "type", "canteenHeadCount"); err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 60 * 5, // Every 5 minutes
				Type:     null.StringFrom("canteenHeadCount"),
			}).Error
		},

		Rollback: func(tx *gorm.DB) error {
			// table
			if err := tx.Migrator().DropTable(&InitialCanteenHeadCount{}); err != nil {
				return err
			}

			// cron
			if err := tx.Delete(&model.Crontab{}, "type = 'canteenHeadCount'").Error; err != nil {
				return err
			}
			return SafeEnumRemove(tx, model.Crontab{}, "type", "canteenHeadCount")
		},
	}
}
