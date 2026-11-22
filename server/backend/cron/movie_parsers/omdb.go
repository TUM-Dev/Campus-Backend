package movie_parsers

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"regexp"
	"strings"

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

var (
	omdbBracketTagRegex  = regexp.MustCompile(`\[[^\]]*\]`)
	omdbYearParensRegex  = regexp.MustCompile(`\(\d{4}\)`)
	omdbNoticePrefixExpr = regexp.MustCompile(`(?i)^(terminverschiebung|vorpremiere|preview|sondervorstellung)\s*:\s*`)
)

func GetOmdbMovie(id string) (*OmdbResults, error) {
	return queryOmdb(url.Values{"r": {"json"}, "v": {"1"}, "i": {id}, "apikey": {os.Getenv("OMDB_API_KEY")}})
}

// GetOmdbMovieByTitle queries OMDb by cleaned movie title and optional release year.
func GetOmdbMovieByTitle(rawTitle, year string) (*OmdbResults, error) {
	title := CleanMovieTitleForSearch(rawTitle)
	if title == "" {
		return nil, errors.New("empty movie title for OMDb lookup")
	}
	params := url.Values{"r": {"json"}, "v": {"1"}, "t": {title}, "apikey": {os.Getenv("OMDB_API_KEY")}}
	if y := strings.TrimSpace(year); y != "" {
		params.Set("y", y)
	}
	return queryOmdb(params)
}

// CleanMovieTitleForSearch strips language tags (e.g. "[OmU]"), notice prefixes (e.g. "Terminverschiebung:"),
// release years in parentheses, and appended subtitles after " – " so OMDb title lookup can match.
func CleanMovieTitleForSearch(title string) string {
	cleaned := omdbBracketTagRegex.ReplaceAllString(title, " ")
	cleaned = omdbYearParensRegex.ReplaceAllString(cleaned, " ")
	if idx := strings.Index(cleaned, " – "); idx > 0 {
		cleaned = cleaned[:idx]
	}
	cleaned = strings.TrimSpace(cleaned)
	cleaned = omdbNoticePrefixExpr.ReplaceAllString(cleaned, "")
	return strings.Join(strings.Fields(cleaned), " ")
}

func queryOmdb(params url.Values) (*OmdbResults, error) {
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
