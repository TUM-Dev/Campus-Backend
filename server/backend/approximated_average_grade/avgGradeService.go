package approximated_average_grade

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
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

	idMap := gradesToIdMap(grades)

	var studyIdGrades []*pb.ApproximatedAverageGradePerStudy

	for id, grades := range *idMap {
		if len(grades) == 0 {
			continue
		}

		steps := service.Repository.GetCalculationSteps(grades)

		gradeEctsSum := 0.0
		ectsSum := int32(0)

		for _, step := range steps {
			gradeEctsSum += step.Grade * float64(step.Ects) * step.Weight
			ectsSum += step.Ects
		}

		average := 0.0

		if ectsSum > 0 {
			average = gradeEctsSum / float64(ectsSum)
		}

		averageGradePerStudy := pb.ApproximatedAverageGradePerStudy{
			StudyId:          id,
			AverageGrade:     average,
			CalculationSteps: steps,
		}

		studyIdGrades = append(studyIdGrades, &averageGradePerStudy)
	}

	return &pb.GetApproximatedAverageGradeReply{
		Studies: studyIdGrades,
	}, nil
}

func gradesToIdMap(grades *model.Grades) *map[string][]model.Grade {
	gradesMap := map[string][]model.Grade{}

	for _, grade := range grades.Grades {
		gradesMap[grade.StudyID] = append(gradesMap[grade.StudyID], grade)
		log.Infof(grade.StudyID)
	}

	return &gradesMap
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		Repository: NewRepository(db),
	}
}
