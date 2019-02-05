package links

import (
	"github.com/jjmschofield/GoCrawl/internal/app/md5"
	"log"
	"net/url"
)

type Link struct {
	Id      string `json:"id"`
	ToURL   string `json:"toUrl"`
	FromURL string `json:"fromUrl"`
	Type    string `json:"type"`
}

func NewAbsLink(fromUrl url.URL, toUrl url.URL) Link {
	var absToUrl url.URL

	if !toUrl.IsAbs() {
		absToUrl = *fromUrl.ResolveReference(&toUrl)
	} else {
		absToUrl = toUrl
	}

	return Link{
		Id:      md5.HashString(fromUrl.String() + toUrl.String()),
		ToURL:   absToUrl.String(),
		FromURL: fromUrl.String(),
		Type:    calcType(fromUrl, absToUrl),
	}
}

func FromHref(pageUrl url.URL, href string) (link Link, err error) {
	toUrl, err := url.Parse(href)

	if err != nil {
		log.Printf("invalid url found %s on %s", href, pageUrl.String())
		return Link{}, err
	}

	return NewAbsLink(pageUrl, *toUrl), nil
}

func FromHrefs(srcUrl url.URL, hrefs []string) (links map[string]Link) {
	links = make(map[string]Link)

	for _, href := range hrefs {
		link, err := FromHref(srcUrl, href)

		if err != nil {
			log.Printf("invalid url found %s on %s", href, srcUrl.String())
		} else {
			links[link.Id] = link
		}
	}

	return links
}
