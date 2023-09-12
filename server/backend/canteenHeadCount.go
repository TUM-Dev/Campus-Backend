package backend

import (
	"context"
	"errors"

	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// GetCanteenHeadCount RPC Endpoint
func (s *CampusServer) GetCanteenHeadCount(_ context.Context, input *pb.GetCanteenHeadCountRequest) (*pb.GetCanteenHeadCountReply, error) {
	data := model.CanteenHeadCount{Count: 0, MaxCount: 0, Percent: -1} // Initialize with an empty (not found) value
	err := s.db.Model(&model.CanteenHeadCount{}).Where(model.CanteenHeadCount{CanteenId: input.CanteenId}).FirstOrInit(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error("while querying the canteen head count for: ", input.CanteenId)
		return nil, errors.New("failed to query head count")
	}

	return &pb.GetCanteenHeadCountReply{
		Count:     data.Count,
		MaxCount:  data.MaxCount,
		Percent:   data.Percent,
		Timestamp: timestamppb.New(data.Timestamp),
	}, nil
}
