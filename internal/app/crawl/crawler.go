package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
)

func FromUrl(startUrl url.URL, crawlWorker PageQueueWorker, nCrawlWorkers int) (crawledPages map[string]pages.Page) {
	var wg sync.WaitGroup

	linkQueue := make(chan links.Link)
	var linkQueueLen int64

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

			for _, link := range page.OutLinks.Internal {
				atomic.AddInt64(&linkQueueLen, 1)
				linkQueue <- link
			}

			if crawlQueueLen < 1 && linkQueueLen < 1 {
				close(crawlQueue)
				close(crawlResults)
				close(linkQueue)
			}
		}
	}()

	go func() {
		for link := range linkQueue {
			log.Printf("%v", link)
			atomic.AddInt64(&linkQueueLen, -1)
		}
	}()

	atomic.AddInt64(&crawlQueueLen, 1)
	crawlQueue <- pages.PageFromUrl(startUrl)

	wg.Wait()

	return crawled
}
