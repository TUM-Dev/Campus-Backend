package cron

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func cronTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return db, mock
}

func TestCronSchedulerDoesNotBlockOrOverlap(t *testing.T) {
	db, mock := cronTestDB(t)
	svc := New(db)
	svc.pollInterval = 20 * time.Millisecond

	releaseSlow := make(chan struct{})
	slowStarted := make(chan struct{}, 1)
	var slowRuns atomic.Int32
	var fastRuns atomic.Int32

	svc.executeOverride = func(job model.Crontab) error {
		switch job.Cron {
		case 1:
			slowRuns.Add(1)
			select {
			case slowStarted <- struct{}{}:
			default:
			}
			<-releaseSlow
		case 2:
			fastRuns.Add(1)
		}
		return nil
	}

	// First poll: both slow (id 1) and fast (id 2) are due.
	mock.ExpectQuery("SELECT .* FROM `crontab`").WillReturnRows(
		sqlmock.NewRows([]string{"cron", "interval", "lastRun", "type", "id"}).
			AddRow(1, 10, 0, NewsType, nil).
			AddRow(2, 10, 0, CanteenHeadcount, nil),
	)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `crontab`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `crontab`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Second poll while slow (id 1) is still running: both returned as due,
	// slow must be skipped while fast runs again immediately.
	mock.ExpectQuery("SELECT .* FROM `crontab`").WillReturnRows(
		sqlmock.NewRows([]string{"cron", "interval", "lastRun", "type", "id"}).
			AddRow(1, 10, 0, NewsType, nil).
			AddRow(2, 10, 0, CanteenHeadcount, nil),
	)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `crontab`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- svc.RunContext(ctx)
	}()

	<-slowStarted
	require.Eventually(t, func() bool { return fastRuns.Load() >= 2 }, time.Second, 5*time.Millisecond)
	require.Equal(t, int32(1), slowRuns.Load(), "slow job must not overlap while still running")

	cancel()
	select {
	case err := <-done:
		t.Fatalf("RunContext returned before in-flight job completed (err=%v)", err)
	case <-time.After(30 * time.Millisecond):
	}

	close(releaseSlow)
	require.NoError(t, <-done)
}

func TestCronSchedulerRecoversFromPanic(t *testing.T) {
	db, mock := cronTestDB(t)
	svc := New(db)
	svc.executeOverride = func(job model.Crontab) error {
		panic("unexpected cron panic")
	}

	mock.ExpectQuery("SELECT .* FROM `crontab`").WillReturnRows(
		sqlmock.NewRows([]string{"cron", "interval", "lastRun", "type", "id"}).
			AddRow(7, 10, 0, MovieType, nil),
	)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `crontab`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	svc.pollOnce(context.Background())
	svc.wg.Wait()

	// Panic must clean up running map so future runs aren't permanently blocked.
	require.True(t, svc.tryMarkRunning(7))
	require.NoError(t, mock.ExpectationsWereMet())
}
