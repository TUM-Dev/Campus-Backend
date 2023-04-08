package approximated_average_grade

import (
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type Repository struct {
	DB *gorm.DB
}

var customContainsCourses = map[string]int32{
	"Praktikum": 10,
	"Seminar":   5,
}

const FallbackEcts = 5

func (repo *Repository) GetCalculationSteps(grades []model.Grade) []*pb.ApproximatedAverageGradeCalculationStep {
	var calculationSteps []*pb.ApproximatedAverageGradeCalculationStep

	for _, grade := range grades {
		var lecture model.CrawlerLecture

		errorMessage := ""
		detailMessage := ""

		tx := repo.DB.
			Where("lecture_id = ?", grade.LectureNumber).
			Or("name like ?", "%"+grade.LectureTitle+"%").
			Or("name like ?", "%"+grade.LectureTitle+"%").
			First(&lecture)

		lectureFound := tx.Error == nil
		customContainsEcts, isCustomCourse := isCustomContainsCourse(grade.LectureTitle)

		if !lectureFound && !isCustomCourse {
			errorMessage = fmt.Sprintf("Could not find lecture %s. Using default ECTS value of %d", grade.LectureTitle, FallbackEcts)
		}

		floatGrade, err := parseGrade(grade.Grade)
		if err != nil {
			errorMessage = "Could not parse grade " + grade.Grade + " to float"
		}

		if floatGrade > 4.0 {
			// grade is not counted because exam not passed
			continue
		}

		floatWeight := 1.0

		if lectureFound {
			floatWeight, err = parseCommaStringToFloat(lecture.Weight)
			if err != nil {
				errorMessage = "Could not parse lecture weight " + lecture.Weight + " to float"
			}
		}

		ectsValue := int32(FallbackEcts)

		if lectureFound {
			var detailMsg *string

			ectsValue, detailMsg, err = parseEctsStringToInt(lecture.ECTS)
			if err != nil {
				errorMessage = "Could not parse lecture ECTS " + lecture.ECTS + " to int"
			}
			if detailMsg != nil {
				detailMessage = *detailMsg
			}
		}

		if !lectureFound && isCustomCourse {
			ectsValue = customContainsEcts
		}

		calculationStep := &pb.ApproximatedAverageGradeCalculationStep{
			LectureTitle: grade.LectureTitle,
			Grade:        floatGrade,
			Weight:       floatWeight,
			Ects:         ectsValue,
		}

		if len(errorMessage) > 0 {
			calculationStep.ErrorMessage = &errorMessage
		}

		if len(detailMessage) > 0 {
			calculationStep.InfoMessage = &detailMessage
		}

		calculationSteps = append(calculationSteps, calculationStep)
	}

	return calculationSteps
}

func isCustomContainsCourse(lectureTitle string) (int32, bool) {
	for key, value := range customContainsCourses {
		if strings.Contains(strings.ToLower(lectureTitle), strings.ToLower(key)) {
			return value, true
		}
	}

	return 0, false
}

func parseEctsStringToInt(stringEcts string) (int32, *string, error) {
	ectsParts := strings.Split(stringEcts, ",")

	if len(ectsParts) > 1 {
		detailMessage := fmt.Sprintf("Found %d ECTS values for one lecture. Selecting the first one of %s.", len(ectsParts), strings.Join(ectsParts, ", "))
		ectsValue, err := strconv.Atoi(ectsParts[0])

		if err != nil {
			return 0, nil, err
		}

		return int32(ectsValue), &detailMessage, nil
	}

	ectsValue, err := strconv.Atoi(stringEcts)
	if err != nil {
		return 0, nil, err
	}

	return int32(ectsValue), nil, nil
}

func parseGrade(gradeString string) (float64, error) {
	switch gradeString {
	case "N":
		return 5.0, nil
	case "B":
		return 5.0, nil
	}

	return parseCommaStringToFloat(gradeString)
}

func parseCommaStringToFloat(stringFloat string) (float64, error) {
	commaReplaced := strings.Replace(stringFloat, ",", ".", -1)

	floatValue, err := strconv.ParseFloat(commaReplaced, 64)
	if err != nil {
		return 0, err
	}

	return floatValue, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
