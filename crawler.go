package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type crawlerShared struct {
	fetcher     Fetcher
	results     chan string
	visitAccess chan map[string]bool
	errors      chan error
	queue       chan string
	wg          *sync.WaitGroup //remove
	statusChan  chan int
}

func handleOneUrl(url string, shared *crawlerShared) {
	visited := <-shared.visitAccess
	if visited[url] {
		//fmt.Println("already visited", url)
		shared.visitAccess <- visited
		return
	}
	fmt.Println("now working on " + string(url))
	visited[url] = true
	shared.visitAccess <- visited

	_, urls, err := shared.fetcher.Fetch(url)
	if err != nil {
		shared.errors <- err
		return
	}
	// Doesn't count as a result unless we could visit the site
	shared.results <- url

	for _, url := range urls {
		shared.queue <- url
	}
	return
}

// crawlHelper uses fetcher to recursively crawlHelper
// pages starting with url, to a maximum of depth.
func crawlHelper(shared *crawlerShared) {
	for {
		shared.statusChan <- -1
		select {
		case url := <-shared.queue:
			shared.statusChan <- 1
			handleOneUrl(url, shared)
		default:
			fmt.Println("default happend")
			return
		}
	}
}

func Crawl(url string, depth int, fetcher Fetcher) ([]string, []error) {
	results := make(chan string, 100000)
	errors := make(chan error, 1000000)
	queue := make(chan string, 1000000)
	statusChannel := make(chan string, 1000000)

	// channel to serialize access to the list of visited sites
	visitAccess := make(chan map[string]bool, 1)
	visitAccess <- make(map[string]bool)

	errorListChannel := make(chan []error, 2)

	numberOfWorkers := 2

	wg := &sync.WaitGroup{} //unnessecary

	// wg.Add(numberOfWorkers)
	// go func(results chan string, wg *sync.WaitGroup) {
	// 	wg.Wait()
	// 	close(results)
	// 	close(errors)
	// }(results, wg)

	shared := &crawlerShared{fetcher, results, visitAccess, errors, queue, wg, statusChannel}
	for i := 0; i < numberOfWorkers; i++ {
		go crawlHelper(shared)
	}

	shared.queue <- url

	// act on the results that should appear in the error and result channels.

	go func(errors chan error) {
		ret := make([]error, 0)
		errorCount := 0
		for err := range errors {
			errorCount++
			fmt.Printf("error: %v. Errors so far: %d\n", err, errorCount)
			ret = append(ret, err)
		}
		errorListChannel <- ret

	}(errors)

	ret := make([]string, 0)
	count := 0
	for url := range results {
		count++
		if count%1 == 0 {
			fmt.Printf("Crawled url #%d: %s\n", count, url)
		}
		ret = append(ret, url)
	}

	errorList := <-errorListChannel
	return ret, errorList
}

func main() {
	urls, errors := Crawl("http://golang.org/", 6, WebFetcher{})

	fmt.Printf("crawled: %d. Failed: %d\n", len(urls), len(errors))
}
