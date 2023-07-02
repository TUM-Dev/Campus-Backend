package migration

import (
	_ "embed"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (m TumDBMigrator) migrate20230618000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230618000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.DbExam{},
				&model.DeviceExam{},
				&model.EncryptedGrade{},
			); err != nil {
				return err
			}

			err := tx.Migrator().DropColumn(&model.IOSDevice{}, "activity_today")
			err = tx.Migrator().DropColumn(&model.IOSDevice{}, "activity_this_week")
			err = tx.Migrator().DropColumn(&model.IOSDevice{}, "activity_this_month")
			err = tx.Migrator().DropColumn(&model.IOSDevice{}, "activity_this_year")

			if err != nil {
				log.WithError(err).Info("Could not drop columns")
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&model.DbExam{}); err != nil {
				return err
			}

			if err := tx.Migrator().DropTable(&model.DeviceExam{}); err != nil {
				return err
			}

			err := tx.Migrator().AddColumn(&model.IOSDevice{}, "activity_today")
			err = tx.Migrator().AddColumn(&model.IOSDevice{}, "activity_this_week")
			err = tx.Migrator().AddColumn(&model.IOSDevice{}, "activity_this_month")
			err = tx.Migrator().AddColumn(&model.IOSDevice{}, "activity_this_year")

			if err != nil {
				return err
			}

			return nil
		},
	}
}
