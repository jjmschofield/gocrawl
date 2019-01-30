package pages

import (
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"github.com/jjmschofield/GoCrawl/internal/pkg/md5"
	"io"
	"log"
	"net/url"
	"strings"
)

type Page struct {
	Id       string
	URL      url.URL
	OutInternalPages []string
	OutLinks []string
	Err      error
}

func PageFromUrl(srcUrl url.URL) Page {
	id, normalizedUrl := CalcPageId(srcUrl)

	return Page{
		Id:  id,
		URL: normalizedUrl,
	}
}

func (page *Page) FetchBody(fetcher PageBodyFetcher) (bodyReader io.ReadCloser, err error) {
	return fetcher(page.URL)
}

func (page *Page) FetchHrefs(fetcher PageBodyFetcher, reader HrefReader) (hrefs []string, err error) {
	bodyReader, err := page.FetchBody(fetcher)

	if err != nil {
		log.Printf("couldn't get html page for %s %v", page.Id, err)
		return nil, err
	}

	return reader(bodyReader), nil
}

func (page *Page) AppendOutInternalPage(outPage Page) []string {
	page.OutInternalPages = append(page.OutInternalPages, outPage.Id)
	return page.OutInternalPages
}

func (page *Page) AppendOutLink(link links.Link) []string {
	page.OutLinks = append(page.OutLinks, link.ToURL.String())
	return page.OutLinks
}

func CalcPageId(srcUrl url.URL) (id string, normalizedUrl url.URL) {
	normalizedUrl = normalizePageUrl(srcUrl)
	id = md5.HashString(normalizedUrl.String())
	return id, normalizedUrl
}

func (page *Page) Print() {
	log.Printf("Page at %s with id %s", page.URL.String(), page.Id)
	if page.Err != nil {
		log.Printf("Page had error %s", page.Err)
	}

	log.Printf("OutInternalPages %v", len(page.OutInternalPages))
	log.Printf("%s", page.OutInternalPages)

	log.Printf("Outlinks %v", len(page.OutLinks))
	log.Printf("%s", page.OutLinks)
}

func normalizePageUrl(srcUrl url.URL) url.URL { //TODO - will this mutate?
	srcUrl.Fragment = ""
	srcUrl.RawQuery = ""
	srcUrl.Path = strings.TrimRight(srcUrl.Path, "/")
	srcUrl.RawPath = strings.TrimRight(srcUrl.RawPath, "/")
	return srcUrl
}
