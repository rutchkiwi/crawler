package main

import (
	"fmt"
	"sync"
	"time"
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
	time.Sleep(300 * time.Millisecond)

	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
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

	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Printf("found: %s %q\n", url, body)
	ch <- url
	for _, u := range urls {
		wg.Add(1)
		go crawlHelper(u, depth-1, fetcher, ch, visitAccess, wg)
	}
	return
}

func Crawl(url string, depth int, fetcher Fetcher) {
	var ch chan string = make(chan string, 10)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	visitAccess := make(chan map[string]bool, 1)
	visitAccess <- make(map[string]bool)

	go crawlHelper(url, depth, fetcher, ch, visitAccess, wg)
	go func(ch chan string, wg *sync.WaitGroup) {
		wg.Wait()
		close(ch)
	}(ch, wg)

	for url := range ch {
		fmt.Printf("Crawled url: %s\n", url)
	}
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
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
