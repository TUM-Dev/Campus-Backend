package model

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Lectures struct {
	XMLName  xml.Name  `xml:"rowset"`
	Lectures []Lecture `xml:"row"`
}

type Lecture struct {
	XMLName       xml.Name `xml:"row"`
	Id            string   `xml:"stp_sp_nr"`
	LectureNumber string   `xml:"stp_lv_nr"`
	LectureTitle  string   `xml:"stp_sp_titel"`
	Duration      string   `xml:"dauer_info"`
	LectureType   string   `xml:"stp_lv_art_name"`
	// possible values: "VO", "TT", "UE", "PR", "SE"
	LectureTypeShort  string `xml:"stp_lv_art_kurz"`
	SemesterYear      string `xml:"sj_name"`
	SemesterShort     string `xml:"semester"`
	SemesterName      string `xml:"semester_name"`
	SemesterId        string `xml:"semester_id"`
	OrganisationId    string `xml:"org_nr_betreut"`
	OrganisationName  string `xml:"org_name_betreut"`
	OrganisationShort string `xml:"org_kennung_betreut"`
	Lecturer          string `xml:"vortragende_mitwirkende"`
}

func (lecture *Lecture) TakesPlaceInWinterSemester() bool {
	return lecture.SemesterShort == "W"
}

func (lecture *Lecture) TakesPlaceInSummerSemester() bool {
	return lecture.SemesterShort == "S"
}

func (lecture *Lecture) ToIOSLecture() (*IOSLecture, error) {
	semester := IOSLectureSemesterWinter
	if lecture.TakesPlaceInSummerSemester() {
		semester = IOSLectureSemesterSummer
	}

	yearString := lecture.SemesterId[:2]

	year, err := strconv.Atoi(yearString)
	if err != nil {
		log.WithError(err).Error("Could not parse year from semester id")
		return nil, err
	}

	log.Infof("Saving lecture %s in year %d and semester %d", lecture.LectureTitle, year, semester)

	return &IOSLecture{
		Id:       lecture.Id,
		Year:     int16(year),
		Semester: semester,
		Title:    lecture.LectureTitle,
	}, nil
}
