package links

import (
	"encoding/json"
	"github.com/jjmschofield/GoCrawl/internal/app/md5"
	"log"
	"net/url"
)

type Link struct {
	Id      string  `json:"id"`
	ToURL   url.URL `json:"url,string"`
	FromURL url.URL `json:"url,string"`
	Type    string  `json:"type"`
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
		ToURL:   absToUrl,
		FromURL: fromUrl,
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

func FromHrefs(srcUrl url.URL, hrefs []string) (links []Link) {
	for _, href := range hrefs {
		link, err := FromHref(srcUrl, href)

		if err != nil {
			log.Printf("invalid url found %s on %s", href, srcUrl.String())
		} else{
			links = append(links, link)
		}
	}

	return links
}

func (link Link) MarshalJSON() ([]byte, error) {
	basicLink := struct {
		Id      string `json:"id"`
		ToURL   string `json:"toUrl"`
		Type    string `json:"type"`
		FromUrl string `json:"fromUrl"`
	}{
		Id:      link.Id,
		ToURL:   link.ToURL.String(),
		Type:    link.Type,
		FromUrl: link.FromURL.String(),
	}

	return json.Marshal(basicLink)
}
