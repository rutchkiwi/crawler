package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawling(t *testing.T) {
	var results chan string
	results = crawl("http://a.com/", fakeFetcher)
	r1 := <-results
	assert.Equal(t, "http://b.com/", r1)
}

// fetcher is a populated fakeFetcher.
var fakeFetcher = FakeFetcher{
	"http://a.com/": `<a href="http://b.com/"></a>`,
	"http://b.com/": `<a href="http://a.com/"></a> bla bla bla <a href="http://b.com/">link text</a>`,
}
