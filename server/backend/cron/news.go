package cron

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/getsentry/sentry-go"
	"github.com/guregu/null"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

const (
	ImageDirectory = "news/newspread/"
	NewspreadHook  = "newspread"
	ImpulsivHook   = "impulsivHook"
	//MAX_IMAGE_RETRYS = 3
)

var ImageContentTypeRegex, _ = regexp.Compile("image/[a-z.]+")

// newsCron fetches news and saves them to the database
func (c *CronService) newsCron(cronjob *model.Crontab) error {
	//check if source id provided for news job is not null
	if !cronjob.ID.Valid {
		cronjobJson, _ := json.Marshal(cronjob)
		log.Println("skipping news job, id of source is null, cronjob: %s", string(cronjobJson))
		return nil
	}
	// get news source for cronjob
	var source model.NewsSource
	err := c.db.Find(&source, cronjob.ID.Int64).Error
	if err != nil {
		log.Printf("error getting news source from database: %v", err)
		sentry.CaptureException(err)
		return err
	}
	// skip sources with null url
	if source.URL.Valid {
		// clean up news older than one year
		err := c.cleanOldNewsForSource(source.Source)
		if err != nil {
			return err
		}
		err = c.parseNewsFeed(source)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseNewsFeed processes a single news feed, extracts titles, content etc and saves it to the database
func (c *CronService) parseNewsFeed(source model.NewsSource) error {
	log.Printf("processing source %s", source.URL.String)
	feed, err := c.gf.ParseURL(source.URL.String)
	if err != nil {
		log.Printf("error parsing rss: %v", err)
		sentry.CaptureException(err)
		return err
	}
	// get all news for this source so we only process new ones, using map for performance reasons
	existingNewsLinksForSource := make([]string, 0)
	if err := c.db.Table("`news`").Select("`link`").Where("`src` = ?", source.Source).Scan(&existingNewsLinksForSource).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("failed to fetch existing news: %v", err)
		sentry.CaptureException(err)
		return err
	}
	var newNews []model.News
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

		if !skipNews(existingNewsLinksForSource, item.Link) {
			// pick the first enclosure that is an image (if any)
			var pickedEnclosure *gofeed.Enclosure
			var enclosureUrl = null.String{NullString: sql.NullString{Valid: true, String: ""}}
			for _, enclosure := range item.Enclosures {
				if strings.HasSuffix(enclosure.URL, "jpg") ||
					strings.HasSuffix(enclosure.URL, "jpeg") ||
					strings.HasSuffix(enclosure.URL, "png") ||
					ImageContentTypeRegex.MatchString(enclosure.Type) {
					pickedEnclosure = enclosure
					break
				}
			}
			var fileId = null.Int{NullInt64: sql.NullInt64{Valid: false}}
			if pickedEnclosure != nil {
				fileId, err = c.saveImage(pickedEnclosure.URL)
				if err != nil {
					log.WithError(err).Error("error saving image")
					sentry.CaptureException(err)
				}
				enclosureUrl = null.String{NullString: sql.NullString{String: pickedEnclosure.URL, Valid: true}}
			}
			bm := bluemonday.StrictPolicy()
			sanitizedDesc := bm.Sanitize(item.Description)

			newsItem := model.News{
				Date:        *item.PublishedParsed,
				Created:     time.Now(),
				Title:       item.Title,
				Description: sanitizedDesc,
				Src:         source.Source,
				Link:        item.Link,
				Image:       enclosureUrl,
				File:        fileId,
			}
			newNews = append(newNews, newsItem)
		}
	}
	if len(newNews) != 0 {
		log.Printf("Inserting %v new news", len(newNews))
		err = c.db.Save(&newNews).Error
		return err
	}
	return nil
}

// saveImage Saves an image to the database so it can be downloaded by another cronjob and returns it's id
func (c *CronService) saveImage(url string) (null.Int, error) {
	targetFileName := fmt.Sprintf("%x.jpg", md5.Sum([]byte(url)))
	var fileId null.Int
	if err := c.db.Model(model.Files{}).
		Where("name = ?", targetFileName).
		Select("file").Scan(&fileId).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Couldn't query database for file: %v", err)
		return null.Int{}, err
	}
	if fileId.Valid { // file already in database -> return for current news.
		return fileId, nil
	}

	// otherwise store in database:
	file := model.Files{
		Name:       targetFileName,
		Path:       ImageDirectory,
		URL:        sql.NullString{String: url, Valid: true},
		Downloaded: sql.NullBool{Bool: false, Valid: true},
	}
	err := c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&file).Error; err != nil {
			log.WithError(err).Error("Could not store new file to database")
			return err
		}
		return nil
	})
	if err != nil {
		return null.Int{}, err
	}
	// creating this int is annoying but i'm too afraid to use real ORM in the model
	return null.Int{NullInt64: sql.NullInt64{Int64: int64(file.File), Valid: true}}, nil
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

func (c *CronService) cleanOldNewsForSource(source int32) error {
	log.Printf("Truncating old entries for source %d\n", source)
	if res := c.db.Delete(&model.News{}, "`src` = ? AND `created` < ?", source, time.Now().Add(time.Hour*24*365*-1)); res.Error == nil {
		log.Infof("cleaned up %v old news", res.RowsAffected)
	} else {
		log.WithError(res.Error).Error("failed to clean up old news")
		sentry.CaptureException(res.Error)
		return res.Error
	}
	return nil
}

// newspreadHook extracts image urls from the body if the feed because such entries are a bit weird
func (c *CronService) newspreadHook(item *gofeed.Item) {
	re := regexp.MustCompile(`https://storage.googleapis.com/tum-newspread-de/assets/[a-z\-0-9]+\.jpeg`)
	extractedImageSlice := re.FindAllString(item.Content, 1)
	extractedImageURL := ""
	if len(extractedImageSlice) != 0 {
		extractedImageURL = extractedImageSlice[0]
	}
	item.Enclosures = []*gofeed.Enclosure{{URL: extractedImageURL}}
	item.Link = extractedImageURL
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
