package ios_grades

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SaveEncryptedGrade(grade *model.EncryptedGrade) error {
	if err := r.DB.Create(grade).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetIOSEncryptedGrades(deviceId string) ([]model.EncryptedGrade, error) {
	var grades []model.EncryptedGrade
	if err := r.DB.Find(&grades, "device_id = ?", deviceId).Error; err != nil {
		return nil, err
	}

	return grades, nil
}

func (r *Repository) DeleteEncryptedGrades(deviceId string) error {
	if err := r.DB.Delete(&model.EncryptedGrade{}, "device_id = ?", deviceId).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) EncryptAndSaveGrades(grades []model.Grade, deviceId string, encryptionKey string) error {
	var encryptedGrades []model.EncryptedGrade

	for _, grade := range grades {
		encryptedGrade := grade.ToEncryptedGrade(deviceId)
		err := encryptedGrade.Encrypt(encryptionKey)

		if err != nil {
			log.WithError(err).Error("Could not encrypt grade")
			continue
		}

		encryptedGrades = append(encryptedGrades, *encryptedGrade)
	}

	if err := r.DB.Create(&encryptedGrades).Error; err != nil {
		log.WithError(err).Error("Could not save grades")
		return err
	}

	return nil
}

func (r *Repository) GetAndDecryptGrades(deviceId string, encryptionKey string) ([]model.EncryptedGrade, error) {
	var encryptedGrades []model.EncryptedGrade
	if err := r.DB.Find(&encryptedGrades, "device_id = ?", deviceId).Error; err != nil {
		return nil, err
	}

	var grades []model.EncryptedGrade

	for _, encryptedGrade := range encryptedGrades {
		err := encryptedGrade.Decrypt(encryptionKey)

		if err != nil {
			log.WithError(err).Error("Could not decrypt grade")
			continue
		}

		grades = append(grades, encryptedGrade)
	}

	return grades, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
