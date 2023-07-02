package ios_exams

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Service struct {
	Repository *Repository
}

func (service *Service) SaveRelevantExamsForDevice(exams []model.Exam, deviceId string) error {
	examsRepo := service.Repository

	relevantExams := FilterRelevantExams(exams)

	err := examsRepo.SaveExams(relevantExams)
	if err != nil {
		log.WithError(err).Error("Could not save exams")
		return err
	}

	err = examsRepo.SaveExamsOfDevice(relevantExams, deviceId)
	if err != nil {
		log.WithError(err).Error("Could not save exams of device")
		return err
	}

	return nil
}

func FilterRelevantExams(exams []model.Exam) []model.Exam {
	var relevantExams []model.Exam

	for _, exam := range exams {
		if isSchedulingRelevantExam(exam) {
			relevantExams = append(relevantExams, exam)
		}
	}

	return relevantExams
}

func isSchedulingRelevantExam(exam model.Exam) bool {
	semester, year := FindCurrentSemester()
	semesterId := buildLectureSemesterId(semester, year)

	if semesterId != exam.LectureSem {
		return false
	}

	return true
}

func buildLectureSemesterId(semester string, year int) string {
	yearString := strconv.Itoa(year)

	switch semester {
	case "winter":
		return yearString + "W"
	case "summer":
		return yearString + "S"
	}

	return ""
}

func FindCurrentSemester() (string, int) {
	return FindSemester(time.Now())
}

// FindSemester returns the semester and year for the given `date`
// e.g. semester = "winter" and year = 20
func FindSemester(date time.Time) (string, int) {
	month := date.Month()

	summerSemesterStart := time.Date(date.Year(), time.April, 1, 0, 0, 0, 0, time.UTC)
	winterSemesterStart := time.Date(date.Year(), time.October, 1, 0, 0, 0, 0, time.UTC)

	if month >= time.January && date.Before(summerSemesterStart) {
		return "winter", (date.Year() - 1) % 100
	}

	if date.After(summerSemesterStart) && date.Before(winterSemesterStart) {
		return "summer", date.Year() % 100
	}

	return "winter", date.Year() % 100
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		Repository: NewRepository(db),
	}
}
