package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"github.com/jjmschofield/GoCrawl/internal/pkg/counters"
	"sync"
)

type CrawlWorkerResult struct {
	crawled pages.Page
	result  PageCrawlResult
}

type PageQueueWorker func(chan pages.Page, chan CrawlWorkerResult, *counters.AtomicInt64, *counters.AtomicInt64, *sync.WaitGroup)

// TODO - someone implementing this has too much responsibility to update counters?
func PageCrawler(queue chan pages.Page, out chan CrawlWorkerResult, qCount *counters.AtomicInt64, workCount *counters.AtomicInt64, wg *sync.WaitGroup) {
	defer wg.Done()
	for page := range queue {
		workCount.Add(1)
		qCount.Sub(1)

		result, err := Crawl(page.URL)

		if err != nil {
			page.Err = err
		}

		page.OutLinks = links.ToLinkGroup(result.OutLinks)
		page.OutPages = result.OutPages

		out <- CrawlWorkerResult{
			crawled: page,
			result:  result,
		}

		workCount.Sub(1)
	}
}
