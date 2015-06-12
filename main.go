package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var results chan SiteInfo
	results = crawl("https://gocardless.com/", WebFetcher{})

	print_results(results, os.Stdout)
}

func print_results(results chan SiteInfo, out io.Writer) {
	for res := range results {
		fmt.Fprintln(out, res)
	}
}
