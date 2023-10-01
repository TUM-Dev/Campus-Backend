package backend

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CampusServer) GetNewsSources(ctx context.Context, _ *pb.GetNewsSourcesRequest) (*pb.GetNewsSourcesReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var sources []model.NewsSource
	if err := s.db.WithContext(ctx).Joins("File").Find(&sources).Error; err != nil {
		log.WithError(err).Error("could not find newsSources")
		return nil, status.Error(codes.Internal, "could not GetNewsSources")
	}

	var resp []*pb.NewsSource
	for _, source := range sources {
		log.WithField("title", source.Title).Trace("sending news source")
		resp = append(resp, &pb.NewsSource{
			Source: fmt.Sprintf("%d", source.Source),
			Title:  source.Title,
			Icon:   source.File.URL.String,
		})
	}
	return &pb.GetNewsSourcesReply{Sources: resp}, nil
}

func (s *CampusServer) GetNews(ctx context.Context, req *pb.GetNewsRequest) (*pb.GetNewsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var newsEntries []model.News
	tx := s.db.WithContext(ctx).Joins("File")
	if req.NewsSource != 0 {
		tx = tx.Where("src = ?", req.NewsSource)
	}
	if req.LastNewsId != 0 {
		tx = tx.Where("news > ?", req.LastNewsId)
	}
	if err := tx.Find(&newsEntries).Error; err != nil {
		log.WithError(err).Error("could not find news item")
		return nil, status.Error(codes.Internal, "could not GetNews")
	}

	resp := make([]*pb.News, len(newsEntries))
	for i, item := range newsEntries {
		log.WithField("title", item.Title).Trace("sending news")
		resp[i] = &pb.News{
			Id:       item.News,
			Title:    item.Title,
			Text:     item.Description,
			Link:     item.Link,
			ImageUrl: item.Image.String,
			Source:   fmt.Sprintf("%d", item.Src),
			Created:  timestamppb.New(item.Created),
			Date:     timestamppb.New(item.Date),
		}
	}
	return &pb.GetNewsReply{News: resp}, nil
}

func (s *CampusServer) GetNewsAlerts(ctx context.Context, req *pb.GetNewsAlertsRequest) (*pb.GetNewsAlertsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var res []*model.NewsAlert
	tx := s.db.WithContext(ctx).Joins("Files").Where("news_alert.to >= NOW()")
	if req.LastNewsAlertId != 0 {
		tx = tx.Where("news_alert.alert > ?", req.LastNewsAlertId)
	}
	if err := tx.Find(&res).Error; errors.Is(err, gorm.ErrRecordNotFound) {
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
