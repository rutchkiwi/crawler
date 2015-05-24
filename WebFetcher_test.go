package main

import "testing"

func TestUrlParsing(t *testing.T) {
	body := `Go offers built-in support for <a href="http://en.wikipedia.org/wiki/Regular_expression">regular expressions</a>.
Here are some examples of  common regexp-related tasks`
	if "http://en.wikipedia.org/wiki/Regular_expression" != findUrls(body)[0] {
		t.Error("url mismatch!")
	}
}

func TestUrlParsingMultipleUrls(t *testing.T) {
	body := `Go offers <a href="http://linkone.com">regular expressions</a>.
Here are some examples of  common regexp-related http://shouldntfindthis tasks
	<a href="http://www.linktwo.se"/>`
	if "http://linkone.com" != findUrls(body)[0] {
		t.Error("url mismatch!")
	}
	if "http://www.linktwo.se" != findUrls(body)[1] {
		t.Error("url mismatch!")
	}
}
