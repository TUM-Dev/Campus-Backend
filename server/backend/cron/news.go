package cron

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/disintegration/imaging"
	"github.com/getsentry/sentry-go"
	"github.com/guregu/null"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"image"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	ImageDirectory   = "news/newspread/"
	NewspreadHook    = "newspread"
	ImpulsivHook     = "impulsivHook"
	MAX_IMAGE_RETRYS = 3
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
	if err := c.db.Table("`news`").Select("`link`").Where("`src` = ?", source.Source).Scan(&existingNewsLinksForSource).Error;
		err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
			var file = null.Int{NullInt64: sql.NullInt64{Valid: false}}
			if pickedEnclosure != nil {
				file, err = c.getDatabaseIdForImageAndDownload(pickedEnclosure.URL)
				if err != nil {
					continue // don't store this entry if file download failed.
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
				File:        file,
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

// getDatabaseIdForImageAndDownload
// Returns a file id that was saved to the database if it already exists. Otherwise creates it and triggers download of the image
func (c *CronService) getDatabaseIdForImageAndDownload(url string) (null.Int, error) {
	targetFileName := fmt.Sprintf("%x.jpg", md5.Sum([]byte(url)))
	var fileId null.Int
	if err := c.db.Model(model.Files{}).Where("name = ?", targetFileName).Select("file").Scan(&fileId).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Couldn't query database for file: %v", err)
		return null.Int{}, err
	}
	if fileId.Valid { // file already in database -> return for current news.
		return fileId, nil
	}

	// otherwise store in database:
	file := model.Files{Name: targetFileName, Path: STORAGE_DIR + ImageDirectory}
	if err := c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&file).Error; err != nil {
			log.Printf("Could not store new file to database: %v", err)
			return err
		}
		return nil
	}); err != nil {
		return null.Int{}, err
	}

	go c.downloadFile(url, targetFileName, MAX_IMAGE_RETRYS)
	return null.Int{NullInt64: sql.NullInt64{Valid: true, Int64: int64(file.File)}}, nil
}

// downloadFile tries downloading a file errorCounter times. After the errorCounters failure the corresponding entry is deleted from the database.
//
// url: download url of the file
// name: target name of the file
// errorCounter: recursively decremented counter for errors.
func (c *CronService) downloadFile(url string, name string, errorCounter int) {
	log.Printf("downloading file %s", url)
	if errorCounter == 0 {
		// unable to download image. Delete from database.
		var newsWithBadFile []model.News
		if err := c.db.Model(&model.News{}).Joins("JOIN files on files.file = news.file WHERE files.name = ?", name).Scan(&newsWithBadFile).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Println("Could not get news with bad files: %v", err)
			sentry.CaptureException(err)
			return
		}
		for i := range newsWithBadFile {
			newsWithBadFile[i].File = null.Int{NullInt64: sql.NullInt64{Valid: false}}
			if err := c.db.Save(&newsWithBadFile).Error; err != nil {
				log.Printf("Couldn't update news entry with bad file: %v", err)
				sentry.CaptureException(err)
			}
		}
		if err := c.db.Delete(&model.Files{}, "name = ?", name).Error; err != nil {
			log.Printf("Couldn't delete bad file entry from database: %v", err)
			sentry.CaptureException(err)
		}
		return
	}
	// wait before next retry n*10 seconds. e.g. first time: 0s,  1. retry -> 10s 2. -> 20, 3. -> 30
	time.Sleep(time.Duration((MAX_IMAGE_RETRYS-errorCounter)*10) * time.Second)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Could not download image file. Remaining attempts: %d, error: %v", errorCounter, err)
		sentry.CaptureException(err)
		c.downloadFile(url, name, errorCounter-1)
		return
	}
	downloadedImg, _, err := image.Decode(resp.Body)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Couldn't decode source image. Remaining attempts: %d, error: %v", errorCounter, err)
		c.downloadFile(url, name, errorCounter-1)
		return
	}

	// in our case resolves to /Storage/news/newspread/1234abc.jpg
	dstFileName := fmt.Sprintf("%s%s%s", STORAGE_DIR, ImageDirectory, name)
	dstImage := imaging.Resize(downloadedImg, 1280, 0, imaging.Lanczos)
	err = imaging.Save(dstImage, dstFileName, imaging.JPEGQuality(75))
	if err != nil {
		log.Printf("Could not save image file. Remaining attempts: %d, error: %v", errorCounter, err)
		sentry.CaptureException(err)
		c.downloadFile(url, name, errorCounter-1)
		return
	}
}

//skipNews returns true if link is in existingLinks or link is invalid
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
		log.Printf("cleaned up %v old news", res.RowsAffected)
	} else {
		log.Printf("failed to clean up old news: %v\n", res.Error)
		sentry.CaptureException(res.Error)
		return res.Error
	}
	return nil
}

func (c *CronService) newspreadHook(item *gofeed.Item) {
	re := regexp.MustCompile("https://storage.googleapis.com/tum-newspread-de/assets/[a-z\\-0-9]+\\.jpeg")
	extractedImageSlice := re.FindAllString(item.Content, 1)
	extractedImageURL := ""
	if len(extractedImageSlice) != 0 {
		extractedImageURL = extractedImageSlice[0]
	}
	item.Enclosures = []*gofeed.Enclosure{{URL: extractedImageURL}}
	item.Link = extractedImageURL
	item.Description = ""
}

//impulsivHook Converts the title of impulsiv news to a human friendly one
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
