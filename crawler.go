package main

import (
	"fmt"
	"sync"
)

type visitedUrls struct {
	sync.RWMutex
	m map[string]bool
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type crawlerShared struct {
	fetcher     Fetcher
	results     chan string
	visitAccess chan map[string]bool
	wg          *sync.WaitGroup
	errors      chan error
}

// crawlHelper uses fetcher to recursively crawlHelper
// pages starting with url, to a maximum of depth.
func crawlHelper(url string, depth int, shared *crawlerShared) {
	defer shared.wg.Done()

	if depth <= 0 {
		return
	}

	visited := <-shared.visitAccess
	if visited[url] {
		fmt.Println("already visited", url)
		shared.visitAccess <- visited
		return
	}
	visited[url] = true
	shared.visitAccess <- visited

	shared.results <- url

	_, urls, err := shared.fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, u := range urls {
		shared.wg.Add(1)
		go crawlHelper(u, depth-1, shared)
	}
	return
}

func Crawl(url string, depth int, fetcher Fetcher) ([]string, []error) {
	var results chan string = make(chan string, 10)

	// channel to serialize access to the list of visited sites
	visitAccess := make(chan map[string]bool, 1)
	visitAccess <- make(map[string]bool)

	// wait group used to shut down the result channel when
	// all goroutines are completed.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(results chan string, wg *sync.WaitGroup) {
		wg.Wait()
		close(results)
	}(results, wg)

	errors := make(chan error)
	shared := &crawlerShared{fetcher, results, visitAccess, wg, errors}

	go crawlHelper(url, depth, shared)

	ret := make([]string, 0)
	count := 0
	for url := range results {
		count++
		fmt.Printf("Crawled url #%d: %s\n", count, url)
		ret = append(ret, url)
	}
	return ret, nil
}

func main() {
	Crawl("http://golang.org/", 6, WebFetcher{})
}
