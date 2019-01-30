package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"log"
	"net/url"
	"sync"
)

// TODO - when we encounter pages with few links we loose the benefit of concurrency eg http://monzo.com/blog/authors/kate-hollowood/11
func FromUrl(startUrl url.URL, crawlWorker PageQueueWorker, nCrawlWorkers int) (crawledPages map[string]pages.Page) {
	var wg sync.WaitGroup

	crawlQueue := make(chan pages.Page, nCrawlWorkers)
	crawlResults := make(chan PageCrawlResult, nCrawlWorkers)

	inProgress := make(map[string]interface{})
	crawled := make(map[string]pages.Page)

	for i := 0; i < nCrawlWorkers; i++ {
		wg.Add(1)
		go crawlWorker(crawlQueue, crawlResults, &wg)
	}

	go func() {
		for crawlResult := range crawlResults {
			crawled[crawlResult.crawledPage.Id] = crawlResult.crawledPage
			delete(inProgress, crawlResult.crawledPage.Id)

			for _, discoveredPage := range crawlResult.discoveredPages {
				_, hasBeenCrawled := crawled[discoveredPage.Id]
				_, isQueued := inProgress[discoveredPage.Id]

				if !hasBeenCrawled && !isQueued {
					inProgress[discoveredPage.Id] = nil
					go enqueuePage(discoveredPage, crawlQueue)
				}
			}

			log.Printf("page crawled %s %s In Progress: %v, Complete: %v", crawlResult.crawledPage.URL.String(), crawlResult.crawledPage.Id, len(inProgress), len(crawled))

			if len(inProgress) < 1 {
				close(crawlQueue)
				close(crawlResults)
			}
		}
	}()

	go enqueuePage(pages.PageFromUrl(startUrl), crawlQueue)

	wg.Wait()

	return crawled
}

func enqueuePage(page pages.Page, queue chan pages.Page){
	queue <- page
}