package backend

import (
	"context"

	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetCanteenHomometer RPC Endpoint
func (s *CampusServer) GetCanteenHomometer(_ context.Context, input *pb.GetCanteenHomometerRequest) (*pb.GetCanteenHomometerReply, error) {
	data := model.CanteenHomometer{Count: 0, MaxCount: 0, Percent: -1} // Initialize with an empty (not found) value
	err := s.db.Model(&model.CanteenHomometer{}).Where(model.CanteenHomometer{CanteenId: input.CanteenId}).FirstOrInit(&data).Error
	if err != nil {
		log.WithError(err).Error("Error while querying the canteen homometer for: ", input.CanteenId)
	}

	return &pb.GetCanteenHomometerReply{
		Count:     data.Count,
		MaxCount:  data.MaxCount,
		Percent:   data.Percent,
		Timestamp: timestamppb.New(data.Timestamp),
	}, nil
}
