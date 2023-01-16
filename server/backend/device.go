package backend

import (
	"context"
	"sync"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// deviceBuffer stores all recent device calls in a buffer and flushes them to the database periodically
type deviceBuffer struct {
	lock     sync.Mutex
	devices  map[string]*model.Devices // key is uuid
	interval time.Duration             // flush interval
}

func (s *CampusServer) RunDeviceFlusher() error {
	for {
		time.Sleep(s.deviceBuf.interval)
		if err := s.deviceBuf.flush(s.db); err != nil {
			log.WithError(err).Error("Error flushing device buffer")
		}
	}
}

// s.deviceBuf.add(md["x-device-id"][0], method[0], osVersion, appVersion)
func (b *deviceBuffer) add(deviceID string, method string, osVersion string, appVersion string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, exists := b.devices[deviceID]; exists {
		b.devices[deviceID].Counter++
	} else {
		b.devices[deviceID] = &model.Devices{
			UUID:       deviceID,
			LastAccess: time.Now(),
			LastAPI:    method,
			OsVersion:  osVersion,
			AppVersion: appVersion,
			Counter:    1,
		}
	}
}

// flush writes all buffered devices to the database
func (b *deviceBuffer) flush(tx *gorm.DB) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	devices := make([]*model.Devices, 0, len(b.devices))
	for _, device := range b.devices {
		devices = append(devices, device)
	}
	if len(b.devices) == 0 {
		return nil
	}
	// store devices in database, update if exists
	err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}}, // conflict on uuid -> update
		DoUpdates: clause.AssignmentColumns([]string{"lastApi", "lastAccess", "osVersion", "appVersion"}),
	}).Create(devices).Error
	if err != nil {
		log.WithError(err).Error("failed to flush device buffer")
	}

	// update number of calls for each device
	err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"counter": gorm.Expr("counter+VALUES(counter)")}),
	}).Create(devices).Error
	if err != nil {
		log.WithError(err).Error("failed to flush device buffer")
	}
	b.devices = make(map[string]*model.Devices)
	return nil
}

var ErrNoDeviceID = status.Errorf(codes.PermissionDenied, "no device id")

// checkDevice checks if the device is approved (TODO: implement)
func (s *CampusServer) checkDevice(ctx context.Context) error {
	var deviceID, method string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Internal, "can't extract metadata from request")
	}
	if len(md["x-device-id"]) == 0 {
		deviceID = "unknown"
		md["x-device-id"] = []string{deviceID}
	}

	// check method header added by middleware. This should always exist.
	if len(md["x-campus-method"]) == 0 {
		log.Info("no method header found for request")
		method = "unknown"
		md["x-campus-method"] = []string{method}
	}

	osVersion := "unknown"
	if len(md["x-os-version"]) > 0 {
		osVersion = md["x-os-version"][0]
	}
	appVersion := "unknown"
	if len(md["x-app-version"]) > 0 {
		appVersion = md["x-app-version"][0]
	}

	// log device to db
	s.deviceBuf.add(deviceID, method, osVersion, appVersion)
	return nil
}
