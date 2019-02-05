package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/caches"
	"github.com/jjmschofield/GoCrawl/internal/app/counters"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"log"
	"net/url"
	"sync"
)

type Crawler struct {
	Config   Config
	counters Counters
	channels channels
	caches   lockingCaches
	workers  workers
	wg       sync.WaitGroup
}

type Config struct {
	CrawlWorkerCount int
}

type workers struct {
	crawl QueueWorker
}

type channels struct {
	workerIn  chan WorkerJob
	workerOut chan WorkerResult
	out       chan pages.Page
}

type lockingCaches struct {
	processing caches.StrThreadSafe
	crawled    caches.StrThreadSafe
}

type Counters struct {
	Discovered    counters.AtomicInt64
	Processing    counters.AtomicInt64
	Crawling      counters.AtomicInt64
	CrawlComplete counters.AtomicInt64
	CrawlsQueued  counters.AtomicInt64
}

func NewCrawler(crawlWorker QueueWorker, out chan pages.Page, config Config) Crawler {
	return Crawler{
		Config: config,
		channels: channels{
			workerIn:  make(chan WorkerJob),
			workerOut: make(chan WorkerResult),
			out:       out,
		},
		caches: lockingCaches{
			processing: caches.NewStrThreadSafe(),
			crawled:    caches.NewStrThreadSafe(),
		},
		workers: workers{
			crawl: crawlWorker,
		},
	}
}

func (c *Crawler) Crawl(startUrl url.URL) Counters {
	c.startWorkers()

	c.startResultWorker()

	c.enqueueUrl(startUrl)

	c.wg.Wait()

	return c.counters
}

func (c *Crawler) startWorkers() {
	for i := 0; i < c.Config.CrawlWorkerCount; i++ {
		c.wg.Add(1)
		go c.workers.crawl(c.channels.workerIn, c.channels.workerOut, &c.counters.Crawling, &c.wg)
	}
}

func (c *Crawler) startResultWorker() {
	c.wg.Add(1)
	go c.crawlResultWorker()
}

func (c *Crawler) crawlResultWorker() {
	defer c.wg.Done()

	for result := range c.channels.workerOut {
		c.caches.crawled.Add(result.crawled.Id)
		c.counters.CrawlComplete.Add(1)

		c.enqueuePageGroup(result.result.OutPages)

		c.counters.Processing.Sub(1)
		c.caches.processing.Remove(result.crawled.Id)

		log.Printf("crawled %s Discovered: %v, Processing: %v, In Scrape Queue: %v, Scraping: %v, Scrape Complete: %v", result.crawled.URL.String(), c.counters.Discovered.Count(), c.counters.Processing.Count(), c.counters.CrawlsQueued.Count(), c.counters.Crawling.Count(), c.counters.CrawlComplete.Count())

		c.channels.out <- result.crawled

		if !c.hasWorkRemaining() {
			c.closeChannels()
		}
	}
}

func (c *Crawler) enqueueJob(job WorkerJob) {
	if !c.caches.crawled.Has(job.pageId) && !c.caches.processing.Has(job.pageId) {
		c.counters.Discovered.Add(1)
		c.counters.Processing.Add(1)

		c.caches.processing.Add(job.pageId)

		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			c.counters.CrawlsQueued.Add(1)
			c.channels.workerIn <- job
			c.counters.CrawlsQueued.Sub(1)
		}()
	}
}

func (c *Crawler) enqueueUrl(srcUrl url.URL){
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
	close(c.channels.out)
}
