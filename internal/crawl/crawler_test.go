package crawl_test

import (
	"github.com/jjmschofield/gocrawl/internal/caches"
	"github.com/jjmschofield/gocrawl/internal/counters"
	. "github.com/jjmschofield/gocrawl/internal/crawl"
	"github.com/jjmschofield/gocrawl/internal/crawl/crawlfakes"
	"github.com/jjmschofield/gocrawl/internal/pages"
	"github.com/jjmschofield/gocrawl/internal/scrape"
	"github.com/jjmschofield/gocrawl/internal/writers"
	"github.com/jjmschofield/gocrawl/internal/writers/writersfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/url"
	"sync"
)

var _ = Describe("PageCrawler", func() {

	Describe("Crawl", func() {
		var startUrl *url.URL

		BeforeEach(func() {
			log.SetOutput(ioutil.Discard)
			startUrl, _ = url.Parse("https://www.google.co.uk")
		})

		It("should start the required numbers of workers", func() {
			// Arrange
			expected := 3
			config, workerFake, _ := createConfig()
			config.WorkerCount = expected
			underTest := NewPageCrawler(config)

			// Act
			underTest.Crawl(*startUrl)

			// Assert
			Expect(workerFake.StartCallCount()).To(Equal(expected))
		})

		It("should enqueue a crawl when provided with a url", func(done Done) {
			// Arrange
			expected := startUrl.String()

			config, workerFake, _ := createConfig()
			workerFake.StartStub = createWorkerStartStub(WorkerResult{}, func(job WorkerJob) {
				// Assert
				Expect(job.URL.String()).To(Equal(expected))
				done <- nil
			})

			underTest := NewPageCrawler(config)

			// Act
			underTest.Crawl(*startUrl)
		})

		It("should wait until all discovered pages have been crawled before returning accurate counters", func() {
			// Arrange
			expectedProcessing := 0
			expectedCrawled := 2
			expectedCrawling := 0

			config, workerFake, _ := createConfig()

			firstResult := WorkerResult{
				Result: scrape.Result{
					OutPages: pages.PageGroup{
						Internal: map[string]string{"some id": "https://www.google.co.uk"},
					},
				}}

			workerFake.StartStub = createWorkerStartStub(firstResult, nil)

			underTest := NewPageCrawler(config)

			// Act
			result := underTest.Crawl(*startUrl)

			// Assert
			Expect(result.Crawling.Count()).To(Equal(int64(expectedCrawling)))
			Expect(result.CrawlComplete.Count()).To(Equal(int64(expectedCrawled)))
			Expect(result.Processing.Count()).To(Equal(int64(expectedProcessing)))
		})

		Describe("When a scrape result is returned", func() {
			It("should enqueue any new internal pages", func() {
				// Arrange
				expectedId := "some new id"
				expectedUrl, _ := url.Parse("https://www.google.co.uk/a")

				config, workerFake, _ := createConfig()

				firstResult := WorkerResult{
					Result: scrape.Result{
						OutPages: pages.PageGroup{
							Internal: map[string]string{expectedId: expectedUrl.String()},
						},
					}}

				i := 0
				workerFake.StartStub = createWorkerStartStub(firstResult, func(job WorkerJob) {
					if i > 0 {
						// Assert
						Expect(job.Id).To(Equal(expectedId))
						Expect(job.URL).To(Equal(*expectedUrl))
					}
					i++
				})

				underTest := NewPageCrawler(config)

				// Act
				underTest.Crawl(*startUrl)
			})
		})

		Describe("when enqueuing a page", func() {
			It("should not enqueue a page which has been crawled", func() {
				// Arrange
				config, workerFake, _ := createConfig()
				alreadyCrawledUrl, _ := url.Parse("https://www.google.co.uk/a")

				id, _ := pages.CalcPageId(*alreadyCrawledUrl)

				config.Caches.Crawled.Add(id)

				workerFake.StartStub = createWorkerStartStub(WorkerResult{
					CrawledId: id,
					Result: scrape.Result{
						OutPages: pages.PageGroup{
							Internal: map[string]string{id: alreadyCrawledUrl.String()},
						},
					},
				}, nil)

				underTest := NewPageCrawler(config)

				// Act
				counters := underTest.Crawl(*startUrl)

				// Assert
				Expect(counters.CrawlComplete.Count()).To(Equal(int64(1)))
			})

			It("should not enqueue a page which is currently being processed", func() {
				config, workerFake, _ := createConfig()
				inProcessUrl, _ := url.Parse("https://www.google.co.uk/a")

				id, _ := pages.CalcPageId(*inProcessUrl)

				config.Caches.Processing.Add(id)

				workerFake.StartStub = createWorkerStartStub(WorkerResult{
					CrawledId: "3f1437859f73b447885255a95afa99a1",
					Result: scrape.Result{
						OutPages: pages.PageGroup{
							Internal: map[string]string{id: inProcessUrl.String()},
						},
					},
				}, nil)

				underTest := NewPageCrawler(config)

				counters := underTest.Crawl(*startUrl)

				Expect(counters.CrawlComplete.Count()).To(Equal(int64(1)))
			})
		})
	})
})

var _ = Describe("NewPageCrawler", func() {
	It("should construct a page crawler using the provided config", func() {
		// Arrange
		crawledCache := caches.NewStrThreadSafe()
		processingCache := caches.NewStrThreadSafe()

		expected := Config{
			Caches: Caches{
				Crawled:    &crawledCache,
				Processing: &processingCache,
			},
			Worker:      &Worker{Scraper: scrape.PageScraper{}},
			WorkerCount: 10,
			Writer:      &writers.FileWriter{FilePath: "some/file/path"},
		}

		// Act
		result := NewPageCrawler(expected)

		// Assert
		Expect(result.Config).To(Equal(expected))
	})
})

var _ = Describe("NewDefaultPageCrawler", func() {
	It("should construct a page crawler using the default config", func() {
		// Arrange
		const expectedWorkerCount = 10
		const expectedOutputPath = "some/file/path"

		crawledCache := caches.NewStrThreadSafe()
		processingCache := caches.NewStrThreadSafe()

		expected := Config{
			Caches: Caches{
				Crawled:    &crawledCache,
				Processing: &processingCache,
			},
			Worker:      &Worker{Scraper: scrape.PageScraper{}},
			WorkerCount: expectedWorkerCount,
			Writer:      &writers.FileWriter{FilePath: expectedOutputPath},
		}

		// Act
		result := NewDefaultPageCrawler(expectedWorkerCount, expectedOutputPath)

		// Assert
		Expect(result.Config).To(Equal(expected))
	})
})

func createConfig() (config Config, workerFake *crawlfakes.FakeQueueWorker, writerFake *writersfakes.FakeWriter) {
	crawledCache := caches.NewStrThreadSafe()
	processingCache := caches.NewStrThreadSafe()

	workerFake = createFakeWorker()
	fakeWriter := &writersfakes.FakeWriter{}

	config = Config{
		Caches: Caches{
			Crawled:    &crawledCache,
			Processing: &processingCache,
		},
		Worker:      workerFake,
		WorkerCount: 1,
		Writer:      fakeWriter,
	}

	return config, workerFake, fakeWriter
}

func createFakeWorker() *crawlfakes.FakeQueueWorker {
	worker := crawlfakes.FakeQueueWorker{}
	worker.StartStub = createWorkerStartStub(WorkerResult{}, nil)
	return &worker
}

func createWorkerStartStub(fake WorkerResult, callback func(job WorkerJob)) func(chans WorkerChannels, queueCounter *counters.AtomicInt64, workCounter *counters.AtomicInt64, wg *sync.WaitGroup) {
	return func(chans WorkerChannels, queueCounter *counters.AtomicInt64, workCounter *counters.AtomicInt64, wg *sync.WaitGroup) {
		defer GinkgoRecover()
		defer wg.Done()

		for job := range chans.In {
			workCounter.Add(1)
			queueCounter.Sub(1)

			if callback != nil {
				callback(job)
			}

			chans.Out <- fake
		}
	}
}
