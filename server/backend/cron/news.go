package cron

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/getsentry/sentry-go"
	"github.com/guregu/null"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gographics/imagick.v2/imagick"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	ImageDirectory = "news/newspread/"
	NewspreadHook  = "newspread"
	ImpulsivHook   = "impulsivHook"
)

var ImageContentTypeRegex, _ = regexp.Compile("image/[a-z.]+")

// newsCron fetches news and saves them to the database
func (c *CronService) newsCron(cronjob model.Crontab) error {
	//check if source id provided for news job is not null
	if !cronjob.ID.Valid {
		log.Println("skipping news job, id of source is null")
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
		err := c.cleanOldNewsForSource(source.URL.String)
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

func (c *CronService) parseNewsFeed(source model.NewsSource) error {
	// initializing here so we don't have to re-init on every image download
	imagick.Initialize()
	defer imagick.Terminate()
	log.Printf("processing source %s\n", source.URL.String)
	feed, err := c.gf.ParseURL(source.URL.String)
	if err != nil {
		log.Printf("error parsing rss: %v\n", err)
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
		if !skipNews(existingNewsLinksForSource, item.Link) {
			// execute special actions for some sources:
			if source.Hook.Valid {
				switch source.Hook.String {
				case NewspreadHook:
					c.newspreadHook(item)
				case ImpulsivHook:
					c.impulsivHook(item)
				}
			}
			// pick the first enclosure that is an image (if any)
			var pickedEnclosure *gofeed.Enclosure
			var enclosureUrl = null.String{NullString: sql.NullString{Valid: false}}
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
				file = c.downloadAndSaveImage(pickedEnclosure.URL)
				enclosureUrl = null.String{NullString: sql.NullString{String: pickedEnclosure.URL, Valid: true}}
			}
			newsItem := model.News{
				Date:        *item.PublishedParsed,
				Created:     time.Now(),
				Title:       item.Title,
				Description: item.Description,
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

// downloadAndSaveImage
// Downloads an image from the url, converts it to 720p width to and saves it to the database.
// Returns id of file in db, null int otherwise.
func (c *CronService) downloadAndSaveImage(url string) null.Int {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Could not download image file %v", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	fileExtension := ""
	switch resp.Header.Get("Content-type") {
	case "image/jpeg":
		fileExtension = ".jpg"
	case "image/png":
		fileExtension = ".png"
	default:
		log.Printf("unsuported content type for image: %s", resp.Header.Get("Content-type"))
		return null.Int{}
	}
	hash := md5.Sum([]byte(url))
	fileName := fmt.Sprintf("%x%s", hash, fileExtension)
	createdFile, err := os.Create(fmt.Sprintf("%s%s%s", STORAGE_DIR, ImageDirectory, fileName))
	if err != nil {
		log.Printf("couldn't create image file: %v\n", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	_, err = io.Copy(createdFile, resp.Body)
	if err != nil {
		log.Printf("couldn't save image file to disk: %v", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	err = createdFile.Close()
	if err != nil {
		log.Printf("couldn't close image file: %v\n", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err = mw.ReadImage(fmt.Sprintf("%s%s%s", STORAGE_DIR, ImageDirectory, fileName))
	if err != nil {
		log.Printf("couldn't open file with imagemagick: %v", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	newHeight := uint((float32(mw.GetImageHeight()) / float32(mw.GetImageWidth())) * 720.0)
	err = mw.ResizeImage(720, newHeight, imagick.FILTER_CATROM, 1)
	if err != nil {
		log.Printf("couldn't resize image: %v", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	err = mw.SetCompressionQuality(75)
	if err != nil {
		log.Printf("couldn't compress image: %v", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	err = mw.WriteImage(fmt.Sprintf("%s%s%s", STORAGE_DIR, ImageDirectory, fileName))
	if err != nil {
		log.Printf("couldn't save compressed image: %v", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	fileForDB := model.Files{Name: fileName, Path: ImageDirectory}
	res := c.db.Create(&fileForDB)
	if res.Error != nil {
		log.Printf("couldn't store file to database: %v", res.Error)
		sentry.CaptureException(res.Error)
	}
	return null.Int{NullInt64: sql.NullInt64{Int64: int64(fileForDB.File), Valid: true}}
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

func (c *CronService) cleanOldNewsForSource(source string) error {
	log.Printf("Truncating old entries for source %s\n", source)
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

}

func (c *CronService) impulsivHook(item *gofeed.Item) {
	// Convert titles such as "123" to "Impulsiv - Ausgabe 123"
	if match, err := regexp.MatchString("[0-9]+", item.Title); err == nil && match {
		item.Title = fmt.Sprintf("Impulsiv - Ausgabe %s", item.Title)
	} else {
		// convert titles such as "Lösungen zur Ausgabe 137" to "Impulsiv - Lösungen zur Ausgabe 137"
		item.Title = fmt.Sprintf("Impulsiv - %s", item.Title)
	}
}
