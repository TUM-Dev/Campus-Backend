package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
)

// Check the database-to-REST path separately from the live importer test.
func TestNewsMySQLGateway(t *testing.T) {
	if os.Getenv("NEWS_TEST_MYSQL") != "1" {
		t.Skip("set NEWS_TEST_MYSQL=1 to test the news API against MySQL")
	}
	db := utils.SetupTestContainer(t.Context(), t)
	icon := model.File{Name: "test.png", Path: "news/sources/"}
	require.NoError(t, db.Create(&icon).Error)
	source := model.NewsSource{Title: "Test source", FileID: icon.File}
	require.NoError(t, db.Omit("File").Create(&source).Error)
	article := model.News{
		Title: "Current news", Description: "Test article", Link: "https://example.com/news",
		Date: time.Now().UTC().Truncate(time.Second), NewsSourceID: source.Source,
	}
	require.NoError(t, db.Omit("NewsSource", "File").Create(&article).Error)

	mux := runtime.NewServeMux()
	require.NoError(t, pb.RegisterCampusHandlerServer(t.Context(), mux, New(db)))
	for _, test := range []struct {
		path string
		want int
	}{
		{"/news", 1},
		{fmt.Sprintf("/news?news_source=%d", source.Source), 1},
		{"/news?oldest_date_at=" + article.Date.Add(-time.Hour).Format(time.RFC3339), 1},
		{"/news?oldest_date_at=" + article.Date.Add(time.Hour).Format(time.RFC3339), 0},
		{fmt.Sprintf("/news?last_news_id=%d", article.News), 0},
	} {
		t.Run(test.path, func(t *testing.T) {
			resp := httptest.NewRecorder()
			mux.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, test.path, nil))
			require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
			var articles []map[string]any
			require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &articles))
			require.Len(t, articles, test.want)
			if test.want > 0 {
				require.Equal(t, article.Title, articles[0]["title"])
				require.Equal(t, source.Title, articles[0]["sourceTitle"])
			}
		})
	}
}

func TestRestoreNewsFeedMigrationSeedsProductionSources(t *testing.T) {
	if os.Getenv("NEWS_TEST_MYSQL") != "1" {
		t.Skip("set NEWS_TEST_MYSQL=1 to test news migration against MySQL")
	}
	db := utils.SetupTestContainer(t.Context(), t)
	icon := model.File{Name: "src_1.png", Path: "news/sources/"}
	require.NoError(t, db.Create(&icon).Error)
	// Seed existing production rows (1 with old URL, 5 & 14 with empty URL, 18 with old HTML URL)
	for _, s := range []model.NewsSource{
		{Source: 1, Title: "TUM News", URL: null.StringFrom("http://www.tum.de/nc/die-tum/aktuelles/?type=100"), FileID: icon.File},
		{Source: 5, Title: "TUM Pressenews", URL: null.StringFrom(""), FileID: icon.File},
		{Source: 14, Title: "TUM Alumni Events", URL: null.StringFrom(""), FileID: icon.File},
		{Source: 18, Title: "TUM Community", URL: null.StringFrom("https://www.community.tum.de/en/?event=true&eventfeed=true"), FileID: icon.File},
	} {
		require.NoError(t, db.Omit("File").Create(&s).Error)
	}
	// Seed an orphaned news cron job pointing to missing source 99
	require.NoError(t, db.Create(&model.Crontab{Interval: 43200, LastRun: 123, Type: null.StringFrom("news"), ID: null.IntFrom(99)}).Error)

	sqlBytes, err := os.ReadFile("migration/sql/20261122183040_restore_tum_news_feed.sql")
	require.NoError(t, err)
	for _, stmt := range strings.Split(string(sqlBytes), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		require.NoError(t, db.Exec(stmt).Error)
	}

	var sources []model.NewsSource
	require.NoError(t, db.Find(&sources).Error)
	require.Len(t, sources, 15)

	var jobs []model.Crontab
	require.NoError(t, db.Where("type = 'news'").Find(&jobs).Error)
	require.Len(t, jobs, 15)
}
