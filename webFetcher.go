package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type WebFetcher struct {
	client http.Client
}

func newWebFetcher() WebFetcher {
	// 2 seconds seems like a nice number
	timeout := time.Duration(2 * time.Second)
	transport := http.Transport{
		//TODO make constant
		MaxIdleConnsPerHost: 30,
	}
	client := http.Client{
		Timeout:   timeout,
		Transport: &transport,
	}
	return WebFetcher{client}
}

func (f WebFetcher) Fetch(url string) (string, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		// getto retry logic
		var err2 error
		resp, err2 = f.client.Get(url)
		if err2 != nil {
			return "", err2
		}
		fmt.Println("RECOVERED!")
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(bodyBytes[:])

	return bodyStr, nil
}
