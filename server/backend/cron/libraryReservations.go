package cron

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

const LibraryUrl = "https://www.ub.tum.de/arbeitsplatz-reservieren?field_teilbibliothek_tid_selective=All&field_tag_value_selective=All&page=%d"

var errNoDataOnPage = errors.New("no matching table on page")

//libraryReservationsCron scans the tum bib website for open timeslots and saves them into the database
func (c CronService) libraryReservationsCron() error {
	for i := 0; i < 5; i++ {
		url := fmt.Sprintf(LibraryUrl, i)
		page, err := http.Get(url)
		if err != nil {
			return err
		}
		if page.StatusCode != http.StatusOK {
			log.WithFields(log.Fields{"url": url, "status": page.Status}).Warn("libraryReservationsCron: non 200 status")
			return nil
		}
		_, err = c.getBibsFromWebsite(page.Body)
		if errors.Is(err, errNoDataOnPage) {
			log.Debug("Reached last page on index %d", i)
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

//getBibsFromWebsite returns a slice of reservable bibDates
func (c *CronService) getBibsFromWebsite(r io.ReadCloser) ([]bibDate, error) {
	document, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	entries := document.Find("table.views-table.cols-4>tbody>tr")
	if entries.Length() == 0 {
		return nil, errNoDataOnPage // this indicates that the last page was reached
	}
	entries.Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find(".views-field-views-conditional>a").Attr("href")
		if exists {
			log.Println(strings.TrimSpace(s.Find(".views-field-field-teilbibliothek").Text()))
			log.Println(strings.TrimSpace(s.Find(".views-field-field-tag").Text()))
			log.Println(strings.TrimSpace(s.Find(".views-field-field-zeitslot").Text()))
			log.Println(strings.TrimSpace(link))
		}
	})
	return nil, nil
}

//bibDate represents a reservable bib date
type bibDate struct {
	Bib  string    // name of the bib (e.g. "Stammgelände")
	From time.Time // start of the time slot
	To   time.Time // start of the time slot
	Key  string    // key for the reservation (e.g. 12345), reservation link would be https://www.ub.tum.de/reserve/12345
}
