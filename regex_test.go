package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindingImageAssets(t *testing.T) {
	assets := findAssets(`<img src="picA.png"/>`)
	assert.Len(t, assets, 1)
}
