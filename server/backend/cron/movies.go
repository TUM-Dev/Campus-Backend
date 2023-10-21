package cron

import (
	"regexp"
	"slices"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/cron/movie_parsers"

	"github.com/guregu/null"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

const (
	MovieImageDirectory = "movie/"
)

func (c *CronService) movieCron() error {
	log.Trace("parsing upcoming feed")
	var allMovieLinks []string
	if err := c.db.Model(&model.Kino{}).Distinct().Pluck("Link", &allMovieLinks).Error; err != nil {
		return err
	}

	channels, err := movie_parsers.GetUpcomingFeed()
	if err != nil {
		return err
	}
	for _, channel := range channels {
		for _, item := range channel.Items {
			logFields := log.Fields{"link": item.Link, "title": item.Title, "date": item.PubDate, "location": item.Location, "url": item.Enclosure.Url}
			if slices.Contains(allMovieLinks, item.Link) {
				log.WithFields(logFields).Trace("Movie already exists")
				continue
			}

			// data cleanup
			date, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				log.WithError(err).WithFields(logFields).Error("Couldn't check if movie already exists")
				continue
			}
			re := regexp.MustCompile(`(?P<date>[\d. ]+): (?P<title>.+)$`)
			matches := re.FindStringSubmatch(item.Title)
			if len(matches) < re.NumSubexp() {
				log.WithFields(logFields).Error("Couldn't parse movie title")
				continue
			}
			item.Title = matches[re.SubexpIndex("title")]

			// populate extra data from omdb
			movieInformation, err := movie_parsers.GetTuFilmWebsiteInformation(item.Link)
			if err != nil {
				log.WithFields(logFields).WithError(err).Error("error while finding imdb id")
				continue
			}
			movie := model.Kino{
				Date:        date,
				Title:       item.Title,
				Location:    null.StringFrom(item.Location),
				Year:        movieInformation.ReleaseYear,
				Runtime:     movieInformation.Runtime,
				Director:    movieInformation.Director,
				Actors:      movieInformation.Actors,
				Description: movieInformation.ShortenedDescription,
				Trailer:     movieInformation.TrailerUrl,
				Link:        item.Link,
			}
			previewFile := model.File{
				Name: item.Title,
				Path: MovieImageDirectory,
				URL:  null.StringFrom(item.Enclosure.Url),
			}
			if movieInformation.ImdbID.ValueOrZero() != "" {
				omdbMovie, err := movie_parsers.GetOmdbMovie(movieInformation.ImdbID.ValueOrZero())
				if err != nil {
					log.WithFields(logFields).WithError(err).Error("error while getting omdb movie")
					continue
				}
				// enrich the movie with data from omdb if present
				movie.Year = null.StringFrom(omdbMovie.ReleaseYear)
				movie.Runtime = null.StringFrom(omdbMovie.Runtime)
				movie.Genre = null.StringFrom(omdbMovie.Genre)
				movie.Director = null.StringFrom(omdbMovie.Director)
				movie.Actors = null.StringFrom(omdbMovie.Actors)
				movie.ImdbRating = null.StringFrom(omdbMovie.ImdbRating)
				movie.Description = omdbMovie.Plot // tu-fim does truncate their plot
			}

			// save the result of the previous steps (🎉)
			if err := c.db.Create(&previewFile).Error; err != nil {
				log.WithFields(logFields).WithError(err).Error("error while creating file")
				continue
			}
			// assign the file_id to make sure the id is assigned
			movie.File = previewFile
			movie.FileID = previewFile.File
			if err := c.db.Create(&movie).Error; err != nil {
				log.WithFields(logFields).WithError(err).Error("error while creating movie")
			}
		}
	}
	return nil
}
