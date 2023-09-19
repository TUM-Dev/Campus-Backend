package model

import (
	"encoding/xml"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_crypto"
)

// IOSGrades is a wrapper for a list of grades => XML stuff
type IOSGrades struct {
	XMLName xml.Name   `xml:"rowset"`
	Grades  []IOSGrade `xml:"row"`
}

type IOSGrade struct {
	XMLName         xml.Name   `xml:"row"`
	Date            customDate `xml:"datum"`
	LectureNumber   string     `xml:"lv_nummer"`
	LectureSemester string     `xml:"lv_semester"`
	LectureTitle    string     `xml:"lv_titel"`
	Examiner        string     `xml:"pruefer_nachname"`
	Grade           string     `xml:"uninotenamekurz"`
	ExamType        string     `xml:"exam_typ_name"`
	Modus           string     `xml:"modus"`
	StudyID         string     `xml:"studienidentifikator"`
	StudyName       string     `xml:"studienbezeichnung"`
	ECTS            string     `xml:"lv_credits"`
}

type customDate struct {
	time.Time
}

func (c *customDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return err
	}

	c.Time = t

	return nil
}

func (grade *IOSGrade) CompareToEncrypted(encryptedGrade *IOSEncryptedGrade) bool {
	return grade.LectureTitle == encryptedGrade.LectureTitle && grade.Grade == encryptedGrade.Grade
}

// IOSEncryptedGrade is a grade that can be encrypted.
// Whether it is currently encrypted or not is indicated by the IsEncrypted field.
type IOSEncryptedGrade struct {
	ID           int64     `gorm:"primaryKey"`
	Device       IOSDevice `gorm:"constraint:OnDelete:CASCADE"`
	DeviceID     string    `gorm:"index;not null"`
	LectureTitle string    `gorm:"not null"`
	Grade        string    `gorm:"not null"`
	IsEncrypted  bool      `gorm:"not null,default:true"`
}

func (e *IOSEncryptedGrade) Encrypt(key string) error {
	encryptedTitle, err := ios_crypto.SymmetricEncrypt(e.LectureTitle, key)
	if err != nil {
		return err
	}

	encryptedGrade, err := ios_crypto.SymmetricEncrypt(e.Grade, key)
	if err != nil {
		return err
	}

	e.LectureTitle = encryptedTitle.String()
	e.Grade = encryptedGrade.String()
	e.IsEncrypted = true

	return nil
}

func (e *IOSEncryptedGrade) Decrypt(key string) error {
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
