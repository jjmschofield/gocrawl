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
	OutLinks PageOutLinks
	Err      error
}

type PageOutLinks struct {
	Internal []links.Link
	External []links.Link
	Tel      []links.Link
	Mailto   []links.Link
	Unknown  []links.Link
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

func (page *Page) AppendOutLink(link links.Link) []links.Link {
	switch {
	case link.Type == links.InternalPageType:
		page.OutLinks.Internal = append(page.OutLinks.Internal, link)
		return page.OutLinks.Internal
	case link.Type == links.ExternalPageType:
		page.OutLinks.External = append(page.OutLinks.External, link)
		return page.OutLinks.Internal
	case link.Type == links.MailtoType:
		page.OutLinks.Mailto = append(page.OutLinks.Mailto, link)
		return page.OutLinks.Internal
	case link.Type == links.TelType:
		page.OutLinks.Tel = append(page.OutLinks.Tel, link)
		return page.OutLinks.Internal
	default:
		page.OutLinks.Unknown = append(page.OutLinks.Unknown, link)
		return page.OutLinks.Internal
	}
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
	log.Printf("Internal Links: %v", len(page.OutLinks.Internal))
	printLinkSlice(page.OutLinks.Internal)
	log.Printf("External Links: %v", len(page.OutLinks.External))
	printLinkSlice(page.OutLinks.External)
	log.Printf("Tel Links: %v", len(page.OutLinks.Tel))
	printLinkSlice(page.OutLinks.Tel)
	log.Printf("Mailto Links: %v", len(page.OutLinks.Mailto))
	printLinkSlice(page.OutLinks.Mailto)
}

func normalizePageUrl(srcUrl url.URL) url.URL { //TODO - will this mutate?
	srcUrl.Fragment = ""
	srcUrl.RawQuery = ""
	srcUrl.Path = strings.TrimRight(srcUrl.Path, "/")
	srcUrl.RawPath = strings.TrimRight(srcUrl.RawPath, "/")
	return srcUrl
}

func printLinkSlice(slice []links.Link) {
	log.Printf("| Link Id | Link Type | Link From -> Link To | Link To Page Id |")
	for _, link := range slice {
		targetPageId, normalizedUrl := CalcPageId(link.ToURL)
		log.Printf("| %s | %s | %s -> %s | %s(%s) |", link.Id, link.Type, link.FromURL.String(), link.ToURL.String(), targetPageId, normalizedUrl.String())
	}
}
