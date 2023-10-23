package migration

import (
	_ "embed"
	"encoding/json"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//migrate20221115000000

//go:embed static_data/iosInitialSchedulingPriorities.json
var iosInitialPrioritiesFile []byte

type InitialIOSDevice struct {
	DeviceID          string    `gorm:"primary_key" json:"deviceId"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"createdAt"`
	PublicKey         string    `gorm:"not null" json:"publicKey"`
	ActivityToday     int32     `gorm:"default:0" json:"activityToday"`
	ActivityThisWeek  int32     `gorm:"default:0" json:"activityThisWeek"`
	ActivityThisMonth int32     `gorm:"default:0" json:"activityThisMonth"`
	ActivityThisYear  int32     `gorm:"default:0" json:"activityThisYear"`
}

// TableName sets the insert table name for this struct type
func (n *InitialIOSDevice) TableName() string {
	return "ios_devices"
}

type InitialIOSSchedulingPriority struct {
	ID       int64 `gorm:"primary_key;auto_increment;not_null" json:"id"`
	FromDay  int   `gorm:"not null" json:"from_day"`
	ToDay    int   `gorm:"not null" json:"to_day"`
	FromHour int   `gorm:"not null" json:"from_hour"`
	ToHour   int   `gorm:"not null" json:"to_hour"`
	Priority int   `gorm:"not null" json:"priority"`
}

// TableName sets the insert table name for this struct type
func (p *InitialIOSSchedulingPriority) TableName() string {
	return "ios_scheduling_priorities"
}

// InitialIOSScheduledUpdateLog logs the last time a device was updated.
type InitialIOSScheduledUpdateLog struct {
	ID        int64            `gorm:"primary_key;auto_increment;not_null" json:"id"`
	DeviceID  string           `gorm:"index:idx_scheduled_update_log_device,unique" json:"deviceId"`
	Device    InitialIOSDevice `gorm:"constraint:OnDelete:CASCADE;" json:"device"`
	Type      string           `gorm:"type:enum ('grades');" json:"type"`
	CreatedAt time.Time        `gorm:"index:idx_scheduled_update_log_created,unique;autoCreateTime" json:"createdAt"`
}

// TableName sets the insert table name for this struct type
func (p *InitialIOSScheduledUpdateLog) TableName() string {
	return "ios_scheduled_update_logs"
}

type InitialIOSDeviceRequestLog struct {
	RequestID   string           `gorm:"primary_key;default:UUID()" json:"requestId"`
	DeviceID    string           `gorm:"size:200;not null" json:"deviceId"`
	Device      InitialIOSDevice `gorm:"constraint:OnDelete:CASCADE;" json:"device"`
	RequestType string           `gorm:"not null;type:enum ('CAMPUS_TOKEN_REQUEST');" json:"requestType"`
	CreatedAt   time.Time        `gorm:"autoCreateTime" json:"createdAt"`
}

// TableName sets the insert table name for this struct type
func (p *InitialIOSDeviceRequestLog) TableName() string {
	return "ios_device_request_logs"
}

type InitialIOSEncryptedGrade struct {
	ID           int64            `gorm:"primaryKey"`
	Device       InitialIOSDevice `gorm:"constraint:OnDelete:CASCADE"`
	DeviceID     string           `gorm:"index;not null"`
	LectureTitle string           `gorm:"not null"`
	Grade        string           `gorm:"not null"`
	IsEncrypted  bool             `gorm:"not null,default:true"`
}

// TableName sets the insert table name for this struct type
func (p *InitialIOSEncryptedGrade) TableName() string {
	return "ios_encrypted_grades"
}

type InitialIOSDevicesActivityReset struct {
	Type      string    `gorm:"primary_key;type:enum('day', 'week', 'month', 'year')"`
	LastReset time.Time `gorm:"autoCreateTime"`
}

// TableName sets the insert table name for this struct type
func (p *InitialIOSDevicesActivityReset) TableName() string {
	return "ios_devices_activity_resets"
}

// migrate20221119131300
// - adds the ability to connect to ios devices
func migrate20221119131300() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221119131300",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&InitialIOSDevice{},
				&InitialIOSSchedulingPriority{},
				&InitialIOSScheduledUpdateLog{},
				&InitialIOSDeviceRequestLog{},
				&InitialIOSEncryptedGrade{},
				&InitialIOSDevicesActivityReset{},
			); err != nil {
				return err
			}

			if err := SafeEnumAdd(tx, &model.Crontab{}, "type", "iosNotifications", "iosActivityReset"); err != nil {
				return err
			}

			var priorities []InitialIOSSchedulingPriority

			if err := json.Unmarshal(iosInitialPrioritiesFile, &priorities); err != nil {
				log.WithError(err).Error("could not unmarshal json")
				return err
			}

			if err := tx.Create(&priorities).Error; err != nil {
				log.WithError(err).Error("could not save priority's")
				return err
			}

			err := tx.Create(&model.Crontab{
				Interval: 60,
				Type:     null.StringFrom("iosNotifications"),
			}).Error

			if err != nil {
				log.WithError(err).Error("could not create crontab")
				return err
			}

			return tx.Create(&model.Crontab{
				Type:     null.StringFrom("iosActivityReset"),
				Interval: 86400,
			}).Error
		},

		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&InitialIOSDevice{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&InitialIOSSchedulingPriority{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&InitialIOSScheduledUpdateLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&InitialIOSDeviceRequestLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&InitialIOSEncryptedGrade{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&InitialIOSDevicesActivityReset{}); err != nil {
				return err
			}

			err := tx.Delete(&model.Crontab{}, "type = 'iosNotifications'").Error
			if err != nil {
				return err
			}

			err = tx.Delete(&model.Crontab{}, "type = 'iosActivityReset'").Error

			if err != nil {
				return err
			}

			return SafeEnumRemove(tx, &model.Crontab{}, "type", "iosNotifications", "iosActivityReset")
		},
	}
}
