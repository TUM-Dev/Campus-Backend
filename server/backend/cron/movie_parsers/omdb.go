package movie_parsers

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"

	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
)

type OmdbResults struct {
	ReleaseYear string `json:"Year"`
	Runtime     string
	Genre       string
	Director    string
	Actors      string
	Plot        string
	ImdbRating  string `json:"imdbRating"`
	Response    string
	Error       string
}

func GetOmdbMovie(id string) (*OmdbResults, error) {
	params := url.Values{"r": {"json"}, "v": {"1"}, "i": {id}, "apikey": {os.Getenv("OMDB_API_KEY")}}
	resp, err := webfetch.Get("https://www.omdbapi.com/?" + params.Encode())
	if err != nil {
		// Transport errors can contain the URL, including the API key.
		return nil, errors.New("OMDb request failed")
	}
	defer func() { _ = resp.Body.Close() }()
	var result OmdbResults
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Response == "False" {
		return nil, errors.New("OMDb could not return movie information")
	}
	return &result, nil
}
