package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawling(t *testing.T) {
	var results chan string
	results = crawl("http://a.com/", fakeFetcher)
	r1 := <-results
	assert.Equal(t, "http://a.com/", r1)
	r2 := <-results
	assert.Equal(t, "http://b.com/", r2)

}

// fetcher is a populated fakeFetcher.
var fakeFetcher = FakeFetcher{
	"http://a.com/": `<a href="http://b.com/">link text</a>`,
	"http://b.com/": `end`,
}
