package main

import "fmt"

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, err error)
}

// FakeFetcher is Fetcher that returns canned results.
type FakeFetcher map[string]string

func (f FakeFetcher) Fetch(url string) (string, error) {
	if body, ok := f[url]; ok {
		return body, nil
	}
	return "", fmt.Errorf("not found: %s", url)
}
