package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindingImageAssetsInString(t *testing.T) {
	assets := findAssets(`<img src="picA.png"/>`)
	assert.Len(t, assets, 1)
}

func TestFindingImageAssetsInStringTricky(t *testing.T) {
	assets := findAssets(`<img alt="a picture" src = "picA.png"/>`)
	assert.Len(t, assets, 1)
}

func TestFindingAllAssetsInString(t *testing.T) {
	assets := findAssets(`<div class="wm-site-info">
<a href="//wikimediafoundation.org/"><img src="/static/images/wikimedia-button.png" srcset="https://upload.wikimedia.org/wikipedia/meta/3/3a/A_Wikimedia_project_1.5x.png 1.5x, https://upload.wikimedia.org/wikipedia/meta/b/b7/A_Wikimedia_project_2x.png 2x" width="88" height="31" alt="A Wikimedia Project"></a>
</div>

<div style="text-align:center"><a href="//wikimediafoundation.org/wiki/Terms_of_Use">Terms of Use</a> | <a href="//wikimediafoundation.org/wiki/Privacy_policy">Privacy Policy</a></div>
<script src="//meta.wikimedia.org/w/load.php?debug=false&amp;lang=en&amp;modules=ext.gadget.wm-portal&amp;only=scripts&amp;skin=vector&amp;*"></script>`)
	assert.Len(t, assets, 2)
}

func TestFindUrlsInString(t *testing.T) {
	str := `<a href="https://test.com/b"><img src="picA.png"/></a><a href="https://test.com/c"><img src="picA.png"/></a>`
	res := findUrls(str)
	assert.Len(t, res, 2)
	assert.Equal(t, "https://test.com/b", res[0])
	assert.Equal(t, "https://test.com/c", res[1])
}

func TestFindUrlsTricky(t *testing.T) {
	str := `<a     href = "https://test.com/b"><img src="picA.png"/></a><a target="_blank" href="/c"><img src="picA.png"/></a>`
	res := findUrls(str)
	assert.Len(t, res, 2)
}

func TestFindSameSiteUrlsInString(t *testing.T) {
	str := `<div class="nav__item u-relative" data-reactid=".2ftrc1lhrsw.0.0.1.0.2"><span data-reactid=".2ftrc1lhrsw.0.0.1.0.2.0"><a id="track-nav-stories" class="u-padding-Vl u-block u-link-invert" href="/stories/" data-reactid=".2ftrc1lhrsw.0.0.1.0.2.0.0"><div class="nav__item-link" data-reactid=".2ftrc1lhrsw.0.0.1.0.2.0.0.0"><span data-reactid=".2ftrc1lhrsw.0.0.1.0.2.0.0.0.0">Stories</span></div></a></span></div>`
	res := findUrls(str)
	assert.Len(t, res, 1)
}
