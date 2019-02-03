package pages

import (
	"github.com/jjmschofield/GoCrawl/internal/app/md5"
	"net/url"
	"strings"
)

func CalcPageId(srcUrl url.URL) (id string, normalizedUrl url.URL) {
	normalizedUrl = normalizePageUrl(srcUrl)
	id = md5.HashString(normalizedUrl.String())
	return id, normalizedUrl
}

func normalizePageUrl(srcUrl url.URL) url.URL {
	srcUrl.Fragment = ""
	srcUrl.Path = strings.TrimRight(srcUrl.Path, "/")
	srcUrl.RawPath = strings.TrimRight(srcUrl.RawPath, "/")
	return srcUrl
}