package pages_test

import (
	. "github.com/jjmschofield/gocrawl/internal/pages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("Id", func() {
	Describe("CalcPageId", func() {
		var (
			underTest func(srcUrl url.URL) (id string, normalizedUrl url.URL)
		)

		BeforeEach(func() {
			underTest = CalcPageId
		})

		It("should return a deterministic hash of a URL as the page id", func() {
			// Arrange
			srcUrl, _ := url.Parse("https://www.google.co.uk")
			exptected := "3f1437859f73b447885255a95afa99a1"

			// Act
			result, _ := underTest(*srcUrl)

			// Result
			Expect(result).To(Equal(exptected))
		})

		It("should remove trailing slashes when calculating a page id and return the normalized url", func() {
			// Arrange
			srcUrl, _ := url.Parse("https://www.google.co.uk/")
			expectedId := "3f1437859f73b447885255a95afa99a1"
			expectedUrl, _ := url.Parse("https://www.google.co.uk")

			// Act
			id, normalizedUrl := underTest(*srcUrl)

			// Result
			Expect(id).To(Equal(expectedId))
			Expect(normalizedUrl).To(Equal(*expectedUrl))

		})

		It("should remove fragments calculating a page id and return the normalized url", func() {
			// Arrange
			srcUrl, _ := url.Parse("https://www.google.co.uk/#some-frag")
			expectedId := "3f1437859f73b447885255a95afa99a1"
			expectedUrl, _ := url.Parse("https://www.google.co.uk")

			// Act
			id, normalizedUrl := underTest(*srcUrl)

			// Result
			Expect(id).To(Equal(expectedId))
			Expect(normalizedUrl).To(Equal(*expectedUrl))
		})
	})
})
