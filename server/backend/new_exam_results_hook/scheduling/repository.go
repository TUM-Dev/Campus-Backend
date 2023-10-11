package scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) StoreExamResultsPublished(examResultsPublished []model.PublishedExamResult) error {
	db := repository.DB

	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("1 = 1").Delete(&model.PublishedExamResult{}).Error

		if err != nil {
			return err
		}

		// disabled logging because this query always prints a warning because it takes longer then normal
		// to execute because we bulk insert a lot of data
		return tx.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Create(examResultsPublished).Error
	})
}

func (repository *Repository) FindAllExamResultsPublished() (*[]model.PublishedExamResult, error) {
	var results []model.PublishedExamResult
	err := repository.DB.Find(&results).Error

	return &results, err
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
