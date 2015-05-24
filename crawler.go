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

// crawlHelper uses fetcher to recursively crawlHelper
// pages starting with url, to a maximum of depth.
func crawlHelper(url string, depth int, fetcher Fetcher, ch chan string, visitAccess chan map[string]bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	visited := <-visitAccess
	if visited[url] {
		fmt.Println("already visited", url)
		visitAccess <- visited
		return
	}
	visited[url] = true
	visitAccess <- visited

	fmt.Println("now visiting url", url)

	ch <- url

	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, u := range urls {
		wg.Add(1)
		go crawlHelper(u, depth-1, fetcher, ch, visitAccess, wg)
	}
	return
}

func Crawl(url string, depth int, fetcher Fetcher) {
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

	go crawlHelper(url, depth, fetcher, results, visitAccess, wg)

	for url := range results {
		fmt.Printf("Crawled url: %s\n", url)
	}
}

func main2() {
	Crawl("http://golang.org/", 4, WebFetcher{})
}

// fetcher is a populated fakeFetcher.
var fakeFetcher = FakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
