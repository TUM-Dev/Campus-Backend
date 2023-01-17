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
	ErrWhileReadingBody     = errors.New("error while reading response body")
	ErrorWhileUnmarshalling = errors.New("error while unmarshalling")
)

func FetchGrades(token string) (*model.IOSGrades, error) {

	requestUrl := CampusApiUrl + CampusGradesPath
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)

	if err != nil {
		log.Errorf("Error while creating request: %s", err)
		return nil, ErrCannotCreateRequest
	}

	q := req.URL.Query()
	q.Add(CampusQueryToken, token)

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Errorf("Error while fetching grades: %s", err)
		return nil, ErrWhileFetchingGrades
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Error while closing body: %s", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("Error while reading response body: %s", err)
		return nil, ErrWhileReadingBody
	}

	var grades model.IOSGrades

	if err := xml.Unmarshal(body, &grades); err != nil {
		log.Errorf("Error while unmarshalling grades: %s", err)
		return nil, ErrorWhileUnmarshalling
	}

	return &grades, nil
}
