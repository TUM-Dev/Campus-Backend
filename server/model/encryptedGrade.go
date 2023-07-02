package model

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_crypto"
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
	IsEncrypted  bool      `gorm:"not null,default:true"`
}

func (e *EncryptedGrade) Encrypt(key string) error {
	encryptedTitle, err := ios_crypto.SymmetricEncrypt(e.LectureTitle, key)
	if err != nil {
		log.WithError(err).Error("Failed to encrypt lecture title")
		return err
	}

	encryptedGrade, err := ios_crypto.SymmetricEncrypt(e.Grade, key)
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
	decryptedTitle, err := ios_crypto.SymmetricDecrypt(ios_crypto.EncryptedString(e.LectureTitle), key)

	if err != nil {
		return err
	}

	decryptedGrade, err := ios_crypto.SymmetricDecrypt(ios_crypto.EncryptedString(e.Grade), key)

	if err != nil {
		return err
	}

	e.LectureTitle = decryptedTitle
	e.Grade = decryptedGrade
	e.IsEncrypted = false

	return nil
}
