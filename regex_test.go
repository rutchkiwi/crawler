package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindingImageAssetsInString(t *testing.T) {
	assets := findAssets(`<img src="picA.png"/>`)
	assert.Len(t, assets, 1)
}

func TestFindUrlsInString(t *testing.T) {
	str := `<a href="https://gocardless.com/b"><img src="picA.png"/></a><a href="https://gocardless.com/c"><img src="picA.png"/></a>`
	res := findUrls(str)
	assert.Len(t, res, 2)
	//assert.Equal(t, "https://gocardless.com/bd", res[0])
}
