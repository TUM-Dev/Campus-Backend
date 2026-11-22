package cron

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

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

func TestDownloadDailyDishesFetchesCurrentAndNextWeek(t *testing.T) {
	old := webfetch.Client
	t.Cleanup(func() { webfetch.Client = old })

	var requestedURLs []string
	webfetch.Client = &http.Client{Transport: dishImportTestTransport(func(r *http.Request) (*http.Response, error) {
		requestedURLs = append(requestedURLs, r.URL.String())
		if strings.HasSuffix(r.URL.Path, "/2026/45.json") {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(`{"days":[{"date":"2026-11-04","dishes":[{"name":"Pasta","dish_type":"Pasta"}]}]}`)),
			}, nil
		}
		// Next week (46) not published yet -> 404 must be tolerated gracefully.
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
			Body:       io.NopCloser(strings.NewReader("not found")),
		}, nil
	})}

	db, mock := dishImportTestDB(t)
	mock.ExpectQuery("SELECT \\* FROM `canteens`").WillReturnRows(
		sqlmock.NewRows([]string{"cafeteria", "name"}).AddRow(7, "MENSA_GARCHING"),
	)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FROM `dishes`").WithArgs("Pasta", int64(7), 1).
		WillReturnRows(sqlmock.NewRows([]string{"dish", "name", "type", "cafeteriaID"}).AddRow(42, "Pasta", "Pasta", 7))
	mock.ExpectQuery("SELECT .* FROM `dishes_of_the_weeks`").WithArgs(2, int64(42), 45, 2026, 1).
		WillReturnRows(sqlmock.NewRows([]string{"dishesOfTheWeek", "dishID", "year", "week", "day"}).AddRow(1, 42, 2026, 45, 2))
	mock.ExpectCommit()

	now := time.Date(2026, time.November, 4, 12, 0, 0, 0, time.UTC)
	require.NoError(t, downloadDailyDishesAt(&CronService{db: db}, now))
	require.Equal(t, []string{
		"https://tum-dev.github.io/eat-api/mensa-garching/2026/45.json",
		"https://tum-dev.github.io/eat-api/mensa-garching/2026/46.json",
	}, requestedURLs)
}

func TestAddDishTagsToMappingDeduplicatesAndExcludes(t *testing.T) {
	db, mock := dishImportTestDB(t)
	// Tag 10 matched twice via two expressions; tag 20 matched and excluded.
	mock.ExpectQuery("SELECT `nameTagID` FROM `included_dish_name_tag_options`").
		WithArgs("currywurst mit chicken hähnchen").
		WillReturnRows(sqlmock.NewRows([]string{"nameTagID"}).AddRow(10).AddRow(10).AddRow(20))
	mock.ExpectQuery("SELECT `nameTagID` FROM `excluded_dish_name_tag_options`").
		WithArgs("currywurst mit chicken hähnchen").
		WillReturnRows(sqlmock.NewRows([]string{"nameTagID"}).AddRow(20))
	// Only tag 10 should be inserted once via FirstOrCreate.
	mock.ExpectQuery("SELECT .* FROM `dish_to_dish_name_tags`").
		WithArgs(int64(99), int64(10), 1).
		WillReturnRows(sqlmock.NewRows([]string{"dishToDishNameTag", "dishID", "nameTagID"}))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `dish_to_dish_name_tags`").
		WithArgs(int64(99), int64(10)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	addDishTagsToMapping(99, "Currywurst mit Chicken Hähnchen", db)
}

func TestDownloadDailyDishesCrossesISOYearBoundary(t *testing.T) {
	old := webfetch.Client
	t.Cleanup(func() { webfetch.Client = old })

	var requestedURLs []string
	webfetch.Client = &http.Client{Transport: dishImportTestTransport(func(r *http.Request) (*http.Response, error) {
		requestedURLs = append(requestedURLs, r.URL.String())
		if strings.HasSuffix(r.URL.Path, "/2026/53.json") {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(`{"days":[{"date":"2026-12-28","dishes":[{"name":"Pasta","dish_type":"Pasta"}]}]}`)),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
			Body:       io.NopCloser(strings.NewReader("not found")),
		}, nil
	})}

	db, mock := dishImportTestDB(t)
	mock.ExpectQuery("SELECT \\* FROM `canteens`").WillReturnRows(
		sqlmock.NewRows([]string{"cafeteria", "name"}).AddRow(7, "MENSA_GARCHING"),
	)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FROM `dishes`").WithArgs("Pasta", int64(7), 1).
		WillReturnRows(sqlmock.NewRows([]string{"dish", "name", "type", "cafeteriaID"}).AddRow(42, "Pasta", "Pasta", 7))
	mock.ExpectQuery("SELECT .* FROM `dishes_of_the_weeks`").WithArgs(0, int64(42), 53, 2026, 1).
		WillReturnRows(sqlmock.NewRows([]string{"dishesOfTheWeek", "dishID", "year", "week", "day"}).AddRow(1, 42, 2026, 53, 0))
	mock.ExpectCommit()

	// 2026-12-28 is Monday of ISO week 2026-W53; +7 days is ISO week 2027-W01.
	now := time.Date(2026, time.December, 28, 12, 0, 0, 0, time.UTC)
	require.NoError(t, downloadDailyDishesAt(&CronService{db: db}, now))
	require.Equal(t, []string{
		"https://tum-dev.github.io/eat-api/mensa-garching/2026/53.json",
		"https://tum-dev.github.io/eat-api/mensa-garching/2027/01.json",
	}, requestedURLs)
}

func TestBackfillMissingDishNameTags(t *testing.T) {
	db, mock := dishImportTestDB(t)
	mock.ExpectQuery("SELECT dish, name FROM `dishes`").WillReturnRows(
		sqlmock.NewRows([]string{"dish", "name"}).AddRow(5, "Kürbis Suppe"),
	)
	mock.ExpectQuery("SELECT `nameTagID` FROM `included_dish_name_tag_options`").
		WithArgs("kürbis suppe").
		WillReturnRows(sqlmock.NewRows([]string{"nameTagID"}).AddRow(3))
	mock.ExpectQuery("SELECT `nameTagID` FROM `excluded_dish_name_tag_options`").
		WithArgs("kürbis suppe").
		WillReturnRows(sqlmock.NewRows([]string{"nameTagID"}))
	mock.ExpectQuery("SELECT .* FROM `dish_to_dish_name_tags`").
		WithArgs(int64(5), int64(3), 1).
		WillReturnRows(sqlmock.NewRows([]string{"dishToDishNameTag", "dishID", "nameTagID"}).AddRow(1, 5, 3))

	require.NoError(t, BackfillMissingDishNameTags(db))
}
