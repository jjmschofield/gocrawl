package main

import (
	"flag"
	"fmt"
	"github.com/jjmschofield/GoCrawl/internal/app/caches"
	"github.com/jjmschofield/GoCrawl/internal/app/crawl"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"github.com/jjmschofield/GoCrawl/internal/app/writers"
	"log"
	"net/url"
	"sync"
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

	out, wg := createWriter(*outFilePath)

	start := time.Now()

	crawledCache := caches.NewStrThreadSafe()
	processingCache := caches.NewStrThreadSafe()

	config := crawl.Config{
		CrawlWorkerCount: *workerCount,
		Caches: crawl.Caches{
			Crawled: &crawledCache,
			Processing: &processingCache,
		},
	}

	crawler := crawl.NewCrawler(crawl.Worker, out, config)
	counters := crawler.Crawl(*crawlUrl)

	wg.Wait()

	end := time.Now()

	fmt.Printf("Scrape Completed in %v ms \n", (end.UnixNano()-start.UnixNano())/int64(time.Millisecond))
	fmt.Printf(" Discovered: %v, \n Crawled: %v \n Parallel Crawls Peak: %v \n Scrape Queue Peak: %v \n Processing Peak: %v \n", counters.Discovered.Count(), counters.CrawlComplete.Count(), counters.Crawling.Peak(), counters.CrawlsQueued.Peak(), counters.Processing.Peak())
}

func createWriter(outFilePath string) (in chan pages.Page, waitGroup *sync.WaitGroup) {
	var wg sync.WaitGroup

	in = make(chan pages.Page)

	writer := writers.FileWriter{FilePath: outFilePath}
	go writer.Write(in, &wg)

	return in, &wg
}
