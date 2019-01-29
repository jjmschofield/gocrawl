package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"net/url"
	"sync"
	"sync/atomic"
)

func FromUrl(startUrl url.URL, crawlWorker PageQueueWorker, nCrawlWorkers int) (crawledPages map[string]pages.Page) {
	var wg sync.WaitGroup

	crawlQueue := make(chan pages.Page)
	crawlResults := make(chan pages.Page)
	var crawlQueueLen int64

	crawled := make(map[string]pages.Page)
	var crawledMutex sync.Mutex

	for i := 0; i < nCrawlWorkers; i++ {
		wg.Add(1)
		go crawlWorker(crawlQueue, crawlResults, &wg)
	}

	go func() {
		for page := range crawlResults {
			crawledMutex.Lock()
			crawled[page.Id] = page
			crawledMutex.Unlock()

			atomic.AddInt64(&crawlQueueLen, -1)

			if crawlQueueLen < 1 {
				close(crawlQueue)
				close(crawlResults)
			}
		}
	}()

	atomic.AddInt64(&crawlQueueLen, 1)
	crawlQueue <- pages.PageFromUrl(startUrl)

	wg.Wait()

	return crawled
}
