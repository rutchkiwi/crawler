package main

import (
	"fmt"
	"regexp"
)

type queue struct{}

func crawl(seedUrl string, fetcher Fetcher) chan string {
	nonVisitedChan := make(chan string, 2)
	resultsChan := make(chan string, 2)

	nonVisitedChan <- seedUrl
	go processUrls(nonVisitedChan, fetcher)
	return resultsChan
}

func processUrls(nonVisitedChan <-chan string, fetcher Fetcher) {
	for url := range nonVisitedChan {
		fmt.Println(url)
		fmt.Println("111")

		body, _ := fetcher.Fetch(url) //handle error!

		// fmt.Printf("body: %v", body)
		urls := findUrls(body)
		for _, url := range urls {
			fmt.Println(url)
		}
	}
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
