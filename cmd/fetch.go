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

	log.Printf( "Page at %s with id %s", page.URL.String(), page.Id)
	log.Printf( "Internal Links: %v", len(page.OutLinks.Internal))
	printLinkSlice(page.OutLinks.Internal)
	log.Printf( "External Links: %v", len(page.OutLinks.External))
	printLinkSlice(page.OutLinks.External)
	log.Printf( "Tel Links: %v", len(page.OutLinks.Tel))
	printLinkSlice(page.OutLinks.Tel)
	log.Printf( "Mailto Links: %v", len(page.OutLinks.Mailto))
	printLinkSlice(page.OutLinks.Mailto)

}

func printLinkSlice(slice []links.Link){
	for _, link := range slice {
		log.Printf("- %s %s %s -> %s", link.Id, link.Type, link.FromURL.String(), link.ToURL.String())
	}
}