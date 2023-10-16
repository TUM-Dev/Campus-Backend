package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func (m TumDBMigrator) migrate20231016000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231016000000",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&model.EncryptedGrade{}, "is_encrypted")
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().AddColumn(&model.EncryptedGrade{}, "is_encrypted")
		},
	}
}
