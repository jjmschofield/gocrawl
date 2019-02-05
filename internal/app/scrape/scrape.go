package scrape

import (
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"net/url"
)

type Result struct {
	OutPages pages.PageGroup
	OutLinks map[string]links.Link
}

func Scrape(target url.URL) (result Result, err error) {
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
		if link.Type == links.InternalPageType{
			pageUrl, _ :=  url.Parse(link.ToURL)
			id, normalizedUrl := pages.CalcPageId(*pageUrl)
			group.Internal[id] = normalizedUrl.String()
		}
	}

	return group
}
