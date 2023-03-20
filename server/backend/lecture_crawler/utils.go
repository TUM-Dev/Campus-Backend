package lecture_crawler

import (
	"io"
	"net/http"
	"net/url"
)

func MakeRequest(url string) (*io.ReadCloser, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, err
	}

	return &resp.Body, nil
}

func ParseUrlParams(rawUrl string) (*url.Values, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	return &params, nil
}
