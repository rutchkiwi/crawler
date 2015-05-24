package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type WebFetcher struct{}

func (f WebFetcher) Fetch(url string) (string, []string, error) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyStr := string(bodyBytes[:])

	fmt.Println(bodyStr)

	return "", nil, fmt.Errorf("not implemented")
}

func findUrls(html string) []string {
	r, _ := regexp.Compile(`"(http://.*?)">`)
	matches := r.FindAllStringSubmatch(html, -1)

	res := make([]string, len(matches))
	for i, match := range matches {
		res[i] = match[1]
	}
	fmt.Println(res)

	return res
}

func main() {
	body := `Go offers built-in support for <a href="http://en.wikipedia.org/wiki/Regular_expression">regular expressions</a>.
Here are some examples of  common regexp-related tasks`
	fmt.Printf("found: %s\n", findUrls(body))
}
