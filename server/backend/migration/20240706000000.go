package migration

import (
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type InitialStudentClubCollection struct {
	ID          string `gorm:"primaryKey;type:varchar(100)"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName sets the insert table name for this struct type
func (n *InitialStudentClubCollection) TableName() string {
	return "student_club_collections"
}

type InitialStudentClub struct {
	gorm.Model
	Name                    string
	Description             null.String
	LinkUrl                 null.String `gorm:"type:varchar(190);unique;"`
	ImageID                 null.Int
	Image                   *File `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageCaption            null.String
	StudentClubCollectionID string                       `gorm:"type:varchar(100)"`
	StudentClubCollection   InitialStudentClubCollection `gorm:"foreignKey:StudentClubCollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName sets the insert table name for this struct type
func (n *InitialStudentClub) TableName() string {
	return "student_clubs"
}

// migrate20240706000000
// Added student club support
func migrate20240706000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240706000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&InitialStudentClub{}, &InitialStudentClubCollection{}); err != nil {
				return err
			}
			if err := SafeEnumAdd(tx, &model.Crontab{}, "type", "scrapeStudentClubs"); err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 60 * 60 * 24, // Every 24h
				Type:     null.StringFrom("scrapeStudentClubs"),
			}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Delete(&model.Crontab{}, "type = 'scrapeStudentClubs'").Error; err != nil {
				return err
			}
			if err := SafeEnumRemove(tx, &model.Crontab{}, "type", "scrapeStudentClubs"); err != nil {
				return err
			}
			if err := tx.Exec("drop table student_club_collections").Error; err != nil {
				return err
			}
			return tx.Exec("drop table student_clubs").Error
		},
	}
}
