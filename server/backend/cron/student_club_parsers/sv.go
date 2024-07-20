package student_club_parsers

import (
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/guregu/null"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
)

type SVStudentClub struct {
	Name         string
	Collection   string
	Description  null.String
	ImageUrl     null.String
	ImageCaption null.String
	LinkUrl      null.String
}

type SVStudentClubCollection struct {
	Name        string
	Description string
}

func DownloadHtml(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("error while getting response for request")
	}
	return resp.Body, nil
}

func ParseStudentClubs(reader io.Reader) ([]SVStudentClub, []SVStudentClubCollection, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.WithError(err).Error("Error while parsing student clubs into goquery document")
		return nil, nil, err
	}
	scrapedClubs := make([]SVStudentClub, 0)
	scrapedClubCollections := make([]SVStudentClubCollection, 0)
	currentCollectionHeader := ""
	filterForRelevantFrames(doc).Each(func(i int, s *goquery.Selection) {
		// section headers and regular headers best differentiated by either having a h2 or h6 tag.
		// h2 => this is a studentClub Collection
		// h6 => this is a studentClub
		if s.Has("h2").Length() >= 1 {
			collection := SVStudentClubCollection{
				Name:        parseTitle(s, "h2"),
				Description: parseDescription(s, "p").ValueOrZero(),
			}
			currentCollectionHeader = collection.Name
			scrapedClubCollections = append(scrapedClubCollections, collection)
		} else if s.Has("h6").Length() >= 1 {
			image, imgCaption := parseImages(s)
			club := SVStudentClub{
				Name:         parseTitle(s, "h6"),
				Collection:   currentCollectionHeader,
				Description:  parseDescription(s, "p"),
				ImageUrl:     image,
				ImageCaption: imgCaption,
				LinkUrl:      parseLink(s),
			}
			scrapedClubs = append(scrapedClubs, club)
		} else {
			logForSelection(log.ErrorLevel, "Cannot parse html into student clubs", s)
		}
	})
	return scrapedClubs, scrapedClubCollections, nil
}

func parseImages(s *goquery.Selection) (null.String, null.String) {
	href, exists := s.Find("img").First().Attr("src")
	if !exists {
		return null.String{}, null.String{}
	}
	sanitisedHref := bluemonday.StrictPolicy().Sanitize(href)
	figcaption := null.String{}
	if s.Has("figcaption").Length() >= 1 {
		caption := s.Find("figcaption").Text()
		figcaption = null.StringFrom(trimSpaces(caption))
	}
	return null.StringFrom("https://sv.tum.de" + sanitisedHref), figcaption
}

func parseLink(s *goquery.Selection) null.String {
	href, exists := s.Find("a").First().Attr("href")
	if !exists {
		return null.String{}
	}
	sanitisedHref := bluemonday.StrictPolicy().Sanitize(href)
	return null.StringFrom(sanitisedHref)
}

func filterForRelevantFrames(doc *goquery.Document) *goquery.Selection {
	passesFilter := false
	return doc.
		Find(".frame").
		FilterFunction(func(_ int, sel *goquery.Selection) bool {
			if sel.HasClass("frame-type-menu_categorized_content") {
				passesFilter = true
				return false
			}
			if sel.HasClass("frame-space-before-extra-large") {
				passesFilter = false
				return false
			}
			// the selections we want are all <div> but there are some other frames like a <nav> element
			if !sel.Is("div") {
				return false
			}
			numberOfRelevantTitles := sel.Has("h2").Length() + sel.Has("h6").Length()
			if numberOfRelevantTitles == 0 {
				return false
			}
			// seems to be an error, idk why this is the case
			if sel.Length() == 0 {
				return false
			}
			return passesFilter
		})
}

func parseTitle(s *goquery.Selection, selector string) string {
	titles := s.Find(selector)
	if titles.Length() > 1 {
		logForSelection(log.WarnLevel, fmt.Sprintf("more than student club's %s found. using first one", selector), s)
	}
	title := titles.First().Text()
	return trimSpaces(sanitise(title))
}

func parseDescription(s *goquery.Selection, selector string) null.String {
	descriptions := s.Find(selector)
	if descriptions.Length() == 0 {
		return null.String{}
	}
	description := descriptions.First().Text()
	description = trimSpaces(sanitise(description))
	return null.StringFrom(description)
}

func sanitise(s string) string {
	s = html.UnescapeString(s)
	return bluemonday.StrictPolicy().Sanitize(s)

}

func trimSpaces(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ") // &nbsp; does not make huge sense given that we layout differently
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.TrimSpace(s)
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}

func logForSelection(level log.Level, message string, selection *goquery.Selection) {
	dumpedHtml, err := selection.Html()
	if err != nil {
		log.WithError(err).Log(level, message)
	} else {
		readableHtml := trimSpaces(dumpedHtml)
		readableHtml = strings.ReplaceAll(readableHtml, "\"", "'")
		log.WithField("html", readableHtml).Error(message)
	}
}
