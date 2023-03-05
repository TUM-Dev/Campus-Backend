package model

import (
	"encoding/xml"
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

func (grade *IOSGrade) CompareToEncrypted(encryptedGrade *IOSEncryptedGrade) bool {
	return grade.LectureTitle == encryptedGrade.LectureTitle && grade.Grade == encryptedGrade.Grade
}
