package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/fetch"
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"io"
	"net/url"
)

type ScrapeResult struct {
	OutPages pages.PageGroup
	OutLinks map[string]links.Link
}

func Scrape(target url.URL) (result ScrapeResult, err error) {
	bodyReader, err := fetch.Body(target)

	if err != nil {
		return ScrapeResult{}, err
	}

	result.OutLinks = extractLinks(target, bodyReader)

	result.OutPages = createOutPages(result.OutLinks)

	return result, nil
}

func extractLinks(target url.URL, bodyReader io.ReadCloser) (extracted map[string]links.Link) {
	hrefs := fetch.ReadHrefs(bodyReader)
	return links.FromHrefs(target, hrefs)
}

func createOutPages(outLinks map[string]links.Link) (group pages.PageGroup) {
	group = pages.PageGroup{
		Internal: make(map[string]string),
	}

	for _, link := range outLinks {
		if link.Type == links.InternalPageType{
			pageUrl, _ :=  url.Parse(link.ToURL)
			id, normalizedUrl := pages.CalcPageId(*pageUrl)
			group.Internal[id] = normalizedUrl.String()
		}
	}

	return group
}
