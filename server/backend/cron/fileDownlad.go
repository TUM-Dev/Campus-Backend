package cron

import (
	"bytes"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"image"
	"io/ioutil"
	"net/http"
	"strings"
)

//fileDownloadCron Downloads all files that are not marked as finished in the database.
func (c *CronService) fileDownloadCron() error {
	var files []model.Files
	err := c.db.Find(&files, "downloaded = 0").Scan(&files).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	for i := range files {
		if files[i].URL.Valid {
			c.downloadFile(files[i].URL.String, files[i].Name)
		}
	}
	return nil
}

// downloadFile tries downloading a file errorCounter times. After the errorCounters failure the corresponding entry is deleted from the database.
// url: download url of the file
// name: target name of the file
func (c *CronService) downloadFile(url string, name string) {
	log.WithField("url", url).Info("downloading file")
	resp, err := http.Get(url)
	if err != nil {
		log.WithError(err).WithField("url", url).Warn("Could not download image")
		sentry.CaptureException(err)
		return
	}
	// read body here because we can't exhaust the io.reader twice
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Warn("Unable to read http body")
		return
	}

	// resize if file is image
	mime := mimetype.Detect(body)
	if strings.HasPrefix(mime.String(), "image/") {
		downloadedImg, _, err := image.Decode(bytes.NewReader(body))
		if err != nil {
			log.WithError(err).WithField("url", url).Warn("Couldn't decode source image")
			sentry.CaptureException(err)
			return
		}

		// in our case resolves to /Storage/news/newspread/1234abc.jpg
		dstFileName := fmt.Sprintf("%s%s%s", STORAGE_DIR, ImageDirectory, name)
		dstImage := imaging.Resize(downloadedImg, 1280, 0, imaging.Lanczos)
		err = imaging.Save(dstImage, dstFileName, imaging.JPEGQuality(75))
		if err != nil {
			log.WithError(err).WithField("url", url).Warn("Could not save image file")
			sentry.CaptureException(err)
			return
		}
	} else {
		// save without resizing image
		err = ioutil.WriteFile(fmt.Sprintf("%s%s", STORAGE_DIR, name), body, 0644)
		if err != nil {
			sentry.CaptureException(err)
			log.WithError(err).Error("Can't save file to disk")
			return
		}
	}
	err = c.db.Model(&model.Files{}).Where("url = ?", url).Update("downloaded", true).Error
	if err != nil {
		sentry.CaptureException(err)
		log.WithError(err).Error("Could not set image to downloaded.")
	}
}
