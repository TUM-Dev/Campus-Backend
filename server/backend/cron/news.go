package cron

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/getsentry/sentry-go"
	"github.com/guregu/null"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"image"
	jpeg "image/jpeg"
	png "image/png"
	"net/http"
	"os"
	"regexp"
	"time"
)

const IMAGE_DIRECTORY = "/resources"

var IMAGE_CONTENT_TYPE_REGEX, _ = regexp.Compile("image/[a-z.]+")

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
		log.Printf("error getting news sources from database: %v", err)
		sentry.CaptureException(err)
	}
	// skip sources with null url
	if source.URL.Valid {
		// clean up news older than one year
		log.Printf("Truncating old entries for source %s\n", source.URL.String)
		if res := c.db.Delete(&model.News{}, "`src` = ? AND `created` < ?", source.URL.String, time.Now().Add(time.Hour*24*365*-1)); res.Error == nil {
			log.Printf("cleaned up %v old news", res.RowsAffected)
		} else {
			log.Printf("failed to clean up old news: %v\n", res.Error)
			sentry.CaptureException(res.Error)
		}
		log.Printf("processing source %s\n", source.URL.String)
		feed, err := c.gf.ParseURL(source.URL.String)
		if err != nil {
			log.Printf("error parsing rss: %v\n", err)
			sentry.CaptureException(err)
		} else {
			c.parseNewsFeed(feed, source)
		}
	}
	return nil
}

func (c *CronService) parseNewsFeed(feed *gofeed.Feed, source model.NewsSource) {
	log.Println(feed.Title)
	// get all news for this source so we only process new ones, using map for performance reasons
	existingNewsLinksForSource := make([]string, 0)
	db := c.db.Debug()
	if err := db.Table("`news`").Select("`link`").Where("`src` = ?", source.Source).Scan(&existingNewsLinksForSource).Error;
		err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("failed to fetch existing news: %v", err)
		sentry.CaptureException(err)
		return
	}
	var newNews []model.News
	for _, item := range feed.Items {
		if !skipNews(existingNewsLinksForSource, item.Link) {
			var pickedEnclosure *gofeed.Enclosure
			for _, enclosure := range item.Enclosures {
				if IMAGE_CONTENT_TYPE_REGEX.MatchString(enclosure.Type) {
					pickedEnclosure = enclosure
				}
			}
			var file = null.Int{NullInt64: sql.NullInt64{Valid: false}}
			if pickedEnclosure != nil {
				file = c.downloadAndSaveImage(pickedEnclosure.URL)
			}
			newsItem := model.News{
				Date:        time.Time{},
				Created:     time.Now(),
				Title:       feed.Title,
				Description: feed.Description,
				Src:         source.Source,
				Link:        feed.Link,
				Image:       null.String{},
				File:        file,
			}
			newNews = append(newNews, newsItem)
		}
	}
	db.Save(&newNews)
}

// downloadAndSaveImage Downloads an image from the url, converts it to and saves it to the database. returns id of file in db
func (c *CronService) downloadAndSaveImage(url string) null.Int {
	file, err := http.Get(url)
	if err != nil {
		log.Printf("Could not download image file %v", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	var decodedImage image.Image
	var fileExtension = ""
	switch file.Header.Get("Content-type") {
	case "image/jpeg":
		fileExtension = ".jpg"
		jpegImage, err := jpeg.Decode(file.Body)
		if err != nil {
			log.Printf("could not decode image: %err", err)
			sentry.CaptureException(err)
			return null.Int{}
		}
		decodedImage = jpegImage
	case "image/png":
		fileExtension = ".png"
		pngImage, err := png.Decode(file.Request.Body) //todo: setup compression
		if err != nil {
			log.Printf("could not decode image: %err", err)
			sentry.CaptureException(err)
			return null.Int{}
		}
		decodedImage = pngImage
	default:
		log.Printf("unsuported content type for image: %s", file.Header.Get("Content-type"))
		return null.Int{}
	}
	hash := md5.Sum([]byte(url))
	fileName := fmt.Sprintf("/%x%s", hash, fileExtension)
	createdFile, err := os.Create(fmt.Sprintf("%s%s", IMAGE_DIRECTORY, fileName))
	if err != nil {
		log.Printf("couldn't create image file: %v\n", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	buf := new(bytes.Buffer)
	switch file.Header.Get("Content-type") {
	case "image/jpeg":
		err := jpeg.Encode(buf, decodedImage, &jpeg.Options{Quality: 75})
		if err != nil {
			log.Printf("Couldn't encode image: %v", err)
			sentry.CaptureException(err)
			return null.Int{}
		}
	case "image/png":
		err := png.Encode(buf, decodedImage)
		if err != nil {
			log.Printf("Couldn't encode image: %v", err)
			sentry.CaptureException(err)
			return null.Int{}
		}
	}
	_, err = createdFile.Write(buf.Bytes())
	if err != nil {
		log.Printf("couldn't save image file to disk: %v", err)
		sentry.CaptureException(err)
		return null.Int{}
	}
	fileForDB := model.Files{Name: fileName, Path: IMAGE_DIRECTORY}
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
