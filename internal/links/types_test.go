package links

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("Types", func() {
	Describe("calcType", func() {

		var (
			underTest func(fromUrl url.URL, toUrl url.URL) string
		)

		BeforeEach(func() {
			underTest = calcType
		})

		Context("when a link is targeting a page", func() {
			Context("and when the urls target the same host and protocol", func() {
				It("should return the internal page type", func() {
					// Arrange
					expected := InternalPageType
					fromUrl, _ := url.Parse("https://www.bbc.co.uk/some-page")
					toUrl, _ := url.Parse("https://www.bbc.co.uk/some/other/page")

					// Act
					result := underTest(*fromUrl, *toUrl)

					// Assert
					Expect(result).To(Equal(expected))
				})
			})

			Context("or when the urls target a different host", func() {
				It("should return the external page type", func() {
					// Arrange
					expected := ExternalPageType
					fromUrl, _ := url.Parse("https://www.bbc.co.uk/some-page")
					toUrl, _ := url.Parse("https://www.bbc.com/some/other/page")

					// Act
					result := underTest(*fromUrl, *toUrl)

					// Assert
					Expect(result).To(Equal(expected))
				})
			})

			Context("or when the urls target a different protocol", func() {
				It("should return the external page type", func() {
					// Arrange
					expected := ExternalPageType
					fromUrl, _ := url.Parse("https://www.bbc.co.uk/some-page")
					toUrl, _ := url.Parse("http://www.bbc.co.uk/some/other/page")

					// Act
					result := underTest(*fromUrl, *toUrl)

					// Assert
					Expect(result).To(Equal(expected))
				})
			})
		})

		Context("when a link is targeting a file", func() {
			Context("and when the file is internal", func() {
				It("should return the internal file type", func() {
					// Arrange
					expected := InternalFileType
					fromUrl, _ := url.Parse("https://www.bbc.co.uk")
					toUrl, _ := url.Parse("https://www.bbc.co.uk/someReport.pdf")

					// Act
					result := underTest(*fromUrl, *toUrl)

					// Assert
					Expect(result).To(Equal(expected))
				})
			})

			Context("and when the file is external", func() {
				It("should return the external file type", func() {
					// Arrange
					expected := ExternalFileType
					fromUrl, _ := url.Parse("https://www.bbc.co.uk")
					toUrl, _ := url.Parse("https://www.google.co.uk/someReport.pdf")

					// Act
					result := underTest(*fromUrl, *toUrl)

					// Assert
					Expect(result).To(Equal(expected))
				})
			})
		})

		Context("when a link is targeting a tel", func() {
			It("should return the tel type", func() {
				// Arrange
				expected := TelType
				fromUrl, _ := url.Parse("https://www.bbc.co.uk")
				toUrl, _ := url.Parse("tel:+440123456789")

				// Act
				result := underTest(*fromUrl, *toUrl)

				// Assert
				Expect(result).To(Equal(expected))
			})
		})

		Context("when a link is targeting a mailto", func() {
			It("should return the mailto type", func() {
				// Arrange
				expected := MailtoType
				fromUrl, _ := url.Parse("https://www.bbc.co.uk")
				toUrl, _ := url.Parse("mailto:someone@someplace.com")

				// Act
				result := underTest(*fromUrl, *toUrl)

				// Assert
				Expect(result).To(Equal(expected))
			})
		})
	})

	Describe("isFile", func() {
		var (
			underTest func(testUrl url.URL) bool
		)

		BeforeEach(func() {
			underTest = isFile
		})

		Context("when the url has no extensions", func() {
			It("should return false", func() {
				// Arrange
				testUrl, _ := url.Parse("https://www.google.co.uk")

				// Act
				result := underTest(*testUrl)

				// Assert
				Expect(result).To(Equal(false))
			})
		})

		Context("when the url has an extensions", func() {
			Context("and when the extension is in the file extension list", func() {
					It("should return true", func() {
						for _, ext := range fileExtensions {
							// Arrange
							testUrl, _ := url.Parse("https://www.google.co.uk/somefile." + ext)

							// Act
							result := underTest(*testUrl)

							// Assert
							Expect(result).To(Equal(true))
						}
					})

					It("should return true for exe", func(){
						// Arrange
						testUrl, _ := url.Parse("https://www.google.co.uk/somefile.exe")

						// Act
						result := underTest(*testUrl)

						// Assert
						Expect(result).To(Equal(true))
					})
			})

			Context("or when the extension is not in the file extension list", func() {
				It("should return false", func() {
					// Arrange
					testUrl, _ := url.Parse("https://www.google.co.uk/some/page.html")

					// Act
					result := underTest(*testUrl)

					// Assert
					Expect(result).To(Equal(false))
				})

				It("should not panic when the ext would be at the end of the file extensions list", func() {
					// Arrange
					testUrl, _ := url.Parse("https://www.google.co.uk/somefile.zzzzzzzzzz")

					// Act
					result := underTest(*testUrl)

					// Assert
					Expect(result).To(Equal(false))
				})
			})
		})
	})
})
