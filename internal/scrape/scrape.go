package scrape

import (
	"github.com/jjmschofield/gocrawl/internal/links"
	"github.com/jjmschofield/gocrawl/internal/pages"
	"net/url"
)

type Result struct {
	OutPages pages.PageGroup
	OutLinks map[string]links.Link
}

//go:generate counterfeiter . Scraper
type Scraper interface {
	Scrape(target url.URL) (result Result, err error)
}

type PageScraper struct{}

func (PageScraper) Scrape(target url.URL) (result Result, err error) {
	outLinks, err := fetchLinks(target)

	if err != nil {
		return Result{}, err
	}

	result = Result{
		OutLinks: outLinks,
		OutPages: createOutPages(outLinks),
	}

	return result, nil
}

func createOutPages(outLinks map[string]links.Link) (group pages.PageGroup) {
	group = pages.PageGroup{
		Internal: make(map[string]string),
	}

	for _, link := range outLinks {
		if link.Type == links.InternalPageType {
			pageUrl, _ := url.Parse(link.ToURL)
			id, normalizedUrl := pages.CalcPageId(*pageUrl)
			group.Internal[id] = normalizedUrl.String()
		}
	}

	return group
}
