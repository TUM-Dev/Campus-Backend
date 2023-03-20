package lecture_crawler

import (
	"github.com/PuerkitoBio/goquery"
)

func GetMaxPageNumber() (int, error) {
	url := BuildLecturesListURL(1)

	resp, err := MakeRequest(url)
	if err != nil {
		return 0, err
	}

	doc, err := goquery.NewDocumentFromReader(*resp)
	if err != nil {
		return 0, err
	}

	result := doc.Find("option")

	return result.Length() / 2, nil
}
