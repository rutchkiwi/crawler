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
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	return WebFetcher{client}
}

func (f WebFetcher) Fetch(url string) (string, error) {
	// todo: needs a timeout
	resp, err := f.client.Get(url)
	if err != nil { //todo: correct?
		var err2 error
		resp, err2 = f.client.Get(url)
		if err2 != nil { //todo: correct?
			return "", err2
		}
		fmt.Println("RECOVERED!")
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyStr := string(bodyBytes[:])

	return bodyStr, nil
}
