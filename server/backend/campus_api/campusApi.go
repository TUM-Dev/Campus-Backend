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
	CampusApiUrl     = "https://campus.tum.de/tumonline"
	CampusQueryToken = "pToken"
	CampusGradesPath = "/wbservicesbasic.noten"
)

var (
	ErrCannotCreateRequest  = errors.New("cannot create http request")
	ErrWhileFetchingGrades  = errors.New("error while fetching grades")
	ErrorWhileUnmarshalling = errors.New("error while unmarshalling")
)

func FetchGrades(token string) (*model.IOSGrades, error) {

	requestUrl := CampusApiUrl + CampusGradesPath
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)

	if err != nil {
		log.WithError(err).Error("Failed to create api-request")
		return nil, ErrCannotCreateRequest
	}

	q := req.URL.Query()
	q.Add(CampusQueryToken, token)

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.WithError(err).Error("failed to fetch grades")
		return nil, ErrWhileFetchingGrades
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.WithError(err).Error("Could not close body")
		}
	}(resp.Body)

	var grades model.IOSGrades
	err = xml.NewDecoder(resp.Body).Decode(&grades)

	if err != nil {
		log.WithError(err).Error("could not unmarshall grades")
		return nil, ErrorWhileUnmarshalling
	}

	return &grades, nil
}
