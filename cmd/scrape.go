package main

import (
	"flag"
	"github.com/jjmschofield/GoCrawl/internal/app/scrape"
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
	result, err := scrape.Scrape(page.URL)

	if err != nil {
		log.Panicf("failed to scrape %s %s", crawlUrl, err)
	}

	log.Printf("%+v", result)
}

