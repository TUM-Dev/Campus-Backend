package cron

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

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

func TestNewsCronSkipsMissingSourceGracefully(t *testing.T) {
	db, mock := newsTestDB(t)
	service := New(db)
	mock.ExpectQuery("SELECT .* FROM `news_sources`").WithArgs(int64(99), 1).
		WillReturnError(gorm.ErrRecordNotFound)
	require.NoError(t, service.newsCron(&model.Crontab{Cron: 7, ID: null.IntFrom(99)}))
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

func TestNewsDeduplicatesCrossSourceAndInFeedLinks(t *testing.T) {
	db, mock := newsTestDB(t)
	service := New(db)
	sharedLink := "https://www.tum.de/aktuelles/shared"
	freshLink := "https://www.tum.de/aktuelles/fresh"
	feedXML := `<rss version="2.0"><channel><title>News</title>` +
		`<item><title>From another source</title><link>` + sharedLink + `</link></item>` +
		`<item><title>Fresh 1</title><link>` + freshLink + `</link></item>` +
		`<item><title>Duplicate in same feed</title><link>` + freshLink + `</link></item>` +
		`</channel></rss>`
	service.client = &http.Client{Transport: newsTestTransport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(feedXML)), Header: make(http.Header)}, nil
	})}

	// Query checks global news.link (no src = ? filter) with deduplicated feed links.
	mock.ExpectQuery("SELECT `link` FROM `news` WHERE `link` IN \\(\\?,\\?\\)").
		WithArgs(sharedLink, freshLink).
		WillReturnRows(sqlmock.NewRows([]string{"link"}).AddRow(sharedLink))
	// Only the single fresh link should be inserted once.
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `news` .* ON DUPLICATE KEY UPDATE").
		WithArgs(sqlmock.AnyArg(), "Fresh 1", "", int64(5), freshLink, "", nil, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `news` WHERE .*`link` NOT IN").
		WithArgs(int64(5), sqlmock.AnyArg(), sharedLink, freshLink).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	require.NoError(t, service.parseNewsFeed(model.NewsSource{Source: 5, URL: null.StringFrom("https://example.com/rss")}))
}

func TestNewsSkipsOversizedLinkAndTruncatesLongTitle(t *testing.T) {
	db, mock := newsTestDB(t)
	service := New(db)
	tooLongLink := "https://example.com/" + strings.Repeat("a", maxNewsLinkRunes)
	validLink := "https://example.com/ok"
	// 200 umlauts = 400 UTF-8 bytes (> 255 tinytext limit).
	longTitle := strings.Repeat("ä", 200)
	feedXML := `<rss version="2.0"><channel><title>News</title>` +
		`<item><title>Too long link</title><link>` + tooLongLink + `</link></item>` +
		`<item><title>` + longTitle + `</title><link>` + validLink + `</link></item>` +
		`</channel></rss>`
	service.client = &http.Client{Transport: newsTestTransport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(feedXML)), Header: make(http.Header)}, nil
	})}

	expectedTitle := truncateUTF8Bytes(longTitle, maxNewsTitleBytes)
	require.LessOrEqual(t, len(expectedTitle), maxNewsTitleBytes)
	require.True(t, utf8.ValidString(expectedTitle))
	require.Equal(t, 254, len(expectedTitle)) // 127 two-byte 'ä' runes

	mock.ExpectQuery("SELECT `link` FROM `news` WHERE `link` IN").
		WillReturnRows(sqlmock.NewRows([]string{"link"}))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `news`").
		WithArgs(sqlmock.AnyArg(), expectedTitle, "", int64(1), validLink, "", nil, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `news`").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	require.NoError(t, service.parseNewsFeed(model.NewsSource{Source: 1, URL: null.StringFrom("https://example.com/rss")}))
}

func TestNewsIsolatesPerItemInsertErrors(t *testing.T) {
	db, mock := newsTestDB(t)
	service := New(db)
	link1 := "https://example.com/bad"
	link2 := "https://example.com/good"
	feedXML := `<rss version="2.0"><channel><title>News</title>` +
		`<item><title>Bad</title><link>` + link1 + `</link></item>` +
		`<item><title>Good</title><link>` + link2 + `</link></item>` +
		`</channel></rss>`
	service.client = &http.Client{Transport: newsTestTransport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(feedXML)), Header: make(http.Header)}, nil
	})}

	mock.ExpectQuery("SELECT `link` FROM `news` WHERE `link` IN").
		WillReturnRows(sqlmock.NewRows([]string{"link"}))
	// First item fails.
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `news`").WillReturnError(errors.New("constraint failed"))
	mock.ExpectRollback()
	// Second item still gets inserted.
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `news`").
		WithArgs(sqlmock.AnyArg(), "Good", "", int64(1), link2, "", nil, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `news`").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := service.parseNewsFeed(model.NewsSource{Source: 1, URL: null.StringFrom("https://example.com/rss")})
	require.Error(t, err)
	require.Contains(t, err.Error(), link1)
}

func TestNewsAtomExtractsContentImageAndOmitsFileAssociationOnInsert(t *testing.T) {
	db, mock := newsTestDB(t)
	service := New(db)
	link := "https://www.tum.de/aktuelles/alle-meldungen/pressemitteilungen/details/test"
	imgURL := "https://www.tum.de/fileadmin/_processed_/c/d/example.jpg"
	atomXML := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>TUM - Aktuelles</title>
  <entry>
    <id>https://www.tum.de/news-1.atom</id>
    <title>TUM Research</title>
    <published>2026-11-19T15:13:00+01:00</published>
    <link rel="alternate" type="text/html" href="` + link + `"/>
    <summary>Summary text</summary>
    <content type="xhtml">
      <div xmlns="http://www.w3.org/1999/xhtml">
        <img src="` + imgURL + `" width="200" height="113"/>
        <p>Body</p>
      </div>
    </content>
  </entry>
</feed>`
	service.client = &http.Client{Transport: newsTestTransport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(atomXML)), Header: make(http.Header)}, nil
	})}

	mock.ExpectQuery("SELECT `link` FROM `news` WHERE `link` IN").
		WithArgs(link).
		WillReturnRows(sqlmock.NewRows([]string{"link"}))
	// saveImage finds existing file row.
	mock.ExpectQuery("SELECT .* FROM `files` WHERE name = \\?").
		WillReturnRows(sqlmock.NewRows([]string{"file", "name", "path", "url", "downloaded"}).
			AddRow(int64(42), "hash.jpg", NewsImageDirectory, imgURL, false))
	// Insert news directly without re-upserting `files` (File association omitted).
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `news` .* ON DUPLICATE KEY UPDATE").
		WithArgs(sqlmock.AnyArg(), "TUM Research", "Summary text", int64(1), link, imgURL, int64(42), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `news`").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	require.NoError(t, service.parseNewsFeed(model.NewsSource{Source: 1, URL: null.StringFrom("https://www.tum.de/news.rss")}))
}

func TestNewsDateRejectsOutOfRangeMySQLYear(t *testing.T) {
	tooOld := time.Date(999, 12, 31, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 11, 19, 12, 0, 0, 0, time.UTC)
	require.Equal(t, updated, newsDate(&gofeed.Item{PublishedParsed: &tooOld, UpdatedParsed: &updated}))
	require.WithinDuration(t, time.Now(), newsDate(&gofeed.Item{PublishedParsed: &tooOld}), time.Second)
}

func TestNewsParsesHeilbronnAndCommunityHTMLPages(t *testing.T) {
	chnHTML := `<!DOCTYPE html><html><head><title>Our Events on Campus | TUM Campus Heilbronn</title></head><body>
<div class="c-events-solr-list__items">
  <article class="c-events-solr-list__item">
    <div class="c-events-solr-list__item-image">
      <img src="/fileadmin/_processed_/7/5/event.webp" />
      <span class="c-date-badge"><time datetime="2026-11-06">06 Nov 26</time></span>
    </div>
    <div class="c-events-solr-list__item-content">
      <h3 class="c-events-solr-list__item-title">
        <a href="https://chn.tum.de/news-events/events/event-detail/event/exhibition-thomas-mann">Exhibition Thomas Mann</a>
      </h3>
      <p class="c-events-solr-list__item-teaser">A joint exhibition by TUM and ETH Zurich</p>
    </div>
  </article>
</div></body></html>`

	feed, err := parseNewsFeedContent("https://chn.tum.de/news-events/events", chnHTML)
	require.NoError(t, err)
	require.Len(t, feed.Items, 1)
	require.Equal(t, "Exhibition Thomas Mann", feed.Items[0].Title)
	require.Equal(t, "https://chn.tum.de/news-events/events/event-detail/event/exhibition-thomas-mann", feed.Items[0].Link)
	require.Equal(t, "A joint exhibition by TUM and ETH Zurich", feed.Items[0].Description)
	require.NotNil(t, feed.Items[0].Image)
	require.Equal(t, "https://chn.tum.de/fileadmin/_processed_/7/5/event.webp", feed.Items[0].Image.URL)
	require.NotNil(t, feed.Items[0].PublishedParsed)
	require.Equal(t, time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC), *feed.Items[0].PublishedParsed)

	commHTML := `<!DOCTYPE html><html><head><title>Events - TUM Community</title></head><body>
<ol class="se-events">
  <li class="se-event" data-start="2027-03-22T16:30:00+01:00">
    <article>
      <p class="se-event__date"><time datetime="2027-03-22T16:30:00+01:00">Mon, 22 Mar 2027</time></p>
      <h2 class="se-event__title"><a href="/en/event/5702/">Intercultural Communication at Work</a></h2>
      <p class="se-event__teaser">Practical workshop for international careers.</p>
      <p class="se-event__meta">Online · Webinar</p>
    </article>
  </li>
</ol></body></html>`

	commFeed, err := parseNewsFeedContent("https://www.community.tum.de/en/events/", commHTML)
	require.NoError(t, err)
	require.Len(t, commFeed.Items, 1)
	require.Equal(t, "Intercultural Communication at Work", commFeed.Items[0].Title)
	require.Equal(t, "https://www.community.tum.de/en/event/5702/", commFeed.Items[0].Link)
	require.Equal(t, "Online · Webinar — Practical workshop for international careers.", commFeed.Items[0].Description)
	require.NotNil(t, commFeed.Items[0].PublishedParsed)
}
