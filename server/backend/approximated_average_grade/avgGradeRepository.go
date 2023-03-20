package approximated_average_grade

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) GetCalculationSteps(grades []model.Grade) []*pb.ApproximatedAverageGradeCalculationStep {
	var calculationSteps []*pb.ApproximatedAverageGradeCalculationStep

	for _, grade := range grades {
		var lecture model.CrawlerLecture

		tx := repo.DB.First(&lecture, "name like ?", "%"+grade.LectureTitle+"%")
		if tx.Error != nil {
			errorMessage := "Could not find a lecture for this grade. Won't be included in the calculation."
			calculationSteps = append(calculationSteps, &pb.ApproximatedAverageGradeCalculationStep{
				LectureTitle: grade.LectureTitle,
				Grade:        1.0,
				Weight:       0,
				Ects:         0,
				ErrorMessage: &errorMessage,
			})
		} else {
			calculationSteps = append(calculationSteps, &pb.ApproximatedAverageGradeCalculationStep{
				LectureTitle: lecture.Name,
				Grade:        1.0,
				Weight:       1,
				Ects:         5,
			})
		}
	}

	return calculationSteps
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
