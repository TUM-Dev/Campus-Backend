package cron

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
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
	omdbKeyExists := strings.TrimSpace(os.Getenv("OMDB_API_KEY")) != ""
	log.Trace("parsing upcoming feed")
	var allMovieLinks []string
	if err := c.db.Model(&model.Movie{}).
		Distinct().
		Pluck("Link", &allMovieLinks).Error; err != nil {
		return err
	}

	var importErrors []error
	channels, err := movie_parsers.GetFeeds()
	if err != nil {
		importErrors = append(importErrors, err)
	}
	for _, channel := range channels {
		for _, item := range channel.Items {
			logFields := log.Fields{"link": item.Link, "title": item.Title, "date": item.PubDate, "location": item.Location, "url": item.Enclosure.Url}
			if slices.Contains(allMovieLinks, item.Link) {
				log.WithFields(logFields).Trace("Movie already exists")
				continue
			}

			// populate extra data from tu-film website
			movieInformation, err := movie_parsers.GetTuFilmWebsiteInformation(item.Link)
			if err != nil {
				log.WithFields(logFields).WithError(err).Error("error while reading movie page")
				importErrors = append(importErrors, fmt.Errorf("movie %s: %w", item.Link, err))
				continue
			}
			date := movieInformation.ScreeningDate
			if date.IsZero() {
				date, err = time.Parse(time.RFC1123Z, item.PubDate)
				if err != nil {
					importErrors = append(importErrors, fmt.Errorf("invalid screening date for %s: %w", item.Link, err))
					continue
				}
			}
			if movieInformation.Location != "" {
				item.Location = movieInformation.Location
			}
			if err := c.storeParsedMovie(item.Title, item.Link, item.Location, date, movieInformation, omdbKeyExists, logFields); err != nil {
				log.WithFields(logFields).WithError(err).Error("error while creating movie")
				importErrors = append(importErrors, err)
			} else {
				allMovieLinks = append(allMovieLinks, item.Link)
			}
		}
	}

	campuskinoScreenings, err := movie_parsers.GetCampuskinoHeilbronnScreenings()
	if err != nil {
		log.WithError(err).Error("error while fetching campuskino heilbronn screenings")
		importErrors = append(importErrors, err)
	}
	for _, screening := range campuskinoScreenings {
		logFields := log.Fields{"link": screening.Link, "title": screening.Title, "date": screening.Info.ScreeningDate, "location": screening.Info.Location}
		if slices.Contains(allMovieLinks, screening.Link) {
			log.WithFields(logFields).Trace("Movie already exists")
			continue
		}
		info := screening.Info
		if err := c.storeParsedMovie(screening.Title, screening.Link, info.Location, info.ScreeningDate, &info, omdbKeyExists, logFields); err != nil {
			log.WithFields(logFields).WithError(err).Error("error while creating campuskino heilbronn movie")
			importErrors = append(importErrors, err)
		} else {
			allMovieLinks = append(allMovieLinks, screening.Link)
		}
	}

	return errors.Join(importErrors...)
}

func (c *CronService) storeParsedMovie(title, link, location string, date time.Time, movieInformation *movie_parsers.TuFilmWebsiteInformation, omdbKeyExists bool, logFields log.Fields) error {
	movie := model.Movie{
		Date:        date,
		Title:       title,
		Location:    null.StringFrom(location),
		Year:        movieInformation.ReleaseYear,
		Runtime:     movieInformation.Runtime,
		Genre:       movieInformation.Genre,
		Director:    movieInformation.Director,
		Actors:      movieInformation.Actors,
		Description: movieInformation.ShortenedDescription,
		Trailer:     movieInformation.TrailerUrl,
		Link:        link,
	}
	if omdbKeyExists && movieInformation.ImdbID.ValueOrZero() != "" {
		omdbMovie, err := movie_parsers.GetOmdbMovie(movieInformation.ImdbID.ValueOrZero())
		if err != nil {
			log.WithFields(logFields).WithError(err).Warn("OMDb enrichment failed; using source movie data")
		} else {
			enrichMovie(&movie, omdbMovie)
		}
	}

	return c.db.Transaction(func(tx *gorm.DB) error {
		file, err := saveImage(tx, movieInformation.ImageUrl)
		if err != nil {
			return err
		}
		movie.File = *file
		movie.FileID = file.File
		return tx.Omit("File").Create(&movie).Error
	})
}

// saveImage saves an image to the database, so it can be downloaded by another cronjob and returns the file
func saveImage(tx *gorm.DB, url string) (*model.File, error) {
	targetFileName := movieImageName(url)
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

func movieImageName(address string) string {
	parsed, err := url.Parse(address)
	extension := ".jpg"
	if err == nil {
		if ext := path.Ext(parsed.Path); ext != "" {
			extension = ext
		}
	}
	return fmt.Sprintf("%x%s", md5.Sum([]byte(address)), extension)
}

// Optional enrichment must not replace useful source data with empty or N/A fields.
func enrichMovie(movie *model.Movie, data *movie_parsers.OmdbResults) {
	fields := []struct {
		target *null.String
		value  string
	}{
		{&movie.Year, data.ReleaseYear}, {&movie.Runtime, data.Runtime},
		{&movie.Genre, data.Genre}, {&movie.Director, data.Director},
		{&movie.Actors, data.Actors}, {&movie.ImdbRating, data.ImdbRating},
	}
	for _, field := range fields {
		if value := strings.TrimSpace(field.value); value != "" && value != "N/A" {
			*field.target = null.StringFrom(value)
		}
	}
	if plot := strings.TrimSpace(data.Plot); plot != "" && plot != "N/A" {
		movie.Description = plot
	}
}
