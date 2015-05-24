package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawling(t *testing.T) {
	urls, _ := Crawl("http://golang.org/", 7, fakeFetcher)
	assert.Equal(t, 4, len(urls), fmt.Sprintf("all urls should be visited. Visited %s", urls))
}

func TestCrawlingErrors(t *testing.T) {
	_, errors := Crawl("http://golang.org/", 7, fakeFetcher)

	assert.Equal(t, 2, len(errors), fmt.Sprintf("all errors should be returned"))
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
			"http://banana",
		},
	},
}
