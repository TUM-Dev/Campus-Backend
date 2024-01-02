package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// IOSSchedulingPriority stores some default priorities for the scheduling of grade updates.
type IOSSchedulingPriority struct {
	ID       int64 `gorm:"primary_key;auto_increment;not_null" json:"id"`
	FromDay  int   `gorm:"not null" json:"from_day"`
	ToDay    int   `gorm:"not null" json:"to_day"`
	FromHour int   `gorm:"not null" json:"from_hour"`
	ToHour   int   `gorm:"not null" json:"to_hour"`
	Priority int   `gorm:"not null" json:"priority"`
}

// migrate20240101000000
// inlined the scheduling priorities to not be configured in the DB
func migrate20240101000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240101000000",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec("DROP table ios_scheduling_priorities").Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().AutoMigrate(&IOSSchedulingPriority{}); err != nil {
				return err
			}
			if err := tx.Create(&IOSSchedulingPriority{
				ID:       1,
				FromDay:  152,
				ToDay:    288,
				FromHour: 0,
				ToHour:   23,
				Priority: 10,
			}).Error; err != nil {
				return err
			}
			if err := tx.Create(&IOSSchedulingPriority{
				ID:       2,
				FromDay:  32,
				ToDay:    106,
				FromHour: 0,
				ToHour:   23,
				Priority: 10,
			}).Error; err != nil {
				return err
			}
			if err := tx.Create(&IOSSchedulingPriority{
				ID:       3,
				FromDay:  1,
				ToDay:    365,
				FromHour: 1,
				ToHour:   6,
				Priority: 5,
			}).Error; err != nil {
				return err
			}
			return nil
		},
	}
}
