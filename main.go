package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	tic := time.Now()
	var results chan SiteInfo
	results, errors := crawl("https://gocardless.com/", newWebFetcher())

	noSuccesses := printResultsDebug(results, os.Stdout)
	noErrors := 0
	for range errors {
		noErrors++
	}
	toc := time.Now()
	timeSpent := toc.Sub(tic)
	fmt.Printf("Crawled %d urls, got %d errors in %v\n", noSuccesses, noErrors, timeSpent)
}

// TODO change to camelcase
func print_results(results chan SiteInfo, out io.Writer) (noSuccesses int) {
	noSuccesses = 0
	for res := range results {
		fmt.Fprintln(out, res)
		noSuccesses++
	}
	return noSuccesses
}

func printResultsDebug(results chan SiteInfo, out io.Writer) (noSuccesses int) {
	noSuccesses = 0
	for range results {
		fmt.Fprintf(out, ".")
		noSuccesses++
	}
	return noSuccesses
}
