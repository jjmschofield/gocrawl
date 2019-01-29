package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"sync"
)

type PageQueueWorker func(chan pages.Page, chan pages.Page, sync.WaitGroup)

func PageCrawler(crawlQueue chan pages.Page, crawlResult chan pages.Page, wg sync.WaitGroup) {
	defer wg.Done()

	for page := range crawlQueue {
		hrefs, err := page.FetchHrefs(pages.FetchPageBody, pages.ReadHrefs)

		if err != nil {
			page.Err = err
			crawlResult <- page
		}

		discoveredLinks := links.FromHrefs(page.URL, hrefs)

		for _, link := range discoveredLinks {
			page.AppendOutLink(link)
		}

		crawlResult <- page
	}
}
