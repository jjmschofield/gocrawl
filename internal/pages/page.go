package pages

import (
	"encoding/json"
	"github.com/jjmschofield/gocrawl/internal/links"
	"log"
	"net/url"
)

type Page struct {
	Id       string
	URL      url.URL
	OutPages PageGroup
	OutLinks map[string]links.Link
	Err      error
}

type PageGroup struct {
	Internal map[string]string `json:"internal"`
}

func PageFromUrl(srcUrl url.URL) Page {
	id, normalizedUrl := CalcPageId(srcUrl)

	return Page{
		Id:  id,
		URL: normalizedUrl,
		OutPages: PageGroup{
			Internal: make(map[string]string),
		},
		OutLinks: make(map[string]links.Link),
	}
}

func (page Page) MarshalJSON() ([]byte, error) {
	basicPage := struct {
		Id       string                `json:"id"`
		URL      string                `json:"url"`
		OutPages PageGroup             `json:"outPages"`
		OutLinks map[string]links.Link `json:"outLinks"`
		Err      error                 `json:"error"`
	}{
		Id:       page.Id,
		URL:      page.URL.String(),
		OutPages: page.OutPages,
		OutLinks: page.OutLinks,
		Err:      page.Err,
	}

	return json.Marshal(basicPage)
}

func (page *Page) Print() {
	log.Printf("Page at %s with id %s", page.URL.String(), page.Id)
	if page.Err != nil {
		log.Printf("Page had error %s", page.Err)
	}
}
