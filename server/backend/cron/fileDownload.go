package cron

import (
	"bytes"
	"errors"
	"image"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// fileDownloadCron Downloads all files that are not marked as finished in the database.
func (c *CronService) fileDownloadCron() error {
	var files []model.Files
	err := c.db.Find(&files, "downloaded = 0").Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error("Could not get files from database")
		return err
	}
	for _, file := range files {
		if !file.URL.Valid {
			log.WithField("fileId", file.File).Info("skipping file without url")
			continue
		}
		log.WithField("url", file.URL.String).Info("downloading file")
		if err := downloadFile(&file); err != nil {
			log.WithError(err).WithField("url", file.URL.String).Warn("Could not download file")
			continue
		}
		if err = c.db.Model(&model.Files{URL: file.URL}).Update("downloaded", true).Error; err != nil {
			log.WithError(err).Error("Could not set image to downloaded.")
			continue
		}
	}
	return nil
}

// downloadFile Downloads a file, marks it downloaded and resizes it if it's an image.
// url: download url of the file
// name: target name of the file
func downloadFile(file *model.Files) error {
	resp, err := http.Get(file.URL.String)
	if err != nil {
		log.WithError(err).WithField("url", file.URL.String).Warn("Could not download image")
		return err
	}
	// read body here because we can't exhaust the io.reader twice
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Warn("Unable to read http body")
		return err
	}

	// in our case resolves to /Storage/news/newspread/1234abc.jpg
	dstPath := path.Join(StorageDir, file.Path, file.Name)
	if strings.HasPrefix(mimetype.Detect(body).String(), "image/") {
		// save with resizing
		downloadedImg, _, err := image.Decode(bytes.NewReader(body))
		if err != nil {
			log.WithError(err).WithField("url", file.URL.String).Warn("Couldn't decode source image")
			return err
		}

		resizedImage := imaging.Resize(downloadedImg, 1280, 0, imaging.Lanczos)
		if err = imaging.Save(resizedImage, dstPath, imaging.JPEGQuality(75)); err != nil {
			log.WithError(err).WithField("url", file.URL.String).Warn("Could not save image file")
			return err
		}
	} else {
		// save without resizing image
		if err := os.WriteFile(dstPath, body, 0644); err != nil {
			log.WithError(err).Error("Can't save file to disk")
			return err
		}
	}
	return nil
}
