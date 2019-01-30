package crawl

import (
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"log"
	"sync"
)

type PageQueueWorker func(chan pages.Page, chan PageCrawlResult, *sync.WaitGroup)

type PageCrawlResult struct {
	crawledPage pages.Page
	discoveredPages map[string]pages.Page
}

func PageCrawler(crawlQueue chan pages.Page, crawlResult chan PageCrawlResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for page := range crawlQueue {
		hrefs, err := page.FetchHrefs(pages.FetchPageBody, pages.ReadHrefs)

		result := PageCrawlResult{
			crawledPage : page,
			discoveredPages: make(map[string]pages.Page),
		}

		if err != nil {
			result.crawledPage.Err = err
			log.Printf("error getting page %s %s %v",page.Id, page.URL.String(), err)
			crawlResult <- result
		}

		discoveredLinks := links.FromHrefs(page.URL, hrefs)

		for _, link := range discoveredLinks {
			if link.Type == links.InternalPageType {
				discoveredPage := pages.PageFromUrl(link.ToURL)
				result.crawledPage.AppendOutInternalPage(discoveredPage)
				result.discoveredPages[discoveredPage.Id] = discoveredPage
			}

			result.crawledPage.AppendOutLink(link)
		}

		crawlResult <- result
	}
}
