package migration

import (
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// IOSScheduledUpdateLog logs the last time a device was updated.
type IOSScheduledUpdateLog struct {
	ID        int64     `gorm:"primary_key;auto_increment;not_null"`
	DeviceID  string    `gorm:"index:idx_scheduled_update_log_device,unique"`
	Device    IOSDevice `gorm:"constraint:OnDelete:CASCADE;"`
	Type      string    `gorm:"type:enum ('grades');"`
	CreatedAt time.Time `gorm:"index:idx_scheduled_update_log_created,unique;autoCreateTime"`
}

type IOSDevicesActivityReset struct {
	Type      string    `gorm:"primary_key;type:enum('day', 'week', 'month', 'year')"`
	LastReset time.Time `gorm:"autoCreateTime"`
}

// IOSEncryptedGrade is a grade that can be encrypted.
// Whether it is currently encrypted or not is indicated by the IsEncrypted field.
type IOSEncryptedGrade struct {
	ID           int64     `gorm:"primaryKey"`
	Device       IOSDevice `gorm:"constraint:OnDelete:CASCADE"`
	DeviceID     string    `gorm:"index;not null"`
	LectureTitle string    `gorm:"not null"`
	Grade        string    `gorm:"not null"`
	IsEncrypted  bool      `gorm:"not null,default:true"`
}

// An IOSDeviceRequestLog is created when the backend wants to request data from the device.
//
// 1. The backend creates a new IOSDeviceRequestLog
//
// 2. The backend sends a background push notification to the device containing
// the RequestID of the IOSDeviceRequestLog.
//
// 3. The device receives the push notification and sends a request to the backend
// containing the RequestID and the data.
type IOSDeviceRequestLog struct {
	RequestID   string    `gorm:"primary_key;default:UUID()"`
	DeviceID    string    `gorm:"size:200;not null"`
	Device      IOSDevice `gorm:"constraint:OnDelete:CASCADE;"`
	RequestType string    `gorm:"not null;type:enum ('CAMPUS_TOKEN_REQUEST');"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

// IOSDevice stores relevant device information.
// E.g. the PublicKey which is used to encrypt push notifications
// The DeviceID can be used to send push notifications via APNs
type IOSDevice struct {
	DeviceID          string    `gorm:"primary_key"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	PublicKey         string    `gorm:"not null"`
	ActivityToday     int32     `gorm:"default:0"`
	ActivityThisWeek  int32     `gorm:"default:0"`
	ActivityThisMonth int32     `gorm:"default:0"`
	ActivityThisYear  int32     `gorm:"default:0"`
}

// migrate20240207000000
// remove unused ios notifications
func migrate20240207000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240207000000",
		Migrate: func(tx *gorm.DB) error {
			// remove cronjobs
			if err := tx.Delete(&model.Crontab{}, "type = 'newExamResultsHook'").Error; err != nil {
				return err
			}
			if err := tx.Delete(&model.Crontab{}, "type = 'iosNotifications'").Error; err != nil {
				return err
			}
			if err := tx.Delete(&model.Crontab{}, "type = 'iosActivityReset'").Error; err != nil {
				return err
			}
			// drop related tables
			if err := tx.Migrator().DropTable(&IOSDeviceRequestLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&IOSEncryptedGrade{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&IOSScheduledUpdateLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&IOSDevicesActivityReset{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&IOSDevice{}); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// create models
			if err := tx.Migrator().AutoMigrate(&IOSDevice{}); err != nil {
				return err
			}
			if err := tx.Migrator().AutoMigrate(&IOSDeviceRequestLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().AutoMigrate(&IOSEncryptedGrade{}); err != nil {
				return err
			}
			if err := tx.Migrator().AutoMigrate(&IOSScheduledUpdateLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().AutoMigrate(&IOSDevicesActivityReset{}); err != nil {
				return err
			}
			// create cron

			if err := tx.Create(&model.Crontab{
				Interval: 60,
				Type:     null.StringFrom("iosNotifications"),
			}).Error; err != nil {
				return err
			}
			if err := tx.Create(&model.Crontab{
				Interval: 86400,
				Type:     null.StringFrom("iosActivityReset"),
			}).Error; err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 600,
				Type:     null.StringFrom("newExamResultsHook"),
			}).Error
		},
	}
}
