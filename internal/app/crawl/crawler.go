package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"net/url"
	"sync"
)

func FromUrl(startUrl url.URL, crawlWorker PageQueueWorker, nCrawlWorkers int) (crawledPages map[string]pages.Page) {
	var wg sync.WaitGroup

	crawlQueue := make(chan pages.Page)
	crawlResults := make(chan pages.Page)

	for i := 0; i < nCrawlWorkers; i++ {
		wg.Add(1)
		go crawlWorker(crawlQueue, crawlResults, wg)
	}

	go func() {
		for page := range crawlResults {
			page.Print()
			close(crawlQueue)
			close(crawlResults)
		}
	}()

	crawlQueue <- pages.PageFromUrl(startUrl)

	wg.Wait()

	return
}
