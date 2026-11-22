// Package webfetch provides bounded requests for external data sources.
package webfetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var Client = &http.Client{Timeout: 30 * time.Second}

// StatusError is returned when an upstream HTTP endpoint responds with a non-200 status.
type StatusError struct {
	StatusCode int
	Status     string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("upstream returned %s", e.Status)
}

// Get returns an OK response. The caller must close its body.
func Get(url string) (*http.Response, error) {
	resp, err := Client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, &StatusError{StatusCode: resp.StatusCode, Status: resp.Status}
	}
	return resp, nil
}

func JSON(url string, target any) error {
	resp, err := Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return json.NewDecoder(resp.Body).Decode(target)
}
