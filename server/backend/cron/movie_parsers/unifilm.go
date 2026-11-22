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

const (
	UnifilmHeilbronnURL         = "https://www.unifilm.de/studentenkinos/Heilbronn"
	unifilmBaseURL              = "https://www.unifilm.de/"
	unifilmDefaultImageFallback = "https://www.tu-film.de/img/film/poster/berraschungsfilm.jpg"
)

var (
	unifilmDateRegex       = regexp.MustCompile(`(\d{1,2})\.(\d{1,2})\.(\d{4})`)
	unifilmTimeRegex       = regexp.MustCompile(`(\d{1,2}:\d{2})`)
	unifilmRuntimeRegex    = regexp.MustCompile(`^(\d+)\s*Min\.?$`)
	unifilmTitleYearRegex  = regexp.MustCompile(`\((\d{4})\)`)
	unifilmFSKRegex        = regexp.MustCompile(`(?i)^FSK\s*\d+$`)
	unifilmCountryCodeList = regexp.MustCompile(`^[A-Z]{2,3}(\s*,\s*[A-Z]{2,3})*$`)
)

// CampuskinoScreening holds a parsed Campuskino Heilbronn screening from unifilm.de.
type CampuskinoScreening struct {
	Title string
	Link  string
	Info  TuFilmWebsiteInformation
}

// GetCampuskinoHeilbronnScreenings fetches and parses the Campuskino Heilbronn schedule from unifilm.de.
func GetCampuskinoHeilbronnScreenings() ([]CampuskinoScreening, error) {
	resp, err := webfetch.Get(UnifilmHeilbronnURL)
	if err != nil {
		return nil, fmt.Errorf("getting campuskino heilbronn page: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.WithError(err).Error("Error while closing body")
		}
	}(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return ParseCampuskinoHeilbronnProgram(doc)
}

// ParseCampuskinoHeilbronnProgram parses all `.film-showcase` entries on the unifilm.de Campuskino Heilbronn page.
func ParseCampuskinoHeilbronnProgram(doc *goquery.Document) ([]CampuskinoScreening, error) {
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return nil, err
	}

	bm := bluemonday.StrictPolicy()
	var screenings []CampuskinoScreening
	var parseErrs []error
	seen := make(map[string]bool)

	doc.Find(".film-showcase").Each(func(_ int, s *goquery.Selection) {
		sid, ok := s.Attr("data-sid")
		if !ok || strings.TrimSpace(sid) == "" {
			sid, _ = s.Attr("data-id")
		}
		sid = strings.TrimSpace(sid)
		if sid == "" {
			return
		}
		link := UnifilmHeilbronnURL + "#sid=" + sid
		if seen[link] {
			return
		}
		seen[link] = true

		title := parseUnifilmTitle(s)
		if title == "" {
			parseErrs = append(parseErrs, fmt.Errorf("campuskino heilbronn screening %s: missing title", sid))
			return
		}

		date, location, err := parseUnifilmScreeningTimeAndLocation(s, berlin)
		if err != nil {
			parseErrs = append(parseErrs, fmt.Errorf("campuskino heilbronn screening %s (%s): %w", sid, title, err))
			return
		}

		director, genre, runtime := parseUnifilmMetadata(s)
		releaseYear := parseUnifilmReleaseYear(title)
		description := parseUnifilmDescription(s, bm)
		imageURL := parseUnifilmImageURL(s, bm)

		screenings = append(screenings, CampuskinoScreening{
			Title: title,
			Link:  link,
			Info: TuFilmWebsiteInformation{
				ScreeningDate:        date,
				Location:             location,
				ImageUrl:             imageURL,
				ShortenedDescription: description,
				ReleaseYear:          releaseYear,
				Director:             director,
				Genre:                genre,
				Runtime:              runtime,
			},
		})
	})

	if len(screenings) == 0 && len(parseErrs) > 0 {
		return nil, errors.Join(parseErrs...)
	}
	return screenings, errors.Join(parseErrs...)
}

func parseUnifilmTitle(s *goquery.Selection) string {
	var parts []string
	s.Find(".text_container h1 span").Each(func(_ int, span *goquery.Selection) {
		t := strings.Join(strings.Fields(span.Text()), " ")
		if t != "" {
			parts = append(parts, t)
		}
	})
	title := strings.Join(parts, " ")
	if title == "" {
		title = strings.Join(strings.Fields(s.Find(".text_container h1").First().Text()), " ")
	}
	subtitle := strings.Join(strings.Fields(s.Find(".text_container h2.film-subtitel").First().Text()), " ")
	if subtitle != "" && !strings.Contains(strings.ToLower(title), strings.ToLower(subtitle)) {
		title = title + " – " + subtitle
	}
	return title
}

func parseUnifilmReleaseYear(title string) null.String {
	m := unifilmTitleYearRegex.FindStringSubmatch(title)
	if len(m) >= 2 {
		return null.StringFrom(m[1])
	}
	return null.String{}
}

func parseUnifilmScreeningTimeAndLocation(s *goquery.Selection, loc *time.Location) (time.Time, string, error) {
	dateText := strings.TrimSpace(s.Find(".film-info-termin .datum").First().Text())
	timeText := strings.TrimSpace(s.Find(".film-info-termin .uhrzeit").First().Text())
	roomText := strings.Join(strings.Fields(s.Find(".film-info-termin .raum").First().Text()), " ")

	dateMatch := unifilmDateRegex.FindStringSubmatch(dateText)
	if len(dateMatch) != 4 {
		return time.Time{}, "", errors.New("campuskino heilbronn screening date missing or invalid")
	}
	timeMatch := unifilmTimeRegex.FindStringSubmatch(timeText)
	clock := "20:30"
	if len(timeMatch) >= 2 {
		clock = timeMatch[1]
	}

	parsed, err := time.ParseInLocation("2.1.2006 15:04", fmt.Sprintf("%s.%s.%s %s", dateMatch[1], dateMatch[2], dateMatch[3], clock), loc)
	if err != nil {
		return time.Time{}, "", err
	}

	if roomText == "" {
		roomText = "CampusKino Heilbronn"
	} else if !strings.Contains(strings.ToLower(roomText), "heilbronn") {
		roomText = roomText + ", Heilbronn"
	}
	return parsed, roomText, nil
}

func parseUnifilmMetadata(s *goquery.Selection) (director, genre, runtime null.String) {
	s.Find("ul.film-info-filmdaten").Not(".film-auszeichnungen").Find("li").Each(func(_ int, li *goquery.Selection) {
		txt := strings.Join(strings.Fields(li.Text()), " ")
		if txt == "" {
			return
		}
		switch {
		case strings.HasPrefix(txt, "R:"):
			if d := strings.TrimSpace(strings.TrimPrefix(txt, "R:")); d != "" {
				director = null.StringFrom(d)
			}
		case unifilmRuntimeRegex.MatchString(txt):
			m := unifilmRuntimeRegex.FindStringSubmatch(txt)
			if len(m) >= 2 {
				runtime = null.StringFrom(m[1] + " min")
			}
		case unifilmFSKRegex.MatchString(txt):
			// ignore FSK age rating
		case unifilmCountryCodeList.MatchString(txt):
			// ignore country code list (e.g. "USA", "DEU, USA")
		default:
			if !genre.Valid {
				genre = null.StringFrom(txt)
			}
		}
	})
	return director, genre, runtime
}

func parseUnifilmDescription(s *goquery.Selection, bm *bluemonday.Policy) string {
	var sections []string

	// Optional notice span right after film-info-termin (e.g. "Terminverschiebung: ...")
	s.Find(".text_container > span").Each(func(_ int, span *goquery.Selection) {
		notice := strings.Join(strings.Fields(bm.Sanitize(span.Text())), " ")
		if notice != "" {
			sections = append(sections, fmt.Sprintf("<b>%s<b>", notice))
		}
	})

	s.Find(".text_container > p").Each(func(_ int, p *goquery.Selection) {
		paragraph := strings.Join(strings.Fields(bm.Sanitize(p.Text())), " ")
		if paragraph != "" {
			sections = append(sections, paragraph)
		}
	})

	result := strings.Join(sections, "\n\n")
	if result == "" {
		return "Surprise yourself"
	}
	return result
}

func parseUnifilmImageURL(s *goquery.Selection, bm *bluemonday.Policy) string {
	src, exists := s.Find("img").First().Attr("src")
	if !exists || strings.TrimSpace(src) == "" {
		return unifilmDefaultImageFallback
	}
	sanitized := bm.Sanitize(strings.TrimSpace(src))
	base, _ := url.Parse(unifilmBaseURL)
	ref, err := url.Parse(sanitized)
	if err != nil {
		return unifilmDefaultImageFallback
	}
	return base.ResolveReference(ref).String()
}
