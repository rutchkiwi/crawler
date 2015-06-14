package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
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

	go dispatcher(urlQueue, foundLinksChannel, assetsChannel, *seedUrl, httpErrors)

	for i := 0; i < noHttpWorkers; i++ {
		go processUrls(urlQueue, fetcher, foundLinksChannel, assetsChannel, httpErrors)
	}

	return assetsChannel, httpErrors
}

func dispatcher(urlQueue chan string, foundLinksChannel chan []string, assetsChannel chan<- SiteInfo, seedUrl url.URL, httpErrors chan error) {

	visitedUrls := make(map[string]bool)

	// we need to keep track of this so we can close channels when there are no more urls to visit
	var jobsInProgress sync.WaitGroup

	visitedUrls[seedUrl.String()] = true
	jobsInProgress.Add(1)
	urlQueue <- seedUrl.String()

	// close channels when job counter reaches zero
	go func() {
		jobsInProgress.Wait()
		close(urlQueue)
		close(foundLinksChannel)
		close(assetsChannel)
		close(httpErrors)
	}()

	for res := range foundLinksChannel {
		for _, newUrl := range res {
			parsedUrl, err := url.Parse(newUrl)
			if err != nil {
				log.Printf("could not parse url %s", newUrl)
				continue
			}
			if !parsedUrl.IsAbs() {
				parsedUrl.Host = seedUrl.Host
				parsedUrl.Scheme = seedUrl.Scheme
			}
			if (!visitedUrls[parsedUrl.String()]) && (parsedUrl.Host == seedUrl.Host) {
				//fmt.Printf("new url %s found \n", parsedUrl.String())
				visitedUrls[parsedUrl.String()] = true
				jobsInProgress.Add(1)
				urlQueue <- parsedUrl.String()
			}
		}
		jobsInProgress.Done()
	}

}

type SiteInfo struct {
	url    string
	assets []string
}

func (info SiteInfo) String() string {
	return fmt.Sprintf("%s has assets: %s", info.url, strings.Join(info.assets, ", "))
}

func processUrls(urlQueue <-chan string, fetcher Fetcher, foundLinksChannel chan<- []string, assetsChannel chan<- SiteInfo, errors chan error) {
	for url := range urlQueue {
		// todo do assets and url search async ?

		//fmt.Printf("processing url %s \n", url)

		body, err := (fetcher).Fetch(url) //handle error!
		if err != nil {
			em := make([]string, 0)
			errors <- err
			// we need to send something so that the dispatcher knows we are done with this url.
			foundLinksChannel <- em
			continue
		}

		assets := findAssets(body)
		info := SiteInfo{url, assets}

		assetsChannel <- info

		urls := findUrls(body)
		foundLinksChannel <- urls
		//fmt.Printf("done processing url %s \n", url)

	}
}

//todo DONT PARSE WITH REGEX

func findAssets(html string) []string {
	r, _ := regexp.Compile(`src="(.*?)"`)
	matches := r.FindAllStringSubmatch(html, -1)

	res := make([]string, len(matches))
	for i, match := range matches {
		res[i] = match[1]
	}

	return res
}

func findUrls(html string) []string {
	r, _ := regexp.Compile(`href="(.*?)"`)
	matches := r.FindAllStringSubmatch(html, -1)

	res := make([]string, len(matches))
	for i, match := range matches {
		res[i] = match[1]
	}

	return res
}
