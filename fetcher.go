package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Fetcher interface {
	Fetch(url string) (body string, err error)
}

type WebFetcher struct {
	client http.Client
}

func newWebFetcher() WebFetcher {
	transport := http.Transport{
		// set up connection pooling
		MaxIdleConnsPerHost: noHttpWorkers,
	}
	client := http.Client{
		// 2 seconds arbitrarily chosen
		Timeout:   time.Duration(2 * time.Second),
		Transport: &transport,
	}
	return WebFetcher{client}
}

func (f WebFetcher) Fetch(url string) (string, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		// ghetto retry logic
		var err2 error
		resp, err2 = f.client.Get(url)
		if err2 != nil {
			return "", err2
		}
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(bodyBytes[:])

	return bodyStr, nil
}
