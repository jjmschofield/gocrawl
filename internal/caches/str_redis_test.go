package caches_test

import (
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	. "github.com/jjmschofield/gocrawl/internal/caches"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StrRedis", func() {

	const setKey = "someKey"

	var (
		underTest StrRedis
		redisMock *miniredis.Miniredis
	)

	BeforeEach(func() {
		var err error
		redisMock, err = miniredis.Run()

		if err != nil {
			panic(err)
		}

		options := &redis.Options{
			Addr:     redisMock.Addr(),
			Password: "", // no password set
			DB:       0,  // use default DB
		}

		underTest = NewStrRedis(setKey, options)
	})

	AfterEach(func() {
		redisMock.Close()
	})

	It("should connect to redis", func() {
		Expect(redisMock.CurrentConnectionCount()).To(Equal(1))
	})

	Describe("Add", func(){
		It("Should add a string to the cache", func(){
			// Arrange
			expectedStr := "some string"
			Expect(redisMock.Exists(setKey)).To(Equal(false))

			// Act
			underTest.Add(expectedStr)

			// Assert
			Expect(redisMock.IsMember(setKey, expectedStr)).To(Equal(true))
		})
	})
	Describe("Remove", func() {
		Context("When the provided string is in the cache", func() {
			It("Should remove the string from the cache", func() {
				// Arrange
				expectedStr := "some string"
				_, err := redisMock.SetAdd(setKey, expectedStr, "some other string")
				if err != nil {
					panic(err)
				}

				Expect(redisMock.IsMember(setKey, expectedStr)).To(Equal(true))

				// Act
				underTest.Remove(expectedStr)

				// Assert
				Expect(redisMock.IsMember(setKey, expectedStr)).To(Equal(false))
			})
		})

		Context("When the provided string is not in the cache", func(){
			It("Should not panic", func(){
				// Arrange
				expectedStr := "some string"
				_, err := redisMock.SetAdd(setKey, "some other string")
				if err != nil {
					panic(err)
				}

				Expect(redisMock.IsMember(setKey, expectedStr)).To(Equal(false))

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
				_, err := redisMock.SetAdd(setKey, expectedStr)
				if err != nil {
					panic(err)
				}

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
				_, err := redisMock.SetAdd(setKey, string(i))
				if err != nil {
					panic(err)
				}
			}

			// Assert
			Expect(underTest.Count()).To(Equal(3))
		})
	})
})
