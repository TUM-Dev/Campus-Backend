package cron

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/guregu/null"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/require"
)

// Opt in with NEWS_FEED_URL; normal tests never depend on an external feed.
// NEWS_TEST_MYSQL=1 also runs the actual importer twice in a disposable database
// using SetupTestContainer (Docker, or an existing server via TEST_MYSQL_DSN).
func TestNewsLiveFeed(t *testing.T) {
	url := os.Getenv("NEWS_FEED_URL")
	if url == "" {
		t.Skip("set NEWS_FEED_URL to inspect a live RSS/Atom feed")
	}
	resp, err := webfetch.Get(url)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	finalURL := url
	if resp.Request != nil && resp.Request.URL != nil {
		finalURL = resp.Request.URL.String()
	}
	feed, err := parseNewsFeedContent(finalURL, string(body))
	require.NoError(t, err)
	require.NotEmpty(t, feed.Items, "feed contains no items")

	links := make(map[string]bool)
	var newest time.Time
	var undated, unlinked int
	for _, item := range feed.Items {
		if item.Link == "" {
			unlinked++
			continue // The importer skips these, too.
		}
		require.False(t, links[item.Link], "duplicate link in feed: %s", item.Link)
		links[item.Link] = true
		// Check the production column limits before attempting a batch insert.
		require.LessOrEqual(t, utf8.RuneCountInString(item.Link), 190, "link exceeds news.link: %s", item.Link)
		require.LessOrEqual(t, len(item.Title), 255, "title exceeds news.title: %s", item.Title)
		require.LessOrEqual(t, len(bluemonday.StrictPolicy().Sanitize(item.Description)), 65535, "description exceeds news.description: %s", item.Link)
		if item.PublishedParsed == nil && item.UpdatedParsed == nil {
			undated++
			continue
		}
		if date := newsDate(item); date.After(newest) {
			newest = date
		}
	}
	t.Logf("url=%s final_url=%s title=%q items=%d importable=%d unlinked=%d undated=%d newest_dated=%s",
		url, resp.Request.URL, feed.Title, len(feed.Items), len(links), unlinked, undated, newest.Format(time.RFC3339))
	require.NotEmpty(t, links, "feed has no importable links")

	t.Run("mysql", func(t *testing.T) {
		if os.Getenv("NEWS_TEST_MYSQL") != "1" {
			t.Skip("set NEWS_TEST_MYSQL=1 to verify database insertion and repeat imports")
		}
		// Replay the exact inspected response so a feed changing between requests
		// cannot invalidate the count or the repeat-import check.
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = io.Copy(w, strings.NewReader(string(body)))
		}))
		defer upstream.Close()
		db := utils.SetupTestContainer(t.Context(), t)
		icon := model.File{Name: "live-test.png", Path: "news/sources/"}
		require.NoError(t, db.Create(&icon).Error)
		source := model.NewsSource{Title: feed.Title, URL: null.StringFrom(upstream.URL), FileID: icon.File}
		require.NoError(t, db.Omit("File").Create(&source).Error)
		service := New(db)
		job := &model.Crontab{ID: null.IntFrom(source.Source)}
		for run := 1; run <= 2; run++ {
			require.NoError(t, service.newsCron(job), "import run %d", run)
			var count int64
			require.NoError(t, db.Model(&model.News{}).Where("src = ?", source.Source).Count(&count).Error)
			require.Equal(t, int64(len(links)), count, "import run %d", run)
			t.Logf("import run %d: %d persisted articles", run, count)
		}
		var stored model.NewsSource
		require.NoError(t, db.First(&stored, source.Source).Error)
		require.Equal(t, source.URL, stored.URL, "import must preserve source configuration")
	})
}
