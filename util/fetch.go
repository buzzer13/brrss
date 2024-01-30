package util

import (
	"errors"
	"net/http"
	"time"
)

type FetchOptions struct {
	Header http.Header
}

func Fetch(method string, url string, options *FetchOptions) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, errors.New("request creation failed")
	}

	for header := range options.Header {
		req.Header.Set(header, options.Header.Get(header))
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return client.Do(req)
}
