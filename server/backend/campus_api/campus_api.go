// Package campus_api handles all requests to the TUM Campus API and decodes the XML responses.
package campus_api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

// FetchExamResultsPublished fetches all published exam results from the TUM Campus API using CAMPUS_API_TOKEN.
func FetchExamResultsPublished() (*model.TUMAPIPublishedExamResults, error) {
	var examResultsPublished model.TUMAPIPublishedExamResults
	token := os.Getenv("CAMPUS_API_TOKEN")
	err := RequestCampusApi("/wbservicesbasic.pruefungenErgebnisse", token, &examResultsPublished)
	if err != nil {
		return nil, err
	}

	return &examResultsPublished, nil
}

func FetchGrades(token string) (*model.IOSGrades, error) {
	var grades model.IOSGrades
	err := RequestCampusApi("/wbservicesbasic.noten", token, &grades)
	if err != nil {
		return nil, err
	}

	return &grades, nil
}

func RequestCampusApi(path string, token string, response any) error {
	requestUrl := fmt.Sprintf("https://campus.tum.de/tumonline%s?pToken=%s", path, token)
	resp, err := http.Get(requestUrl)
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

	if err = xml.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.WithError(err).WithField("path", path).Error("Error while unmarshalling")
		return errors.New("error while unmarshalling")
	}

	return nil
}
