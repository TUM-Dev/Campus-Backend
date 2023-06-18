// Package campus_api handles all requests to the TUM Campus API and decodes the XML responses.
package campus_api

import (
	"encoding/xml"
	"errors"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const (
	CampusApiUrl               = "https://campus.tum.de/tumonline"
	CampusQueryToken           = "pToken"
	CampusGradesPath           = "/wbservicesbasic.noten"
	CampusExamResultsPublished = "/wbservicesbasic.pruefungenErgebnisse"
)

var (
	ErrCannotCreateRequest  = errors.New("cannot create http request")
	ErrorWhileUnmarshalling = errors.New("error while unmarshalling")
)

func FetchExamResultsPublished(token string) (*model.TUMAPIExamResultsPublished, error) {
	var examResultsPublished model.TUMAPIExamResultsPublished
	err := RequestCampusApi(CampusExamResultsPublished, token, &examResultsPublished)
	if err != nil {
		return nil, err
	}

	return &examResultsPublished, nil
}

func FetchGrades(token string) (*model.IOSGrades, error) {
	var grades model.IOSGrades
	err := RequestCampusApi(CampusGradesPath, token, &grades)
	if err != nil {
		return nil, err
	}

	return &grades, nil
}

func RequestCampusApi(path string, token string, response any) error {
	requestUrl := "https://exams.free.beeceptor.com/"
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)

	if err != nil {
		log.Errorf("Error while creating request: %s", err)
		return ErrCannotCreateRequest
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Errorf("Error while fetching %s: %s", path, err)
		return errors.New("error while fetching " + path)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Error while closing body: %s", err)
		}
	}(resp.Body)

	err = xml.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		log.Errorf("Error while unmarshalling %s: %s", path, err)
		return ErrorWhileUnmarshalling
	}

	return nil
}
