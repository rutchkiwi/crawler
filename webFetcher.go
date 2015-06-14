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
	// 2 seconds seems good, could perhaps be optimized.
	timeout := time.Duration(2 * time.Second)
	transport := http.Transport{
		//TODO: make constant
		MaxIdleConnsPerHost: noHttpWorkers,
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
		// ghetto retry logic
		var err2 error
		resp, err2 = f.client.Get(url)
		if err2 != nil {
			fmt.Println("F")
			return "", err2
		}
		fmt.Println("r")
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(bodyBytes[:])

	return bodyStr, nil
}
