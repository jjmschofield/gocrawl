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

	crawlQueue := make(chan pages.Page)
	crawlResults := make(chan pages.Page)

	inProgress := make(map[string]pages.Page)
	var inProgressMutex sync.Mutex

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

			inProgressMutex.Lock()
			delete(inProgress, page.Id)

			for _, link := range page.OutLinks.Internal {
				newPage := pages.PageFromUrl(link.ToURL)

				_, hasBeenCrawled := crawled[newPage.Id]
				_, isQueued := inProgress[newPage.Id]

				if !hasBeenCrawled && !isQueued {

					inProgress[newPage.Id] = newPage

					go func() {
						crawlQueue <- newPage
					}()
				}
			}
			inProgressMutex.Unlock()

			log.Printf("page crawled %s In Progress: %v, Complete: %v", page.URL.String(), len(inProgress), len(crawled))

			if len(inProgress) < 1{
				close(crawlQueue)
				close(crawlResults)
			}
		}
	}()
	
	crawlQueue <- pages.PageFromUrl(startUrl) // TODO - find a way to split this into nice enqueue / producer functions to make the code sane

	wg.Wait()

	return crawled
}
