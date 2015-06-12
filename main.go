package main

import (
	"fmt"
	"io"
)

func main() {
	var results chan SiteInfo
	results = crawl("http://a.com", fakeFetcher)

	for e := range results {
		fmt.Println(e)
	}
}

func print_results(results chan SiteInfo, out io.Writer) {
	for res := range results {
		fmt.Fprintln(out, res)
	}
}
