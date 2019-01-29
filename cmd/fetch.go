package main

import (
	"flag"
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"log"
	"net/url"
)

func main(){
	crawlUrlRaw := flag.String("url", "https://monzo.com", "an absolute url eg http://www.google.co.uk")

	flag.Parse()

	crawlUrl, err := url.Parse(*crawlUrlRaw)

	if err != nil{
		log.Panic(err)
	}

	page := pages.PageFromUrl(*crawlUrl)
	hrefs, err := page.FetchHrefs(pages.FetchPageBody, pages.ReadHrefs)
	discoveredLinks := links.FromHrefs(page.URL, hrefs)

	if err != nil{
		log.Panic(err)
	}

	for _, link := range discoveredLinks {
		page.OutLinks = append(page.OutLinks, link.Id)
	}

	for _, linkId := range page.OutLinks {
		toUrl :=  discoveredLinks[linkId].ToURL
		log.Printf("Page id %s with url %s link id %s type %s to %s", page.Id, page.URL.String(), linkId, discoveredLinks[linkId].Type, toUrl.String())
	}
}