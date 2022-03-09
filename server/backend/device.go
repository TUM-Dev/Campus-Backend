package backend

import (
	"context"
	"errors"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
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
		err := s.deviceBuf.flush(s.db)
		if err != nil {
			log.WithError(err).Error("Error flushing device buffer")
		}
	}
}

func (b *deviceBuffer) add(device model.Devices) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, exists := b.devices[device.UUID]; exists {
		b.devices[device.UUID].Counter++
	} else {
		b.devices[device.UUID] = &device
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

var ErrNoDeviceID = errors.New("no device id")

// checkDevice checks if the device is approved (TODO: implement)
func (s *CampusServer) checkDevice(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Internal, "can't extract metadata from request")
	}
	log.Println()
	if len(md["x-device-id"]) == 0 && md["x-forwarded-for"][0] != "::1" {
		return ErrNoDeviceID
	}

	method := md["x-campus-method"]
	if len(method) == 0 {
		return status.Error(codes.Internal, "can't extract method from request")
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
	s.deviceBuf.add(model.Devices{
		UUID:       md["x-device-id"][0],
		LastAPI:    method[0],
		LastAccess: time.Now(),
		OsVersion:  osVersion,
		AppVersion: appVersion,
		Counter:    1, // initially 1 call, incremented on subsequent calls
	})
	return nil
}
