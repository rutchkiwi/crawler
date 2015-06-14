package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type queue struct{}

func crawl(seedUrl string, fetcher Fetcher) (chan SiteInfo, chan error) {
	nonVisited := make(chan string, 1000000) // todo should be lower
	resultsChan := make(chan []string)       // todo should be lower

	outputAssets := make(chan SiteInfo)     // todo should be lower
	httpErrors := make(chan error, 1000000) // todo should be lower

	go dispatcher(nonVisited, resultsChan, outputAssets, seedUrl, httpErrors)

	for i := 0; i < 30; i++ {
		go processUrls(nonVisited, fetcher, resultsChan, outputAssets, httpErrors)
	}

	return outputAssets, httpErrors
}

func dispatcher(nonVisited chan string, resultsChan chan []string, outputAssets chan<- SiteInfo, seedUrl string, httpErrors chan error) {
	// because we added a url before
	seen := make(map[string]bool)
	var wg sync.WaitGroup

	wg.Add(1)

	rootUrl, err := url.Parse(seedUrl)
	if err != nil {
		log.Fatal(err)
	}

	seen[rootUrl.String()] = true
	nonVisited <- rootUrl.String()
	//log.Println("Will only allow urls on domain " + allowedDomain)

	// close channels when job counter reaches zero
	go func() {
		wg.Wait()
		close(nonVisited)
		close(resultsChan)
		close(outputAssets)
		close(httpErrors)
	}()

	for res := range resultsChan {
		for _, newUrl := range res {
			parsedUrl, err := url.Parse(newUrl)
			if err != nil {
				log.Printf("could not parse url %s", newUrl)
				continue
			}
			if !parsedUrl.IsAbs() {
				parsedUrl.Host = rootUrl.Host
				parsedUrl.Scheme = rootUrl.Scheme
			}
			if (!seen[parsedUrl.String()]) && (parsedUrl.Host == rootUrl.Host) {
				//fmt.Printf("new url %s found \n", parsedUrl.String())
				seen[parsedUrl.String()] = true
				wg.Add(1)
				nonVisited <- parsedUrl.String()
			}
		}
		wg.Done()
	}

}

type SiteInfo struct {
	url    string
	assets []string
}

func (info SiteInfo) String() string {
	return fmt.Sprintf("%s has assets: %s", info.url, strings.Join(info.assets, ", "))
}

func processUrls(nonVisited <-chan string, fetcher Fetcher, resultsChan chan<- []string, outputAssets chan<- SiteInfo, errors chan error) {
	for url := range nonVisited {
		// todo do assets and url search async ?

		//fmt.Printf("processing url %s \n", url)

		body, err := fetcher.Fetch(url) //handle error!
		if err != nil {
			em := make([]string, 0)
			// we need to send something so that the dispatcher knows we are done with this job.
			resultsChan <- em
			//fmt.Printf("ERROR for %s %v", url, err)
			errors <- err
			continue
		}

		assets := findAssets(body)
		info := SiteInfo{url, assets}

		// bad name outputAssets
		outputAssets <- info

		urls := findUrls(body)
		resultsChan <- urls
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
