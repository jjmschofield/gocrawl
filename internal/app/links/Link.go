package links

import (
	"github.com/jjmschofield/GoCrawl/internal/pkg/md5"
	"log"
	"net/url"
)

type Link struct {
	Id      string
	ToURL   url.URL
	FromURL url.URL
	Type    LinkType
}

func NewAbsLink(fromUrl url.URL, toUrl url.URL) Link { // TODO - needs to be normalized to exclude query params
	if !toUrl.IsAbs() {
		toUrl.Scheme = fromUrl.Scheme
		toUrl.Host = fromUrl.Host
	}

	return Link{
		Id:      md5.HashString(fromUrl.String() + toUrl.String()),
		ToURL:   toUrl,
		FromURL: fromUrl,
		Type: calcType(fromUrl, toUrl),
	}
}

func FromHref(pageUrl url.URL, href string) (link Link, err error) {
	toUrl, err := url.Parse(href)

	if err != nil {
		log.Printf("invalid url found %s on %s", href, pageUrl.String())
		return Link{}, err
	}

	return NewAbsLink(pageUrl, *toUrl), nil // TODO - will this mutate or is it de-referencing
}

func FromHrefs(pageUrl url.URL, hrefs []string) (links map[string]Link) {
	links = make(map[string]Link)

	for _, href := range hrefs {
		link, err := FromHref(pageUrl, href)

		if err != nil {
			log.Printf("invalid url found %s on %s", href, pageUrl.String())
		}

		links[link.Id] = link
	}

	return links
}
