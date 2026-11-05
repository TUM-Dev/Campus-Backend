package cron

import (
	"errors"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dishImportTestTransport func(*http.Request) (*http.Response, error)

func (f dishImportTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func dishImportTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mock.ExpectationsWereMet()); _ = conn.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return db, mock
}

func TestMenuImportRepairsMissingWeeklyEntries(t *testing.T) {
	db, mock := dishImportTestDB(t)
	menu := CanteenDays{Days: []CanteenDate{{Date: "2026-11-04", Dates: []CanteenDish{{Name: "Pasta", DishType: "Pasta"}}}}}
	// An existing dish can still be missing its weekly mapping. The second run
	// must not add another mapping, even if other cafeterias were imported first.
	for run := 0; run < 2; run++ {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT .* FROM `dishes`").WithArgs("Pasta", int64(7), 1).
			WillReturnRows(sqlmock.NewRows([]string{"dish", "name", "type", "cafeteriaID"}).AddRow(42, "Pasta", "Pasta", 7))
		rows := sqlmock.NewRows([]string{"dishesOfTheWeek", "dishID", "year", "week", "day"})
		if run == 1 {
			rows.AddRow(1, 42, 2026, 45, 2)
		}
		mock.ExpectQuery("SELECT .* FROM `dishes_of_the_weeks`").WithArgs(2, int64(42), 45, 2026, 1).WillReturnRows(rows)
		if run == 0 {
			mock.ExpectExec("INSERT INTO `dishes_of_the_weeks`").WithArgs(2026, 45, 2, int64(42)).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()
		require.NoError(t, storeCanteenMenu(db, 7, menu, 2026, 45))
	}
}

func TestInvalidMenuDoesNotWriteDishes(t *testing.T) {
	db, mock := dishImportTestDB(t)
	for _, date := range []string{"not-a-date", "2026-11-11"} {
		mock.ExpectBegin()
		mock.ExpectRollback()
		require.Error(t, storeCanteenMenu(db, 7, CanteenDays{Days: []CanteenDate{{Date: date}}}, 2026, 45))
	}
}

func TestCanteenDirectoryNetworkFailure(t *testing.T) {
	old := webfetch.Client
	t.Cleanup(func() { webfetch.Client = old })
	webfetch.Client = &http.Client{Transport: dishImportTestTransport(func(*http.Request) (*http.Response, error) { return nil, errors.New("offline") })}
	require.ErrorContains(t, downloadCanteenNames(&CronService{}), "offline")
}
