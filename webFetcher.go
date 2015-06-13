package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type WebFetcher struct{}

func (f WebFetcher) Fetch(url string) (string, error) {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	// todo: needs a timeout
	resp, err := client.Get(url)
	if err != nil { //todo: correct?
		fmt.Println("REQUEST TIMED OUT")
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyStr := string(bodyBytes[:])

	return bodyStr, nil
}
