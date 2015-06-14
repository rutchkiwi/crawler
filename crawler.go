package main

import (
	"log"
	"net/url"
	"os"
	"sync"
)

const noHttpWorkers = 30

func crawl(seedUrlString string, fetcher Fetcher) (chan SiteInfo, chan error) {
	// these do not need to be buffered, but it might help a little bit with performance
	foundLinksChannel := make(chan []string, 10)
	assetsChannel := make(chan SiteInfo, 10)

	// these must have a large enough buffer to contain all seen urls or the program will deadlock
	urlQueue := make(chan string, 1000000)
	httpErrors := make(chan error, 1000000)

	seedUrl, err := url.Parse(seedUrlString)
	if err != nil {
		log.Fatal("Could not parse seed url %s", seedUrlString)
		os.Exit(1)
	}

	go dispatcher(*seedUrl, urlQueue, foundLinksChannel, assetsChannel, httpErrors)

	for i := 0; i < noHttpWorkers; i++ {
		go processUrls(urlQueue, fetcher, foundLinksChannel, assetsChannel, httpErrors)
	}

	return assetsChannel, httpErrors
}

func dispatcher(
	seedUrl url.URL,
	urlQueue chan string,
	foundLinksChannel chan []string,
	assetsChannel chan SiteInfo,
	httpErrors chan error) {

	visitedUrls := make(map[string]bool)

	// we need to keep track of this so we can close channels when there are no more urls to visit
	var jobsInProgress sync.WaitGroup

	visitedUrls[seedUrl.String()] = true
	jobsInProgress.Add(1)
	urlQueue <- seedUrl.String()

	// close channels when there are no more urls to visit
	go func() {
		jobsInProgress.Wait()
		close(urlQueue)
		close(foundLinksChannel)
		close(assetsChannel)
		close(httpErrors)
	}()

	// Pick up lists of urls from the foundLinksChannel,
	// check if they should be visited, and if so, puts them on the urlQueue.
	// This approach allows us to handle the list of visited sites in a single place, and
	// makes it easier to keep track of when all relevant urls are visited.
	for res := range foundLinksChannel {
		for _, newUrl := range res {
			parsedUrl, err := url.Parse(newUrl)
			if err != nil {
				// could not parse this url. Ignore it.
				continue
			}
			if !parsedUrl.IsAbs() {
				// handle relative urls
				parsedUrl.Host = seedUrl.Host
				parsedUrl.Scheme = seedUrl.Scheme
			}
			if (!visitedUrls[parsedUrl.String()]) && (parsedUrl.Host == seedUrl.Host) {
				visitedUrls[parsedUrl.String()] = true
				jobsInProgress.Add(1)
				urlQueue <- parsedUrl.String()
			}
		}
		jobsInProgress.Done()
	}
}

func processUrls(urlQueue <-chan string, fetcher Fetcher, foundLinksChannel chan<- []string, assetsChannel chan<- SiteInfo, errors chan error) {
	for url := range urlQueue {

		body, err := (fetcher).Fetch(url)
		if err != nil {
			em := make([]string, 0)
			errors <- err
			// we need to put something on the foundLinksChannel otherwise the dispatcher
			// will not know that we are done with this url.
			foundLinksChannel <- em
			continue
		}

		// No need to block the http fetching when doing the html parsing.
		// It's important to pass url and body as parameters, as this ensures they will be
		// copied by value into the spwned goroutine
		go func(url string, body string) {
			assets := findAssets(body)
			assetsChannel <- SiteInfo{url, assets}
			foundLinksChannel <- findUrls(body)
		}(url, body)
	}
}
