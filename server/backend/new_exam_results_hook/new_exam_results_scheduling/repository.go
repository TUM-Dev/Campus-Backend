package new_exam_results_scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) StoreExamResultsPublished(examResultsPublished []model.ExamResultPublished) error {
	db := repository.DB

	db.Where("1 = 1").Delete(&model.ExamResultPublished{})

	return db.
		Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
		Create(examResultsPublished).Error
}

func (repository *Repository) FindAllExamResultsPublished() (*[]model.ExamResultPublished, error) {
	db := repository.DB

	var results []model.ExamResultPublished

	err := db.Find(&results).Error

	return &results, err
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
