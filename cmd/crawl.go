package main

import (
	"flag"
	"fmt"
	"github.com/jjmschofield/GoCrawl/internal/app/crawl"
	"log"
	"net/url"
	"time"
)

func main() {
	crawlUrlRaw := flag.String("url", "https://monzo.com", "an absolute url eg http://www.google.co.uk")
	workerCount := flag.Int("workers", 100, "Number of crawl workers to run")
	outFilePath := flag.String("out", "data", "A file path to send results to")

	flag.Parse()

	crawlUrl, err := url.Parse(*crawlUrlRaw)

	if err != nil {
		log.Panic(err)
	}

	start := time.Now()

	crawler := crawl.NewDefaultCrawler(*workerCount, *outFilePath)

	counters := crawler.Crawl(*crawlUrl)

	end := time.Now()

	fmt.Printf("Scrape Completed in %v ms \n", (end.UnixNano()-start.UnixNano())/int64(time.Millisecond))
	fmt.Printf(" Discovered: %v, \n Crawled: %v \n Parallel Crawls Peak: %v \n Scrape Queue Peak: %v \n Processing Peak: %v \n", counters.Discovered.Count(), counters.CrawlComplete.Count(), counters.Crawling.Peak(), counters.CrawlsQueued.Peak(), counters.Processing.Peak())
}
