package cron

import (
	"bytes"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/server/model"
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

// fileDownloadCron Downloads all files that are not marked as finished in the database.
func (c *CronService) fileDownloadCron() error {
	var files []model.Files
	err := c.db.Find(&files, "downloaded = 0").Scan(&files).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	for i := range files {
		if files[i].URL.Valid {
			c.downloadFile(files[i])
		}
	}
	return nil
}

// downloadFile Downloads a file, marks it downloaded and resizes it if it's an image.
// url: download url of the file
// name: target name of the file
func (c *CronService) downloadFile(file model.Files) {
	if !file.URL.Valid {
		log.WithField("fileId", file.File).Info("skipping file without url")
	}
	url := file.URL.String
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
		dstFileName := fmt.Sprintf("%s%s", file.Path, file.Name)
		dstImage := imaging.Resize(downloadedImg, 1280, 0, imaging.Lanczos)
		err = imaging.Save(dstImage, STORAGE_DIR+dstFileName, imaging.JPEGQuality(75))
		if err != nil {
			log.WithError(err).WithField("url", url).Warn("Could not save image file")
			sentry.CaptureException(err)
			return
		}
	} else {
		// save without resizing image
		err = ioutil.WriteFile(fmt.Sprintf("%s%s", file.Path, file.Name), body, 0644)
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
