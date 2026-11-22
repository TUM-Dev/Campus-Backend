package cron

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	NewsImageDirectory = "news/newspread/"
	NewspreadHook      = "newspread"
	ImpulsivHook       = "impulsivHook"

	// Production column limits on the `news` table (see migration baseline).
	maxNewsLinkRunes        = 190
	maxNewsTitleBytes       = 255
	maxNewsDescriptionBytes = 65535
)

var (
	ImageContentTypeRegex = regexp.MustCompile("image/[a-z.]+")
	newspreadImageRegex   = regexp.MustCompile(`(?:^|[\s"'(])` + `(https://storage\.googleapis\.com/tum-newspread-de/assets/[a-z\-0-9]+\.jpeg)` + `(?:$|[\s"')])`)
	impulsivNumberRegex   = regexp.MustCompile("[0-9]+")
)

// newsCron fetches news and saves them to the database
func (c *CronService) newsCron(cronjob *model.Crontab) error {
	//check if source id provided for news job is not null
	if !cronjob.ID.Valid {
		fields := log.Fields{
			"Cron":     cronjob.Cron,
			"Interval": cronjob.Interval,
			"LastRun":  cronjob.LastRun,
			"Type":     cronjob.Type,
			"ID":       cronjob.ID,
		}
		log.WithFields(fields).Warn("skipping news job, id of source is null")
		return nil
	}
	// get news source for cronjob
	var source model.NewsSource
	err := c.db.First(&source, cronjob.ID.Int64).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithFields(log.Fields{"cron": cronjob.Cron, "source": cronjob.ID.Int64}).
				Warn("skipping news job, news source does not exist")
			return nil
		}
		log.WithError(err).Error("getting news source from database")
		return err
	}
	// skip sources with null url or empty url
	if source.URL.Valid && strings.TrimSpace(source.URL.String) != "" {
		return c.parseNewsFeed(source)
	}
	return nil
}

// fetchNewsFeed fetches a URL and parses it either as an RSS/Atom/JSON feed or
// as a supported TUM HTML news/events listing (Heilbronn, TUM Community, SV).
func (c *CronService) fetchNewsFeed(feedURL string) (*gofeed.Feed, error) {
	client := c.client
	if client == nil {
		client = webfetch.Client
	}
	req, err := http.NewRequest(http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Gofeed/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	finalURL := feedURL
	if resp.Request != nil && resp.Request.URL != nil {
		finalURL = resp.Request.URL.String()
	}
	return parseNewsFeedContent(finalURL, string(body))
}

func parseNewsFeedContent(baseURL, body string) (*gofeed.Feed, error) {
	// gofeed initializes translators lazily and is not safe to share between
	// concurrent jobs. Each parse gets its own parser instance.
	feed, err := gofeed.NewParser().ParseString(body)
	if err == nil {
		return feed, nil
	}
	if htmlFeed, htmlErr := parseHTMLNewsFeed(baseURL, body); htmlErr == nil && len(htmlFeed.Items) > 0 {
		return htmlFeed, nil
	}
	return nil, err
}

// parseHTMLNewsFeed extracts news and event items from TUM TYPO3 SolR pages
// (chn.tum.de newsroom and events), TUM Community event listings
// (community.tum.de), and TYPO3 news lists (sv.tum.de).
func parseHTMLNewsFeed(baseURL, body string) (*gofeed.Feed, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	base, _ := url.Parse(baseURL)
	resolveURL := func(raw string) string {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return ""
		}
		u, err := url.Parse(raw)
		if err != nil {
			return ""
		}
		if base != nil {
			return base.ResolveReference(u).String()
		}
		return u.String()
	}

	feed := &gofeed.Feed{
		Title: strings.TrimSpace(doc.Find("title").First().Text()),
		Link:  baseURL,
	}

	// 1. TUM Campus Heilbronn (chn.tum.de) SolR Events & Newsroom cards
	doc.Find("article.c-events-solr-list__item, article.c-news-solr-list__item").Each(func(_ int, s *goquery.Selection) {
		titleLink := s.Find("h3.c-events-solr-list__item-title a, h3.c-news-solr-list__item-title a").First()
		title := strings.TrimSpace(titleLink.Text())
		href, _ := titleLink.Attr("href")
		link := resolveURL(href)
		if title == "" || link == "" {
			return
		}
		desc := strings.TrimSpace(s.Find("p.c-events-solr-list__item-teaser, p.c-news-solr-list__item-teaser").First().Text())
		item := &gofeed.Item{
			Title:       title,
			Link:        link,
			Description: desc,
		}
		if dt, ok := s.Find("time[datetime]").First().Attr("datetime"); ok {
			if parsed, ok := parseFlexibleHTMLTime(dt); ok {
				item.PublishedParsed = &parsed
			}
		}
		if imgSrc, ok := s.Find("img[src]").First().Attr("src"); ok {
			if fullImg := resolveURL(imgSrc); isHTTPImageURL(fullImg) {
				item.Image = &gofeed.Image{URL: fullImg}
			}
		}
		feed.Items = append(feed.Items, item)
	})
	if len(feed.Items) > 0 {
		return feed, nil
	}

	// 2. TUM Community (community.tum.de/en/events/ or /veranstaltungen/)
	doc.Find("li.se-event").Each(func(_ int, s *goquery.Selection) {
		titleLink := s.Find(".se-event__title a, h2 a, h3 a").First()
		title := strings.TrimSpace(titleLink.Text())
		href, _ := titleLink.Attr("href")
		link := resolveURL(href)
		if title == "" || link == "" {
			return
		}
		teaser := strings.TrimSpace(s.Find(".se-event__teaser").First().Text())
		meta := strings.TrimSpace(s.Find(".se-event__meta").First().Text())
		desc := teaser
		if meta != "" {
			if desc != "" {
				desc = meta + " — " + desc
			} else {
				desc = meta
			}
		}
		item := &gofeed.Item{
			Title:       title,
			Link:        link,
			Description: desc,
		}
		dt, ok := s.Attr("data-start")
		if !ok || strings.TrimSpace(dt) == "" {
			dt, _ = s.Find("time[datetime]").First().Attr("datetime")
		}
		if parsed, ok := parseFlexibleHTMLTime(dt); ok {
			item.PublishedParsed = &parsed
		}
		feed.Items = append(feed.Items, item)
	})
	if len(feed.Items) > 0 {
		return feed, nil
	}

	// 3. Generic TYPO3 news list view (e.g. sv.tum.de)
	doc.Find("div.news-list-view div.article[data-testid='news-item']").Each(func(_ int, s *goquery.Selection) {
		titleLink := s.Find(".news-header-headline a, h3 a").First()
		title := strings.TrimSpace(titleLink.Text())
		href, _ := titleLink.Attr("href")
		link := resolveURL(href)
		if title == "" || link == "" {
			return
		}
		teaserSel := s.Find(".teaser-text p").First().Clone()
		teaserSel.Find(".news-list-morelink").Remove()
		desc := strings.TrimSpace(teaserSel.Text())
		item := &gofeed.Item{
			Title:       title,
			Link:        link,
			Description: desc,
		}
		if dt, ok := s.Find("time[datetime]").First().Attr("datetime"); ok {
			if parsed, ok := parseFlexibleHTMLTime(dt); ok {
				item.PublishedParsed = &parsed
			}
		}
		if imgSrc, ok := s.Find("img[src]").First().Attr("src"); ok {
			if fullImg := resolveURL(imgSrc); isHTTPImageURL(fullImg) {
				item.Image = &gofeed.Image{URL: fullImg}
			}
		}
		feed.Items = append(feed.Items, item)
	})

	if len(feed.Items) == 0 {
		return nil, errors.New("no supported news or event items found in HTML")
	}
	return feed, nil
}

func parseFlexibleHTMLTime(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006.01.02",
		"02.01.2006",
	} {
		if t, err := time.Parse(layout, raw); err == nil && isValidMySQLDateTime(&t) {
			return t, true
		}
	}
	return time.Time{}, false
}

// parseNewsFeed processes a single news feed, extracts titles, content etc and saves it to the database
func (c *CronService) parseNewsFeed(source model.NewsSource) error {
	log.WithFields(log.Fields{"source": source.Source, "url": source.URL.String}).Trace("processing newsfeed")
	feed, err := c.fetchNewsFeed(source.URL.String)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"source": source.Source, "url": source.URL.String}).Error("parsing rss")
		return err
	}

	activeLinks := make([]string, 0, len(feed.Items))
	activeLinkSet := make(map[string]struct{}, len(feed.Items))
	for _, item := range feed.Items {
		if item.Link != "" {
			if _, seen := activeLinkSet[item.Link]; !seen {
				activeLinkSet[item.Link] = struct{}{}
				activeLinks = append(activeLinks, item.Link)
			}
		}
	}

	// `news.link` has a global UNIQUE index across all sources, so look up all
	// matching links regardless of source to prevent cross-source duplicate errors.
	existingLinks := make(map[string]struct{}, len(activeLinks))
	if len(activeLinks) > 0 {
		var found []string
		if err := c.db.Table("`news`").
			Select("`link`").
			Where("`link` IN ?", activeLinks).
			Find(&found).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).WithField("source", source.Source).Error("failed to fetch existing news")
			return err
		}
		for _, link := range found {
			existingLinks[link] = struct{}{}
		}
	}

	sanitizer := bluemonday.StrictPolicy()
	var insertErrors []error
	insertedCount := 0
	for _, item := range feed.Items {
		// execute special actions for some sources:
		if source.Hook.Valid {
			switch source.Hook.String {
			case NewspreadHook:
				c.newspreadHook(item)
			case ImpulsivHook:
				c.impulsivHook(item)
			}
		}

		if item.Link == "" {
			continue
		}
		if _, exists := existingLinks[item.Link]; exists {
			continue
		}
		// Track links seen in this feed so duplicate entries within one feed don't fail.
		existingLinks[item.Link] = struct{}{}

		if linkRunes := utf8.RuneCountInString(item.Link); linkRunes > maxNewsLinkRunes {
			log.WithFields(log.Fields{
				"source": source.Source,
				"link":   item.Link,
				"length": linkRunes,
			}).Warn("Skipping news item whose link exceeds column limit")
			continue
		}

		var imageURL = null.StringFrom("")
		var file *model.File
		var fileID null.Int
		if pickedImageURL := pickNewsImageURL(item, source.Hook.String); pickedImageURL != "" {
			file, err = c.saveImage(pickedImageURL)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{"source": source.Source, "url": pickedImageURL}).Error("can't save news image")
			} else {
				fileID = null.IntFrom(file.File)
			}
			imageURL = null.StringFrom(pickedImageURL)
		}

		rawDesc := item.Description
		if strings.TrimSpace(rawDesc) == "" && source.Hook.String != NewspreadHook {
			rawDesc = extractHTMLDescriptionFallback(item.Content, item.Title)
		}

		newsItem := model.News{
			Date:         newsDate(item),
			Created:      time.Now(),
			Title:        truncateUTF8Bytes(item.Title, maxNewsTitleBytes),
			Description:  truncateUTF8Bytes(sanitizer.Sanitize(rawDesc), maxNewsDescriptionBytes),
			NewsSourceID: source.Source,
			NewsSource:   source,
			Link:         item.Link,
			Image:        imageURL,
			FileID:       fileID,
			File:         file,
		}

		// Insert per-item with ON DUPLICATE KEY DO NOTHING so a concurrent source
		// race or one bad article never aborts the rest of the feed.
		if err := c.db.Omit("NewsSource", "File").
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&newsItem).Error; err != nil {
			log.WithError(err).WithFields(log.Fields{
				"source": source.Source,
				"link":   newsItem.Link,
			}).Error("Inserting new news failed")
			insertErrors = append(insertErrors, fmt.Errorf("news %s: %w", newsItem.Link, err))
			continue
		}
		insertedCount++
	}
	if insertedCount > 0 {
		log.WithFields(log.Fields{"source": source.Source, "count": insertedCount}).Trace("Inserted new news")
	}

	cleanErr := c.cleanOldNewsForSource(source.Source, activeLinks)
	return errors.Join(append(insertErrors, cleanErr)...)
}

// extractHTMLDescriptionFallback extracts paragraph text when an Atom entry
// provides <content> without a <summary> (such as sv.tum.de's Atom feeds).
func extractHTMLDescriptionFallback(content, title string) string {
	if strings.TrimSpace(content) == "" {
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return content
	}
	doc.Find("h1, h2, h3").Each(func(_ int, s *goquery.Selection) {
		if strings.TrimSpace(s.Text()) == strings.TrimSpace(title) {
			s.Remove()
		}
	})
	return strings.TrimSpace(doc.Text())
}

// saveImage saves an image to the database, so it can be downloaded by another cronjob and returns its id
func (c *CronService) saveImage(url string) (*model.File, error) {
	targetFileName := fmt.Sprintf("%x.jpg", md5.Sum([]byte(url)))
	var file model.File
	// path intentionally omitted in query to allow for deduplication
	if err := c.db.First(&file, "name = ?", targetFileName).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).WithField("targetFileName", targetFileName).Error("Couldn't query database for file")
		return nil, err
	} else if err == nil {
		return &file, nil
	}

	// does not exist, store in database
	file = model.File{
		Name:       targetFileName,
		Path:       NewsImageDirectory,
		URL:        null.StringFrom(url),
		Downloaded: null.BoolFrom(false),
	}
	if err := c.db.Create(&file).Error; err != nil {
		log.WithError(err).Error("Could not store new file to database")
		return nil, err
	}
	return &file, nil
}

func (c *CronService) cleanOldNewsForSource(source int64, activeLinks []string) error {
	log.WithField("source", source).Trace("Truncating old entries")
	oneYearAgo := time.Now().Add(time.Hour * 24 * 365 * -1)
	query := c.db.Where("`src` = ? AND `created` < ?", source, oneYearAgo)
	if len(activeLinks) > 0 {
		query = query.Where("`link` NOT IN ?", activeLinks)
	}
	if res := query.Delete(&model.News{}); res.Error == nil {
		log.WithField("RowsAffected", res.RowsAffected).Info("cleaned up old news")
	} else {
		log.WithError(res.Error).Error("failed to clean up old news")
		return res.Error
	}
	return nil
}

// newspreadHook extracts image urls from the body if the feed because such entries are a bit weird
func (c *CronService) newspreadHook(item *gofeed.Item) {
	extractedImageSlice := newspreadImageRegex.FindAllStringSubmatch(item.Content, 1)
	extractedImageURL := ""
	if len(extractedImageSlice) != 0 && len(extractedImageSlice[0]) > 1 {
		extractedImageURL = extractedImageSlice[0][1]
	}
	item.Enclosures = []*gofeed.Enclosure{{URL: extractedImageURL}}
	item.Description = ""
}

// impulsivHook Converts the title of impulsiv news to a human friendly one
func (c *CronService) impulsivHook(item *gofeed.Item) {
	// Convert titles such as "123" to "Impulsiv - Ausgabe 123"
	match := impulsivNumberRegex.FindAllString(item.Title, -1)
	if len(match) == 1 && match[0] == item.Title {
		item.Title = fmt.Sprintf("Impulsiv - Ausgabe %s", item.Title)
	} else {
		// convert titles such as "Lösungen zur Ausgabe 137" to "Impulsiv - Lösungen zur Ausgabe 137"
		item.Title = fmt.Sprintf("Impulsiv - %s", item.Title)
	}
}

func isValidMySQLDateTime(t *time.Time) bool {
	if t == nil || t.IsZero() {
		return false
	}
	y := t.Year()
	return y >= 1000 && y <= 9999
}

// Atom feeds may only specify updated; some RSS sources omit dates entirely.
func newsDate(item *gofeed.Item) time.Time {
	if isValidMySQLDateTime(item.PublishedParsed) {
		return *item.PublishedParsed
	}
	if isValidMySQLDateTime(item.UpdatedParsed) {
		return *item.UpdatedParsed
	}
	return time.Now()
}

func pickNewsImageURL(item *gofeed.Item, hook string) string {
	for _, enclosure := range item.Enclosures {
		if enclosure == nil || enclosure.URL == "" {
			continue
		}
		if hasRasterImageExtension(enclosure.URL) || ImageContentTypeRegex.MatchString(enclosure.Type) {
			return enclosure.URL
		}
	}
	if hook == NewspreadHook {
		return ""
	}
	if item.Image != nil && isHTTPImageURL(item.Image.URL) {
		return item.Image.URL
	}
	for _, snippet := range []string{item.Content, item.Description} {
		if imgURL := extractFirstHTMLImageURL(snippet); imgURL != "" {
			return imgURL
		}
	}
	return ""
}

func hasRasterImageExtension(raw string) bool {
	p := raw
	if u, err := url.Parse(strings.TrimSpace(raw)); err == nil && u.Path != "" {
		p = u.Path
	}
	switch strings.ToLower(path.Ext(p)) {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	default:
		return false
	}
}

func isHTTPImageURL(raw string) bool {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}
	// saveImage stores downloads as .jpg; skip SVGs from HTML fallbacks.
	return strings.ToLower(path.Ext(u.Path)) != ".svg"
}

func extractFirstHTMLImageURL(snippet string) string {
	if !strings.Contains(strings.ToLower(snippet), "<img") {
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(snippet))
	if err != nil {
		return ""
	}
	var found string
	doc.Find("img").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if src, ok := s.Attr("src"); ok && isHTTPImageURL(src) {
			found = strings.TrimSpace(src)
			return false
		}
		return true
	})
	return found
}

// truncateUTF8Bytes trims s to at most maxBytes without splitting a UTF-8 rune.
func truncateUTF8Bytes(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	if maxBytes <= 0 {
		return ""
	}
	cut := maxBytes
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return s[:cut]
}
