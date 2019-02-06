package main

import (
	"flag"
	"fmt"
	"github.com/jjmschofield/GoCrawl/internal/crawl"
	"log"
	"net/url"
	"time"
)

func main() {
	crawlUrlRaw := flag.String("url", "https://monzo.com", "an absolute url, including protocol and hostname")
	workerCount := flag.Int("workers", 100, "Number of crawl workers to run")
	outFilePath := flag.String("dir", "data", "A file path to send results to")

	flag.Parse()

	crawlUrl, err := url.Parse(*crawlUrlRaw)

	if err != nil {
		log.Panic(err)
	}
	start := time.Now()

	counters := Crawl(*crawlUrl, *workerCount, *outFilePath)

	end := time.Now()

	fmt.Printf("Scrape Completed in %v ms \n", (end.UnixNano()-start.UnixNano())/int64(time.Millisecond))
	fmt.Printf(" Discovered: %v, \n Crawled: %v \n Parallel Crawls Peak: %v \n Scrape Queue Peak: %v \n Processing Peak: %v \n", counters.Discovered.Count(), counters.CrawlComplete.Count(), counters.Crawling.Peak(), counters.CrawlsQueued.Peak(), counters.Processing.Peak())
}

func Crawl(crawlUrl url.URL, workerCount int, outFilePath string) crawl.Counters {
	crawler := crawl.NewDefaultPageCrawler(workerCount, outFilePath)
	counters := crawler.Crawl(crawlUrl)
	return counters
}
