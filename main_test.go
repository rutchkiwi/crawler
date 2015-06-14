package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fakeFetcher = FakeFetcher{
	"http://gocardless.com/a":  `<a href="https://gocardless.com/b"><img src="picA.png"/><img src="picB.png"/></a>`,
	"https://gocardless.com/b": `<a href="http://gocardless.com/a"></a><img src="picC.png"/><a href="https://gocardless.com/b">link text</a>`,
}

type FakeFetcher map[string]string

func (f FakeFetcher) Fetch(url string) (string, error) {
	if body, ok := f[url]; ok {
		return body, nil
	}
	return "", fmt.Errorf("not found: %s", url)
}

func TestCrawlsRightUrls(t *testing.T) {
	urls := getUrlsFromCrawling(fakeFetcher, "http://gocardless.com/a")

	assert.Contains(t, urls, "http://gocardless.com/a")
	assert.Contains(t, urls, "https://gocardless.com/b")
	assert.Len(t, urls, 2)
}

func TestFindImageAssets(t *testing.T) {
	assets := getAssetsFromCrawling(fakeFetcher, "http://gocardless.com/a")

	assert.Len(t, assets, 2)
	assert.Contains(t, assets, []string{"picA.png", "picB.png"})
	assert.Contains(t, assets, []string{"picC.png"})
}

func TestError(t *testing.T) {
	_, errors := crawl("http://nothingHere.se", fakeFetcher)
	err := <-errors
	assert.EqualError(t, err, "not found: http://nothingHere.se")
}

func TestDontGoOutsideOriginalDomain(t *testing.T) {
	var fakeFetcherWithLinkOutsideDomain = FakeFetcher{
		"http://gocardless.com/a": `<a href="https://google.com"><img src="picA.png"/></a>`,
		"https://google.com":      `welcome to google <img src="google.png"/>`,
	}
	urls := getUrlsFromCrawling(fakeFetcherWithLinkOutsideDomain, "http://gocardless.com/a")

	assert.Len(t, urls, 1)
	assert.Contains(t, urls, "http://gocardless.com/a")
}

func TestRelativeUrls(t *testing.T) {
	var fakeFetcherWithRelativeUrl = FakeFetcher{
		"http://gocardless.com/a": `<a href="/b"><img src="picA.png"/></a>`,
		"http://gocardless.com/b": `<a href="http://gocardless.com/a"></a><img src="picB.png"/><a href="http://gocardless.com/b">link text</a>`,
	}
	urls := getUrlsFromCrawling(fakeFetcherWithRelativeUrl, "http://gocardless.com/a")

	assert.Len(t, urls, 2)
	assert.Contains(t, urls, "http://gocardless.com/a")
	assert.Contains(t, urls, "http://gocardless.com/b")
}

func TestManyLinks(t *testing.T) {
	var manyLinksFetcher = FakeFetcher{
		"http://test.com/a": `<a href="/b"></a> <a href="/c"></a> <a href="/d"></a>`,
		"http://test.com/b": `<a href="/1b"></a> <a href="/1c"></a> <a href="/1d"></a>`,
		"http://test.com/c": `<a href="/b"></a> <a href="/c"></a> <a href="/d"></a>`,
		"http://test.com/d": `<a href="/b"></a> <a href="/c"></a> <a href="/d"></a>`,
	}
	urls := getUrlsFromCrawling(manyLinksFetcher, "http://test.com/a")

	assert.Len(t, urls, 4)
	assert.Contains(t, urls, "http://test.com/a")
	assert.Contains(t, urls, "http://test.com/b")
	assert.Contains(t, urls, "http://test.com/c")
	assert.Contains(t, urls, "http://test.com/d")

}
func getUrlsFromCrawling(fetcher Fetcher, url string) (urls []string) {
	urls = make([]string, 0)

	var results chan SiteInfo
	results, _ = crawl(url, fetcher)

	for r := range results {
		urls = append(urls, r.url)
	}
	return urls
}

func getAssetsFromCrawling(fetcher Fetcher, url string) (assets [][]string) {
	assets = make([][]string, 0)

	var results chan SiteInfo
	results, _ = crawl(url, fetcher)

	for r := range results {
		assets = append(assets, r.assets)
	}
	return assets
}
