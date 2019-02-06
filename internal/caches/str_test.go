package caches_test

import (
	. "github.com/jjmschofield/gocrawl/internal/caches"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StrThreadSafe", func() {
	var (
		underTest StrThreadSafe
	)

	BeforeEach(func() {
		underTest = NewStrThreadSafe()
	})

	Describe("Add", func(){
		It("Should add a string to the cache", func(){
			// Arrange
			expectedStr := "some string"
			Expect(underTest.Has(expectedStr)).To(Equal(false))

			// Act
			underTest.Add(expectedStr)

			// Assert
			Expect(underTest.Has(expectedStr)).To(Equal(true))
		})
	})

	Describe("Remove", func(){
		Context("When the provided string is in the cache", func(){
			It("Should remove the string from the cache", func(){
				// Arrange
				expectedStr := "some string"
				underTest.Add("some string")
				Expect(underTest.Has(expectedStr)).To(Equal(true))

				// Act
				underTest.Remove(expectedStr)

				// Assert
				Expect(underTest.Has(expectedStr)).To(Equal(false))
			})
		})

		Context("When the provided string is not in the cache", func(){
			It("Should not panic", func(){
				// Arrange
				expectedStr := "some string"
				Expect(underTest.Has(expectedStr)).To(Equal(false))

				// Act
				underTest.Remove(expectedStr)
			})
		})
	})

	Describe("Has", func(){
		Context("When the provided string is in the cache", func(){
			It("Should return true", func(){
				// Arrange
				expectedStr := "some string"
				underTest.Add(expectedStr)

				// Act
				result := underTest.Has(expectedStr)

				// Assert
				Expect(result).To(Equal(true))
			})
		})

		Context("When the provided string is in the cache", func(){
			It("Should return false", func(){
				// Arrange
				expectedStr := "some string"

				// Act
				result := underTest.Has(expectedStr)

				// Assert
				Expect(result).To(Equal(false))
			})
		})
	})

	Describe("Count", func(){
		It("Should return the number of items in the cache", func(){
			// Arrange
			expected := 3
			Expect(underTest.Count()).To(Equal(0))

			// Act
			for i := 0; i < expected; i++{
				underTest.Add(string(i))
			}

			// Assert
			Expect(underTest.Count()).To(Equal(3))
		})
	})
})
