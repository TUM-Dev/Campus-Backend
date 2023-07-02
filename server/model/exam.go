package model

import (
	"encoding/xml"
	"strconv"
	"time"
)

const (
	LectureSemesterWinter = "winter"
	LectureSemesterSummer = "summer"
)

type Exams struct {
	XMLName xml.Name `xml:"pruefungen"`
	Exams   []Exam   `xml:"pruefung"`
}

type Exam struct {
	XMLName       xml.Name   `xml:"pruefung"`
	Date          customDate `xml:"datum"`
	ExamID        string     `xml:"pv_term_nr"`
	LectureTitle  string     `xml:"lv_titel"`
	LectureNumber string     `xml:"lv_nummer"`
	LectureSem    string     `xml:"lv_semester"`
	LectureType   string     `xml:"lv_typ"`
}

func (exam *Exam) LectureSemesterYear() (int, error) {
	return strconv.Atoi(exam.LectureSem[0:2])
}

func (exam *Exam) LectureSemesterType() string {
	switch exam.LectureSem[2:] {
	case "S":
		return LectureSemesterSummer
	case "W":
		return LectureSemesterWinter
	}

	return LectureSemesterWinter
}

func (exam *Exam) TakesPlaceInWinterSemester() bool {
	return exam.LectureSemesterType() == LectureSemesterWinter
}

func (exam *Exam) TakesPlaceInSummerSemester() bool {
	return exam.LectureSemesterType() == LectureSemesterSummer
}

func (exam *Exam) ToDbExam() DbExam {
	return DbExam{
		Id:         exam.ExamID,
		Semester:   exam.LectureSemesterType(),
		Title:      exam.LectureTitle,
		Date:       exam.Date.Time,
		LastUpdate: time.Now(),
	}
}

type DbExam struct {
	Id         string    `gorm:"primaryKey"`
	Semester   string    `gorm:"type:enum ('winter', 'summer');"`
	LastUpdate time.Time `gorm:"default:now()"`
	Title      string    `gorm:"not null"`
	Date       time.Time `gorm:"not null"`
}

func (exam *DbExam) TableName() string {
	return "exams"
}
