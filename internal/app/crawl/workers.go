package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"github.com/jjmschofield/GoCrawl/internal/app/counters"
	"sync"
)

type CrawlWorkerResult struct {
	crawled pages.Page
	result  PageCrawlResult
}

type CrawlQueueWorker func(chan pages.Page, chan CrawlWorkerResult, *counters.AtomicInt64, *counters.AtomicInt64, *sync.WaitGroup)

// TODO - someone implementing this has too much responsibility to update counters?
func CrawlWorker(queue chan pages.Page, out chan CrawlWorkerResult, qCount *counters.AtomicInt64, workCount *counters.AtomicInt64, wg *sync.WaitGroup) {
	defer wg.Done()
	for page := range queue {
		workCount.Add(1)
		qCount.Sub(1)

		crawlResult, err := Crawl(page.URL)

		workerResult := CrawlWorkerResult{
			crawled: page,
			result:  crawlResult,
		}

		if err != nil {
			workerResult.crawled.Err = err
		} else {
			workerResult.crawled.OutLinks = crawlResult.OutLinks
			workerResult.crawled.OutPages = crawlResult.OutPages
		}

		out <- workerResult
		workCount.Sub(1)
	}
}
