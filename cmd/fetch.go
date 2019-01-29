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

	if err != nil{
		log.Panic(err)
	}

	discoveredLinks := links.FromHrefs(page.URL, hrefs)

	for _, link := range discoveredLinks {
		page.AppendOutLink(link)
	}

	page.Print()
}

