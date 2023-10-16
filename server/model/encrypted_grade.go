package model

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/crypto"
	log "github.com/sirupsen/logrus"
)

// EncryptedGrade is a grade that can be encrypted.
// Whether it is currently encrypted or not is indicated by the IsEncrypted field.
type EncryptedGrade struct {
	ID           uint      `gorm:"primaryKey"`
	Device       IOSDevice `gorm:"constraint:OnDelete:CASCADE"`
	DeviceID     string    `gorm:"index;not null"`
	LectureTitle string    `gorm:"not null"`
	Grade        string    `gorm:"not null"`
	IsEncrypted  bool      `gorm:"-"`
}

func (e *EncryptedGrade) Encrypt(key string) error {
	encryptedTitle, err := crypto.SymmetricEncrypt(e.LectureTitle, key)
	if err != nil {
		log.WithError(err).Error("Failed to encrypt lecture title")
		return err
	}

	encryptedGrade, err := crypto.SymmetricEncrypt(e.Grade, key)
	if err != nil {
		log.WithError(err).Error("Failed to encrypt grade")
		return err
	}

	e.LectureTitle = encryptedTitle.String()
	e.Grade = encryptedGrade.String()
	e.IsEncrypted = true

	return nil
}

func (e *EncryptedGrade) Decrypt(key string) error {
	decryptedTitle, err := crypto.SymmetricDecrypt(crypto.EncryptedString(e.LectureTitle), key)
	if err != nil {
		return err
	}

	decryptedGrade, err := crypto.SymmetricDecrypt(crypto.EncryptedString(e.Grade), key)
	if err != nil {
		log.WithError(err).Error("Failed to decrypt grade")
		return err
	}

	e.LectureTitle = decryptedTitle
	e.Grade = decryptedGrade
	e.IsEncrypted = false

	return nil
}
