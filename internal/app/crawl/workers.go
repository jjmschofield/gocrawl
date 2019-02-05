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
	crawled string
	result  scrape.Result
}

type QueueWorker func(queue chan WorkerJob, out chan WorkerResult, write chan pages.Page, workCount *counters.AtomicInt64, wg *sync.WaitGroup)

func Worker(queue chan WorkerJob, out chan WorkerResult, write chan pages.Page, workCount *counters.AtomicInt64, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range queue {
		workCount.Add(1)

		page := pages.PageFromUrl(job.pageUrl)

		scrapeResult, err := scrape.Scrape(page.URL)

		if err != nil {
			page.Err = err
		} else {
			page.OutLinks = scrapeResult.OutLinks
			page.OutPages = scrapeResult.OutPages
		}

		write <- page

		workerResult := WorkerResult{
			crawled: page.Id,
			result:  scrapeResult,
		}

		out <- workerResult

		workCount.Sub(1)
	}
}
