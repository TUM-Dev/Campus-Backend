package cron

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
