// Package ios_lectures provides functionality for saving and finding lectures.
package ios_lectures

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var (
	// relevant lecture types for scheduling are just lectures that will publish a grade
	// like Lectures (VO), Practicals (PR) and Seminars (SE)
	schedulingRelevantLectureTypes = [3]string{
		"VO", "PR", "SE",
	}
	notSchedulingRelevantLectureOrganisations = [1]string{
		"TUVBSBV",
	}
)

type Service struct {
	Repository *Repository
}

// SaveRelevantLecturesForDevice filters the lectures for relevant ones and saves them in the database.
// For information about what is relevant see isSchedulingRelevantLecture.
func (service *Service) SaveRelevantLecturesForDevice(lectures []model.Lecture, deviceId string) error {
	lecturesRepo := service.Repository

	relevantLectures := FilterRelevantLectures(lectures)

	err := lecturesRepo.SaveLecturesAsIOSLectures(relevantLectures)
	if err != nil {
		log.WithError(err).Error("Could not save lectures")
		return err
	}

	err = lecturesRepo.SaveLecturesOfDevice(relevantLectures, deviceId)
	if err != nil {
		log.WithError(err).Error("Could not save lectures of device")
		return err
	}

	return nil
}

func FilterRelevantLectures(lectures []model.Lecture) []model.Lecture {
	var relevantLectures []model.Lecture

	for _, lecture := range lectures {
		if isSchedulingRelevantLecture(lecture) {
			relevantLectures = append(relevantLectures, lecture)
		}
	}

	return relevantLectures
}

// A Lecture isSchedulingRelevantLecture if the following conditions are met:
// - LectureType is in schedulingRelevantLectureTypes e.g. a lecture (VO), practical (PR) or seminar (SE)
// - Organisation is not in notSchedulingRelevantLectureOrganisations e.g. TUVBSBV => Fachschaft
// - SemesterId is the current semester
func isSchedulingRelevantLecture(lecture model.Lecture) bool {
	correctLectureType := false

	for _, lectureType := range schedulingRelevantLectureTypes {
		if lectureType == lecture.LectureTypeShort {
			correctLectureType = true
			break
		}
	}

	if !correctLectureType {
		return false
	}

	correctOrganisation := true
	for _, organisation := range notSchedulingRelevantLectureOrganisations {
		if organisation == lecture.OrganisationShort {
			correctOrganisation = false
			break
		}
	}

	if !correctOrganisation {
		return false
	}

	semester, year := FindCurrentSemester()
	semesterId := buildLectureSemesterId(semester, year)

	if semesterId != lecture.SemesterId {
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
