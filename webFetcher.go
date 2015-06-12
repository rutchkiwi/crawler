package main

import (
	"io/ioutil"
	"net/http"
)

type WebFetcher struct{}

func (f WebFetcher) Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil { //todo: correct?
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyStr := string(bodyBytes[:])

	return bodyStr, nil
}
