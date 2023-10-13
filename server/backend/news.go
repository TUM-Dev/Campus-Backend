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

func (s *CampusServer) ListNewsSources(ctx context.Context, _ *pb.ListNewsSourcesRequest) (*pb.ListNewsSourcesReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var sources []model.NewsSource
	if err := s.db.WithContext(ctx).Joins("File").Find(&sources).Error; err != nil {
		log.WithError(err).Error("could not find newsSources")
		return nil, status.Error(codes.Internal, "could not ListNewsSources")
	}

	var resp []*pb.NewsSource
	for _, source := range sources {
		log.WithField("title", source.Title).Trace("sending news source")
		resp = append(resp, &pb.NewsSource{
			Source:  fmt.Sprintf("%d", source.Source),
			Title:   source.Title,
			IconUrl: fmt.Sprintf("https://api.tum.app/files/%s%s", source.File.Path, source.File.Name),
		})
	}
	return &pb.ListNewsSourcesReply{Sources: resp}, nil
}

func (s *CampusServer) ListNews(ctx context.Context, req *pb.ListNewsRequest) (*pb.ListNewsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var newsEntries []model.News
	tx := s.db.WithContext(ctx).Joins("File")
	if req.NewsSource != 0 {
		tx = tx.Where("src = ?", req.NewsSource)
	}
	if req.OldestDateAt.GetSeconds() != 0 || req.OldestDateAt.GetNanos() != 0 {
		tx = tx.Where("date > ?", req.OldestDateAt.AsTime())
	}
	if req.LastNewsId != 0 {
		tx = tx.Where("news > ?", req.LastNewsId)
	}
	if err := tx.Find(&newsEntries).Error; err != nil {
		log.WithError(err).Error("could not find news item")
		return nil, status.Error(codes.Internal, "could not ListNews")
	}

	resp := make([]*pb.News, len(newsEntries))
	for i, item := range newsEntries {
		log.WithField("title", item.Title).Trace("sending news")
		imgUrl := ""
		if item.File != nil {
			imgUrl = fmt.Sprintf("https://api.tum.app/files/%s%s", item.File.Path, item.File.Name)
		}
		resp[i] = &pb.News{
			Id:       item.News,
			Title:    item.Title,
			Text:     item.Description,
			Link:     item.Link,
			ImageUrl: imgUrl,
			Source:   fmt.Sprintf("%d", item.Src),
			Created:  timestamppb.New(item.Created),
			Date:     timestamppb.New(item.Date),
		}
	}
	return &pb.ListNewsReply{News: resp}, nil
}

func (s *CampusServer) ListNewsAlerts(ctx context.Context, req *pb.ListNewsAlertsRequest) (*pb.ListNewsAlertsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var res []*model.NewsAlert
	tx := s.db.WithContext(ctx).Joins("File").Where("news_alert.to >= NOW()")
	if req.LastNewsAlertId != 0 {
		tx = tx.Where("news_alert.alert > ?", req.LastNewsAlertId)
	}
	if err := tx.Find(&res).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "no news alerts")
	} else if err != nil {
		log.WithError(err).Error("could not ListNewsAlerts")
		return nil, status.Error(codes.Internal, "could not ListNewsAlerts")
	}

	var alerts []*pb.NewsAlert
	for _, alert := range res {
		alerts = append(alerts, &pb.NewsAlert{
			ImageUrl: fmt.Sprintf("https://api.tum.app/files/%s%s", alert.File.Path, alert.File.Name),
			Link:     alert.Link.String,
			Created:  timestamppb.New(alert.Created),
			From:     timestamppb.New(alert.From),
			To:       timestamppb.New(alert.To),
		})
	}
	return &pb.ListNewsAlertsReply{Alerts: alerts}, nil
}
