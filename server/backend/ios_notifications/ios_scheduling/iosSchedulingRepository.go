package ios_scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) FindSchedulingPriorities() ([]model.IOSSchedulingPriority, error) {
	db := repository.DB

	var priorities []model.IOSSchedulingPriority

	err := db.Find(&priorities).Error

	return priorities, err
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
