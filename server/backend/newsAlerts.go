package backend

import (
	"context"
	"errors"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *CampusServer) GetTopNews(ctx context.Context, _ *pb.GetTopNewsRequest) (*pb.GetTopNewsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var res *model.NewsAlert
	err := s.db.Joins("Files").Where("NOW() between `from` and `to`").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "no current active top news")
	} else if err != nil {
		log.WithError(err).Error("could not GetTopNews")
		return nil, status.Error(codes.Internal, "could not GetTopNews")
	}

	return &pb.GetTopNewsReply{
		ImageUrl: res.Files.URL.String,
		Link:     res.Link.String,
		Created:  timestamppb.New(res.Created),
		From:     timestamppb.New(res.From),
		To:       timestamppb.New(res.To),
	}, nil
}
