package counters_test

import (
	. "github.com/jjmschofield/GoCrawl/internal/counters"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("atomicInt64", func() {
	var (
		underTest  AtomicInt64
	)

	BeforeEach(func() {
		underTest = AtomicInt64{}
	})

	Describe("Add", func() {
		Context("With a count of zero", func() {
			It("should add the provided delta and return the new count", func() {
				// Arrange
				delta := int64(100)
				expected := delta
				Expect(underTest.Count()).To(Equal(int64(0)))

				// Act
				result := underTest.Add(delta)

				// Assert
				Expect(result).To(Equal(expected))
				Expect(underTest.Count()).To(Equal(expected))
			})
		})

		Context("With a non-zero count", func() {
			It("should add the provided delta and return the new count", func() {
				// Arrange
				starting := int64(50)
				delta := int64(100)
				expected := starting + delta
				underTest.Add(starting)
				Expect(underTest.Count()).To(Equal(starting))

				// Act
				result := underTest.Add(delta)

				// Assert
				Expect(result).To(Equal(expected))
				Expect(underTest.Count()).To(Equal(expected))
			})
		})

		Context("With a peak lower then the new count", func() {
			It("it should update peak to match the new count", func() {
				// Arrange
				delta := int64(100)
				expected := delta
				Expect(underTest.Peak()).To(Equal(int64(0)))

				// Act
				underTest.Add(delta)

				// Assert
				Expect(underTest.Peak()).To(Equal(expected))
			})
		})

		Context("With a peak higher then the new count", func() {
			It("it should retain the old peak", func() {
				// Arrange
				starting := int64(50)
				delta := int64(1)

				underTest.Add(starting)
				underTest.Sub(starting)
				Expect(underTest.Peak()).To(Equal(starting))

				// Act
				underTest.Add(delta)

				// Assert
				Expect(underTest.Peak()).To(Equal(starting))
			})
		})
	})

	Describe("Sub", func() {
		Context("With any count value" , func() {
			It("should subtract the provided delta and return the new count", func() {
				// Arrange
				delta := int64(100)
				expected := -delta
				Expect(underTest.Count()).To(Equal(int64(0)))

				// Act
				result := underTest.Sub(delta)

				// Assert
				Expect(result).To(Equal(expected))
				Expect(underTest.Count()).To(Equal(expected))
			})
		})
	})

	Describe("Count", func() {
		Context("With any count value" , func() {
			It("should return the current count value", func() {
				// Arrange
				delta := int64(100)
				expected := delta

				// Act
				underTest.Add(delta)

				// Assert
				Expect(underTest.Count()).To(Equal(expected))
			})
		})
	})

	Describe("Peak", func() {
		Context("With any peak value" , func() {
			It("should return the current peak value", func() {
				// Arrange
				delta := int64(100)
				expected := delta

				// Act
				underTest.Add(delta)
				underTest.Sub(delta)

				// Assert
				Expect(underTest.Peak()).To(Equal(expected))
			})
		})
	})
})