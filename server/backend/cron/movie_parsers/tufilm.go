package movie_parsers

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/PuerkitoBio/goquery"
	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/guregu/null"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
)

type TuFilmWebsiteInformation struct {
	ScreeningDate        time.Time
	Location             string
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
	resp, err := webfetch.Get(url)
	if err != nil {
		return nil, fmt.Errorf("getting movie page: %w", err)
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
	date, location, err := parseScreening(doc)
	if err != nil {
		return nil, err
	}
	return &TuFilmWebsiteInformation{
		ScreeningDate:        date,
		Location:             location,
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
	if doc.Find(".programentry td.poster").Length() > 0 {
		var director, actors, runtime null.String
		doc.Find(".programentry td.poster p").Each(func(_ int, p *goquery.Selection) {
			text := strings.TrimSpace(p.Text())
			switch {
			case strings.HasPrefix(text, "Regie:"):
				director = null.StringFrom(strings.TrimSpace(strings.TrimPrefix(text, "Regie:")))
			case strings.HasPrefix(text, "Schauspieler:"):
				actors = null.StringFrom(strings.TrimSpace(strings.TrimPrefix(text, "Schauspieler:")))
			case regexp.MustCompile(`^\d+ Minuten$`).MatchString(text):
				runtime = null.StringFrom(strings.TrimSuffix(text, " Minuten") + " min")
			}
		})
		return director, actors, runtime
	}
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
	releasePlaceYear := doc.Find(".title h4, .programentry h3").Text()
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
	teaser := strings.TrimSpace(doc.Find("div.teaser, .programentry .teaser").Text())
	if teaser != "" {
		result += fmt.Sprintf("<b>%s<b>\n", bm.Sanitize(teaser))
	}

	doc.Find("div.description, .programentry td.text").Each(func(i int, s *goquery.Selection) {
		// images which are randomly inserted in the text is great for the website, but not for us
		s.Find("img").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
		// paragraphs are used to separate the text
		s.Find("p:not(.teaser):not(.comment)").Each(func(i int, s *goquery.Selection) {
			content := strings.TrimSpace(s.Text())
			content = strings.ReplaceAll(content, "\n", "")
			if content != "" {
				result += "\n" + bm.Sanitize(content)
			}
		})
	})
	// cover the case where we have no teaser
	result = strings.TrimSpace(result)

	comment := doc.Find("div.comment, .programentry .comment").Text()
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
	href, exists := doc.Find("img.poster, .programentry td.poster img").First().Attr("src")
	if !exists {
		return "https://www.tu-film.de/img/film/poster/berraschungsfilm.jpg"
	}
	sanitisedHref := bluemonday.StrictPolicy().Sanitize(href)
	sanitisedHref = strings.ReplaceAll(sanitisedHref, "poster/.sized.", "poster/")
	sanitisedHref = strings.ReplaceAll(sanitisedHref, "poster/.thumb.", "poster/")
	return resolveURL(sanitisedHref)
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
	matches := re.FindStringSubmatch(href)
	if len(matches) < 2 {
		return null.String{}
	}
	return null.StringFrom(matches[re.SubexpIndex("imdb_id")])
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

// GetFeeds discovers the current semester from the site navigation. The old
// /programm/index/*.rss feeds were removed when tu-film moved to a static site.
func GetFeeds() ([]MovieChannel, error) {
	home, err := getDocument("https://www.tu-film.de/")
	if err != nil {
		return nil, err
	}
	href, ok := home.Find("nav a[href^='/program/']").First().Attr("href")
	if !ok {
		return nil, errors.New("tu-film program link missing")
	}
	program, err := getDocument(resolveURL(href))
	if err != nil {
		return nil, err
	}
	return parseProgram(program)
}

func getDocument(address string) (*goquery.Document, error) {
	resp, err := webfetch.Get(address)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return goquery.NewDocumentFromReader(resp.Body)
}

func resolveURL(reference string) string {
	base, _ := url.Parse("https://www.tu-film.de/")
	ref, err := url.Parse(reference)
	if err != nil {
		return ""
	}
	return base.ResolveReference(ref).String()
}

func parseProgram(doc *goquery.Document) ([]MovieChannel, error) {
	channel := MovieChannel{}
	seen := map[string]bool{}
	doc.Find(".plan td.title a[href]").Each(func(_ int, a *goquery.Selection) {
		href, _ := a.Attr("href")
		link := resolveURL(href)
		if seen[link] {
			return
		}
		seen[link] = true
		channel.Items = append(channel.Items, MovieItems{Title: strings.Join(strings.Fields(a.Parent().Text()), " "), Link: link})
	})
	if len(channel.Items) == 0 {
		return nil, errors.New("tu-film program contains no screenings")
	}
	return []MovieChannel{channel}, nil
}

// parseScreening uses the actual local screening time, rather than RSS pubDate
// (the replacement site's RSS entries are dated January 1970).
func parseScreening(doc *goquery.Document) (time.Time, string, error) {
	event := doc.Find("p.event")
	if event.Length() == 0 {
		if doc.Find(".programentry").Length() > 0 {
			return time.Time{}, "", errors.New("tu-film screening information missing")
		}
		return time.Time{}, "", nil
	} // legacy page
	event = event.Clone()
	event.Find("br").ReplaceWithHtml(" ")
	text := strings.Join(strings.Fields(event.Text()), " ")
	re := regexp.MustCompile(`(\d{1,2})\. (\p{L}+) (\d{4}) um (\d{1,2}:\d{2}) Uhr`)
	parts := re.FindStringSubmatch(text)
	if len(parts) != 5 {
		return time.Time{}, "", errors.New("tu-film screening date missing or invalid")
	}
	months := map[string]string{"Januar": "01", "Februar": "02", "März": "03", "April": "04", "Mai": "05", "Juni": "06", "Juli": "07", "August": "08", "September": "09", "Oktober": "10", "November": "11", "Dezember": "12"}
	month, ok := months[parts[2]]
	if !ok {
		return time.Time{}, "", errors.New("unknown screening month")
	}
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return time.Time{}, "", err
	}
	date, err := time.ParseInLocation("2 01 2006 15:04", parts[1]+" "+month+" "+parts[3]+" "+parts[4], berlin)
	return date, strings.TrimSpace(event.Find("a").Text()), err
}
