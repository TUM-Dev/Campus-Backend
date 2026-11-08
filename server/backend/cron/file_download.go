package cron

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Bound each pass so a backlog cannot monopolize the cron scheduler. Prioritize
// files with fewer attempts so broken upstream URLs do not starve new images.
const fileDownloadBatchSize = 20

func (c *CronService) fileDownloadCron() error {
	var files []model.File
	if err := c.db.Where("downloaded = 0 AND url IS NOT NULL").
		Order("downloads ASC").Order("file ASC").Limit(fileDownloadBatchSize).
		Find(&files).Error; err != nil {
		return err
	}

	var downloadErrors []error
	for _, file := range files {
		if err := c.downloadPendingFile(file); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"file": file.File, "url": file.URL.String,
				"dstPath": path.Join(StorageDir, file.Path, file.Name),
			}).Warn("Could not download pending file")
			downloadErrors = append(downloadErrors, fmt.Errorf("file %d: %w", file.File, err))
		}
	}
	return errors.Join(downloadErrors...)
}

func (c *CronService) downloadPendingFile(file model.File) error {
	// Persist progress independently of other files; never hold a database
	// transaction open while waiting for upstream HTTP requests.
	if err := c.db.Model(&file).UpdateColumn("downloads", gorm.Expr("downloads + 1")).Error; err != nil {
		return err
	}
	dstPath := path.Join(StorageDir, file.Path, file.Name)
	if err := ensureFileDoesNotExist(dstPath); err != nil {
		return err
	}
	var downloadErr error
	if strings.HasPrefix(file.Path, "news/") {
		downloadErr = downloadFileWithClient(file.URL.String, dstPath, webfetch.PublicClient)
	} else {
		downloadErr = downloadFile(file.URL.String, dstPath)
	}
	if downloadErr != nil {
		return downloadErr
	}
	mime, err := mimetype.DetectFile(dstPath)
	if err != nil {
		return err
	}
	if !mime.Is("image/svg+xml") && strings.HasPrefix(mime.String(), "image/") {
		if err := resizeImage(dstPath); err != nil {
			return err
		}
	}
	return c.db.Model(&file).UpdateColumn("downloaded", true).Error
}

// ensureFileDoesNotExist makes sure that the file does not exist, but the directory in which it should be does
func ensureFileDoesNotExist(dstPath string) error {
	if _, err := os.Stat(dstPath); err == nil {
		// file already exists
		return os.Remove(dstPath)
	}
	return os.MkdirAll(path.Dir(dstPath), 0755)
}

// resizeImage resizes the image if it's an image to 1280px width keeping the aspect ratio
func resizeImage(dstPath string) error {
	img, err := imaging.Open(dstPath)
	if err != nil {
		return err
	}
	resizedImage := imaging.Resize(img, 1280, 0, imaging.Lanczos)
	return imaging.Save(resizedImage, dstPath, imaging.JPEGQuality(75))
}

// downloadFile Downloads a file from the given url and saves it to the given path
func downloadFile(url string, dstPath string) error {
	return downloadFileWithClient(url, dstPath, webfetch.Client)
}

func downloadFileWithClient(url string, dstPath string, client *http.Client) error {
	fields := log.Fields{"url": url, "dstPath": dstPath}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return fmt.Errorf("upstream returned %s", resp.Status)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.WithError(err).WithFields(fields).Error("Error while closing body")
		}
	}(resp.Body)

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
