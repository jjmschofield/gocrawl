package crawl

import (
	"github.com/go-redis/redis"
	"github.com/jjmschofield/gocrawl/internal/caches"
	"github.com/jjmschofield/gocrawl/internal/counters"
	"github.com/jjmschofield/gocrawl/internal/pages"
	"github.com/jjmschofield/gocrawl/internal/queue"
	"github.com/jjmschofield/gocrawl/internal/scrape"
	"log"
	"net/url"
)

//go:generate counterfeiter . Crawler
type Crawler interface {
	Crawl(startUrl url.URL) Counters
}

type PageCrawler struct {
	Config   Config
	Counters Counters
	caches   Caches
	out      chan pages.Page
	queue    queue.Queue
	scraper  scrape.Scraper
	worker   queue.QueueWorker
}

type Config struct {
	Caches      Caches
	Queue       queue.Queue
	Scraper     scrape.Scraper
	Worker      queue.QueueWorker
	WorkerCount int
}

type Caches struct {
	Crawled  caches.ThreadSafeCache
	Crawling caches.ThreadSafeCache
}

type Counters struct {
	Discovered counters.AtomicInt64 // Pages discovered so far
	Crawled    counters.AtomicInt64 // Pages that we have Crawled
	Crawling   counters.AtomicInt64 // Pages that we have either in queue or are being scraped
	Queued     *counters.AtomicInt64 // Pages currently queued for crawling
	Scraping   *counters.AtomicInt64 // Pages that we are actively being crawled right now crawling
}

func (c *PageCrawler) Crawl(startUrl url.URL) chan pages.Page {
	go c.start(startUrl)
	return c.out
}

func (c *PageCrawler) start(startUrl url.URL) {
	results, _ := c.queue.Start(c.worker, c.Config.WorkerCount)

	c.enqueueUrl(startUrl)

	for result := range results {
		c.enqueuePageGroup(result.Page.OutPages)

		c.caches.Crawled.Add(result.Page.Id)
		c.caches.Crawling.Remove(result.Page.Id)

		c.out <- result.Page

		c.Counters.Crawled.Add(1)
		c.Counters.Crawling.Sub(1)

		log.Printf(
			"Crawled %s Discovered: %v, Scraping: %v, In Queue: %v, Scraping: %v, Crawled: %v",
			result.Page.Id,
			c.Counters.Discovered.Count(),
			c.Counters.Crawling.Count(),
			c.Counters.Queued.Count(),
			c.Counters.Scraping.Count(),
			c.Counters.Crawled.Count())

		if !c.hasWorkRemaining() {
			c.close()
		}
	}
}

func (c *PageCrawler) enqueueJob(job queue.WorkerJob) {
	if !c.caches.Crawled.Has(job.Id) && !c.caches.Crawling.Has(job.Id) {
		c.Counters.Discovered.Add(1)
		c.Counters.Crawling.Add(1)

		c.caches.Crawling.Add(job.Id)

		err := c.queue.Push(job)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (c *PageCrawler) enqueueUrl(srcUrl url.URL) {
	pageId, normalizedUrl := pages.CalcPageId(srcUrl)

	job := queue.WorkerJob{
		Id:  pageId,
		URL: normalizedUrl,
	}

	c.enqueueJob(job)
}

func (c *PageCrawler) enqueuePageGroup(pageGroup pages.PageGroup) {
	for pageId, href := range pageGroup.Internal {
		pageUrl, _ := url.Parse(href)

		job := queue.WorkerJob{
			Id:  pageId,
			URL: *pageUrl,
		}

		c.enqueueJob(job)
	}
}

func (c *PageCrawler) hasWorkRemaining() bool {
	if c.Counters.Crawling.Count() > 0 {
		return true
	}

	if c.Counters.Queued.Count() > 0 {
		return true
	}

	if c.Counters.Scraping.Count() > 0 {
		return true
	}

	return false
}

func (c *PageCrawler) close() {
	c.queue.Stop()
	close(c.out)
}

func NewPageCrawler(config Config) PageCrawler {
	return PageCrawler{
		caches: config.Caches,
		Config: config,
		Counters: Counters{
			Scraping: config.Queue.Counters().Work,
			Queued:   config.Queue.Counters().Queue,
		},
		out:    make(chan pages.Page),
		queue:  config.Queue,
		worker: config.Worker,
	}
}

func NewDefaultPageCrawler(workerCount int, filePath string) PageCrawler {
	crawledCache := caches.NewStr()
	processingCache := caches.NewStr()

	config := Config{
		Caches: Caches{
			Crawled:  &crawledCache,
			Crawling: &processingCache,
		},
		Worker:      &queue.Worker{Scraper: scrape.PageScraper{}},
		WorkerCount: workerCount,
		Queue:       queue.NewBasicQueue(),
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
			Crawled:  &crawledCache,
			Crawling: &processingCache,
		},
		Worker:      &queue.Worker{Scraper: scrape.PageScraper{}},
		WorkerCount: workerCount,
		Queue:       queue.NewBasicQueue(),
	}

	crawler := NewPageCrawler(config)

	return crawler
}
