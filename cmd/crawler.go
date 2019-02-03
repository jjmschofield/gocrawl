package main

import (
	"flag"
	"fmt"
	"github.com/jjmschofield/GoCrawl/internal/app/crawl"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"github.com/jjmschofield/GoCrawl/internal/app/writers"
	"log"
	"net/url"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	crawlUrlRaw := flag.String("url", "https://monzo.com", "an absolute url eg http://www.google.co.uk")
	workerCount := flag.Int("workers", 10, "Number of crawl workers to run")
	outFilePath := flag.String("file", "", "A file path to send results to, if not set will print to stdout")

	flag.Parse()

	crawlUrl, err := url.Parse(*crawlUrlRaw)

	if err != nil {
		log.Panic(err)
	}

	var wg sync.WaitGroup

	out := make(chan pages.Page)

	if len(*outFilePath)  > 0 {
		writer := writers.FileWriter{FilePath: *outFilePath}
		go writer.Write(out, &wg)
	} else{
		go writers.StdoutWriter(out, &wg)
	}

	crawler := crawl.NewCrawler(crawl.CrawlWorker, out, crawl.CrawlerConfig{CrawlWorkerCount: *workerCount})
	counters := crawler.Crawl(*crawlUrl)
	
	wg.Wait()

	end := time.Now()

	fmt.Printf("Crawl Completed in %v ms \n", (end.UnixNano()-start.UnixNano())/int64(time.Millisecond))
	fmt.Printf(" Discovered: %v, \n Crawled: %v \n Parallel Crawls Peak: %v \n Crawl Queue Peak: %v \n Processing Peak: %v \n", counters.Discovered.Count(), counters.CrawlComplete.Count(), counters.Crawling.Peak(), counters.CrawlsQueued.Peak(), counters.Processing.Peak())
}