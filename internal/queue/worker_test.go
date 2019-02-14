package queue_test

import (
	"errors"
	"github.com/jjmschofield/gocrawl/internal/counters"
	"github.com/jjmschofield/gocrawl/internal/links"
	"github.com/jjmschofield/gocrawl/internal/pages"
	. "github.com/jjmschofield/gocrawl/internal/queue"
	"github.com/jjmschofield/gocrawl/internal/scrape"
	"github.com/jjmschofield/gocrawl/internal/scrape/scrapefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("Worker", func() {
	var channels Channels
	var queueCounter counters.AtomicInt64
	var workCounter counters.AtomicInt64

	BeforeEach(func() {
		channels = Channels{
			Jobs:    make(chan WorkerJob),
			Results: make(chan WorkerResult),
		}

		workCounter = counters.AtomicInt64{}
		queueCounter = counters.AtomicInt64{}
	})

	It("should scrape a page for the given url", func() {
		// Arrange
		fakeScraper := scrapefakes.FakeScraper{}

		fakeScraper.ScrapeReturns(createFakeScrapeResult(), nil)

		underTest := Worker{Scraper: &fakeScraper}
		go underTest.Start(channels, &queueCounter, &workCounter)

		expected, _ := url.Parse("https://www.google.co.uk")

		// Act
		channels.Jobs <- WorkerJob{
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

	It("should output the the scraped page to the result channel", func() {
		// Arrange
		targetUrl, _ := url.Parse("https://www.google.co.uk")
		fakeScrapeResult := createFakeScrapeResult()
		expected := pages.PageFromUrl(*targetUrl)
		expected.OutLinks = createFakeScrapeResult().OutLinks
		expected.OutPages = createFakeScrapeResult().OutPages

		fakeScraper := scrapefakes.FakeScraper{}

		fakeScraper.ScrapeReturns(fakeScrapeResult, nil)

		go outHandlerFake(channels, func(result WorkerResult) { return })

		underTest := Worker{Scraper: &fakeScraper}
		go underTest.Start(channels, &queueCounter, &workCounter)

		// Act
		channels.Jobs <- WorkerJob{
			Id:  "some id",
			URL: *targetUrl,
		}

		// Assert
		outHandlerFake(channels, func(result WorkerResult) {
			Expect(result.Page).To(Equal(expected))
		})
	})

	Context("when the scrape is not successful", func() {
		It("should output the scraped pages along with the error to the write channel", func() {
			// Arrange
			targetUrl, _ := url.Parse("https://www.google.co.uk")
			expected := pages.PageFromUrl(*targetUrl)
			expected.Err = errors.New("some error")

			fakeScraper := scrapefakes.FakeScraper{}
			fakeScraper.ScrapeReturns(createFakeScrapeResult(), expected.Err)

			go outHandlerFake(channels, func(result WorkerResult) { return })

			underTest := Worker{Scraper: &fakeScraper}
			go underTest.Start(channels, &queueCounter, &workCounter)

			// Act
			channels.Jobs <- WorkerJob{
				Id:  "some id",
				URL: *targetUrl,
			}

			// Assert
			outHandlerFake(channels, func(result WorkerResult) {
				Expect(result.Page.Err).To(Equal(expected.Err))
			})
		})
	})
})

func outHandlerFake(channels Channels, callback func(page WorkerResult)) {
	for result := range channels.Results {
		callback(result)
		close(channels.Jobs)
		close(channels.Results)
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
