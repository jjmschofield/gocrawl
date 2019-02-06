package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCrawl(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping benchmarking tests")
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Crawl Suite")
}
