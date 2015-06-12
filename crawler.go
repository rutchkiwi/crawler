package main

import (
	"fmt"
	"regexp"
	"sync"
)

type queue struct{}

func crawl(seedUrl string, fetcher Fetcher) chan string {
	nonVisited := make(chan string, 100)    // todo should be lower
	resultsChan := make(chan []string, 100) // todo should be lower

	outputUrls := make(chan string, 100) // todo should be lower

	go dispatcher(nonVisited, resultsChan, outputUrls, seedUrl)
	go processUrls(nonVisited, fetcher, resultsChan)
	return outputUrls
}

func dispatcher(nonVisited chan string, resultsChan chan []string, output chan<- string, seedUrl string) {
	// because we added a url before
	seen := make(map[string]bool)
	var wg sync.WaitGroup

	wg.Add(1)
	seen[seedUrl] = true
	nonVisited <- seedUrl

	// close channels when job counter reaches zero
	go func() {
		wg.Wait()
		close(nonVisited)
		close(resultsChan)
		close(output)
	}()

	for res := range resultsChan {
		for _, url := range res {
			if !seen[url] {
				fmt.Printf("new url %s found \n", url)
				seen[url] = true
				wg.Add(1)
				nonVisited <- url
				output <- url
			}
		}
		wg.Done()
	}

}

func processUrls(nonVisited <-chan string, fetcher Fetcher, resultsChan chan<- []string) {
	for url := range nonVisited {
		fmt.Printf("processing url %s \n", url)

		body, _ := fetcher.Fetch(url) //handle error!

		// fmt.Printf("body: %v", body)
		urls := findUrls(body)
		resultsChan <- urls
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
