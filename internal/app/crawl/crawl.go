package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/fetch"
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"io"
	"net/url"
)

type PageCrawlResult struct {
	OutPages pages.PageGroup
	OutLinks []links.Link
}

func Crawl(target url.URL) (result PageCrawlResult, err error) {
	bodyReader, err := fetch.Body(target)

	if err != nil {
		return PageCrawlResult{}, err
	}

	result.OutLinks = extractLinks(target, bodyReader)

	result.OutPages = createPages(result.OutLinks)

	return result, nil
}

func extractLinks(target url.URL, bodyReader io.ReadCloser) (extracted []links.Link) {
	hrefs := fetch.ReadHrefs(bodyReader)

	return links.FromHrefs(target, hrefs)
}

func createPages(outLinks []links.Link) (group pages.PageGroup) {
	var internal []pages.Page

	for _, link := range outLinks {
		if link.Type == links.InternalPageType {
			internal = append(internal, pages.PageFromUrl(link.ToURL))
		}
	}

	return pages.ToPageGroup(internal)
}
