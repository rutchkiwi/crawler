package main

import (
	"fmt"
	"strings"
)

type SiteInfo struct {
	url    string
	assets []string
}

func (info SiteInfo) String() string {
	return fmt.Sprintf("%s has assets: %s", info.url, strings.Join(info.assets, ", "))
}
