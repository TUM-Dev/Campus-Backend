package backend

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *CampusServer) ListNewsSources(ctx context.Context, _ *pb.ListNewsSourcesRequest) (*pb.ListNewsSourcesReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	sources, err := s.getNewsSources(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not ListNewsSources")
	}

	var resp []*pb.NewsSource
	for _, source := range sources {
		resp = append(resp, &pb.NewsSource{
			Source:  fmt.Sprintf("%d", source.Source),
			Title:   source.Title,
			IconUrl: source.File.FullExternalUrl(),
		})
	}
	return &pb.ListNewsSourcesReply{Sources: resp}, nil
}

const CacheKeyAllNewsSources = "all_news_sources"

func (s *CampusServer) getNewsSources(ctx context.Context) ([]model.NewsSource, error) {
	if newsSources, ok := s.newsSourceCache.Get(CacheKeyAllNewsSources); ok {
		return newsSources, nil
	}
	var sources []model.NewsSource
	if err := s.db.WithContext(ctx).Joins("File").Find(&sources).Error; err != nil {
		return nil, err
	}
	s.newsSourceCache.Add(CacheKeyAllNewsSources, sources)
	return sources, nil
}

func (s *CampusServer) ListNews(ctx context.Context, req *pb.ListNewsRequest) (*pb.ListNewsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var newsEntries, err = s.getNews(ctx, req.NewsSource, req.LastNewsId, req.OldestDateAt.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, "could not ListNews")
	}

	var resp []*pb.News
	for _, item := range newsEntries {
		imgUrl := ""
		if item.File != nil {
			imgUrl = item.File.FullExternalUrl()
		}
		resp = append(resp, &pb.News{
			Id:            item.News,
			Title:         item.Title,
			Text:          item.Description,
			Link:          item.Link,
			ImageUrl:      imgUrl,
			SourceId:      fmt.Sprintf("%d", item.NewsSource.Source),
			SourceTitle:   item.NewsSource.Title,
			SourceIconUrl: item.NewsSource.File.FullExternalUrl(),
			Created:       timestamppb.New(item.Created),
			Date:          timestamppb.New(item.Date),
		})
	}
	return &pb.ListNewsReply{News: resp}, nil
}

func (s *CampusServer) getNews(ctx context.Context, sourceID int32, lastNewsID int32, oldestDateAt time.Time) ([]model.News, error) {
	cacheKey := fmt.Sprintf("%d_%d_%d", sourceID, oldestDateAt.Second(), lastNewsID)

	if news, ok := s.newsCache.Get(cacheKey); ok {
		return news, nil
	}

	var news []model.News
	tx := s.db.WithContext(ctx).
		Joins("File").
		Joins("NewsSource").
		Joins("NewsSource.File")
	if sourceID != 0 {
		tx = tx.Where("src = ?", sourceID)
	}
	if oldestDateAt.Unix() != 0 {
		tx = tx.Where("date > ?", oldestDateAt)
	}
	if lastNewsID != 0 {
		tx = tx.Where("news > ?", lastNewsID)
	}
	if err := tx.Find(&news).Error; err != nil {
		return nil, err
	}
	s.newsCache.Add(cacheKey, news)
	return news, nil
}

func (s *CampusServer) ListNewsAlerts(ctx context.Context, req *pb.ListNewsAlertsRequest) (*pb.ListNewsAlertsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var res []*model.NewsAlert
	tx := s.db.WithContext(ctx).
		Joins("File").
		Where("news_alerts.to >= NOW()")
	if req.LastNewsAlertId != 0 {
		tx = tx.Where("news_alerts.news_alert > ?", req.LastNewsAlertId)
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
			ImageUrl: alert.File.FullExternalUrl(),
			Link:     alert.Link.String,
			Created:  timestamppb.New(alert.Created),
			From:     timestamppb.New(alert.From),
			To:       timestamppb.New(alert.To),
		})
	}
	return &pb.ListNewsAlertsReply{Alerts: alerts}, nil
}
