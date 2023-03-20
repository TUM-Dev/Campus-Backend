package lecture_crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/TUM-Dev/Campus-Backend/server/backend/utils"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

func GetLecturesFromURLs(urls []string) (*[]model.CrawlerLecture, error) {
	var lectures []model.CrawlerLecture

	utils.RunTasksInRoutines(&urls, func(url string) {
		lecture, err := GetLectureFromURL(url)
		if err != nil {
			log.Errorf("Error while getting lecture from url: %v", err)
			return
		}

		lectures = append(lectures, *lecture)
	}, MaxRoutineCount)

	return &lectures, nil
}

func GetLectureFromURL(url string) (*model.CrawlerLecture, error) {
	lecture := model.CrawlerLecture{}

	params, err := ParseUrlParams(url)
	if err != nil {
		log.Errorf("Error while parsing detail url params: %v", err)
		return nil, err
	}

	lecture.Id = params.Get("pKnotenNr")

	resp, err := MakeRequest(url)
	if err != nil {
		log.Errorf("Error while making detail request: %v", err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(*resp)
	if err != nil {
		log.Errorf("Error while parsing detail response: %v", err)
		return nil, err
	}

	doc.Find("label").
		Each(func(i int, s *goquery.Selection) {
			forId, ok := s.Attr("for")

			if !ok {
				return
			}

			labelTrimmed := strings.Trim(s.Text(), " ")

			valueForLabel := findValueForLabelId(s, forId)

			switch labelTrimmed {
			case "Name":
				lecture.Name = valueForLabel
			case "Modul-Kennung":
				lecture.LectureId = valueForLabel
			case "ECTS-Credits":
				lecture.ECTS = valueForLabel
			case "Organisation":
				lecture.OrganisationName = valueForLabel
			case "Organisationskennung":
				lecture.OrganisationId = valueForLabel
			case "Gewichtungsfaktor":
				lecture.Weight = valueForLabel
			}
		})

	return &lecture, nil
}

func findValueForLabelId(s *goquery.Selection, labelId string) string {
	spanElement := fmt.Sprintf("span[id=%s]", labelId)
	return s.Parent().Parent().Find(spanElement).Text()
}
