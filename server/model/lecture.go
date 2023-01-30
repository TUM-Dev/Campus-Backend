package model

import "encoding/xml"

type Lectures struct {
	XMLName  xml.Name  `xml:"rowset"`
	Lectures []Lecture `xml:"row"`
}

type Lecture struct {
	XMLName           xml.Name `xml:"row"`
	Id                string   `xml:"stp_sp_nr"`
	LectureNumber     string   `xml:"stp_lv_nr"`
	LectureTitle      string   `xml:"stp_lv_titel"`
	Duration          int      `xml:"dauer_info"`
	LectureType       string   `xml:"stp_lv_art_name"`
	LectureTypeShort  string   `xml:"stp_lv_art_kurz"`
	SemesterYear      string   `xml:"sj_name"`
	SemesterShort     string   `xml:"semester"`
	SemesterName      string   `xml:"semester_name"`
	SemesterId        string   `xml:"semester_id"`
	OrganisationId    string   `xml:"org_nr_betreut"`
	OrganisationName  string   `xml:"org_name_betreut"`
	OrganisationShort string   `xml:"org_kennung_betreut"`
	Lecturer          string   `xml:"vortragende_mitwirkende"`
}

func (lecture *Lecture) TakesPlaceInWinterSemester() bool {
	return lecture.SemesterShort == "W"
}

func (lecture *Lecture) TakesPlaceInSummerSemester() bool {
	return lecture.SemesterShort == "S"
}
