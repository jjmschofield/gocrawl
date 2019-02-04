package pages_test

import (
	. "github.com/onsi/ginkgo"
	"net/url"

	. "github.com/onsi/gomega"
	//
	. "github.com/jjmschofield/GoCrawl/internal/app/pages"
)

var _ = Describe("Group", func() {
	Describe("ToLinkGroup", func(){
		var (
			underTest func(internal []Page) (group PageGroup)
		)

		BeforeEach(func(){
			underTest = ToPageGroup
		})

		It("should parse internal pages into a the internal slice", func(){
			// Arrange
			srcUrl, _ := url.Parse("https://www.google.co.uk")
			expected := PageFromUrl(*srcUrl)
			pages := []Page{expected}

			// Act
			result := underTest(pages)

			// Assert
			Expect(result.Internal[0]).To(Equal(expected))
		})
	})

	Describe("MarshalJson", func(){
		It("should marshal links into ids", func(){
			// Arrange
			srcUrl, _ := url.Parse("https://www.google.co.uk")
			pages := []Page{PageFromUrl(*srcUrl)}
			expected := []byte(`{"internal":["3f1437859f73b447885255a95afa99a1"]}`)
			underTest := ToPageGroup(pages)

			// Act
			result, _ := underTest.MarshalJSON()

			// Assert
			Expect(result).To(Equal(expected))
		})
	})
})
