package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawlsRightUrls(t *testing.T) {
	var results chan SiteInfo
	results = crawl("http://gocardless.com/a", fakeFetcher1)

	r1 := <-results
	assert.Equal(t, "http://gocardless.com/a", r1.url)
	r2, _ := <-results
	assert.Equal(t, "https://gocardless.com/b", r2.url)
	_, more := <-results
	assert.False(t, more)
}

func TestFindImageAssets(t *testing.T) {
	var results chan SiteInfo
	results = crawl("http://gocardless.com/a", fakeFetcher1)

	r1 := <-results
	assert.Equal(t, "http://gocardless.com/a", r1.url)
	assert.Len(t, r1.assets, 1)
	if len(r1.assets) == 1 {
		assert.Equal(t, "picA.png", r1.assets[0])
	}
	r2, _ := <-results
	assert.Equal(t, "https://gocardless.com/b", r2.url)
	assert.Len(t, r2.assets, 1)
	if len(r2.assets) == 1 {
		assert.Equal(t, "picB.png", r2.assets[0])
	}
	_, more := <-results
	assert.False(t, more)
}

func TestDontGoOutsideOriginalDomain(t *testing.T) {
	var results chan SiteInfo
	results = crawl("http://gocardless.com/a", fakeFetcher2)

	r1 := <-results
	assert.Equal(t, "http://gocardless.com/a", r1.url)
	assert.Len(t, r1.assets, 1)
	if len(r1.assets) == 1 {
		assert.Equal(t, "picA.png", r1.assets[0])
	}
	_, more := <-results
	assert.False(t, more)
}

func TestRelativeUrls(t *testing.T) {
	var results chan SiteInfo
	results = crawl("http://gocardless.com/a", fakeFetcherWithRelativeUrl)

	r1 := <-results
	assert.Equal(t, "http://gocardless.com/a", r1.url)
	r2, _ := <-results
	assert.Equal(t, "http://gocardless.com/b", r2.url)
	_, more := <-results
	assert.False(t, more)
}

var fakeFetcher1 = FakeFetcher{
	"http://gocardless.com/a":  `<a href="https://gocardless.com/b"><img src="picA.png"/></a>`,
	"https://gocardless.com/b": `<a href="http://gocardless.com/a"></a><img src="picB.png"/><a href="https://gocardless.com/b">link text</a>`,
}

var fakeFetcher2 = FakeFetcher{
	"http://gocardless.com/a": `<a href="https://google.com"><img src="picA.png"/></a>`,
	"https://google.com":      `welcome to google <img src="google.png"/>`,
}

var fakeFetcherWithRelativeUrl = FakeFetcher{
	"http://gocardless.com/a": `<a href="/b"><img src="picA.png"/></a>`,
	"http://gocardless.com/b": `<a href="http://gocardless.com/a"></a><img src="picB.png"/><a href="http://gocardless.com/b">link text</a>`,
}
