package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var results chan SiteInfo
	results = crawl("https://gocardless.com/", WebFetcher{})

	printResultsDebug(results, os.Stdout)
}

// TODO change to camelcase
func print_results(results chan SiteInfo, out io.Writer) {
	for res := range results {
		fmt.Fprintln(out, res)
	}
}

func printResultsDebug(results chan SiteInfo, out io.Writer) {
	for res := range results {
		fmt.Fprintln(out, res.url)
	}
}
