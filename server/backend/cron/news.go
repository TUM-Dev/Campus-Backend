package cron

import (
	"crypto/md5"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	NewsImageDirectory = "news/newspread/"
	NewspreadHook      = "newspread"
	ImpulsivHook       = "impulsivHook"
	//MAX_IMAGE_RETRYS = 3
)

var ImageContentTypeRegex, _ = regexp.Compile("image/[a-z.]+")

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
		log.WithError(err).Error("getting news source from database")
		return err
	}
	// skip sources with null url
	if source.URL.Valid {
		return c.parseNewsFeed(source)
	}
	return nil
}

// parseNewsFeed processes a single news feed, extracts titles, content etc and saves it to the database
func (c *CronService) parseNewsFeed(source model.NewsSource) error {
	log.WithField("url", source.URL.String).Trace("processing newsfeed")
	// gofeed initializes translators lazily and is not safe to share between
	// concurrent jobs. Each source gets a parser; connections share one client.
	parser := gofeed.NewParser()
	parser.Client = c.client
	if parser.Client == nil {
		parser.Client = webfetch.Client
	}
	feed, err := parser.ParseURL(source.URL.String)
	if err != nil {
		log.WithError(err).Error("parsing rss")
		return err
	}
	// get all news for this source so we only process new ones, using map for performance reasons
	existingNewsLinksForSource := make([]string, 0)
	if err := c.db.Table("`news`").
		Select("`link`").
		Find(&existingNewsLinksForSource, "`src` = ?", source.Source).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error("failed to fetch existing news")
		return err
	}
	var newNews []model.News
	activeLinks := make([]string, 0, len(feed.Items))
	for _, item := range feed.Items {
		if item.Link != "" {
			activeLinks = append(activeLinks, item.Link)
		}
		// execute special actions for some sources:
		if source.Hook.Valid {
			switch source.Hook.String {
			case NewspreadHook:
				c.newspreadHook(item)
			case ImpulsivHook:
				c.impulsivHook(item)
			}
		}

		if skipNews(existingNewsLinksForSource, item.Link) {
			continue
		}

		// pick the first enclosure that is an image (if any)
		var pickedEnclosure *gofeed.Enclosure
		for _, enclosure := range item.Enclosures {
			if strings.HasSuffix(enclosure.URL, "jpg") ||
				strings.HasSuffix(enclosure.URL, "jpeg") ||
				strings.HasSuffix(enclosure.URL, "png") ||
				ImageContentTypeRegex.MatchString(enclosure.Type) {
				pickedEnclosure = enclosure
				break
			}
		}
		var enclosureUrl = null.StringFrom("")
		var file *model.File
		var fileID null.Int
		if pickedEnclosure != nil {
			file, err = c.saveImage(pickedEnclosure.URL)
			if err != nil {
				log.WithError(err).WithField("url", pickedEnclosure.URL).Error("can't save news image")
			} else {
				fileID = null.IntFrom(file.File)
			}
			enclosureUrl = null.StringFrom(pickedEnclosure.URL)
		}

		newsItem := model.News{
			Date:         newsDate(item),
			Created:      time.Now(),
			Title:        item.Title,
			Description:  bluemonday.StrictPolicy().Sanitize(item.Description),
			NewsSourceID: source.Source,
			NewsSource:   source,
			Link:         item.Link,
			Image:        enclosureUrl,
			FileID:       fileID,
			File:         file,
		}
		newNews = append(newNews, newsItem)
	}
	if ammountOfNewNews := len(newNews); ammountOfNewNews != 0 {
		err = c.db.Omit("NewsSource").Create(&newNews).Error
		if err != nil {
			log.WithField("ammountOfNewNews", ammountOfNewNews).Error("Inserting new news failed")
		} else {
			log.WithField("ammountOfNewNews", ammountOfNewNews).Trace("Inserting new news")
		}
		if err != nil {
			return err
		}
	}
	return c.cleanOldNewsForSource(source.Source, activeLinks)
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

// skipNews returns true if link is in existingLinks or link is invalid
func skipNews(existingLinks []string, link string) bool {
	if link == "" {
		return true
	}
	for _, l := range existingLinks {
		if l == link {
			return true
		}
	}
	return false
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
	re := regexp.MustCompile(`(?:^|[\s"'(])` + `(https://storage\.googleapis\.com/tum-newspread-de/assets/[a-z\-0-9]+\.jpeg)` + `(?:$|[\s"')])`)
	extractedImageSlice := re.FindAllStringSubmatch(item.Content, 1)
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
	re := regexp.MustCompile("[0-9]+")
	match := re.FindAllString(item.Title, -1)
	if len(match) == 1 && match[0] == item.Title {
		item.Title = fmt.Sprintf("Impulsiv - Ausgabe %s", item.Title)
	} else {
		// convert titles such as "Lösungen zur Ausgabe 137" to "Impulsiv - Lösungen zur Ausgabe 137"
		item.Title = fmt.Sprintf("Impulsiv - %s", item.Title)
	}
}

// Atom feeds may only specify updated; some RSS sources omit dates entirely.
func newsDate(item *gofeed.Item) time.Time {
	if item.PublishedParsed != nil && !item.PublishedParsed.IsZero() {
		return *item.PublishedParsed
	}
	if item.UpdatedParsed != nil && !item.UpdatedParsed.IsZero() {
		return *item.UpdatedParsed
	}
	return time.Now()
}
