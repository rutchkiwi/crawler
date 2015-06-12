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

func TestFindSameSiteUrlsInString(t *testing.T) {
	str := `<div class="nav__item u-relative" data-reactid=".2ftrc1lhrsw.0.0.1.0.2"><span data-reactid=".2ftrc1lhrsw.0.0.1.0.2.0"><a id="track-nav-stories" class="u-padding-Vl u-block u-link-invert" href="/stories/" data-reactid=".2ftrc1lhrsw.0.0.1.0.2.0.0"><div class="nav__item-link" data-reactid=".2ftrc1lhrsw.0.0.1.0.2.0.0.0"><span data-reactid=".2ftrc1lhrsw.0.0.1.0.2.0.0.0.0">Stories</span></div></a></span></div>`
	res := findUrls(str)
	assert.Len(t, res, 1)
}
