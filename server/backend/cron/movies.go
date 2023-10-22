package cron

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/guregu/null"

	"github.com/PuerkitoBio/goquery"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

type MovieItems struct {
	Title     string `xml:"title"`
	Link      string `xml:"link"`
	PubDate   string `xml:"pubDate"`
	Location  string `xml:"location"`
	Enclosure struct {
		Url    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
}

type MovieChannel struct {
	Items []MovieItems `xml:"item"`
}

const (
	MovieImageDirectory = "movie/"
)

func (c *CronService) movieCron() error {
	log.Trace("parsing upcoming feed")
	channels, err := parseUpcomingFeed()
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
			imdbID, err := extractImdbIDFromTUFilmWebsite(item.Link)
			if err != nil {
				log.WithFields(logFields).WithError(err).Error("error while finding imdb id")
				continue
			}
			omdbMovie, err := getOmdbMovie(imdbID)
			if err != nil {
				log.WithFields(logFields).WithError(err).Error("error while getting omdb movie")
				continue
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

type omdbResults struct {
	ReleaseYear string `json:"Year"`
	Runtime     string
	Genre       string
	Director    string
	Actors      string
	Plot        string
	ImdbRating  string `json:"imdbRating"`
}

func getOmdbMovie(id string) (*omdbResults, error) {
	url := fmt.Sprintf("https://www.omdbapi.com/?r=json&v=1&i=%s&apikey=%s", id, os.Getenv("OMDB_API_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		log.WithField("url", url).WithError(err).Error("Error while getting response for request")
		return nil, err
	}
	// check if the api key is valid
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("missing or invalid api key for omdb (environment variable OMDB_API_KEY)")
	}
	// other errors
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Warn("Unable to read http body")
			return nil, err
		} else {
			log.WithField("status", resp.StatusCode).WithField("status", resp.Status).WithField("body", string(body)).Error("error while getting omdb movie")
			return nil, errors.New("error while getting omdb movie")
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.WithField("url", url).WithError(err).Error("Error while closing body")
		}
	}(resp.Body)
	// parse the response body
	var res omdbResults
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.WithField("url", url).WithError(err).Error("Error while unmarshalling omdbResults")
		return nil, err
	}
	return &res, nil
}

// extractImdbIDFromTUFilmWebsite scrapes the imdb id and fullDescription from the tu-film website
// url: url of the tu-film website, e.g. https://www.tu-film.de/programm/view/1204
func extractImdbIDFromTUFilmWebsite(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("error while getting response for request")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing body")
		}
	}(resp.Body)
	// parse the response body
	return parseImdbIDFromReader(resp.Body)
}

func parseImdbIDFromReader(body io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.WithError(err).Error("Error while parsing document")
		return "", err
	}

	// extract the imdb link
	imdbLinks := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		href, hrefExists := s.Attr("href")
		return hrefExists && strings.Contains(href, "imdb.com/title/")
	})
	if imdbLinks.Length() == 0 {
		return "", errors.New("no imdb link found")
	}
	if imdbLinks.Length() > 1 {
		log.Warn("more than one imdb link found. using first one")
	}
	// extract the imdb id from the link
	href, _ := imdbLinks.First().Attr("href")
	re := regexp.MustCompile(`https?://www.imdb.com/title/(?P<imdb_id>[^/]+)/?`)
	return re.FindStringSubmatch(href)[re.SubexpIndex("imdb_id")], nil
}

// parseUpcomingFeed downloads a file from a given url and returns the path to the file
func parseUpcomingFeed() ([]MovieChannel, error) {
	resp, err := http.Get("https://www.tu-film.de/programm/index/upcoming.rss")
	if err != nil {
		log.WithError(err).Error("Error while getting response for request")
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing body")
		}
	}(resp.Body)
	//Parse the data into a struct
	var upcomingMovies struct {
		Channels []MovieChannel `xml:"channel"`
	}
	err = xml.NewDecoder(resp.Body).Decode(&upcomingMovies)
	if err != nil {
		log.WithError(err).Error("Error while unmarshalling UpcomingFeed")
		return nil, err
	}
	return upcomingMovies.Channels, nil
}
