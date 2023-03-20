package lecture_crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/TUM-Dev/Campus-Backend/server/backend/utils"
	log "github.com/sirupsen/logrus"
)

func GetAllLectureLinks(maxPages int) (*[]string, error) {
	var links []string

	utils.RunXTasksInRoutines(maxPages, func(i int) {
		linksOfPage, err := GetAllLectureLinksOfPage(i + 1)
		if err != nil {
			log.Errorf("Error while getting links of page %d: %v", i+1, err)
			return
		}

		links = append(links, *linksOfPage...)
	}, MaxRoutineCount)

	return &links, nil
}

func GetAllLectureLinksOfPage(pageNumber int) (*[]string, error) {
	resp, err := MakeRequest(BuildLecturesListURL(pageNumber))
	if err != nil {
		log.Errorf("Error while making lectures request: %v", err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(*resp)
	if err != nil {
		log.Errorf("Error while parsing lectures response: %v", err)
		return nil, err
	}

	var links []string

	doc.Find("tr.coRow td a[href]").
		FilterFunction(func(i int, s *goquery.Selection) bool {
			return s.Text() != ""
		}).
		Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			fullLink := GenerateLinkForKnotenNr(link)

			links = append(links, fullLink)
		})

	return &links, nil
}

func GenerateLinkForKnotenNr(knotenNr string) string {
	return BaseURL + "/" + knotenNr
}
