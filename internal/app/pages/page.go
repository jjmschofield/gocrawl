package pages

import (
	"encoding/json"
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"log"
	"net/url"
)

type Page struct {
	Id       string
	URL      url.URL
	OutPages PageGroup
	OutLinks links.LinkGroup
	Err      error
}

func PageFromUrl(srcUrl url.URL) Page {
	id, normalizedUrl := CalcPageId(srcUrl)

	return Page{
		Id:  id,
		URL: normalizedUrl,
	}
}

func (page Page) MarshalJSON() ([]byte, error) {
	basicPage := struct {
		Id       string          `json:"id"`
		URL      string          `json:"url"`
		OutPages PageGroup		 `json:"outPages"`
		OutLinks links.LinkGroup `json:"outLinks"`
		Err      error           `json:"error"`
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
	log.Printf("Internal Links: %v", len(page.OutLinks.Internal))
	printLinkSlice(page.OutLinks.Internal)
	log.Printf("External Links: %v", len(page.OutLinks.External))
	printLinkSlice(page.OutLinks.External)
	log.Printf("Tel Links: %v", len(page.OutLinks.Tel))
	printLinkSlice(page.OutLinks.Tel)
	log.Printf("Mailto Links: %v", len(page.OutLinks.Mailto))
	printLinkSlice(page.OutLinks.Mailto)
}

func printLinkSlice(slice []links.Link) {
	log.Printf("| Link Id | Link Type | Link From -> Link To | Link To Page Id |")
	for _, link := range slice {
		targetPageId, normalizedUrl := CalcPageId(link.ToURL)
		log.Printf("| %s | %s | %s -> %s | %s(%s) |", link.Id, link.Type, link.FromURL.String(), link.ToURL.String(), targetPageId, normalizedUrl.String())
	}
}
