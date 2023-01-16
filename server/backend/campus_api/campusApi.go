package campus_api

import (
	"encoding/xml"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const (
	CAMPUS_API_URL = "https://campus.tum.de/tumonline"
)

func FetchGrades(token string) (*model.IOSGrades, error) {

	req, err := http.NewRequest(http.MethodGet, CAMPUS_API_URL+"/wbservicesbasic.noten", nil)

	if err != nil {
		log.Errorf("Error while creating request: %s", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("pToken", token)

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Errorf("Error while fetching grades: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("Error while reading response body: %s", err)
		return nil, err
	}

	var grades model.IOSGrades

	xml.Unmarshal(body, &grades)

	return &grades, nil
}
