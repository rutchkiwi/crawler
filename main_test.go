package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawlsRightUrls(t *testing.T) {
	var results chan siteInfo
	results = crawl("http://a.com", fakeFetcher)
	r1 := <-results
	assert.Equal(t, "http://a.com", r1.url)
	r2 := <-results
	assert.Equal(t, "http://b.com", r2.url)
	// check empty now
}

// fetcher is a populated fakeFetcher.
var fakeFetcher = FakeFetcher{
	"http://a.com": `<a href="http://b.com"><img src="picA.png"</a>`,
	"http://b.com": `<a href="http://a.com"></a><img src="picB.png"/><a href="http://b.com">link text</a>`,
}
