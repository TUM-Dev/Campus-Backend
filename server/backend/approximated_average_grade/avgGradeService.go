package approximated_average_grade

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrCouldNotFetchGrades           = status.Error(codes.Internal, "Could not fetch grades from TUMOnline")
	ErrCouldNotFindALectureForAGrade = status.Error(codes.Internal, "Could not find a lecture for a grade")
	ErrBasicInternalError            = status.Error(codes.Internal, "An internal error occurred")
)

type Service struct {
	Repository *Repository
}

func (service *Service) HandleGetApproximatedAverageGrade(req *pb.GetApproximatedAverageGradeRequest) (*pb.GetApproximatedAverageGradeReply, error) {
	grades, err := campus_api.FetchGrades(req.CampusToken)
	if err != nil {
		return nil, ErrCouldNotFetchGrades
	}

	steps := service.Repository.GetCalculationSteps(grades.Grades)

	sum := 0.0
	for _, step := range steps {
		sum += step.Grade * float64(step.Ects) * step.Weight
	}

	average := sum / float64(len(steps))

	return &pb.GetApproximatedAverageGradeReply{
		AverageGrade:     average,
		CalculationSteps: steps,
	}, nil
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		Repository: NewRepository(db),
	}
}
