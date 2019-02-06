package crawl_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCrawl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crawl Suite")
}
