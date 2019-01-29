package main

import (
	"flag"
	"github.com/jjmschofield/GoCrawl/internal/app/crawl"
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

	pages := crawl.FromUrl(*crawlUrl, crawl.PageCrawler, 1)

	for _, page := range pages {
		page.Print()
	}
}
