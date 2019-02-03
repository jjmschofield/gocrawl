package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"github.com/jjmschofield/GoCrawl/internal/app/caches"
	"github.com/jjmschofield/GoCrawl/internal/app/counters"
	"log"
	"net/url"
	"sync"
)

type Crawler struct {
	Config   CrawlerConfig
	counters Counters
	channels channels
	caches   lockingCaches
	workers  workers
	wg       sync.WaitGroup
}

type CrawlerConfig struct {
	CrawlWorkerCount int
}

type workers struct {
	crawl CrawlQueueWorker
}

type channels struct {
	workerIn  chan pages.Page
	workerOut chan CrawlWorkerResult
	out       chan pages.Page
}

type lockingCaches struct {
	processing caches.LockingStr
	crawled    caches.LockingStr
}

type Counters struct {
	Discovered    counters.AtomicInt64
	Processing    counters.AtomicInt64
	Crawling      counters.AtomicInt64
	CrawlComplete counters.AtomicInt64
	CrawlsQueued  counters.AtomicInt64
}

func NewCrawler(crawlWorker CrawlQueueWorker, out chan pages.Page, config CrawlerConfig) Crawler {
	return Crawler{
		Config: config,
		channels: channels{
			workerIn:  make(chan pages.Page, config.CrawlWorkerCount),
			workerOut: make(chan CrawlWorkerResult, config.CrawlWorkerCount),
			out:       out,
		},
		caches: lockingCaches{
			processing: caches.NewLockingStr(),
			crawled:    caches.NewLockingStr(),
		},
		workers: workers{
			crawl: crawlWorker,
		},
	}
}

func (c *Crawler) Crawl(startUrl url.URL) Counters {
	c.startWorkers()

	c.startResultWorker()

	c.enqueueCrawl(pages.PageFromUrl(startUrl))

	c.wg.Wait()

	return c.counters
}

func (c *Crawler) startWorkers() {
	for i := 0; i < c.Config.CrawlWorkerCount; i++ {
		c.wg.Add(1)
		go c.workers.crawl(c.channels.workerIn, c.channels.workerOut, &c.counters.CrawlsQueued, &c.counters.Crawling, &c.wg)
	}
}

func (c *Crawler) startResultWorker(){
	c.wg.Add(1)
	go c.crawlResultWorker()
}

func (c *Crawler) crawlResultWorker() {
	defer c.wg.Done()

	for result := range c.channels.workerOut {
		c.caches.crawled.Add(result.crawled.Id)
		c.counters.CrawlComplete.Add(1)

		c.counters.Processing.Sub(1)
		c.caches.processing.Remove(result.crawled.Id)

		c.enqueueNewPages(result.result.OutPages.Internal)

		log.Printf("crawled crawled %s Discovered: %v, Processing: %v, In Scrape Queue: %v, Crawling: %v, Scrape Complete: %v", result.crawled.URL.String(), c.counters.Discovered.Count(), c.counters.Processing.Count(), c.counters.CrawlsQueued.Count(), c.counters.Crawling.Count(), c.counters.CrawlComplete.Count())

		c.channels.out <- result.crawled

		if !c.hasWorkRemaining() {
			c.closeChannels()
		}
	}
}

func (c *Crawler) enqueueCrawl(page pages.Page) {
	c.counters.Discovered.Add(1)
	c.counters.Processing.Add(1)

	c.caches.processing.Add(page.Id)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.counters.CrawlsQueued.Add(1)
		c.channels.workerIn <- page
	}()
}

func (c *Crawler) enqueueNewPages(pageMap []pages.Page) {
	for _, page := range pageMap {
		if !c.caches.crawled.Has(page.Id) && !c.caches.processing.Has(page.Id) {
			c.enqueueCrawl(page)
		}
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
	close(c.channels.out)
}
