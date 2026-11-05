package cron

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type newsTestTransport func(*http.Request) (*http.Response, error)

func (f newsTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func newsTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mock.ExpectationsWereMet()); _ = conn.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return db, mock
}

func TestNewsDateFallback(t *testing.T) {
	published := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	updated := published.Add(time.Hour)
	require.Equal(t, published, newsDate(&gofeed.Item{PublishedParsed: &published, UpdatedParsed: &updated}))
	require.Equal(t, updated, newsDate(&gofeed.Item{UpdatedParsed: &updated}))
	require.WithinDuration(t, time.Now(), newsDate(&gofeed.Item{}), time.Second)
}

func TestNewsFeedWithoutPublicationDate(t *testing.T) {
	db, mock := newsTestDB(t)
	client := &http.Client{Transport: newsTestTransport(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`<rss version="2.0"><channel><title>News</title><item><title>Undated</title><link>https://example.com/article</link></item></channel></rss>`)), Header: make(http.Header)}, nil
	})}
	service := &CronService{db: db, client: client}
	mock.ExpectQuery("SELECT `link` FROM `news`").WillReturnRows(sqlmock.NewRows([]string{"link"}))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `news`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `news`").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	require.NoError(t, service.parseNewsFeed(model.NewsSource{Source: 1, URL: null.StringFrom("https://example.com/rss")}))
}

func TestNewsCronKeepsExistingItemsWhenFeedFails(t *testing.T) {
	db, mock := newsTestDB(t)
	requested := false
	service := New(db)
	service.client = &http.Client{Transport: newsTestTransport(func(*http.Request) (*http.Response, error) {
		requested = true
		return nil, errors.New("upstream offline")
	})}
	mock.ExpectQuery("SELECT .* FROM `news_sources`").WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"source", "url"}).AddRow(1, "https://example.com/rss"))
	require.Error(t, service.newsCron(&model.Crontab{ID: null.IntFrom(1)}))
	require.True(t, requested)
	// No DELETE should be issued when the feed cannot be refreshed.
}

func TestNewsCleanupPreservesLinksStillInFeed(t *testing.T) {
	db, mock := newsTestDB(t)
	service := New(db)
	link := "https://example.com/still-active"
	service.client = &http.Client{Transport: newsTestTransport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`<rss version="2.0"><channel><title>News</title><item><title>Old but active</title><link>` + link + `</link></item></channel></rss>`)), Header: make(http.Header)}, nil
	})}
	// A link can have been created over a year ago and still be in the feed.
	mock.ExpectQuery("SELECT `link` FROM `news`").WillReturnRows(sqlmock.NewRows([]string{"link"}).AddRow(link))
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `news` WHERE .*`link` NOT IN").WithArgs(int64(1), sqlmock.AnyArg(), link).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	require.NoError(t, service.parseNewsFeed(model.NewsSource{Source: 1, URL: null.StringFrom("https://example.com/rss")}))
}
