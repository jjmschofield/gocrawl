package main_test

import (
	. "github.com/jjmschofield/gocrawl"
	"github.com/jjmschofield/gocrawl/internal/crawl"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("crawl", func() {
	Measure("benchmark over 3 runs", func(b Benchmarker) {
		var counters crawl.Counters

		b.Time("runtime", func() {
			crawlUrl, _ := url.Parse("https://www.monzo.com")
			counters = Crawl(*crawlUrl, 100, "./data", "")
		})

		Expect(counters.Discovered.Count()).To(BeNumerically(">", 1))

		b.RecordValue("pages discovered", float64(counters.Discovered.Count()))
		b.RecordValue("pages crawled", float64(counters.CrawlComplete.Count()))
		b.RecordValue("queue peak", float64(counters.CrawlsQueued.Peak()))
		b.RecordValue("parallel scrape peak", float64(counters.Crawling.Peak()))
	}, 3)
})
