package cron

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/TUM-Dev/Campus-Backend/server/backend/cron/student_club_parsers"
	"gorm.io/gorm"

	"github.com/guregu/null"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

const (
	StudentClubImageDirectory = "student_club/"
)

func (c *CronService) studentClubCron() error {
	body, err := student_club_parsers.DownloadHtml("https://www.sv.tum.de/sv/hochschulgruppen/")
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.WithError(err).Error("Error while closing body")
		}
	}(body)
	if err != nil {
		return err
	}
	scrapedClubs, scrapedCollections, err := student_club_parsers.ParseStudentClubs(body)
	if err != nil {
		return err
	}

	// save the result of the previous steps (🎉)
	if err := c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("1 = 1").Delete(&model.StudentClub{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1 = 1").Delete(&model.StudentClubCollection{}).Error; err != nil {
			return err
		}
		for _, scrapedCollection := range scrapedCollections {
			collection := model.StudentClubCollection{
				ID:          scrapedCollection.Name,
				Description: scrapedCollection.Description,
			}
			if err := tx.Create(&collection).Error; err != nil {
				return err
			}
		}
		for _, scrapedClub := range scrapedClubs {
			club := model.StudentClub{
				Name:                    scrapedClub.Name,
				Description:             scrapedClub.Description,
				LinkUrl:                 scrapedClub.LinkUrl,
				StudentClubCollectionID: scrapedClub.Collection,
			}
			if scrapedClub.ImageUrl.Valid {
				file, err := saveImageTo(tx, scrapedClub.ImageUrl.String, StudentClubImageDirectory)
				if err != nil {
					return err
				}
				// assign the file_id to make sure the id is assigned
				club.Image = file
				club.ImageID = null.IntFrom(file.File)
			}
			if err := tx.Create(&club).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.WithError(err).Error("error while creating movie")
	}
	return nil
}

// saveImage saves an image to the database, so it can be downloaded by another cronjob and returns the file
func saveImageTo(tx *gorm.DB, url string, path string) (*model.File, error) {
	seps := strings.SplitAfter(url, ".")
	fileExtension := seps[len(seps)-1]
	targetFileName := fmt.Sprintf("%x.%s", md5.Sum([]byte(url)), fileExtension)
	var file model.File
	// path intentionally omitted in query to allow for deduplication
	if err := tx.First(&file, "name = ?", targetFileName).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).WithField("targetFileName", targetFileName).Error("Couldn't query database for file")
		return nil, err
	} else if err == nil {
		return &file, nil
	}

	// does not exist, store in database
	file = model.File{
		Name:       targetFileName,
		Path:       path,
		URL:        null.StringFrom(url),
		Downloaded: null.BoolFrom(false),
	}
	if err := tx.Create(&file).Error; err != nil {
		log.WithError(err).Error("Could not store new file to database")
		return nil, err
	}
	return &file, nil
}
