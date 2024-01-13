package util

import (
	"errors"
	"net/http"
	"strings"
)

type FetchOptions struct {
	Headers []string
}

func Fetch(method string, url string, options *FetchOptions) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, errors.New("request creation failed")
	}

	for _, reqHeader := range options.Headers {
		hdr := strings.SplitN(reqHeader, ":", 2)

		if len(hdr) != 2 {
			return nil, errors.New("invalid header")
		}

		req.Header.Set(strings.TrimSpace(hdr[0]), strings.TrimSpace(hdr[1]))
	}

	return client.Do(req)
}
