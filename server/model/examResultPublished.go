package model

import (
	"encoding/xml"
	"time"
)

type campusApiBool bool

type TUMAPIExamResultsPublished struct {
	XMLName     xml.Name                    `xml:"pruefungen"`
	ExamResults []TUMAPIExamResultPublished `xml:"pruefung"`
}

type TUMAPIExamResultPublished struct {
	Exam
	Published campusApiBool `xml:"note_veroeffentlicht"`
}

func (examResult *TUMAPIExamResultPublished) ToDBExamResult() *ExamResultPublished {
	return &ExamResultPublished{
		Date:         examResult.Date.Time,
		ExamID:       examResult.ExamID,
		LectureTitle: examResult.LectureTitle,
		LectureType:  examResult.LectureType,
		LectureSem:   examResult.LectureSem,
		Published:    bool(examResult.Published),
	}
}

type ExamResultPublished struct {
	Date         time.Time `json:"date"`
	ExamID       string    `gorm:"primary_key" json:"examId"`
	LectureTitle string    `json:"lectureTitle"`
	LectureType  string    `json:"lectureType"`
	LectureSem   string    `json:"lectureSem"`
	Published    bool      `json:"published"`
}

func (ExamResultPublished) TableName() string {
	return "exam_results_published"
}

func (p *campusApiBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	if err := d.DecodeElement(&value, &start); err != nil {
		return err
	}
	switch value {
	case "J":
		*p = true
	default:
		*p = false
	}
	return nil
}

type NewExamsPublishedHookPayload struct {
	PublishedExams []ExamResultPublished `json:"publishedExams"`
}
