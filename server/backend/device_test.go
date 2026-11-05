package backend

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// fakeStream lets tests set the method returned by grpc.Method(ctx).
type fakeStream struct{ method string }

func (f *fakeStream) Method() string                 { return f.method }
func (f *fakeStream) SetHeader(_ metadata.MD) error  { return nil }
func (f *fakeStream) SendHeader(_ metadata.MD) error { return nil }
func (f *fakeStream) SetTrailer(_ metadata.MD) error { return nil }

func TestCheckDeviceRecordsMetadata(t *testing.T) {
	s := CampusServer{deviceBuf: newDeviceBuffer()}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-device-id", "dev-1",
		"x-os-version", "iOS 26",
		"x-app-version", "5.1.0",
	))
	ctx = grpc.NewContextWithServerTransportStream(ctx, &fakeStream{method: "/api.Campus/ListNewsSources"})

	require.NoError(t, s.checkDevice(ctx))
	require.NoError(t, s.checkDevice(ctx))

	require.Len(t, s.deviceBuf.devices, 1)
	d := s.deviceBuf.devices["dev-1"]
	require.NotNil(t, d)
	require.Equal(t, "dev-1", d.UUID)
	require.Equal(t, "/api.Campus/ListNewsSources", d.LastAPI)
	require.Equal(t, "iOS 26", d.OsVersion)
	require.Equal(t, "5.1.0", d.AppVersion)
	require.Equal(t, 2, int(d.Counter))
}

func TestCheckDeviceMissingMetadataFallsBackToUnknown(t *testing.T) {
	s := CampusServer{deviceBuf: newDeviceBuffer()}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})

	require.NoError(t, s.checkDevice(ctx))

	d := s.deviceBuf.devices["unknown"]
	require.NotNil(t, d)
	require.Equal(t, "unknown", d.LastAPI)
	require.Equal(t, "unknown", d.OsVersion)
	require.Equal(t, "unknown", d.AppVersion)
}

func TestCheckDeviceWithoutMetadataFails(t *testing.T) {
	s := CampusServer{deviceBuf: newDeviceBuffer()}
	require.Error(t, s.checkDevice(context.Background()))
}
