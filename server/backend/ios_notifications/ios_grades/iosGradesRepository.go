package ios_grades

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SaveEncryptedGrade(grade *model.IOSEncryptedGrade) error {
	if err := r.DB.Create(grade).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetIOSEncryptedGrades(deviceId string) ([]model.IOSEncryptedGrade, error) {
	var grades []model.IOSEncryptedGrade
	if err := r.DB.Find(&grades, "device_id = ?", deviceId).Error; err != nil {
		return nil, err
	}

	return grades, nil
}

func (r *Repository) DeleteEncryptedGrades(deviceId string) error {
	if err := r.DB.Delete(&model.IOSEncryptedGrade{}, "device_id = ?", deviceId).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
