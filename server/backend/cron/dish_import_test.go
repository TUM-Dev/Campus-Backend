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

func TestMenuImportISOYearBoundaries(t *testing.T) {
	for _, tt := range []struct {
		date            string
		year, week, day int
	}{
		{"2025-12-28", 2025, 52, 6},
		{"2025-12-29", 2026, 1, 0},
		{"2026-12-28", 2026, 53, 0},
		{"2027-01-01", 2026, 53, 4},
		{"2027-01-04", 2027, 1, 0},
	} {
		t.Run(tt.date, func(t *testing.T) {
			db, mock := dishImportTestDB(t)
			menu := CanteenDays{Days: []CanteenDate{{Date: tt.date, Dates: []CanteenDish{{Name: "Pasta", DishType: "Pasta"}}}}}
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT .* FROM `dishes`").WithArgs("Pasta", int64(7), 1).
				WillReturnRows(sqlmock.NewRows([]string{"dish", "name", "type", "cafeteriaID"}).AddRow(42, "Pasta", "Pasta", 7))
			mock.ExpectQuery("SELECT .* FROM `dishes_of_the_weeks`").
				WithArgs(tt.day, int64(42), tt.week, tt.year, 1).
				WillReturnRows(sqlmock.NewRows([]string{"dishesOfTheWeek", "dishID", "year", "week", "day"}))
			mock.ExpectExec("INSERT INTO `dishes_of_the_weeks`").
				WithArgs(tt.year, tt.week, tt.day, int64(42)).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			require.NoError(t, storeCanteenMenu(db, 7, menu, tt.year, tt.week))
		})
	}
}

func TestCanteenDirectoryNetworkFailure(t *testing.T) {
	old := webfetch.Client
	t.Cleanup(func() { webfetch.Client = old })
	webfetch.Client = &http.Client{Transport: dishImportTestTransport(func(*http.Request) (*http.Response, error) { return nil, errors.New("offline") })}
	require.ErrorContains(t, downloadCanteenNames(&CronService{}), "offline")
}
