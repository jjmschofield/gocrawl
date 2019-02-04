package pages_test

import (
	"errors"
	"github.com/jjmschofield/GoCrawl/internal/app/links"
	. "github.com/jjmschofield/GoCrawl/internal/app/pages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/url"
)

var _ = Describe("Page", func() {
	Describe("PageFromUrl", func(){
		var (
			underTest func(srcUrl url.URL) Page
		)

		BeforeEach(func() {
			underTest = PageFromUrl
		})

		It("should create a page from a url", func(){
			// Arrange
			expectedUrl, _ := url.Parse("https://www.google.co.uk")
			expectedId := "3f1437859f73b447885255a95afa99a1"

			// Act
			result := underTest(*expectedUrl)

			// Assert
			Expect(result.Id).To(Equal(expectedId))
			Expect(result.URL).To(Equal(*expectedUrl))
		})
	})

	Describe("MarshalJson", func(){
		It("should marshal to json with sensible defaults", func(){
			// Arrange
			srcUrl, _ := url.Parse("https:///www.google.co.uk")
			underTest := PageFromUrl(*srcUrl)
			expected := []byte(`{"id":"e23861440e88ba6d0510254bdd3fe614","url":"https:///www.google.co.uk","outPages":{"internal":[]},"outLinks":{"internal":[],"internalFiles":[],"external":[],"externalFiles":[],"tel":[],"mailto":[],"unknown":[]},"error":null}`)

			// Act
			result, _ := underTest.MarshalJSON()

			// Assert
			Expect(result).To(Equal(expected))
		})
	})

	Describe("Print", func(){
		BeforeEach(func(){
			log.SetOutput(ioutil.Discard)
		})

		It("should print information to stdout without panicing", func(){
			// Arrange
			srcUrl, _ := url.Parse("https:///www.google.co.uk")
			underTest := PageFromUrl(*srcUrl)
			underTest.OutLinks = links.ToLinkGroup(links.FromHrefs(*srcUrl, []string{"https:///www.google.co.uk/internal"}))
			underTest.Err = errors.New("some error")

			// Act
			underTest.Print()
		})
	})
})
