package md5_test

import (
	. "github.com/jjmschofield/gocrawl/internal/md5"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Md5", func() {
	Describe("HashString", func(){
		It("Should return an md5 hash of a string", func(){
			// Arrange
			input := "please hash this string, thank you!"
			expected := "601fc460b75c44d0457e8b4c00459ce3"

			// Act
			result := HashString(input)

			// Assert
			Expect(result).To(Equal(expected))
		})
	})
})
