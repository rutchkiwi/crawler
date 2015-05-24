package main

import "fmt"

// FakeFetcher is Fetcher that returns canned results.
type FakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f FakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}
