package main

import (
	"flag"
	"github.com/jjmschofield/GoCrawl/internal/app/crawl"
	"log"
	"net/url"
	"time"
)

func main(){
	start := time.Now()

	crawlUrlRaw := flag.String("url", "https://monzo.com", "an absolute url eg http://www.google.co.uk")

	flag.Parse()

	crawlUrl, err := url.Parse(*crawlUrlRaw)

	if err != nil{
		log.Panic(err)
	}

	pages := crawl.FromUrl(*crawlUrl, crawl.PageCrawler, 200)

	end:= time.Now()

	log.Printf("Complete and found %v pages in %v ms" , len(pages), (end.UnixNano() - start.UnixNano()) / int64(time.Millisecond))
}
