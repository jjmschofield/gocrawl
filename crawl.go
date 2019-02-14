package main

import (
	"flag"
	"fmt"
	"github.com/jjmschofield/gocrawl/internal/crawl"
	"github.com/jjmschofield/gocrawl/internal/writers"
	"log"
	"net/url"
	"time"
)

func main() {
	crawlUrlRaw := flag.String("url", "https://monzo.com", "an absolute url, including protocol and hostname")
	workerCount := flag.Int("workers", 50, "Number of crawl workers to run")
	outFilePath := flag.String("dir", "data", "A relative file path to send results to")
	redisAddr := flag.String("redis", "", "An optional redis address to make use of redis rather then in memory queues and caches eg: localhost:6379")

	flag.Parse()

	crawlUrl, err := url.Parse(*crawlUrlRaw)

	if err != nil {
		log.Panic(err)
	}
	start := time.Now()

	counters := Crawl(*crawlUrl, *workerCount, *outFilePath, *redisAddr)

	end := time.Now()

	fmt.Printf("Scrape Completed in %v ms \n", (end.UnixNano()-start.UnixNano())/int64(time.Millisecond))
	fmt.Printf(" Discovered: %v, \n Crawled: %v \n Parallel Crawls Peak: %v \n Scrape Queue Peak: %v \n Scraping Peak: %v \n", counters.Discovered.Count(), counters.Crawled.Count(), counters.Scraping.Peak(), counters.Queued.Peak(), counters.Crawling.Peak())
}

func Crawl(crawlUrl url.URL, workerCount int, outFilePath string, redisAddr string) crawl.Counters {
	var crawler crawl.PageCrawler

	if len(redisAddr) < 1 {
		crawler = crawl.NewDefaultPageCrawler(workerCount, outFilePath)
	} else{
		crawler = crawl.NewRedisPageCrawler(workerCount, outFilePath, redisAddr)
	}

	out := crawler.Crawl(crawlUrl)

	writer := writers.FileWriter{FilePath: outFilePath}

	writer.Start(out)

	return crawler.Counters
}
