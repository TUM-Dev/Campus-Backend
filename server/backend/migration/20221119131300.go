package migration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/model"
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
			); err != nil {
				return err
			}

			if cTypes, err := tx.Migrator().ColumnTypes(&model.Crontab{}); err == nil {
				types, _ := getCrontabEnumTypes(cTypes)

				if !containsIOSNotifications(types) {
					if err := tx.Exec(fmt.Sprintf("alter table test.crontab modify type enum (%s) null;", typesToString(types))).Error; err != nil {
						log.Info(err.Error())
						return err
					}
				}
			}

			if path, err := filepath.Abs("backend/ios_notifications/ios_scheduling/iosInitialSchedulingPriorities.json"); err == nil {
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

			return tx.Create(&model.Crontab{
				Interval: 60,
				Type:     null.String{NullString: sql.NullString{String: cron.IOS_NOTIFICATIONS, Valid: true}},
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

			if cTypes, err := tx.Migrator().ColumnTypes(&model.Crontab{}); err == nil {
				types, _ := getCrontabEnumTypes(cTypes)

				if !containsIOSNotifications(types) {
					if err := tx.Exec(fmt.Sprintf("alter table test.crontab modify type enum (%s) null;", rollbackTypesToString(types))).Error; err != nil {
						log.Info(err.Error())
						return err
					}
				}
			}

			return tx.Delete(&model.Crontab{}, "type = ? AND interval = ?", cron.IOS_NOTIFICATIONS, 60).Error
		},
	}
}

func typesToString(types []string) string {
	var str string
	for _, t := range types {
		str += fmt.Sprintf("%s,", t)
	}

	str += "'iosNotifications'"

	return str
}

func rollbackTypesToString(types []string) string {
	var str string

	for i := 0; i < len(types)-1; i++ {
		str += fmt.Sprintf("%s,", types[i])
	}

	return strings.TrimRight(str, ",")
}

func containsIOSNotifications(types []string) bool {
	for _, t := range types {
		if t == "'iosNotifications'" {
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
