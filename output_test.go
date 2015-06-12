package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {

	results := make(chan SiteInfo, 10)
	results <- SiteInfo{"pageA", []string{"asset1", "asset2"}}
	close(results)

	buf := new(bytes.Buffer)
	print_results(results, buf)
	lines := readBufferIntoLines(buf)
	assert.Equal(t, lines[0], "pageA has assets: asset1, asset2")
}

func readBufferIntoLines(b *bytes.Buffer) []string {
	var lines []string
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
