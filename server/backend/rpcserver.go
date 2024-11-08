package backend

import (
	"sync"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/hashicorp/golang-lru/v2/expirable"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampusServer struct {
	pb.UnimplementedCampusServer
	feedbackEmailLastReuestAt *sync.Map
	db                        *gorm.DB
	deviceBuf                 *deviceBuffer // deviceBuf stores all devices from recent request and flushes them to db
	newsSourceCache           *expirable.LRU[string, []model.NewsSource]
	newsCache                 *expirable.LRU[string, []model.News]
	moviesCache               *expirable.LRU[string, []model.Movie]
}

// Verify that CampusServer implements the pb.CampusServer interface
var _ pb.CampusServer = (*CampusServer)(nil)

func New(db *gorm.DB) *CampusServer {
	log.Trace("Server starting up")
	return &CampusServer{
		db:                        db,
		deviceBuf:                 newDeviceBuffer(),
		feedbackEmailLastReuestAt: &sync.Map{},
		newsSourceCache:           expirable.NewLRU[string, []model.NewsSource](1, nil, time.Hour*6),
		newsCache:                 expirable.NewLRU[string, []model.News](1024, nil, time.Minute*30),
	}
}
