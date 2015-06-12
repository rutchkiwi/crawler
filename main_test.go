package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawlsRightUrls(t *testing.T) {
	var results chan SiteInfo
	results = crawl("http://a.com", fakeFetcher)

	r1 := <-results
	assert.Equal(t, "http://a.com", r1.url)
	r2, _ := <-results
	assert.Equal(t, "http://b.com", r2.url)
	_, more := <-results
	assert.False(t, more)
}

func TestFindImageAssets(t *testing.T) {
	var results chan SiteInfo
	results = crawl("http://a.com", fakeFetcher)

	r1 := <-results
	assert.Equal(t, "http://a.com", r1.url)
	assert.Len(t, r1.assets, 1)
	assert.Equal(t, "picA.png", r1.assets[0])
	r2, _ := <-results
	assert.Equal(t, "http://b.com", r2.url)
	assert.Len(t, r2.assets, 1)
	assert.Equal(t, "picB.png", r2.assets[0])
	_, more := <-results
	assert.False(t, more)
}

// fetcher is a populated fakeFetcher.
var fakeFetcher = FakeFetcher{
	"http://a.com": `<a href="http://b.com"><img src="picA.png"/></a>`,
	"http://b.com": `<a href="http://a.com"></a><img src="picB.png"/><a href="http://b.com">link text</a>`,
}
