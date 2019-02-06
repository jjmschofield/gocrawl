package scrape_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ScrapeFetch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scrape Suite")
}
