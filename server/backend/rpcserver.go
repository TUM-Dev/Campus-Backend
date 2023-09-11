package backend

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns/ios_apns_jwt"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"net"
	"sync"
	"time"
)

func (s *CampusServer) GRPCServe(l net.Listener) error {
	grpcServer := grpc.NewServer()
	pb.RegisterCampusServer(grpcServer, s)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return grpcServer.Serve(l)
}

type CampusServer struct {
	pb.UnimplementedCampusServer
	db                      *gorm.DB
	deviceBuf               *deviceBuffer // deviceBuf stores all devices from recent request and flushes them to db
	iOSNotificationsService *IOSNotificationsService
}

// Verify that CampusServer implements the pb.CampusServer interface
var _ pb.CampusServer = (*CampusServer)(nil)

func New(db *gorm.DB) *CampusServer {
	log.Trace("Server starting up")
	initTagRatingOptions(db)

	return &CampusServer{
		db: db,
		deviceBuf: &deviceBuffer{
			lock:     sync.Mutex{},
			devices:  make(map[string]*model.Devices),
			interval: time.Minute,
		},
		iOSNotificationsService: NewIOSNotificationsService(),
	}
}

func NewIOSNotificationsService() *IOSNotificationsService {
	if err := ios_apns.ValidateRequirementsForIOSNotificationsService(); err != nil {
		log.Warn(err)

		return &IOSNotificationsService{
			APNSToken: nil,
			IsActive:  false,
		}
	}

	token, err := ios_apns_jwt.NewToken()

	if err != nil {
		log.Fatal(err)
	}

	return &IOSNotificationsService{
		APNSToken: token,
		IsActive:  true,
	}
}

func (s *CampusServer) GetNewsSources(ctx context.Context, _ *emptypb.Empty) (newsSources *pb.NewsSourceArray, err error) {
	if err = s.checkDevice(ctx); err != nil {
		return
	}

	var sources []model.NewsSource
	if err := s.db.Find(&sources).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var resp []*pb.NewsSource
	for _, source := range sources {
		var icon model.Files
		if err := s.db.Where("file = ?", source.Icon).First(&icon).Error; err != nil {
			icon = model.Files{File: 0}
		}
		log.WithField("Title", source.Title).Info("sending news source")
		resp = append(resp, &pb.NewsSource{
			Source: fmt.Sprintf("%d", source.Source),
			Title:  source.Title,
			Icon:   icon.URL.String,
		})
	}
	return &pb.NewsSourceArray{Sources: resp}, nil
}

func (s *CampusServer) GetTopNews(ctx context.Context, _ *emptypb.Empty) (*pb.GetTopNewsReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	var res *model.NewsAlert
	err := s.db.Joins("Company").Where("NOW() between `from` and `to`").Limit(1).First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Errorf("Failed to fetch top news")
	} else if res != nil {
		return &pb.GetTopNewsReply{
			//ImageUrl: res.Name,
			Link: res.Link.String,
			To:   timestamppb.New(res.To),
		}, nil
	}
	return &pb.GetTopNewsReply{}, nil
}

func (s *CampusServer) GetIOSNotificationsService() *IOSNotificationsService {
	return s.iOSNotificationsService
}
