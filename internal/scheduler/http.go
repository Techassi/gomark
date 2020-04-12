package scheduler

import (
	"errors"
	"net/http"
	"net/url"
)

var (
	ErrorHTTPStatusCode = errors.New("HTTP status code greater than 203")
)

func (s *Scheduler) fetch(u string) (*http.Response, error) {
	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method:     "GET",
		URL:        url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return res, err
	}

	if res.StatusCode > 203 {
		return res, ErrorHTTPStatusCode
	}

	return res, nil
}
