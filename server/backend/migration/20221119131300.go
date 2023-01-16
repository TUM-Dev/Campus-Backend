package migration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//migrate20221115000000

func (m TumDBMigrator) migrate20221119131300() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221119131300",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.IOSDevice{},
				&model.IOSDeviceUsageLog{},
				&model.Crontab{},
				&model.IOSSchedulingPriority{},
				&model.IOSScheduledUpdateLog{},
				&model.IOSDeviceRequestLog{},
				&model.IOSEncryptedGrade{},
				&model.IOSLog{},
			); err != nil {
				return err
			}

			if cTypes, err := tx.Migrator().ColumnTypes(&model.Crontab{}); err == nil {
				types, _ := getCrontabEnumTypes(cTypes)

				if !contains(types, "iosNotifications") {
					if err := tx.Exec(fmt.Sprintf("alter table campus_db.crontab modify type enum (%s) null;", typesToString(types, "iosNotifications"))).Error; err != nil {
						log.Info(err.Error())
						return err
					}
				}

				if !contains(types, "iosActivityReset") {
					if err := tx.Exec(fmt.Sprintf("alter table campus_db.crontab modify type enum (%s) null;", typesToString(types, "iosActivityReset"))).Error; err != nil {
						log.Info(err.Error())
						return err
					}
				}
			}

			if path, err := filepath.Abs("/static_data/iosInitialSchedulingPriorities.json"); err == nil {
				file, err := os.Open(path)
				defer file.Close()

				if err != nil {
					log.Info(err.Error())
					return err
				}

				byteValue, err := io.ReadAll(file)

				if err != nil {
					log.Info(err.Error())
					return err
				}

				var priorities []model.IOSSchedulingPriority

				unmarshalErr := json.Unmarshal(byteValue, &priorities)

				if err != nil {
					log.Info(unmarshalErr.Error())
					return unmarshalErr
				}

				if err := tx.Create(&priorities).Error; err != nil {
					log.Info(err.Error())
					return err
				}
			}

			err := tx.Create(&model.Crontab{
				Interval: 60,
				Type:     null.String{NullString: sql.NullString{String: cron.IOS_NOTIFICATIONS, Valid: true}},
			}).Error

			if err != nil {
				log.Error(err.Error())
				return err
			}

			return tx.Create(&model.Crontab{
				Type: null.String{
					NullString: sql.NullString{
						String: cron.IOS_ACTIVITY_RESET,
						Valid:  true,
					},
				},
				Interval: 86400,
			}).Error
		},

		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&model.IOSDevice{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSDeviceUsageLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSSchedulingPriority{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSScheduledUpdateLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSDeviceRequestLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSEncryptedGrade{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSLog{}); err != nil {
				return err
			}

			if cTypes, err := tx.Migrator().ColumnTypes(&model.Crontab{}); err == nil {
				types, _ := getCrontabEnumTypes(cTypes)

				if !contains(types, "iosNotifications") {
					if err := tx.Exec(fmt.Sprintf("alter table campus_db.crontab modify type enum (%s) null;", rollback(types, "iosNotifications"))).Error; err != nil {
						log.Info(err.Error())
						return err
					}
				}

				if !contains(types, "iosActivityReset") {
					if err := tx.Exec(fmt.Sprintf("alter table campus_db.crontab modify type enum (%s) null;", rollback(types, "iosActivityReset"))).Error; err != nil {
						log.Info(err.Error())
						return err
					}
				}
			}

			err := tx.Delete(&model.Crontab{}, "type = ? AND interval = ?", cron.IOS_NOTIFICATIONS, 60).Error
			if err != nil {
				log.Error(err.Error())
			}

			return tx.Delete(&model.Crontab{}, "type = ? AND interval = ?", cron.IOS_ACTIVITY_RESET, 86400).Error
		},
	}
}

func typesToString(types []string, plus string) string {
	var str string
	for _, t := range types {
		str += fmt.Sprintf("%s,", t)
	}

	str += fmt.Sprintf("'%s'", plus)

	return str
}

func rollback(types []string, search string) string {
	var str string

	for _, eType := range types {
		if eType != fmt.Sprintf("'%s'", search) {
			str += fmt.Sprintf("%s,", eType)
		}
	}

	return strings.TrimRight(str, ",")
}

func contains(types []string, search string) bool {
	for _, t := range types {
		if t == fmt.Sprintf("'%s'", search) {
			return true
		}
	}

	return false
}

func getCrontabEnumTypes(types []gorm.ColumnType) ([]string, error) {
	for _, t := range types {
		if t.Name() == "type" {
			if cType, ok := t.ColumnType(); ok {
				leftTrimmed := strings.TrimLeft(cType, "enum")
				trimmed := strings.Trim(leftTrimmed, "()")

				return strings.Split(trimmed, ","), nil
			} else {
				return []string{}, errors.New("could not get column type")
			}
		}
	}

	return []string{}, errors.New("could not find column type")
}
