// Package campus_api handles all requests to the TUM Campus API and decodes the XML responses.
package campus_api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
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

func FetchExamResultsPublished(token string) (*model.TUMAPIPublishedExamResults, error) {
	var examResultsPublished model.TUMAPIPublishedExamResults
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
	requestUrl := fmt.Sprintf("%s%s?%s=%s", CampusApiUrl, path, CampusQueryToken, token)
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)

	if err != nil {
		log.WithError(err).Error("Error while creating request")
		return ErrCannotCreateRequest
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.WithError(err).WithField("path", path).Error("Error while fetching url")
		return errors.New("error while fetching " + path)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing body")
		}
	}(resp.Body)

	err = xml.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		log.WithError(err).WithField("path", path).Error("Error while unmarshalling")
		return ErrorWhileUnmarshalling
	}

	return nil
}
