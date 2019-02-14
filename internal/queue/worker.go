package queue

import (
	"github.com/jjmschofield/gocrawl/internal/counters"
	"github.com/jjmschofield/gocrawl/internal/pages"
	"github.com/jjmschofield/gocrawl/internal/scrape"
	"net/url"
)

//go:generate counterfeiter . QueueWorker
type QueueWorker interface {
	Start(chans Channels, qCounter *counters.AtomicInt64, workCounter *counters.AtomicInt64)
}

type WorkerJob struct {
	Id  string
	URL url.URL // TODO - change to href
}

type WorkerResult struct {
	Page   pages.Page
}

type Worker struct {
	Scraper scrape.Scraper
}

func (w *Worker) Start(chans Channels, qCounter *counters.AtomicInt64, workCounter *counters.AtomicInt64) {
	for job := range chans.Jobs {
		workCounter.Add(1)

		page := pages.PageFromUrl(job.URL)

		scrapeResult, err := w.Scraper.Scrape(page.URL)

		if err != nil {
			page.Err = err
		} else {
			page.OutLinks = scrapeResult.OutLinks
			page.OutPages = scrapeResult.OutPages
		}

		workerResult := WorkerResult{
			Page:   page,
		}

		chans.Results <- workerResult
		workCounter.Sub(1)
	}
}
