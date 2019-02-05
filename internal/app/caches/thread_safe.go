package caches

type ThreadSafeCache interface {
	Add(str string)
	Remove(str string)
	Has(str string) bool
	Count() int
}
