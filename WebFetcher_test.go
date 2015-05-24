package main

import "testing"

func TestUrlParsing(t *testing.T) {
	body := `Go offers built-in support for <a href="http://en.wikipedia.org/wiki/Regular_expression">regular expressions</a>.
Here are some examples of  common regexp-related tasks`
	if "http://en.wikipedia.org/wiki/Regular_expression" != findUrls(body)[0] {
		t.Error("url mismatch!")
	}
}
