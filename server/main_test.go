package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestHealthHandler(t *testing.T) {
	t.Run("nil db returns 503", func(t *testing.T) {
		rec := httptest.NewRecorder()
		healthHandler(nil).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
		require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})

	t.Run("healthy db returns 200", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		t.Cleanup(func() { _ = sqlDB.Close() })
		mock.ExpectPing() // gorm.Open ping
		mock.ExpectPing() // healthHandler ping

		db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		healthHandler(db).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "healthy", rec.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("unhealthy db returns 503", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		t.Cleanup(func() { _ = sqlDB.Close() })
		mock.ExpectPing() // gorm.Open ping
		mock.ExpectPing().WillReturnError(errors.New("db down"))

		db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		healthHandler(db).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
		require.Equal(t, http.StatusServiceUnavailable, rec.Code)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPanicRecoveryInterceptors(t *testing.T) {
	t.Run("unary panic becomes Internal error", func(t *testing.T) {
		resp, err := UnaryPanicRecovery(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/api.Campus/Test"}, func(context.Context, any) (any, error) {
			panic("boom")
		})
		require.Nil(t, resp)
		require.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("stream panic becomes Internal error", func(t *testing.T) {
		err := StreamPanicRecovery(nil, nil, &grpc.StreamServerInfo{FullMethod: "/api.Campus/Test"}, func(any, grpc.ServerStream) error {
			panic("stream boom")
		})
		require.Equal(t, codes.Internal, status.Code(err))
	})
}

func TestStreamRequestLoggerLevels(t *testing.T) {
	hook := test.NewGlobal()
	t.Cleanup(hook.Reset)

	err := StreamRequestLogger(nil, nil, &grpc.StreamServerInfo{FullMethod: "/api.Campus/CreateFeedback"}, func(any, grpc.ServerStream) error {
		return nil
	})
	require.NoError(t, err)
	require.NotNil(t, hook.LastEntry())
	require.Equal(t, log.InfoLevel, hook.LastEntry().Level)

	hook.Reset()
	expectedErr := status.Error(codes.InvalidArgument, "bad input")
	err = StreamRequestLogger(nil, nil, &grpc.StreamServerInfo{FullMethod: "/api.Campus/CreateFeedback"}, func(any, grpc.ServerStream) error {
		return expectedErr
	})
	require.Equal(t, expectedErr, err)
	require.NotNil(t, hook.LastEntry())
	require.Equal(t, log.WarnLevel, hook.LastEntry().Level)
}
