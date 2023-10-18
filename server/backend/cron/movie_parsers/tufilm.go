package movie_parsers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type TuFilmWebsiteInformation struct {
	ImdbID               string
	TrailerUrl           string
	ImageUrl             string
	ShortenedDescription string
	Location             string
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

	bm := bluemonday.StrictPolicy()
	return &TuFilmWebsiteInformation{
		ImdbID:               bm.Sanitize(parseImdbID(doc)),
		TrailerUrl:           bm.Sanitize(parseTrailerUrl(doc)),
		ImageUrl:             bm.Sanitize(parseImageUrl(doc)),
		ShortenedDescription: parseShortenedDescription(doc),
		Location:             bm.Sanitize(parseLocation(doc)),
	}, nil
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

	return result
}

func parseImageUrl(doc *goquery.Document) string {
	href, exists := doc.Find("img.poster").First().Attr("src")
	if !exists {
		return ""
	}
	return "https://www.tu-film.de" + href
}

func parseTrailerUrl(doc *goquery.Document) string {
	trailerLinks := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Text() == "Zum Trailer"
	})
	if trailerLinks.Length() == 0 {
		return ""
	}
	if trailerLinks.Length() > 1 {
		log.Warn("more than one trailer link found. using first one")
	}
	// extract the imdb id from the link
	href, exists := trailerLinks.First().Attr("href")
	if !exists {
		log.Error("'Zum Trailer' does not have a link")
		return ""
	}
	href = strings.Replace(href, "http://", "https://", 1)
	return href
}

func parseLocation(doc *goquery.Document) string {
	locationLinks := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		href, hrefExists := s.Attr("href")
		return hrefExists && strings.Contains(href, "https://goo.gl/maps/")
	})
	if locationLinks.Length() == 0 {
		return ""
	}
	if locationLinks.Length() > 1 {
		log.Warn("more than one location link found. using first one")
	}
	return locationLinks.First().Text()
}

func parseImdbID(doc *goquery.Document) string {
	imdbLinks := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		href, hrefExists := s.Attr("href")
		return hrefExists && strings.Contains(href, "imdb.com/title/")
	})
	if imdbLinks.Length() == 0 {
		return ""
	}
	if imdbLinks.Length() > 1 {
		log.Warn("more than one imdb link found. using first one")
	}
	// extract the imdb id from the link
	href, _ := imdbLinks.First().Attr("href")
	re := regexp.MustCompile(`https?://www.imdb.com/title/(?P<imdb_id>[^/]+)/?`)
	return re.FindStringSubmatch(href)[re.SubexpIndex("imdb_id")]
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

// GetUpcomingFeed downloads a file from a given url and returns the path to the file
func GetUpcomingFeed() ([]MovieChannel, error) {
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
