package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/caches"
	"github.com/jjmschofield/GoCrawl/internal/app/counters"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"github.com/jjmschofield/GoCrawl/internal/app/writers"
	"log"
	"net/url"
	"sync"
)

type Crawler struct {
	Config   Config
	counters Counters
	channels channels
	caches   Caches
	workers  workers
	writer   writers.Writer
	wg       sync.WaitGroup
}

type Config struct {
	CrawlWorkerCount int
	Caches           Caches
}

type Caches struct {
	Processing caches.ThreadSafeCache
	Crawled    caches.ThreadSafeCache
}

type workers struct {
	crawl QueueWorker
}

type channels struct {
	workerIn  chan WorkerJob
	workerOut chan WorkerResult
	write     chan pages.Page
}

type Counters struct {
	Discovered    counters.AtomicInt64
	Processing    counters.AtomicInt64
	Crawling      counters.AtomicInt64
	CrawlComplete counters.AtomicInt64
	CrawlsQueued  counters.AtomicInt64
}

func NewCrawler(worker QueueWorker, writer writers.Writer, config Config) Crawler {
	return Crawler{
		Config: config,
		channels: channels{
			workerIn:  make(chan WorkerJob),
			workerOut: make(chan WorkerResult),
			write:     make(chan pages.Page),
		},
		caches: config.Caches,
		workers: workers{
			crawl: worker,
		},
		writer: writer,
	}
}

func NewDefaultCrawler(workerCount int, filePath string) Crawler {
	crawledCache := caches.NewStrThreadSafe()
	processingCache := caches.NewStrThreadSafe()

	config := Config{
		CrawlWorkerCount: workerCount,
		Caches: Caches{
			Crawled:    &crawledCache,
			Processing: &processingCache,
		},
	}

	writer := writers.FileWriter{FilePath: filePath}

	crawler := NewCrawler(Worker, &writer, config)

	return crawler
}

func (c *Crawler) Crawl(startUrl url.URL) Counters {
	c.startWriter()

	c.startWorkers()

	c.startResultHandler()

	c.enqueueUrl(startUrl)

	c.wg.Wait()

	return c.counters
}

func (c *Crawler) startWriter() {
	go c.writer.Start(c.channels.write)
}

func (c *Crawler) startWorkers() {
	for i := 0; i < c.Config.CrawlWorkerCount; i++ {
		c.wg.Add(1)
		go c.workers.crawl(c.channels.workerIn, c.channels.workerOut, c.channels.write, &c.counters.Crawling, &c.wg)
	}
}

func (c *Crawler) startResultHandler() {
	c.wg.Add(1)
	go c.crawlResultHandler()
}

func (c *Crawler) crawlResultHandler() {
	defer c.wg.Done()

	for result := range c.channels.workerOut {
		c.caches.Crawled.Add(result.crawled)
		c.counters.CrawlComplete.Add(1)

		c.enqueuePageGroup(result.result.OutPages)

		c.counters.Processing.Sub(1)
		c.caches.Processing.Remove(result.crawled)

		log.Printf(
			"Crawled %s Discovered: %v, Processing: %v, In Scrape Queue: %v, Scraping: %v, Scrape Complete: %v",
			result.crawled,
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

func (c *Crawler) enqueueJob(job WorkerJob) {
	if !c.caches.Crawled.Has(job.pageId) && !c.caches.Processing.Has(job.pageId) {
		c.counters.Discovered.Add(1)
		c.counters.Processing.Add(1)

		c.caches.Processing.Add(job.pageId)

		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			c.counters.CrawlsQueued.Add(1)
			c.channels.workerIn <- job
			c.counters.CrawlsQueued.Sub(1)
		}()
	}
}

func (c *Crawler) enqueueUrl(srcUrl url.URL) {
	pageId, normalizedUrl := pages.CalcPageId(srcUrl)

	job := WorkerJob{
		pageId:  pageId,
		pageUrl: normalizedUrl,
	}

	c.enqueueJob(job)
}

func (c *Crawler) enqueuePageGroup(pageGroup pages.PageGroup) {
	for pageId, href := range pageGroup.Internal {
		pageUrl, _ := url.Parse(href)

		job := WorkerJob{
			pageId:  pageId,
			pageUrl: *pageUrl,
		}
		c.enqueueJob(job)
	}
}

func (c *Crawler) hasWorkRemaining() bool {
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

func (c *Crawler) closeChannels() {
	close(c.channels.workerIn)
	close(c.channels.workerOut)
	close(c.channels.write)
}
