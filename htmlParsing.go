package main

import (
	"strings"

	// there is probably a nicer libray to do this kind of stuff
	"golang.org/x/net/html"
)

func findAssets(bodyHtml string) []string {
	links := make([]string, 0)
	page := html.NewTokenizer(strings.NewReader(bodyHtml))
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := page.Token()
		for _, attr := range token.Attr {
			if attr.Key == "src" {
				links = append(links, attr.Val)
				break
			}
		}
	}
	return links
}

func findUrls(bodyHtml string) []string {
	links := make([]string, 0)
	page := html.NewTokenizer(strings.NewReader(bodyHtml))
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := page.Token()
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
	}
	return links
}
