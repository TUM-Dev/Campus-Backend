package backend

import (
	"context"
	"errors"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *CampusServer) GetTopNewsAlert(ctx context.Context, _ *emptypb.Empty) (*pb.GetTopNewsAlertReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var res model.NewsAlert
	err := s.db.Joins("Files").First(&res, "NOW() between `from` and `to`").Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "no current active top news")
	} else if err != nil {
		log.WithError(err).Error("could not GetTopNewsAlert")
		return nil, status.Error(codes.Internal, "could not GetTopNewsAlert")
	}

	return &pb.GetTopNewsAlertReply{Alert: &pb.NewsAlert{
		ImageUrl: res.Files.URL.String,
		Link:     res.Link.String,
		Created:  timestamppb.New(res.Created),
		From:     timestamppb.New(res.From),
		To:       timestamppb.New(res.To),
	}}, nil
}

func (s *CampusServer) GetNewsAlert(ctx context.Context, req *pb.GetNewsAlertRequest) (*pb.GetNewsAlertReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	res := model.NewsAlert{NewsAlert: req.Id}
	err := s.db.Joins("Files").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "no news alert")
	} else if err != nil {
		log.WithError(err).Error("could not GetNewsAlert")
		return nil, status.Error(codes.Internal, "could not GetNewsAlert")
	}

	return &pb.GetNewsAlertReply{Alert: &pb.NewsAlert{
		ImageUrl: res.Files.URL.String,
		Link:     res.Link.String,
		Created:  timestamppb.New(res.Created),
		From:     timestamppb.New(res.From),
		To:       timestamppb.New(res.To),
	}}, nil
}

func (s *CampusServer) GetNewsAlerts(ctx context.Context, _ *emptypb.Empty) (*pb.GetNewsAlertsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var res []*model.NewsAlert
	err := s.db.Joins("Files").Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "no news alerts")
	} else if err != nil {
		log.WithError(err).Error("could not GetNewsAlerts")
		return nil, status.Error(codes.Internal, "could not GetNewsAlerts")
	}

	var alerts []*pb.NewsAlert
	for _, alert := range res {
		alerts = append(alerts, &pb.NewsAlert{
			ImageUrl: alert.Files.URL.String,
			Link:     alert.Link.String,
			Created:  timestamppb.New(alert.Created),
			From:     timestamppb.New(alert.From),
			To:       timestamppb.New(alert.To),
		})
	}
	return &pb.GetNewsAlertsReply{Alerts: alerts}, nil
}
