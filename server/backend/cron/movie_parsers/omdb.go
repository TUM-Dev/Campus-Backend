package movie_parsers

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

type OmdbResults struct {
	ReleaseYear string `json:"Year"`
	Runtime     string
	Genre       string
	Director    string
	Actors      string
	Plot        string
	ImdbRating  string `json:"imdbRating"`
}

func GetOmdbMovie(id string) (*OmdbResults, error) {
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
	var res OmdbResults
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.WithField("url", url).WithError(err).Error("Error while unmarshalling omdbResults")
		return nil, err
	}
	return &res, nil
}
