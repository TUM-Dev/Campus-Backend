package lecture_crawler

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const BaseURL = "https://campus.tum.de/tumonline"
const MaxRoutineCount = 25

type Crawler struct {
	DB *gorm.DB
}

func (crawler *Crawler) Crawl() {
	maxPageNumber, err := GetMaxPageNumber()
	if err != nil {
		log.Errorf("Error while getting max page number: %v", err)
		return
	}

	log.Infof("Start crawling %d pages. Using %d threads.", maxPageNumber, MaxRoutineCount)

	links, err := GetAllLectureLinks(maxPageNumber)
	if err != nil {
		log.Errorf("Error while getting all lecture links: %v", err)
		return
	}

	log.Infof("Crawled %d lecture page links", len(*links))
	log.Infof("Will now attempt to crawl %d lecture detail pages", len(*links))

	lectures, err := GetLecturesFromURLs(*links)
	if err != nil {
		log.Errorf("Error while getting lectures from urls: %v", err)
		return
	}

	log.Infof("Crawled %d lecture detail pages.", len(*lectures))
	log.Infof("Will now attempt to save %d lectures to the database.", len(*lectures))

	for _, lecture := range *lectures {
		crawler.DB.Save(&lecture)
	}
}

func New(db *gorm.DB) *Crawler {
	return &Crawler{
		DB: db,
	}
}
