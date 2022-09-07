package LRZscraper

import (
	"fmt"
	"strings"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type AP struct {
	address string
	room    string
	name    string
	status  string
	typ     string
	load    string
}

func Nothing() {

}

// It scrapes the html table data from "https://wlan.lrz.de/apstat/""
func ScrapeApstat() []AP {
	apstatURL := "https://wlan.lrz.de/apstat/"

	c := colly.NewCollector(
		colly.AllowedDomains("wlan.lrz.de"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	accessPoints := make([]AP, 0)

	// it uses jQuery selectors to scrape the table with id "aptable" row by row
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.DOM.Find("table#aptable > tbody > tr").Each(func(i int, s *goquery.Selection) {
			// scrape overview data for each access point
			address := s.ChildrenFiltered("td:nth-child(1)").Text()
			room := s.ChildrenFiltered("td:nth-child(2)").Text()
			apName := s.ChildrenFiltered("td:nth-child(3)").Text()
			apStatus := s.ChildrenFiltered("td:nth-child(4)").Text()
			apStatus = strings.TrimSpace(apStatus)
			apType := s.ChildrenFiltered("td:nth-child(5)").Text()
			load := s.ChildrenFiltered("td:nth-child(6)").Text()

			ap := AP{address, room, apName, apStatus, apType, load}
			accessPoints = append(accessPoints, ap)
		})
	})

	c.Visit(apstatURL)
	c.Wait()

	return accessPoints
}

func getTotAvg(load string) string {
	rgx := regexp.MustCompile(`\((.*?)\)`)
	loadStr := rgx.FindStringSubmatch(load)
	loads := strings.Split(loadStr[1], " - ")
	totAvg := loads[3]
	return totAvg
}

