package cron

import (
	"errors"
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

// fileDownloadCron downloads all files that are not marked as finished in the database
func (c *CronService) fileDownloadCron() error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		var files []model.Files
		err := tx.Find(&files, "downloaded = 0 AND url IS NOT NULL").Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("Could not get files from database")
			return err
		}
		for _, file := range files {
			// in our case resolves to /Storage/news/newspread/1234abc.jpg
			dstPath := path.Join(StorageDir, file.Path, file.Name)
			fields := log.Fields{"url": file.URL.String, "dstPath": dstPath}
			log.WithFields(fields).Info("downloading file")

			if err = tx.Model(&model.Files{File: file.File}).Update("downloads", file.Downloads+1).Error; err != nil {
				log.WithError(err).WithFields(fields).Error("Could not set update the download-count")
				continue
			}

			log.WithFields(fields).Info("ensuring file does not exist")
			if err := ensureFileDoesNotExist(dstPath); err != nil {
				log.WithError(err).WithFields(fields).Warn("Could not ensure file does not exist")
				continue
			}
			log.WithFields(fields).Info("downloading file")
			if err := downloadFile(file.URL.String, dstPath); err != nil {
				log.WithError(err).WithFields(fields).Warn("Could not download file")
				continue
			}
			log.WithFields(fields).Info("maybeResizeImage")
			if err := maybeResizeImage(dstPath); err != nil {
				log.WithError(err).WithFields(fields).Warn("Could not resize image")
				continue
			}
			log.WithFields(fields).Info("updating downloaded")
			// everything went well => we can mark the file as downloaded
			if err = tx.Model(&model.Files{URL: file.URL}).Update("downloaded", true).Error; err != nil {
				log.WithError(err).WithFields(fields).Error("Could not set image to downloaded.")
				continue
			}
		}
		return nil
	})
}

// ensureFileDoesNotExist makes sure that the file does not exist, but the directory in which it should be does
func ensureFileDoesNotExist(dstPath string) error {
	if _, err := os.Stat(dstPath); err == nil {
		// file already exists
		return os.Remove(dstPath)
	}
	return os.MkdirAll(path.Dir(dstPath), 0755)
}

// maybeResizeImage resizes the image if it's an image to 1280px width keeping the aspect ratio
func maybeResizeImage(dstPath string) error {
	mime, err := mimetype.DetectFile(dstPath)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(mime.String(), "image/") {
		return nil
	}

	img, err := imaging.Open(dstPath)
	if err != nil {
		return err
	}
	resizedImage := imaging.Resize(img, 1280, 0, imaging.Lanczos)
	return imaging.Save(resizedImage, dstPath, imaging.JPEGQuality(75))
}

// downloadFile Downloads a file from the given url and saves it to the given path
func downloadFile(url string, dstPath string) error {
	fields := log.Fields{"url": url, "dstPath": dstPath}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.WithError(err).WithFields(fields).Error("Error while closing body")
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return err
	}

	// save the file to disk
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.WithError(err).WithFields(fields).Error("Error while closing file")
		}
	}(out)
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}
