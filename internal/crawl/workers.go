package crawl

import (
	"github.com/jjmschofield/gocrawl/internal/counters"
	"github.com/jjmschofield/gocrawl/internal/pages"
	"github.com/jjmschofield/gocrawl/internal/scrape"
	"net/url"
	"sync"
)

//go:generate counterfeiter . QueueWorker
type QueueWorker interface {
	Start(chans WorkerChannels, qCounter *counters.AtomicInt64, workCounter *counters.AtomicInt64, wg *sync.WaitGroup)
}

type WorkerJob struct {
	Id  string
	URL url.URL // TODO - change to href
}

type WorkerResult struct {
	CrawledId string
	Result    scrape.Result
}

type WorkerChannels struct {
	In    chan WorkerJob
	Out   chan WorkerResult
	Write chan pages.Page
}

type Worker struct {
	Scraper scrape.Scraper
}

// TODO - remove worker group
// TODO - remove writer
func (w *Worker) Start(chans WorkerChannels, queueCounter *counters.AtomicInt64, workCounter *counters.AtomicInt64, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range chans.In {
		workCounter.Add(1)
		queueCounter.Sub(1)

		page := pages.PageFromUrl(job.URL)

		scrapeResult, err := w.Scraper.Scrape(page.URL)

		if err != nil {
			page.Err = err
		} else {
			page.OutLinks = scrapeResult.OutLinks
			page.OutPages = scrapeResult.OutPages
		}

		chans.Write <- page

		workerResult := WorkerResult{
			CrawledId: page.Id,
			Result:    scrapeResult,
		}

		workCounter.Sub(1)
		chans.Out <- workerResult
	}
}
