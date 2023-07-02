// Package campus_api handles all requests to the TUM Campus API and decodes the XML responses.
package campus_api

import (
	"encoding/xml"
	"errors"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

var (
	CampusApiUrl = url.URL{
		Scheme: "https",
		Host:   "campus.tum.de",
		Path:   "/tumonline",
	}
)

const (
	CampusQueryToken           = "pToken"
	CampusGradesPath           = "/wbservicesbasic.noten"
	CampusExamResultsPublished = "/wbservicesbasic.pruefungenErgebnisse"
	CampusPersonalExamsPath    = "/wbservicesbasic.pruefungenEigene"
)

var (
	ErrCannotCreateRequest  = errors.New("cannot create http request")
	ErrorWhileUnmarshalling = errors.New("error while unmarshalling")
)

func FetchExamResultsPublished(token string) (*model.TUMAPIExamResultsPublished, error) {
	var examResultsPublished model.TUMAPIExamResultsPublished

	url, _ := url.Parse("https://exams2.free.beeceptor.com")

	err := RequestCampusApiWithBaseUrl(url, "", token, &examResultsPublished)
	if err != nil {
		return nil, err
	}

	return &examResultsPublished, nil
}

func FetchGrades(token string) (*model.Grades, error) {
	var grades model.Grades
	err := RequestCampusApi(CampusGradesPath, token, &grades)
	if err != nil {
		return nil, err
	}

	return &grades, nil
}

func FetchPersonalExams(token string) (*model.Exams, error) {
	var exams model.Exams
	err := RequestCampusApi(CampusPersonalExamsPath, token, &exams)

	if err != nil {
		return nil, err
	}

	return &exams, nil
}

func RequestCampusApi(path string, token string, response any) error {
	return RequestCampusApiWithBaseUrl(&CampusApiUrl, path, token, response)
}

func RequestCampusApiWithBaseUrl(baseUrl *url.URL, path string, token string, response any) error {
	requestUrl := baseUrl.JoinPath(path)

	query := requestUrl.Query()
	query.Add(CampusQueryToken, token)

	requestUrl.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, requestUrl.String(), nil)

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
