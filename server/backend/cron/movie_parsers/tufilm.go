package movie_parsers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/guregu/null"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
)

type TuFilmWebsiteInformation struct {
	ImdbID               null.String
	TrailerUrl           null.String
	ImageUrl             string
	ShortenedDescription string
	ReleaseYear          null.String
	Director             null.String
	Actors               null.String
	Runtime              null.String
}

// GetTuFilmWebsiteInformation scrapes the tu-film website for all usefully information
// url: url of the tu-film website, e.g. https://www.tu-film.de/programm/view/1204
func GetTuFilmWebsiteInformation(url string) (*TuFilmWebsiteInformation, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("error while getting response for request")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing body")
		}
	}(resp.Body)
	// parse the response body
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.WithError(err).Error("Error while parsing document")
		return nil, err
	}

	return parseWebsiteInformation(doc)
}

func parseWebsiteInformation(doc *goquery.Document) (*TuFilmWebsiteInformation, error) {
	Director, Actors, Runtime := parseDirectorActorsRuntime(doc)
	return &TuFilmWebsiteInformation{
		ImdbID:               parseImdbID(doc),
		TrailerUrl:           parseTrailerUrl(doc),
		ImageUrl:             parseImageUrl(doc),
		ShortenedDescription: parseShortenedDescription(doc),
		ReleaseYear:          parseReleaseYear(doc),
		Director:             Director,
		Actors:               Actors,
		Runtime:              Runtime,
	}, nil
}

func parseDirectorActorsRuntime(doc *goquery.Document) (null.String, null.String, null.String) {
	bm := bluemonday.StrictPolicy()
	rawTable := doc.Find("td.film").Text()
	rawTable = bm.Sanitize(rawTable)
	rawTable = strings.TrimSpace(rawTable)
	for strings.Contains(rawTable, "  ") {
		rawTable = strings.ReplaceAll(rawTable, "  ", " ")
		rawTable = strings.ReplaceAll(rawTable, "\n ", "\n")
	}
	re := regexp.MustCompile(`Regie: (?P<director>.+)\nSchauspieler: (?P<actors>.+)\n(?P<runtime>\d+) Minuten`)
	matches := re.FindStringSubmatch(rawTable)
	if len(matches) < re.NumSubexp() {
		return null.String{}, null.String{}, null.String{}
	}
	director, actors, runtime := matches[re.SubexpIndex("director")], matches[re.SubexpIndex("actors")], matches[re.SubexpIndex("runtime")]
	return null.StringFrom(director), null.StringFrom(actors), null.StringFrom(runtime + " min")
}

func parseReleaseYear(doc *goquery.Document) null.String {
	releasePlaceYear := doc.Find(".title h4").Text()
	re := regexp.MustCompile(`(?P<release_place>.+) \((?P<release_year>\d{4})\)`)
	match := re.FindStringSubmatch(releasePlaceYear)
	index := re.SubexpIndex("release_year")
	if len(match) < index+1 {
		return null.String{}
	}
	return null.StringFrom(match[index])
}

func parseShortenedDescription(doc *goquery.Document) string {
	bm := bluemonday.StrictPolicy()
	result := ""
	teaser := strings.TrimSpace(doc.Find("div.teaser").Text())
	if teaser != "" {
		result += fmt.Sprintf("<b>%s<b>\n", bm.Sanitize(teaser))
	}

	doc.Find("div.description").Each(func(i int, s *goquery.Selection) {
		// images which are randomly inserted in the text is great for the website, but not for us
		s.Find("img").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
		// paragraphs are used to separate the text
		s.Find("p").Each(func(i int, s *goquery.Selection) {
			content := strings.TrimSpace(s.Text())
			content = strings.ReplaceAll(content, "\n", "")
			if content != "" {
				result += "\n" + bm.Sanitize(content)
			}
		})
	})
	// cover the case where we have no teaser
	result = strings.TrimSpace(result)

	comment := doc.Find("div.comment").Text()
	comment = strings.TrimSpace(comment)
	comment = strings.ReplaceAll(comment, "\n", "")
	if comment != "" {
		result += fmt.Sprintf("\n\n<i>%s<i>", bm.Sanitize(comment))
	}
	// clean up the result
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}
	if result == "" {
		return "Surprise yourself"
	}
	return result
}

func parseImageUrl(doc *goquery.Document) string {
	href, exists := doc.Find("img.poster").First().Attr("src")
	if !exists {
		return "https://www.tu-film.de/img/film/poster/.sized.berraschungsfilm.jpg"
	}
	sanitisedHref := bluemonday.StrictPolicy().Sanitize(href)
	return "https://www.tu-film.de" + sanitisedHref
}

func parseTrailerUrl(doc *goquery.Document) null.String {
	trailerLinks := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Text() == "Zum Trailer"
	})
	if trailerLinks.Length() == 0 {
		return null.String{}
	}
	if trailerLinks.Length() > 1 {
		log.Warn("more than one trailer link found. using first one")
	}
	// extract the imdb id from the link
	href, exists := trailerLinks.First().Attr("href")
	if !exists {
		log.Error("'Zum Trailer' does not have a link")
		return null.String{}
	}
	href = strings.Replace(href, "http://", "https://", 1)
	href = bluemonday.StrictPolicy().Sanitize(href)
	return null.StringFrom(href)
}

func parseImdbID(doc *goquery.Document) null.String {
	imdbLinks := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		href, hrefExists := s.Attr("href")
		return hrefExists && strings.Contains(href, "imdb.com/title/")
	})
	if imdbLinks.Length() == 0 {
		return null.String{}
	}
	if imdbLinks.Length() > 1 {
		log.Warn("more than one imdb link found. using first one")
	}
	// extract the imdb id from the link
	href, _ := imdbLinks.First().Attr("href")
	re := regexp.MustCompile(`.*imdb.com/title/(?P<imdb_id>[a-zA-Z0-9]+)/?`)
	return null.StringFrom(re.FindStringSubmatch(href)[re.SubexpIndex("imdb_id")])
}

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

// GetFeeds gets all feeds from the tu-film website
func GetFeeds() ([]MovieChannel, error) {
	var channels []MovieChannel
	if newChannels, err := GetFeed("https://www.tu-film.de/programm/index/upcoming.rss"); err != nil {
		return nil, err
	} else {
		channels = append(channels, newChannels...)
	}
	for i := 2001; i <= 2023; i++ {
		for _, semester := range []string{"ws", "ss"} {
			if newChannels, err := GetFeed(fmt.Sprintf("https://www.tu-film.de/programm/index/%s%d.rss", semester, i)); err != nil {
				log.WithError(err).Warn("Error while getting old movie feed")
			} else {
				channels = append(channels, newChannels...)
			}
		}
	}
	return channels, nil
}

func GetFeed(url string) ([]MovieChannel, error) {
	resp, err := http.Get(url)
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
	var newMovies struct {
		Channels []MovieChannel `xml:"channel"`
	}
	err = xml.NewDecoder(resp.Body).Decode(&newMovies)
	if err != nil {
		log.WithError(err).Error("Error while unmarshalling UpcomingFeed")
		return nil, err
	}
	return newMovies.Channels, nil
}
