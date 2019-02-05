package links_test

import (
	. "github.com/jjmschofield/GoCrawl/internal/app/links"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("Link", func() {
	Describe("NewAbsLink", func(){
		var (
			underTest func(fromUrl url.URL, toUrl url.URL) Link
		)

		BeforeEach(func(){
			underTest = NewAbsLink
		})

		It("should calculate an id", func(){
			// Arrange
			fromUrl, _ := url.Parse("https://www.google.co.uk")
			toUrl, _ := url.Parse("https://www.google.co.uk/some/path")
			expected := "a71e371fbedd579f156709922c7c0070"

			// Act
			link := underTest(*fromUrl, *toUrl)

			// Asset
			Expect(link.Id).To(Equal(expected))
		})

		It("should assign the toUrl", func(){
			// Arrange
			expected := "https://www.google.co.uk/some/path"
			fromUrl, _ := url.Parse("https://www.google.co.uk")
			toUrl, _ := url.Parse(expected)

			// Act
			link := underTest(*fromUrl, *toUrl)

			// Asset
			Expect(link.ToURL).To(Equal(expected))
		})

		It("should assign the fromUrl", func(){
			// Arrange
			expected := "https://www.google.co.uk"
			fromUrl, _ := url.Parse(expected)
			toUrl, _ := url.Parse("https://www.google.co.uk/some/path")

			// Act
			link := underTest(*fromUrl, *toUrl)

			// Asset
			Expect(link.FromURL).To(Equal(expected))
		})

		It("should calculate a link type", func(){
			// Arrange
			fromUrl, _ := url.Parse("https://www.google.co.uk")
			toUrl, _ := url.Parse("https://www.google.co.uk/some/path")
			expected := InternalPageType

			// Act
			link := underTest(*fromUrl, *toUrl)

			// Asset
			Expect(link.Type).To(Equal(expected))
		})

		It("should resolve relative toUrl's to be absolute urls with the base of the from url", func(){
			// Arrange
			fromUrl, _ := url.Parse("https://www.google.co.uk/some/src/path")
			toUrl, _ := url.Parse("../dest/path")
			expected := "https://www.google.co.uk/some/dest/path"

			// Act
			link := underTest(*fromUrl, *toUrl)

			// Asset
			Expect(link.ToURL).To(Equal(expected))
		})
	})

	Describe("FromHref", func(){
		var (
			underTest func(pageUrl url.URL, href string) (link Link, err error)
		)

		BeforeEach(func(){
			underTest = FromHref
		})

		Context("when the provided href is a valid url", func(){
			It("should create a new absolute link", func(){
				// Arrange
				fromUrl, _ := url.Parse("https://www.google.co.uk")
				href := "https://www.google.co.uk/some/path"
				toUrl, _ := url.Parse(href)
				expected := NewAbsLink(*fromUrl, *toUrl)

				// Act
				link, _ := underTest(*fromUrl, "https://www.google.co.uk/some/path")

				// Assert
				Expect(link).To(Equal(expected))
			})
		})

		Context("or when the provided href is not valid url", func(){
			It("should return the error", func(){
				// Arrange
				fromUrl, _ := url.Parse("https://www.google.co.uk")

				// Act
				_, err := underTest(*fromUrl, ":")

				// Assert
				Expect(err.Error()).To(Equal("parse :: missing protocol scheme"))
			})
		})
	})
})
