package backend

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *CampusServer) GetNewsSources(ctx context.Context, _ *emptypb.Empty) (newsSources *pb.NewsSourceReply, err error) {
	if err = s.checkDevice(ctx); err != nil {
		return
	}

	var sources []model.NewsSource
	if err := s.db.Joins("Files").Find(&sources).Error; err != nil {
		log.WithError(err).Error("could not find newsSources")
		return nil, status.Error(codes.Internal, "could not GetNewsSources")
	}

	var resp []*pb.NewsSource
	for _, source := range sources {
		log.WithField("title", source.Title).Trace("sending news source")
		resp = append(resp, &pb.NewsSource{
			Source: fmt.Sprintf("%d", source.Source),
			Title:  source.Title,
			Icon:   source.Files.URL.String,
		})
	}
	return &pb.NewsSourceReply{Sources: resp}, nil
}

func (s *CampusServer) GetTopNews(ctx context.Context, _ *emptypb.Empty) (*pb.GetTopNewsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var res *model.NewsAlert
	err := s.db.Joins("Files").Where("NOW() between `from` and `to`").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "no currenty active top news")
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
