package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron/movie_parsers"
	"time"

	"github.com/guregu/null"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

const (
	MovieImageDirectory = "movie/"
)

func (c *CronService) movieCron() error {
	log.Trace("parsing upcoming feed")
	channels, err := movie_parsers.GetUpcomingFeed()
	if err != nil {
		return err
	}
	for _, channel := range channels {
		for _, item := range channel.Items {
			logFields := log.Fields{"link": item.Link, "title": item.Title, "date": item.PubDate, "location": item.Location, "url": item.Enclosure.Url}
			var exists bool
			if err := c.db.Model(model.Kino{}).Select("count(*) > 0").Find(&exists, "link = ?", item.Link).Error; err != nil {
				log.WithError(err).WithFields(logFields).Error("Cound lot check if movie already exists")
				continue
			}
			if exists {
				log.WithFields(logFields).Trace("Movie already exists")
				continue
			}

			// data cleanup
			date, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				log.WithError(err).WithFields(logFields).Error("Couldn't check if movie already exists")
				continue
			}

			// populate extra data from omdb
			movieInformation, err := movie_parsers.GetTuFilmWebsiteInformation(item.Link)
			if err != nil {
				log.WithFields(logFields).WithError(err).Error("error while finding imdb id")
				continue
			}
			var omdbMovie *movie_parsers.OmdbResults
			if movieInformation.ImdbID != "" {
				omdbMovie, err = movie_parsers.GetOmdbMovie(movieInformation.ImdbID)
				if err != nil {
					log.WithFields(logFields).WithError(err).Error("error while getting omdb movie")
					continue
				}
			}

			// add a file to preview (downloaded in another cronjob)
			file := model.File{
				Name: item.Title,
				Path: MovieImageDirectory,
				URL:  null.StringFrom(item.Enclosure.Url),
			}
			if err := c.db.Create(&file).Error; err != nil {
				log.WithFields(logFields).WithError(err).Error("error while creating file")
				continue
			}

			// save the result of the previous steps (🎉)
			movie := model.Kino{
				Date:        date,
				Title:       item.Title,
				Year:        omdbMovie.ReleaseYear,
				Runtime:     omdbMovie.Runtime,
				Genre:       omdbMovie.Genre,
				Director:    omdbMovie.Director,
				Actors:      omdbMovie.Actors,
				ImdbRating:  omdbMovie.ImdbRating,
				Description: omdbMovie.Plot, // we get this from imdb as tu-fim does truncate their plot
				FileID:      file.File,
				File:        file,
				Link:        item.Link,
			}
			if err := c.db.Create(&movie).Error; err != nil {
				log.WithFields(logFields).WithError(err).Error("error while creating movie")
				continue
			} else {
				log.WithFields(logFields).Debug("created movie")
			}
		}
	}
	return nil
}
