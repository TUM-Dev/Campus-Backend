package cron

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/TUM-Dev/Campus-Backend/server/backend/cron/movie_parsers"

	"github.com/guregu/null"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

const (
	MovieImageDirectory = "movie/"
)

func (c *CronService) movieCron() error {
	_, omdbKeyExists := os.LookupEnv("OMDB_API_KEY")
	if !omdbKeyExists {
		log.Info("Skipping movieCron as no OMDB_API_KEY was set")
		return nil
	}
	log.Trace("parsing upcoming feed")
	var allMovieLinks []string
	if err := c.db.Model(&model.Movie{}).
		Distinct().
		Pluck("Link", &allMovieLinks).Error; err != nil {
		return err
	}

	channels, err := movie_parsers.GetFeeds()
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

			// populate extra data from tu-film website
			movieInformation, err := movie_parsers.GetTuFilmWebsiteInformation(item.Link)
			if err != nil {
				log.WithFields(logFields).WithError(err).Error("error while finding imdb id")
				continue
			}
			movie := model.Movie{
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
			if err := c.db.Transaction(func(tx *gorm.DB) error {
				file, err := saveImage(tx, movieInformation.ImageUrl)
				if err != nil {
					return err
				}
				// assign the file_id to make sure the id is assigned
				movie.File = *file
				movie.FileID = file.File
				return tx.Create(&movie).Error
			}); err != nil {
				log.WithFields(logFields).WithError(err).Error("error while creating movie")
			}
		}
	}
	return nil
}

// saveImage saves an image to the database, so it can be downloaded by another cronjob and returns the file
func saveImage(tx *gorm.DB, url string) (*model.File, error) {
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
		Path:       MovieImageDirectory,
		URL:        null.StringFrom(url),
		Downloaded: null.BoolFrom(false),
	}
	if err := tx.Create(&file).Error; err != nil {
		log.WithError(err).Error("Could not store new file to database")
		return nil, err
	}
	return &file, nil
}
