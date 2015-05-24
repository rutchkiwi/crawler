package main

import "fmt"

type WebFetcher struct{}

func (f WebFetcher) Fetch(url string) (string, []string, error) {
	return "", nil, fmt.Errorf("not implemented")
}
