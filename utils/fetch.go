// Package utils provides a fetch library and possibly other libs in the future
// to abstract away basic API responsibilities. Note that the response body is
// read in this util, and therefore will be closed upon return the *http.Response
// struct.
package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Fetch optimized 10kft fetch helper
func (opts FetchOpts) Fetch() (resp *http.Response, err error) {
	c := &http.Client{}
	payload := strings.NewReader(opts.Body)

	req, err := http.NewRequest(opts.Method, opts.URL, payload)
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	for key, value := range opts.Headers {
		req.Header.Add(key, value)
	}

	resp, err = c.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode == 429 && opts.MaxRetries > 0 {
		opts.MaxRetries--
		time.Sleep(time.Second * 10)
		resp, err = opts.Fetch()
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if opts.MaxRetries > 0 {
			opts.MaxRetries--
			time.Sleep(time.Second * 2)
			resp, err = opts.Fetch()
		} else {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				err = fmt.Errorf("Non OK status code %v and could not parse response text", resp.StatusCode)
				return resp, err
			}

			err = fmt.Errorf("Non OK status Code: %v, body: %v", resp.StatusCode, string(b))

			resp.Body.Close()

			return resp, err
		}
	}

	return
}

// NewFetchOpts opts
func NewFetchOpts(url, method, body string, headers map[string]string, maxRetries int) (FetchOpts, error) {
	var err error
	opts := FetchOpts{}
	if url == "" {
		err = errors.New("URL cannot be empty")
		return opts, err
	}
	opts.URL = url

	if method == "" {
		method = "GET"
	}
	opts.Method = method

	opts.Body = body

	opts.Headers = headers

	opts.MaxRetries = maxRetries

	return opts, nil
}

// FetchOpts optimized 10kft fetch helper
type FetchOpts struct {
	URL        string
	Method     string
	Body       string
	Headers    map[string]string
	MaxRetries int
}
