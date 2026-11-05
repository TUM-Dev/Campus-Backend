package backend

import (
	"context"
	"sync"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// deviceBuffer stores all recent device calls in a buffer and flushes them to the database periodically
type deviceBuffer struct {
	lock    sync.Mutex
	devices map[string]*model.Device // key is uuid
}

func newDeviceBuffer() *deviceBuffer {
	return &deviceBuffer{
		lock:    sync.Mutex{},
		devices: make(map[string]*model.Device),
	}
}

const FlushingInterval = time.Minute

func (s *CampusServer) RunDeviceFlusher() error {
	for {
		time.Sleep(FlushingInterval)
		if err := s.deviceBuf.flush(s.db); err != nil {
			log.WithError(err).Error("Error flushing device buffer")
		}
	}
}

func (b *deviceBuffer) add(deviceID string, method string, osVersion string, appVersion string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, exists := b.devices[deviceID]; exists {
		b.devices[deviceID].Counter++
	} else {
		b.devices[deviceID] = &model.Device{
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
	devices := make([]*model.Device, 0, len(b.devices))
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
	b.devices = make(map[string]*model.Device)
	return nil
}

// checkDevice records the calling device (x-device-id, x-os-version and
// x-app-version metadata) together with the called method in the device
// buffer, which is periodically flushed to the devices table.
// It does not authenticate or reject devices.
func (s *CampusServer) checkDevice(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Internal, "can't extract metadata from request")
	}

	method, ok := grpc.Method(ctx)
	if !ok {
		method = "unknown"
	}

	s.deviceBuf.add(
		firstOrUnknown(md, "x-device-id"),
		method,
		firstOrUnknown(md, "x-os-version"),
		firstOrUnknown(md, "x-app-version"),
	)
	return nil
}

// firstOrUnknown returns the first value of the metadata key, or "unknown" if it is missing or empty.
func firstOrUnknown(md metadata.MD, key string) string {
	if v := md.Get(key); len(v) > 0 && v[0] != "" {
		return v[0]
	}
	return "unknown"
}

func (s *CampusServer) CreateDevice(_ context.Context, req *pb.CreateDeviceRequest) (*pb.CreateDeviceReply, error) {
	if err := ValidateCreateDevice(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	switch req.DeviceType {
	case pb.DeviceType_ANDROID:
		return nil, status.Error(codes.Unimplemented, "android device creation not implemented")
	case pb.DeviceType_IOS:
		return nil, status.Error(codes.Unimplemented, "ios device creation not implemented")
	case pb.DeviceType_WINDOWS:
		return nil, status.Error(codes.Unimplemented, "windows device creation not implemented")
	}

	return nil, status.Error(codes.InvalidArgument, "invalid device type")
}

func (s *CampusServer) DeleteDevice(_ context.Context, req *pb.DeleteDeviceRequest) (*pb.DeleteDeviceReply, error) {
	if err := ValidateDeleteDevice(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	switch req.DeviceType {
	case pb.DeviceType_ANDROID:
		return nil, status.Error(codes.Unimplemented, "android device remove not implemented")
	case pb.DeviceType_IOS:
		return nil, status.Error(codes.Unimplemented, "ios device remove not implemented")
	case pb.DeviceType_WINDOWS:
		return nil, status.Error(codes.Unimplemented, "windows device remove not implemented")
	}

	return nil, status.Error(codes.InvalidArgument, "invalid device type")
}
