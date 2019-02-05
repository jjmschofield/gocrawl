package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/counters"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"github.com/jjmschofield/GoCrawl/internal/app/scrape"
	"net/url"
	"sync"
)

type WorkerJob struct {
	pageId string
	pageUrl url.URL
}

type WorkerResult struct {
	crawled pages.Page
	result  scrape.Result
}

type QueueWorker func(queue chan WorkerJob, out chan WorkerResult, workCount *counters.AtomicInt64, wg *sync.WaitGroup)

func Worker(queue chan WorkerJob, out chan WorkerResult, workCount *counters.AtomicInt64, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range queue {
		workCount.Add(1)

		crawled := pages.PageFromUrl(job.pageUrl)

		scrapeResult, err := scrape.Scrape(crawled.URL)

		workerResult := WorkerResult{
			crawled: crawled,
			result:  scrapeResult,
		}

		if err != nil {
			workerResult.crawled.Err = err
		} else {
			workerResult.crawled.OutLinks = scrapeResult.OutLinks
			workerResult.crawled.OutPages = scrapeResult.OutPages
		}

		out <- workerResult
		workCount.Sub(1)
	}
}
