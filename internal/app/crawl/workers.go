package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"github.com/jjmschofield/GoCrawl/internal/pkg/counters"
	"sync"
)

type PageQueueWorker func(chan pages.Page, chan PageCrawlResult, *counters.AtomicInt64, *counters.AtomicInt64, *sync.WaitGroup)

type PageCrawlResult struct {
	crawled    pages.Page
	discovered map[string]pages.Page
}

// TODO - someone implementing this has too much responsibility to update counters?
func PageCrawler(queue chan pages.Page, result chan PageCrawlResult, qCount *counters.AtomicInt64, queueCount *counters.AtomicInt64, wg *sync.WaitGroup) {
	defer wg.Done()
	for page := range queue {
		discovered := make(map[string]pages.Page)

		queueCount.Add(1)
		qCount.Sub(1)

		_, err := page.FetchLinks(pages.FetchPageBody, pages.ReadHrefs)

		if err != nil {
			page.Err = err
		} else {
			for _, link := range page.OutLinks.Internal {
				discoveredPage := pages.PageFromUrl(link.ToURL)
				discovered[discoveredPage.Id] = discoveredPage
			}
		}

		result <- PageCrawlResult{
			crawled:    page,
			discovered: discovered,
		}

		queueCount.Sub(1)
	}
}
