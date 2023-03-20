package model

type CrawlerLecture struct {
	Name             string
	LectureId        string
	ECTS             string
	Id               string `gorm:"primaryKey"`
	OrganisationName string
	OrganisationId   string
	Weight           string
}
