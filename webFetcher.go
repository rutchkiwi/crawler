package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type WebFetcher struct{}

func (f WebFetcher) Fetch(url string) (string, []string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, fmt.Errorf("not found: %s", url)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyStr := string(bodyBytes[:])

	return bodyStr, findUrls(bodyStr), nil
}

func findUrls(html string) []string {
	r, _ := regexp.Compile(`href="(http://.*?)"`)
	matches := r.FindAllStringSubmatch(html, -1)

	res := make([]string, len(matches))
	for i, match := range matches {
		res[i] = match[1]
	}

	return res
}
