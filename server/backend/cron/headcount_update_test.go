package cron

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func headcountTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mock.ExpectationsWereMet()); _ = conn.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return db, mock
}

func TestHeadcountCanReturnToZero(t *testing.T) {
	db, mock := headcountTestDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `canteen_head_counts` SET `count`=\\?,`max_count`=\\?,`percent`=\\?,`timestamp`=\\? WHERE `canteen_head_counts`\\.`canteen_id` = \\?").
		WithArgs(uint32(0), uint32(1000), float32(0), sqlmock.AnyArg(), "mensa-garching").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	require.NoError(t, updateDb(&CanteenApInformation{CanteenId: "mensa-garching", MaxCount: 1000}, 0, db))
}

func TestRecordAndPruneHeadCountHistory(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	berlin, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)
	t1 := time.Date(2026, time.November, 17, 12, 2, 0, 0, berlin)
	t2 := time.Date(2026, time.November, 17, 12, 10, 0, 0, berlin)

	require.NoError(t, recordHeadCountHistory(db, "mensa-garching", 400, 40, t1))
	require.NoError(t, recordHeadCountHistory(db, "mensa-garching", 600, 60, t2))

	var row model.CanteenHeadCountHistory
	require.NoError(t, db.First(&row, "canteen_id = ? AND sample_date = ? AND bucket = ?", "mensa-garching", "2026-11-17", 48).Error)
	require.Equal(t, int32(2), row.SampleCount)
	require.InDelta(t, 500.0, row.AvgCount, 0.01)
	require.InDelta(t, 50.0, row.AvgPercent, 0.01)

	// Records older than 180 days are pruned while 90-day-old and recent records survive.
	require.NoError(t, recordHeadCountHistory(db, "mensa-garching", 200, 20, t1.AddDate(0, 0, -90)))
	require.NoError(t, recordHeadCountHistory(db, "mensa-garching", 100, 10, t1.AddDate(0, 0, -190)))
	require.NoError(t, pruneOldHeadCountHistory(db, t1))
	var count int64
	require.NoError(t, db.Model(&model.CanteenHeadCountHistory{}).Count(&count).Error)
	require.Equal(t, int64(2), count)
}
