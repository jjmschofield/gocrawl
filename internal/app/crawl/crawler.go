package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"github.com/jjmschofield/GoCrawl/internal/pkg/caches"
	"github.com/jjmschofield/GoCrawl/internal/pkg/counters"
	"log"
	"net/url"
	"sync"
)

type Crawler struct {
	Config   CrawlerConfig
	counters Counters
	channels cChans
	caches   cCaches
	workers  cWorkers
	wg       sync.WaitGroup
}

type CrawlerConfig struct {
	CrawlWorkerCount int
}

type cWorkers struct {
	crawl PageQueueWorker
}

type cChans struct {
	crawlQueue  chan pages.Page
	crawlResult chan PageCrawlResult
	out         chan pages.Page
}

type cCaches struct {
	processing caches.StrLocking
	crawled    caches.StrLocking
}

type Counters struct {
	Discovered    counters.AtomicInt64
	Processing    counters.AtomicInt64
	Crawling      counters.AtomicInt64
	CrawlComplete counters.AtomicInt64
	CrawlsQueued  counters.AtomicInt64
}

func NewCrawler(crawlWorker PageQueueWorker, out chan pages.Page, config CrawlerConfig) Crawler {
	return Crawler{
		Config: config,
		channels: cChans{
			crawlQueue:  make(chan pages.Page, config.CrawlWorkerCount),
			crawlResult: make(chan PageCrawlResult, config.CrawlWorkerCount),
			out:         out,
		},
		caches: cCaches{
			processing: caches.NewStrBlocking(),
			crawled:    caches.NewStrBlocking(),
		},
		workers: cWorkers{
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
		go c.workers.crawl(
			c.channels.crawlQueue,
			c.channels.crawlResult,
			&c.counters.CrawlsQueued,
			&c.counters.Crawling,
			&c.wg,
		)
	}
}

func (c *Crawler) startResultWorker(){
	c.wg.Add(1)
	go c.crawlResultWorker()
}

func (c *Crawler) crawlResultWorker() {
	defer c.wg.Done()

	for result := range c.channels.crawlResult {
		c.caches.crawled.Add(result.crawled.Id)
		c.counters.CrawlComplete.Add(1)

		c.counters.Processing.Sub(1)
		c.caches.processing.Remove(result.crawled.Id)

		c.enqueueNewPages(result.discovered)

		log.Printf("page crawled %s Discovered: %v, Processing: %v, In Crawl Queue: %v, Crawling: %v, Crawl Complete: %v", result.crawled.URL.String(), c.counters.Discovered.Count(), c.counters.Processing.Count(), c.counters.CrawlsQueued.Count(), c.counters.Crawling.Count(), c.counters.CrawlComplete.Count())

		c.channels.out <- result.crawled

		if !c.hasWorkRemaining() {
			c.closeChannels()
		}
	}
}

func (c *Crawler) enqueueCrawl(page pages.Page) {
	c.counters.Discovered.Add(1)
	c.counters.Processing.Add(1)
	c.counters.CrawlsQueued.Add(1)

	c.caches.processing.Add(page.Id)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.channels.crawlQueue <- page
	}()
}

func (c *Crawler) enqueueNewPages(pageMap map[string]pages.Page) {
	for id, page := range pageMap {
		if !c.caches.crawled.Has(id) && !c.caches.processing.Has(id) {
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
	close(c.channels.crawlQueue)
	close(c.channels.crawlResult)
	close(c.channels.out)
}
