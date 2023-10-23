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

type IOSDevice struct {
	DeviceID          string    `gorm:"primary_key" json:"deviceId"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"createdAt"`
	PublicKey         string    `gorm:"not null" json:"publicKey"`
	ActivityToday     int32     `gorm:"default:0" json:"activityToday"`
	ActivityThisWeek  int32     `gorm:"default:0" json:"activityThisWeek"`
	ActivityThisMonth int32     `gorm:"default:0" json:"activityThisMonth"`
	ActivityThisYear  int32     `gorm:"default:0" json:"activityThisYear"`
}

type IOSSchedulingPriority struct {
	ID       int64 `gorm:"primary_key;auto_increment;not_null" json:"id"`
	FromDay  int   `gorm:"not null" json:"from_day"`
	ToDay    int   `gorm:"not null" json:"to_day"`
	FromHour int   `gorm:"not null" json:"from_hour"`
	ToHour   int   `gorm:"not null" json:"to_hour"`
	Priority int   `gorm:"not null" json:"priority"`
}

// IOSScheduledUpdateLog logs the last time a device was updated.
type IOSScheduledUpdateLog struct {
	ID        int64     `gorm:"primary_key;auto_increment;not_null" json:"id"`
	DeviceID  string    `gorm:"index:idx_scheduled_update_log_device,unique" json:"deviceId"`
	Device    IOSDevice `gorm:"constraint:OnDelete:CASCADE;" json:"device"`
	Type      string    `gorm:"type:enum ('grades');" json:"type"`
	CreatedAt time.Time `gorm:"index:idx_scheduled_update_log_created,unique;autoCreateTime" json:"createdAt"`
}

type IOSDeviceRequestLog struct {
	RequestID   string    `gorm:"primary_key;default:UUID()" json:"requestId"`
	DeviceID    string    `gorm:"size:200;not null" json:"deviceId"`
	Device      IOSDevice `gorm:"constraint:OnDelete:CASCADE;" json:"device"`
	RequestType string    `gorm:"not null;type:enum ('CAMPUS_TOKEN_REQUEST');" json:"requestType"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type IOSEncryptedGrade struct {
	ID           int64     `gorm:"primaryKey"`
	Device       IOSDevice `gorm:"constraint:OnDelete:CASCADE"`
	DeviceID     string    `gorm:"index;not null"`
	LectureTitle string    `gorm:"not null"`
	Grade        string    `gorm:"not null"`
	IsEncrypted  bool      `gorm:"not null,default:true"`
}

type IOSDevicesActivityReset struct {
	Type      string    `gorm:"primary_key;type:enum('day', 'week', 'month', 'year')"`
	LastReset time.Time `gorm:"autoCreateTime"`
}

// migrate20221119131300
// - adds the ability to connect to ios devices
func migrate20221119131300() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221119131300",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&IOSDevice{},
				&IOSSchedulingPriority{},
				&IOSScheduledUpdateLog{},
				&IOSDeviceRequestLog{},
				&IOSEncryptedGrade{},
				&IOSDevicesActivityReset{},
			); err != nil {
				return err
			}

			if err := SafeEnumAdd(tx, &model.Crontab{}, "type", "iosNotifications", "iosActivityReset"); err != nil {
				return err
			}

			var priorities []IOSSchedulingPriority

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
			if err := tx.Migrator().DropTable(&IOSDevice{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&IOSSchedulingPriority{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&IOSScheduledUpdateLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&IOSDeviceRequestLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&IOSEncryptedGrade{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&IOSDevicesActivityReset{}); err != nil {
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
