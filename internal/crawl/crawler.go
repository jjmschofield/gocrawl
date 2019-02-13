package crawl

import (
	"github.com/go-redis/redis"
	"github.com/jjmschofield/gocrawl/internal/caches"
	"github.com/jjmschofield/gocrawl/internal/counters"
	"github.com/jjmschofield/gocrawl/internal/pages"
	"github.com/jjmschofield/gocrawl/internal/scrape"
	"github.com/jjmschofield/gocrawl/internal/writers"
	"log"
	"net/url"
	"sync"
)

//go:generate counterfeiter . Crawler
type Crawler interface {
	Crawl(startUrl url.URL) Counters
}

type PageCrawler struct {
	Config   Config
	caches   Caches
	channels WorkerChannels
	counters Counters
	scraper  scrape.Scraper
	worker   QueueWorker
	writer   writers.Writer
}

type Config struct {
	Caches      Caches
	WorkerCount int
	Scraper     scrape.Scraper
	Worker      QueueWorker
	Writer      writers.Writer
}

type Caches struct {
	Processing caches.ThreadSafeCache
	Crawled    caches.ThreadSafeCache
}

type Counters struct {
	Discovered    counters.AtomicInt64 // Pages discovered so far
	Processing    counters.AtomicInt64 // Pages that we need to complete processing
	Crawling      counters.AtomicInt64 // Pages that we are currently crawling
	CrawlComplete counters.AtomicInt64 // Pages that we have CrawledId
	CrawlsQueued  counters.AtomicInt64 // Pages currently queued for crawling
}

func (c *PageCrawler) Crawl(startUrl url.URL) Counters {
	var wg sync.WaitGroup

	go c.writer.Start(c.channels.Write)

	go c.resultReceiver()

	c.startWorkers(c.Config.WorkerCount, &wg)

	c.enqueueUrl(startUrl)

	wg.Wait()

	return c.counters
}

func (c *PageCrawler) startWorkers(workerCount int, wg *sync.WaitGroup) {
	wg.Add(workerCount)
	for i := 0; i < c.Config.WorkerCount; i++ {
		go c.worker.Start(c.channels, &c.counters.CrawlsQueued, &c.counters.Crawling, wg)
	}
}

func (c *PageCrawler) enqueueJob(job WorkerJob) {
	if !c.caches.Crawled.Has(job.Id) && !c.caches.Processing.Has(job.Id) {
		c.counters.Processing.Add(1)
		c.counters.Discovered.Add(1)

		c.caches.Processing.Add(job.Id)

		go func() {
			c.counters.CrawlsQueued.Add(1)
			c.channels.In <- job
		}()
	}
}

func (c *PageCrawler) resultReceiver() {
	for result := range c.channels.Out {
		c.caches.Crawled.Add(result.CrawledId)
		c.counters.CrawlComplete.Add(1)

		c.enqueuePageGroup(result.Result.OutPages)

		c.counters.Processing.Sub(1)
		c.caches.Processing.Remove(result.CrawledId)

		log.Printf(
			"Crawled %s Discovered: %v, Processing: %v, In Scrape Queue: %v, Scraping: %v, Scrape Complete: %v",
			result.CrawledId,
			c.counters.Discovered.Count(),
			c.counters.Processing.Count(),
			c.counters.CrawlsQueued.Count(),
			c.counters.Crawling.Count(),
			c.counters.CrawlComplete.Count())

		if !c.hasWorkRemaining() {
			c.closeChannels()
		}
	}
}

func (c *PageCrawler) enqueueUrl(srcUrl url.URL) {
	pageId, normalizedUrl := pages.CalcPageId(srcUrl)

	job := WorkerJob{
		Id:  pageId,
		URL: normalizedUrl,
	}

	c.enqueueJob(job)
}

func (c *PageCrawler) enqueuePageGroup(pageGroup pages.PageGroup) {
	for pageId, href := range pageGroup.Internal {
		pageUrl, _ := url.Parse(href)

		job := WorkerJob{
			Id:  pageId,
			URL: *pageUrl,
		}

		c.enqueueJob(job)
	}
}

func (c *PageCrawler) hasWorkRemaining() bool {
	if c.counters.Processing.Count() > 0 {
		return true
	}

	if c.counters.CrawlsQueued.Count() > 0 {
		return true
	}

	if c.counters.Crawling.Count() > 0 {
		return true
	}

	return false
}

func (c *PageCrawler) closeChannels() {
	close(c.channels.In)
	close(c.channels.Out)
	close(c.channels.Write)
}

func NewPageCrawler(config Config) PageCrawler {
	return PageCrawler{
		Config: config,
		caches: config.Caches,
		channels: WorkerChannels{
			In:    make(chan WorkerJob),
			Out:   make(chan WorkerResult),
			Write: make(chan pages.Page),
		},
		worker: config.Worker,
		writer: config.Writer,
	}
}

func NewDefaultPageCrawler(workerCount int, filePath string) PageCrawler {
	crawledCache := caches.NewStr()
	processingCache := caches.NewStr()

	config := Config{
		Caches: Caches{
			Crawled:    &crawledCache,
			Processing: &processingCache,
		},
		Worker:      &Worker{Scraper: scrape.PageScraper{}},
		WorkerCount: workerCount,
		Writer:      &writers.FileWriter{FilePath: filePath},
	}

	crawler := NewPageCrawler(config)

	return crawler
}

func NewRedisPageCrawler(workerCount int, filePath string, redisAddr string) PageCrawler {
	options := &redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	client := redis.NewClient(options)
	client.FlushAll()

	err := client.Close()
	if err != nil {
		panic(err)
	}

	crawledCache := caches.NewStrRedis("crawled", options)
	processingCache := caches.NewStrRedis("processing", options)

	config := Config{
		Caches: Caches{
			Crawled:    &crawledCache,
			Processing: &processingCache,
		},
		Worker:      &Worker{Scraper: scrape.PageScraper{}},
		WorkerCount: workerCount,
		Writer:      &writers.FileWriter{FilePath: filePath},
	}

	crawler := NewPageCrawler(config)

	return crawler
}

