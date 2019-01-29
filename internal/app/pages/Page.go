package pages

import (
	"github.com/jjmschofield/GoCrawl/internal/pkg/md5"
	"io"
	"log"
	"net/url"
	"strings"
)

type Page struct {
	Id       string
	URL      url.URL
	OutLinks []string
	InLinks  []string
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

func NormalizePageUrl(srcUrl url.URL) url.URL { //TODO - will this mutate?
	srcUrl.Fragment = ""
	srcUrl.RawQuery = ""
	strings.TrimRight(srcUrl.RawPath, "/")
	return srcUrl
}

func CalcPageId(srcUrl url.URL) (id string, normalizedUrl url.URL) {
	normalizedUrl = NormalizePageUrl(srcUrl)
	id = md5.HashString(normalizedUrl.String())
	return id, normalizedUrl
}
