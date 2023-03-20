package backend

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/approximated_average_grade"
)

func (s *CampusServer) GetApproximatedAverageGradeService() *approximated_average_grade.Service {
	return approximated_average_grade.NewService(s.db)
}

func (s *CampusServer) GetApproximatedAverageGrade(_ context.Context, req *pb.GetApproximatedAverageGradeRequest) (*pb.GetApproximatedAverageGradeReply, error) {
	service := s.GetApproximatedAverageGradeService()
	return service.HandleGetApproximatedAverageGrade(req)
}
