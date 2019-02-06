package crawl_test

import (
	"errors"
	"github.com/jjmschofield/gocrawl/internal/counters"
	. "github.com/jjmschofield/gocrawl/internal/crawl"
	"github.com/jjmschofield/gocrawl/internal/links"
	"github.com/jjmschofield/gocrawl/internal/pages"
	"github.com/jjmschofield/gocrawl/internal/scrape"
	"github.com/jjmschofield/gocrawl/internal/scrapeapefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"sync"
)

var _ = Describe("Workers", func() {
	var channels WorkerChannels
	var queueCounter counters.AtomicInt64
	var workCounter counters.AtomicInt64

	BeforeEach(func() {
		channels = WorkerChannels{
			In:    make(chan WorkerJob),
			Out:   make(chan WorkerResult),
			Write: make(chan pages.Page),
		}

		workCounter = counters.AtomicInt64{}
		queueCounter = counters.AtomicInt64{}
	})

	It("should scrape a page for the given url", func() {
		// Arrange
		var wg sync.WaitGroup
		wg.Add(1)

		fakeScraper := scrapefakes.FakeScraper{}

		fakeScraper.ScrapeReturns(createFakeScrapeResult(), nil)

		go writeHandlerFake(channels, func(page pages.Page) { return })

		underTest := Worker{Scraper: &fakeScraper}
		go underTest.Start(channels, &queueCounter, &workCounter, &wg)

		expected, _ := url.Parse("https://www.google.co.uk")

		// Act
		channels.In <- WorkerJob{
			Id:  "some id",
			URL: *expected,
		}

		// Assert
		outHandlerFake(channels, func(_ WorkerResult) {
			result := fakeScraper.ScrapeArgsForCall(0)
			Expect(fakeScraper.ScrapeCallCount()).To(Equal(1))
			Expect(result.String()).To(Equal(expected.String()))
		})
	})

	It("should output the crawl results to the out channel", func() {
		// Arrange
		var wg sync.WaitGroup
		wg.Add(1)

		targetUrl, _ := url.Parse("https://www.google.co.uk")

		expected := createFakeScrapeResult()

		fakeScraper := scrapefakes.FakeScraper{}
		fakeScraper.ScrapeReturns(expected, nil)

		go writeHandlerFake(channels, func(page pages.Page) { return })

		underTest := Worker{Scraper: &fakeScraper}
		go underTest.Start(channels, &queueCounter, &workCounter, &wg)

		// Act
		channels.In <- WorkerJob{
			Id:  "some id",
			URL: *targetUrl,
		}

		// Assert
		outHandlerFake(channels, func(result WorkerResult) {
			Expect(result.Result).To(Equal(expected))
		})
	})

	It("should output the crawled pages id the out channel", func() {
		// Arrange
		var wg sync.WaitGroup
		wg.Add(1)

		expected := "3f1437859f73b447885255a95afa99a1"
		targetUrl, _ := url.Parse("https://www.google.co.uk")

		fakeScraper := scrapefakes.FakeScraper{}
		fakeScraper.ScrapeReturns(createFakeScrapeResult(), nil)

		go writeHandlerFake(channels, func(page pages.Page) { return })

		underTest := Worker{Scraper: &fakeScraper}
		go underTest.Start(channels, &queueCounter, &workCounter, &wg)

		// Act
		channels.In <- WorkerJob{
			Id:  "some id",
			URL: *targetUrl,
		}

		// Assert
		outHandlerFake(channels, func(result WorkerResult) {
			Expect(result.CrawledId).To(Equal(expected))
		})
	})

	Context("when the scrape is successful", func() {
		It("should output the the scraped page to the write channel", func() {
			// Arrange
			var wg sync.WaitGroup
			wg.Add(1)

			targetUrl, _ := url.Parse("https://www.google.co.uk")
			fakeScrapeResult := createFakeScrapeResult()
			expected := pages.PageFromUrl(*targetUrl)
			expected.OutLinks = createFakeScrapeResult().OutLinks
			expected.OutPages = createFakeScrapeResult().OutPages

			fakeScraper := scrapefakes.FakeScraper{}

			fakeScraper.ScrapeReturns(fakeScrapeResult, nil)

			go outHandlerFake(channels, func(result WorkerResult) { return })

			underTest := Worker{Scraper: &fakeScraper}
			go underTest.Start(channels, &queueCounter, &workCounter, &wg)

			// Act
			channels.In <- WorkerJob{
				Id:  "some id",
				URL: *targetUrl,
			}

			// Assert
			writeHandlerFake(channels, func(page pages.Page) {
				Expect(page).To(Equal(expected))
			})
		})
	})

	Context("when the scrape is not successful", func() {
		It("should output the scraped pages along with the error to the write channel", func() {
			// Arrange
			var wg sync.WaitGroup
			wg.Add(1)

			targetUrl, _ := url.Parse("https://www.google.co.uk")
			expected := pages.PageFromUrl(*targetUrl)
			expected.Err = errors.New("some error")

			fakeScraper := scrapefakes.FakeScraper{}
			fakeScraper.ScrapeReturns(createFakeScrapeResult(), expected.Err)

			go outHandlerFake(channels, func(result WorkerResult) { return })

			underTest := Worker{Scraper: &fakeScraper}
			go underTest.Start(channels, &queueCounter, &workCounter, &wg)

			// Act
			channels.In <- WorkerJob{
				Id:  "some id",
				URL: *targetUrl,
			}

			// Assert
			writeHandlerFake(channels, func(page pages.Page) {
				Expect(page.Err).To(Equal(expected.Err))
			})
		})
	})
})

func outHandlerFake(channels WorkerChannels, callback func(page WorkerResult)) {
	for result := range channels.Out {
		callback(result)
		close(channels.In)
		close(channels.Out)
		close(channels.Write)
	}
}

func createFakeScrapeResult() scrape.Result {
	someUrl, _ := url.Parse("https://www.bbc.co.uk")
	return scrape.Result{
		OutPages: pages.PageGroup{
			Internal: map[string]string{"some id": "some url"},
		},
		OutLinks: map[string]links.Link{"some link id": links.NewAbsLink(*someUrl, *someUrl)},
	}
}

func writeHandlerFake(channels WorkerChannels, callback func(page pages.Page)) {
	for page := range channels.Write {
		callback(page)
	}
}
